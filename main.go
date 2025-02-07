package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/GrapefruitCat030/gfc_dcache/server"
)

type config struct {
	cacheType   string
	cacheTtl    int
	nodeAddr    string
	clusterAddr string
}

var cfg config

var rootCmd = &cobra.Command{
	Use:   "gfc_dcache",
	Short: "GFC DCache is a distributed cache system",
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.Run(cfg.cacheTtl, cfg.cacheType, cfg.nodeAddr, cfg.clusterAddr); err != nil {
			log.Println(err)
		}
	},
}

func main() {
	rootCmd.Flags().StringVarP(&cfg.cacheType, "cache", "c", "memory", "Type of cache to use (e.g., memory, leveldb)")
	rootCmd.Flags().IntVarP(&cfg.cacheTtl, "ttl", "t", 30, "Time to live for cache items in seconds")
	/*
		Due to the characteristics of the Gossip protocol, any node in the cluster will gradually
		spread the new node information to the entire cluster after receiving it. Therefore, it does not
		matter which node the cluster parameter selects, only needs to be a node that already exists in the cluster.
	*/
	rootCmd.Flags().StringVarP(&cfg.nodeAddr, "node", "n", "127.0.0.1", "node address")
	rootCmd.Flags().StringVarP(&cfg.clusterAddr, "cluster", "l", "", "random cluster address(cause gossip protocol, it)")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
