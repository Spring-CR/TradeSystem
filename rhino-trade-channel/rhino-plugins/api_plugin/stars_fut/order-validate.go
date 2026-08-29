package stars_fut

import (
	"rhino-common/domain_error"
	"rhino-core/schema"
)

func (a *StarFurAPIAdapter) RefineAndValidate(tradeOrder *schema.TradeOrder, trade bool) *domain_error.Error {
	return nil
}