package handler

import (
	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
	"github.com/GrapefruitCat030/gfc_dcache/pkg/protocol"
)

func HandleDel(req *protocol.Request) *protocol.Response {
	err := cache.GlobalCache().Delete(string(req.Key))
	if err != nil {
		return &protocol.Response{IsError: true, Data: []byte(err.Error())}
	}
	return &protocol.Response{IsError: false, Data: []byte("del success")}
}
