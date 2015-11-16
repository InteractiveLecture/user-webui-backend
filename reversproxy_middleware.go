package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func proxyHandlerFunc(hostname string) {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   hostname,
	})
	proxy.Director = proxyDirector
	http.Handle("/", proxy)
	http.ListenAndServe(":8080", nil)
}

func proxyDirector(r *http.Request) {
	r.Host = "www.nooooooooooooooo.com"
	r.URL.Host = r.Host
	r.URL.Scheme = "http"
}

// New let you use the Middleware
func New(next http.Handler, hostname string) http.Handler {
	proxyHandlerFunc(hostname)
	return next
}
