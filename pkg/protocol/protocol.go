package protocol

import (
	"bytes"
	"fmt"
	"strconv"
)

type OpType string

const (
	SET OpType = "SET"
	GET OpType = "GET"
	DEL OpType = "DEL"
)

// 请求结构体
type Request struct {
	Op    OpType
	Key   []byte
	Value []byte // 可选，用于SET操作
}

// 响应结构体
type Response struct {
	IsError bool
	Data    []byte
}

// 将字节数组序列化为 length SP content 格式
func encodeBytesArray(data []byte) []byte {
	return []byte(fmt.Sprintf("%d %s", len(data), data))
}

// 解析字节数组格式（length SP content）
func decodeBytesArray(data []byte) ([]byte, []byte, error) {
	parts := bytes.SplitN(data, []byte(" "), 2)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid format")
	}
	length, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		return nil, nil, err
	}
	content := parts[1]
	if len(content) < length {
		return nil, nil, fmt.Errorf("content length mismatch")
	}
	// 返回解析出的内容和剩余的数据
	return content[:length], content[length:], nil
}

// 序列化请求
func (req *Request) Encode() ([]byte, error) {
	if req.Op == "" {
		return nil, fmt.Errorf("missing op")
	}
	var buf bytes.Buffer
	buf.WriteString(string(req.Op))
	buf.WriteByte(' ')
	buf.Write(encodeBytesArray(req.Key))
	if req.Value != nil {
		buf.WriteByte(' ')
		buf.Write(encodeBytesArray(req.Value))
	}
	return buf.Bytes(), nil
}

// 解析请求
func DecodeRequest(data []byte) (*Request, error) {
	req := &Request{}
	for _, op := range []OpType{SET, GET, DEL} {
		if bytes.HasPrefix(data, []byte(op)) {
			req.Op = op
			data = bytes.TrimPrefix(data, []byte(op))
			data = bytes.TrimLeft(data, " ")
			break
		}
	}
	key, remaining, err := decodeBytesArray(data)
	if err != nil {
		return nil, err
	}
	req.Key = key
	if len(remaining) > 0 {
		remaining = bytes.TrimLeft(remaining, " ")
		value, _, err := decodeBytesArray(remaining)
		if err != nil {
			return nil, err
		}
		req.Value = value
	}
	return req, nil
}

// 序列化响应
func (resp *Response) Encode() []byte {
	if resp.IsError {
		return append([]byte("-"), encodeBytesArray(resp.Data)...)
	}
	return encodeBytesArray(resp.Data)
}

// 解析响应
func DecodeResponse(data []byte) (*Response, error) {
	resp := &Response{}
	if bytes.HasPrefix(data, []byte("-")) {
		resp.IsError = true
		data = bytes.TrimPrefix(data, []byte("-"))
	}
	content, _, err := decodeBytesArray(data)
	if err != nil {
		return nil, err
	}
	resp.Data = content
	return resp, nil
}

func DecodeRequestWithLeftover(data []byte) (*Request, int, error) {
	if len(data) == 0 {
		return nil, 0, fmt.Errorf("no data")
	}
	origData := data
	req := &Request{}
	// 先解析操作类型
	var matched bool
	for _, op := range []OpType{SET, GET, DEL} {
		if bytes.HasPrefix(data, []byte(op)) {
			req.Op = op
			data = bytes.TrimPrefix(data, []byte(op))
			data = bytes.TrimLeft(data, " ")
			matched = true
			break
		}
	}
	if !matched {
		return nil, 0, fmt.Errorf("invalid op")
	}
	// 解析key
	key, remaining, err := decodeBytesArray(data)
	if err != nil {
		return nil, 0, err
	}
	req.Key = key
	data = remaining
	// 如果还有数据，解析value
	data = bytes.TrimLeft(data, " ")
	if len(data) > 0 {
		value, leftover, err := decodeBytesArray(data)
		if err != nil {
			return nil, 0, err
		}
		req.Value = value
		data = leftover
	}
	used := len(origData) - len(data)
	return req, used, nil
}

func DecodeResponseWithLeftover(data []byte) (*Response, int, error) {
	origData := data
	resp := &Response{}
	if bytes.HasPrefix(data, []byte("-")) {
		resp.IsError = true
		data = bytes.TrimPrefix(data, []byte("-"))
	}
	content, leftover, err := decodeBytesArray(data)
	if err != nil {
		return nil, 0, err
	}
	resp.Data = content
	used := len(origData) - len(leftover)
	return resp, used, nil
}
