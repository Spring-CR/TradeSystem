package enum

type SyncType int

const (
	SyncType_DSP                 SyncType = 0 // 大数据DSP接口同步
	SyncType_HTTP_HOOK           SyncType = 1 // http回调
	SyncType_CSV                 SyncType = 2 // 直接提供csv内容
	SyncType_PAGING_HTTP_HOOK    SyncType = 3 // 分页http回调
	SyncType_ITERATIVE_HTTP_HOOK SyncType = 4 // 迭代http回调
)

type DataSyncLogPhase int

const (
	DataSyncLogPhase_New        DataSyncLogPhase = 0
	DataSyncLogPhase_Processing DataSyncLogPhase = 1
	DataSyncLogPhase_Complete   DataSyncLogPhase = 2
	DataSyncLogPhase_Fail       DataSyncLogPhase = 3
	DataSyncLogPhase_Cancel     DataSyncLogPhase = 4
)
