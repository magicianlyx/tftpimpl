package tftpimpl

func Uint16ToBytes(v uint16) (b []byte) {
	b = make([]byte, 2)
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	b[0] = byte(v >> 8)
	b[1] = byte(v)
	return b
}


func BytesToUint16(b []byte)uint16{
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[1]) | uint16(b[0])<<8
}

