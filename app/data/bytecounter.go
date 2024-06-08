package data

import (
	"log"
	"sync"
)

var bLock sync.RWMutex

type ByteCounter struct {
	TotalBytes int64
}

func (b *ByteCounter) AddBytes(numbytes int64) {
  bLock.Lock()
	b.TotalBytes += numbytes

  log.Printf("Added Bytes: %d", numbytes)
  log.Printf("Total Count: %d", b.TotalBytes)
  bLock.Unlock()
}

func (b *ByteCounter) Total() int64 {
  var total int64 =  0
  bLock.RLock()
  total = b.TotalBytes
  log.Printf("Total Ret: %d", b.TotalBytes)
  bLock.RUnlock()
	return total
}
