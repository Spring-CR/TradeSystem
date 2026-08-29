package schema

type UtilFixMessage struct {
	ID          int64
	MsgSide     int    // 1-toAdmin, 2-toApp, 3-fromAdmin, 4-fromApp
	MsgType     string `sql:"size: 2"`
	MsgTime     int64
	Data        []byte
	ChannelCode string `sql:"size: 32"`
}
