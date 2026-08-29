package domain_error

func init() {
	registerErrMsg(map[string]string{
		COLLECTION_NOT_FOUND_ERR_CODE: `找不到数据集:%s`,
	})
}
