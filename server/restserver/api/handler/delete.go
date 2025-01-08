package handler

import (
	"log"
	"net/http"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
	"github.com/gorilla/mux"
)

func DelCacheHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	if err := cache.GlobalCache().Delete(key); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
