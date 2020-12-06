package main

import (
	"flag"
	"log"
	"time"

	"github.com/ds-vologdin/otus-software-architect/task06/app/auth-sidecar/config"
	"github.com/ds-vologdin/otus-software-architect/task06/app/auth-sidecar/server"
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

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Init HTTP server: %v", err)
	}
	srv.Run()
}
