package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/GrapefruitCat030/gfc_dcache/server"
	"github.com/GrapefruitCat030/gfc_dcache/server/restserver"
	"github.com/GrapefruitCat030/gfc_dcache/server/selfserver"
)

type config struct {
	serverType  string
	cacheType   string
	clusterAddr string
}

var cfg config

var rootCmd = &cobra.Command{
	Use:   "gfc_dcache",
	Short: "GFC DCache is a distributed cache system",
	Run: func(cmd *cobra.Command, args []string) {
		var srv server.Server
		switch cfg.serverType {
		case "self":
			srv = &selfserver.SelfServer{}
		case "http":
			srv = &restserver.RESTserver{}
		default:
			log.Fatalf("Unknown server type: %s", cfg.serverType)
		}
		if err := server.Run(srv, cfg.cacheType); err != nil {
			log.Println(err)
		}
	},
}

func main() {
	rootCmd.Flags().StringVarP(&cfg.serverType, "server", "s", "http", "Type of server to run (e.g., http, self)")
	rootCmd.Flags().StringVarP(&cfg.cacheType, "cache", "c", "memory", "Type of cache to use (e.g., memory, leveldb)")
	/*
		Due to the characteristics of the Gossip protocol, any node in the cluster will gradually
		spread the new node information to the entire cluster after receiving it. Therefore, it does not
		matter which node the cluster parameter selects, only needs to be a node that already exists in the cluster.
	*/
	rootCmd.Flags().StringVarP(&cfg.clusterAddr, "cluster", "l", "", "random cluster address(cause gossip protocol, it)")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
