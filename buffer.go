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

// LimitedBuffer 有限缓冲
type LimitedBuffer interface {
	io.ReadWriter
	IsFull() bool  // for write
	IsEmpty() bool // for read
	Reset()
	Capacity() int
}
