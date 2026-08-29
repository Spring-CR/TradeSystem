package enum

// 属性值的类型枚举
type DataArchivingLogPhase int

const (
	DataArchivingLogPhase_New                DataArchivingLogPhase = 0
	DataArchivingLogPhase_Ready              DataArchivingLogPhase = 1
	DataArchivingLogPhase_CompleteDailyTable DataArchivingLogPhase = 2
	DataArchivingLogPhase_MergeToHisTable    DataArchivingLogPhase = 3
)

type DataPurgingLogPhase int

const (
	DataPurgingLogPhase_New            DataPurgingLogPhase = 0 // 新建记录
	DataPurgingLogPhase_Ready          DataPurgingLogPhase = 1 // 登记了需要删除的记录的主键
	DataPurgingLogPhase_ResetDB        DataPurgingLogPhase = 2 // 完成数据库的表truncate和记录重录
	DataPurgingLogPhase_ResetKafka     DataPurgingLogPhase = 3 // 完成Kafka消息队列重置
	DataPurgingLogPhase_ResetMem       DataPurgingLogPhase = 4 // 完成内存模型重置
	DataPurgingLogPhase_CustomizedTask DataPurgingLogPhase = 5 // 完成项目自定义的任务
)
