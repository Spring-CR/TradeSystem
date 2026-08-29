package ficc_v2

import (
	"log"
	"rhino-core/adapter_registry"
)

// 作者：林春泉

func init() {
    log.Printf("======> RegisterAdapterFunction NewTitansFiccOrderPositionAdapter v2")
    adapter_registry.RegisterAdapterFunction(NewTitansFiccOrderPositionAdapter)
}