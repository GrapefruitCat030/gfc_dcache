package client

import (
	"bufio"
	"errors"
	"net"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/protocol"
)

type TcpClient struct {
	net.Conn
	r *bufio.Reader
}

func (c *TcpClient) sendRequest(req *protocol.Request) error {
	data, err := req.Encode()
	if err != nil {
		return err
	}
	_, err = c.Write(data)
	return err
}

func (c *TcpClient) recvResponse() (*protocol.Response, error) {
	buf := make([]byte, 1024)
	n, err := c.r.Read(buf)
	if err != nil {
		return nil, err
	}
	return protocol.DecodeResponse(buf[:n])
}

func (c *TcpClient) sendGet(key string) error {
	req := &protocol.Request{
		Op:  protocol.GET,
		Key: []byte(key),
	}
	return c.sendRequest(req)
}

func (c *TcpClient) sendSet(key, value string) error {
	req := &protocol.Request{
		Op:    protocol.SET,
		Key:   []byte(key),
		Value: []byte(value),
	}
	return c.sendRequest(req)
}

func (c *TcpClient) sendDel(key string) error {
	req := &protocol.Request{
		Op:  protocol.DEL,
		Key: []byte(key),
	}
	return c.sendRequest(req)
}

func (c *TcpClient) Do(op *Operation) error {
	var err error
	switch op.Name {
	case "get":
		err = c.sendGet(op.Key)
	case "set":
		err = c.sendSet(op.Key, op.Value)
	case "del":
		err = c.sendDel(op.Key)
	default:
		panic("unknown op name " + op.Name)
	}
	if err != nil {
		return err
	}
	resp, err := c.recvResponse()
	if err != nil {
		return err
	}
	if resp.IsError {
		return errors.New(string(resp.Data))
	}
	op.Value = string(resp.Data)
	return nil
}

func (c *TcpClient) PipelinedDo(ops []*Operation) error {
	if len(ops) == 0 {
		return nil
	}
	for _, op := range ops {
		var err error
		switch op.Name {
		case "get":
			err = c.sendGet(op.Key)
		case "set":
			err = c.sendSet(op.Key, op.Value)
		case "del":
			err = c.sendDel(op.Key)
		default:
			panic("unknown op name " + op.Name)
		}
		if err != nil {
			return err
		}
	}
	for _, op := range ops {
		resp, err := c.recvResponse()
		if err != nil {
			return err
		}
		if resp.IsError {
			return errors.New(string(resp.Data))
		}
		op.Value = string(resp.Data)
	}
	return nil
}

func newTCPClient(serverAddr string) *TcpClient {
	c, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(c)
	return &TcpClient{c, r}
}
