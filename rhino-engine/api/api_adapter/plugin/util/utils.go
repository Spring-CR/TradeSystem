package util

import (
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
)

func GetIntValueInField(msgProps map[string]interface{}, field string) (val int, ok bool, de *domain_error.Error) {
	var err error
	var v interface{}
	v, ok, err = attrutil.GetAttrValue(msgProps, field, enum.AttrValueType_INT)
	if err != nil {
		de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err)
		return
	}
	val = v.(int)
	return
}

func GetStringValueInField(msgProps map[string]interface{}, field string) (val string, ok bool, de *domain_error.Error) {
	var err error
	var v interface{}
	v, ok, err = attrutil.GetAttrValue(msgProps, field, enum.AttrValueType_STRING)
	if err != nil {
		de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err, v)
		return
	}
	val = v.(string)
	return
}

func GetFloatValueInField(msgProps map[string]interface{}, field string) (val float64, ok bool, de *domain_error.Error) {
	var err error
	var v interface{}
	v, ok, err = attrutil.GetAttrValue(msgProps, field, enum.AttrValueType_FLOAT)
	if err != nil {
		de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err)
		return
	}
	val = v.(float64)
	return
}

func GetStringValuePerhapsInTwoFields(msgProps map[string]interface{}, fieldOne string, fieldTwo string) (val string, ok bool, de *domain_error.Error) {
	var err error
	var v interface{}
	v, ok, err = attrutil.GetAttrValue(msgProps, fieldOne, enum.AttrValueType_STRING)
	if err != nil {
		de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err)
		return
	}
	if !ok {
		v, ok, err = attrutil.GetAttrValue(msgProps, fieldTwo, enum.AttrValueType_STRING)
		if err != nil {
			de = domain_error.Build(domain_error.CONVERT_NEW_ORDER_SINGLE_MSG_ERR_CODE, err)
			return
		}
	}
	val = v.(string)
	return
}