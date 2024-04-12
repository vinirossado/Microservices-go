package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "8081"

type Config struct{}

func main() {
	app := &Config{}

	log.Printf("Starting broker service on port %s\n", webPort)

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := svr.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
