package api

import (
	"log"
	"rhino-common/domain_error"
	"rhino-core/adapter_registry"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain"
	"rhino-core/types"
	"rhino-engine/api/api_adapter"
	"rhino-engine/api/fix"
	"rhino-engine/api/fix_api_adapter"
	"rhino-engine/api/http"
	"rhino-engine/api/stream"
	"rhino-engine/api/stream/kafka"
	"time"
)

type TradeEngineApiServer struct {
	applicationCfg *domain_cfg.ApplicationCfg
	streamAPI      *stream.StreamAPI
	httAPI         *http.HttpApiServer
	fixAPI         *fix.FixApi
}

func NewTradeEngineApiServer(applicationCfg *domain_cfg.ApplicationCfg, workerCount int, workerJobBuf int) (*TradeEngineApiServer, *domain_error.Error) {
	inst := &TradeEngineApiServer{applicationCfg: applicationCfg}

	adapterPath := applicationCfg.GetApiAdapterPath()
	var apiAdapter api_adapter.APIAdapter
	if adapterPath != "" {

		// 从注册表获取适配器的构造函数（目前，对于apiAdpater，是无参数函数，有参函数需要根据特殊情况来处理了）
		_apiAdater, de, err := adapter_registry.CallAdapterFunction(adapterPath, applicationCfg)
		if err != nil {
			panic(err)
		}

		if de != domain_error.NilDomainError {
			log.Printf("de:%+v\n", de)
			return nil, de.(*domain_error.Error)
		}

		// 获取apiAdapter
		apiAdapter = _apiAdater.(api_adapter.APIAdapter)
		log.Println("finish get apiAdapter")
	}

	kafkaBrokers := applicationCfg.GetApiKafkaBrokers()
	var onTradeResps []func(tradeResp *types.TradeActionRespReturn) bool
	if kafkaBrokers != "" {
		client := kafka.NewKafkaClient(applicationCfg)
		inst.streamAPI = stream.NewStreamAPI(apiAdapter, client)
		onTradeResps = append(onTradeResps, inst.streamAPI.OnTradeResp)
	}

	fixServerAdapterPath := applicationCfg.GetFixServerAdapterPath()
	var fixApiAdapter fix_api_adapter.FixApiAdapter
	if fixServerAdapterPath != "" {
		_fixApiAdapter, de, err := adapter_registry.CallAdapterFunction(fixServerAdapterPath, applicationCfg)
		if err != nil {
			panic(err)
		}

		if de != domain_error.NilDomainError {
			domain_error.ProcessSevereError(true, 5, de.(*domain_error.Error), nil, "fail to create fixApiAdapter")
			return nil, de.(*domain_error.Error)
		}

		// 获取apiAdapter
		fixApiAdapter = _fixApiAdapter.(fix_api_adapter.FixApiAdapter)
		log.Println("finish get fixApiAdapter")
		inst.fixAPI = fix.NewFixApi(fixApiAdapter, apiAdapter)
		onTradeResps = append(onTradeResps, inst.fixAPI.OnTradeResp)
	}

	// 创建engine
	orderEngine, de := order_domain.NewOrderEngine(applicationCfg, onTradeResps)
	if de != nil {
		return nil, de
	}
	log.Println("finish create orderEngine!")

	if inst.streamAPI != nil {
		// 创建kafka客户端
		begin := time.Now()
		log.Printf("=====> finish create KafkaClient, time cost:%v\n", time.Since(begin))
		inst.streamAPI.Start(orderEngine, 20000, workerCount, workerJobBuf)
	}
	if inst.fixAPI != nil {
		log.Printf("=====> finish create fixserver\n")
		inst.fixAPI.Start(orderEngine)
	}

	inst.httAPI = http.NewHttpApiServer(applicationCfg.GetHttpAPIPort(), apiAdapter, orderEngine, workerCount)

	log.Println("finish NewTradeEngineApiServer!")

	return inst, nil
}

// 启用http API
func (s *TradeEngineApiServer) Start() error {
	return s.httAPI.Start()
}
