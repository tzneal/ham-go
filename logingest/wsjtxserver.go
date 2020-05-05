package logingest

import (
	"log"
	"net"
)

type WSJTXServer struct {
	Messages chan WSJTXMessage

	conn     *net.UDPConn
	shutdown chan struct{}
}

func NewWSJTXServer(address string) (*WSJTXServer, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err

	}
	s := &WSJTXServer{
		conn:     conn,
		shutdown: make(chan struct{}),
		Messages: make(chan WSJTXMessage),
	}
	return s, nil
}

func (s *WSJTXServer) run() {
	buf := make([]byte, 8192)
	for {
		select {
		case <-s.shutdown:
			if s.conn != nil {
				s.conn.Close()
				s.conn = nil
			}
			return
		default:
			n, _, err := s.conn.ReadFromUDP(buf)
			if err == nil {
				msg, err := WSJTXDecode(buf[0:n])
				if err != nil {
					log.Printf("error decoding WJST-X message: %s", err)
				} else if msg != nil {
					s.Messages <- msg
				}
			}

		}
	}
}

// Close gracefully shuts down the server
func (s *WSJTXServer) Close() error {
	close(s.shutdown)
	return s.conn.Close()
}

// Run is a non-blocking call that starts the server
func (s *WSJTXServer) Run() {
	go s.run()
}
