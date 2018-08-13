package main

import (
	"os"
	"net/http"
	"os/signal"
)

type Application struct {
	logger  Logger
	setting Setting
	servers []*http.Server
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
	servers := set.Servers()
	app.servers = make([]*http.Server, len(servers))
	for i, s := range servers {
		app.servers[i] = app.setting2http(s)
	}
	return nil
}

// setting server to http server
func (app *Application) setting2http(s ServerSetting) *http.Server {
	mux := http.NewServeMux()
	for _, l := range s.Locations() {
		h := app.location2handler(l)
		mux.Handle(l.Pattern(), LoggingHandler(s.Name(), app.logger.Print, h))
	}
	return &http.Server{Addr: s.Listen(), Handler: mux}
}

// setting location to http handler
func (app *Application) location2handler(l LocationSetting) http.Handler {
	if l.IsRoot() {
		return StaticHandler(l.Root())
	} else if l.IsProxy() {
		return ProxyHandler(l.Proxy())
	} else if l.IsProxies() {
		proxies := l.Proxies()
		if len(proxies) == 1 {
			return ProxyHandler(proxies[0])
		} else {
			hs := ProxiesHandler(proxies)
			switch l.Mode() {
			case "round":
				return RoundHandler(hs)
			case "random":
				return RandomHandler(hs)
			case "iphash":
				return IphashHandler(hs)
			}
		}
	}
	return nil
}

// check error and exit
func (app *Application) check(err error) {
	if err != nil {
		println(err)
		os.Exit(1)
	}
}

// Start all of server
func (app *Application) Start() {
	for _, s := range app.servers {
		go s.ListenAndServe()
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<- c
}

// Stop all of server
func (app *Application) Stop() {
	for _, s := range app.servers {
		s.Close()
	}
	app.logger.Close()
}

func main() {
	app := &Application{}
	err := app.Init()
	app.check(err)
	app.Start()
	app.Stop()
}
