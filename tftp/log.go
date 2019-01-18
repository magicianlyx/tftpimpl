package tftp

import (
	"net"
	"log"
	"strconv"
)

const (
	ActionSend = "Send"
	ActionRecv = "Receive"
)

func PrintDynamic(action string, addr *net.UDPAddr, op DatagramOp) {
	if op == nil {
		return
	}
	if op.Pack() == nil {
		return
	}
	if len(op.Pack()) < op.Size() {
		return
	}
	log.Printf("   |action=>%-10s |addr=>%-20s |op=>%-10s|data=>%v|",
		action,
		addr.IP.String()+":"+strconv.Itoa(addr.Port),
		op.Op(),
		op.Pack()[:op.Size()],
	)
}
