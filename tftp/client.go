package tftp

import (
	"net"
	"fmt"
	"errors"
	"os"
	"io/ioutil"
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
	wrq := NewWRQDatagram(fileName, modeOCTET, map[string]string{"tsize": "0"})
	rrqData := wrq.Pack()
	rrqData = BytesFill(rrqData, DatagramSize)
	fmt.Println(rrqData)
	err := c.sendData(remoteHost, remotePort, rrqData)
	if err != nil {
		return err
	}
	fs, err := os.Open(fileName)
	fileData, err := ioutil.ReadAll(fs)
	datas, n := SplitDataSegment(fileData, DataBlockSize)
	for {
		_, addr, data, err := c.waitRecv()
		i := ParseDatagram(data)
		switch v := i.(type) {
		case *ACKDatagram:
			fmt.Println(v.BlockId)
			if n >= int(v.BlockId+1) {
				dg := NewDATADatagram(v.BlockId+1, datas[int(v.BlockId)])
				bs := dg.Pack()
				c.sendData(remoteHost, remotePort, bs)
			}
			break
		default:
			fmt.Println(v)
			fmt.Println("---------")
		}
		_ = addr
		_ = err
	}
	
	return nil
}

func (c *Client) TestRead(remoteHost string, remotePort int, fileName string) (error) {
	rrq := NewRRQDatagram(fileName, modeOCTET, map[string]string{"tsize": "4"})
	rrqData := rrq.Pack()
	rrqData = BytesFill(rrqData, DatagramSize)
	err := c.sendData(remoteHost, remotePort, rrqData)
	if err != nil {
		return err
	}
	for {
		_, addr, data, err := c.waitRecv()
		fmt.Println("address:", addr.IP.String(), addr.Port)
		i := ParseDatagram(data)
		switch v := i.(type) {
		case *ACKDatagram:
			fmt.Println("ACKDatagram: ", v)
			break
		case *DATADatagram:
			fmt.Println("DATADatagram: ", v)
			fmt.Println(BytesFill(NewACKDatagram(v.BlockId).Pack(), DatagramSize))
			c.sendData(addr.IP.String(), addr.Port, BytesFill(NewACKDatagram(v.BlockId).Pack(), DatagramSize))
			break
		case *OACKDatagram:
			fmt.Println("OACKDatagram: ", v)
			c.sendData(addr.IP.String(), addr.Port, NewACKDatagram(0).Pack())
			break
		default:
			// fmt.Println("default: ", v)
		}
		_ = addr
		_ = err
	}
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
