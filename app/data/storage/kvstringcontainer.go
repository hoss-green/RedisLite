package storage

import (
	"strings"

	"redislite/app/commands/parsers/utils"
	"redislite/app/data/storage/datatyperrors"
	"redislite/app/data/storage/datatypes"
)

func (s *DataStore) SetKvString(key string, dataObject DataItem) {
	s.SetKvStrings([]string{key}, []DataItem{dataObject})
}

func (s *DataStore) SetKvStrings(keys []string, dataObjects []DataItem) {
	datalock.Lock()
	for index := 0; index < len(keys); index++ {
		dataObject := dataObjects[index]
		k := strings.ToUpper(keys[index])
		dataObject.key = k
		di := &DataItem{
			key:            k,
			dataType:       datatypes.DATA_TYPE_STRING,
			ExpiryTimeNano: dataObject.ExpiryTimeNano,
			Value:          []byte(dataObject.Value),
		}
		(*s.items)[k] = *di
	}
	datalock.Unlock()
}

func (s *DataStore) GetKvString(key string) (DataItem, error) {
	k := strings.ToUpper(key)
	datalock.RLock()
	di, ok := (*s.items)[k]
	datalock.RUnlock()
	if utils.Expired(di.ExpiryTimeNano) {
		return DataItem{}, &datatyperrors.ExpiredKeyError{}
	}
	if !ok {
		return DataItem{}, &datatyperrors.KeyNotFoundError{}
	}
	if di.dataType != datatypes.DATA_TYPE_STRING {
		return DataItem{}, &datatyperrors.WrongtypeError{}
	}
	return di, nil
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
