package handler

import (
	"bytes"
	"log"
	"net/http"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
	"github.com/GrapefruitCat030/gfc_dcache/pkg/cluster"
)

func RebalanceHandler(w http.ResponseWriter, r *http.Request) {
	sc := cache.GlobalCache().NewScanner()
	defer sc.Close()
	cli := &http.Client{}
	for sc.Scan() {
		key := sc.Key()
		nodeAddr, ok := cluster.GlobalNode().ShouldProcess(key)
		if !ok {
			req, err := http.NewRequest(http.MethodPut, "http://"+nodeAddr+"/cache/"+key, bytes.NewReader(sc.Value()))
			if err != nil {
				log.Printf("Rebalance: Failed to create request. %v", err)
				continue
			}
			resp, err := cli.Do(req)
			if err != nil {
				log.Printf("Rebalance: Failed to send request to %v. %v", nodeAddr, err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("Rebalance: Failed to set cache on %v. %v", nodeAddr, resp.StatusCode)
				continue
			}
			cache.GlobalCache().Delete(key)
		}
	}
}
