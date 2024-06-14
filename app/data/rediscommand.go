package data

type RedisCommand struct {
	RawCommand   string
	Command      string
	ParamLength  int
	Params       []string
	MessageBytes int64
}
