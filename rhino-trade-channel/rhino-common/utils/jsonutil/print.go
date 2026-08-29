package jsonutil

import (
	"log"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func Print(msgPrefix string, obj interface{}) {
	jsData, _ := json.MarshalIndent(obj, "", "  ")
	log.Println(msgPrefix + string(jsData))
}

func PrintSimple(msgPrefix string, obj interface{}) {
	jsData, _ := json.Marshal(obj)
	log.Println(msgPrefix + string(jsData))
}
