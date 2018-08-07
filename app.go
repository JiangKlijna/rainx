package main

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
