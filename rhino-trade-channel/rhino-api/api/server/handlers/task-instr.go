package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"rhino-api/api/api_const"
	"rhino-api/api/options"
	"rhino-common/context"
	"rhino-common/domain_error"
	"rhino-common/server/middleware"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/json_extend"
	"rhino-common/utils/request"
	domain_task_instr "rhino-instr/domain/task-instr"
	"rhino-instr/store"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 下达指令
func IssueInstruction(c *gin.Context) {
	
	opt := &options.IssueInstructionOption{}
	if !middleware.BindInputOption(c, opt) {
		return
	}

	var datetimes []int64
	if opt.CreateTime > 0 {
		datetimes = append(datetimes, opt.CreateTime)
	}

	dailyInstrNo, indexDailyModify, batchSerialNo, dateNum, de := domain_task_instr.CreateInsertTaskInstr(
		opt.AccountNo, opt.CombiNo, opt.InstrType,
		opt.BeginDate, opt.EndDate, opt.BeginTime, opt.EndTime,
		opt.DirectOperator,
		opt.BusinessType,
		opt.LimitOperator,
		opt.Stocks,
		datetimes...,
	)

	if middleware.ProcessDomainError(de, c) {
		return
	}

	ret := &options.IssueInstructionRet{
		DailyInstrNo     : dailyInstrNo,
		IndexDailyModify : indexDailyModify,
		BatchSerialNo    : batchSerialNo,
		DateNum          : dateNum,
	}

	middleware.ResponseJson(c, ret)
}

// 指令查询
func FindInstruction(c *gin.Context) {

	opt := &options.CommonQueryOption{}
	if !middleware.BindInputOption(c, opt) {
		return
	}

	// 对指令编号，做特殊化处理
	for i, fc := range opt.FieldConditions {
		if fc.Field == "instr_code" {
			strs := strings.Split(fc.Value.(string), "-")
			if len(strs) != 3{
				de := domain_error.Build(domain_error.ILLEGAL_INSTR_CODE_ERR_CODE, nil)
				if middleware.ProcessDomainError(de, c) {
					return
				}
			}

			date, err := strconv.Atoi(strs[0])
			if err != nil {
				de := domain_error.Build(domain_error.ILLEGAL_INSTR_CODE_ERR_CODE, err)
				if middleware.ProcessDomainError(de, c) {
					return
				}
			}

			dailyInstrNo, err := strconv.Atoi(strs[1])
			if err != nil {
				de := domain_error.Build(domain_error.ILLEGAL_INSTR_CODE_ERR_CODE, err)
				if middleware.ProcessDomainError(de, c) {
					return
				}
			}

			indexDailyModify, err := strconv.Atoi(strs[1])
			if err != nil {
				de := domain_error.Build(domain_error.ILLEGAL_INSTR_CODE_ERR_CODE, err)
				if middleware.ProcessDomainError(de, c) {
					return
				}
			}

			opt.FieldConditions[i].Field = "date"
			opt.FieldConditions[i].ValueType = 0
			opt.FieldConditions[i].Value = date

			opt.FieldConditions = append(opt.FieldConditions, &dbutil.FieldCondition{
				Field: "daily_instr_no",
				ValueType: 0,
				Value: dailyInstrNo,
			})

			opt.FieldConditions = append(opt.FieldConditions, &dbutil.FieldCondition{
				Field: "index_daily_modify",
				ValueType: 0,
				Value: indexDailyModify,
			})

			break
		}
	}

	jsData, _ := json.MarshalIndent(opt, "", "  ")
	log.Printf("query_option:%s\n", jsData)

	result, total, de := domain_task_instr.FindTaskInstrs(opt.FieldConditions, opt.Limit, opt.Offset)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	for _, r := range result {
		r.InstrCode = fmt.Sprintf("%v-%v-%v", r.Date, r.DailyInstrNo, r.IndexDailyModify)
		r.InstStockCode = fmt.Sprintf("%v-%v-%v-%v", r.Date, r.DailyInstrNo, r.IndexDailyModify, r.StockSerialNo)
	}

	v := json_extend.TransformToJsonOfUnderscoreMap(result)
	pagingRecord := options.PagingRecord{Data: v, Total: total}
	middleware.ResponseJson(c, pagingRecord)
}

func StatisTaskInstrStock(c *gin.Context) {
	date, de := request.GetQueryAsInt(c, api_const.ParamDate, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	
	dailyInstrNo, de := request.GetQueryAsInt(c, api_const.ParamDailyInstrNo, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	
	indexDailyModify, de := request.GetQueryAsInt(c, api_const.ParamIndexDailyModify, false) 
	if middleware.ProcessDomainError(de, c) {
		return
	}
	
	stockSerialNo, de := request.GetQueryAsInt(c, api_const.ParamStockSerialNo, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	_, _, _, _, _, _, _, _, _, err := store.StatisTaskInstrStock(context.DB, date, int64(dailyInstrNo), int64(indexDailyModify), int64(stockSerialNo), true)
	if err != nil {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		if middleware.ProcessDomainError(de, c) {
			return
		}
	}
}