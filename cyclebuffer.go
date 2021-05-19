package limitedbuffer

type cycleBuffer struct {
	capacity int    // max length
	cbuf     []byte // fixed buffer
	rpos     int    // reader position
	wpos     int    // writer position
}

func (c *cycleBuffer) isEmpty() bool {
	return c.rpos == c.wpos
}

func (c *cycleBuffer) isFull() bool {
	if c.rpos > 0 {
		return c.wpos+1 == c.rpos
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
	return
}

func (c *cycleBuffer) Read(p []byte) (n int, err error) {
	n, err = c.read(p)
	return
}

func (c *cycleBuffer) write(p []byte) (n int, err error) {
	if c.isFull() {
		n = 0
		err = ErrNotAvailableSpace
		return
	}

	n1 := 0
	n2 := 0
	tail := c.capacity - c.wpos

	n1 = copy(c.cbuf[c.wpos:c.capacity], p)
	if n1 == tail && c.rpos > 0 {
		n2 = copy(c.cbuf[0:c.rpos], p[n1:])
	}
	if n2 > 0 {
		c.wpos = n2
	} else {
		c.wpos += n1
	}
	n = n1 + n2
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
		capacity: capacity,
		cbuf:     make([]byte, capacity),
		rpos:     0,
		wpos:     0,
	}
	return cb
}
