package main

import (
	"log"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
	"github.com/GrapefruitCat030/gfc_dcache/server"
	"github.com/GrapefruitCat030/gfc_dcache/server/restserver"
)

func main() {
	if err := server.Run(
		&restserver.RESTserver{},
		cache.CacheTypeMemory,
	); err != nil {
		log.Println(err)
	}
}
