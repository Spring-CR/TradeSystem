package fixutil

import (
	"rhino-common/domain_error"
	"rhino-common/enum"
	"strconv"
	"strings"

	"github.com/quickfixgo/quickfix"
)

func GetValueByTag(message *quickfix.Message, tagName string, tag quickfix.Tag, fieldMap quickfix.FieldMap, valType enum.AttrValueType) (val interface{}, de *domain_error.Error) {
	switch valType {
	case enum.AttrValueType_INT:
		value, rejErr := fieldMap.GetInt(tag)
		if rejErr != nil {
			de = domain_error.Build(domain_error.ILLEGAL_FIT_TAG_ERR_CODE, rejErr, strings.Title(tagName), tag)
			return
		}
		val = value
		return
	case enum.AttrValueType_FLOAT:
		value, rejErr := fieldMap.GetString(tag)
		if rejErr != nil {
			de = domain_error.Build(domain_error.ILLEGAL_FIT_TAG_ERR_CODE, rejErr, strings.Title(tagName), tag)
			return
		}
		fVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			de = domain_error.Build(domain_error.ILLEGAL_FIT_TAG_ERR_CODE, err, strings.Title(tagName), tag)
			return
		}
		val = fVal
		return
	case enum.AttrValueType_BOOL:
		value, rejErr := fieldMap.GetString(tag)
		value = strings.ToUpper(value)
		if rejErr != nil {
			de = domain_error.Build(domain_error.ILLEGAL_FIT_TAG_ERR_CODE, rejErr, strings.Title(tagName), tag)
			return
		}
		switch value {
		case "Y", "TRUE":
			val = true
		case "N", "FALSE":
			val = false
		default:
			de = domain_error.Build(domain_error.ILLEGAL_FIT_TAG_ERR_CODE, rejErr, strings.Title(tagName), tag)
			return
		}
		return
	case enum.AttrValueType_STRING:
		value, rejErr := fieldMap.GetString(tag)
		if rejErr != nil {
			de = domain_error.Build(domain_error.ILLEGAL_FIT_TAG_ERR_CODE, rejErr, strings.Title(tagName), tag)
			return
		}
		val = value
		return
	}
	return
}

func SetValueByTag(msgProps map[string]interface{}, message *quickfix.Message, tagName string, tag quickfix.Tag, fieldMap quickfix.FieldMap, valType enum.AttrValueType) (de *domain_error.Error) {
	var val interface{}
	val, de = GetValueByTag(message, tagName, tag, fieldMap, valType)
	if de != nil {
		return
	}
	msgProps[tagName] = val
	return
}

func ConvertDomainErrToFixRejection(de*domain_error.Error, refTagID quickfix.Tag) (rejErr quickfix.MessageRejectError) {
	if de == nil {
		return
	}
	reason, _ := strconv.Atoi(de.Code)
	return quickfix.NewMessageRejectError(de.Msg, reason, &refTagID)
}