package order_archive

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/bean"
	"rhino-common/utils/dbutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"strings"
)

func (a *OrderArchiver) copyDailyTableStructure(originalTableName string, archivingLog *schema.DataArchivingLog) (string, error) {
	targetTableName := fmt.Sprintf(originalTableName+"_%s_%s_%s", archivingLog.SystemCode, archivingLog.BusinessCode, archivingLog.ArchivingDate)
	if archivingLog.TaskName != "" {
		targetTableName = fmt.Sprintf(originalTableName+"_%s_%s_%s_%s", archivingLog.SystemCode, archivingLog.BusinessCode, archivingLog.TaskName, archivingLog.ArchivingDate)
	}
	return targetTableName, dbutil.CopyTableStructure(a.applicationCfg.GetAppDB(), a.applicationCfg.GetCentralDB(), originalTableName, targetTableName, true)
}

func (a *OrderArchiver) createDailyTables(archivingLog *schema.DataArchivingLog) (
	dailyGroupTradeOrderTableName,
	dailyTradeActionLatestRespTableName,
	dailyTradeActionRespTableName,
	dailyTradeOrderTableName, extendDailyTradeOrderTableName, extendDailyTradeActionRespTableName string,
	err error) {

	// 进行表的复制（在这个过程中，如果目标表存在，则会被清理并重新创建）
	dailyGroupTradeOrderTableName, err = a.copyDailyTableStructure("group_trade_orders", archivingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createDailyTables")
		return
	}
	dailyTradeActionLatestRespTableName, err = a.copyDailyTableStructure("trade_action_latest_resps", archivingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createDailyTables")
		return
	}
	dailyTradeActionRespTableName, err = a.copyDailyTableStructure("trade_action_resps", archivingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createDailyTables")
		return
	}
	dailyTradeOrderTableName, err = a.copyDailyTableStructure("trade_orders", archivingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createDailyTables")
		return
	}

	extendAttrItems := a.applicationCfg.GetExtendAttrItems()

	// 创建TradeOrder扩展表
	tableNameTmpl := "trade_orders_extend_%s_%s_%s"
	extendDailyTradeOrderTableName = fmt.Sprintf(tableNameTmpl, archivingLog.SystemCode, archivingLog.BusinessCode, archivingLog.ArchivingDate)
	if archivingLog.TaskName != "" {
		tableNameTmpl := "trade_orders_extend_%s_%s_%s_%s"
		extendDailyTradeOrderTableName = fmt.Sprintf(tableNameTmpl, archivingLog.SystemCode, archivingLog.BusinessCode, archivingLog.TaskName, archivingLog.ArchivingDate)
	}
	err = dbutil.CreateSimpleTableWithSysID(a.applicationCfg.GetCentralDB(), extendDailyTradeOrderTableName, true)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createDailyTables")
		return
	}
	extendAttrItems = domain_cfg.ConfigExtendAttrItemsForTradeOrderExtending(extendAttrItems)
	err = ExtendTableColumn(a.applicationCfg.GetCentralDB(), extendDailyTradeOrderTableName, extendAttrItems)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createDailyTables")
		return
	}
	// 设置插入语句
	a._insertDataToExtendTradeOrderTableWithIDSql = a.getInsertSqlTemplate("%s", true, extendAttrItems)
	a._insertDataToExtendTradeOrderTableWithoutIDSql = a.getInsertSqlTemplate("%s", false, extendAttrItems)

	// 创建TradeActionResp扩展表
	tableNameTmpl = "trade_action_resps_extend_%s_%s_%s"
	extendDailyTradeActionRespTableName = fmt.Sprintf(tableNameTmpl, archivingLog.SystemCode, archivingLog.BusinessCode, archivingLog.ArchivingDate)
	if archivingLog.TaskName != "" {
		tableNameTmpl = "trade_action_resps_extend_%s_%s_%s_%s"
		extendDailyTradeActionRespTableName = fmt.Sprintf(tableNameTmpl, archivingLog.SystemCode, archivingLog.BusinessCode, archivingLog.TaskName, archivingLog.ArchivingDate)
	}
	err = dbutil.CreateSimpleTableWithSysID(a.applicationCfg.GetCentralDB(), extendDailyTradeActionRespTableName, true)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createDailyTables")
		return
	}
	extendAttrItems = domain_cfg.ConfigExtendAttrItemsForTradeActionRespExtending(extendAttrItems)
	for _, extendAttrItem := range a.applicationCfg.GetTradeActionRespAttrItems() {
		item := &schema.ExtendAttrItem{}
		bean.Copy(extendAttrItem).To(item)
		extendAttrItems = append(extendAttrItems, item)
	}
	err = ExtendTableColumn(a.applicationCfg.GetCentralDB(), extendDailyTradeActionRespTableName, extendAttrItems)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createDailyTables")
		return
	}
	// 设置插入语句
	a._insertDataToExtendTradeActionRespTableWithIDSql = a.getInsertSqlTemplate("%s", true, extendAttrItems)
	a._insertDataToExtendTradeActionRespTableWithoutIDSql = a.getInsertSqlTemplate("%s", false, extendAttrItems)

	log.Printf("_insertDataToExtendTradeOrderTableWithIDSql:%s\n", a._insertDataToExtendTradeOrderTableWithIDSql)
	log.Printf("_insertDataToExtendTradeOrderTableWithoutIDSql:%s\n", a._insertDataToExtendTradeOrderTableWithoutIDSql)
	log.Printf("_insertDataToExtendTradeActionRespTableWithIDSql:%s\n", a._insertDataToExtendTradeActionRespTableWithIDSql)
	log.Printf("_insertDataToExtendTradeActionRespTableWithoutIDSql:%s\n", a._insertDataToExtendTradeActionRespTableWithoutIDSql)

	return
}

func (a *OrderArchiver) getInsertSqlTemplate(tableName string, withID bool, extendAttrItems []*schema.ExtendAttrItem) string {
	sqlBuf := bytes.NewBufferString("INSERT INTO " + tableName + "(\n")
	if withID {
		sqlBuf.WriteString("f_id\n")
	}
	for i, extendAttrItem := range extendAttrItems {
		if i == 0 && !withID {
			sqlBuf.WriteString(extendAttrItem.AttrName + "\n")
		} else {
			sqlBuf.WriteString("," + extendAttrItem.AttrName + "\n")
		}
	}
	sqlBuf.WriteString(") VALUES (")

	n := len(extendAttrItems)
	if withID {
		// 加1是因为添加了id列
		n += 1
	}
	n_1 := n - 1
	for i := 0; i < n; i++ {
		sqlBuf.WriteString("?")
		if i == n_1 {
			sqlBuf.WriteString(")")
		} else {
			sqlBuf.WriteString(",")
		}
	}
	return sqlBuf.String()
}

func (a *OrderArchiver) getInsertExtendTradeOrderArgs(order *schema.TradeOrder, withID bool, clOrdID string, archivingLog *schema.DataArchivingLog) (args []interface{}, err error) {

	if clOrdID == "" {
		clOrdID = order.ClOrdID
	}

	extendAttrMap := order.ExtendAttrMap
	if extendAttrMap == nil {
		return nil, errors.New("extendAttrMap is empty for order " + order.AppOrdID)
	}

	if withID {
		args = append(args, order.ID)
	}

	// 加入应用自定义字段
	extendAttrItems := a.applicationCfg.GetExtendAttrItems()
	n := len(extendAttrItems)
	for i := 0; i < n; i++ {
		extendAttrItem := extendAttrItems[i]
		val, ok, err := attrutil.GetAttrValue(extendAttrMap, extendAttrItem.AttrName, enum.AttrValueType(extendAttrItem.AttrValueType))
		if err != nil {
			log.Printf("error occur while get value of %s, ok = %v\n", extendAttrItem.AttrName, ok)
			return nil, err
		}
		args = append(args, val)
	}

	// 加入系统字段（对照initMemDb的字段次序，位置不能乱）
	// args 是用于插入新记录的参数列表
	args = append(args, order.OrdStatusUpdateTime)
	args = append(args, order.OrdStatus)
	args = append(args, order.DBInsertTime)
	args = append(args, order.TransactTime)
	args = append(args, order.Reviewer)
	args = append(args, order.ApproveStatus)
	args = append(args, order.LastShares)
	args = append(args, order.LastPx)
	args = append(args, order.LeavesQty)
	args = append(args, order.CumQty)
	args = append(args, order.AvgPx)
	args = append(args, order.OrdRejReason)
	args = append(args, order.AppOrdID)
	args = append(args, clOrdID)
	args = append(args, order.OrdCreateTime)
	args = append(args, order.AlgParams)
	args = append(args, order.TradeDate)

	// 增加了四个用户
	// args = append(args, order.OrdCreator)
	// args = append(args, order.OrdDraftUpdateUser)
	// args = append(args, order.OrdDraftDelUser)
	// args = append(args, order.OrdExecUser)

	if archivingLog !=nil {
		args = append(args, archivingLog.ArchivingDate, archivingLog.TaskName)
	}

	return
}

func (a *OrderArchiver) getInsertExtendTradeActionRespArgs(order *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, withID bool, archivingLog *schema.DataArchivingLog) (args []interface{}, err error) {

	args, err = a.getInsertExtendTradeOrderArgs(order, withID, tradeActionResp.ClOrdID, nil)
	if err != nil {
		return
	}

	// 加入系统字段（对照createOrderResponseTable的字段次序，位置不能乱）
	args = append(args, tradeActionResp.OrigClOrdID)
	args = append(args, tradeActionResp.ExecID)
	args = append(args, tradeActionResp.ExecRefID)
	args = append(args, tradeActionResp.ExecTransType)
	args = append(args, tradeActionResp.ExecType)
	args = append(args, tradeActionResp.MsgTime)
	args = append(args, tradeActionResp.ChannelCode)

	// 加入扩展字段
	if tradeActionResp.ExtendAttrMap == nil {
		tradeActionResp.RecoverExtendAttrMap()
	}
	for _, extendAttrItem := range a.applicationCfg.GetTradeActionRespAttrItems() {
		val, ok, err := attrutil.GetAttrValue(tradeActionResp.ExtendAttrMap, extendAttrItem.AttrName, enum.AttrValueType(extendAttrItem.AttrValueType))
		if err != nil {
			log.Printf("error occur while get value of %s, ok = %v\n", extendAttrItem.AttrName, ok)
		}
		args = append(args, val)
	}

	if archivingLog !=nil {
		args = append(args, archivingLog.ArchivingDate, archivingLog.TaskName)
	}

	return
}

// 该函数能根据extendAttrItems的扩展字段的声明，扩展由db、table参数指定的MySql表。
// 注意：
// 1、检查原表中是否已经存在了同名的字段，如果存在了则该扩展项，继续处理下一个扩展项；
// 2、中途如果出错，请直接返回即可；
// 3、建索引，请确保索引名称的唯一性，基于表名+字段名+是否唯一性索引来命名；
// 4、不需要进行字段名的合法性校验，由应用层保证名称合法性；
func ExtendTableColumn(db *sql.DB, table string, extendAttrItems []*schema.ExtendAttrItem) error {
	// 遍历所有需要扩展的字段
	for _, item := range extendAttrItems {
		// 1. 检查字段是否已存在
		exists, err := columnExists(db, table, item.AttrName)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		// 2. 添加新字段
		if err := addColumn(db, table, item); err != nil {
			return err
		}

		// 3. 处理索引
		if item.Unique {
			if err := createIndex(db, table, item, true); err != nil {
				return err
			}
		} else if item.Index {
			if err := createIndex(db, table, item, false); err != nil {
				return err
			}
		}
	}
	return nil
}

// 检查字段是否已存在
func columnExists(db *sql.DB, table, column string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = ?
		AND COLUMN_NAME = ?
	`
	var count int
	err := db.QueryRow(query, table, column).Scan(&count)
	return count > 0, err
}

// 添加字段到表
func addColumn(db *sql.DB, table string, item *schema.ExtendAttrItem) error {
	// 构造字段类型
	var columnType string
	switch enum.AttrValueType(item.AttrValueType) {
	case enum.AttrValueType_STRING:
		length := item.AttrValueLen
		if length == 0 {
			length = 512
		}
		columnType = fmt.Sprintf("VARCHAR(%d)", length)
	case enum.AttrValueType_INT:
		length := item.AttrValueLen
		if length <= 11 {
			columnType = "INTEGER"
		} else {
			columnType = "BIGINT"
		}
	case enum.AttrValueType_FLOAT:
		columnType = "DOUBLE"
	case enum.AttrValueType_BOOL:
		columnType = "BOOLEAN"
	default:
		return fmt.Errorf("unknown AttrValueType: %d", item.AttrValueType)
	}

	// 执行ALTER TABLE
	query := fmt.Sprintf(
		"ALTER TABLE `%s` ADD COLUMN `%s` %s",
		table,
		item.AttrName,
		columnType,
	)
	_, err := db.Exec(query)
	return err
}

func shortenIndexName(s string) string {
	parts := strings.Split(s, "_")
	n := len(parts)
	newParts := make([]string, n)

	for i := 0; i < n; i++ {
		part := parts[i]
		if i == n-4 || i == n-5 {
			// 倒数第4、5个子串，取前6位
			if len(part) > 6 {
				newParts[i] = part[:6]
			} else {
				newParts[i] = part
			}
		} else if i < n-5 {
			// 倒数第6及更前的子串，取第一个字符
			if len(part) > 0 {
				newParts[i] = string(part[0])
			} else {
				newParts[i] = ""
			}
		} else {
			// 其他子串保持不变
			newParts[i] = part
		}
	}

	return strings.Join(newParts, "_")
}

// 创建索引（unique参数控制是否创建唯一索引）
// 修改后的createIndex函数（整合字段类型判断）
func createIndex(db *sql.DB, table string, item *schema.ExtendAttrItem, unique bool) error {
	// 生成唯一索引名称（table_column_uniq 或 table_column_idx）
	indexType := "idx"
	if unique {
		indexType = "uniq"
	}
	indexName := fmt.Sprintf("%s_%s_%s",
		strings.ToLower(table),
		strings.ToLower(item.AttrName),
		indexType,
	)

	indexName = shortenIndexName(indexName)

	// 处理字段类型和长度
	columnDef := item.AttrName
	if enum.AttrValueType(item.AttrValueType) == enum.AttrValueType_STRING {
		// 计算实际使用的字段长度
		fieldLength := item.AttrValueLen
		if fieldLength == 0 {
			fieldLength = 512 // 使用默认长度
		}

		// 当长度超过191时使用前缀索引（兼容utf8mb4字符集）
		if fieldLength > 191 {
			columnDef = fmt.Sprintf("%s(191)", item.AttrName)
		}
	}

	// 构造索引类型语句
	var query string
	if unique {
		query = fmt.Sprintf(
			"CREATE UNIQUE INDEX `%s` ON `%s` (%s)",
			indexName, table, columnDef,
		)
	} else {
		query = fmt.Sprintf(
			"CREATE INDEX `%s` ON `%s` (%s)",
			indexName, table, columnDef,
		)
	}

	_, err := db.Exec(query)
	return err
}

func (a *OrderArchiver) copyHistoricalTableStructure(originalTableName string, archivingLog *schema.DataArchivingLog) (string, error) {
	targetTableName := fmt.Sprintf(originalTableName+"_%s_%s", archivingLog.SystemCode, archivingLog.BusinessCode)
	return targetTableName, dbutil.CopyTableStructure(a.applicationCfg.GetAppDB(), a.applicationCfg.GetCentralDB(), originalTableName, targetTableName, false)
}

func (a *OrderArchiver) createHistoricalTables(archivingLog *schema.DataArchivingLog) (
	historicalGroupTradeOrderTableName,
	historicalTradeActionLatestRespTableName,
	historicalTradeActionRespTableName,
	historicalTradeOrderTableName, extendHistoricalTradeOrderTableName, extendHistoricalTradeActionRespTableName string,
	err error) {

	// 进行表的复制（在这个过程中，如果目标表存在，则会被清理并重新创建）
	historicalGroupTradeOrderTableName, err = a.copyHistoricalTableStructure("group_trade_orders", archivingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createHistoricalTables")
		return
	}
	historicalTradeActionLatestRespTableName, err = a.copyHistoricalTableStructure("trade_action_latest_resps", archivingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createHistoricalTables")
		return
	}
	historicalTradeActionRespTableName, err = a.copyHistoricalTableStructure("trade_action_resps", archivingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createHistoricalTables")
		return
	}
	historicalTradeOrderTableName, err = a.copyHistoricalTableStructure("trade_orders", archivingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createHistoricalTables")
		return
	}

	extendAttrItems := a.applicationCfg.GetExtendAttrItems()

	// 创建TradeOrder扩展表
	tableNameTmpl := "trade_orders_extend_%s_%s"
	extendHistoricalTradeOrderTableName = fmt.Sprintf(tableNameTmpl, archivingLog.SystemCode, archivingLog.BusinessCode)
	err = dbutil.CreateSimpleTableWithSysID(a.applicationCfg.GetCentralDB(), extendHistoricalTradeOrderTableName, false)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createHistoricalTables")
		return
	}
	extendAttrItems = domain_cfg.ConfigExtendAttrItemsForTradeOrderExtending(extendAttrItems)
	err = ExtendTableColumn(a.applicationCfg.GetCentralDB(), extendHistoricalTradeOrderTableName, extendAttrItems)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createHistoricalTables")
		return
	}

	// 创建TradeActionResp扩展表
	tableNameTmpl = "trade_action_resps_extend_%s_%s"
	extendHistoricalTradeActionRespTableName = fmt.Sprintf(tableNameTmpl, archivingLog.SystemCode, archivingLog.BusinessCode)
	err = dbutil.CreateSimpleTableWithSysID(a.applicationCfg.GetCentralDB(), extendHistoricalTradeActionRespTableName, false)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createHistoricalTables")
		return
	}
	extendAttrItems = domain_cfg.ConfigExtendAttrItemsForTradeActionRespExtending(extendAttrItems)
	for _, extendAttrItem := range a.applicationCfg.GetTradeActionRespAttrItems() {
		item := &schema.ExtendAttrItem{}
		bean.Copy(extendAttrItem).To(item)
		extendAttrItems = append(extendAttrItems, item)
	}
	err = ExtendTableColumn(a.applicationCfg.GetCentralDB(), extendHistoricalTradeActionRespTableName, extendAttrItems)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createHistoricalTables")
		return
	}

	// 历史表增加两个个字段：f_archive_date、f_task_name
	extendAttrItems = []*schema.ExtendAttrItem{{
		AttrName:      "f_archiving_date",
		AttrValueType: int(enum.AttrValueType_STRING),
		AttrValueLen:  8,
		Index:         true,
	}, {
		AttrName:      "f_task_name",
		AttrValueType: int(enum.AttrValueType_STRING),
		AttrValueLen:  32,
	}}
	tables := []string{
		historicalGroupTradeOrderTableName,
		historicalTradeActionLatestRespTableName,
		historicalTradeActionRespTableName,
		historicalTradeOrderTableName, extendHistoricalTradeOrderTableName, extendHistoricalTradeActionRespTableName,
	}
	for _, table := range tables {
		err = ExtendTableColumn(a.applicationCfg.GetCentralDB(), table, extendAttrItems)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderArchiver::createHistoricalTables")
			return
		}
	}

	return
}
