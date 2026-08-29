package fix

import (
	"log"
	"rhino-common/domain_error"
	"strconv"

	"github.com/quickfixgo/quickfix"
	"golang.org/x/time/rate"
)

func (a *FixApi) createLimitors(sessionSettings map[quickfix.SessionID]*quickfix.SessionSettings) map[quickfix.SessionID]*rate.Limiter {

	rateLimitMap := make(map[quickfix.SessionID]*rate.Limiter)

	for sessionID, sessionSetting := range sessionSettings {
		val, _ := sessionSetting.Setting("RateLimit")
		if val == "" {
			continue
		}
		rateLimit, _ := strconv.ParseFloat(val, 64)
		if rateLimit > 0 {
			rateLimitMap[sessionID] = rate.NewLimiter(rate.Limit(rateLimit), int(rateLimit))
			log.Printf("======> sessionID:%v, rateLimit:%v\n", sessionID, rateLimit)
		}
	}

	return rateLimitMap
}

func (a *FixApi) isRateLimitExceeded(sessionID quickfix.SessionID) (de *domain_error.Error) {

	limiter, ok := a.rateLimitMap[sessionID]
	if !ok {
		return
	}

	if limiter.Allow() {
		return
	}

	return domain_error.Build(domain_error.RATE_LIMIT_EXCEEDED_ERR_CODE, nil, int(limiter.Limit()))
}
