package main

import (
	"net/url"
	"net/http"
	"net/http/httputil"
)

// static
func StaticHandler(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}

// proxy
func ProxyHandler(path string) http.Handler {
	targetUrl, _ := url.Parse(path)
	return httputil.NewSingleHostReverseProxy(targetUrl)
}
