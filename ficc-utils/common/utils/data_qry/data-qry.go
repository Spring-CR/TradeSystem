package data_qry

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"ficc-utils/common/utils/synclogs"
)
var (
	symbolsT0 CachedData[[]SymbolT0]
)

type CachedData[T any] struct {
	Data T
	UpdateTime *time.Time
}

type SymbolT0 struct {
	Symbol string `json:"symbol"`
	SymbolName string `json:"symbolName"`
}

type Symbol struct {
	Symbol string `json:"symbol"`
	SymbolName string `json:"SecurityName"`
	OverSoldIndicator string `json:"OverSoldIndicator"`//T0_LABELS
}

// 查询数据 测试环境"http://10.51.136.21:6093/api/v1/data_qry"
func QryTableData(dataQryUrl, table, key string, v any ) error {
	qryUrl := fmt.Sprintf("%s?collection=%s&key=%s&maxRecords=%d", dataQryUrl, table, key, 1000000)
	resp, err := http.Get(qryUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		log.Printf("qryTableData %s body:%s\n", qryUrl, body)
		return err
	}

	jsonData, _ := json.MarshalIndent(v, "", "   ")
	log.Printf("qryTableData %s result:%s\n", qryUrl, jsonData)
	return nil
}

// 根据syncLogs缓存数据
func SyncCachedData[T any](syncTable string, cache *CachedData[T], qryFunc func() (T, error)) error {
	// synclogs 判断数据是否有更新
	syncTime := synclogs.GetTableSyncTime(syncTable)
	if cache.UpdateTime != nil && syncTime != nil && !syncTime.After(*cache.UpdateTime) {
		return nil
	}
	data, err := qryFunc()
	if err != nil {
		return err
	}
	cache.Data = data
	cache.UpdateTime = syncTime
	return nil
}

//查询交易时间
func QryTradingTime(dataQryUrl string) (openTime, closeTime time.Time, err error) {
	output := struct {
		Code string `json:"code"`
		Msg string `json:"msg"`
		Stack string `json:"stack"`
		Total int32 `json:"total"`
		Data []struct {
			ConfigItemName string `json:"ConfigItemName"`
			ConfigItemValue string `json:"ConfigItemValue"`
		} `json:"data"`
	}{}

	err = QryTableData(dataQryUrl, "ApplicationConfig", "TradingTime", &output)
	if err != nil {
		return openTime, closeTime, err
	}

	if output.Code != "" {
		return openTime, closeTime, errors.New(output.Msg)
	}
	if len(output.Data) < 1 {
		return openTime, closeTime, errors.New("QryTradingTime result is empty")
	}

	tradingTime := strings.Split(output.Data[0].ConfigItemValue, "-") //09:00-17:30
	if len(tradingTime) != 2 {
		return openTime, closeTime, errors.New("QryTradingTime result is invalid")
	}

	openTime, err = time.Parse("15:04", tradingTime[0])
	if err != nil {
		return openTime, closeTime, fmt.Errorf("QryTradingTime error:%v", err.Error())
	}
	closeTime, err = time.Parse("15:04", tradingTime[1])
	if err != nil {
		return openTime, closeTime, fmt.Errorf("QryTradingTime error:%v", err.Error())
	}

	log.Printf("QryTradingTime: %s - %s", openTime.Format(time.TimeOnly), closeTime.Format(time.TimeOnly))
	return openTime, closeTime, nil
}

func qrySymbolsT0(dataQryUrl string) ([]SymbolT0, error) {
	var qryData map[string][]Symbol
	var data = []SymbolT0{}
	err := QryTableData(dataQryUrl, "securities", "ALL_MAP_DATA", &qryData)
	if err != nil {
		return nil, err
	}

	symbolsSet := map[string]bool{}
	for _, v := range qryData {
		for _, v1 := range v {
			if v1.OverSoldIndicator == "Y" && !symbolsSet[v1.Symbol] {
				symbolsSet[v1.Symbol] = true
				data = append(data, SymbolT0{Symbol: v1.Symbol, SymbolName: v1.SymbolName})
			}
		}
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Symbol < data[j].Symbol
	})
	return data, nil
}

func syncSymbolsT0(dataQryUrl string) error {
	return SyncCachedData("securities", &symbolsT0, func() ([]SymbolT0, error) {
		return qrySymbolsT0(dataQryUrl)
	})
}

func GetSymbolsT0(dataQryUrl string) ([]SymbolT0, error) {
	err := syncSymbolsT0(dataQryUrl)
	if err != nil {
		return nil, err
	}
	return symbolsT0.Data, nil
}