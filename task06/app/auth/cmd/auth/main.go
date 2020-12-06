package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ds-vologdin/otus-software-architect/task06/app/auth/config"
	"github.com/ds-vologdin/otus-software-architect/task06/app/auth/providers/account"
	"github.com/ds-vologdin/otus-software-architect/task06/app/auth/server"
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

	accountProvider := account.NewAccountProvider(cfg.AccountService)

	address := fmt.Sprintf(":%d", cfg.Server.Port)
	srv, err := server.NewServer(cfg, accountProvider)
	if err != nil {
		log.Fatalf("Init HTTP server: %v", err)
	}
	log.Printf("Run HTTP server: %v", address)
	srv.Run()
}
