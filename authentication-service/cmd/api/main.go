package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "8082"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service...")

	conn := connectToDB()

	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := svr.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := OpenDB(dsn)

		if err != nil {
			log.Println("Error connecting to DB, retrying in 2 seconds")
			counts++

		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}
		time.Sleep(2 * time.Second)
		continue
	}
}
