package main

import "tftpimpl/tftp"

func main() {
	//c := tftp.NewClient()
	////c.GetFile("192.168.1.111",69,"2.txt","1.txt")
	//c.PutFile("192.168.1.111",69,"2.txt","1.txt")

	s, err := tftp.NewServer(70)
	if err != nil {
		panic(err)
	}
	s.Listen()
}
