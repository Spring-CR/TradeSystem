package handlers

import (
	"rhino-api/api/api_const"
	"rhino-common/server/middleware"
	"rhino-common/utils/request"
	domain_asset_unit "rhino-instr/domain/asset-unit"
	"rhino-instr/schema"

	"github.com/gin-gonic/gin"
)

func CreateAssetUnit(c *gin.Context) {
	opt := &schema.AssetUnit{}
	if !middleware.BindInputOption(c, opt) {
		return
	}
	de := domain_asset_unit.CreateAssetUnit(opt.AccountNo, opt.AccountName, opt.CombiNo, opt.CombiName)
	if middleware.ProcessDomainError(de, c) {
		return
	}

}

func FindAllAssetUnits(c *gin.Context) {
	assetUnits, de := domain_asset_unit.FindAllAssetUnits()
	if middleware.ProcessDomainError(de, c) {
		return
	}
	if len(assetUnits) == 0 {
		assetUnits = make([]*schema.AssetUnit, 0)
	}
	middleware.ResponseJson(c, assetUnits)
}

func DeleteAssetUnit(c *gin.Context) {
	accountNo, de := request.GetQueryAsString(c, api_const.ParamAccountNo, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	combiNo, de := request.GetQueryAsString(c, api_const.ParamCombiNo, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	de = domain_asset_unit.DeleteAssetUnit(accountNo, combiNo)
	if middleware.ProcessDomainError(de, c) {
		return
	}
}