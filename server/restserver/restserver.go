package restserver

import (
	"context"
	"net/http"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
	"github.com/GrapefruitCat030/gfc_dcache/server/restserver/api/route"
)

type RESTserver struct {
	http.Server
}

func (s *RESTserver) InitServer() {
	r := route.NewRouter()
	s.Server = http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}

func (s *RESTserver) StartServer() error {
	return s.ListenAndServe()
}

func (s *RESTserver) ShutdownServer() error {
	if err := cache.GlobalCache().Close(); err != nil {
		return err
	}
	return s.Shutdown(context.TODO())
}
