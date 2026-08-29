package adapter_registry

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
)

// 注册器
var funcRegistry = make(map[string]reflect.Value)


// 注册适配器函数
func RegisterAdapterFunction(fn interface{}) {
    packageName := getAdapterFunctionPath(fn)
	log.Printf("registry package %s\n", packageName)
    funcRegistry[packageName] = reflect.ValueOf(fn)
}

// 获取包名
func getAdapterFunctionPath(fn interface{}) string {
    funcPath := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	i := strings.LastIndex(funcPath, ".")
	if i >= 0 {
		funcPath = funcPath[:i]
	}
    return funcPath
}

// 调用函数
func CallAdapterFunction(pkgName string, args ...interface{}) (interface{}, interface{}, error) {
    log.Printf(">>>CallAdapterFunction, pkgName:%s, len(args)=%d\n", pkgName, len(args))
    if fn, ok := funcRegistry[pkgName]; ok {
        reflectArgs := make([]reflect.Value, len(args))
        for i, arg := range args {
            reflectArgs[i] = reflect.ValueOf(arg)
        }
        result := fn.Call(reflectArgs)
        
        if len(result) == 2 {
            return result[0].Interface(), result[1].Interface(), nil
        }
    }
    return nil, nil, fmt.Errorf("适配器函数在包 %s 未找到", pkgName)
}