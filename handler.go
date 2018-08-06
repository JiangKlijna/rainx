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

// iphash
func IphashHandlerh(hs []http.Handler) http.Handler {
	n := uint32(len(hs))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hs[IphashByAddress(r.RemoteAddr)%n].ServeHTTP(w, r)
	})
}
