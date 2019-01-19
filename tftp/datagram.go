package tftp

import (
	"math"
	"bytes"
	"strings"
	"fmt"
)

type Datagram struct {
	data   []byte
	length int // put依赖
	index  int // get依赖
}

func newEmptyDatagram() *Datagram {
	d := &Datagram{}
	d.data = []byte{}
	return d
}

func newDatagramByBytes(b []byte) *Datagram {
	d := &Datagram{}
	d.Put(b)
	return d
}

func (d *Datagram) Put(b []byte) {
	if d == nil {
		return
	}
	addL := len(b)
	d.data = append(d.data, b...)
	d.length += addL
	return
}

func (d *Datagram) GetReset() {
	d.index = 0
}

func (d *Datagram) Get(n int) ([]byte) {
	min := func(v1, v2 int) int {
		return int(math.Min(float64(v1), float64(v2)))
	}
	dboundary := min(d.index+n, d.length)
	b := d.data[d.index:dboundary]
	d.index = dboundary
	return b
}

func (d *Datagram) GetAll() ([]byte) {
	if d.index > d.length {
		return []byte{}
	}
	b := d.data[d.index:]
	d.index = d.length
	return b
}

func (d *Datagram) SplitByByte(b byte) ([][]byte) {
	bss := bytes.Split(d.data[d.index:], []byte{b})
	d.index = d.length
	return bss
}

func (d *Datagram) ToBytes() ([]byte) {
	if d == nil {
		return nil
	}
	return d.data
}

// /-----------------------

type DatagramOp interface {
	Op() string
	Size() int
	Pack() ([]byte)
	Unpack([]byte)
}

type Options map[string]string

type RRQDatagram struct {
	OpCode   uint16
	FileName string
	Mode     string
	Options  Options
}

func NewRRQDatagram(fileName string, mode string, options Options) (*RRQDatagram) {
	if !CheckMode(mode) {
		return nil
	}
	return &RRQDatagram{opRRQ, fileName, mode, options}
}

func (d *RRQDatagram) Pack() ([]byte) {
	if d.OpCode != opRRQ {
		return nil
	}
	if !CheckMode(strings.ToLower(d.Mode)) {
		return nil
	}
	d.Mode = strings.ToLower(d.Mode)
	datagram := newEmptyDatagram()
	datagram.Put(uint16ToBytes(d.OpCode))
	datagram.Put([]byte(d.FileName))
	datagram.Put([]byte{0})
	datagram.Put([]byte(d.Mode))
	datagram.Put([]byte{0})
	if d.Options != nil {
		for k, v := range d.Options {
			datagram.Put([]byte(k))
			datagram.Put([]byte{0})
			datagram.Put([]byte(v))
			datagram.Put([]byte{0})
		}
	}
	return datagram.ToBytes()
}

func (d *RRQDatagram) Unpack(b []byte) {
	if b == nil {
		return
	}
	GetLastNotZero := func(bs []byte) int {
		for i := len(bs) - 1; i >= 0; i-- {
			if bs[i] != byte(0) {
				return i
			}
		}
		return -1
	}
	b = b[:GetLastNotZero(b)+1]
	dg := newDatagramByBytes(b)
	d.OpCode = bytesToUint16(dg.Get(2))
	if d.OpCode != opRRQ {
		return
	}
	
	bss := dg.SplitByByte(byte(0))
	if len(bss) < 2 {
		return
	}
	d.FileName = string(bss[0])
	d.Mode = string(bss[1])
	if !CheckMode(strings.ToLower(d.Mode)) {
		return
	}
	d.Options = map[string]string{}
	for i := 2; i+1 < len(bss); i += 2 {
		key := string(bss[i])
		val := string(bss[i+1])
		d.Options[key] = val
	}
	return
}

func (d *RRQDatagram) Op() string {
	return "RRQ"
}

func (d *RRQDatagram) Size() int {
	return len(d.Pack())
}

type WRQDatagram struct {
	OpCode   uint16
	FileName string
	Mode     string
	Options  Options
}

func NewWRQDatagram(fileName string, mode string, options Options) (*WRQDatagram) {
	if !CheckMode(mode) {
		return nil
	}
	return &WRQDatagram{opWRQ, fileName, mode, options}
}

func (d *WRQDatagram) Pack() ([]byte) {
	if d.OpCode != opWRQ {
		return nil
	}
	if !CheckMode(strings.ToLower(d.Mode)) {
		return nil
	}
	datagram := newEmptyDatagram()
	datagram.Put(uint16ToBytes(d.OpCode))
	datagram.Put([]byte(d.FileName))
	datagram.Put([]byte{0})
	datagram.Put([]byte(d.Mode))
	datagram.Put([]byte{0})
	if d.Options != nil {
		for k, v := range d.Options {
			datagram.Put([]byte(k))
			datagram.Put([]byte{0})
			datagram.Put([]byte(v))
			datagram.Put([]byte{0})
		}
	}
	return datagram.ToBytes()
}

func (d *WRQDatagram) Unpack(b []byte) {
	dg := newDatagramByBytes(b)
	op := bytesToUint16(dg.Get(2))
	if op != opWRQ {
		return
	}
	d.OpCode = op
	
	bss := dg.SplitByByte(byte(0))
	if len(bss) < 2 {
		return
	}
	d.FileName = string(bss[0])
	d.Mode = string(bss[1])
	if !CheckMode(strings.ToLower(d.Mode)) {
		return
	}
	d.Options = map[string]string{}
	for i := 2; i+1 < len(bss); i += 2 {
		key := string(bss[i])
		val := string(bss[i+1])
		d.Options[key] = val
	}
	return
}

func (d *WRQDatagram) Op() string {
	return "WRQ"
}

func (d *WRQDatagram) Size() int {
	return len(d.Pack())
}

type DATADatagram struct {
	OpCode  uint16
	BlockId uint16
	Data    []byte
}

func NewDATADatagram(blockId uint16, data []byte) (*DATADatagram) {
	if blockId <= 0 {
		return nil
	}
	return &DATADatagram{opDATA, blockId, data}
}
func (d *DATADatagram) Pack() ([]byte) {
	if d.OpCode != opDATA {
		return nil
	}
	if d.BlockId <= 0 {
		return nil
	}
	if len(d.Data) > DataBlockSize {
		return nil
	}
	dg := newEmptyDatagram()
	dg.Put(uint16ToBytes(d.OpCode))
	dg.Put(uint16ToBytes(d.BlockId))
	dg.Put(d.Data)
	return dg.ToBytes()
}

func (d *DATADatagram) Unpack(b []byte) {
	dg := newDatagramByBytes(b)
	d.OpCode = bytesToUint16(dg.Get(2))
	if !CheckOpCode(d.OpCode) {
		return
	}
	d.BlockId = bytesToUint16(dg.Get(2))
	d.Data = dg.GetAll()
	if len(d.Data) > DataBlockSize {
		return
	}
	return
}

func (d *DATADatagram) Op() string {
	return "DATA"
}

func (d *DATADatagram) Size() int {
	return len(d.Pack())
}

type ACKDatagram struct {
	OpCode  uint16
	BlockId uint16
}

func NewACKDatagram(blockId uint16) (*ACKDatagram) {
	if blockId < 0 {
		return nil
	}
	return &ACKDatagram{opACK, blockId}
}

func (d *ACKDatagram) Pack() ([]byte) {
	if d.OpCode != opACK {
		return nil
	}
	dg := newEmptyDatagram()
	dg.Put(uint16ToBytes(d.OpCode))
	dg.Put(uint16ToBytes(d.BlockId))
	return dg.ToBytes()
}

func (d *ACKDatagram) Unpack(b []byte) () {
	dg := newDatagramByBytes(b)
	d.OpCode = bytesToUint16(dg.Get(2))
	if d.OpCode != opACK {
		return
	}
	d.BlockId = bytesToUint16(dg.Get(2))
	return
}

func (d *ACKDatagram) Op() string {
	return "ACK"
}
func (d *ACKDatagram) Size() int {
	return len(d.Pack())
}

type ERRDatagram struct {
	OpCode  uint16
	ErrCode uint16
	ErrMsg  string
}

func (d *ERRDatagram) Error() string {
	return fmt.Sprintf("%v", d)
}

func NewERRDatagram(errCode uint16, errMsg string) (*ERRDatagram) {
	return &ERRDatagram{opERR, errCode, errMsg}
}

func (d *ERRDatagram) Pack() ([]byte) {
	if d.OpCode != opERR {
		return nil
	}
	errMsgBytes := []byte(d.ErrMsg)
	dg := newEmptyDatagram()
	dg.Put(uint16ToBytes(d.OpCode))
	dg.Put(uint16ToBytes(d.ErrCode))
	dg.Put(errMsgBytes)
	dg.Put([]byte{0})
	return dg.ToBytes()
}

func (d *ERRDatagram) Unpack(b []byte) {
	dg := newDatagramByBytes(b)
	d.OpCode = bytesToUint16(dg.Get(2))
	if d.OpCode != opERR {
		return
	}
	d.ErrCode = bytesToUint16(dg.Get(2))
	bss := dg.SplitByByte(byte(0))
	d.ErrMsg = string(bss[0])
	return
}

func (d *ERRDatagram) Op() string {
	return "ERR"
}

func (d *ERRDatagram) Size() int {
	return len(d.Pack())
}

type OACKDatagram struct {
	OpCode  uint16
	Options Options
}

func NewOACKDatagram(options Options) (*OACKDatagram) {
	return &OACKDatagram{opOACK, options}
}

func (d *OACKDatagram) Pack() ([]byte) {
	if d.OpCode != opOACK {
		return nil
	}
	dg := newEmptyDatagram()
	for k, v := range d.Options {
		dg.Put([]byte(k))
		dg.Put([]byte{0})
		dg.Put([]byte(v))
		dg.Put([]byte{0})
	}
	return dg.ToBytes()
}

func (d *OACKDatagram) Unpack(b []byte) {
	dg := newDatagramByBytes(b)
	d.OpCode = bytesToUint16(dg.Get(2))
	if d.OpCode != opOACK {
		return
	}
	d.Options = map[string]string{}
	bss := dg.SplitByByte(byte(0))
	for i := 0; i+1 < len(bss); i += 2 {
		d.Options[string(bss[i])] = string(bss[i+1])
	}
	return
}

func (d *OACKDatagram) Op() string {
	return "OACK"
}

func (d *OACKDatagram) Size() int {
	return len(d.Pack())
}

func ParseDatagram(data []byte) (DatagramOp) {
	dg := newDatagramByBytes(data)
	opCode := bytesToUint16(dg.Get(2))
	if opCode == opRRQ {
		rrq := &RRQDatagram{}
		rrq.Unpack(data)
		return rrq
	} else if opCode == opWRQ {
		wrq := &WRQDatagram{}
		wrq.Unpack(data)
		return wrq
	} else if opCode == opDATA {
		dd := &DATADatagram{}
		dd.Unpack(data)
		return dd
	} else if opCode == opACK {
		ack := &ACKDatagram{}
		ack.Unpack(data)
		return ack
	} else if opCode == opERR {
		ed := &ERRDatagram{}
		ed.Unpack(data)
		return ed
	} else if opCode == opOACK {
		od := &OACKDatagram{}
		od.Unpack(data)
		return od
	} else {
		return nil
	}
}