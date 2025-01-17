package selfserver

import (
	"fmt"
	"net"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/protocol"
)

// asyncProcessor 异步处理器
type asyncProcessor struct {
	respChan chan chan *protocol.Response
	conn     net.Conn
}

// process 负责reply消息
func (ap *asyncProcessor) process() {
	defer ap.conn.Close() // asyncProcessor退出时关闭连接
	for {                 // 从respChan中读取响应, 始终保持与请求相同的顺序
		ch, ok := <-ap.respChan
		if !ok {
			return
		}
		resp := <-ch
		if _, err := ap.conn.Write(resp.Encode()); err != nil {
			fmt.Printf("ap conn write error: %v\n", err)
			return
		}
	}
}

func (ap *asyncProcessor) close() {
	close(ap.respChan)
}

func newConnProcessor(conn net.Conn) *asyncProcessor {
	ap := &asyncProcessor{
		respChan: make(chan chan *protocol.Response),
		conn:     conn,
	}
	go ap.process()
	return ap
}
