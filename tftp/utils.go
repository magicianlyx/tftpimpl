package tftp

func Uint16ToBytes(v uint16) (b []byte) {
	b = make([]byte, 2)
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	b[0] = byte(v >> 8)
	b[1] = byte(v)
	return b
}

func BytesToUint16(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[1]) | uint16(b[0])<<8
}

func BytesFill(bs []byte, size int) ([]byte) {
	if size <= len(bs) {
		return bs
	}
	nbs := make([]byte, size)
	copy(nbs[:len(bs)], bs)
	return nbs
}
