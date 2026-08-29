package synclogs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	syncLock sync.RWMutex
	syncLogs []SyncLog
	syncLogsQryUrl string
	syncLogsInterval time.Duration
)

type SyncLog struct {
	TableName string
	ReportTime int64
}

func GetTableSyncTime(table string) *time.Time {
	l, err := GetSyncLogs()
	if err != nil {
		return nil
	}
	for _, v := range l {
		if v.TableName == table {
			tm := time.UnixMilli(v.ReportTime)
			return &tm
		}
	}
	return nil
}
func GetSyncLogs() ([]SyncLog, error) {
	syncLock.RLock()
	l := syncLogs
	syncLock.RUnlock()
	if len(l) > 0 {
		return l, nil
	}

	err := updateSyncLogs()
	if err != nil {
		return nil, err
	}

	syncLock.RLock()
	l = syncLogs
	syncLock.RUnlock()
	if len(l) > 0 {
		return l, nil
	}
	return nil, fmt.Errorf("SyncLogs empty")
}

func qrySyncLogs() (result []SyncLog, err error) {
	resp, err := http.Get(syncLogsQryUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	output := []SyncLog{}
	err = json.Unmarshal(body, &output)
	if err != nil {
		log.Printf("qrySyncLogs body:%s\n", body)
		return nil, err
	}

	log.Printf("qrySyncLogs output:\n")
	for _, v := range output {
		log.Printf("qrySyncLogs, table:%s, sync:%s", v.TableName, time.UnixMilli(v.ReportTime).Format(time.DateTime+".000"))
	}
	log.Printf("qrySyncLogs output end\n")
	return output, nil
}

func updateSyncLogs() error {
	l, err := qrySyncLogs()
	if err != nil {
		return err
	}
	syncLock.Lock()
	syncLogs = l
	syncLock.Unlock()
	return nil
}

func InitSyncLogs(qryUrl, interval string) {
	syncLogsQryUrl = qryUrl
	d, err := time.ParseDuration(interval)
	if err != nil {
		log.Printf("InitSyncLogs error:%v, SyncLogsInterval use default:30min", err)
		syncLogsInterval = 30 * time.Minute
	} else {
		syncLogsInterval = d
	}
	log.Printf("InitSyncLogs SyncLogsInterval:%.2fmin", syncLogsInterval.Minutes())

	go func() {
		for {
			err := updateSyncLogs()
			if err != nil {
				log.Printf("updateSyncLogs error:%v\n", err)
			}
			time.Sleep(syncLogsInterval)
		}
	}()
}