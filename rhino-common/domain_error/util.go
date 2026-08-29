package domain_error

import (
	"fmt"
	"log"
	"time"
)

func ProcessSevereError(doPanic bool, sleepSecondBeforePanic int, de *Error, err error, errMsg string) (errHappen bool) {
	errHappen = de != nil || err != nil
	if de == nil && err != nil {
		err = fmt.Errorf(errMsg+", error:%v", err)
		de = builder.Build(GENERIC_ERR_CODE, err)
		log.Printf("ProcessSevereError error:%s\n", de.ErrorString())
	}
	if errHappen {
		go func() {
			ReportIfErrorHappen(de)
			if doPanic {
				time.Sleep(time.Duration(sleepSecondBeforePanic) * time.Second)
				panic(de.ErrorString())
			}
		}()
	}
	return
}

func ProcessSevereError3(doPanic bool, sleepSecondBeforePanic int, de *Error, err error, errMsg string) (errHappen bool) {
	errHappen = de != nil || err != nil
	var stack string
	if errHappen {
		stack = GetStack()
	}
	go func(stack string) {
		if de != nil {
			de.Stack = stack
			ReportIfErrorHappen(de)
		} else if err != nil {
			err = fmt.Errorf(errMsg+", error:%v", err)
			de = builder.Build(GENERIC_ERR_CODE, err)
			de.Stack = stack
			ReportIfErrorHappen(de)
		}
		if de != nil {
			log.Println(de.ErrorString())
			if doPanic {
				time.Sleep(time.Duration(sleepSecondBeforePanic) * time.Second)
				panic(de.ErrorString())
			}
		}
	}(stack)
	return
}

func ProcessSevereError_2(doPanic bool, sleepSecondBeforePanic int, de *Error, err error, errMsg string) (errHappen bool) {
	if de != nil {
		errHappen = true
		ReportIfErrorHappen(de)
		if doPanic {
			time.Sleep(time.Duration(sleepSecondBeforePanic) * time.Second)
			panic(de.ErrorString())
		}
		return
	} else if err != nil {
		errHappen = true
		de = builder.Build(GENERIC_ERR_CODE, err)
		ReportIfErrorHappen(de)
		if doPanic {
			time.Sleep(time.Duration(sleepSecondBeforePanic) * time.Second)
			panic(de.ErrorString())
		}
	}
	return
}

func ExtendErrorMessage(summary string, originalErr error) error {
	if originalErr == nil {
		return nil
	}
	return fmt.Errorf("%s, %v", summary, originalErr)
}