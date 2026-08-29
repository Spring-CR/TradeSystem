package json_extend_test

import (
	"encoding/json"
	"log"
	"rhino-common/utils/json_extend"
	"rhino-instr/schema"
	"testing"
)

func TestTransformToJsonOfUnderscoreMap(t *testing.T){
	data := &schema.TaskInstr{}
	v := json_extend.TransformToJsonOfUnderscoreMap(data)
	jsData, _ := json.MarshalIndent(v, "", "  ")
	log.Printf("jsData:%s\n", jsData)
}