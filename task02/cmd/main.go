package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ds-vologdin/otus-software-architect/task02/config"
	"github.com/ds-vologdin/otus-software-architect/task02/server"
	"github.com/ds-vologdin/otus-software-architect/task02/users/service"
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

	userService, err := service.NewUserService(conf.Database.DSN)
	if err != nil {
		log.Fatalf("Init UserService: %v", err)
	}

	address := fmt.Sprintf(":%d", conf.Server.Port)
	srv, err := server.NewServer(address, userService)
	if err != nil {
		log.Fatalf("Init HTTP server: %v", err)
	}
	log.Printf("Run HTTP server: %v", address)
	srv.Run()
}
