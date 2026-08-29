package dbutil_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"rhino-common/utils/dbutil"
	"rhino-instr/schema"
	"testing"
)

// go test -v common-query_test.go
// go test

func TestGetFieldTypeMap(t *testing.T) {
	typ := reflect.TypeOf((*schema.TaskInstr)(nil)).Elem()
	m := dbutil.GetFieldTypeMap(typ)
	jsData, _ := json.MarshalIndent(m, "", " ")
	fmt.Printf("%s\n", jsData)

	m = dbutil.GetFieldTypeMap(typ)
	jsData, _ = json.MarshalIndent(m, "", " ")
	fmt.Printf("%s\n", jsData)
}

func TestCheckValueType(t *testing.T) {
	typ := reflect.TypeOf((*schema.TaskInstr)(nil)).Elem()
	m := dbutil.GetFieldTypeMap(typ)

	fieldConditions := []*dbutil.FieldCondition{
		{
			Field     : "my_field",
			ValueType : 0,
			Value     : 1.0,
		},
	}
	err := dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error :%+v\n", err)

	fieldConditions = []*dbutil.FieldCondition{
		{
			Field     : "my_field",
			ValueType : 1,
			Value     : 1.0,
		},
	}
	err = dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error0:%+v\n", err)

	fieldConditions = []*dbutil.FieldCondition{
		{
			Field     : "batch_serial_no",
			ValueType : 0,
			Value     : 1.2,
		},
	}
	err = dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error1:%+v\n", err)

	fieldConditions = []*dbutil.FieldCondition{
		{
			Field     : "batch_serial_no",
			ValueType : 1,
			Value     : 1.2,
		},
	}
	err = dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error2:%+v\n", err)

	fieldConditions = []*dbutil.FieldCondition{
		{
			Field     : "batch_serial_no",
			ValueType : 1,
			Value     : []interface{}{1.2},
		},
	}
	err = dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error3:%+v\n", err)

	fieldConditions = []*dbutil.FieldCondition{
		{
			Field     : "batch_serial_no",
			ValueType : 1,
			Value     : []interface{}{1, 2, 3},
		},
	}
	err = dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error3.1:%+v\n", err)

	fieldConditions = []*dbutil.FieldCondition{
		{
			Field     : "batch_serial_no",
			ValueType : 2,
			Value     : 1.2,
		},
	}
	err = dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error4:%+v\n", err)

	fieldConditions = []*dbutil.FieldCondition{
		{
			Field     : "batch_serial_no",
			ValueType : 2,
			Value     : []interface{}{1.2},
		},
	}
	err = dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error5:%+v\n", err)

	////
	fieldConditions = []*dbutil.FieldCondition{
		{
			Field     : "batch_serial_no",
			ValueType : 0,
			Value     : 1,
		},
	}
	err = dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error1:%+v\n", err)

	fieldConditions = []*dbutil.FieldCondition{
		{
			Field     : "batch_serial_no",
			ValueType : 1,
			Value     : 1,
		},
	}
	err = dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error2:%+v\n", err)

	fieldConditions = []*dbutil.FieldCondition{
		{
			Field     : "batch_serial_no",
			ValueType : 1,
			Value     : []interface{}{1},
		},
	}
	err = dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error3:%+v\n", err)

	fieldConditions = []*dbutil.FieldCondition{
		{
			Field     : "batch_serial_no",
			ValueType : 2,
			Value     : 1,
		},
	}
	err = dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error4:%+v\n", err)

	fieldConditions = []*dbutil.FieldCondition{
		{
			Field     : "batch_serial_no",
			ValueType : 2,
			Value     : []interface{}{1},
		},
	}
	err = dbutil.CheckValueType(m, fieldConditions)
	fmt.Printf("error5:%+v\n", err)
}

func TestPrepareQueryStatment(t *testing.T) {
	
	typ := reflect.TypeOf((*schema.TaskInstr)(nil)).Elem()
	fieldConditions := []*dbutil.FieldCondition{
		{
			Field     : "date",
			ValueType : 1,
			Value     : []interface{}{20240501,20240531},//[]int{20240501, 20240531},
		},
		{
			Field     : "direct_operator",
			ValueType : 0,
			Value     : "d001",
		},
		{
			Field     : "direct_operator",
			ValueType : 2,
			Value     : []interface{}{"o001", "o002", "o002"},//[]string{"operator1", "operator2", "operator3"},
		},
	}

	whereClause, args, err := dbutil.PrepareQueryStatment(typ, fieldConditions)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		return
	}

	fmt.Printf("whereClause:%s\nargs:%+v\n", whereClause, args)

	jsData, _ := json.Marshal(fieldConditions)
	fmt.Printf("jsData:%s\n", jsData)
	
	var fieldConditions2[]*dbutil.FieldCondition
	err = json.Unmarshal(jsData, &fieldConditions2)
	if err != nil {
		fmt.Printf("error1:%v\n", err)
		return
	}

	whereClause, args, err = dbutil.PrepareQueryStatment(typ, fieldConditions)
	if err != nil {
		fmt.Printf("error1:%v\n", err)
		return
	}

	fmt.Printf("whereClause2:%s\nargs2:%+v\n", whereClause, args)
}