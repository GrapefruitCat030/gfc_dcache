package main

import (
	"fmt"
	"log"
	"time"

	"github.com/GrapefruitCat030/gfc_dcache/pkg/protocol"
)

func main() {
	client := protocol.NewClient("localhost:8080", 5*time.Second)
	defer client.Close()

	// SET 操作
	err := client.Set([]byte("name"), []byte("gopher"))
	if err != nil {
		log.Fatal(err)
	}

	// GET 操作
	value, err := client.Get([]byte("name"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Got value: %s\n", value)

	// DEL 操作
	err = client.Del([]byte("name"))
	if err != nil {
		log.Fatal(err)
	}
}
