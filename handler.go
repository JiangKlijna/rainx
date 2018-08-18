package main

import (
	"time"
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

// proxy array
func ProxiesHandler(path []string) []http.Handler {
	hs := make([]http.Handler, len(path))
	for i, s := range path {
		hs[i] = ProxyHandler(s)
	}
	return hs
}

// iphash
func IphashHandler(hs []http.Handler) http.Handler {
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

// logging
func LoggingHandler(tag string, printf func(string, ...interface{}), next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Add("Server", Version)
		next.ServeHTTP(w, r)
		go printf("%d-%d-%d %d:%d:%d [%s] Comleted %s %s in %v from %s\n",
			start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute(), start.Second(),
			tag, r.Method, r.URL.Path, time.Since(start), r.RemoteAddr)
	})
}
