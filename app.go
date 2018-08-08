package main

import "os"

type Application struct {
	log Logger
	set Setting
}

// Init log & set
func (app *Application) Init() error {
	log, err := NewLogger("access.log", L_DEBUG)
	if err != nil {
		return err
	}
	app.log = log
	set, err := NewSetting("setting.json")
	if err != nil {
		return err
	}
	app.set = set
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
