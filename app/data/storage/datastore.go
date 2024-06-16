package storage

import (
	"redislite/app/commands/parsers/utils"
	"redislite/app/data/storage/datatypes"
	"strings"
	"sync"
)

var datalock sync.RWMutex

type DataStore struct {
	items *map[string]DataItem
}

type DataItem struct {
	key      string
	dataType datatypes.DataType
	// Value    any
	Value          []byte
	ExpiryTimeNano int64
}

func CreateDataItem(key string, dataType datatypes.DataType, value []byte, expiryTimeNano int64) DataItem {
	return DataItem{
		key:            key,
		dataType:       dataType,
		Value:          value,
		ExpiryTimeNano: expiryTimeNano,
	}
}

func (di *DataItem) GetKey() string {
	return di.key
}

func (s *DataStore) HandleExpired() {
  delItems := []string{}
  for k,v := range (*s.items) {
    if utils.Expired(v.ExpiryTimeNano) {
      delItems = append(delItems, k)
    }
  }

  for index := 0; index < len(delItems); index ++ {
    delete((*s.items), delItems[index])
  }
}

func CreateKvStore() DataStore {
	dataitems := make(map[string]DataItem)
	return DataStore{
		items: &dataitems,
	}
}

func (s *DataStore) LoadFromDb(data map[string]DataItem) {
	datalock.Lock()
	for k, v := range data {
		// key := strings.ToUpper(k)
		di := &DataItem{
			key:      k,
			dataType: datatypes.DATA_TYPE_STRING,
			Value:    []byte(v.Value),
			// ExpiryTimeNano: v.ExpiryTimeNano,
		}
		(*s.items)[k] = *di
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

// func (s *DataStore) GetKvStream(key string) (kvstream.KvStream, bool) {
// 	datalock.RLock()
// 	di, ok := (*s.items)[key]
// 	datalock.RUnlock()
// 	return di.value.(kvstream.KvStream), ok
// }
