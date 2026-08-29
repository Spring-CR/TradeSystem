package datamap

import (
	"log"
	"time"
)

func (r *AutoSyncRepo) Reset() {
	if r.tableMaps != nil {
		for _, m := range r.tableMaps {
			m.Reset()
		}
	}

	// wait for data ready
	for {
		allReady := true
		for _, m := range r.tableMaps {
			if m.dataSyncLog == nil {
				log.Printf("======> %s is not ready, keep waiting..\n", m.tableConfig.TableName)
				allReady = false
			} else {
				log.Printf("======> %s is ready\n", m.tableConfig.TableName)
			}
		}
		if allReady {
			break
		}
		time.Sleep(10*time.Second)
	}
}