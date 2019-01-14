package tftpimpl

import (
	"math"
	"bytes"
	"strings"
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
	Pack() ([]byte, error)
	Unpack([]byte) (error)
}

type Options map[string]string

type RRQDatagram struct {
	OpCode   uint16
	FileName string
	Mode     string
	Options  Options
}

func NewRRQDatagram(fileName string, mode string, options Options) (*RRQDatagram, error) {
	if !CheckMode(mode) {
		return nil, ErrParam
	}
	return &RRQDatagram{opRRQ, fileName, mode, options}, nil
}

func (d *RRQDatagram) Pack() ([]byte, error) {
	if d.OpCode != opRRQ {
		return nil, ErrParam
	}
	if !CheckMode(strings.ToLower(d.Mode)) {
		return nil, ErrParam
	}
	d.Mode = strings.ToLower(d.Mode)
	datagram := newEmptyDatagram()
	datagram.Put(Uint16ToBytes(d.OpCode))
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
	return datagram.ToBytes(), nil
}

func (d *RRQDatagram) Unpack(b []byte) (error) {
	dg := newDatagramByBytes(b)
	d.OpCode = BytesToUint16(dg.Get(2))
	if d.OpCode != opRRQ {
		return ErrDatagram
	}
	
	bss := dg.SplitByByte(byte(0))
	if len(bss) < 2 {
		return ErrDatagram
	}
	d.FileName = string(bss[0])
	d.Mode = string(bss[1])
	if !CheckMode(strings.ToLower(d.Mode)) {
		return ErrDatagram
	}
	d.Options = map[string]string{}
	for i := 2; i+1 < len(bss); i += 2 {
		key := string(bss[i])
		val := string(bss[i+1])
		d.Options[key] = val
	}
	return nil
}

type WRQDatagram struct {
	OpCode   uint16
	FileName string
	Mode     string
	Options  Options
}

func NewWRQDatagram(fileName string, mode string, options Options) (*WRQDatagram, error) {
	if !CheckMode(mode) {
		return nil, ErrParam
	}
	return &WRQDatagram{opWRQ, fileName, mode, options}, nil
}

func (d *WRQDatagram) Pack() ([]byte, error) {
	if d.OpCode != opWRQ {
		return nil, ErrParam
	}
	if !CheckMode(strings.ToLower(d.Mode)) {
		return nil, ErrParam
	}
	datagram := newEmptyDatagram()
	datagram.Put(Uint16ToBytes(d.OpCode))
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
	return datagram.ToBytes(), nil
}

func (d *WRQDatagram) Unpack(b []byte) (error) {
	dg := newDatagramByBytes(b)
	op := BytesToUint16(dg.Get(2))
	if op != opWRQ {
		return ErrDatagram
	}
	d.OpCode = op
	
	bss := dg.SplitByByte(byte(0))
	if len(bss) < 2 {
		return ErrDatagram
	}
	d.FileName = string(bss[0])
	d.Mode = string(bss[1])
	if !CheckMode(strings.ToLower(d.Mode)) {
		return ErrDatagram
	}
	d.Options = map[string]string{}
	for i := 2; i+1 < len(bss); i += 2 {
		key := string(bss[i])
		val := string(bss[i+1])
		d.Options[key] = val
	}
	return nil
}

type DATADatagram struct {
	OpCode  uint16
	BlockId uint16
	Data    []byte
}

func NewDATADatagram(blockId uint16, data []byte) (*DATADatagram, error) {
	if blockId <= 0 {
		return nil, ErrParam
	}
	return &DATADatagram{opDATA, blockId, data}, nil
}
func (d *DATADatagram) Pack() ([]byte, error) {
	if d.OpCode != opDATA {
		return nil, ErrParam
	}
	if d.BlockId <= 0 {
		return nil, ErrParam
	}
	if len(d.Data) > DataBlockSize {
		return nil, ErrDataTooLong
	}
	dg := newEmptyDatagram()
	dg.Put(Uint16ToBytes(d.OpCode))
	dg.Put(Uint16ToBytes(d.BlockId))
	dg.Put(d.Data)
	return dg.ToBytes(), nil
}

func (d *DATADatagram) Unpack(b []byte) (error) {
	dg := newDatagramByBytes(b)
	d.OpCode = BytesToUint16(dg.Get(2))
	if !CheckOpCode(d.OpCode) {
		return ErrDatagram
	}
	d.BlockId = BytesToUint16(dg.Get(2))
	d.Data = dg.GetAll()
	if len(d.Data) > DataBlockSize {
		return ErrDataTooLong
	}
	return nil
}

type ACKDatagram struct {
	OpCode  uint16
	BlockId uint16
}

func NewACKDatagram(blockId uint16) (*ACKDatagram, error) {
	if blockId <= 0 {
		return nil, ErrParam
	}
	return &ACKDatagram{opACK, blockId}, nil
}

func (d *ACKDatagram) Pack() ([]byte, error) {
	if d.OpCode != opACK {
		return nil, ErrParam
	}
	dg := newEmptyDatagram()
	dg.Put(Uint16ToBytes(d.OpCode))
	dg.Put(Uint16ToBytes(d.BlockId))
	return dg.ToBytes(), nil
}

func (d *ACKDatagram) Unpack(b []byte) (error) {
	dg := newDatagramByBytes(b)
	d.OpCode = BytesToUint16(dg.Get(2))
	if d.OpCode != opACK {
		return ErrDatagram
	}
	d.BlockId = BytesToUint16(dg.Get(2))
	return nil
}

type ERRDatagram struct {
	OpCode  uint16
	ErrCode uint16
	ErrMsg  string
}

func NewERRDatagram(errCode uint16, errMsg string) (*ERRDatagram, error) {
	return &ERRDatagram{opERR, errCode, errMsg}, nil
}

func (d *ERRDatagram) Pack() ([]byte, error) {
	if d.OpCode != opERR {
		return nil, ErrParam
	}
	errMsgBytes := []byte(d.ErrMsg)
	dg := newEmptyDatagram()
	dg.Put(Uint16ToBytes(d.OpCode))
	dg.Put(Uint16ToBytes(d.ErrCode))
	dg.Put(errMsgBytes)
	dg.Put([]byte{0})
	return dg.ToBytes(), nil
}

func (d *ERRDatagram) Unpack(b []byte) (error) {
	dg := newDatagramByBytes(b)
	d.OpCode = BytesToUint16(dg.Get(2))
	if d.OpCode != opERR {
		return ErrDatagram
	}
	d.ErrCode = BytesToUint16(dg.Get(2))
	bss := dg.SplitByByte(byte(0))
	d.ErrMsg = string(bss[0])
	return nil
}

type OACKDatagram struct {
	OpCode  uint16
	Options Options
}

func NewOACKDatagram(options Options) (*OACKDatagram, error) {
	return &OACKDatagram{opOACK, options}, nil
}

func (d *OACKDatagram) Pack() ([]byte, error) {
	if d.OpCode != opOACK {
		return nil, ErrParam
	}
	dg := newEmptyDatagram()
	for k, v := range d.Options {
		dg.Put([]byte(k))
		dg.Put([]byte{0})
		dg.Put([]byte(v))
		dg.Put([]byte{0})
	}
	return dg.ToBytes(), nil
}

func (d *OACKDatagram) Unpack(b []byte) (error) {
	dg := newDatagramByBytes(b)
	d.OpCode = BytesToUint16(dg.Get(2))
	if d.OpCode != opOACK {
		return ErrDatagram
	}
	d.Options = map[string]string{}
	bss := dg.SplitByByte(byte(0))
	for i := 0; i+1 < len(bss); i += 2 {
		d.Options[string(bss[i])] = string(bss[i+1])
	}
	return nil
}

func ParseDatagram(data []byte) (interface{}, error) {
	dg := newDatagramByBytes(data)
	opCode := BytesToUint16(dg.Get(2))
	if opCode == opRRQ {
		rrq := &RRQDatagram{}
		err := rrq.Unpack(data)
		return rrq, err
	} else if opCode == opWRQ {
		wrq := &WRQDatagram{}
		err := wrq.Unpack(data)
		return wrq, err
	} else if opCode == opDATA {
		dd := &DATADatagram{}
		err := dd.Unpack(data)
		return dd, err
	} else if opCode == opACK {
		ack := &ACKDatagram{}
		err := ack.Unpack(data)
		return ack, err
	} else if opCode == opERR {
		ed := &ERRDatagram{}
		err := ed.Unpack(data)
		return ed, err
	} else if opCode == opOACK {
		od := &OACKDatagram{}
		err := od.Unpack(data)
		return od, err
	} else {
		return nil, ErrDatagram
	}
	
}
