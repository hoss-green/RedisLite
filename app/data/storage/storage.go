package storage

import (
	"sync"

	"redislite/app/data/datatypes/kvstring"
)

type Storage struct {
	// KvStrings KvStringContainer
	KvStreams KvStreamContainer // map[string]kvstream.KvStream
}

var kvstreamlock sync.RWMutex
var kvstringlock sync.RWMutex
var setup bool

func CreateDataStorage() Storage {
	if setup {
		panic("CreateDataStorage called twice")
	}

	return Storage{
    // KvStore: CreateKvStore(),
    // kvContainer: KvContainer{},
		// KvStrings: CreateKvStringContainer(),
		KvStreams: CreateKvStreamContainer(),
	}
}

func (s *Storage) LoadFromDb(data map[string]kvstring.KvString) {
	// kvstringlock.Lock()
	// s.KvStrings = CreateKvStringContainerWithData(data)
	s.KvStreams = CreateKvStreamContainer()
	// kvstringlock.Unlock()
}

