package main

import (
	"tftpimpl/tftp"
)

func main() {
	c := tftp.NewClient()
	c.SendTest("192.168.2.35",69,"2.txt")
	// c.SendFile("192.168.2.35",69,"2.txt")
}
