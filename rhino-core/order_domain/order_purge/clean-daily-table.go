package order_purge

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type tableEntry struct {
	tableName string
	dateNum   int
}

func cleanDailyTable(db *sql.DB, maxReserve int, systemCode, businessCode string) {
	// 查询并打印表名
	rows, err := db.Query("SHOW TABLES") // 执行 SQL 命令
	if err != nil {
		log.Printf("查询失败:%v\n", err)
		return
	}
	defer rows.Close() // 确保关闭结果集

	dateStr := time.Now().Format("20060102")
	currDateNum, err := strconv.Atoi(dateStr)
	if err != nil {
		return
	}

	tableMap := make(map[string][]*tableEntry)

	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			log.Printf("读取失败:%v\n", err)
			return
		}

		//log.Printf("tableName:%s\n", tableName)

		shortTableName := tableName
		idx := strings.LastIndex(tableName, ".")
		if idx >= 0 {
			shortTableName = tableName[idx+1:]
		}

		if len(shortTableName) <= 9 {
			continue
		}

		if !strings.Contains(shortTableName, "_"+systemCode+"_"+businessCode+"_") {
			continue
		}

		if !strings.HasPrefix(shortTableName, "group_trade_orders_") &&
			!strings.HasPrefix(shortTableName, "trade_action_latest_resps_") &&
			!strings.HasPrefix(shortTableName, "trade_action_resps_") &&
			!strings.HasPrefix(shortTableName, "trade_orders_") {
			continue
		}

		dateNum, err := strconv.Atoi(shortTableName[len(shortTableName)-8:])
		if err != nil {
			continue
		}

		if dateNum < 20250101 || dateNum > currDateNum {
			continue
		}

		tableKey := shortTableName[:len(shortTableName)-9]
		tableEntries := tableMap[tableKey]
		tableEntries = append(tableEntries, &tableEntry{tableName: tableName, dateNum: dateNum})
		tableMap[tableKey] = tableEntries
	}
	err = rows.Close()
	if err != nil {
		log.Printf("rows close with error:%v\n", err)
		return
	}

	for k, v := range tableMap {
		sort.Slice(tableMap[k], func(i, j int) bool {
			return v[i].dateNum < v[j].dateNum
		})
	}

	for k, v := range tableMap {
		n := len(v)
		if n <= maxReserve {
			delete(tableMap, k)
			continue
		}
		v = v[:n-maxReserve]
		tableMap[k] = v
	}

	log.Println("delete the following table")

	for k, v := range tableMap {
		n := len(v)
		if n > maxReserve {
			v = v[:n-maxReserve]
		}
		tableMap[k] = v
	}

	for k, v := range tableMap {
		log.Printf("key:%s\n", k)
		for _, t := range v {
			log.Printf(" ===> delete %s\n", t.tableName)
			if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", t.tableName)); err != nil {
				log.Printf("drop table %s failed, error: %v", t.tableName, err)
			}
		}
	}
}