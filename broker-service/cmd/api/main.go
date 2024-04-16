package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer func(rabbitConn *amqp.Connection) {
		err := rabbitConn.Close()
		if err != nil {

		}
	}(rabbitConn)

	app := &Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = svr.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = conn
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}
		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("RabbitMQ not ready, backing off for", backoff)
		time.Sleep(backoff)
		continue
	}

	return connection, nil
}
