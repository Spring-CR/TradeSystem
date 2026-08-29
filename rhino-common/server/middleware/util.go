package middleware

import (
	"fmt"
	"log"
	"net/http"
	domainErr "rhino-common/domain_error"
	"rhino-common/utils/byteutils"
	"strings"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type ErrorResult struct {
	Code    int
	Message string
}

const (
	API_HEADER_Name          = `X-Api-Name`
	API_HEADER_ERROR_CODE    = `X-Api-Err-Code`
	API_HEADER_USER_ID       = "X-Api-User-ID"
	API_HEADER_PRIVATE_TOKEN = "x-Api-Private-Token"
)

func ProcessDomainError(err *domainErr.Error, c *gin.Context) (errHappen bool) {
	if err != nil {
		c.Header(API_HEADER_ERROR_CODE, err.Code)
		c.JSON(http.StatusInternalServerError, err)
		errHappen = true
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		errHappen = false
	}
	return errHappen
}

func ProcessGenericError(err error, c *gin.Context) (errHappen bool) {
	if err != nil {
		c.Header(API_HEADER_ERROR_CODE, domainErr.GENERIC_ERR_CODE)
		c.JSON(http.StatusInternalServerError, domainErr.Build(domainErr.GENERIC_ERR_CODE, err))
		errHappen = true
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		errHappen = false
	}
	return errHappen
}

func BindInputOption(c *gin.Context, option interface{}) (sucess bool) {
	sucess = true
	err := c.Bind(option)
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, &ErrorResult{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Bad Request\n%s\n", err.Error())})
		c.AbortWithStatus(http.StatusInternalServerError)
		sucess = false
	}
	return sucess
}

func ResponseJson(c *gin.Context, data interface{}) {
	c.Header("Content-Type", "application/json")
	c.JSON(200, data)
}

func ResponseOk(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.String(200, `{"status": "ok"}`)
}

func ResponseJsonStrict(c *gin.Context, data interface{}) {
	c.Header("Content-Type", "application/json")
	jsonData, _ := json.Marshal(data)
	//c.JSON(200, data)
	log.Printf("jsonData:\n%s\n", jsonData)
	c.String(200, byteutils.GetZeroCopyString(jsonData))
}

func ResponseBadRequest(c *gin.Context, errMsg string) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusBadRequest, &ErrorResult{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("Bad Request\n%s\n", errMsg)})
}

func SplitNameListString(nameListStr string) (names []string, newNameListStr string) {
	nameListStr = strings.TrimSpace(nameListStr)
	nameListStr = strings.Replace(nameListStr, ",", " ", -1)
	nameListStr = strings.Replace(nameListStr, ";", " ", -1)
	strs := strings.Split(nameListStr, " ")
	n := len(strs)
	for i := 0; i < n; i++ {
		str := strings.TrimSpace(strs[i])
		if len(str) > 0 {
			names = append(names, str)
		}
	}
	if len(names) > 0 {
		newNameListStr = strings.Join(names, ",")
	}
	return names, newNameListStr
}
