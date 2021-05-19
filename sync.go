package limitedbuffer

import "sync"

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

// NewSyncCycleBuffer New Sync CycleBuffer
func NewSyncCycleBuffer(capacity int) LimitedBuffer {
	return &WithSync{
		mu: sync.Mutex{},
		lb: NewCycleBuffer(capacity),
	}
}
