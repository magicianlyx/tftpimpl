package main

import (
	"tftpimpl/tftp"
)

func main() {
	c := tftp.NewClient()
	c.GetFile("192.168.2.35",69,"1.txt","3.txt")
	
	// s, err := tftp.NewServer(70)
	// if err != nil {
	// 	panic(err)
	// }
	// s.Listen()
}
