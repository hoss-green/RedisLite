package storage

import (
	"strings"

	"redislite/app/data/datatypes/kvstring"
	"redislite/app/data/storage/datatyperrors"
	"redislite/app/data/storage/datatypes"
)

func (s *DataStore) SetKvString(key string, dataObject kvstring.KvString) {
	s.SetKvStrings([]string{key}, []kvstring.KvString{dataObject})
}

func (s *DataStore) SetKvStrings(keys []string, dataObjects []kvstring.KvString) {
	datalock.Lock()
	for index := 0; index < len(keys); index++ {
		dataObject := dataObjects[index]
		k := strings.ToUpper(keys[index])
		dataObject.Key = k
		di := &DataItem{
			key:      k,
			dataType: datatypes.DATA_TYPE_STRING,
			value:    dataObject,
		}
		(*s.items)[k] = *di
	}
	datalock.Unlock()
}

func (s *DataStore) GetKvString(key string) (kvstring.KvString, error) {
	k := strings.ToUpper(key)
	datalock.RLock()
	di, ok := (*s.items)[k]
	datalock.RUnlock()
	if !ok {
		return kvstring.KvString{}, &datatyperrors.KeyNotFoundError{}
	}
	if di.dataType != datatypes.DATA_TYPE_STRING {
		return kvstring.KvString{}, &datatyperrors.WrongtypeError{}
	}
	return di.value.(kvstring.KvString), nil
}

func (s *DataStore) DelKvString(key string) {
	k := strings.ToUpper(key)
	datalock.Lock()
	delete((*s.items), k)
	datalock.Unlock()
}

func (s *DataStore) CountKvString() int {
	datalock.RLock()
	value := len(*s.items)
	datalock.RUnlock()
	return value
}
