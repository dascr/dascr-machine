package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dascr/dascr-machine/service/config"
	"github.com/dascr/dascr-machine/service/connector"
	"github.com/dascr/dascr-machine/service/ui"

	"github.com/tarm/serial"
)

func main() {
	var ip string
	var err error
	// Read eth name from ENV
	netInt := config.MustGet("INT")
	if netInt != "any" {
		ip, err = config.ReadSystemIP(netInt)
		if err != nil {
			log.Panic(err)
		}
	} else {
		ip = "0.0.0.0"
	}

	// Init config
	err = config.Init()
	if err != nil {
		panic(err)
	}

	// Setup WebServer for UI
	ui.Webs = ui.WebServer{
		IP:   ip,
		Port: 3000,
	}

	// Setup Connector service
	connector.Serv.Config = &serial.Config{
		Name:        config.Config.Machine.Serial,
		Baud:        9600,
		ReadTimeout: time.Second * 2,
	}

	go func() {
		err := ui.Webs.Start()
		if err != http.ErrServerClosed {
			log.Panicf("ERROR: Error starting web server: %+v", err)
		}
	}()

	go func() {
		err := connector.Serv.Start()
		if err != nil {
			log.Panicf("ERROR: Error starting the connector service: %+v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15)
	defer cancel()

	// Stop service
	ui.Webs.Stop(ctx)
	connector.Serv.Stop()

	// Output
	log.Println("Got ctrl+c, shutting down ...")
	os.Exit(0)
}
