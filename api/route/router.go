package route

import (
	"net/http"

	"github.com/GrapefruitCat030/gfc_dcache/api/handler"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/status", handler.GetStatusHandler).Methods(http.MethodGet)
	r.HandleFunc("/cache/{key}", handler.GetCacheHandler).Methods(http.MethodGet)
	r.HandleFunc("/cache/{key}", handler.SetCacheHandler).Methods(http.MethodPost)
	r.HandleFunc("/cache/{key}", handler.DelCacheHandler).Methods(http.MethodDelete)
	return r
}
