package selfserver

import (
	"context"
	"log"
	"net"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/protocol"
	"github.com/GrapefruitCat030/gfc_dcache/server/selfserver/api/handler"
	"github.com/GrapefruitCat030/gfc_dcache/server/selfserver/api/route"
)

type SelfServer struct {
	listener net.Listener
	addr     string
	router   *route.Router
	ctx      context.Context
	cancel   context.CancelFunc
}

func (s *SelfServer) InitServer() {
	ctx, cancel := context.WithCancel(context.Background())
	s.addr = ":8080"
	s.ctx = ctx
	s.cancel = cancel

	s.router = route.NewRouter()
	s.router.Register(protocol.GET, handler.HandleGet)
	s.router.Register(protocol.SET, handler.HandleSet)
	s.router.Register(protocol.DEL, handler.HandleDel)
}

func (s *SelfServer) StartServer() error {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				log.Printf("accept error: %v\n", err)
				continue
			}
			go s.handleConnection(conn)
		}
	}
}

func (s *SelfServer) ShutdownServer() error {
	s.cancel()
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *SelfServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		num, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("connection closed\n")
				return
			}
			log.Printf("connection read error: %v\n", err)
			continue
		}
		req, err := protocol.DecodeRequest(buf[:num])
		if err != nil {
			resp := &protocol.Response{IsError: true, Data: []byte(err.Error())}
			conn.Write(resp.Encode())
			continue
		}
		resp := s.handleRequest(req)
		conn.Write(resp.Encode())
	}
}

func (s *SelfServer) handleRequest(req *protocol.Request) *protocol.Response {
	return s.router.Dispatch(req)
}
