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

type Article struct {
	Host string
	Port int
	Size int
	Data []byte
}

type Client struct {
	localAddr     *net.UDPAddr
	timeout       int
	conn          *net.UDPConn
	printDetail   bool
	retry         int // udp传输出错重试次数
}

func NewClient() *Client {
	c := &Client{}
	addr := &net.UDPAddr{}
	c.localAddr = addr
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	c.conn = conn
	c.printDetail = true
	c.retry = 3
	return c
}

func (c *Client) GetFile(remoteHost string, remotePort int, fileName string, localFileName string) (error) {
	fileData := []byte{}
	err := c.sendData(remoteHost, remotePort, NewRRQDatagram(fileName, modeOctet, nil))
	if err != nil {
		return err
	}
	for {
		n, addr, data, err := c.waitRecv()
		if err != nil {
			
			return err
		}
		i := ParseDatagram(data)
		switch v := i.(type) {
		case *DATADatagram:
			dataSegment := v.Data[:n]
			fileData = append(fileData, dataSegment...)
			c.sendData(addr.IP.String(), addr.Port, NewACKDatagram(v.BlockId))
			if n < DataBlockSize {
				return saveFileData(localFileName, fileData)
			}
			break
		case *OACKDatagram:
			c.sendData(addr.IP.String(), addr.Port, NewACKDatagram(0))
			break
		case *ERRDatagram:
			return v
			break
		default:
			break
		}
	}
}

func (c *Client) PutFile(remoteHost string, remotePort int, fileName string, LocalFileName string) (error) {
	fs, err := os.Open(LocalFileName)
	if err != nil {
		return err
	}
	fileData, err := ioutil.ReadAll(fs)
	if err != nil {
		return err
	}
	bss, blockCount := splitDataSegment(fileData, DataBlockSize)
	err = c.sendData(remoteHost, remotePort, NewWRQDatagram(fileName, modeOctet, nil))
	if err != nil {
		return err
	}
	for {
		_, addr, data, err := c.waitRecv()
		if err != nil {
			
			return err
		}
		i := ParseDatagram(data)
		switch v := i.(type) {
		case *ACKDatagram:
			var segmentData = []byte{}
			if uint16(blockCount) > v.BlockId {
				segmentData = bss[v.BlockId]
			}
			c.sendData(addr.IP.String(), addr.Port, NewDATADatagram(v.BlockId+1, segmentData))
		case *OACKDatagram:
			var segmentData = []byte{}
			if uint16(blockCount) > 0 {
				segmentData = bss[0]
			}
			c.sendData(addr.IP.String(), addr.Port, NewDATADatagram(1, segmentData))
		case *ERRDatagram:
			return v
		default:
		
		}
	}
}

func (c *Client) waitRecv() (int, *net.UDPAddr, []byte, error) {
	data := make([]byte, DatagramSize)
	n, addr, err := c.conn.ReadFromUDP(data)
	data = data[:n]
	i := ParseDatagram(data)
	if c.printDetail {
		PrintDynamic(ActionRecv, addr, i)
	}
	return n, addr, data, err
}

// 向远程服务器发送数据包（发送失败会重试，retry次）
func (c *Client) sendData(remoteHost string, remotePort int, op DatagramOp) (error) {
	addr := fmt.Sprintf("%s:%d", remoteHost, remotePort)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	for i := 0; i < c.retry; i++ {
		_, err = c.conn.WriteToUDP(op.Pack(), udpAddr)
		if err != nil {
		
		} else {
			err = nil
			if c.printDetail {
				PrintDynamic(ActionSend, udpAddr, op)
			}
			break
		}
	}
	return err
}
