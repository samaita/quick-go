package main

import (
	"log"

	"github.com/samaita/quick-go/config"
)

var (
	cfg *config.Config
)

func main() {
	var (
		err error
	)
	// init config
	cfg, err = config.LoadConfig("")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

}
