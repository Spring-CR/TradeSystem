package stream

import (
	"log"
	"rhino-common/domain_error"
	"rhino-core/order_domain"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-engine/api/api_adapter"
)

type StreamAPI struct {
	apiAdapter       api_adapter.APIAdapter
	engine           *order_domain.OrderEngine
	client           StreamClient
	ingressMsgBuffer int
	workers          []*streamAPIWorker // 工作协程数
	respKeys         map[string]bool    // 已经发送过的常规成交回报的回报key
	respReqMsgSeq    *respReqMsgSeq     // 已经发送过的异常（非常规，由domain_error.Error转化）成交回报所关联的请求消息的msgSeq，只有对streamApi有意义
}

type respReqMsgSeq struct {
	respReqMsgSeqs map[int64]bool
	minSeq         int64
	maxSeq         int64
}

func NewStreamAPI(apiAdapter api_adapter.APIAdapter, client StreamClient) *StreamAPI {
	inst := &StreamAPI{apiAdapter: apiAdapter, client: client, respReqMsgSeq: &respReqMsgSeq{}}
	inst.onReset()
	return inst
}

func NewStreamAPI1(apiAdapter api_adapter.APIAdapter, engine *order_domain.OrderEngine, client StreamClient, ingressMsgBuffer int, workerCount int, workerJobBuf int) *StreamAPI {
	inst := &StreamAPI{apiAdapter: apiAdapter, engine: engine, client: client, ingressMsgBuffer: ingressMsgBuffer}
	// 确保 size 是 2 的幂
	if workerCount&(workerCount-1) != 0 {
		panic("workerCount must be a power of 2")
	}

	for i := 0; i < workerCount; i++ {
		inst.workers = append(inst.workers, newStreamApIWorker(apiAdapter, engine, client, workerJobBuf))
	}

	inst.respReqMsgSeq = &respReqMsgSeq{}

	inst.onReset()

	engine.GetOrderOrchestrator().GetOrderCache().SetAfterResetFunc(inst.onReset)

	return inst
}

func (s *StreamAPI) GetOrderByAppOrdID(appOrdID string) (*schema.TradeOrder, bool) {
	return s.engine.GetOrderByAppOrdID(appOrdID)
}

// var (
// 	shardingDuration int64
// )

//	func GetShardingDuration() int64 {
//		return shardingDuration
//	}
func (s *StreamAPI) Start(engine *order_domain.OrderEngine, ingressMsgBuffer int, workerCount int, workerJobBuf int) {
	//inst := &StreamAPI{apiAdapter: apiAdapter, engine: engine, client: client, ingressMsgBuffer: ingressMsgBuffer}
	s.engine = engine
	s.ingressMsgBuffer = ingressMsgBuffer
	
	// 确保 size 是 2 的幂
	if workerCount&(workerCount-1) != 0 {
		panic("workerCount must be a power of 2")
	}

	for i := 0; i < workerCount; i++ {
		s.workers = append(s.workers, newStreamApIWorker(s.apiAdapter, engine, s.client, workerJobBuf))
	}

	//s.respReqMsgSeq = &respReqMsgSeq{}
	//s.onReset()

	engine.GetOrderOrchestrator().GetOrderCache().SetAfterResetFunc(s.onReset)
	
	ch := s.client.PrepareIngressMessageChannel(s.ingressMsgBuffer)

	//workerCount := len(s.workers)
	n := workerCount - 1

	workSharding := s.apiAdapter.GetWorkerSharding(workerCount, s.GetOrderByAppOrdID)
	if workSharding == nil {
		workSharding = func(rawMsg []byte, decodeOrder *schema.TradeOrder, cumTaskCount int) (workerIndex int) {
			return cumTaskCount & n // 使用按位与选择 worker
		}
	}

	go func() {
		taskCount := 0 // 记录任务数量
		workerIndex := 0

		if len(s.respReqMsgSeq.respReqMsgSeqs) > 0 {
		for {
			msg := <-ch
			if s.respReqMsgSeq.respReqMsgSeqs[msg.MsgSeq] {
				if msg.MsgSeq >= s.respReqMsgSeq.maxSeq {
					log.Printf("break pre-looping for msgSeq1 %d\n", msg.MsgSeq)
					break
				} else {
					continue
				}
			}

			//workerIndex := taskCount & n // 使用按位与选择 worker
			//begin := time.Now()
			workerIndex = workSharding(msg.Data, nil, taskCount)
			//time.Sleep(10*time.Microsecond)
			//atomic.AddInt64(&shardingDuration, int64(time.Since(begin)))

			worker := s.workers[workerIndex]
			taskCount++
			msg.WorkerAffinity = workerIndex // 标记好worker序号，用于订单的强顺序执行
			worker.process(msg)

			if msg.MsgSeq >= s.respReqMsgSeq.maxSeq {
				log.Printf("break pre-looping for msgSeq2 %d\n", msg.MsgSeq)
				break
			}
		}
		}

		for {
			msg := <-ch
			//workerIndex := taskCount & n // 使用按位与选择 worker
			//begin := time.Now()
			workerIndex = workSharding(msg.Data, nil, taskCount)
			//time.Sleep(10*time.Microsecond)
			//atomic.AddInt64(&shardingDuration, int64(time.Since(begin)))

			worker := s.workers[workerIndex]
			taskCount++
			msg.WorkerAffinity = workerIndex // 标记好worker序号，用于订单的强顺序执行
			worker.process(msg)
		}
	}()
}

// Todo：后续需要完善写入回报时的应对策略
// 1、写入失败时，追加到本地磁盘
// 2、周期性检查磁盘，尝试把本地缓存的未报数据再度推送
// 3、提供一种机制，让客户端可以获得最新的状态
// 4、成功推送消息之后，返回true，其他情况返回false
func (k *StreamAPI) OnTradeResp(tradeResp *types.TradeActionRespReturn) bool {

	key := k.getRespKey(tradeResp)
	if k.respKeys[key] {
		return false
	}

	log.Printf("send for key:%s\n", key)

	// 1、构造回报
	data, de := k.apiAdapter.ConvertTradeResponseMessage(tradeResp)
	if domain_error.ReportIfErrorHappen(de) {
		return false
	}
	// 2、写入回报
	_, err := k.client.SendMessage(data)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fait")
		return false
	} else {
		k.respKeys[key] = true
	}

	return true
}

func (k *StreamAPI) getRespKey(tradeResp *types.TradeActionRespReturn) string {
	// 对于OrderCancelReject，ExecID为空，但是没关系，因为此时ClOrdID代表cancel请求的id，具有唯一性
	return tradeResp.CurrentTradeActionResp.GetCacheKey()
}

// Todo: 测试造成的性能影响
func (k *StreamAPI) onReset() {
	k.respKeys = map[string]bool{}
	k.respReqMsgSeq.respReqMsgSeqs = map[int64]bool{}
	k.respReqMsgSeq.minSeq = 0
	k.respReqMsgSeq.maxSeq = 0
	// 从kafka重新拉取key
	keys, seqs := k.client.GetHistoricalSentKeysAndReqMsgSeqs()
	log.Printf("GetHistoricalSentKeysAndReqMsgSeqs, keys size:%d\n", len(keys))
	for _, key := range keys {
		k.respKeys[key] = true
	}
	if len(seqs) > 0 {
		k.respReqMsgSeq.minSeq = seqs[0]
		k.respReqMsgSeq.maxSeq = seqs[len(seqs)-1]
		for _, seq := range seqs {
			k.respReqMsgSeq.respReqMsgSeqs[seq] = true
		}
	}
}
