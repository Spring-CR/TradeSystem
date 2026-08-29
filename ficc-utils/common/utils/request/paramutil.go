package request

import (
	"ficc-utils/common/domain_error"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetParamAsString(c*gin.Context, key string, allowEmpty bool)(val string, err *domain_error.Error){
	val = c.Params.ByName(key)
	err = checkEmpty(key, val, allowEmpty)
	if err!=nil{
		return val, err
	}
	return val, nil
}

func GetQueryAsString(c*gin.Context, key string, allowEmpty bool)(val string, err *domain_error.Error){
	val = c.Query(key)
	err = checkEmpty(key, val, allowEmpty)
	if err!=nil{
		return val, err
	}
	return val, nil
}

func GetQueryAsStringSlice(c*gin.Context, key string, allowEmpty bool)(val []string, err *domain_error.Error){
	str := c.Query(key)
	err = checkEmpty(key, str, allowEmpty)
	if err!=nil{
		return val, err
	}
	strs := strings.Split(str,",")
	for _, s := range strs{
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			val = append(val, s)
		}
	}
	if !allowEmpty && len(val)==0{
		return val, domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE,nil, key)
	}
	return val, nil
}

func GetParamAsInt(c*gin.Context, key string, allowEmpty bool)(val int, err *domain_error.Error){
	strVal := c.Params.ByName(key)
	err = checkEmpty(key, strVal, allowEmpty)
	if err!=nil{
		return val, err
	}
	if allowEmpty && len(strVal) == 0 {
		return 0, nil
	}
	v, e := strconv.Atoi(strVal)
	if e!=nil{
		return val, domain_error.Build(domain_error.API_PARAM_PARSING_ERR_CODE, e, key, strVal)
	}
	return v, nil
}

func GetQueryAsInt(c*gin.Context, key string, allowEmpty bool)(val int, err *domain_error.Error){
	strVal := c.Query(key)
	err = checkEmpty(key, strVal, allowEmpty)
	if err!=nil{
		return val, err
	}
	if allowEmpty && len(strVal) == 0 {
		return 0, nil
	}
	v, e := strconv.Atoi(strVal)
	if e!=nil{
		return val, domain_error.Build(domain_error.API_PARAM_PARSING_ERR_CODE, e, key, strVal)
	}
	return v, nil
}

func GetQueryAsBool(c*gin.Context, key string, allowEmpty bool)(val bool, err *domain_error.Error){
	strVal := c.Query(key)
	err = checkEmpty(key, strVal, allowEmpty)
	if err!=nil{
		return val, err
	}
	return strings.EqualFold(strVal, "true") || strVal == "1", nil
}

func GetParamAsFloat(c*gin.Context, key string, allowEmpty bool)(val float64, err *domain_error.Error){
	strVal := c.Params.ByName(key)
	err = checkEmpty(key, strVal, allowEmpty)
	if err!=nil{
		return val, err
	}
	if allowEmpty && len(strVal) == 0 {
		return 0, nil
	}
	v, e :=  strconv.ParseFloat(strVal, 64)
	if e!=nil{
		return val, domain_error.Build(domain_error.API_PARAM_PARSING_ERR_CODE, e, key, strVal)
	}
	return v, nil
}

func GetQueryAsFloat(c*gin.Context, key string, allowEmpty bool)(val float64, err *domain_error.Error){
	strVal := c.Query(key)
	err = checkEmpty(key, strVal, allowEmpty)
	if err!=nil{
		return val, err
	}
	if allowEmpty && len(strVal) == 0 {
		return 0, nil
	}
	v, e := strconv.ParseFloat(strVal, 64)	
	if e!=nil{
		return val, domain_error.Build(domain_error.API_PARAM_PARSING_ERR_CODE, e, key, strVal)
	}
	return v, nil
}

/*
func GetParamAsBool(c*gin.Context, key string)(val bool){
	strVal := c.Params.ByName(key)
	if strings.EqualFold("1", strVal) || strings.EqualFold("true", strings.ToLower(strVal)){
		return true
	}
	return false
}*/

func checkEmpty(key, value string, allowEmpty bool)*domain_error.Error{
	if !allowEmpty && len(value)==0{
		return domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE,nil, key)
	}
	return nil
}

func AddQueryParam(req*http.Request, key string, value interface{}){
	q := req.URL.Query()
	q.Add(key, fmt.Sprintf("%v",value))
	req.URL.RawQuery = q.Encode()
}