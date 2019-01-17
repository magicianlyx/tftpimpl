package tftp

import (
	"testing"
	"time"
	"fmt"
	"net"
)

func TestUDP(t *testing.T) {
	s, err := NewServer(69, func(n int, addr *net.UDPAddr, data []byte) {
		fmt.Println(fmt.Sprintf("%s\r\n", data))
	})
	if err != nil {
		t.Fatalf("%s", err)
		return
	}
	go s.Listen()
	
	time.Sleep(time.Second * 3)
	c := NewClient()
	err = c.sendData("192.168.1.117", 69, []byte("asdf"))
	if err != nil {
		t.Fatalf("%s", err)
		return
	}
	time.Sleep(time.Hour)
}
