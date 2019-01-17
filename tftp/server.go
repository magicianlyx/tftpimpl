package tftp

// type Server struct {
// 	localHost   string
// 	localPort   int
// 	timeout     int
// 	conn        *net.UDPConn
// 	printDetail bool
// 	// dataHandle  func(n int, addr *net.UDPAddr, data []byte)
// 	msgChan map[string]chan []byte
// 	retry   int
// }
// 
// func (c *Server) sendData(remoteHost string, remotePort int, data []byte) (error) {
// 	addr := fmt.Sprintf("%s:%d", remoteHost, remotePort)
// 	udpAddr, err := net.ResolveUDPAddr("udp", addr)
// 	if err != nil {
// 		err = HandleError(err)
// 		return err
// 	}
// 	_, err = c.conn.WriteToUDP(data, udpAddr)
// 	if err != nil {
// 		err = HandleError(err)
// 		return err
// 	}
// 	return nil
// }
// 
// func (c *Server) waitRecv() (int, *net.UDPAddr, []byte, error) {
// 	data := make([]byte, DatagramSize)
// 	n, addr, err := c.conn.ReadFromUDP(data)
// 	if err != nil {
// 		return 0, nil, nil, err
// 	}
// 	return n, addr, data, nil
// }
// 
// // IPv4 ("192.0.2.1")
// // IPv6 ("2001:db8::68")
// func NewServer(port int) (*Server, error) {
// 	if port > 65535 || port <= 0 {
// 		port = 69
// 	}
// 	s := &Server{}
// 	s.localHost = "0.0.0.0"
// 	s.localPort = port
// 	s.printDetail = true
// 	s.retry = 3
// 	// s.dataHandle = dataHandle
// 	s.msgChan = map[string]chan []byte{}
// 	addr := fmt.Sprintf("%s:%d", s.localHost, s.localPort)
// 	udpAddr, err := net.ResolveUDPAddr("udp", addr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	conn, err := net.ListenUDP("udp", udpAddr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	s.conn = conn
// 	return s, nil
// }
// 
// func (s *Server) Listen() {
// 	for {
// 		var bs = make([]byte, DatagramSize)
// 		n, cliAddr, err := s.conn.ReadFromUDP(bs)
// 		if err != nil {
// 			err = HandleError(err)
// 			continue
// 		}
// 		err = s.dataHandle(n, cliAddr, bs)
// 		err = HandleError(err)
// 		_ = n
// 		_ = cliAddr
// 	}
// }
// 
// func (s *Server) dataHandle(n int, addr *net.UDPAddr, data []byte) error {
// 	i, err := ParseDatagram(data)
// 	if err != nil {
// 		return err
// 	}
// 	// var blockId uint16 = 0
// 	switch v := i.(type) {
// 	case *RRQDatagram:
// 		if ok := FileIsExist(v.FileName); ok {
// 			// ack, err := NewACKDatagram(0)
// 			// if err != nil {
// 			// 	HandleError(err)
// 			// 	break
// 			// } else {
// 			// 	pushData, err := ack.Pack()
// 			// 	if err != nil {
// 			// 		HandleError(err)
// 			// 		break
// 			// 	}
// 			// 	err = RetryFunc(func() error {
// 			// 		return s.sendData(addr.IP.String(), addr.Port, pushData)
// 			// 	}, s.retry)
// 			// 	if err != nil {
// 			// 		HandleError(err)
// 			// 		break
// 			// 	}
// 			// }
// 		} else {
// 			errData, err := NewERRDatagram(tftpErrFileNotFound, "file not found")
// 			if err != nil {
// 				HandleError(err)
// 				break
// 			} else {
// 				pushData, err := errData.Pack()
// 				if err != nil {
// 					HandleError(err)
// 					break
// 				}
// 				err = RetryFunc(func() error {
// 					return s.sendData(addr.IP.String(), addr.Port, pushData)
// 				}, s.retry)
// 				if err != nil {
// 					HandleError(err)
// 					break
// 				}
// 			}
// 		}
// 		break
// 	case *WRQDatagram:
// 		if ok := FileIsExist(v.FileName); ok {
// 			errData, err := NewERRDatagram(tftpErrFileAlreadyExists, "file already exists")
// 			if err != nil {
// 				HandleError(err)
// 				break
// 			} else {
// 				pushData, err := errData.Pack()
// 				if err != nil {
// 					HandleError(err)
// 					break
// 				}
// 				err = RetryFunc(func() error {
// 					return s.sendData(addr.IP.String(), addr.Port, pushData)
// 				}, s.retry)
// 				if err != nil {
// 					HandleError(err)
// 					break
// 				}
// 			}
// 		} else {
// 			ack, err := NewACKDatagram()
// 			if err != nil{
// 				HandleError(err)
// 				break
// 			}else{
// 			
// 			}
// 		}
// 		break
// 	case *ACKDatagram:
// 		break
// 	case *DATADatagram:
// 		break
// 	default:
// 		return ErrDatagram
// 	}
// }
// 
// func FileIsExist(fileName string) bool {
// 	fs, err := os.Open(fileName)
// 	defer fs.Close()
// 	if err != nil && os.IsNotExist(err) {
// 		return false
// 	}
// 	return true
// }
