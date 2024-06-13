package datatyperrors

type WrongtypeError struct {
}

func (e *WrongtypeError) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

type KeyNotFoundError struct {
}

func (e *KeyNotFoundError) Error() string {
	return "Key Not Found Error"
}
