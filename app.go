package main

import (
	"os"
	"net/http"
)

type Application struct {
	logger  Logger
	setting Setting
	servers []http.Server
}

// Init log & set
func (app *Application) Init() error {
	log, err := NewLogger("access.log", L_DEBUG)
	if err != nil {
		return err
	}
	app.logger = log
	set, err := NewSetting("setting.json")
	if err != nil {
		return err
	}
	err = set.IsValid()
	if err != nil {
		return err
	}
	app.setting = set
	return nil
}

// Start all of server
func (app *Application) Start() {
}

func main() {
	app := &Application{}
	err := app.Init()
	if err != nil {
		println(err)
		os.Exit(1)
	}
	app.Start()
}
