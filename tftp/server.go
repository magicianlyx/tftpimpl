package tftp

import (
	"net"
	"fmt"
	"time"
	"os"
	"io"
)

type Server struct {
	localHost   string
	localPort   int
	timeout     time.Duration
	conn        *net.UDPConn
	printDetail bool
	msgChan     chan DatagramOp
	retry       int
	connInfos   *ConnInfos
}

// IPv4 ("192.0.2.1")
// IPv6 ("2001:db8::68")
func NewServer(port int) (*Server, error) {
	if port > 65535 || port <= 0 {
		port = 69
	}
	s := &Server{}
	s.localHost = "0.0.0.0"
	s.localPort = port
	s.printDetail = true
	s.retry = 3
	s.connInfos = NewConnInfos()
	s.msgChan = make(chan DatagramOp, 100)
	addr := fmt.Sprintf("%s:%d", s.localHost, s.localPort)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	s.conn = conn
	return s, nil
}

func (s *Server) sendData(remoteHost string, remotePort int, op DatagramOp) (error) {
	addr := fmt.Sprintf("%s:%d", remoteHost, remotePort)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	for i := 0; i < s.retry; i++ {
		_, err = s.conn.WriteToUDP(op.Pack(), udpAddr)
		if err == nil {
			if s.printDetail {
				PrintDynamic(ActionSend, udpAddr, op)
			}
			break
		}
	}
	return err
}
func (s *Server) Listen() {
	for {
		var bs = make([]byte, DatagramSize)
		n, cliAddr, err := s.conn.ReadFromUDP(bs)
		if err != nil {
			continue
		}
		bs = bs[:n]
		do := ParseDatagram(bs)
		if s.printDetail {
			PrintDynamic(ActionRecv, cliAddr, do)
		}
		s.DatagramAnalysis(cliAddr.IP.String(), cliAddr.Port, do)
	}
}

func (s *Server) DatagramAnalysis(addrIp string, addrPort int, do DatagramOp) {
	switch v := do.(type) {
	case *RRQDatagram:
		if fs, err := os.OpenFile(v.FileName, 0666, os.ModePerm); err == nil {
			if v.Mode == modeOctet {
				fileData := make([]byte, DataBlockSize)
				n, err := fs.Read(fileData)
				if err == nil {
					fileData = fileData[:n]
				}
				s.sendData(addrIp, addrPort, NewDATADatagram(1, fileData))
			}
			s.connInfos.Set(EncodeAddr(addrIp, addrPort), NewConnInfo(fs, v.Options, v.Mode, DataBlockSize))
		} else {
			s.sendData(addrIp, addrPort, NewERRDatagram(tftpErrFileNotFound, "file not found"))
		}
		break
	case *WRQDatagram:
		if fs, err := os.Create(v.FileName); err != nil {
			s.sendData(addrIp, addrPort, NewERRDatagram(tftpErrFileAlreadyExists, "file already exists"))
		} else {
			if v.Options == nil {
				s.sendData(addrIp, addrPort, NewOACKDatagram(v.Options))
			} else {
				s.sendData(addrIp, addrPort, NewACKDatagram(0))
			}
			s.connInfos.Set(EncodeAddr(addrIp, addrPort), NewConnInfo(fs, v.Options, v.Mode, DataBlockSize))
		}
		break
	case *DATADatagram:
		if fc := s.connInfos.Get(EncodeAddr(addrIp, addrPort)); fc == nil {
			s.sendData(addrIp, addrPort, NewERRDatagram(tftpErrIllegalOperation, "illegal operation"))
		} else {
			fc.FileInfo.Write(v.Data)
			if len(v.Data) < fc.DataBlockSize {
				fc.FileInfo.Close()
				s.connInfos.Del(EncodeAddr(addrIp, addrPort))
			}
			s.sendData(addrIp, addrPort, NewACKDatagram(fc.BlockId))
			fc.BlockId += 1
			s.connInfos.Set(EncodeAddr(addrIp, addrPort), fc)
		}
		break
	case *ACKDatagram:
		if fc := s.connInfos.Get(EncodeAddr(addrIp, addrPort)); fc == nil {
			s.sendData(addrIp, addrPort, NewERRDatagram(tftpErrIllegalOperation, "illegal operation"))
		} else {
			fileData := make([]byte, DataBlockSize)
			fc.FileInfo.Seek(int64(fc.DataBlockSize)*int64(fc.BlockId), io.SeekStart)
			n, err := fc.FileInfo.Read(fileData)
			if err == nil {
				fileData = fileData[:n]
			}
			if n<DataBlockSize{
				s.connInfos.Del(EncodeAddr(addrIp, addrPort))
			}
			s.sendData(addrIp, addrPort, NewDATADatagram(v.BlockId+1, fileData))
			fc.BlockId = v.BlockId + 1
			s.connInfos.Set(EncodeAddr(addrIp, addrPort), fc)
		}
		break
	default:
		break
	}
}

func EncodeAddr(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}
