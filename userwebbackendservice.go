package main

import (
	//"github.com/gorilla/context"

	"log"
	"net/http"
	"time"

	"github.com/InteractiveLecture/middlewares/aclware"
	"github.com/InteractiveLecture/middlewares/jwtware"
	"github.com/InteractiveLecture/serviceclient"
	"github.com/InteractiveLecture/serviceclient/cacheadapter"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	adapter := cacheadapter.New("discovery:8500", 10*time.Second, 5*time.Second, 3)
	serviceclient.Configure(adapter, "acl-service", "authentication-service")
	//
	r.Methods("GET").
		Path("/topics").
		Handler(jwtware.New(
			authware.New(
				New(nil, "lecture-service"), nil, "user"
					)
				)
			)
			
	r.Methods("POST").
		Path("/topics")
	//Handler()
	log.Println("listening on 8000")
	// Bind to a port and pass our router in
	http.ListenAndServe(":8000", r)
}
