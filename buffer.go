package limitedbuffer

import (
	"errors"
	"io"
	"sync"
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

// WithSync Sync Limited Buffer.
type WithSync struct {
	lb LimitedBuffer
	mu sync.Mutex
}

func (w *WithSync) Read(p []byte) (n int, err error) {
	w.mu.Lock()
	n, err = w.lb.Read(p)
	w.mu.Unlock()
	return
}

func (w *WithSync) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	n, err = w.lb.Write(p)
	w.mu.Unlock()
	return
}

// IsFull if full cannot write
func (w *WithSync) IsFull() bool {
	w.mu.Lock()
	ret := w.lb.IsFull()
	w.mu.Unlock()
	return ret
}

// IsEmpty if empty cannot read
func (w *WithSync) IsEmpty() bool {
	w.mu.Lock()
	ret := w.lb.IsEmpty()
	w.mu.Unlock()
	return ret
}

// Reset reset read and write position to 0
func (w *WithSync) Reset() {
	w.mu.Lock()
	w.lb.Reset()
	w.mu.Unlock()
}

// Capacity return fixed buffer size
func (w *WithSync) Capacity() int {
	w.mu.Lock()
	c := w.lb.Capacity()
	w.mu.Unlock()
	return c
}

// NewCycleBuffer New CyclyBuffer
func NewCycleBuffer(capacity int) LimitedBuffer {
	return newCycleBuffer(capacity)
}

// NewSyncCycleBuffer New Sync CycleBuffer
func NewSyncCycleBuffer(capacity int) LimitedBuffer {
	return &WithSync{
		mu: sync.Mutex{},
		lb: NewCycleBuffer(capacity),
	}
}
