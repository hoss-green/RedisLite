package storage

import (
	"strings"
)

func (s *DataStore) DeleteKey(key string) int64 {
	return s.DeleteKeys([]string{key})
}

func (s *DataStore) DeleteKeys(keys []string) int64 {
	var count int64 = 0
	for _, key := range keys {
		k := strings.ToUpper(key)
		datalock.RLock()
		_, ok := (*s.items)[k]
		if ok {
			count += 1
		}
		datalock.RUnlock()
		datalock.Lock()
		delete((*s.items), k)
		datalock.Unlock()
	}

	return count
}

func (s *DataStore) CopyKey(sourceKey string, targetKey string) int64 {
	k := strings.ToUpper(sourceKey)
	targetKey = strings.ToUpper(targetKey)
	datalock.RLock()
	source, ok := (*s.items)[k]
	datalock.RUnlock()
	if !ok {
		return 0
	}
  di := source
  datalock.Lock()
	(*s.items)[targetKey] = di
  datalock.Unlock()
	return 1
}

func (s *DataStore) ExistsKey(sourceKey string) int64 {
	k := strings.ToUpper(sourceKey)
	datalock.RLock()
	_, ok := (*s.items)[k]
	datalock.RUnlock()
	if !ok {
		return 0
	}
	return 1
}

func (s *DataStore) ExistKey(key string) int64 {
	return s.ExistKeys([]string{key})
}

func (s *DataStore) ExistKeys(keys []string) int64 {
	var count int64 = 0
	for _, key := range keys {
		k := strings.ToUpper(key)
		datalock.RLock()
		_, ok := (*s.items)[k]
		if ok {
			count += 1
		}
		datalock.RUnlock()
	}

	return count
}
