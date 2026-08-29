package position

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-core/order_domain/order_position_manager"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func (s *MemPosition) ProcessPositionChangeEvent(event *order_position_manager.PositionChangeEvent) {
	jsData, _ := json.Marshal(event)
	log.Printf("ProcessPositionChangeEvent: event=%s\n", jsData)
	switch event.InsertOrUpdate {
	case 0: // insert
		s.insertPosition(event.PositionData)
	case 1: // update
		s.updatePosition(event.PositionData)
	case 2: // update
		s.deletePosition(event.PositionData)
	}
}

func (s *MemPosition) insertPosition(positionData map[string]interface{}) {
	var args []interface{}
	for _, positionAttrItem := range s.positionAttrItems {
		val, ok := positionData[positionAttrItem.AttrName]
		if !ok {
			continue
		}
		args = append(args, val)
	}
	insertPositionSql := s.getPositionInsertSql(positionData)
	_, err := s.dbWrite.Exec(insertPositionSql, args...)
	if err != nil {
		log.Printf("===>insertPositionSql:%s\n", insertPositionSql)
		log.Printf("===>args:%v\n", args)
		log.Printf("===> error:%v\n", err)
		jsData, _ := json.Marshal(positionData)
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to insertPosition, positionData:%s\n", jsData))
	} else {
		jsData, _ := json.Marshal(positionData)
		log.Printf("success insert %s\n", jsData)
	}
}

func (s *MemPosition) updatePosition(positionData map[string]interface{}) {

	unqKey := ""
	for _, positionAttrItem := range s.positionAttrItems {
		if positionAttrItem.Unique {
			unqKey = positionAttrItem.AttrName
			break
		}
	}

	if unqKey == "" {
		return
	}

	keyValue, ok := positionData[unqKey]
	if !ok {
		return
	}

	updateSql := "UPDATE " + s.positionTableName + " SET"
	var args []interface{}
	for k, v := range positionData {
		if k == unqKey {
			continue
		}
		updateSql += fmt.Sprintf(" %s=?,", k)
		args = append(args, v)
	}

	updateSql = updateSql[:len(updateSql)-1] + " WHERE " + unqKey + "=?"
	args = append(args, keyValue)

	_, err := s.dbWrite.Exec(updateSql, args...)
	if err != nil {
		log.Printf("===>updateSql:%s\n", updateSql)
		log.Printf("===>args:%v\n", args)
		log.Printf("===> error:%v\n", err)
		jsData, _ := json.Marshal(positionData)
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to updatePosition, positionData:%s\n", jsData))
	} else {
		jsData, _ := json.Marshal(positionData)
		log.Printf("success update %s\n", jsData)
	}
}


func (s *MemPosition) deletePosition(positionData map[string]interface{}) {

	unqKey := ""
	for _, positionAttrItem := range s.positionAttrItems {
		if positionAttrItem.Unique {
			unqKey = positionAttrItem.AttrName
			break
		}
	}

	if unqKey == "" {
		log.Printf("no unqKey for positionData:%v\n", positionData)
		return
	}

	keyValue, ok := positionData[unqKey]
	if !ok {
		log.Printf("no keyValue for positionData:%v\n", positionData)
		return
	}

	deleteSql := "DELETE FROM " + s.positionTableName + " WHERE " + unqKey + "=?"
	var args []interface{}
	args = append(args, keyValue)

	_, err := s.dbWrite.Exec(deleteSql, args...)
	if err != nil {
		log.Printf("===>deleteSql:%s\n", deleteSql)
		log.Printf("===>args:%v\n", args)
		log.Printf("===> error:%v\n", err)
		jsData, _ := json.Marshal(positionData)
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to deletePosition, positionData:%s\n", jsData))
	} else {
		jsData, _ := json.Marshal(positionData)
		log.Printf("success delete %s\n", jsData)
	}
}
