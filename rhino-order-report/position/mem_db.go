package position

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-core/schema"
)

func (s *MemPosition) initDBVar() {
	s.positionAttrItems = s.applicationCfg.GetPositionAttrItems()
	s.dbConfig = `
PRAGMA journal_mode = WAL;
PRAGMA busy_timeout = 5000;
PRAGMA synchronous = OFF;
PRAGMA cache_size = 1000000000;
PRAGMA foreign_keys = true;
PRAGMA temp_store = memory;`
// 	s.dbConfig = `
// `
	s.createPositionTableInitSql = `
CREATE TABLE position_%s_%s(
`
}

// https://www.runoob.com/sqlite/sqlite-data-types.html
func (s *MemPosition) initMemDb() {

	s.initDBVar()

	// 坑，配置prefix
	dbutil.ConfigDBFieldPrefix("")

	var err error
	//s.dbWrite, err = sql.Open("sqlite3", ":memory:?cache=shared")
	s.dbWrite, err = sql.Open("sqlite3", ":memory:?cache=shared&other=param1")
	//s.dbWrite, err = sql.Open("sqlite3", "/tmp/data.db")
	if err != nil {
		log.Fatal(err)
	}
	// 设置写db
	s.dbWrite.SetMaxOpenConns(1)

	s.systemCode, s.businessCode = s.applicationCfg.GetSystemAndBusinessCodes()

	createPositionTableInitSql := fmt.Sprintf(s.createPositionTableInitSql, s.systemCode, s.businessCode)
	s.positionTableName = fmt.Sprintf("position_%s_%s", s.systemCode, s.businessCode)

	initSql := bytes.NewBufferString(s.dbConfig + createPositionTableInitSql)

	
	for i, positionAttrItem := range s.positionAttrItems {
		if i != 0 {
			initSql.WriteString(",")
		}
		if positionAttrItem.Unique {
			initSql.WriteString(positionAttrItem.AttrName + " " + s.getDbFieldType(positionAttrItem.AttrValueType) + " UNIQUE \n")
		} else {
			initSql.WriteString(positionAttrItem.AttrName + " " + s.getDbFieldType(positionAttrItem.AttrValueType) + "\n")
		}
	}

	initSql.WriteString(");\n")

	for _, positionAttrItem := range s.positionAttrItems {
		// 创建索引
		if positionAttrItem.Index {
			initSql.WriteString("CREATE INDEX index_" + s.positionTableName + "_" + positionAttrItem.AttrName + " ON " + s.positionTableName + "(" + positionAttrItem.AttrName + ");\n")
		}
	}

	createPositionTableSqlStr := initSql.String()
	log.Printf("createPositionTableSqlStr:%s\n", createPositionTableSqlStr)

	// 初始化DB
	_, err = s.dbWrite.Exec(createPositionTableSqlStr)
	if err != nil {
		log.Fatal(err)
	}

	// // 设置读db
	// s.dbRead.SetMaxOpenConns(max(4, runtime.NumCPU()))
	s.dbRead = s.dbWrite // SQLite内存模式无法读写分离

	var version string
    err = s.dbRead.QueryRow("SELECT sqlite_version();").Scan(&version)
    if err != nil {
        log.Fatal("查询版本失败: ", err)
    }

    fmt.Printf("SQLite 版本: %s\n", version)

	log.Println("success create memory database!")
}

func (s *MemPosition) getPositionInsertSql(positionData map[string]interface{}) string{
	insertPositionSql := `INSERT INTO position_%s_%s(`
	insertPositionSql = fmt.Sprintf(insertPositionSql, s.systemCode, s.businessCode)
	sqlBuf := bytes.NewBufferString(insertPositionSql)

	var positionAttrItems []*schema.PositionAttrItem
	for _, positionAttrItem := range s.positionAttrItems {
		if _, ok := positionData[positionAttrItem.AttrName]; ok {
			positionAttrItems = append(positionAttrItems, positionAttrItem)
		}
	}

	unqIdx := -1
	for i, positionAttrItem := range positionAttrItems {
		if i != 0 {
			sqlBuf.WriteString(",")
		}
		if positionAttrItem.Unique {
			unqIdx = i
		}
		sqlBuf.WriteString(positionAttrItem.AttrName)
	}
	sqlBuf.WriteString(") VALUES (")
	n := len(positionAttrItems)
	for i := 0; i < n; i++ {
		sqlBuf.WriteString("?")
		if i == n-1 {
			sqlBuf.WriteString(")")
		} else {
			sqlBuf.WriteString(",")
		}
	}
	sqlBuf.WriteString("\nON CONFLICT (" + positionAttrItems[unqIdx].AttrName + ")\n")

	sqlBuf.WriteString("DO UPDATE SET ")

	for i, positionAttrItem := range positionAttrItems {
		if i == unqIdx {
			continue
		}
		sqlBuf.WriteString(positionAttrItem.AttrName + " = excluded."+positionAttrItem.AttrName+",")
	}

	insertPositionSql = sqlBuf.String()
	insertPositionSql = insertPositionSql[:len(insertPositionSql)-1]

	log.Printf("======>getPositionInsertSql:%s\n", insertPositionSql)
	return insertPositionSql
}

func (s *MemPosition) getDbFieldType(attrValueType int) string {
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
