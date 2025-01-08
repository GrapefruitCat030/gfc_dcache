// example_test.go
package protocol

import (
	"fmt"
	"testing"
)

func TestProtocol(t *testing.T) {
	// 创建SET命令
	req := &Request{
		Op:    SET,
		Key:   []byte("mykey"),
		Value: []byte("myvalue"),
	}

	// 序列化
	data, err := req.Encode()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Encoded: %s\n", data)

	// 反序列化
	decoded, err := DecodeRequest(data)
	if err != nil {
		t.Fatal(err)
	}

	// 验证
	if string(decoded.Key) != "mykey" || string(decoded.Value) != "myvalue" {
		t.Fatal("decode failed")
	}
}
