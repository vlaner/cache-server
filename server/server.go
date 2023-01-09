package server

import (
	"log"
	"net"

	"github.com/vlaner/cache-server/cache"
	"github.com/vlaner/cache-server/commands"
)

type Server struct {
	listerAddr string
	ln         net.Listener
	cache      cache.Cache
	stop       chan struct{}
}

func New(listerAddr string) *Server {
	return &Server{
		listerAddr: listerAddr,
		ln:         nil,
		cache:      cache.New(),
		stop:       make(chan struct{}, 1),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listerAddr)
	if err != nil {
		return err
	}

	s.ln = ln

	go s.serve()

	return nil
}

func (s *Server) Stop() {
	s.stop <- struct{}{}
	close(s.stop)
	s.ln.Close()
}

func (s *Server) serve() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.stop:
				return
			default:
				log.Printf("ERROR ACCEPTING CONNECTION: %v", err)
				continue
			}
		}

		go s.handleConnection(conn)

	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("ERROR READING FROM CONNECTION: %v", err)
			break
		}
		payload, err := commands.Parse(string(buf[:n]))
		if err != nil {
			_, err = conn.Write([]byte(err.Error() + "\n"))
			if err != nil {
				log.Printf("ERROR WRITING TO CONNECTION: %v", err)
				break
			}
			continue
		}
		var val string
		switch payload.Cmd {
		case commands.COMMAND_GET:
			val, err = s.cache.Get(payload.Key)
		case commands.COMMAND_SET:
			val, err = s.cache.Set(payload.Key, payload.Value)
		case commands.COMMAND_DEL:
			val, err = s.cache.Del(payload.Key)
		case commands.COMMAND_EXPIRE:
			val, err = s.cache.Expire(payload.Key, payload.Expire)
		}
		if err != nil {
			_, err = conn.Write([]byte(err.Error() + "\n"))
			if err != nil {
				log.Printf("ERROR WRITING TO CONNECTION: %v", err)
				break
			}
			continue
		}

		_, err = conn.Write([]byte(val + "\n"))
		if err != nil {
			log.Printf("ERROR WRITING TO CONNECTION: %v", err)
			break
		}
	}

}
