package protocol

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	conn    net.Conn
	addr    string
	timeout time.Duration
}

// 创建新客户端
func NewClient(addr string, timeout time.Duration) *Client {
	return &Client{
		addr:    addr,
		timeout: timeout,
	}
}

// 连接服务器
func (c *Client) Connect() error {
	if c.conn != nil {
		return nil
	}

	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return fmt.Errorf("connect error: %v", err)
	}
	c.conn = conn
	return nil
}

// 关闭连接
func (c *Client) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

// 发送请求并获取响应
func (c *Client) do(req *Request) (*Response, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}

	// 设置读写超时
	c.conn.SetDeadline(time.Now().Add(c.timeout))

	// 编码并发送请求
	data, err := req.Encode()
	if err != nil {
		return nil, err
	}

	if _, err := c.conn.Write(data); err != nil {
		c.Close()
		return nil, fmt.Errorf("write error: %v", err)
	}

	// 读取响应
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("read error: %v", err)
	}

	// 解码响应
	resp, err := DecodeResponse(buf[:n])
	if err != nil {
		return nil, fmt.Errorf("decode response error: %v", err)
	}

	return resp, nil
}

// Set 操作
func (c *Client) Set(key, value []byte) error {
	req := &Request{
		Op:    SET,
		Key:   key,
		Value: value,
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	if resp.IsError {
		return fmt.Errorf("server error: %s", resp.Data)
	}
	return nil
}

// Get 操作
func (c *Client) Get(key []byte) ([]byte, error) {
	req := &Request{
		Op:  GET,
		Key: key,
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if resp.IsError {
		return nil, fmt.Errorf("server error: %s", resp.Data)
	}
	return resp.Data, nil
}

// Del 操作
func (c *Client) Del(key []byte) error {
	req := &Request{
		Op:  DEL,
		Key: key,
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	if resp.IsError {
		return fmt.Errorf("server error: %s", resp.Data)
	}
	return nil
}
