package schema

type TradeAlgorithm struct {
	ID              int64
	ChannelCode     string `sql:"unique: pk_ta, size: 32"`
	AlgorithmCode   string `sql:"index: pk_ta, size: 32"`
	AlgorithmZhName string `sql:"size: 128"`
	AlgorithmEnName string `sql:"size: 64"`
	Description     string
}
