package domain_asset_unit

import (
	"rhino-common/context"
	"rhino-common/domain_error"
	"rhino-common/utils/dbutil"
	"rhino-instr/schema"
	"rhino-instr/store"
)

func CreateAssetUnit(accountNo string, accountName string, combiNo string, combiName string) (de *domain_error.Error) {
	au := &schema.AssetUnit{AccountNo: accountNo, AccountName: accountName, CombiNo: combiNo, CombiName: combiName}
	err := store.InsertAssetUnit(context.DB, au)
	if err != nil {
		de = domain_error.Build(domain_error.CANNOT_CREATE_ASSET_UNIT_ERR_CODE, err, accountNo, accountName, combiNo, combiName)
		return
	}
	return
}

func FindAllAssetUnits() (result []*schema.AssetUnit, de *domain_error.Error) {

	var err error
	result, err = store.FindAllAssetUnits(context.DB)
	if dbutil.IsDbRecordEmptyError(err) {
		err = nil
	}

	if err != nil {
		de = domain_error.Build(domain_error.CANNOT_FIND_ALL_ASSET_UNITS_ERR_CODE, err)
		return
	}

	return
}

func DeleteAssetUnit(accountNo string, combiNo string) (de *domain_error.Error) {
	err := store.DeleteAssetUnitByAccountNoAndCombiNo(context.DB, accountNo, combiNo)
	if err != nil {
		de = domain_error.Build(domain_error.CANNOT_DELETE_ASSET_UNIT_ERR_CODE, err, accountNo, combiNo)
		return
	}
	return
}
