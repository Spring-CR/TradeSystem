package ficc_fut

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"rhino-common/utils/dbutil"
	ficc_fut_posi "rhino-plugins/order_position_plugin/ficc_fut"
	"strconv"
)

var (
	query = `select Account, Account as CounterpartyID, Counterparty, Symbol2, SymbolName, Currency, PlanCode, UltraContractCode, SecurityExchange2 as SecurityExchange, futType as SecurityType, ProductCode, ContractMultiplier, marginRatio as LongMarginRatio, marginRatio as ShortMarginRatio, CAST(f_trade_date AS CHAR) AS ContractBaseDate,
	CASE 
        WHEN side = '1' THEN f_last_shares 
        ELSE -f_last_shares 
    END AS InitNetPosition,
    CASE 
        WHEN side = '1' THEN f_last_shares * ContractMultiplier
        ELSE 0 
    END AS LongPriceCost,
    CASE 
        WHEN side = '2' THEN f_last_shares * ContractMultiplier
        ELSE 0 
    END AS ShortPriceCost,
	f_last_px as lastPrice, commissionType, commissionValue, exchangeRateCNY
from trade_action_resps_extend_ficc_fut where f_last_shares > 0 and exchangeRateCNY > 0 and ContractMultiplier > 0 and f_trade_date > ? and f_channel_code = ? order by f_ord_status_update_time`
)

// findPositionRecordFromDB 从数据库查询持仓记录
// tradeDate: 整型日期（如20260813），会直接用于 f_trade_date > tradeDate 条件
// channelCode: 渠道编码
func (a *FiccFutDataSyncAdapter) findPositionRecordFromDB(db *sql.DB, tradeDate int, channelCode string, exchangeArea string) ([]*ficc_fut_posi.PositionRecord, error) {

	if db == nil {
		return nil, fmt.Errorf("db connection is nil")
	}

	rows, err := db.Query(query, tradeDate, channelCode)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var records []*ficc_fut_posi.PositionRecord
	for rows.Next() {
		// 使用 sql.Null 类型处理可能为 NULL 的字段，转换时自动取零值
		var (
			account            sql.NullInt64
			counterpartyID     sql.NullInt64
			counterparty       sql.NullString
			symbol2            sql.NullString
			symbolName         sql.NullString
			currency           sql.NullString
			planCode           sql.NullString
			ultraContractCode  sql.NullString
			securityExchange   sql.NullString
			securityType       sql.NullString
			productCode        sql.NullString
			contractMultiplier sql.NullFloat64
			longMarginRatio    sql.NullFloat64
			shortMarginRatio   sql.NullFloat64
			contractBaseDate   sql.NullString
			initNetPosition    sql.NullFloat64
			longPriceCost      sql.NullFloat64
			shortPriceCost     sql.NullFloat64
			lastPrice          sql.NullFloat64
			commissionType     sql.NullString
			commissionValue    sql.NullFloat64
			exchangeRateCNY    sql.NullFloat64
		)

		err := rows.Scan(
			&account,
			&counterpartyID,
			&counterparty,
			&symbol2,
			&symbolName,
			&currency,
			&planCode,
			&ultraContractCode,
			&securityExchange,
			&securityType,
			&productCode,
			&contractMultiplier,
			&longMarginRatio,
			&shortMarginRatio,
			&contractBaseDate,
			&initNetPosition,
			&longPriceCost,
			&shortPriceCost,
			&lastPrice,
			&commissionType,
			&commissionValue,
			&exchangeRateCNY,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row failed: %w", err)
		}

		rec := &ficc_fut_posi.PositionRecord{
			Account:            int(account.Int64), // NULL -> 0
			CounterpartyID:     int(counterpartyID.Int64),
			Counterparty:       counterparty.String, // NULL -> ""
			Symbol2:            symbol2.String,
			SymbolName:         symbolName.String,
			Currency:           currency.String,
			PlanCode:           planCode.String,
			UltraContractCode:  ultraContractCode.String,
			SecurityExchange:   securityExchange.String,
			SecurityType:       securityType.String,
			ProductCode:        productCode.String,
			ContractMultiplier: contractMultiplier.Float64, // NULL -> 0.0
			LongMarginRatio:    longMarginRatio.Float64,
			ShortMarginRatio:   shortMarginRatio.Float64,
			ExchangeArea:       exchangeArea,
			ContractBaseDate:   contractBaseDate.String,
			InitNetPosition:    initNetPosition.Float64,
			LongPriceCost:      longPriceCost.Float64,
			ShortPriceCost:     shortPriceCost.Float64,
		}

		tmpLongPriceCost := rec.LongPriceCost
		tmpShortPriceCost := rec.ShortPriceCost

		price := lastPrice.Float64
		rec.LongPriceCost *= price
		rec.ShortPriceCost *= price

		if rec.LongPriceCost > 0 {
			switch commissionType.String {
			case "FEE_BY_RATE":
				price *= 1 + commissionValue.Float64
			case "FEE_BY_PER_SHARE":
				price += commissionValue.Float64
			}
		}
		if rec.ShortPriceCost > 0 {
			switch commissionType.String {
			case "FEE_BY_RATE":
				price *= 1 - commissionValue.Float64
			case "FEE_BY_PER_SHARE":
				price -= commissionValue.Float64
			}
		}

		if price <= 0 {
			price = lastPrice.Float64
		}

		rec.LongPriceWithFeeCost = tmpLongPriceCost * price
		rec.ShortPriceWithFeeCost = tmpShortPriceCost * price

		rec.LongPriceCNYWithFeeCost = rec.LongPriceWithFeeCost * exchangeRateCNY.Float64
		rec.ShortPriceCNYWithFeeCost = rec.ShortPriceWithFeeCost * exchangeRateCNY.Float64

		records = append(records, rec)
	}
	if err = rows.Err(); err != nil {
		log.Printf("fail to findPositionRecordFromDB, sql:%s, tradeDate:%v, channelCode:%v\n", query, tradeDate, channelCode)
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return records, nil
}

/*
整体逻辑：
1、同方向，代表持续开仓，直接净持仓、持仓成本相加
2、反方向，首先净持仓相加；如果净持仓方向和basePosition的净持仓方向相同，说明平仓未平完，持仓均价是不变的，所以basePosition的持仓成本等于原均价*最新净持仓的绝对值；

	如果净持仓方向和basePosition的净持仓方向相反，说明basePosition已经平仓完，并且反向开仓了，这时持仓均价需要等于tradePosition的持仓均价，basePosition的持仓成本等于tradePosition持仓均价*最新净持仓的绝对值；
*/
func (a *FiccFutDataSyncAdapter) updateByTradePosition(basePosition, tradePosition *ficc_fut_posi.PositionRecord) {

	if tradePosition.ContractBaseDate > basePosition.ContractBaseDate {
		basePosition.ContractBaseDate = tradePosition.ContractBaseDate
	}

	// 当basePosition、tradePosition净持仓方向相同时，NetPosition、LongPriceCost、ShortPriceCost可以直接相加
	if basePosition.InitNetPosition >= 0 && tradePosition.InitNetPosition > 0 || basePosition.InitNetPosition <= 0 && tradePosition.InitNetPosition < 0 {

		basePosition.InitNetPosition += tradePosition.InitNetPosition
		basePosition.LongPriceCost += tradePosition.LongPriceCost
		basePosition.ShortPriceCost += tradePosition.ShortPriceCost
		basePosition.LongPriceWithFeeCost += tradePosition.LongPriceWithFeeCost
		basePosition.ShortPriceWithFeeCost += tradePosition.ShortPriceWithFeeCost
		basePosition.LongPriceCNYWithFeeCost += tradePosition.LongPriceCNYWithFeeCost
		basePosition.ShortPriceCNYWithFeeCost += tradePosition.ShortPriceCNYWithFeeCost

		return
	}

	// 无论如何，净持仓都是：basePosition.InitNetPosition + tradePosition.InitNetPosition
	newNetPosition := basePosition.InitNetPosition + tradePosition.InitNetPosition

	// 当newNetPosition == 0，均价、成本等，全部归零
	if newNetPosition == 0 {

		basePosition.InitNetPosition = 0
		basePosition.LongPriceCost = 0
		basePosition.ShortPriceCost = 0
		basePosition.LongPriceWithFeeCost = 0
		basePosition.ShortPriceWithFeeCost = 0
		basePosition.LongPriceCNYWithFeeCost = 0
		basePosition.ShortPriceCNYWithFeeCost = 0

		return
	}

	// LongPriceCost、ShortPriceCost只能有一个有值（需要去验证）
	baseAvgPrice := (basePosition.LongPriceCost + basePosition.ShortPriceCost) / math.Abs(basePosition.InitNetPosition)
	tradeAvgPrice := (tradePosition.LongPriceCost + tradePosition.ShortPriceCost) / math.Abs(tradePosition.InitNetPosition)

	baseAvgPriceWithFee := (basePosition.LongPriceWithFeeCost + basePosition.ShortPriceWithFeeCost) / math.Abs(basePosition.InitNetPosition)
	tradeAvgPriceWithFee := (tradePosition.LongPriceWithFeeCost + tradePosition.ShortPriceWithFeeCost) / math.Abs(tradePosition.InitNetPosition)

	baseAvgPriceCNYWithFee := (basePosition.LongPriceCNYWithFeeCost + basePosition.ShortPriceCNYWithFeeCost) / math.Abs(basePosition.InitNetPosition)
	tradeAvgPriceCNYWithFee := (tradePosition.LongPriceCNYWithFeeCost + tradePosition.ShortPriceCNYWithFeeCost) / math.Abs(tradePosition.InitNetPosition)

	// newNetPosition > 0
	if newNetPosition > 0 {

		if basePosition.InitNetPosition >= 0 {

			// 同向，部分平仓引起持仓减少，这时均价取basePosition的原均价维持不变，成本用原均价*新的净持仓绝对值
			basePosition.InitNetPosition = newNetPosition
			basePosition.LongPriceCost = baseAvgPrice * newNetPosition
			basePosition.ShortPriceCost = 0
			basePosition.LongPriceWithFeeCost = baseAvgPriceWithFee * newNetPosition
			basePosition.ShortPriceWithFeeCost = 0
			basePosition.LongPriceCNYWithFeeCost = baseAvgPriceCNYWithFee * newNetPosition
			basePosition.ShortPriceCNYWithFeeCost = 0

		} else {

			// 反向，全部平仓引且反向开仓，这时均价取tradePosition的原均价维持不变，成本用原均价*新的净持仓绝对值
			basePosition.InitNetPosition = newNetPosition
			basePosition.LongPriceCost = tradeAvgPrice * newNetPosition
			basePosition.ShortPriceCost = 0
			basePosition.LongPriceWithFeeCost = tradeAvgPriceWithFee * newNetPosition
			basePosition.ShortPriceWithFeeCost = 0
			basePosition.LongPriceCNYWithFeeCost = tradeAvgPriceCNYWithFee * newNetPosition
			basePosition.ShortPriceCNYWithFeeCost = 0

		}

		return
	}

	// newNetPosition < 0
	if basePosition.InitNetPosition <= 0 {

		// 同向
		basePosition.InitNetPosition = newNetPosition
		basePosition.LongPriceCost = 0
		basePosition.ShortPriceCost = baseAvgPrice * -newNetPosition
		basePosition.LongPriceWithFeeCost = 0
		basePosition.ShortPriceWithFeeCost = baseAvgPriceWithFee * -newNetPosition
		basePosition.LongPriceCNYWithFeeCost = 0
		basePosition.ShortPriceCNYWithFeeCost = baseAvgPriceCNYWithFee * -newNetPosition

	} else {

		// 反向
		basePosition.InitNetPosition = newNetPosition
		basePosition.LongPriceCost = 0
		basePosition.ShortPriceCost = tradeAvgPrice * -newNetPosition
		basePosition.LongPriceWithFeeCost = 0
		basePosition.ShortPriceWithFeeCost = tradeAvgPriceWithFee * -newNetPosition
		basePosition.LongPriceCNYWithFeeCost = 0
		basePosition.ShortPriceCNYWithFeeCost = tradeAvgPriceCNYWithFee * -newNetPosition
	}
}

func (a *FiccFutDataSyncAdapter) updateByTradeActionResps(tableConfig *dbutil.TableConfig, positions []*ficc_fut_posi.PositionRecord) ([]*ficc_fut_posi.PositionRecord, error) {
	// 重组positionMap
	m := make(map[string]*ficc_fut_posi.PositionRecord)
	maxTradeDate := ""
	for _, position := range positions {
		if position.ContractBaseDate > maxTradeDate {
			maxTradeDate = position.ContractBaseDate
		}
		key := fmt.Sprintf("%v-%v", position.CounterpartyID, position.Symbol2)
		m[key] = position
	}

	sinceTradeDate, err1 := strconv.Atoi(maxTradeDate)
	if err1 != nil {
		return nil, err1
	}

	var channelCode string
	var exchangeArea string
	switch tableConfig.TableAlias {
	case "PositionBaseDms":
		channelCode = "olts-fut"
		exchangeArea = "dms"
	case "PositionBaseOvs":
		channelCode = "stars-fut"
		exchangeArea = "ovs"
	}

	// 从成交回报合成PositionRecord
	tradePositionsTmp, err1 := a.findPositionRecordFromDB(a.dataSyncConfig.GetCentralDB(), sinceTradeDate, channelCode, exchangeArea)
	if err1 != nil {
		return nil, err1
	}

	var tradePositions []*ficc_fut_posi.PositionRecord
	for _, position := range tradePositionsTmp {
		key := fmt.Sprintf("%v-%v", position.CounterpartyID, position.Symbol2)
		if _, ok := m[key]; !ok {
			positions = append(positions, position)
			m[key] = position
		} else {
			tradePositions = append(tradePositions, position)
		}
	}

	for _, tradePosition := range tradePositions {
		key := fmt.Sprintf("%v-%v", tradePosition.CounterpartyID, tradePosition.Symbol2)
		basePosition, ok := m[key]
		if !ok {
			continue
		}
		a.updateByTradePosition(basePosition, tradePosition)
	}

	return positions, nil
}
