package kvstream

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type KvStream struct {
	StreamData *[]KvStreamData
}

type KvStreamData struct {
	Id    KvStreamId
	Items []KvPair
}

type KvPair struct {
  Key            string
	Value          string
}

type KvStreamId struct {
	Timestamp int64
	Sequence  int64
}

func (kvid KvStreamId) ToString() string {
  return createIdString(kvid.Timestamp, kvid.Sequence) 
}
// 1 greater than, 0 equal to, -1 less than
func (kvid KvStreamId) Compare(id KvStreamId) int {

	// equal
	if kvid.Timestamp == id.Timestamp && kvid.Sequence == id.Sequence {
		return 0
	}

	//greater than
	if kvid.Timestamp > id.Timestamp || (kvid.Timestamp == id.Timestamp && kvid.Sequence > id.Sequence) {
    return 1
	}

	//less than
	return -1
}

// func (kvs *KvStream) GetRangeFromStream(startts int64, startseq int64, endts int64, endseq int64) {
//
// }
//
// func (kvs *KvStream) GetStartRangeFromStream(endts int64, endseq int64) {
//
// }
//
// func (kvs *KvStream) GetEndRangeFromStream(startts int64, startseq int64) {
//
// }

func (kvs *KvStream) AppendItemAutoId(data []KvPair) (KvStreamData, error) {
	timestamp := time.Now().UTC().UnixMilli()
	return kvs.AppendItemAutoSequence(timestamp, data)
}

func (kvs *KvStream) AppendItemAutoSequence(msTimestamp int64, data []KvPair) (KvStreamData, error) {
	kvlen := len(*kvs.StreamData)
	if kvlen > 0 {
		lastElement := (*kvs.StreamData)[kvlen-1]
		if lastElement.Id.Timestamp == msTimestamp {
      kvsd := KvStreamData{
				Id: KvStreamId{
					Timestamp: msTimestamp,
					Sequence:  lastElement.Id.Sequence + 1,
				},
				Items: []KvPair{},
			}
			*kvs.StreamData = append(*kvs.StreamData, kvsd)
			return kvsd, nil
		}
	}

	var seq int64 = 0
	if msTimestamp == 0 {
		seq = 1
	}
	//if nothing exists or the timestamps don't match just create a new one
  kvsd := KvStreamData{
		Id: KvStreamId{
			Timestamp: msTimestamp,
			Sequence:  seq,
		},
		Items: []KvPair{},
	}
	*kvs.StreamData = append(*kvs.StreamData, kvsd)

	return kvsd, nil

}

func (kvs *KvStream) AppendItemFullExplicitId(msTimestamp int64, seq int64, data []KvPair) (KvStreamData, error) {
	kvlen := len(*kvs.StreamData)
	if msTimestamp == 0 && seq == 0 {
		return KvStreamData{}, errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}

	log.Println("len: ", kvlen)
	if kvlen == 0 {
    kvsd := KvStreamData{
			Id: KvStreamId{
				Timestamp: msTimestamp,
				Sequence:  seq,
			},
			Items: data,
		}
		*kvs.StreamData = append(*kvs.StreamData, kvsd)
		return kvsd, nil
	}

	lastElement := (*kvs.StreamData)[kvlen-1]
	log.Println(lastElement)
	if msTimestamp < lastElement.Id.Timestamp || (msTimestamp == lastElement.Id.Timestamp && seq <= lastElement.Id.Sequence) {
		return KvStreamData{}, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	log.Println("C")
  kvsd := KvStreamData{
		Id: KvStreamId{
			Timestamp: msTimestamp,
			Sequence:  seq,
		},
		Items: data,
	}
	*kvs.StreamData = append(*kvs.StreamData, kvsd)
	return kvsd, nil

}

func createIdString(timestamp int64, sequence int64) string {
	return fmt.Sprintf("%d-%d", timestamp, sequence)
}
