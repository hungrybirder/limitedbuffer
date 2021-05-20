package limitedbuffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	b1  = []byte("X")
	b2  = []byte("YY")
	b3  = []byte("ZZZ")
	b4  = []byte("VVVV")
	b8  = []byte("ABCDEFGH")
	b10 = []byte("0123456789")
)

func makeByteSlice(n int) []byte {
	return make([]byte, n)
}

func TestNewCycleBuffer(t *testing.T) {
	assert := assert.New(t)
	const bufSize = 8
	cb := NewCycleBuffer(bufSize)
	status := cb.Status()
	assert.Equal(bufSize, status.Capacity())
}

func TestCycleBufferReadWrite(t *testing.T) {
	assert := assert.New(t)
	const bufSize = 8
	var status BufferStatus
	var err error
	var wn int
	var rn int
	cb := NewCycleBuffer(bufSize)
	status = cb.Status()
	assert.Equal(bufSize, status.Capacity())
	t.Logf("#0 %v", cb)

	status = cb.Status()
	wn, err = cb.Write(b10)
	assert.Nil(err)
	assert.Equal(status.Capacity(), wn)
	t.Logf("#1 %v", cb)

	var r1 = makeByteSlice(1)
	rn, err = cb.Read(r1) // read '0'
	assert.Nil(err)
	assert.Equal(rn, len(r1))
	assert.Equal(r1[0], b10[0])
	t.Logf("#2 %v", cb)

	// cannot write
	_, err = cb.Write(b10)
	assert.NotNil(err)

	var r3 = makeByteSlice(3)
	rn, err = cb.Read(r3)
	assert.Nil(err)
	assert.Equal(rn, len(r3)) // read '123'
	assert.Equal(r3[0:3], b10[1:4])
	t.Logf("#3 %v", cb)

	status = cb.Status()
	wn, err = cb.Write(b10)
	assert.Nil(err)
	assert.Equal(status.FreeWriteSpace(), wn)
	assert.True(cb.IsFull())
	t.Logf("#4 %v", cb)

	var r12 = makeByteSlice(12)
	status = cb.Status()
	rn, err = cb.Read(r12)
	assert.Nil(err)
	assert.Equal(status.UnreadSize(), rn)
	t.Logf("#5 %v", cb)
}

func TestCycleBufferReadWrite2(t *testing.T) {
	assert := assert.New(t)
	const bufSize = 8
	var status BufferStatus
	var err error
	var wn int
	// var rn int
	cb := NewCycleBuffer(bufSize)
	status = cb.Status()
	assert.Equal(bufSize, status.Capacity())
	t.Logf("NewCycleBuffer(%d), %v", bufSize, cb)

	t.Logf("Before Write 3 Bytes to %v", cb)
	cb.Write([]byte("ABC"))
	// cb.Write([]byte("A"))
	t.Logf("After Write 3 Bytes to %v", cb)

	var r3 = makeByteSlice(3)
	cb.Read(r3)
	t.Logf("After read 3 byte from %v", cb)
	status = cb.Status()

	t.Logf("Try Write 8 Bytes to %v", cb)
	wn, err = cb.Write([]byte("12345678"))
	t.Logf("After Try Write 8 Bytes to %v", cb)
	assert.Nil(err)
	assert.Equal(status.FreeWriteSpace(), wn)
}
