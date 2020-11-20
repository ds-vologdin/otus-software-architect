package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ds-vologdin/otus-software-architect/task05/auth/config"
	"github.com/ds-vologdin/otus-software-architect/task05/auth/providers/account"
	"github.com/ds-vologdin/otus-software-architect/task05/auth/server"
)

const (
	maxShutdownTime = 10 * time.Second
)

var (
	configFile = flag.String("config", "app.yaml", "config file")
)

func main() {
	flag.Parse()

	conf, err := config.ReadConfig(*configFile)
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	accountProvider := account.NewAccountProvider(conf.AccountService)

	address := fmt.Sprintf(":%d", conf.Server.Port)
	srv, err := server.NewServer(address, accountProvider)
	if err != nil {
		log.Fatalf("Init HTTP server: %v", err)
	}
	log.Printf("Run HTTP server: %v", address)
	srv.Run()
}
