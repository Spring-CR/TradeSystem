package position

import (
	"database/sql"
	"log"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
)

type MemPosition struct {
	systemCode                 string
	businessCode               string
	applicationCfg             *domain_cfg.ApplicationCfg
	positionAttrItems          []*schema.PositionAttrItem
	dbConfig                   string
	createPositionTableInitSql string
	//insertPositionSql          string
	positionTableName          string
	dbWrite                    *sql.DB
	dbRead                     *sql.DB
}

func NewMemPosition(applicationCfg *domain_cfg.ApplicationCfg) *MemPosition {
	log.Println("======> start to NewMemPosition")
	inst := &MemPosition{}
	inst.applicationCfg = applicationCfg
	inst.initMemDb()
	return inst
}

// func (s *MemPosition) Reset(){
// 	s.initMemDb()
// 	s.applicationCfg.GetAutoSyncRepo().Reset()
// }
