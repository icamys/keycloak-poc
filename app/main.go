package main

import (
	"log"
	"net/http"

	"app/middleware"
)

func main() {
	config, err := Parse()
	if err != nil {
		log.Fatalf("There was an error while parsing config: %v", err)
	}

	middleware.InitializeOauthServer(&middleware.AuthorizationConfig{
		ClientId:     config.Authentication.ClientId,
		ClientSecret: config.Authentication.ClientSecret,
		Realm:        config.Authentication.Realm,
		Hostname:     config.Authentication.Hostname,
	})

	router := http.NewServeMux()

	router.Handle("/public", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Hello from a public endpoint! You don't need to be authenticated to see this."}`))
	}))

	router.Handle("/private/dashboard", middleware.Authorization(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Hello from the private dashboard!"}`))
		}),
	))

	router.Handle("/private/products", middleware.Authorization(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Hello from a private products page!"}`))
		}),
	))

	log.Print("Listening on http://" + config.HTTP.ListenAddr)
	if err := http.ListenAndServe(config.HTTP.ListenAddr, router); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
