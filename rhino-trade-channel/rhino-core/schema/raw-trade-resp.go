package schema

type RawTradeResp struct {
	ID               int64
	ClOrdID          string
	SecondaryClOrdID string
	MsgSeq           int64
	MsgTime          int64
	MsgReceivedTime  int64
	ConsumedFlag     int
	JsonMsg          string
}
