package route

import (
	"net/http"

	"github.com/GrapefruitCat030/gfc_dcache/server/restserver/api/handler"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/status", handler.GetStatusHandler).Methods(http.MethodGet)
	r.HandleFunc("/cache/{key}", handler.GetCacheHandler).Methods(http.MethodGet)
	r.HandleFunc("/cache/{key}", handler.SetCacheHandler).Methods(http.MethodPut)
	r.HandleFunc("/cache/{key}", handler.DelCacheHandler).Methods(http.MethodDelete)
	r.HandleFunc("/cluster", handler.GetClusterHandler).Methods(http.MethodGet)
	r.HandleFunc("/rebalance", handler.RebalanceHandler).Methods(http.MethodPost)
	return r
}
