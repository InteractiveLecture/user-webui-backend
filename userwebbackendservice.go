package main

import (
	//"github.com/gorilla/context"

	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/InteractiveLecture/middlewares/jwtware"
	"github.com/InteractiveLecture/servicecache"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	servicecache.Configure("discovery:8500", 10*time.Second, "authentication-service", "lecture-service")
	//serviceclient.Configure(cacheadapter.New("discovery:8500", 10*time.Second, 5*time.Second, 3))
	servicecache.Start(3, 5*time.Second)
	log.Println("listening on 8000")

	// Posten der Login-daten
	r.Methods("POST").
		Path("/login").
		Handler(createProxy(
		"authentication-service",
		"/oauth/token"))

	// TOPIC Anfragen
	//---------------------

	// TODO Authware und Groupware Nutzung bestimmen

	// Liste aller Topics
	r.Methods("GET").
		Path("/topics").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics")))

	// Erzeuge ein neues Topic
	r.Methods("POST").
		Path("/topics").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics")))

	// Ein spezielles Topic nachfragen
	r.Methods("GET").
		Path("/topics/{id}").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics/{id}")))

	// Verändert ein vorhandenes Topic
	r.Methods("PUT").
		Path("/topics/{id}").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics/{id}")))

	// Löscht ein vorhandenes Topic
	r.Methods("DELETE").
		Path("/topics/{id}").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics/{id}")))

	// Erfragt ein neues Rootmodul
	r.Methods("GET").
		Path("/topics/{id}/modules").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics/{id}/modules")))

	// Erzeugt ein neues Rootmodul
	r.Methods("POST").
		Path("/topics/{id}/modules").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics/{id}/modules")))

	// Liste aller officer eines Topics anfordern
	r.Methods("GET").
		Path("/topics/{id}/officers").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics/{id}/officers")))

	// Einen officer für ein Topic löschen
	r.Methods("DELETE").
		Path("/topics/{id}/officers/{userId}").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics/{id}/officers/{userId}")))

	// Einen officer einem Topic hinfügen
	r.Methods("POST").
		Path("/topics/{id}/officers").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics/{id}/officers")))

	// Liste aller assistant eines Topics anfordern
	r.Methods("GET").
		Path("/topics/{id}/assistants").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics/{id}/assistants")))

	// Einen assistant für ein Topic löschen
	r.Methods("DELETE").
		Path("/topics/{id}/assistants/{userId}").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics/{id}/assistants/{userId}")))

	// Einen assistant einem Topic hinfügen
	r.Methods("POST").
		Path("/topics/{id}/assistants").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/topics/{id}/assistants")))

	// HINTS Anfragen
	//---------------------

	// Einen einzelnen Hint anfragen
	r.Methods("GET").
		Path("/hint/{id}").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/hint/{id}")))

	// Fügt einen weiteren Hint hinzu
	r.Methods("POST").
		Path("/hint/{id}").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/hint/{id}")))

	// Löscht einen Hint
	r.Methods("DELETE").
		Path("/hint/{id}").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/hint/{id}")))

	// Verändert einen Hint
	r.Methods("PUT").
		Path("/hint/{id}").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/hint/{id}")))

	// Konsumiert den angebenen Hint
	r.Methods("POST").
		Path("/hint/{id}/consume").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/hint/{id}/consume")))

	// USER Anfragen
	//---------------------

	// Einen einzelnen User anfragen
	r.Methods("GET").
		Path("/users/{id}").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/users/{id}")))

	// Fügt einen weiteren User hinzu
	r.Methods("POST").
		Path("/users").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/users")))

	// Löscht einen user
	r.Methods("DELETE").
		Path("/users/{id}").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/users/{id}")))

	// Verändert einen user
	r.Methods("PUT").
		Path("/users/{id}").
		Handler(jwtware.New(createProxy(
		"lecture-service",
		"/users/{id}")))

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
