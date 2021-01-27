package ui

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/dascr/dascr-machine/service/config"
	"github.com/dascr/dascr-machine/service/connector"
	"github.com/gorilla/mux"
)

// Webs is the global webserver
var Webs WebServer

// WebServer will hold the information to run the UIs webserver
type WebServer struct {
	IP         string
	Port       int
	HTTPServer *http.Server
}

// Start will start the ui webserver
func (ws *WebServer) Start() error {
	// Setup routing
	mux := mux.NewRouter()
	mux.HandleFunc("/admin", ws.admin)
	mux.HandleFunc("/updateMachine", ws.updateMachine)
	mux.HandleFunc("/updateScoreboard", ws.updateScoreboard)
	mux.HandleFunc("/debugSerial", ws.debugSerial)
	add := fmt.Sprintf("%+v:%+v", ws.IP, ws.Port)
	ws.HTTPServer = &http.Server{
		Addr:    add,
		Handler: mux,
	}
	log.Printf("Navigate to http://%+v:%+v/admin to configure your darts machine", ws.IP, ws.Port)

	if err := ws.HTTPServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// Stop will stop the ui webserver
func (ws *WebServer) Stop(ctx context.Context) {
	ws.HTTPServer.Shutdown(ctx)
	log.Println("Webserver stopped")
}

// Admin route
func (ws *WebServer) admin(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/admin.html")
	if err != nil {
		log.Printf("Error when reading template admin.html %+v", err)
		http.Error(w, "Error when reading template admin.html", http.StatusBadRequest)
		return
	}

	err = t.Execute(w, &config.Config)
	if err != nil {
		log.Printf("Cannot execute template admin.html: %+v", err)
		http.Error(w, "Error when executing template admin.html", http.StatusBadRequest)
		return
	}
}

// Update routes
func (ws *WebServer) updateMachine(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error when updating machine settings", http.StatusBadRequest)
		return
	}

	delayTime, err := strconv.Atoi(r.FormValue("delay"))
	if err != nil {
		log.Printf("Invalid form input: %+v", err)
		http.Error(w, "Error when updating machine settings", http.StatusBadRequest)
		return
	}
	if config.Config.Machine.WaitingTime != delayTime {
		config.Config.Machine.WaitingTime = delayTime
	}
	if config.Config.Machine.Serial != r.FormValue("serial") {
		config.Config.Machine.Serial = r.FormValue("serial")
	}

	err = config.SaveConfig()
	if err != nil {
		log.Printf("Error writing config file after update: %+v", err)
		http.Error(w, "Error when updating machine settings", http.StatusBadRequest)
		return
	}

	// restart connector service
	connector.Serv.Restart()

	// Redirect back to admin
	http.Redirect(w, r, "admin", http.StatusSeeOther)
}

func (ws *WebServer) updateScoreboard(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error when updating scoreborad settings", http.StatusBadRequest)
		return
	}

	https := r.FormValue("sbprot") != ""
	if config.Config.Scoreboard.HTTPS != https {
		config.Config.Scoreboard.HTTPS = https
	}

	if config.Config.Scoreboard.Host != r.FormValue("sbhost") {
		config.Config.Scoreboard.Host = r.FormValue("sbhost")
	}

	if config.Config.Scoreboard.Port != r.FormValue("sbport") {
		config.Config.Scoreboard.Port = r.FormValue("sbport")
	}

	if config.Config.Scoreboard.Game != r.FormValue("sbgame") {
		config.Config.Scoreboard.Game = r.FormValue("sbgame")
	}

	err = config.SaveConfig()
	if err != nil {
		log.Printf("Error writing config file after update: %+v", err)
		http.Error(w, "Error when updating scoreborad settings", http.StatusBadRequest)
		return
	}

	// restart connector service
	connector.Serv.Restart()

	// Redirect back to admin
	http.Redirect(w, r, "admin", http.StatusSeeOther)
}

func (ws *WebServer) debugSerial(w http.ResponseWriter, r *http.Request) {
	input, ok := r.URL.Query()["debug"]
	if !ok || len(input[0]) < 1 {
		log.Println("No debug string provided, use debug='command'")
		return
	}

	if input[0] == "read" {
		output := connector.Serv.Read()
		log.Printf("output in debug function: %+v", output)
	} else {
		connector.Serv.Write(input[0])
	}
}
