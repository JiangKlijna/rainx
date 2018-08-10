package main

import (
	"fmt"
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
	for i, s := range path  {
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
func LoggingHandler(print func(v ...interface{}), next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		str := fmt.Sprintf("%s Comleted %s %s in %v from %s\n",
			start.Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			time.Since(start),
			r.RemoteAddr)
		go print(str)
	})
}
