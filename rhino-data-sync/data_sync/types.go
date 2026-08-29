package data_sync

import (
	"crypto/md5"
	"fmt"
	"rhino-common/utils/timeutil"
	"strings"
	"time"
)

type DspOption struct {
	BaseUrl   string `json:"baseUrl"`
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
	QueryId   string `json:"queryId"`
}

type CsvQryResp struct {
	Data string `json:"data"`
}

func generateMD5(str string) string {
	hashInBytes := md5.Sum([]byte(str))
	// 使用 %x 格式动词将字节数组直接格式化为十六进制字符串
	return strings.ToUpper(fmt.Sprintf("%x", hashInBytes))
}

func (o *DspOption) GenerateSecret() (sign string, timestamp int64) {
	timestamp = timeutil.ConvertTimeToMilliseconds(time.Now())
	sign = fmt.Sprintf("%s%d", o.AppSecret, timestamp)
	sign = generateMD5(sign)
	return
}
