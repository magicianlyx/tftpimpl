package main

import (
	"tftpimpl/tftp"
)

func main() {
	c := tftp.NewClient()
	c.TestRead("192.168.2.35",69,"1.txt")
	
	// s, err := tftp.NewServer(70)
	// if err != nil {
	// 	panic(err)
	// }
	// s.Listen()
}
