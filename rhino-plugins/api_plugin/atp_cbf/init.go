package atp_cbf

import (
	"log"
	"rhino-core/adapter_registry"
)

// 作者：林春泉

func init() {
    log.Printf("======> RegisterAdapterFunction NewTitansCrossBorderFutureAPIAdapter")
    adapter_registry.RegisterAdapterFunction(NewTitansCrossBorderFutureAPIAdapter)
}