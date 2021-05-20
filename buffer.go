package limitedbuffer

import (
	"errors"
	"io"
)

var (
	// ErrNotAvailableData 没有可读的数据，请忙写数据
	ErrNotAvailableData error = errors.New("NotAvailableData")

	// ErrNotAvailableSpace 没有可写空间，请尽快读数据
	ErrNotAvailableSpace error = errors.New("NotAvailableSpace")
)

// BufferStatus 状态
type BufferStatus struct {
	capacity       int
	unreadSize     int
	freeWriteSpace int
}

// Capacity 缓冲区大小
func (b BufferStatus) Capacity() int {
	return b.capacity
}

// UnreadSize 可读数据量
func (b BufferStatus) UnreadSize() int {
	return b.unreadSize
}

// FreeWriteSpace 可写数据量
func (b BufferStatus) FreeWriteSpace() int {
	return b.freeWriteSpace
}

// LimitedBuffer 有限缓冲
type LimitedBuffer interface {
	io.ReadWriter
	IsEmpty() bool        // for read
	IsFull() bool         // for write
	Capacity() int        // max buffer size
	Status() BufferStatus // return current buffer status
	Reset()               // reset buffer
}
