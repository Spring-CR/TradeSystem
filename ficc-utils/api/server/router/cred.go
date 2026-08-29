package router

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"ficc-utils/api/api_const"
	"ficc-utils/common/domain_error"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type CounterpartyPhoneNum struct {
	CounterpartyID int
	PhoneNum       string
}

type CredProvider struct {
	credLock           *sync.RWMutex
	dataQryUrl         string
	queryCycle         time.Duration
	local              *time.Location
	lastQueryResult    map[string][]*CounterpartyPhoneNum // 记录上次的查询结果
	tokenMap           map[string]map[int]bool            // key = phone + yyyyMMdd的md5编码，value = 交易对手id
	lastQueryRawResult []byte                             // 记录上次的查询结果，http://10.51.136.21:6093/api/v1/data_qry?collection=CounterpartyPhoneNum&key=ALL_MAP_DATA的查询结果
	lastQueryDate      string                             // 记录上一次的查询日期
}

func NewCredProvider(dataQryUrl string, queryCycle time.Duration, local *time.Location) *CredProvider {
	inst := &CredProvider{
		credLock:        &sync.RWMutex{},
		dataQryUrl:      dataQryUrl,
		queryCycle:      queryCycle,
		local:           local,
		lastQueryResult: map[string][]*CounterpartyPhoneNum{},
		tokenMap:        map[string]map[int]bool{},
	}
	inst.doRefresh()
	return inst
}

func (p *CredProvider) CredFunc(c *gin.Context) (ok bool) {

	tokens := c.Request.Header["X-Api-Token"]
	if len(tokens) == 0 {
		return
	}

	token := tokens[0]

	p.credLock.RLock()
	cptyMap, ok := p.tokenMap[token]
	p.credLock.RUnlock()

	if !ok {
		return
	}

	_account, ok := c.GetQuery(api_const.ParamAccount)
	if ok {
		account, err := strconv.Atoi(_account)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "fail to strconv the account value "+_account)
			return
		}
		return cptyMap[account]
	}

	ok = true

	return
}

func (p *CredProvider) doRefresh() {
	go func() {
		for {
			p.refreshCounterpartyPhoneNum()
			time.Sleep(p.queryCycle)
		}
	}()
}

func (p *CredProvider) refreshCounterpartyPhoneNum() {
	resp, err := http.Get(p.dataQryUrl)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to invoke "+p.dataQryUrl)
		return
	}
	body := resp.Body
	if body == nil {
		return
	}
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to read result after invoke "+p.dataQryUrl)
		return
	}

	date := time.Now().In(p.local).Format(time.DateOnly)

	// 结果一样，啥也不做
	if bytes.Equal(data, p.lastQueryRawResult) && date == p.lastQueryDate {
		return
	}

	err = json.Unmarshal(data, &p.lastQueryResult)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to unmarshal result after invoke "+p.dataQryUrl)
		return
	}

	tokenMap := map[string]map[int]bool{}

	for _, valList := range p.lastQueryResult {
		for _, val := range valList {
			counterpartyID := val.CounterpartyID
			phoneNum := val.PhoneNum

			log.Printf("===> counterpartyID:%d, phoneNum:%s\n", counterpartyID, phoneNum)

			k := generateMD5(fmt.Sprintf("%s%s", phoneNum, date))
			v := tokenMap[k]
			if v == nil {
				v = map[int]bool{}
			}
			v[counterpartyID] = true

			tokenMap[k] = v
		}
	}

	p.credLock.Lock()
	defer p.credLock.Unlock()

	// 重置tokenMap
	p.tokenMap = tokenMap
	// 更新lastQueryRawResult
	p.lastQueryRawResult = data
	p.lastQueryDate = date
}

func (p *CredProvider) getTokenKey(phoneNum string) (string, string, map[int]bool, bool) {
	date := time.Now().In(p.local).Format(time.DateOnly)
	k := generateMD5(fmt.Sprintf("%s%s", phoneNum, date))
	v, ok := p.tokenMap[k]
	return k, date, v, ok
}

func generateMD5(str string) string {
	hashInBytes := md5.Sum([]byte(str))
	// 使用 %x 格式动词将字节数组直接格式化为十六进制字符串
	return strings.ToUpper(fmt.Sprintf("%x", hashInBytes))
}
