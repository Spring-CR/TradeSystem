package fix

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"rhino-common/domain_error"
	"rhino-common/utils/fixutil"
	"rhino-core/order_domain"
	"rhino-core/types"
	"rhino-engine/api/api_adapter"
	"rhino-engine/api/fix_api_adapter"
	"strconv"
	"strings"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/store/file"
	"github.com/quickfixgo/tag"
	"golang.org/x/time/rate"
)

type FixApi struct {
	fixApiAdapter fix_api_adapter.FixApiAdapter
	acceptor      *quickfix.Acceptor
	apiAdapter    api_adapter.APIAdapter
	engine        *order_domain.OrderEngine
	systemCode    string
	businessCode  string
	sessionIDs    []quickfix.SessionID
	rateLimitMap  map[quickfix.SessionID]*rate.Limiter
}

func NewFixApi(fixApiAdapter fix_api_adapter.FixApiAdapter, apiAdapter api_adapter.APIAdapter) *FixApi {
	//systemCode, businessCode := engine.GetSystemAndBusinessCodes()
	inst := &FixApi{fixApiAdapter: fixApiAdapter, apiAdapter: apiAdapter}
	
	return inst
}

func (a *FixApi) Start(engine *order_domain.OrderEngine) {

	a.initAcceptor()

	a.engine = engine
	a.systemCode, a.businessCode = engine.GetSystemAndBusinessCodes()
}

func (a *FixApi) initAcceptor() {
	// 启动fix acceptor
	configFileContent, err := os.ReadFile(a.fixApiAdapter.GetConfigPath())
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("fail to read fixserver config from %s\n", a.fixApiAdapter.GetConfigPath()))
	}
	appSettings, err := quickfix.ParseSettings(bytes.NewReader(configFileContent))
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("fail to parse fixserver setting from %s\n", a.fixApiAdapter.GetConfigPath()))
	}
	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("fail to create fileLogFactory from %s\n", a.fixApiAdapter.GetConfigPath()))
	}

	acceptor, err := quickfix.NewAcceptor(a, file.NewStoreFactory(appSettings), appSettings, fileLogFactory)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, fmt.Sprintf("fail to create fix acceptor from %s\n", a.fixApiAdapter.GetConfigPath()))
	}

	sessionSettings := appSettings.SessionSettings()
	for sessionID := range sessionSettings {
		log.Printf("detect sessionID: %v\n", sessionID)
		a.sessionIDs = append(a.sessionIDs, sessionID)
	}

	a.rateLimitMap = a.createLimitors(sessionSettings)

	log.Println("create acceptor!")

	// err = acceptor.Start()
	// if err != nil {
	// 	domain_error.ProcessSevereError(true, 0, nil, err, fmt.Sprintf("unable to start FIX acceptor from %s\n", a.fixApiAdapter.GetConfigPath()))
	// }
	a.acceptor = acceptor

	a.schedule()
}

func (a *FixApi) OnCreate(sessionID quickfix.SessionID) {

}

// OnLogon notification of a session successfully logging on.
func (a *FixApi) OnLogon(sessionID quickfix.SessionID) {

}

// OnLogout notification of a session logging off or disconnecting.
func (a *FixApi) OnLogout(sessionID quickfix.SessionID) {

}

// ToAdmin notification of admin message being sent to target.
func (a *FixApi) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
	log.Printf("ToAdmin: %s", fixutil.ConvertFIXDataToString(message.Bytes()))
}

// ToApp notification of app message being sent to target.
func (a *FixApi) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {
	log.Printf("ToApp: %s", fixutil.ConvertFIXDataToString(message.Bytes()))
	return nil
}

// FromAdmin notification of admin message being received from target.
func (a *FixApi) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	logText := fixutil.ConvertFIXDataToString(message.Bytes())
	msgType, _ := message.MsgType()
	if msgType == string(enum.MsgType_LOGON) {
		logText = desensitizeFIXLog(logText)
	}
	log.Printf("FromAdmin: %s", logText)
	if !message.IsMsgTypeOf(string(enum.MsgType_LOGON)) {
		return nil
	}

	username, _ := message.Body.GetString(tag.Username)
	password, _ := message.Body.GetString(tag.Password)

	// 进行身份验证
	if !a.fixApiAdapter.LoginValidate(username, password, sessionID) {
		//return quickfix.NewMessageRejectError("Authentication failed", 0, nil)
		return quickfix.RejectLogon{Text: "Authentication failed"}
	}

	return nil
}

func desensitizeFIXLog(fixLog string) string {
    fields := strings.Split(fixLog, "|")
    
    for i, field := range fields {
        if strings.HasPrefix(field, "554=") {
            // 替换密码字段，保留标签但隐藏值
            fields[i] = "554=**"
        }
    }
    
    return strings.Join(fields, "|")
}

// FromApp notification of app message being received from target.
func (a *FixApi) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	log.Printf("FromApp: %s", fixutil.ConvertFIXDataToString(message.Bytes()))

	msgType, _ := message.MsgType()
	switch msgType {
	case string(enum.MsgType_ORDER_SINGLE):
		return a.newOrderSingle(message, sessionID, a.isRateLimitExceeded(sessionID))
	case string(enum.MsgType_ORDER_CANCEL_REQUEST):
		return a.cancelOrder(message, sessionID, a.isRateLimitExceeded(sessionID))
	}

	return nil
}

func (a *FixApi) processDomainError(de *domain_error.Error, message *quickfix.Message) quickfix.MessageRejectError {

	if de == nil {
		return nil
	}
	errCode, err := strconv.Atoi(de.Code)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to get errCode:"+de.Code)
	}

	return quickfix.NewMessageRejectError(de.Msg, errCode, nil)
}

func (a *FixApi) OnTradeResp(tradeResp *types.TradeActionRespReturn) bool {
	message, de := a.fixApiAdapter.ConvertTradeResponseMessage(tradeResp)
	if de != nil {
		domain_error.ProcessSevereError(false, 0, de, nil, "fail to ConvertTradeResponseMessage")
		return false
	}
	if message == nil {
		return false
	}
	err := quickfix.Send(message)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to ConvertTradeResponseMessage")
		// 插入store
		// 不需要，因为quickfix.Send会先存储，后发送消息
		return false
	}
	return true
}
