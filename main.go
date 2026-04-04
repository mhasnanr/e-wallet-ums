package main

import (
	"log"

	"github.com/mhasnanr/ewallet-ums/bootstrap"
	"github.com/mhasnanr/ewallet-ums/cmd"
)

func main() {
	if err := bootstrap.SetupZapLogger(); err != nil {
		log.Fatal("failed to initialize logger")
	}

	if err := bootstrap.SetupConfig(".env"); err != nil {
		log.Fatalf("failed to load config file")
	}

	bootstrap.SetupDatabase()

	go cmd.ServeGRPC()
	cmd.ServeHTTP(bootstrap.DB)
}
