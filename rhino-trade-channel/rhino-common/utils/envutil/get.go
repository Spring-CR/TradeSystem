package envutil

import (
	"strconv"
	"fmt"
	"time"
	"os"
	"strings"
)

func GetIntValue(envMap map[string]string, envKey string, defaultValue string)int{
	valStr := envMap[envKey]
	if len(valStr)==0{
		valStr = defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err!=nil{
		panic(fmt.Sprintf("%s is not int format : %+v", valStr, err))
		time.Sleep(5*time.Second)
	}
	return val
}

func GetStrValue(envMap map[string]string, envKey string, defaultValue string)string{
	valStr := envMap[envKey]
	if len(valStr)==0 {
		valStr = defaultValue
	}
	return valStr
}


func GetEnvVarValue(envName , defaultValue string) string{
	value := os.Getenv(envName)
	if strings.EqualFold("", value){
		value = defaultValue
	}
	return value
}
