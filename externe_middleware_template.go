package main

import "net/http"

func searchUrlMiddleware(nextHandler http.Handler) http.Handler {
	result := func(w http.ResponseWriter, r *http.Request) {

	}
	return http.HandlerFunc(result)
}
