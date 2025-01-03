package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
	"github.com/gorilla/mux"
)

func SetCacheHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(val) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Set the value to the cache
	err = cache.GlobalCache().Set(key, val)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
