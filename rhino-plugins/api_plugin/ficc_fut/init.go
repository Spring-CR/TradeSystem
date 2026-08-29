package ficc_fut

import "rhino-core/adapter_registry"

func init() {
    adapter_registry.RegisterAdapterFunction(NewFiccFurAPIAdapter)
}