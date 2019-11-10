package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pdk/forum/conf"
	"github.com/pdk/forum/srv"
	"github.com/pdk/forum/store"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("usage: forum config.json")
		os.Exit(1)
	}
	configFileName := os.Args[1]

	log.Printf("starting forum...")

	config, err := conf.ReadConfiguration(configFileName)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("configuration loaded: %v", config)

	db, err := store.NewConnection(config.Database)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %s", err)
	}

	server, err := srv.NewServer(db, config.AssetsDir)
	if err != nil {
		log.Fatalf("failed to initialize server: %s", err)
	}

	server.ListenAndServe(config.ListenAddress)
}
