package ficc

import (
	"log"
	"rhino-core/adapter_registry"
)

// 作者：林春泉

func init() {
    log.Printf("======> RegisterAdapterFunction NewFiccFutOrderStatusAdapter")
    adapter_registry.RegisterAdapterFunction(NewFiccFutOrderStatusAdapter)
}