package main

import (
	"net/http"

	"github.com/dascr/dascr-machine/service/config"
	"github.com/dascr/dascr-machine/service/connector"
	"github.com/dascr/dascr-machine/service/logger"
	"github.com/dascr/dascr-machine/service/ui"
)

func main() {
	var err error

	// Init config
	err = config.Init()
	if err != nil {
		logger.Panic(err)
	}

	// Setup WebServer for UI
	web := ui.New()

	connector.MachineConnector = connector.New()

	err = connector.MachineConnector.Start()
	if err != nil {
		logger.Warnf("Error starting the connector service: %+v", err)
	}

	err = web.Start()
	if err != http.ErrServerClosed {
		logger.Panicf("Error starting web server: %+v", err)
	}

	/*
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
	*/
}
