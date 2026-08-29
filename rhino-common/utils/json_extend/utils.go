package json_extend

import (
	"reflect"

	"github.com/acsellers/inflections"
)

// TransformToJsonOfUnderscoreMap 将任意结构体或结构体数组转换为包含下划线命名字段的JSON兼容的map或数组
func TransformToJsonOfUnderscoreMap(v interface{}) interface{} {
    val := reflect.ValueOf(v)
    // 处理结构体指针
    for val.Kind() == reflect.Ptr {
        val = val.Elem()
    }
    switch val.Kind() {
    case reflect.Struct:
        underscoreMap := make(map[string]interface{})
        typ := val.Type()
        for i := 0; i < val.NumField(); i++ {
            field := typ.Field(i)
            fieldValue := val.Field(i).Interface()
            underscoreKey := inflections.Underscore(field.Name)
            underscoreMap[underscoreKey] = TransformToJsonOfUnderscoreMap(fieldValue)
        }
        return underscoreMap

    case reflect.Slice, reflect.Array:
        length := val.Len()
        result := make([]interface{}, length)
        for i := 0; i < length; i++ {
            elem := val.Index(i).Interface()
            result[i] = TransformToJsonOfUnderscoreMap(elem)
        }
        return result

    default:
        return v
    }
}