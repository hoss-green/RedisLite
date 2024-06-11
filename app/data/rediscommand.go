package data

type RedisCommand struct {
	RawCommand   string
	Command      string
	ParamLength  int
	Params       []string
	MessageBytes int64
}

type ParamOptions struct {
	// EX seconds -- Set the specified expire time, in seconds (a positive integer).
	// PX milliseconds -- Set the specified expire time, in milliseconds (a positive integer).
	// EXAT timestamp-seconds -- Set the specified Unix time at which the key will expire, in seconds (a positive integer).
	// PXAT timestamp-milliseconds -- Set the specified Unix time at which the key will expire, in milliseconds (a positive integer).
	// represents EX/PX/EXAT/PXAT
	Expiry int64
	// NX -- Only set the key if it does not already exist.
	// XX -- Only set the key if it already exists.
	// KEEPTTL -- Retain the time to live associated with the key.
	// GET -- Return the old string stored at key, or nil if key did not exist. An error is returned and SET aborted if the value stored at key is not a string.

}