package tftp

import (
	"testing"
	"fmt"
	"github.com/json-iterator/go"
)

// 测试使用
func ToJson(v interface{}) string {
	bs, err := jsoniter.MarshalIndent(v, "", "   ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(bs)
}

func TestRRQDatagram(t *testing.T) {
	rrq := &RRQDatagram{
		opRRQ,
		"filename",
		modeOCTET,
		nil,
	}
	bs, err := rrq.Pack()
	if err != nil {
		t.Fatalf("Pack: %v", err)
	}
	rrqup := &RRQDatagram{}
	err = rrqup.Unpack(bs)
	if err != nil {
		t.Fatalf("Unpack:%v", err)
	} else {
		fmt.Println(ToJson(rrqup))
	}

}

func TestWRQDatagram(t *testing.T) {
	wrq := &WRQDatagram{
		opWRQ,
		"filename",
		modeOCTET,
		nil,
	}
	bs, err := wrq.Pack()
	if err != nil {
		t.Fatalf("Pack: %v", err)
	}
	wrqup := &WRQDatagram{}
	err = wrqup.Unpack(bs)
	if err != nil {
		t.Fatalf("Unpack:%v", err)
	} else {
		fmt.Println(ToJson(wrqup))
	}
}

func TestDATADatagram(t *testing.T) {
	ddg := &DATADatagram{
		opDATA,
		1,
		[]byte("asdf"),
	}
	bs, err := ddg.Pack()
	if err != nil {
		t.Fatalf("Pack: %v", err)
	}
	ddgup := &DATADatagram{}
	err = ddgup.Unpack(bs)
	if err != nil {
		t.Fatalf("Unpack:%v", err)
	} else {
		fmt.Printf("%s\r\n", ddgup.Data)
	}
}

func TestACKDatagram(t *testing.T) {
	ddg := &ACKDatagram{
		opACK,
		1,
	}
	bs, err := ddg.Pack()
	if err != nil {
		t.Fatalf("Pack: %v", err)
	}
	ddgup := &ACKDatagram{}
	err = ddgup.Unpack(bs)
	if err != nil {
		t.Fatalf("Unpack:%v", err)
	} else {
		fmt.Printf("%s\r\n", ToJson(ddgup))
	}
}

func TestERRDatagram(t *testing.T) {
	ddg := &ERRDatagram{
		opERR,
		tftpErrFileNotFound,
		"error msg",
	}
	bs, err := ddg.Pack()
	if err != nil {
		t.Fatalf("Pack: %v", err)
	}
	ddgup := &ERRDatagram{}
	err = ddgup.Unpack(bs)
	if err != nil {
		t.Fatalf("Unpack:%v", err)
	} else {
		fmt.Printf("%s\r\n", ToJson(ddgup))
	}
}
