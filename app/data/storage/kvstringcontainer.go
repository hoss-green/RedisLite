package storage

import (
	"strings"

	"redislite/app/data/datatypes/kvstring"
	"redislite/app/data/storage/datatypes"
)

func (s *DataStore) SetKvString(key string, dataObject kvstring.KvString) {
	k := strings.ToUpper(key)
	datalock.Lock()
	dataObject.Key = k
	di := &DataItem{
		key:      k,
		dataType: datatypes.DATA_TYPE_STRING,
		value:    dataObject,
	}
	(*s.items)[k] = *di
	datalock.Unlock()
}

func (s *DataStore) GetKvString(key string) (kvstring.KvString, bool) {
	k := strings.ToUpper(key)
	datalock.RLock()
	di, ok := (*s.items)[k]
	datalock.RUnlock()
	if !ok {
		return kvstring.KvString{}, false
	}
	return di.value.(kvstring.KvString), ok
}

func (s *DataStore) DelKvString(key string) (kvstring.KvString, bool) {
	k := strings.ToUpper(key)
	datalock.RLock()
	di, ok := (*s.items)[k]
	datalock.RUnlock()
	if !ok {
		return kvstring.KvString{}, false
	}
	return di.value.(kvstring.KvString), ok
}

func (s *DataStore) CountKvString() int {
	datalock.RLock()
	value := len(*s.items)
	datalock.RUnlock()
	return value
}
