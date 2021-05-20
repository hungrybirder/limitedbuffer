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

// BufferStatus Buffer Status
type BufferStatus interface {
	Capacity() int
	UnreadSize() int
	FreeWriteSpace() int
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
