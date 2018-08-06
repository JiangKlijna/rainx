package main

import (
	"net/http"
)

// static
func StaticHandler(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}
