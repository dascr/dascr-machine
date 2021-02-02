package ui

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/dascr/dascr-machine/service/config"
	"github.com/dascr/dascr-machine/service/connector"
	"github.com/dascr/dascr-machine/service/logger"
)

// Webs is the global webserver
var Webs WebServer

// Static will provide the embedded files as http.FS
//go:embed static
var webui embed.FS

// WebServer will hold the information to run the UIs webserver
type WebServer struct {
	IP         string
	Port       int
	HTTPServer *http.Server
}

// Start will start the ui webserver
func (ws *WebServer) Start() error {
	staticdir, err := fs.Sub(webui, "static")
	if err != nil {
		return err
	}
	// Setup routing
	http.Handle("/", http.FileServer(http.FS(staticdir)))
	http.HandleFunc("/admin", ws.admin)
	http.HandleFunc("/updateMachine", ws.updateMachine)
	http.HandleFunc("/updateScoreboard", ws.updateScoreboard)
	add := fmt.Sprintf("%+v:%+v", ws.IP, ws.Port)
	ws.HTTPServer = &http.Server{
		Addr: add,
	}
	logger.Infof("Navigate to http://%+v:%+v to configure your darts machine", ws.IP, ws.Port)

	if err := ws.HTTPServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// Stop will stop the ui webserver
func (ws *WebServer) Stop(ctx context.Context) {
	ws.HTTPServer.Shutdown(ctx)
	logger.Info("Webserver stopped")
}

// Admin route
func (ws *WebServer) admin(w http.ResponseWriter, r *http.Request) {
	file, err := webui.ReadFile("static/admin.html")
	if err != nil {
		logger.Errorf("Error when reading template admin.html %+v", err)
		http.Error(w, "Error when reading template admin.html", http.StatusBadRequest)
		return
	}

	t := template.New("adminPage")
	_, err = t.Parse(string(file))
	if err != nil {
		logger.Errorf("Error when reading template admin.html %+v", err)
		http.Error(w, "Error when reading template admin.html", http.StatusBadRequest)
		return
	}

	err = t.Execute(w, &config.Config)
	if err != nil {
		logger.Errorf("Cannot execute template admin.html: %+v", err)
		http.Error(w, "Error when executing template admin.html", http.StatusBadRequest)
		return
	}
}

// Update routes
func (ws *WebServer) updateMachine(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logger.Error(err)
		http.Error(w, "Error when updating machine settings", http.StatusBadRequest)
		return
	}

	delayTime, err := strconv.Atoi(r.FormValue("delay"))
	if err != nil {
		logger.Errorf("Invalid form input: %+v", err)
		http.Error(w, "Error when updating machine settings", http.StatusBadRequest)
		return
	}

	usThreshold, err := strconv.Atoi(r.FormValue(("thresh")))
	if err != nil {
		logger.Errorf("Invalid form input: %+v", err)
		http.Error(w, "Error when updating machine settings", http.StatusBadRequest)
		return
	}

	if config.Config.Machine.WaitingTime != delayTime {
		config.Config.Machine.WaitingTime = delayTime
	}
	if config.Config.Machine.Piezo != usThreshold {
		config.Config.Machine.Piezo = usThreshold
	}
	if config.Config.Machine.Serial != r.FormValue("serial") {
		config.Config.Machine.Serial = r.FormValue("serial")
	}

	err = config.SaveConfig()
	if err != nil {
		logger.Errorf("Error writing config file after update: %+v", err)
		http.Error(w, "Error when updating machine settings", http.StatusBadRequest)
		return
	}

	// restart connector service
	connector.Serv.Restart()

	// Redirect back to admin
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (ws *WebServer) updateScoreboard(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logger.Error(err)
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

	if config.Config.Scoreboard.User != r.FormValue("sbuser") {
		config.Config.Scoreboard.User = r.FormValue("sbuser")
	}

	if config.Config.Scoreboard.Pass != r.FormValue("sbpass") {
		config.Config.Scoreboard.Pass = r.FormValue("sbpass")
	}

	err = config.SaveConfig()
	if err != nil {
		logger.Errorf("Error writing config file after update: %+v", err)
		http.Error(w, "Error when updating scoreborad settings", http.StatusBadRequest)
		return
	}

	// restart connector service
	connector.Serv.Restart()

	// Redirect back to admin
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
