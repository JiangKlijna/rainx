package main

import (
	"net/url"
	"net/http"
	"math/rand"
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

// random
func RandomHandler(hs []http.Handler) http.Handler {
	n := len(hs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hs[rand.Intn(n)].ServeHTTP(w, r)
	})
}

// round
func RoundHandler(hs []http.Handler) http.Handler {
	n := len(hs)
	i := 0
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if i == n {
			i = 0
		}
		hs[i].ServeHTTP(w, r)
		i++
	})
}
