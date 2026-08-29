package order_check

import (
	"ficc-utils/api/order_check/ctpty_data"
	"ficc-utils/api/order_check/open_order"
	"ficc-utils/api/order_check/over_sold"
	"ficc-utils/api/order_check/securities"
	"ficc-utils/common/utils/config"
	"ficc-utils/common/utils/data_qry"
	"ficc-utils/common/utils/timeutil"
	"fmt"
	"log"
	"time"
)
func RunOrderChecks(config *config.Config) {
	for _, task := range config.OrderCheckTasks {
		switch task.Task {
			case "over_sold":
				registerTask(task, func() error{
					return over_sold.Task_CheckOversoldPosition(config.PositionServiceUrl, config.DataQryUrl, config.WebhookUrl)
				})
			case "ctpty_data":
				registerTask(task, func() error{
					return ctpty_data.Task_TrsCtptyData(config.PositionServiceUrl, config.DataQryUrl, config.WebhookUrl)
				})
			case "open_orders":
				if task.StartTime == "" {
					_, closeTime, err := data_qry.QryTradingTime(config.DataQryUrl)
					if err != nil {
						log.Printf("Register Task %s QryTradingTime error: %v, closeTime use default:17:30:00", task.Task, err)
						closeTime, _ = time.Parse(time.TimeOnly, "17:30:00")
					}
					task.StartTime = closeTime.Format(time.TimeOnly)
					if task.EndTime == "" {
						task.EndTime = closeTime.Add(15 * time.Minute).Format(time.TimeOnly)
					}
					if task.Interval == "" {
						task.Interval = "30s"
					}
				}
				registerTask(task, func() error{
					return open_order.Task_CheckOpenOrders(config.CurrentTradeOrdersServiceUrl, config.WebhookUrl)
				})
			case "symbols_t0":
				registerTask(task, func() error{
					return securities.Task_CheckSymbolsT0(config.DataQryUrl, config.WebhookUrl)
				})
			default:
				log.Printf("Unspported task: %s", task.Task)
		}
	}
}

func registerTask(sched config.OrderCheckTaskConfig, task func() error) {
	log.Printf("Register Task:%+v", sched)
	startTime, endTime, interval, err := parseOrderCheckTaskConfig(sched)
	if err != nil {
		log.Printf("Register Task %s config error: %v", sched.Task, err)
		return
	}

	go func() {
		for {
			now := time.Now().In(timeutil.CnTimeLocation)
			nowTime, _ := time.Parse(time.TimeOnly, now.Format(time.TimeOnly))

			if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday || nowTime.After(endTime) {
				time.Sleep(startTime.Add(24*time.Hour).Sub(nowTime))
				continue
			} else if nowTime.Before(startTime) {
				time.Sleep(startTime.Sub(nowTime))
				continue
			}

			log.Printf("Run Task:%s", sched.Task)
			for {
				err := task()
				if err != nil {
					log.Printf("Run Task:%s error: %v, retry after 10s", sched.Task, err)
					time.Sleep(time.Second * 10)
					continue
				}
				break
			}
			log.Printf("Run Task:%s end", sched.Task)
			time.Sleep(interval)
		}
	}()
}

func parseOrderCheckTaskConfig(taskConfig config.OrderCheckTaskConfig) (startTime, endTime time.Time, interval time.Duration, err error) {
	startTime, err = time.Parse(time.TimeOnly, taskConfig.StartTime)
	if err != nil {
		err = fmt.Errorf("parse StartTime error: %s", err)
		return
	}

	endTime, err = time.Parse(time.TimeOnly, taskConfig.EndTime)
	if err != nil {
		err = fmt.Errorf("parse EndTime error: %s", err)
		return
	}

	interval, err = time.ParseDuration(taskConfig.Interval)
	if err != nil {
		err = fmt.Errorf("parse Interval error: %s", err)
		return
	}

	return
}