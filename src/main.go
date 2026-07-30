package main

import (
	"flag"
	"log"

	"github.com/cpd007/myredis/config"
	"github.com/cpd007/myredis/server"
)

func setUpFlags() {
	flag.StringVar(&config.Config.Host, "host", "0.0.0.0", "Host for redis server")
	flag.IntVar(&config.Config.Port, "port", 7379, "Port for redis server")
	flag.Parse()
}

func main() {
	setUpFlags()
	log.Println("starting the server")
	server.Start()
}
