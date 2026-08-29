package ficc

import (
	"errors"
	"fmt"
	"math"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"strings"
	"time"
)

// 计息基准类型
type DayCountBasis string

const (
	ActAct     DayCountBasis = "A/A"
	ActActISDA DayCountBasis = "A/A(ISDA)"
	Act365F    DayCountBasis = "A/365F"
	Act365     DayCountBasis = "A/365"
	Act360     DayCountBasis = "A/360"
	Thirty360  DayCountBasis = "30/360"
)

// 息票类型
type CouponType string

const (
	FixedRate  CouponType = "1" // 固定利率
	Floating   CouponType = "2" // 浮动利率
	ZeroCoupon CouponType = "3" // 零息
)

// 市场类型
type MarketType string

const (
	Interbank MarketType = "yhj" // 银行间市场
	Exchange  MarketType = "jys" // 交易所市场
)

// 债券详情
type BondDetail struct {
	FaceValue       float64       `json:"face_value"`       // 面值
	CouponType      CouponType    `json:"coupon_type"`      // 息票类型
	CouponRate      float64       `json:"coupon_rate"`      // 票面利率（%）
	IssuePrice      float64       `json:"issue_price"`      // 发行价格
	StartDate       time.Time     `json:"start_date"`       // 起息日
	EndDate         time.Time     `json:"end_date"`         // 到期日
	PaymentFreq     int           `json:"payment_freq"`     // 年付息次数
	DayCountBasis   DayCountBasis `json:"day_count_basis"`  // 计息基准
	Market          MarketType    `json:"market"`           // 市场类型
	BondType        string        `json:"bond_type"`        // 债券类型
	CleanPrice      float64       `json:"clean_price"`      // 净价（输入）
	DirtyPrice      float64       `json:"dirty_price"`      // 全价（输出）
	AccruedInterest float64       `json:"accrued_interest"` // 应计利息（输出）
}

// 债券定价引擎
type BondPricingEngine struct {
	SettleDate time.Time
	Detail     BondDetail
	// 内部计算中间变量
	ts float64
	t  float64
}

// 解析数据库格式的日期时间字符串
func parseDBDateTime(dateStr string) (time.Time, error) {
	dateTrim := strings.TrimSpace(dateStr)
	// 尝试解析带时分秒的格式
	parseLayouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	var t time.Time
	var err error
	for _, layout := range parseLayouts {
		t, err = time.Parse(layout, dateTrim)
		if err == nil {
			return t.UTC(), nil // 统一转换为UTC时区
		}
	}
	return time.Time{}, fmt.Errorf("不支持的日期格式：%s（支持格式：2023-01-01 00:00:00 或 2023-01-01）", dateStr)
}

// 转换计息基准字符串为枚举类型
func convertDayCountBasis(dayCountBasisParam string) (DayCountBasis, error) {
	switch dayCountBasisParam {
	case "ACT/ACT(ISMA)":
		return ActAct, nil
	case "":
		return ActAct, nil
	case "Actual/360":
		return Act360, nil
	case "ACT/365NoLeap":
		return Act365F, nil
	case "ACT/365":
		return Act365, nil
	}

	dayCountBasisParam = strings.ReplaceAll(dayCountBasisParam, " ", "")
	dayCountBasisParam = strings.ReplaceAll(dayCountBasisParam, "实际", "ACT")
	dayCountBasisParam = strings.ReplaceAll(dayCountBasisParam, "固定", "FIXED")
	dayCountBasisParam = strings.ToUpper(dayCountBasisParam)
	dayCountBasisParam = strings.ReplaceAll(dayCountBasisParam, "ACTUAL", "a")
	dayCountBasisParam = strings.ReplaceAll(dayCountBasisParam, "ACT", "a")
	dayCountBasisParam = strings.ReplaceAll(dayCountBasisParam, "A/", "ACT/")
	dayCountBasisParam = strings.ReplaceAll(dayCountBasisParam, "/A", "/ACT")
	dayCountBasisParam = strings.ReplaceAll(dayCountBasisParam, "a", "ACT")
	dayCountBasisParam = strings.ReplaceAll(dayCountBasisParam, "/", "")

	switch dayCountBasisParam {
	case "ACTACT", "ACTACT(BOND)", "ACTACT(ISMA)", "ACTACT(AFB)":
		return ActAct, nil
	case "ACTACT(ISDA)", "ACTACT(STARTDATE)":
		return ActActISDA, nil
	case "ACT365F", "ACT365(FIXED)":
		return Act365F, nil
	case "ACT365":
		return Act365, nil
	case "ACT360":
		return Act360, nil
	case "30360", "30360(BONDBASIS)", "30E360(EUROBONDBASIS)":
		return Thirty360, nil
	default:
		return ActAct, nil
	}
}

// 四舍五入
func roundDecimal(value float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(value*pow) / pow
}

// 判断闰年
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// 获取N个月后的日期(逻辑已优化)
func getNMonthDate(date time.Time, months int) time.Time {
	// 计算新的年月
	year := date.Year()
	month := int(date.Month()) + months

	newYear := year + (month-1)/12
	newMonth := (month-1)%12 + 1

	// 处理天数，防止超出当月天数
	day := date.Day()

	// 获取目标月份的天数
	daysInMonth := time.Date(newYear, time.Month(newMonth+1), 0, 0, 0, 0, 0, time.UTC).Day()

	if day > daysInMonth {
		day = daysInMonth
	}

	return time.Date(newYear, time.Month(newMonth), day, 0, 0, 0, 0, time.UTC)
}

// 计算结算日在当前付息周期内的已计息天数（区分银行间 / 交易所规则）
func identifySpecialDay(start, end time.Time, market MarketType) bool {
	hasSpecialDay := false
	if isLeapYear(end.Year()) {
		feb29EndYear := time.Date(end.Year(), 2, 29, 0, 0, 0, 0, start.Location())
		if start.Before(feb29EndYear) || start.Equal(feb29EndYear) { // start <= 2月29日
			if feb29EndYear.Before(end) { // 2月29日 < end
				hasSpecialDay = true
			}
		}
		if market == Exchange && feb29EndYear.Equal(end) {
			hasSpecialDay = true
		}
	}

	if !isLeapYear(end.Year()) && isLeapYear(start.Year()) {
		feb29StartYear := time.Date(start.Year(), 2, 29, 0, 0, 0, 0, start.Location())
		if start.Before(feb29StartYear) || start.Equal(feb29StartYear) { // start <= 2月29日
			if feb29StartYear.Before(end) { // 2月29日 < end
				hasSpecialDay = true
			}
		}
		if market == Exchange && feb29StartYear.Equal(end) {
			hasSpecialDay = true
		}
	}

	return hasSpecialDay
}

// 获取N年后的日期
func getNYearDate(date time.Time, years int) time.Time {
	return getNMonthDate(date, years*12)
}

// 根据起息日和到期日获取债券全部的利息“日期流”
func getDateListByStartAndEnd(startDate, endDate time.Time, freq int) []time.Time {
	dateList := []time.Time{}
	currentDate := startDate
	n := 12 / freq
	i := 0
	for currentDate.Before(endDate) {
		nextDate := getNMonthDate(startDate, i+n)
		dateList = append(dateList, nextDate)
		i += n
		currentDate = nextDate
	}
	return dateList
}

// 计算结算日在当前付息周期内的已计息天数
func calculateT(settleDate, couponPeriodStart time.Time, market MarketType) int {
	var t int
	// 计算日期差的天数（小时/24转换，避免纳秒精度问题）
	daysDiff := int(math.Round(settleDate.Sub(couponPeriodStart).Hours() / 24))

	if market == Interbank {
		t = daysDiff
	} else {
		t = daysDiff + 1
		if identifySpecialDay(couponPeriodStart, settleDate, market) {
			t -= 1
		}
	}
	return t
}

// 计算结算日在当前付息周期内的已计息天数（区分银行间 / 交易所规则）
func calculateTy(settleDate, startDate time.Time, mode MarketType, basis DayCountBasis, bondType string) float64 {
	if basis == ActAct || basis == ActActISDA || mode == Exchange || bondType == "IBNCD" {
		yearBeg := startDate
		n := 0
		var yearEnd time.Time

		for i := 0; i < 1000; i++ {
			yearEnd = getNYearDate(yearBeg, 1)
			// 对齐 Python 的 year_range[0] <= settle_date < year_range[1]
			if (settleDate.Equal(yearBeg) || settleDate.After(yearBeg)) && settleDate.Before(yearEnd) {
				break
			}
			if n > 100 {
				break
			}
			yearBeg = yearEnd
			n += 1
		}
		return yearEnd.Sub(yearBeg).Hours() / 24
	}

	if (basis == Act365F || basis == Act365) && mode == Interbank {
		return 365.0
	}

	if (basis == Act360 || basis == Thirty360) && mode == Interbank {
		return 360.0
	}

	return 365.0
}

// 按计息基准计算 “日计数分数（DayCountFraction）”（已计息天数 / 付息周期总天数），是应计利息计算的核心
func calcDaycount(settleDate, couponPeriodStart time.Time, periodDaycount float64, engine *BondPricingEngine) (float64, float64, error) {
	t := float64(calculateT(settleDate, couponPeriodStart, engine.Detail.Market))
	f := engine.Detail.PaymentFreq
	basis := engine.Detail.DayCountBasis
	mode := engine.Detail.Market

	if f < 1 {
		f = 1
	}

	// 2. 交易所规则
	if mode == Exchange {
		fraction := t / (365.0 / float64(f))
		return fraction, t, nil
	}

	// 3. 银行间规则
	var fraction float64
	switch basis {
	case ActActISDA:
		settleYearLeap := isLeapYear(settleDate.Year())
		startYearLeap := isLeapYear(couponPeriodStart.Year())

		if (settleYearLeap && !startYearLeap) || (!settleYearLeap && startYearLeap) {
			fraction = t / (365.0 / float64(f))
		} else {
			yearStart := time.Date(settleDate.Year(), 1, 1, 0, 0, 0, 0, settleDate.Location())
			n1 := float64(yearStart.Sub(couponPeriodStart).Hours() / 24)
			n2 := float64(settleDate.Sub(yearStart).Hours() / 24)

			b1 := 365.0 / float64(f)
			if startYearLeap {
				b1 = 366.0 / float64(f)
			}
			b2 := 365.0 / float64(f)
			if settleYearLeap {
				b2 = 366.0 / float64(f)
			}

			t = n1 + n2
			fraction = n1/b1 + n2/b2
		}

	case ActAct:
		if periodDaycount == 0 {
			return 0, 0, errors.New("periodDaycount 不能为0（ActAct基准）")
		}
		fraction = t / periodDaycount

	case Act365:
		fraction = t / (365.0 / float64(f))

	case Act365F:
		if identifySpecialDay(couponPeriodStart, settleDate, mode) {
			t -= 1
		}
		fraction = t / (365.0 / float64(f))

	case Act360:
		fraction = t / (360.0 / float64(f))

	case Thirty360:
		yearDiff := float64(settleDate.Year() - couponPeriodStart.Year())
		monthDiff := float64(int(settleDate.Month()) - int(couponPeriodStart.Month()))
		dayDiff := float64(settleDate.Day() - couponPeriodStart.Day())
		t = 360*yearDiff + 30*monthDiff + dayDiff
		fraction = t / (360.0 / float64(f))

	default:
		return 0, 0, fmt.Errorf("不支持的计息基准: %s", basis)
	}

	return fraction, t, nil
}

// 计算当前付息周期的总天数、付息周期区间、剩余付息次数、全付息日期列表
func calculateTs(settleDate, startDate, endDate time.Time, engine *BondPricingEngine) (float64, [2]time.Time, int, []time.Time, error) {
	f := engine.Detail.PaymentFreq
	basis := engine.Detail.DayCountBasis
	mode := engine.Detail.Market
	bondType := engine.Detail.BondType

	paymentDates := getDateListByStartAndEnd(startDate, endDate, f)
	dateList := append([]time.Time{startDate}, paymentDates...)

	if len(dateList) < 2 {
		return 0, [2]time.Time{}, 0, nil, errors.New("付息日期列表长度不足，无法计算付息周期")
	}

	n := 0
	maxIter := 1000
	var dateRange [2]time.Time // 存储当前付息周期[起始, 结束]
	for i := 0; i < maxIter; i++ {
		// 边界保护：避免数组越界（对齐Python的n == len(dateList)-1判断）
		if n >= len(dateList)-1 {
			dateRange[0] = dateList[len(dateList)-2]
			dateRange[1] = dateList[len(dateList)-1]
			break
		}

		// 构建当前周期[dateList[n], dateList[n+1]]
		dateRange[0] = dateList[n]
		dateRange[1] = dateList[n+1]

		// 判断结算日是否在当前周期内：date_range[0] <= settle_date < date_range[1]
		cond1 := settleDate.Equal(dateRange[0]) || settleDate.After(dateRange[0])
		cond2 := settleDate.Before(dateRange[1])
		if cond1 && cond2 {
			break // 找到目标周期，退出循环
		}

		// 已到最后一个周期，强制退出
		if n == len(dateList)-1 {
			break
		}

		// 未找到则继续下一个周期
		n++
	}

	// 计算剩余付息次数（对齐Python的last_n = len(dateList)-1 -n）
	lastN := len(dateList) - 1 - n

	// 计算当前付息周期天数ts
	var ts float64
	if f < 1 {
		return 0, [2]time.Time{}, 0, nil, errors.New("请检查付息次数（f必须≥1）")
	}

	switch {
	// 场景1：ActAct/ISDA / 交易所 / IBNCD债券
	case basis == ActAct || basis == ActActISDA || mode == Exchange || bondType == "IBNCD":
		// 计算周期实际天数（算头不算尾）
		ts = float64(dateRange[1].Sub(dateRange[0]).Hours() / 24)

	// 场景2：银行间 + A/365F/A/365
	case (basis == Act365F || basis == Act365) && mode == Interbank:
		ts = 365.0 / float64(f)

	// 场景3：银行间 + A/360/30/360
	case (basis == Act360 || basis == Thirty360) && mode == Interbank:
		ts = 360.0 / float64(f)

	// 兜底
	default:
		return 0, [2]time.Time{}, 0, nil, fmt.Errorf("不支持的计息规则组合: basis=%s, mode=%s, bondType=%s", basis, mode, bondType)
	}

	return ts, dateRange, lastN, dateList, nil
}

// 计算应计利息
func (b *BondPricingEngine) calculateAccruedInterest() (float64, error) {
	detail := b.Detail
	reserve := 15

	var (
		ts, t, ty, coupon float64
		fraction, ai      float64
		couponPeriod      [2]time.Time
		lastN             int
		dateList          []time.Time
		couponPer         []float64
	)

	ts = float64(detail.EndDate.Sub(detail.StartDate).Hours() / 24)
	t = float64(calculateT(b.SettleDate, detail.StartDate, detail.Market))
	ty = calculateTy(b.SettleDate, detail.StartDate, detail.Market, detail.DayCountBasis, detail.BondType)

	if detail.CouponType == ZeroCoupon && detail.PaymentFreq == 0 {
		if detail.IssuePrice < detail.FaceValue {
			coupon = detail.FaceValue - detail.IssuePrice
		}

		if detail.Market == Exchange && identifySpecialDay(detail.StartDate, detail.EndDate, detail.Market) {
			ts -= 1
		}

		fraction = t / ts
		ai = coupon * fraction
	} else if (detail.CouponType == ZeroCoupon && detail.PaymentFreq == 1) ||
		(detail.CouponType == FixedRate && detail.PaymentFreq == 0) {

		var err error
		fraction, t, err = calcDaycount(b.SettleDate, detail.StartDate, ty, b)
		if err != nil {
			return 0, fmt.Errorf("计算daycount失败: %v", err)
		}

		couponRate := detail.CouponRate / 100.0 * detail.FaceValue
		ai = couponRate * fraction
	} else if detail.CouponType == FixedRate && detail.PaymentFreq > 0 {
		var err error
		ts, couponPeriod, lastN, dateList, err = calculateTs(b.SettleDate, detail.StartDate, detail.EndDate, b)
		if err != nil {
			return 0, fmt.Errorf("计算ts失败: %v", err)
		}

		couponRate := detail.CouponRate / 100.0 * detail.FaceValue / float64(detail.PaymentFreq)
		if len(dateList) > 1 {
			couponPer = make([]float64, len(dateList)-1)
			for i := range couponPer {
				couponPer[i] = couponRate
			}
		}
		fraction, t, err = calcDaycount(b.SettleDate, couponPeriod[0], ts, b)
		if err != nil {
			return 0, fmt.Errorf("计算daycount失败: %v", err)
		}
		ai = couponPer[len(couponPer)-lastN] * fraction
	} else {
		return 0, errors.New("请检查息票类别")
	}

	ai = roundDecimal(ai, reserve)
	b.ts = ts
	b.t = t
	detail.AccruedInterest = ai
	detail.DirtyPrice = detail.CleanPrice + ai
	b.Detail = detail

	return ai, nil
}

// 根据净价计算全价
func (b *BondPricingEngine) CalculateDirtyPriceFromCleanPrice() (float64, float64, error) {
	if b.Detail.CleanPrice <= 0 {
		return 0, 0, errors.New("净价必须大于0")
	}

	// 计算应计利息
	accruedInterest, err := b.calculateAccruedInterest()
	if err != nil {
		return 0, 0, fmt.Errorf("计算应计利息失败: %v", err)
	}

	// 计算全价：全价 = 净价 + 应计利息
	dirtyPrice := b.Detail.CleanPrice + accruedInterest

	// 四舍五入
	dirtyPrice = roundDecimal(dirtyPrice, 8)
	accruedInterest = roundDecimal(accruedInterest, 8)

	b.Detail.AccruedInterest = accruedInterest
	b.Detail.DirtyPrice = dirtyPrice

	return dirtyPrice, accruedInterest, nil
}

// 创建债券定价引擎
func NewBondPricingEngine(settleDate string, detail BondDetail) (*BondPricingEngine, error) {
	// 解析结算日期
	settle, err := time.Parse("2006-01-02", settleDate)
	if err != nil {
		return nil, fmt.Errorf("结算日期格式错误: %v", err)
	}

	// 验证债券详情
	if detail.FaceValue <= 0 {
		return nil, errors.New("面值必须大于0")
	}

	if detail.StartDate.After(detail.EndDate) {
		return nil, errors.New("起息日不能晚于到期日")
	}

	if settle.Before(detail.StartDate) {
		return nil, errors.New("结算日不能早于起息日")
	}

	return &BondPricingEngine{
		SettleDate: settle,
		Detail:     detail,
	}, nil
}

func transformBondInfo(settleDateParam, couponTypeParam, startDateParam, maturityDateParam, dayCountBasisParam, exchMarketParam, bondTypeParam string, faceValueParam int,
	issuePriceParam, couponRateParam, cashTimesParam, cleanPriceParam float64) (BondDetail, string, error) {
	var detail BondDetail
	var err error
	settleDate := ""
	settleDateParse, err := time.Parse("20060102", settleDateParam)
	if err != nil {
		return detail, settleDate, fmt.Errorf("结算日解析失败: %v（输入值：%s）", err, settleDateParam)
	}
	settleDate = settleDateParse.Format("2006-01-02")

	detail.FaceValue = float64(faceValueParam)

	switch strings.TrimSpace(strings.ToUpper(couponTypeParam)) {
	case "FIXED":
		detail.CouponType = FixedRate
	case "ZERO":
		detail.CouponType = ZeroCoupon
	case "FLOAT":
		detail.CouponType = Floating
	default:
		return detail, settleDate, fmt.Errorf("不支持的息票类型：%s（仅支持FIXED/ZERO/FLOAT）", couponTypeParam)
	}

	detail.CouponRate = couponRateParam * 100.0

	detail.IssuePrice = issuePriceParam

	detail.StartDate, err = parseDBDateTime(startDateParam)
	if err != nil {
		return detail, settleDate, fmt.Errorf("起息日转换失败: %v（输入值：%s）", err, startDateParam)
	}

	detail.EndDate, err = parseDBDateTime(maturityDateParam)
	if err != nil {
		return detail, settleDate, fmt.Errorf("到期日转换失败: %v（输入值：%s）", err, maturityDateParam)
	}
	if detail.EndDate.Before(detail.StartDate) {
		return detail, settleDate, errors.New("到期日不能早于起息日")
	}

	if cashTimesParam < 0 {
		return detail, settleDate, errors.New("年付息次数不能为负数")
	}
	detail.PaymentFreq = int(cashTimesParam)

	detail.DayCountBasis, err = convertDayCountBasis(dayCountBasisParam)
	if err != nil {
		return detail, settleDate, fmt.Errorf("计息基准转换失败: %v（输入值：%s）", err, dayCountBasisParam)
	}
	if exchMarketParam == "NIB" {
		detail.Market = Interbank
	} else {
		detail.Market = Exchange
	}

	detail.BondType = strings.TrimSpace(bondTypeParam)
	if detail.BondType == "" {
		return detail, settleDate, errors.New("债券类型不能为空")
	}

	detail.CleanPrice = cleanPriceParam
	if detail.CleanPrice <= 0 {
		return detail, settleDate, errors.New("净价必须大于0")
	}

	detail.DirtyPrice = 0.0
	detail.AccruedInterest = 0.0

	return detail, settleDate, nil
}
/*
func full_price_calculator_simple_cases() {
	bondDetail, settleDate, err := transformBondInfo(
		"20230102",
		"FIXED",
		"2023-01-01 00:00:00",
		"2024-01-01 00:00:00",
		"ACTACT",
		"NIB",
		"CENTRAL_BANK_BILL",
		100,
		100.0,
		0.035,
		1.0,
		100.0,
	)
	if err != nil {
		fmt.Printf("转换债券信息失败: %v\n", err)
		return
	}
	// fmt.Printf("转换后的债券信息: %+v\n", bondDetail)
	// fmt.Printf("转换后的结算日: %s\n", settleDate)
	engine, err := NewBondPricingEngine(settleDate, bondDetail)
	if err != nil {
		fmt.Printf("创建定价引擎失败: %v\n", err)
		return
	}
	dirtyPrice, accrued, err := engine.CalculateDirtyPriceFromCleanPrice()
	if err != nil {
		fmt.Printf("计算全价失败: %v\n", err)
		return
	}
	fmt.Printf("转换后的全价: %.8f\n", dirtyPrice)
	fmt.Printf("转换后的应计利息: %.8f\n", accrued)
}*/

func ComputeDirtyPrice(applicationCfg *domain_cfg.ApplicationCfg, settlType string, tradeOrder *schema.TradeOrder, parValue int, cleanPrice float64, symbolData map[string]interface{}) (dirtyPrice float64, err error) {

	var settleDate string

	if settlType == "T+0" {
		settleDate = timeutil.ConvertMillisecondsToTime(tradeOrder.TransactTime).In(timeutil.CnTimeLocation).Format("20060102")
	} else {
		currDate := timeutil.ConvertMillisecondsToTime(tradeOrder.TransactTime).In(timeutil.CnTimeLocation).Format(time.DateOnly)
		settleDate = getTradeDate(applicationCfg.GetApplicationCfgItemMap(), currDate, 1)
		settleDate = strings.ReplaceAll(settleDate, "-", "")
	}

	couponType, _, err1 := attrutil.GetAttrValue(symbolData, "CouponType", enum.AttrValueType_STRING)
	if err != nil {
		err = err1
		return
	}
	startDate, _, err1 := attrutil.GetAttrValue(symbolData, "StartDate", enum.AttrValueType_STRING)
	if err != nil {
		err = err1
		return
	}
	maturityDate, _, err1 := attrutil.GetAttrValue(symbolData, "MaturityDate", enum.AttrValueType_STRING)
	if err != nil {
		err = err1
		return
	}
	dayCount, _, err1 := attrutil.GetAttrValue(symbolData, "DayCount", enum.AttrValueType_STRING)
	if err != nil {
		err = err1
		return
	}
	exchMarket, _, err1 := attrutil.GetAttrValue(symbolData, "SecurityExchange", enum.AttrValueType_STRING)
	if err != nil {
		err = err1
		return
	}
	bondType, _, err1 := attrutil.GetAttrValue(symbolData, "BondType", enum.AttrValueType_STRING)
	if err != nil {
		err = err1
		return
	}
	issuePrice, _, err1 := attrutil.GetAttrValue(symbolData, "IssuePrice", enum.AttrValueType_FLOAT)
	if err != nil {
		err = err1
		return
	}
	couponRate, _, err1 := attrutil.GetAttrValue(symbolData, "CouponRate", enum.AttrValueType_FLOAT)
	if err != nil {
		err = err1
		return
	}
	cashTimes, _, err1 := attrutil.GetAttrValue(symbolData, "CashTimes", enum.AttrValueType_FLOAT)
	if err != nil {
		err = err1
		return
	}

	bondDetail, settleDate, err1 := transformBondInfo(
		settleDate,
		couponType.(string),
		startDate.(string),
		maturityDate.(string),
		dayCount.(string),
		exchMarket.(string),
		bondType.(string),
		parValue,
		issuePrice.(float64),
		couponRate.(float64),
		cashTimes.(float64),
		cleanPrice,
	)
	if err1 != nil {
		err = fmt.Errorf("转换债券信息失败: %v", err1)
		return
	}

	engine, err1 := NewBondPricingEngine(settleDate, bondDetail)
	if err1 != nil {
		err = fmt.Errorf("创建定价引擎失败: %v", err1)
		return
	}
	dirtyPrice, _, err1 = engine.CalculateDirtyPriceFromCleanPrice()
	if err1 != nil {
		err = fmt.Errorf("计算全价失败: %v", err)
		return
	}
	
	return
}
