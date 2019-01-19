package tftp

import (
	"net"
	"fmt"
	"errors"
	"os"
	"io/ioutil"
	"time"
)

var (
	ErrTimeOut       = errors.New("wait server time out")
	ErrIllegalTftpOp = errors.New("illegal tftp operation")
)

type Article struct {
	Host string
	Port int
	Size int
	Data []byte
}

type Client struct {
	timeout     time.Duration // 从发送到接受回复的最长时间间隔
	conn        *net.UDPConn  // 连接tftp服务器的udp回话连接
	printDetail bool          // 是否打印细节
	retry       int           // udp传输出错重试次数
}

func NewClient() *Client {
	c := &Client{}
	addr := &net.UDPAddr{}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	c.timeout = time.Millisecond * 500
	c.conn = conn
	c.printDetail = true
	c.retry = 3
	return c
}

func (c *Client) GetFile(remoteHost string, remotePort int, fileName string, localFileName string) (error) {
	fileData := []byte{}
	n, addr, do, err := c.sendAndRecv(remoteHost, remotePort, NewRRQDatagram(fileName, modeOctet, nil))
	if err != nil {
		return err
	}
	for {
		switch v := do.(type) {
		case *DATADatagram:
			dataSegment := v.Data
			fileData = append(fileData, dataSegment...)
			if n < DataBlockSize {
				// 最后一段数据
				c.sendData(addr.IP.String(), addr.Port, NewACKDatagram(v.BlockId))
				return saveFileData(localFileName, fileData)
			} else {
				// 不是最后一段数据
				n, addr, do, err = c.sendAndRecv(addr.IP.String(), addr.Port, NewACKDatagram(v.BlockId))
			}
			break
		case *OACKDatagram:
			n, addr, do, err = c.sendAndRecv(addr.IP.String(), addr.Port, NewACKDatagram(0))
			break
		case *ERRDatagram:
			return v
		default:
			return ErrIllegalTftpOp
		}
	}
}

func (c *Client) PutFile(remoteHost string, remotePort int, fileName string, LocalFileName string) (error) {
	fs, err := os.Open(LocalFileName)
	defer fs.Close()
	if err != nil {
		return err
	}
	fileData, err := ioutil.ReadAll(fs)
	if err != nil {
		return err
	}
	bss, blockCount := splitDataSegment(fileData, DataBlockSize)
	// 最后一段数据长度为数据块大小时 后面添加一个空的字节数组
	if len(bss[len(bss)-1]) == DataBlockSize {
		bss = append(bss, []byte{})
	}
	_, addr, do, err := c.sendAndRecv(remoteHost, remotePort, NewWRQDatagram(fileName, modeOctet, nil))
	if err != nil {
		return err
	}
	for {
		switch v := do.(type) {
		case *ACKDatagram:
			var segmentData = []byte{}
			if v.BlockId >= uint16(blockCount) {
				// 发送完毕
				return nil
			} else {
				segmentData = bss[v.BlockId]
			}
			_, addr, do, err = c.sendAndRecv(addr.IP.String(), addr.Port, NewDATADatagram(v.BlockId+1, segmentData))
		case *OACKDatagram:
			// 将会发送第一个数据报文
			var segmentData = []byte{}
			if uint16(blockCount) > 0 {
				segmentData = bss[0]
			}
			_, addr, do, err = c.sendAndRecv(addr.IP.String(), addr.Port, NewDATADatagram(1, segmentData))
		case *ERRDatagram:
			return v
		default:
			return ErrDatagram
		}
	}
}

// 向服务器发送tftp报文并等待响应报文
func (c *Client) sendAndRecv(remoteHost string, remotePort int, op DatagramOp) (int, *net.UDPAddr, DatagramOp, error) {
	var n int
	var addr *net.UDPAddr
	var data DatagramOp
	var err error
	var recv = false // 是否接收成功 recv为false时有可能已经接收到响应 recv为true时一定接收到响应
	err = c.sendData(remoteHost, remotePort, op)
	if err != nil {
		return n, nil, nil, err
	}
	waitRecvTimeout(
		func() {
			n, addr, data, err = c.waitRecv()
			recv = true
		},
		c.timeout,
	)
	if recv {
		return n, addr, data, nil
	} else {
		return n, nil, nil, ErrTimeOut
	}
}

// 阻塞 等待服务器发送数据包
func (c *Client) waitRecv() (int, *net.UDPAddr, DatagramOp, error) {
	data := make([]byte, DatagramSize)
	n, addr, err := c.conn.ReadFromUDP(data)
	data = data[:n]
	i := ParseDatagram(data)
	if c.printDetail {
		PrintDynamic(ActionRecv, addr, i)
	}
	return n, addr, i, err
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
		if err == nil {
			if c.printDetail {
				PrintDynamic(ActionSend, udpAddr, op)
			}
			break
		}
	}
	return err
}
