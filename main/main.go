package main

import "fmt"

func main() {
	is := make([]byte, 2)
	copy(is[:1], []byte{byte(9)})
	copy(is[1:2], []byte{byte(8)})
	fmt.Println(is)
}
