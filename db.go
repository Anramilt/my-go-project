package main

import (
	"fmt"
	"log"

	//"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB //Переменная сессии БД

const (
	host     = "localhost"
	port     = 5432
	user     = "testadmin"
	password = "12345678"
	dbname   = "godb"
)

func ConnectDB() *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error connected")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connected")
	}

	fmt.Println("Successfully connected!")
	return db
}

func GetDB() *sqlx.DB {
	return db
}
