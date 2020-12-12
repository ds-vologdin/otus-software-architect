package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ds-vologdin/otus-software-architect/task06/app/billing/bill"

	"github.com/ds-vologdin/otus-software-architect/task06/app/billing/config"
	"github.com/ds-vologdin/otus-software-architect/task06/app/billing/server"
)

const (
	maxShutdownTime = 10 * time.Second
)

var (
	configFile = flag.String("config", "app.yaml", "config file")
)

func main() {
	flag.Parse()

	cfg, err := config.ReadConfig(*configFile)
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	billService, err := bill.NewBillService(cfg.Database)
	if err != nil {
		log.Fatalf("create a bill service: %v", err)
	}

	address := fmt.Sprintf(":%d", cfg.Server.Port)
	srv, err := server.NewServer(cfg, billService)
	if err != nil {
		log.Fatalf("Init HTTP server: %v", err)
	}
	log.Printf("Run HTTP server: %v", address)
	srv.Run()
}
