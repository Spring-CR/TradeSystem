package ficc

import "rhino-core/adapter_registry"

func init() {
    adapter_registry.RegisterAdapterFunction(NewTitansFiccAPIAdapter)
}