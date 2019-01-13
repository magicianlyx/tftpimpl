package tftpimpl

import (
	"testing"
	"time"
	"fmt"
)

func TestUDP(t *testing.T) {
	s,err := NewServer(69, func(bytes []byte) {
		fmt.Println(fmt.Sprintf("%s\r\n", bytes))
	})
	if err != nil {
		t.Fatalf("%s", err)
		return
	}
	go s.Listen()

	time.Sleep(time.Second*3)
	c := NewClient()
	err = c.Send("192.168.1.117", 69, []byte("asdf"))
	if err != nil {
		t.Fatalf("%s", err)
		return
	}
	time.Sleep(time.Hour)
}
