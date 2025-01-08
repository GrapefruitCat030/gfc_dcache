package handler

import (
	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
	"github.com/GrapefruitCat030/gfc_dcache/pkg/protocol"
)

func HandleGet(req *protocol.Request) *protocol.Response {
	val, err := cache.GlobalCache().Get(string(req.Key))
	if err != nil {
		return &protocol.Response{IsError: true, Data: []byte("key not found")}
	}
	return &protocol.Response{IsError: false, Data: val}
}
