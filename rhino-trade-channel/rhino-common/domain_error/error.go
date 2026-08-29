package domain_error

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"rhino-common/context"
	"rhino-common/context/constant"
	"rhino-common/utils/emailalert"
	"rhino-core/schema"
	"runtime"
	"strconv"
	"time"
)

var errorMessagesInCN = map[string]string{}

const (
	ERROR   int = 0
	WARNING int = 1
)

var (
	CnTimeLocation, _   = time.LoadLocation("Asia/Shanghai") // 显式加载中国时区
	TransactTimeLayout  = "20060102-15:04:05.000"
	ErrorNotifyFunction func(*Error)
)

type Error struct {
	Code      string             `json:"code"`
	Msg       string             `json:"msg"`
	Stack     string             `json:"stack"`
	Err       error              `json:"details"`
	Level     int                `json:"level"` // 0 - error; 1 - warning
	Order     *schema.TradeOrder `json:"order"`
	Timestamp string             `json:"timestamp"`
}

func (e *Error) ErrorString() string {
	return fmt.Sprintf("ERROR ==> code:%s, message:%s, details:%v, stack:%s", e.Code, e.Msg, e.Err, e.Stack)
}

func (e *Error) SimpleErrorString() string {
	errCode := e.Code
	errMsg := e.Msg
	if e.Err != nil {
		errMsg += ", " + e.Err.Error()
	}
	return fmt.Sprintf("error code:%s, message:%s, \nstack:%s\n", errCode, errMsg, e.Stack)
}

func (e *Error) SimpleErrorStringForUser() string {
	return fmt.Sprintf("error:%s\n", e.Msg)
}

func (e *Error) MarshalJSON() (data []byte, err error) {
	type ErrorAlias Error
	errStruct := &struct {
		Err string `json:"details"`
		*ErrorAlias
	}{
		ErrorAlias: (*ErrorAlias)(e),
	}
	if e.Err != nil {
		errStruct.Err = e.Err.Error()
	}
	return json.Marshal(errStruct)
}

func (e *Error) UnmarshalJSON(data []byte) (err error) {
	type ErrorAlias Error
	errStruct := &struct {
		Err string `json:"details"`
		*ErrorAlias
	}{
		ErrorAlias: (*ErrorAlias)(e),
	}

	if err = json.Unmarshal(data, &errStruct); err != nil {
		return err
	}

	if len(errStruct.Err) > 0 {
		e.Err = errors.New(errStruct.Err)
	}

	return nil
}

func (e *Error) Refine(level int, order *schema.TradeOrder) {
	e.Level = level
	e.Order = order
	if ErrorNotifyFunction != nil {
		ErrorNotifyFunction(e)
	}
}

type errorBuilder interface {
	Build(errorCode string, err error, args ...interface{}) *Error
	BuildWithDetails(level int, order *schema.TradeOrder, errorCode string, err error, args ...interface{}) *Error
	Panic(context string, e *Error)
}

type errorBuilderInCN struct {
	errorBuilder
}

func (eb *errorBuilderInCN) Build(errorCode string, err error, args ...interface{}) *Error {
	msg, ok := errorMessagesInCN[errorCode]
	if ok {
		if len(args) > 0 {
			msg = fmt.Sprintf(msg, args...)
		}
	} else {
		msg = fmt.Sprintf("%s", err)
	}

	stackBuf := &bytes.Buffer{}
	for i := 1; i <= 4; i++ {
		_, fn, line, _ := runtime.Caller(i + 1)
		stackBuf.WriteString(fn)
		stackBuf.WriteByte(':')
		stackBuf.WriteString(strconv.Itoa(line))
		stackBuf.WriteByte('\n')
	}
	stack := stackBuf.String()
	if os.Getenv("IgnorePrintErrBuildDetails") != "1" {
		log.Printf("BUILD ERROR:: code:%s, message:%s, native_error:%v, log_stack:\n%s", errorCode, msg, err, stack)
	}
	de := &Error{Code: errorCode, Msg: msg, Err: err, Stack: stack, Timestamp: time.Now().In(CnTimeLocation).Format(TransactTimeLayout)}
	if ErrorNotifyFunction != nil {
		ErrorNotifyFunction(de)
	}
	return de
}

func (eb *errorBuilderInCN) BuildWithDetails(level int, order *schema.TradeOrder, errorCode string, err error, args ...interface{}) *Error {
	msg, ok := errorMessagesInCN[errorCode]
	if ok {
		if len(args) > 0 {
			msg = fmt.Sprintf(msg, args...)
		}
	} else {
		msg = fmt.Sprintf("%s", err)
	}

	stackBuf := &bytes.Buffer{}
	for i := 1; i <= 4; i++ {
		_, fn, line, _ := runtime.Caller(i + 1)
		stackBuf.WriteString(fn)
		stackBuf.WriteByte(':')
		stackBuf.WriteString(strconv.Itoa(line))
		stackBuf.WriteByte('\n')
	}
	stack := stackBuf.String()
	if os.Getenv("IgnorePrintErrBuildDetails") != "1" {
		log.Printf("BUILD ERROR:: code:%s, message:%s, native_error:%v, log_stack:\n%s", errorCode, msg, err, stack)
	}
	de := &Error{Level: level, Order: order, Code: errorCode, Msg: msg, Err: err, Stack: stack, Timestamp: time.Now().In(CnTimeLocation).Format(TransactTimeLayout)}
	if ErrorNotifyFunction != nil {
		ErrorNotifyFunction(de)
	}
	return de
}

func (eb *errorBuilderInCN) Panic(context string, e *Error) {
	if e != nil {
		errorString := "CONTEXT:: " + context + ", " + e.ErrorString()
		panic(errorString)
	}
}

var builder errorBuilder

func init() {
	switch context.Lang {
	case constant.Lang_CN:
		builder = &errorBuilderInCN{}
	default:
		builder = &errorBuilderInCN{}
	}
}

func Build(errorCode string, err error, args ...interface{}) *Error {
	return builder.Build(errorCode, err, args...)
}

func BuildWithDetails(level int, order *schema.TradeOrder, errorCode string, err error, args ...interface{}) *Error {
	return builder.BuildWithDetails(level, order, errorCode, err, args...)
}

func Panic(context string, e *Error) {
	builder.Panic(context, e)
}

func ReportIfErrorHappen(e *Error) (errHappen bool) {
	if e != nil {
		errHappen = true
		err := emailalert.Send(e.Msg, e.ErrorString())
		if err != nil {
			log.Printf("send error by email error:%v\n", err)
		}
	}
	return errHappen
}

func GetStack() string {
	stackBuf := &bytes.Buffer{}
	for i := 1; i <= 5; i++ {
		_, fn, line, _ := runtime.Caller(i + 1)
		stackBuf.WriteString(fn)
		stackBuf.WriteByte(':')
		stackBuf.WriteString(strconv.Itoa(line))
		stackBuf.WriteByte('\n')
	}
	stack := stackBuf.String()
	return stack
}

func WrapErrWithStack(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("error:%+v, stack:%s", err, GetStack())
}

func registerErrMsg(errMsgs map[string]string) {
	for k, v := range errMsgs {
		_, ok := errorMessagesInCN[k]
		if ok {
			log.Fatalf("duplicated error message key : %s\n", k)
		}
		errorMessagesInCN[k] = v
	}
}

func (e *Error) ToSimpleError() error {
	if e.Err != nil {
		return e.Err
	}
	return errors.New(e.ErrorString())
}
