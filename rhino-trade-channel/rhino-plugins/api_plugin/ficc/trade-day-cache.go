package ficc

import (
	"log"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"strings"
	"time"
)

var (
	t1DayCache map[string]string = make(map[string]string)
)

func refreshT1DayCache(configMap map[string]*schema.ApplicationCfgItem) {
	begin := time.Now()
	t1DayCache = make(map[string]string)
	for i := -3; i < 7; i++ {
		currDate := time.Now().Add(time.Duration(i*24) * time.Hour).In(timeutil.CnTimeLocation).Format(time.DateOnly)
		settleDate := getTradeDate(configMap, currDate, 1)
		settleDate = strings.ReplaceAll(settleDate, "-", "")
		t1DayCache[currDate] = settleDate
	}
	js, _ := json.Marshal(t1DayCache)
	log.Printf("After refreshT1DayCache, t1DayCache=%s, time cost=%v\n", js, time.Since(begin))
}
