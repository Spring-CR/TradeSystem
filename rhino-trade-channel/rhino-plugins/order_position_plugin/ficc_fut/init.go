package ficc_fut

import (
	"log"
	"rhino-core/adapter_registry"
)

// 作者：林春泉

func init() {
    log.Printf("======> RegisterAdapterFunction NewFiccFutOrderPositionAdapter")
    adapter_registry.RegisterAdapterFunction(NewFiccFutOrderPositionAdapter)
}