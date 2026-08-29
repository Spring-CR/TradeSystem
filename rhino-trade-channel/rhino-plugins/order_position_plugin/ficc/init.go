package ficc

import (
	"log"
	"rhino-core/adapter_registry"
)

// 作者：林春泉

func init() {
    log.Printf("======> RegisterAdapterFunction NewTitansFiccOrderPositionAdapter")
    adapter_registry.RegisterAdapterFunction(NewTitansFiccOrderPositionAdapter)
}