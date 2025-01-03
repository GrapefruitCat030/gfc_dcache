package main

import (
	"log"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/server"
)

func main() {
	server.InitServer("memory")
	if err := server.StartServer(); err != nil {
		log.Println(err)
	}
}
