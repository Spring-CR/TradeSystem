package citic_qfii

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

//const URL_REALTIME = "http://q3inner-dev.gf.com.cn/v1/realtime"
//const URL_REALTIME = "http://q3inner.gf.com.cn/v1/realtime" // 生产环境
//const URL_REALTIME = "http://10.51.135.40:8080/v1/realtime" // 交易所仿真

type QuoteResponse struct {
	Error struct {
		Code int    `json:"code"`
		Desc string `json:"desc"`
	} `json:"error"`
	Quote struct {
		Data struct {
			Now   int `json:"now"`
			Stock struct {
				Position []struct {
					Bid int `json:"bid"`
					//BidSize int `json:"bid_size"`
					Ask int `json:"ask"`
					//AskSize int `json:"ask_size"`
				} `json:"position"`
			} `json:"stock"`
		} `json:"data"`
	} `json:"quote"`
}

func reqQuoteInfo(url string, code string, exchange int) (quoteResp QuoteResponse, err error) {
	cli := &http.Client{}
	jsonBody := []byte(fmt.Sprintf(`{"id":{"code": "%s", "exchange": %d}}`, code, exchange))
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Add("X-DeviceType", "0x3")
	req.Header.Add("X-SoftwareVer", "0.1")
	resp, err := cli.Do(req)
	if err != nil {
		return quoteResp, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &quoteResp)
	return quoteResp, err
}

// 计算保护限价  side假设1买 2卖
func calcProtectPrice(url string, code string, exchange int, side string) (protectPrice int, err error) {
	quoteResp, err := reqQuoteInfo(url, code, exchange)
	if err != nil || quoteResp.Error.Code != 0 {
		return protectPrice, errors.New(fmt.Sprintln(err, quoteResp.Error.Code, quoteResp.Error.Desc))
	}
	fmt.Printf("quoteResp:%+v\n", quoteResp)

	quotePrices := quoteResp.Quote.Data.Stock.Position
	for i := len(quotePrices); i > 0; i-- {
		if side == "1" {
			protectPrice = quotePrices[i-1].Ask
		} else {
			protectPrice = quotePrices[i-1].Bid
		}
		if protectPrice > 0 {
			fmt.Printf("获取第%d档位行情价格=%d\n", i, quotePrices[i-1].Ask)
			return protectPrice, nil
		}
	}
	if quoteResp.Quote.Data.Now > 0 {
		return quoteResp.Quote.Data.Now, nil
	}
	return protectPrice, errors.New("获取最新行情价格失败！")
}

// 计算保护限价  side假设1买 2卖
func calcShProtectPrice(url string, code string, side string) (protectPrice int, err error) {
	return calcProtectPrice(url, code, 101, side)
}
