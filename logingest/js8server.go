package logingest

import (
	"log"
	"net"
)

type JS8Server struct {
	Messages chan JS8Message

	conn     *net.UDPConn
	shutdown chan struct{}
}

func NewJS8Server(address string) (*JS8Server, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err

	}
	s := &JS8Server{
		conn:     conn,
		shutdown: make(chan struct{}),
		Messages: make(chan JS8Message),
	}
	return s, nil
}

func (s *JS8Server) run() {
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
				msg, err := JS8Decode(buf[0:n])
				if err != nil {
					log.Printf("erorr decoding JS8Call QSO: %s", err)
				} else {
					s.Messages <- msg
				}
			}

		}
	}
}

// Close gracefully shuts down the server
func (s *JS8Server) Close() error {
	close(s.shutdown)
	return s.conn.Close()
}

// Run is a non-blocking call that starts the server
func (s *JS8Server) Run() {
	go s.run()
}
