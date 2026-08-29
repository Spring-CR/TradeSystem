package order_domain

import (
	"log"
	"rhino-common/domain_error"
	"rhino-core/adapter_registry"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-trade-channel/channel"
)

type OrderEngine struct {
	applicationCfg    *domain_cfg.ApplicationCfg
	orderAcceptor     OrderAcceptorInterface
	orderExecutor     OrderExecutorInterface
	orderOrchestrator *OrderOrchestrator
	tradeRespCh       chan *types.TradeActionRespReturn
	onTradeResps      []func(*types.TradeActionRespReturn) bool
}

func NewOrderEngine(applicationCfg *domain_cfg.ApplicationCfg, onTradeResps []func(*types.TradeActionRespReturn) bool) (engine *OrderEngine, de *domain_error.Error) {

	log.Printf("start to NewOrderEngine, propograteTradeResp:%v\n", onTradeResps != nil)

	var tradeRespCh chan *types.TradeActionRespReturn
	if onTradeResps != nil {
		tradeRespCh = make(chan *types.TradeActionRespReturn, 10240)
	}

	// 前移
	engine = &OrderEngine{applicationCfg: applicationCfg, tradeRespCh: tradeRespCh, onTradeResps: onTradeResps}
	if len(onTradeResps) > 0 {
		engine.propograteTraceableTradeActionResp()
		log.Println("start to propograteTradeResp...")
	}

	channelMap := map[string]channel.TradeChannelInterface{}
	orderOrchestrator := NewOrderOrchestrator(applicationCfg, tradeRespCh)

	// create trade channel config
	tradeChannelDetails := applicationCfg.GetTradeChannels()
	tradeChannelCfgMap := make(map[string]*domain_cfg.TradeChannelCfg)
	for _, tradeChannelDetails := range tradeChannelDetails {
		log.Printf("start to create channel: %s\n", tradeChannelDetails.TradeChannel.ChannelCode)
		var tradeChannelCfg *domain_cfg.TradeChannelCfg
		tradeChannelCfg, de = domain_cfg.NewTradeChannelCfg(applicationCfg, tradeChannelDetails)
		if de != nil {
			return
		}

		_channel, de, err := adapter_registry.CallAdapterFunction(tradeChannelDetails.TradeChannel.AdapterPath, tradeChannelCfg, func(tradeActionResp *schema.TradeActionResp) {
			orderOrchestrator.updateOrderStatus(tradeActionResp)
		})
		if err != nil {
			panic(err)
		}

		if de != domain_error.NilDomainError {
			log.Printf("de:%+v\n", de)
			return nil, de.(*domain_error.Error)
		}

		tradeChannel := _channel.(channel.TradeChannelInterface)

		channelMap[tradeChannelDetails.TradeChannel.ChannelCode] = tradeChannel
		tradeChannelCfgMap[tradeChannelDetails.TradeChannel.ChannelCode] = tradeChannelCfg

		// 在这里才加入监听器太晚了，会导致回报不被处理
		// log.Printf("AddTradeActionRespListener for %s\n", tradeChannelDetails.TradeChannel.ChannelCode)
		// tradeChannel.AddTradeActionRespListener(func(tradeActionResp *schema.TradeActionResp) {
		//	orderOrchestrator.updateOrderStatus(tradeActionResp)
		// })
	}

	exector := NewOrderExecutor(applicationCfg, channelMap, orderOrchestrator)
	acceptor := NewOrderAcceptor(128, exector, orderOrchestrator)

	//engine = &OrderEngine{applicationCfg: applicationCfg, orderAcceptor: acceptor, orderExecutor: exector, orderOrchestrator: orderOrchestrator, tradeRespCh: tradeRespCh, onTradeResps: onTradeResps}

	// if len(onTradeResps) > 0 {
	// 	engine.propograteTraceableTradeActionResp()
	// 	log.Println("start to propograteTradeResp...")
	// }

	engine.orderAcceptor = acceptor
	engine.orderExecutor = exector
	engine.orderOrchestrator = orderOrchestrator

	log.Println("new order engine!")

	// 跟踪收市通知事件
	exector.trackingMarketCloseEvent(tradeChannelCfgMap)

	return
}

func (e *OrderEngine) GetOrderAcceptor() OrderAcceptorInterface {
	return e.orderAcceptor
}

func (e *OrderEngine) propograteTraceableTradeActionResp() {
	if e.tradeRespCh == nil {
		return
	}
	go func() {
		for {
			tradeResp := <-e.tradeRespCh
			for _, onTradeResp := range e.onTradeResps {
				ok := onTradeResp(tradeResp)
				if !ok {
					break
				}
			}
			if tradeResp.WaitGroup != nil {
				log.Println("WaitGroup.Done() in propograteTraceableTradeActionResp")
				tradeResp.WaitGroup.Done()
			}
		}
	}()
}

func (e *OrderEngine) GetSystemAndBusinessCodes() (string, string) {
	return e.applicationCfg.GetSystemAndBusinessCodes()
}

func (e *OrderEngine) GetOrderByAppOrdID(appOrdID string) (order *schema.TradeOrder, ok bool) {
	return e.orderOrchestrator.GetOrderByAppOrdID(appOrdID)
}

func (e *OrderEngine) GetTraceableOrderByAppOrdID(appOrdID string) (order *types.TraceableTradeOrder, ok bool) {
	return e.orderOrchestrator.GetTraceableOrderByAppOrdID(appOrdID)
}

func (e *OrderEngine) ForceArchiving() {
	e.orderOrchestrator.ForceArchiving()
}

func (e *OrderEngine) ForcePurging() {
	e.orderOrchestrator.ForcePurging()
}

func (e *OrderEngine) GetOrderOrchestrator() *OrderOrchestrator {
	return e.orderOrchestrator
}

func (e *OrderEngine) Dump() {
	e.orderOrchestrator.GetOrderCache().Dump()
}

func (e *OrderEngine) GetApiToken() string {
	return e.applicationCfg.GetApiToken()
}
