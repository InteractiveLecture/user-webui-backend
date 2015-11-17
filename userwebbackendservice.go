package main

import (
	//"github.com/gorilla/context"

	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/InteractiveLecture/servicecache"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	servicecache.Configure("discovery:8500", 10*time.Second, "authentication-service", "lecture-service")
	//serviceclient.Configure(cacheadapter.New("discovery:8500", 10*time.Second, 5*time.Second, 3))
	servicecache.Start(3, 5*time.Second)
	log.Println("listening on 8000")

	r.Methods("GET").
		Path("/login").
		Handler(createProxy(
		"authentication-service",
		"/oauth/token"))

	// Bind to a port and pass our router in
	http.ListenAndServe(":8000", r)
}

func createProxy(service, servicePath string) http.Handler {
	address, _ := servicecache.GetServiceAddress(service)
	targetURL, err := url.Parse(fmt.Sprintf("http://%s%s", address, servicePath))
	if err != nil {
		panic(err)
	}
	handler := httputil.NewSingleHostReverseProxy(targetURL)
	handler.Director = func(r *http.Request) {
		r.Host = address
		r.URL = targetURL
	}
	return handler
}
