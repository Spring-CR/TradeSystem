package timeutil

import (
	"errors"
	"fmt"
	"rhino-common/domain_error"
	"strconv"
	"strings"
	"time"
)

const (
	DateSeconds        = 86400
	TransactTimeLayout = "20060102-15:04:05.000"
)

var (
	secondsEastOfUTC = int((8 * time.Hour).Seconds()) // 中国标准时间为UTC+8
	tzChina          = time.FixedZone("China Standard Time", secondsEastOfUTC)
	timeZoneMap2     = make(map[int]*time.Location)
	timeZoneMap      = make(map[string]*time.Location)
	CnTimeZoneName   = "Asia/Shanghai"
	CnTimeLocation   *time.Location
	UsTimeZoneName   = "America/New_York"
	UsTimeLocation   *time.Location
)

func init() {

	for i := 0; i < 24; i++ {
		timeZoneMap2[i] = GetTimeZone2(i)
	}

	var err error

	CnTimeLocation, err = time.LoadLocation(CnTimeZoneName) // 显式加载中国时区
	if err != nil {
		panic("fail to get time LoadLocation for " + CnTimeZoneName)
	}

	UsTimeLocation, err = time.LoadLocation(UsTimeZoneName) // 显式加载美国时区
	if err != nil {
		panic("fail to get time LoadLocation for " + UsTimeZoneName)
	}

	timeZoneMap[CnTimeZoneName] = CnTimeLocation
	timeZoneMap[UsTimeZoneName] = UsTimeLocation
}

func GetTimeZone2(zoneNum int) *time.Location {
	return time.FixedZone(fmt.Sprintf("%d Time Zone", zoneNum), int((time.Duration(zoneNum) * time.Hour).Seconds()))
}

func GetTimeZone(timeZone string) *time.Location {
	loc, ok := timeZoneMap[timeZone]
	if !ok {
		domain_error.ProcessSevereError(true, 5, nil, errors.New("unwarm up timezone "+timeZone), "unwarm up timezone")
	}
	return loc
}

func GetDateStrForTimeZone2(timestamp int64, zoneNum int) (string, error) {
	if zoneNum < 0 || zoneNum > 23 {
		return "", fmt.Errorf("unsupport time zone number: %d", zoneNum)
	}
	t := ConvertMillisecondsToTime(timestamp)
	dateStr := t.In(timeZoneMap2[zoneNum]).Format("20060102")
	return dateStr, nil
}

func GetDateStrForTimeZone(timestamp int64, timeZone string) (string, error) {
	t := ConvertMillisecondsToTime(timestamp)
	loc, ok := timeZoneMap[timeZone]
	var err error
	if !ok {
		loc, err = time.LoadLocation(timeZone)
		if err != nil {
			return "", err
		}
	}
	dateStr := t.In(loc).Format("20060102")
	return dateStr, err
}

func ParseTimeStrToTime(layout, timeStr string, timeZone string) (time.Time, error) {
	loc, ok := timeZoneMap[timeZone]
	var err error
	if !ok {
		loc, err = time.LoadLocation(timeZone)
		if err != nil {
			return time.Time{}, err
		}
	}
	return time.ParseInLocation(layout, timeStr, loc)
}

func ParseTimeStrToTimeByTimeLocation(layout, timeStr string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(layout, timeStr, loc)
}

func GetDateNum(ts ...time.Time) int {
	var t time.Time
	if len(ts) > 0 {
		t = ts[0]
	} else {
		t = time.Now()
	}
	dateStr := t.In(tzChina).Format("20060102")
	num, _ := strconv.Atoi(dateStr)
	return num
}

func GetTimeNum(ts ...time.Time) int {
	var t time.Time
	if len(ts) > 0 {
		t = ts[0]
	} else {
		t = time.Now()
	}
	dateStr := t.In(tzChina).Format("150405")
	num, _ := strconv.Atoi(dateStr)
	return num
}

func GetTimeByDateNum(dateNum int) (time.Time, error) {
	dateStr := strconv.Itoa(dateNum)
	return time.ParseInLocation(time.DateOnly, dateStr, tzChina)
}

// ConvertMicrosecondsToTime 接收微秒时间戳并返回时间对象
func ConvertMicrosecondsToTime(microseconds int64) time.Time {
	// 将微秒转换为纳秒
	nanoseconds := microseconds * int64(time.Microsecond)
	// 使用纳秒创建time.Time对象
	return time.Unix(0, nanoseconds)
}

// ConvertTimeToMicroseconds 接收一个time.Time对象并返回微秒时间戳
func ConvertTimeToMicroseconds(t time.Time) int64 {
	// 获取纳秒级时间戳
	nanoseconds := t.UnixNano()
	// 将纳秒转换为微秒（1微秒 = 1000纳秒）
	microseconds := nanoseconds / int64(time.Microsecond)
	return microseconds
}

// ConvertMillisecondsToTime 接收毫秒时间戳并返回时间对象
func ConvertMillisecondsToTime(milliseconds int64) time.Time {
	// 将毫秒转换为纳秒
	nanoseconds := milliseconds * int64(time.Millisecond)
	// 使用纳秒创建time.Time对象
	return time.Unix(0, nanoseconds)
}

// ConvertTimeToMilliseconds 接收一个time.Time对象并返回毫秒时间戳
func ConvertTimeToMilliseconds(t time.Time) int64 {
	// 获取纳秒级时间戳
	nanoseconds := t.UnixNano()
	// 将纳秒转换为毫秒（1毫秒 = 1,000,000纳秒）
	milliseconds := nanoseconds / int64(time.Millisecond)
	return milliseconds
}

// GetCumulativeSecondsFromSimpleTimeString: 将一个 HH:mm:ss 格式的时间字符串转换为累计的秒数
// timeStr: HH:mm:ss
func GetCumulativeSecondsFromSimpleTimeString(timeStr string) (int, error) {
	err := errors.New("illegal time string " + timeStr)
	if len(timeStr) != 8 {
		return 0, err
	}
	strs := strings.Split(timeStr, ":")
	if len(strs) != 3 {
		return 0, err
	}
	hours, e := strconv.Atoi(strs[0])
	if e != nil {
		return 0, err
	}
	minutes, e := strconv.Atoi(strs[1])
	if e != nil {
		return 0, err
	}
	seconds, e := strconv.Atoi(strs[2])
	if e != nil {
		return 0, err
	}
	return 3600*hours + 60*minutes + seconds, nil
}

func GetCurrentCumulativeSeconds2(timeZone int) int {
	t := time.Now()
	timeStr := t.In(GetTimeZone2(timeZone)).Format(time.TimeOnly)
	cumSeconds, _ := GetCumulativeSecondsFromSimpleTimeString(timeStr)
	return cumSeconds
}

func GetCumulativeSeconds2(t time.Time, timeZone int) int {
	timeStr := t.In(GetTimeZone2(timeZone)).Format(time.TimeOnly)
	cumSeconds, _ := GetCumulativeSecondsFromSimpleTimeString(timeStr)
	return cumSeconds
}

func GetCurrentCumulativeSeconds(timeZone string) int {
	t := time.Now()

	loc, ok := timeZoneMap[timeZone]
	if !ok {
		domain_error.ProcessSevereError(true, 5, nil, errors.New("unwarm up timezone "+timeZone), "unwarm up timezone")
	}

	timeStr := t.In(loc).Format(time.TimeOnly)
	cumSeconds, _ := GetCumulativeSecondsFromSimpleTimeString(timeStr)
	return cumSeconds
}

func GetCumulativeSeconds(t time.Time, timeZone string) int {

	loc, ok := timeZoneMap[timeZone]
	if !ok {
		domain_error.ProcessSevereError(true, 5, nil, errors.New("unwarm up timezone "+timeZone), "unwarm up timezone")
	}

	timeStr := t.In(loc).Format(time.TimeOnly)
	cumSeconds, _ := GetCumulativeSecondsFromSimpleTimeString(timeStr)
	return cumSeconds
}

// 注意：为了极端速度，此方法不做验证
// 因使用time format且无错的结果作为入参数
func Parse8BitDateStrToNum(dateStr string) int64 {
	n := int64(dateStr[0]-'0')*1e7 +
		int64(dateStr[1]-'0')*1e6 +
		int64(dateStr[2]-'0')*1e5 +
		int64(dateStr[3]-'0')*1e4 +
		int64(dateStr[4]-'0')*1e3 +
		int64(dateStr[5]-'0')*1e2 +
		int64(dateStr[6]-'0')*1e1 +
		int64(dateStr[7]-'0')*1e0
	return n
}

// 是否夏令时间
func IsDST() bool {
	// 获取该时区当前时间
	currentTime := time.Now().In(UsTimeLocation)
	// 判断是否处于夏令时
	return currentTime.IsDST()
}

func WarmUpTimeLocation(timezone string) {
	_, ok := timeZoneMap[timezone]
	if ok {
		// 不必重复设置
		return
	}
	loc, err := time.LoadLocation(timezone) // 显式加载美国时区
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, errors.New("fail to get time LoadLocation for "+timezone), "fail to get time LoadLocation for "+timezone)
	}

	timeZoneMap[timezone] = loc
}

func Until(timeStr, timeLayout, timeZone string) (time.Duration, error) {
	// 计算休眠时间并开始休眠
	loc := GetTimeZone(timeZone)
	if loc == nil {
		return 0, fmt.Errorf("fail to get timeZone for %s", timeZone)
	}
	dateStr := time.Now().In(loc).Format(time.DateOnly)
	dateTimeStr := dateStr + " " + timeStr
	endtime, err := time.ParseInLocation(time.DateOnly+" "+timeLayout, dateTimeStr, loc)
	if err != nil {
		return 0, fmt.Errorf("time ParseInLocation error")
	}

	if time.Now().After(endtime) {
		endtime = endtime.Add(24 * time.Hour)
	}

	return time.Until(endtime), nil
}

func GetTimeMinutes(t time.Time) int {
	return t.Hour()*60 + t.Minute()
}
