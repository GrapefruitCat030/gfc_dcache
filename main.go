package main

import (
	"log"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
	"github.com/GrapefruitCat030/gfc_dcache/server"
	"github.com/GrapefruitCat030/gfc_dcache/server/selfserver"
)

func main() {
	if err := server.Run(
		&selfserver.SelfServer{},
		cache.CacheTypeLevelDB,
	); err != nil {
		log.Println(err)
	}
}
