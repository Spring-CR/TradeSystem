package timeutil

import (
	"log"
	"time"

	"github.com/manucorporat/try"
)

func StarTimeRangeSchedule(begin, end, timeLayout string, local *time.Location, onScheduleStar func(duringTimeRange bool), onTimeBeigin, onTimeEnd func()) (err error) {

	_onScheduleStar := func(duringTimeRange bool) {
		try.This(func() { onScheduleStar(duringTimeRange) }).Catch(func(err try.E) {
			log.Printf("fail to run onScheduleStar, err:%v\n", err)
		})
	}
	_onTimeBeigin := func() {
		try.This(func() { onTimeBeigin() }).Catch(func(err try.E) {
			log.Printf("fail to run onTimeBeigin, err:%v\n", err)
		})
	}
	_onTimeEnd := func() {
		try.This(func() { onTimeEnd() }).Catch(func(err try.E) {
			log.Printf("fail to run onTimeEnd, err:%v\n", err)
		})
	}

	dateTimeLayout := time.DateOnly + " " + timeLayout

	dateStr := time.Now().In(local).Format(time.DateOnly) + " "
	timeBegin, err := time.ParseInLocation(dateTimeLayout, dateStr+begin, local)
	if err != nil {
		return err
	}
	timeEnd, err := time.ParseInLocation(dateTimeLayout, dateStr+end, local)
	if err != nil {
		return err
	}
	if !timeEnd.After(timeBegin) {
		timeEnd = timeEnd.Add(24 * time.Hour)
	}

	log.Printf("======> StarTimeRangeSchedule, begin: %s, end: %s, timeLayout: %s, dateTimeLayout: %s, timeBegin: %v, timeEnd: %v\n", begin, end, timeLayout, dateTimeLayout, timeBegin, timeEnd)

	timeNow := time.Now()
	var waitForBegin time.Duration
	var waitForEnd time.Duration

	// 考虑临界点
	if timeNow.Equal(timeBegin) {
		onTimeBeigin()
		waitForBegin = time.Until(timeBegin.Add(24 * time.Hour))
		waitForEnd = waitForBegin + timeEnd.Sub(timeBegin)
	} else if timeNow.Equal(timeEnd) {
		onTimeEnd()
		waitForEnd = time.Until(timeEnd.Add(24 * time.Hour))
		waitForBegin = waitForEnd - timeEnd.Sub(timeBegin)
	} else if timeNow.After(timeBegin) && timeNow.Before(timeEnd) {
		_onScheduleStar(true)
		waitForBegin = time.Until(timeBegin.Add(24 * time.Hour))
		waitForEnd = time.Until(timeEnd)
	} else {
		_onScheduleStar(false)
		waitForBegin = time.Until(timeBegin.Add(24 * time.Hour))
		waitForEnd = time.Until(timeEnd.Add(24 * time.Hour))

		if waitForBegin > 24 * time.Hour {
			waitForBegin -= 24 * time.Hour
		}

		if waitForEnd > 24 * time.Hour {
			waitForEnd -= 24 * time.Hour
		}
	}

	// on time begin
	go func() {
		for {
			log.Printf("===>waitForBegin:%v\n", waitForBegin)
			time.Sleep(waitForBegin)
			_onTimeBeigin()
			// 计算休眠时长
			dateStr := time.Now().Add(24*time.Hour).In(local).Format(time.DateOnly) + " "
			timeBegin, _ = time.ParseInLocation(dateTimeLayout, dateStr+begin, local)
			waitForBegin = time.Until(timeBegin)
		}
	}()

	// on time end
	go func() {
		for {
			log.Printf("===>waitForEnd:%v\n", waitForEnd)
			time.Sleep(waitForEnd)
			_onTimeEnd()
			// 计算休眠时长
			dateStr := time.Now().Add(24*time.Hour).In(local).Format(time.DateOnly) + " "
			timeEnd, _ = time.ParseInLocation(dateTimeLayout, dateStr+end, local)
			waitForEnd = time.Until(timeEnd)
		}
	}()

	return
}
