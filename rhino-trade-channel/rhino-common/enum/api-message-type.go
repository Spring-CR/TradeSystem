package enum

type ApiMessageType int

const (
	ApiMessageType_NewOrderSingle     ApiMessageType = 0
	ApiMessageType_OrderCancelRequest ApiMessageType = 1
)
