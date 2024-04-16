package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"os"
	"time"
)

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
	log.Println("Connected to RabbitMQ")
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	for {
		conn, err := amqp.Dial("amqp://guest:guest@localhost")
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
