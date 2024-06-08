package storage

import (
	// "log"
	"strings"
	"sync"

	"redislite/app/data/datatypes/kvstream"
	"redislite/app/data/datatypes/kvstring"
	"redislite/app/data/storage/datatypes"
)

var datalock sync.RWMutex

type DataStore struct {
	items *map[string]DataItem
}

type DataItem struct {
	key      string
	dataType datatypes.DataType
	value    any
}

func CreateKvStore() DataStore {
	dataitems := make(map[string]DataItem)
	return DataStore{
		items: &dataitems,
	}
}

func (s *DataStore) LoadFromDb(data map[string]kvstring.KvString) {
	datalock.Lock()
	for k, v := range data {
    key := strings.ToUpper(k)
		di := &DataItem{
			key:      key,
			dataType: datatypes.DATA_TYPE_STRING,
			value:    v,
		}
		(*s.items)[key] = *di
	}
	datalock.Unlock()
	// kvstringlock.Unlock()
}

func (s *DataStore) KvType(key string) string {
	k := strings.ToUpper(key)
	datalock.RLock()
	di, ok := (*s.items)[k]
	datalock.RUnlock()
	if !ok {
		return "none"
	}

	switch di.dataType {
	case datatypes.DATA_TYPE_STRING:
		return "string"
	case datatypes.DATA_TYPE_STREAM:
		return "stream"
	default:
		return "none"
	}
}

func (s *DataStore) GetKvStream(key string) (kvstream.KvStream, bool) {
	datalock.RLock()
	di, ok := (*s.items)[key]
	datalock.RUnlock()
	return di.value.(kvstream.KvStream), ok
}
