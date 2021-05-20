package limitedbuffer

import "fmt"

// CycleBufferStatus cycleBuffer Status
type CycleBufferStatus struct {
	capacity       int
	unreadSize     int
	freeWriteSpace int

	totalRead  int64
	totalWrite int64
}

// Capacity 缓冲区大小
func (cb *CycleBufferStatus) Capacity() int {
	return cb.capacity
}

// UnreadSize 可读数据量
func (cb *CycleBufferStatus) UnreadSize() int {
	return cb.unreadSize
}

// FreeWriteSpace 可写数据量
func (cb *CycleBufferStatus) FreeWriteSpace() int {
	return cb.freeWriteSpace
}

// TotalRead 已读量
func (cb *CycleBufferStatus) TotalRead() int64 {
	return cb.totalRead
}

// TotalWrite 已写量
func (cb *CycleBufferStatus) TotalWrite() int64 {
	return cb.totalWrite
}

var _ BufferStatus = (*CycleBufferStatus)(nil)

type cycleBuffer struct {
	capacity int    // max length
	cbuf     []byte // fixed buffer
	rpos     int    // reader position
	wpos     int    // writer position

	totalRead  int64
	totalWrite int64
}

func (c *cycleBuffer) status() BufferStatus {
	var unreadSize int = 0
	var freeWriteSpace int = 0
	var capacity = c.capacity
	if c.isEmpty() {
		unreadSize = 0
	} else {
		if c.wpos > c.rpos {
			unreadSize = c.wpos - c.rpos
		} else {
			unreadSize = capacity - c.rpos + c.wpos
		}
	}
	if c.isFull() {
		freeWriteSpace = 0
	} else {
		if c.rpos == 0 || c.rpos == 1 {
			freeWriteSpace = capacity - c.wpos
		} else {
			if c.wpos == c.rpos {
				freeWriteSpace = capacity - 1
			} else if c.wpos > c.rpos {
				freeWriteSpace = capacity - c.wpos + c.rpos - 1
			} else {
				freeWriteSpace = c.rpos - c.wpos - 1
			}
		}
	}
	return &CycleBufferStatus{
		capacity:       capacity,
		unreadSize:     unreadSize,
		freeWriteSpace: freeWriteSpace,
		totalRead:      c.totalRead,
		totalWrite:     c.totalWrite,
	}
}

// Status 当前状态
func (c *cycleBuffer) Status() BufferStatus {
	return c.status()
}

func (c *cycleBuffer) String() string {
	status := c.status()
	return fmt.Sprintf(
		"<cycleBuffer(cap=%d rpos=%d wpos=%d unread=%d freewrite=%d totalRead=%d totalWrite=%d)> at %p",
		c.capacity, c.rpos, c.wpos,
		status.UnreadSize(), status.FreeWriteSpace(),
		c.totalRead, c.totalWrite, c)
}

func (c *cycleBuffer) isEmpty() bool {
	return c.rpos == c.wpos
}

// 如果 rpos == 1 or 0, wpos == capacity，则不能写
func (c *cycleBuffer) isFull() bool {
	if c.rpos > 1 {
		return c.wpos == c.rpos-1
	}
	return c.wpos == c.capacity
}

func (c *cycleBuffer) reset() {
	c.rpos = 0
	c.wpos = 0
}

func (c *cycleBuffer) IsFull() bool {
	ret := c.isFull()
	return ret
}

func (c *cycleBuffer) IsEmpty() bool {
	ret := c.isEmpty()
	return ret
}

func (c *cycleBuffer) Reset() {
	c.reset()
}

func (c *cycleBuffer) read(p []byte) (n int, err error) {
	if c.isEmpty() {
		n = 0
		err = ErrNotAvailableData
		return
	}

	n1 := 0
	n2 := 0
	if c.wpos > c.rpos {
		n1 = copy(p, c.cbuf[c.rpos:c.wpos])
	} else {
		n1 = copy(p, c.cbuf[c.rpos:c.capacity])
		tail := c.capacity - c.rpos
		if n1 == tail && c.wpos > 0 {
			n2 = copy(p[n1:], c.cbuf[0:c.wpos])
		}
	}

	if n2 > 0 {
		c.rpos = n2
	} else {
		c.rpos += n1
	}
	n = n1 + n2
	if n > 0 {
		c.totalRead += int64(n)
	}
	return
}

func (c *cycleBuffer) Read(p []byte) (n int, err error) {
	n, err = c.read(p)
	return
}

// 写操作，最多能循环写到 rpos - 1
//
// 如果 rpos == 0， wpos == 0
// 最多可以写到 capacity， 写完后 wpos == capacity
//
// 如果 rpos > 1, wpos >= rpos
// 最多可以写 capacity - wpos + rpos - 1，写完后 wpos = rpos - 1
func (c *cycleBuffer) write(p []byte) (n int, err error) {
	if c.isFull() {
		n = 0
		err = ErrNotAvailableSpace
		return
	}

	if c.isEmpty() && c.wpos == c.capacity {
		c.reset()
	}

	n1 := 0
	n2 := 0
	tail := c.capacity - c.wpos

	n1 = copy(c.cbuf[c.wpos:c.capacity], p)
	if n1 == tail && c.rpos > 1 {
		n2 = copy(c.cbuf[0:c.rpos-1], p[n1:])
	}
	if n2 > 0 {
		c.wpos = n2
	} else {
		c.wpos += n1
	}
	n = n1 + n2
	if n > 0 {
		c.totalWrite += int64(n)
	}
	return
}

func (c *cycleBuffer) Write(p []byte) (n int, err error) {
	n, err = c.write(p)
	return
}

func (c *cycleBuffer) Capacity() int {
	return c.capacity
}

var _ LimitedBuffer = (*cycleBuffer)(nil)

func newCycleBuffer(capacity int) *cycleBuffer {
	cb := &cycleBuffer{
		capacity:   capacity,
		cbuf:       make([]byte, capacity),
		rpos:       0,
		wpos:       0,
		totalRead:  0,
		totalWrite: 0,
	}
	return cb
}

// NewCycleBuffer New CyclyBuffer
func NewCycleBuffer(capacity int) LimitedBuffer {
	return newCycleBuffer(capacity)
}
