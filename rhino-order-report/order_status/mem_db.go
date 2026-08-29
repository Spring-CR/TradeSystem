package order_status

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
)

/*
var (
	extendAttrItems []*schema.ExtendAttrItem
	dbConfig        = `
PRAGMA journal_mode = WAL;
PRAGMA busy_timeout = 5000;
PRAGMA synchronous = OFF;
PRAGMA cache_size = 1000000000;
PRAGMA foreign_keys = true;
PRAGMA temp_store = memory;`

	createOrderTableInitSql = `
CREATE TABLE trade_orders_%s_%s(
id INTEGER PRIMARY KEY NOT NULL
`
	insertOrderSql = `INSERT INTO trade_orders_%s_%s(
id
`
	updateOrderSql = `UPDATE trade_orders_%s_%s SET
`
	fullUpdateOrderSql = `UPDATE trade_orders_%s_%s SET
`

	insertOrderRespSql = `INSERT INTO trade_action_resps_%s_%s(
`
)*/

func (s *OrderStatusReplica) initDBVar() {
	s.extendAttrItems = nil
	s.dbConfig = `
PRAGMA journal_mode = WAL;
PRAGMA busy_timeout = 5000;
PRAGMA synchronous = OFF;
PRAGMA cache_size = 1000000000;
PRAGMA foreign_keys = true;
PRAGMA temp_store = memory;`
	s.createOrderTableInitSql = `
CREATE TABLE trade_orders_%s_%s(
id INTEGER PRIMARY KEY NOT NULL
`
	s.insertOrderSql = `INSERT INTO trade_orders_%s_%s(
id
`
	s.updateOrderSql = `UPDATE trade_orders_%s_%s SET
`
	s.fullUpdateOrderSql = `UPDATE trade_orders_%s_%s SET
`
	s.insertOrderRespSql = `INSERT INTO trade_action_resps_%s_%s(
`

	s.orderTableName = ""
	s.orderRespTableName = ""
	s.hisOrderTableName = fmt.Sprintf("trade_orders_extend_%s_%s", s.systemCode, s.businessCode)
	s.hisOrderRespTableName = fmt.Sprintf("trade_action_resps_extend_%s_%s", s.systemCode, s.businessCode)
}

// https://www.runoob.com/sqlite/sqlite-data-types.html
func (s *OrderStatusReplica) initMemDb() {

	s.initDBVar()

	// 坑，配置prefix
	dbutil.ConfigDBFieldPrefix("")

	var err error
	s.dbWrite, err = sql.Open("sqlite3", ":memory:?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	// 设置写db
	s.dbWrite.SetMaxOpenConns(1)

	s.systemCode, s.businessCode = s.applicationCfg.GetSystemAndBusinessCodes()

	s.createOrderTableInitSql = fmt.Sprintf(s.createOrderTableInitSql, s.systemCode, s.businessCode)
	s.orderTableName = fmt.Sprintf("trade_orders_%s_%s", s.systemCode, s.businessCode)

	initSql := bytes.NewBufferString(s.dbConfig + s.createOrderTableInitSql)

	s.extendAttrItems = s.applicationCfg.GetExtendAttrItems()
	s.tradeRespAttrItems = s.applicationCfg.GetTradeActionRespAttrItems()

	/*
		// 添加状态属性
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_ord_status_update_time",
			AttrValueType: int(enum.AttrValueType_INT),
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_ord_status",
			AttrValueType: int(enum.AttrValueType_STRING),
			Index:         true,
		})

		// 添加新的属性请从这里开始向下加入 =====================>
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_db_insert_time",
			AttrValueType: int(enum.AttrValueType_INT),
			Index: true,
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_last_shares",
			AttrValueType: int(enum.AttrValueType_INT),
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_last_px",
			AttrValueType: int(enum.AttrValueType_FLOAT),
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_leaves_qty",
			AttrValueType: int(enum.AttrValueType_INT),
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_cum_qty",
			AttrValueType: int(enum.AttrValueType_INT),
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_avg_px",
			AttrValueType: int(enum.AttrValueType_FLOAT),
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_ord_rej_reason",
			AttrValueType: int(enum.AttrValueType_STRING),
		})
		// 添加新的属性请从这里开始向上加入 <=====================

		// id属性，用于关联
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_app_ord_id",
			AttrValueType: int(enum.AttrValueType_STRING),
			Index:         true,
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_cl_ord_id",
			AttrValueType: int(enum.AttrValueType_STRING),
			Index:         true,
		})
	*/

	s.extendAttrItems = domain_cfg.ConfigExtendAttrItemsForTradeOrderExtending(s.extendAttrItems)

	s.extendAttrMap = make(map[string]*schema.ExtendAttrItem)

	for _, extendAttrItem := range s.extendAttrItems {
		s.extendAttrMap[extendAttrItem.AttrName] = extendAttrItem
		if extendAttrItem.Unique {
			initSql.WriteString("," + extendAttrItem.AttrName + " " + s.getDbFieldType(extendAttrItem.AttrValueType) + " UNIQUE \n")
		} else {
			initSql.WriteString("," + extendAttrItem.AttrName + " " + s.getDbFieldType(extendAttrItem.AttrValueType) + "\n")
		}
	}

	initSql.WriteString(");\n")

	for _, extendAttrItem := range s.extendAttrItems {
		// 创建索引
		if extendAttrItem.Index {
			initSql.WriteString("CREATE INDEX index_" + s.orderTableName + "_" + extendAttrItem.AttrName + " ON " + s.orderTableName + "(" + extendAttrItem.AttrName + ");\n")
		}
	}

	createOrderSqlStr := initSql.String()
	log.Printf("createOrderSqlStr:%s\n", createOrderSqlStr)

	s.refineOrderInsertSql()
	s.refineOrderUpdateSql()
	s.refineFullOrderUpdateSql()

	// 初始化DB
	_, err = s.dbWrite.Exec(createOrderSqlStr)
	if err != nil {
		log.Fatal(err)
	}

	s.createOrderResponseTable(s.extendAttrItems, s.tradeRespAttrItems)

	// s.dbRead, err = sql.Open("sqlite3", ":memory:")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // 设置读db
	// s.dbRead.SetMaxOpenConns(max(4, runtime.NumCPU()))
	s.dbRead = s.dbWrite // SQLite内存模式无法读写分离

	log.Println("success create memory database!")

	//panic("only for test!")
}

func (s *OrderStatusReplica) createOrderResponseTable(extendAttrItems []*schema.ExtendAttrItem, tradeRespAttrItems []*schema.TradeActionRespAttrItem) {

	creatOrderResponseTableSql := `
CREATE TABLE trade_action_resps_%s_%s(
id INTEGER PRIMARY KEY AUTOINCREMENT
`
	creatOrderResponseTableSql = fmt.Sprintf(creatOrderResponseTableSql, s.systemCode, s.businessCode)
	s.orderRespTableName = fmt.Sprintf("trade_action_resps_%s_%s", s.systemCode, s.businessCode)

	initSql := bytes.NewBufferString(creatOrderResponseTableSql)

	/*
		// 去掉Order表里的唯一性约束
		for _, extendAttrItem := range extendAttrItems {
			if extendAttrItem.Unique {
				extendAttrItem.Unique = false
				extendAttrItem.Index = true
			}
		}

		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_orig_cl_ord_id",
			AttrValueType: int(enum.AttrValueType_STRING),
			Index:         true,
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_exec_id",
			AttrValueType: int(enum.AttrValueType_STRING),
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_exec_ref_id",
			AttrValueType: int(enum.AttrValueType_STRING),
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_exec_trans_type",
			AttrValueType: int(enum.AttrValueType_STRING),
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_exec_type",
			AttrValueType: int(enum.AttrValueType_STRING),
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_msg_time",
			AttrValueType: int(enum.AttrValueType_INT),
		})
		extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
			AttrName:      "f_channel_code",
			AttrValueType: int(enum.AttrValueType_STRING),
		})
	*/

	extendAttrItems = domain_cfg.ConfigExtendAttrItemsForTradeActionRespExtending(extendAttrItems)

	// 添加订单属性、状态属性、f_orig_cl_ord_id
	for _, extendAttrItem := range extendAttrItems {
		if extendAttrItem.Unique {
			initSql.WriteString("," + extendAttrItem.AttrName + " " + s.getDbFieldType(extendAttrItem.AttrValueType) + " UNIQUE \n")
		} else {
			initSql.WriteString("," + extendAttrItem.AttrName + " " + s.getDbFieldType(extendAttrItem.AttrValueType) + "\n")
		}
	}

	for _, extendAttrItem := range tradeRespAttrItems {
		if extendAttrItem.Unique {
			initSql.WriteString("," + extendAttrItem.AttrName + " " + s.getDbFieldType(extendAttrItem.AttrValueType) + " UNIQUE \n")
		} else {
			initSql.WriteString("," + extendAttrItem.AttrName + " " + s.getDbFieldType(extendAttrItem.AttrValueType) + "\n")
		}
	}

	initSql.WriteString(");\n")

	for _, extendAttrItem := range extendAttrItems {
		// 创建索引
		if extendAttrItem.Index {
			initSql.WriteString("CREATE INDEX index_" + s.orderRespTableName + "_" + extendAttrItem.AttrName + " ON " + s.orderRespTableName + "(" + extendAttrItem.AttrName + ");\n")
		}
	}

	for _, extendAttrItem := range tradeRespAttrItems {
		// 创建索引
		if extendAttrItem.Index {
			initSql.WriteString("CREATE INDEX index_" + s.orderRespTableName + "_" + extendAttrItem.AttrName + " ON " + s.orderRespTableName + "(" + extendAttrItem.AttrName + ");\n")
		}
	}

	// 联合索引
	initSql.WriteString("CREATE UNIQUE INDEX index_" + s.orderRespTableName + "_pk_tar" + " ON " + s.orderRespTableName + "(f_cl_ord_id,f_exec_id,f_channel_code);\n")

	creatOrderResponseTableSql = initSql.String()
	log.Printf("======> creatOrderResponseTableSql:%s\n", creatOrderResponseTableSql)

	// 创建订单执行回报表
	_, err := s.dbWrite.Exec(creatOrderResponseTableSql)
	if err != nil {
		log.Fatal(err)
	}

	s.refineOrderRespInsertSql(s.systemCode, s.businessCode, extendAttrItems, tradeRespAttrItems)
}

func (s *OrderStatusReplica) refineOrderInsertSql() {
	s.insertOrderSql = fmt.Sprintf(s.insertOrderSql, s.systemCode, s.businessCode)
	sqlBuf := bytes.NewBufferString(s.insertOrderSql)
	for _, extendAttrItem := range s.extendAttrItems {
		sqlBuf.WriteString("," + extendAttrItem.AttrName + "\n")
	}
	sqlBuf.WriteString(") VALUES (")
	n := len(s.extendAttrItems) + 1 // 加1是因为添加了id列
	n_1 := n - 1
	for i := 0; i < n; i++ {
		sqlBuf.WriteString("?")
		if i == n_1 {
			sqlBuf.WriteString(")")
		} else {
			sqlBuf.WriteString(",")
		}
	}
	s.insertOrderSql = sqlBuf.String()

	log.Printf("======>refineOrderInsertSql:%s\n", s.insertOrderSql)
}

func (s *OrderStatusReplica) refineOrderUpdateSql() {
	s.updateOrderSql = fmt.Sprintf(s.updateOrderSql, s.systemCode, s.businessCode)
	sqlBuf := bytes.NewBufferString(s.updateOrderSql)
	hit := false
	i := 0
	for _, extendAttrItem := range s.extendAttrItems {

		if extendAttrItem.AttrName != "f_ord_status_update_time" && !hit {
			continue
		} else {
			hit = true
		}

		if extendAttrItem.AttrName == "f_app_ord_id" {
			break
		}

		if i != 0 {
			sqlBuf.WriteString(",")
		}
		sqlBuf.WriteString(extendAttrItem.AttrName + "=?\n")
		i++
	}
	sqlBuf.WriteString("WHERE f_app_ord_id=?")

	s.updateOrderSql = sqlBuf.String()

	log.Printf("======>refineOrderUpdateSql:%s\n", s.updateOrderSql)
}

func (s *OrderStatusReplica) refineFullOrderUpdateSql() {
	s.fullUpdateOrderSql = fmt.Sprintf(s.fullUpdateOrderSql, s.systemCode, s.businessCode)
	sqlBuf := bytes.NewBufferString(s.fullUpdateOrderSql)
	for i, extendAttrItem := range s.extendAttrItems {
		if i != 0 {
			sqlBuf.WriteString(",")
		}
		sqlBuf.WriteString(extendAttrItem.AttrName + "=?\n")
		i++
	}
	sqlBuf.WriteString("WHERE f_app_ord_id=?")

	s.fullUpdateOrderSql = sqlBuf.String()

	log.Printf("======>refineFullOrderUpdateSql:%s\n", s.updateOrderSql)
}

func (s *OrderStatusReplica) refineOrderRespInsertSql(systemCode, businessCode string, extendAttrItems []*schema.ExtendAttrItem, tradeRespAttrItems []*schema.TradeActionRespAttrItem) {
	s.insertOrderRespSql = fmt.Sprintf(s.insertOrderRespSql, systemCode, businessCode)
	sqlBuf := bytes.NewBufferString(s.insertOrderRespSql)
	for i, extendAttrItem := range extendAttrItems {
		if i != 0 {
			sqlBuf.WriteString(",")
		}
		sqlBuf.WriteString(extendAttrItem.AttrName + "\n")
	}
	for i, extendAttrItem := range tradeRespAttrItems {
		if len(extendAttrItems) > 0 || i != 0 {
			sqlBuf.WriteString(",")
		}
		sqlBuf.WriteString(extendAttrItem.AttrName + "\n")
	}

	sqlBuf.WriteString(") VALUES (")
	n := len(extendAttrItems) + len(tradeRespAttrItems) // 不加加1是因为id采用自增策略
	n_1 := n - 1
	for i := 0; i < n; i++ {
		sqlBuf.WriteString("?")
		if i == n_1 {
			sqlBuf.WriteString(")")
		} else {
			sqlBuf.WriteString(",")
		}
	}
	s.insertOrderRespSql = sqlBuf.String()

	log.Printf("======>refineOrderRespInsertSql:%s\n", s.insertOrderRespSql)
}

func (s *OrderStatusReplica) getDbFieldType(attrValueType int) string {
	t := enum.AttrValueType(attrValueType)
	switch t {
	case enum.AttrValueType_INT:
		return "INTEGER"
	case enum.AttrValueType_FLOAT:
		return "REAL"
	case enum.AttrValueType_BOOL:
		return "BOOLEAN"
	default:
		return "TEXT"
	}
}

func (s *OrderStatusReplica) printDBTables() {

	log.Println("start to printReadDBTables")
	// 获取所有表名
	rows, err := s.dbRead.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 打印所有表名
	fmt.Println("Tables in the database:")
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(tableName)
	}

	log.Println("start to printWriteDBTables")
	// 获取所有表名
	rows, err = s.dbWrite.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 打印所有表名
	fmt.Println("Tables in the database:")
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(tableName)
	}
}
