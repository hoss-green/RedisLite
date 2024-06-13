package getset 

// option 1
type setType int8

const (
	NormalSetType setType = 0
	NXSetIfKeyNotExist setType = 1 //NX
	XXOnlySetIfExists   setType = 2 //XX
)

//
