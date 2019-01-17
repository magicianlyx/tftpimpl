package tftp

import (
	"net"
	"fmt"
	"errors"
)

var (
	ErrNotFromSource = errors.New("data not from source")
)

type Client struct {
	localPort   int
	timeout     int
	conn        *net.UDPConn
	printDetail bool
	retry       int // udp传输出错重试次数
}

func NewClient() *Client {
	c := &Client{}
	addr := &net.UDPAddr{}
	c.localPort = addr.Port
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	c.conn = conn
	return c
}

func (c *Client) sendData(remoteHost string, remotePort int, data []byte) (error) {
	addr := fmt.Sprintf("%s:%d", remoteHost, remotePort)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		err = HandleError(err)
		return err
	}
	_, err = c.conn.WriteToUDP(data, udpAddr)
	if err != nil {
		err = HandleError(err)
		return err
	}
	return nil
}

func (c *Client) waitRecv() (int, *net.UDPAddr, []byte, error) {
	data := make([]byte, DatagramSize)
	n, addr, err := c.conn.ReadFromUDP(data)
	if err != nil {
		return 0, nil, nil, err
	}
	return n, addr, data, nil
}

func (c *Client) SendTest(remoteHost string, remotePort int, fileName string) (error) {
	wrq, err := NewWRQDatagram(fileName, modeOCTET, nil)
	if err != nil {
		return err
	}
	rrqData, err := wrq.Pack()
	if err != nil {
		return err
	}
	err = c.sendData(remoteHost, remotePort, rrqData)
	if err != nil {
		return err
	}
	_, addr, data, err := c.waitRecv()
	i, err := ParseDatagram(data)
	switch v := i.(type) {
	case *ACKDatagram:
		fmt.Println(v.OpCode)
		fmt.Println(v.BlockId)
		break
	default:
		fmt.Println(v)
	}
	_ = addr
	
	return nil
}

func RetryFunc(f func() error, time int) error {
	var err error
	for i := 0; i < time; time++ {
		err = f()
		if err == nil {
			return nil
		}
	}
	return err
}
