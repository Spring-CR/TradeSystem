package domain_error

import (
	"bytes"
	"encoding/json"
	"errors"
	"ficc-utils/common/context"
	"ficc-utils/common/context/constant"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

var errorMessagesInCN = map[string]string{}

type Error struct {
	Code  string `json:"code"`
	Msg   string `json:"msg"`
	Stack string `json:"stack"`
	Err   error  `json:"details"`
}

func (e *Error) ErrorString() string {
	return fmt.Sprintf("ERROR ==> code:%s, message:%s, details:%v, stack:%s", e.Code, e.Msg, e.Err, e.Stack)
}

func (e *Error) SimpleErrorString() string{
	errCode := e.Code
	errMsg := e.Msg
	if e.Err != nil {
		errMsg += ", " + e.Err.Error()
	}
	return fmt.Sprintf("error code:%s, message:%s, \nstack:%s\n", errCode, errMsg, e.Stack)
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

type errorBuilder interface {
	Build(errorCode string, err error, args ...interface{}) *Error
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
	if "1" != os.Getenv("IgnorePrintErrBuildDetails") {
		log.Printf("BUILD ERROR:: code:%s, message:%s, native_error:%v, log_stack:\n%s", errorCode, msg, err, stack)
	}
	return &Error{Code: errorCode, Msg: msg, Err: err, Stack: stack}
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

func Panic(context string, e *Error) {
	builder.Panic(context, e)
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
