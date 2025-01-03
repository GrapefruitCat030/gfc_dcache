package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
)

func GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	stat := cache.GlobalCache().GetStatus()
	v, err := json.Marshal(stat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(v)
}
