package server

import (
	"net/http"

	"github.com/GrapefruitCat030/gfc_dcache/api/route"
	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
)

type Server struct {
	http.Server
}

var globalServer *Server

func InitServer(cacheName string) {
	cache.InitCache(cacheName)
	r := route.NewRouter()
	globalServer = &Server{
		Server: http.Server{
			Addr:    ":8080",
			Handler: r,
		},
	}
}

func StartServer() error {
	return globalServer.ListenAndServe()
}
