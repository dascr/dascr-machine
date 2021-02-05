package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/dascr/dascr-machine/service/config"
	"github.com/dascr/dascr-machine/service/connector"
	"github.com/dascr/dascr-machine/service/logger"
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
			logger.Panic(err)
		}
	} else {
		ip = "0.0.0.0"
	}

	// Init config
	err = config.Init()
	if err != nil {
		logger.Panic(err)
	}

	// Setup WebServer for UI
	ui.Webs = ui.WebServer{
		IP:   ip,
		Port: 3000,
	}

	// Setup Connector service
	connector.Serv.Config = &serial.Config{
		Name: config.Config.Machine.Serial,
		Baud: 9600,
	}

	go func() {
		err := ui.Webs.Start()
		if err != http.ErrServerClosed {
			logger.Warnf("Error starting web server: %+v", err)
		}
	}()

	go func() {
		err := connector.Serv.Start()
		if err != nil {
			logger.Warnf("Error starting the connector service: %+v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15)
	defer cancel()

	// Output
	logger.Info("Got ctrl+c, shutting down ...")

	// Stop service
	ui.Webs.Stop(ctx)
	connector.Serv.Stop()

	// Exit
	os.Exit(0)
}
