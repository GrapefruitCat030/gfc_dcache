package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/cluster"
)

func GetClusterHandler(w http.ResponseWriter, r *http.Request) {
	m := cluster.GlobalNode().MemberList()
	b, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
