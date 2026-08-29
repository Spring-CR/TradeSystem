package utils

import "reflect"

func CopyStruct(src, dst any) {
	srcVal := reflect.Indirect(reflect.ValueOf(src))
	dstVal := reflect.Indirect(reflect.ValueOf(dst))

	// 只处理结构体
	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		return
	}

	dstType := dstVal.Type()
	// 遍历目标结构体字段
	for i := 0; i < dstType.NumField(); i++ {
		dstField := dstType.Field(i)
		// 找源结构体同名字段
		srcFieldVal := srcVal.FieldByName(dstField.Name)
		if !srcFieldVal.IsValid() {
			continue
		}
		// 类型相同才能赋值
		if srcFieldVal.Type() == dstField.Type {
			dstVal.Field(i).Set(srcFieldVal)
		}
	}
}