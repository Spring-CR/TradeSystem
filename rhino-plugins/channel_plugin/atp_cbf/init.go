package atp_cbf

import "rhino-core/adapter_registry"

// 作者：林春泉

func init() {
    adapter_registry.RegisterAdapterFunction(NewAtpFIXChannel)
}