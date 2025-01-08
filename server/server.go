package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/cache"
)

type Server interface {
	InitServer()
	StartServer() error
	ShutdownServer() error
}

var globalServer Server

func Run(s Server, cacheType string) error {
	cache.InitCache(cacheType)
	globalServer = s
	globalServer.InitServer()
	if err := globalServer.StartServer(); err != nil {
		return err
	}
	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case sig := <-sigChan:
			log.Printf("Received signal %v, initiating shutdown...", sig)
			if err := globalServer.ShutdownServer(); err != nil {
				return err
			}
			log.Println("Server shutdown completed")
			return nil
			// TODO: hot reload
		}
	}
}
