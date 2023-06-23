package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"webapp/platform/authenticator"
	"webapp/platform/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := router.New(auth)

	log.Print("Server listening on http://" + os.Getenv("HTTP_LISTEN_ADDR"))
	if err := http.ListenAndServe(os.Getenv("HTTP_LISTEN_ADDR"), rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
