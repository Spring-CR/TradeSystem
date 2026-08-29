package fix

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/utils/timeutil"
	"strings"
	"time"

	"github.com/quickfixgo/quickfix"
)

func (a *FixApi) schedule() {
	begin, end, layout, local := a.fixApiAdapter.GetFixPortOpenTimeRange()
	timeutil.StarTimeRangeSchedule(begin, end, layout, local, func(duringTimeRange bool) {
		if duringTimeRange {
			for {
				err := a.acceptor.Start()
				if err != nil {
					domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("unable to start FIX acceptor from %s\n", a.fixApiAdapter.GetConfigPath()))
					time.Sleep(5 * time.Second)
				} else {
					return
				}
			}
		} else {
			for _, sessionID := range a.sessionIDs {
				for {
					log.Printf("start to reset session:%v, a==nil:%v\n", sessionID, a==nil)
					err := quickfix.ResetSession(sessionID)
					log.Printf("finish reset session:%v, error:%v\n", sessionID, err)
					if err != nil && strings.HasPrefix(err.Error(), "Unknown session") {
						err = nil
					}
					if err != nil {
						domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("unable to ResetSession %v\n", sessionID))
						time.Sleep(5 * time.Second)
					} else {
						break
					}
				}
				log.Printf("reset session:%v\n", sessionID)
			}
		}
	}, func() {
		for {
			err := a.acceptor.Start()
			if err != nil {
				domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("unable to start FIX acceptor from %s\n", a.fixApiAdapter.GetConfigPath()))
				time.Sleep(5 * time.Second)
			} else {
				return
			}
		}
	}, func() {
		// 清理session
		normalReset := false
		for _, sessionID := range a.sessionIDs {
			for {
				log.Printf("start to reset session:%v\n", sessionID)
				err := quickfix.ResetSession(sessionID)
				log.Printf("finish reset session:%v, error:%v\n", sessionID, err)
				if err != nil && strings.HasPrefix(err.Error(), "Unknown session") {
					err = nil
				} else {
					normalReset = true
				}
				if err != nil {
					domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("unable to ResetSession %v\n", sessionID))
					time.Sleep(5 * time.Second)
				} else {
					break
				}
			}
			log.Printf("reset session:%v\n", sessionID)
		}
		if normalReset {
			a.acceptor.Stop()
		}
	})
}
