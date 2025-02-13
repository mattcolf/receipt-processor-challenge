package main

import (
	"log"
	"net/http"

	"github.com/mattcolf/receipt-processor-challenge/api"
)

func main() {
	config := api.LoadConfig()
	api := api.SetupApi()

	server := &http.Server{
		Handler:      api.Router,
		Addr:         config.ServerBindAddress,
		WriteTimeout: config.ServerWriteTimeout,
		ReadTimeout:  config.ServerReadTimeout,
	}

	// TODO: add graceful shutdown

	log.Fatal(server.ListenAndServe())
}
