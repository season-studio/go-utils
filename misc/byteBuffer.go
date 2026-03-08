package misc

type ByteBuffer struct {
	buf       []byte
	blockSize int
}

func CreateByteBuffer(startSize int, blockSize int) *ByteBuffer {
	if startSize < 0 || blockSize <= 0 {
		return nil
	}
	return &ByteBuffer{
		buf:       make([]byte, 0, startSize),
		blockSize: blockSize,
	}
}

func (b *ByteBuffer) Reset() {
	b.buf = b.buf[:0]
}

func (b *ByteBuffer) Write(data byte) {
	blockSize := b.blockSize
	needed := len(b.buf) + 1
	if needed > cap(b.buf) {
		newCap := ((needed + blockSize - 1) / blockSize) * blockSize
		newBuf := make([]byte, len(b.buf), newCap)
		copy(newBuf, b.buf)
		b.buf = newBuf
	}
	b.buf = append(b.buf, data)
}

func (b *ByteBuffer) WriteBytes(data []byte) {
	blockSize := b.blockSize
	needed := len(b.buf) + len(data)
	if needed > cap(b.buf) {
		newCap := ((needed + blockSize - 1) / blockSize) * blockSize
		newBuf := make([]byte, len(b.buf), newCap)
		copy(newBuf, b.buf)
		b.buf = newBuf
	}
	b.buf = append(b.buf, data...)
}

func (b *ByteBuffer) Bytes() []byte {
	return b.buf
}
