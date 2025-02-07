package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
	"github.com/GrapefruitCat030/gfc_dcache/pkg/cluster"
	"github.com/GrapefruitCat030/gfc_dcache/server/restserver"
	"github.com/GrapefruitCat030/gfc_dcache/server/selfserver"
)

type Server interface {
	InitServer()
	StartServer() error
	ShutdownServer() error
}

func Run(cacheTtl int, cacheType, nodeAddr, clusterAddr string) error {
	cache.InitCache(cacheType, cacheTtl)
	cluster.InitNode(nodeAddr, clusterAddr)
	// start TCP + HTTP server
	servers := []Server{
		&selfserver.SelfServer{},
		&restserver.RESTserver{},
	}
	for _, s := range servers {
		s.InitServer()
		go func(s Server) {
			if err := s.StartServer(); err != nil {
				log.Println(err)
			}
		}(s)
	}
	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case sig := <-sigChan:
			log.Printf("Received signal %v, initiating shutdown...", sig)
			for _, s := range servers {
				if err := s.ShutdownServer(); err != nil {
					log.Println(err)
				}
			}
			log.Println("Server shutdown completed")
			if err := cache.ShotdownCache(); err != nil {
				log.Println(err)
			}
			log.Printf("Cache shutdown completed")
			return nil
			// TODO: hot reload
		}
	}
}
