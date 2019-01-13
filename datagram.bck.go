package tftpimpl

//type Datagram struct {
//	data   []byte
//	size   int
//	length int // put依赖
//	index  int // get依赖
//}
//
//func newEmptyDatagram() *Datagram {
//	d := &Datagram{}
//	d.data = make([]byte, DatagramSize)
//	d.size = DatagramSize
//	d.length = 0
//	return d
//}
//
//func newDatagramByBytes(b []byte) *Datagram {
//	d := &Datagram{}
//	d.Put(b)
//	return d
//}
//func newDatagramWithSize(size int) *Datagram {
//	d := &Datagram{}
//	d.size = size
//	d.data = make([]byte, size)
//	return d
//}
//
//func (d *Datagram) Put(b []byte) {
//	min := func(v1, v2 int) int {
//		return int(math.Min(float64(v1), float64(v2)))
//	}
//	if d == nil {
//		return
//	}
//	bboundary := min(d.size, d.length+len(b))
//	dboundary := min(len(b), d.size-d.length)
//	n := copy(d.data[d.length:dboundary], b[:bboundary])
//	d.length += n
//	return
//}
//
//func (d *Datagram) Get(n int) ([]byte) {
//	min := func(v1, v2 int) int {
//		return int(math.Min(float64(v1), float64(v2)))
//	}
//	dboundary := min(d.index+n, d.size)
//	b := d.data[d.index:dboundary]
//	d.index = dboundary
//	return b
//}
//
//func (d *Datagram) GetAll() ([]byte) {
//	b := d.data[d.index:d.size]
//	d.index = d.size
//	return b
//}
//
//func (d *Datagram) SplitByByte(b byte) ([][]byte) {
//	return bytes.Split(d.data[d.index:], []byte{b})
//}
//
//func (d *Datagram) ToBytes() ([]byte) {
//	if d == nil {
//		return nil
//	}
//	return d.data
//}
