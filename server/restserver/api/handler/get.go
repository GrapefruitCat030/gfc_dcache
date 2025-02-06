package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
	"github.com/GrapefruitCat030/gfc_dcache/pkg/cluster"

	"github.com/gorilla/mux"
)

func GetCacheHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val, err := cache.GlobalCache().Get(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(val) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(val)
}

func GetClusterHandler(w http.ResponseWriter, r *http.Request) {
	m := cluster.GlobalNode().MemberList()
	b, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
