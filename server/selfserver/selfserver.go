package selfserver

import (
	"context"
	"errors"
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
	s.addr = ":8081"
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
				if errors.Is(err, net.ErrClosed) {
					log.Printf("tcp: listener closed")
					return nil
				}
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
	asyncProcessor := newConnProcessor(conn)
	defer asyncProcessor.close()
	readBuf := make([]byte, 1024*10)
	var dataBuf []byte
	for {
		num, err := conn.Read(readBuf)
		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("connection closed\n")
				return
			}
			log.Printf("connection read error: %v\n", err)
			continue
		}
		dataBuf = append(dataBuf, readBuf[:num]...)
		for {
			req, used, err := protocol.DecodeRequestWithLeftover(dataBuf)
			if err != nil {
				// 无法完整解析时退出循环，等待更多数据
				break
			}
			// 成功解析请求后，将已用字节剔除
			dataBuf = dataBuf[used:]
			s.handleRequest(asyncProcessor, req)
		}
	}
}

func (s *SelfServer) handleRequest(ap *asyncProcessor, req *protocol.Request) {
	ch := make(chan *protocol.Response)
	ap.respChan <- ch
	go func() {
		ch <- s.router.Dispatch(req)
	}()
}
