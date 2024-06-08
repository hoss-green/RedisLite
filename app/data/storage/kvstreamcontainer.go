package storage

import (
	"errors"
	// "log"
	"strconv"
	"strings"

	"redislite/app/data/datatypes/kvstream"
	"redislite/app/data/datatypes/kvstring"
	"redislite/app/data/storage/datatypes"
)

type KvStreamContainer struct {
	Items map[string]kvstream.KvStream
}

func CreateKvStreamContainer() KvStreamContainer {
	return KvStreamContainer{Items: (make(map[string]kvstream.KvStream))}
}

func (s *DataStore) appendStreamAutoId(streamkey string, data []kvstring.KvString) (kvstream.KvStreamData, error) {
	str := s.getOrCreateStream(streamkey)
	kvstreamlock.Lock()
	kvsd, err := str.AppendItemAutoId(data)
	kvstreamlock.Unlock()
	return kvsd, err
}

func (s *DataStore) appendStreamAutoSequencing(streamkey string, id int64, data []kvstring.KvString) (kvstream.KvStreamData, error) {
	str := s.getOrCreateStream(streamkey)
	kvstreamlock.Lock()
	kvsd, err := str.AppendItemAutoSequence(id, data)
	kvstreamlock.Unlock()
	return kvsd, err
}

func (s *DataStore) appendStreamFullExplicitId(streamkey string, id int64, seq int64, data []kvstring.KvString) (kvstream.KvStreamData, error) {
	str := s.getOrCreateStream(streamkey)
	kvstreamlock.Lock()
	kvsd, err := str.AppendItemFullExplicitId(id, seq, data)
	kvstreamlock.Unlock()
	return kvsd, err
}

func (s *DataStore) getOrCreateStream(streamkey string) *kvstream.KvStream {
	streamkey = strings.ToUpper(streamkey)
	kvs, ok := s.getKvStream(streamkey)
	if !ok {
		kvs = s.createKvStream(streamkey)
	}

	return &kvs
}

func (s *DataStore) createKvStream(key string) kvstream.KvStream {
	k := strings.ToUpper(key)
	kvs := kvstream.KvStream{
		StreamData: &[]kvstream.KvStreamData{},
	}
	datalock.Lock()
	di := &DataItem{
		key:      k,
		dataType: datatypes.DATA_TYPE_STREAM,
		value:    kvs,
	}
	(*s.items)[k] = *di
	datalock.Unlock()
	return kvs
}

func (s *DataStore) getKvStream(key string) (kvstream.KvStream, bool) {
	k := strings.ToUpper(key)
	datalock.RLock()
	di, ok := (*s.items)[k]
	datalock.RUnlock()
	if !ok || di.dataType == datatypes.DATA_TYPE_STRING {
		return kvstream.KvStream{}, false
	}

	return di.value.(kvstream.KvStream), ok
}

// check the key is correctly formatted and append the stream with the key
func (s *DataStore) AppendKvStream(streamkey string, key string, data []kvstring.KvString) (kvstream.KvStreamData, error) {
	//fully autogenerate key
	if key == "*" {
		return s.appendStreamAutoId(streamkey, data)
	}

	//check the key is in the correct format
	idparts := strings.Split(key, "-")
	if len(idparts) != 2 || len(idparts[0]) == 0 || len(idparts[1]) == 0 {
		return kvstream.KvStreamData{}, errors.New("wrong number of arguments for command")
	}

	id, err := strconv.ParseInt(idparts[0], 10, 64)
	if err != nil {
		return kvstream.KvStreamData{}, errors.New("wrong number of arguments for command")
	}

	//if the sequence is generated
	if idparts[1] == "*" {
		return s.appendStreamAutoSequencing(streamkey, id, data)
	}

	seq, err := strconv.ParseInt(idparts[1], 10, 64)
	if err != nil {
		return kvstream.KvStreamData{}, errors.New("wrong number of arguments for command")
	}

	//if the key is totally explicit
	return s.appendStreamFullExplicitId(streamkey, id, seq, data)
}
func (s *DataStore) GetStreamRead(streamkey string, startidstring string, inclusive bool) ([]kvstream.KvStreamData, error) {
	streamdata, err := s.GetStreamRange(streamkey, startidstring, "+")
	if err != nil {
		return []kvstream.KvStreamData{}, err
	}

	startId, err := parseKvStreamId(startidstring)
	if err != nil {
		return []kvstream.KvStreamData{}, err
	}

	if len(streamdata) == 0 {
		return []kvstream.KvStreamData{}, nil
	}

	// if inclusive && streamdata[0].Id.Compare(startId) == 0 {
	// 	return streamdata[0:], nil
	// }

	if streamdata[0].Id.Compare(startId) <= 0 {
		if inclusive {
			return streamdata, nil
		}
		return streamdata[1:], nil
	}
	return streamdata, nil
}
func (s *DataStore) GetStreamRange(streamkey string, startidstring string, endidstring string) ([]kvstream.KvStreamData, error) {
	datastream, exists := s.getKvStream(streamkey)
	if !exists {
		return []kvstream.KvStreamData{}, nil
	}

	foundstart := startidstring == "-"
	foundend := endidstring == "+"

	var err error
	var startId kvstream.KvStreamId
	var endId kvstream.KvStreamId
	//if we have both range keys just return the whole array
	if foundstart && foundend {
		return *datastream.StreamData, nil
	}

	startindex := 0
	endindex := len(*datastream.StreamData)

	if !foundstart {
		startId, err = parseKvStreamId(startidstring)
		if err != nil {
			return []kvstream.KvStreamData{}, err
		}
	}

	if !foundend {
		endId, err = parseKvStreamId(endidstring)
		if err != nil {
			return []kvstream.KvStreamData{}, err
		}
	}

searchloop:
	for index, kvsdata := range *datastream.StreamData {
		// log.Println(kvsdata)
		if !foundstart {
			if kvsdata.Id.Compare(startId) <= 0 {
				startindex = index
			} else {
				foundstart = true
			}
		}

		if !foundend {
			if kvsdata.Id.Compare(endId) > 0 {
				endindex = index
				foundend = true
				break searchloop
			}
		}
	}

	//if the key is totally explicit
	return (*datastream.StreamData)[startindex:endindex], nil
}

func parseKvStreamId(input string) (kvstream.KvStreamId, error) {
	idparts := strings.Split(input, "-")
	if len(idparts) != 2 || len(idparts[0]) == 0 || len(idparts[1]) == 0 {
		return kvstream.KvStreamId{}, errors.New("Cannot Parse Id")
	}

	timestamp, err := strconv.ParseInt(idparts[0], 10, 64)
	if err != nil {
		return kvstream.KvStreamId{}, errors.New("timestamp incorrectly formatted")
	}

	sequence, err := strconv.ParseInt(idparts[1], 10, 64)
	if err != nil {
		return kvstream.KvStreamId{}, errors.New("sequence id incorrectly formatted")
	}

	return kvstream.KvStreamId{Timestamp: timestamp, Sequence: sequence}, nil
}
