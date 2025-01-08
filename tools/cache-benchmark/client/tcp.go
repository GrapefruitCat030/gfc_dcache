package client

import (
	"bufio"
	"errors"
	"net"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/protocol"
)

type tcpClient struct {
	net.Conn
	r *bufio.Reader
}

func (c *tcpClient) sendRequest(req *protocol.Request) error {
	data, err := req.Encode()
	if err != nil {
		return err
	}
	_, err = c.Write(data)
	return err
}

func (c *tcpClient) recvResponse() (*protocol.Response, error) {
	buf := make([]byte, 1024)
	n, err := c.r.Read(buf)
	if err != nil {
		return nil, err
	}
	return protocol.DecodeResponse(buf[:n])
}

func (c *tcpClient) sendGet(key string) error {
	req := &protocol.Request{
		Op:  protocol.GET,
		Key: []byte(key),
	}
	return c.sendRequest(req)
}

func (c *tcpClient) sendSet(key, value string) error {
	req := &protocol.Request{
		Op:    protocol.SET,
		Key:   []byte(key),
		Value: []byte(value),
	}
	return c.sendRequest(req)
}

func (c *tcpClient) sendDel(key string) error {
	req := &protocol.Request{
		Op:  protocol.DEL,
		Key: []byte(key),
	}
	return c.sendRequest(req)
}

func (c *tcpClient) Run(cmd *Cmd) {
	var err error
	switch cmd.Name {
	case "get":
		err = c.sendGet(cmd.Key)
	case "set":
		err = c.sendSet(cmd.Key, cmd.Value)
	case "del":
		err = c.sendDel(cmd.Key)
	default:
		panic("unknown cmd name " + cmd.Name)
	}

	if err != nil {
		cmd.Error = err
		return
	}

	resp, err := c.recvResponse()
	if err != nil {
		cmd.Error = err
		return
	}

	if resp.IsError {
		cmd.Error = errors.New(string(resp.Data))
	} else {
		cmd.Value = string(resp.Data)
	}
}

func (c *tcpClient) PipelinedRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}
	for _, cmd := range cmds {
		var err error
		switch cmd.Name {
		case "get":
			err = c.sendGet(cmd.Key)
		case "set":
			err = c.sendSet(cmd.Key, cmd.Value)
		case "del":
			err = c.sendDel(cmd.Key)
		default:
			panic("unknown cmd name " + cmd.Name)
		}
		if err != nil {
			cmd.Error = err
			return
		}
	}
	for _, cmd := range cmds {
		resp, err := c.recvResponse()
		if err != nil {
			cmd.Error = err
			return
		}
		if resp.IsError {
			cmd.Error = errors.New(string(resp.Data))
		} else {
			cmd.Value = string(resp.Data)
		}
	}
}

func newTCPClient(server string) *tcpClient {
	c, err := net.Dial("tcp", server)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(c)
	return &tcpClient{c, r}
}
