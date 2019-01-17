package tftp

//
// type Server struct {
// 	localHost   string
// 	localPort   int
// 	timeout     int
// 	conn        *net.UDPConn
// 	printDetail bool
// 	dataHandle  func(n int, addr *net.UDPAddr, data []byte)
// }
//
// // IPv4 ("192.0.2.1")
// // IPv6 ("2001:db8::68")
// func NewServer(port int, dataHandle func(int, *net.UDPAddr, []byte)) (*Server, error) {
// 	if port > 65535 || port <= 0 {
// 		port = 69
// 	}
// 	s := &Server{}
// 	s.localHost = "0.0.0.0"
// 	s.localPort = port
// 	s.printDetail = true
// 	s.dataHandle = dataHandle
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
// 		s.dataHandle(n, cliAddr, bs)
// 		_ = n
// 		_ = cliAddr
// 	}
// }
//
// type Client struct {
// 	remoteHost  string
// 	remotePort  int
// 	localPort   int
// 	timeout     int
// 	conn        *net.UDPConn
// 	printDetail bool
// }
//
// func NewClient() *Client {
// 	c := &Client{}
// 	addr := &net.UDPAddr{}
// 	c.localPort = addr.Port
// 	conn, err := net.ListenUDP("udp", addr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	c.conn = conn
// 	return c
// }
//
// func (c *Client) SendData(remoteHost string, remotePort int, data []byte) (error) {
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
// func (c *Client) WaitRecv() (int, *net.UDPAddr, []byte, error) {
// 	data := make([]byte, DatagramSize)
// 	n, addr, err := c.conn.ReadFromUDP(data)
// 	if err != nil {
// 		return 0, nil, nil, err
// 	}
// 	return n, addr, data, nil
// }
//
//
// func (c *Client) WriteFile(filename string, mode string) (error) {
// 	// send RRQ
// 	rrq, err := NewRRQDatagram(filename, mode, nil)
// 	if err != nil {
// 		return HandleError(err)
// 	}
// 	data, err := rrq.Pack()
// 	if err != nil {
// 		return HandleError(err)
// 	}
// 	c.conn.Write(data)
//
// }
