package datatypes

type DataType byte

const (
	DATA_TYPE_NONE   DataType = 0x00
	DATA_TYPE_STRING DataType = 0x01
	DATA_TYPE_STREAM DataType = 0x02
)
