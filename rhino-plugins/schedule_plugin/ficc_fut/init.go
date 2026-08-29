package ficc_fut

import (
	"log"
	"rhino-core/adapter_registry"
)

// 作者：林春泉

func init() {
    log.Printf("======> RegisterAdapterFunction NewFiccFutScheduleAdapter")
    adapter_registry.RegisterAdapterFunction(NewFiccFutScheduleAdapter)
}