package wsjtx

import (
	"net"
)

type Server struct {
	Messages chan Message

	conn     *net.UDPConn
	shutdown chan struct{}
}

func NewServer(address string) (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err

	}
	s := &Server{
		conn:     conn,
		shutdown: make(chan struct{}),
		Messages: make(chan Message),
	}
	return s, nil
}

func (s *Server) run() {
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
				msg, err := Decode(buf[0:n])
				if err != nil {
					//log.Fatalf("%s", err)
				} else {
					s.Messages <- msg
				}
			}

		}
	}
}

// Close gracefully shuts down the server
func (s *Server) Close() error {
	close(s.shutdown)
	return s.conn.Close()
}

// Run is a non-blocking call that starts the server
func (s *Server) Run() {
	go s.run()
}
