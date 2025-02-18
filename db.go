package main

import (
	"fmt"

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

// Структура для хранения данных из таблицы testone
type TestOne struct {
	ID  int    `db:"id"`
	One string `db:"one"`
}

// Структура для хранения эхо
type EchoMessage struct {
	ID      int    `db:"id"`
	Message string `db:"message"`
}

func ConnectDB() error {
	var err error
	db, err = sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname))
	fmt.Println("Successfully connected!")
	return err
}

func getEchomessageList() ([]EchoMessage, error) {
	query := `SELECT id, message FROM echomessage`
	var messages []EchoMessage
	err := db.Select(&messages, query)
	if err != nil {
		return nil, fmt.Errorf("error selecting rows in echomessage: %w", err)
	}
	return messages, nil
}

func addEchotoDB(message string) error {
	query := `INSERT INTO echomessage (message) VALUES ($1)`
	_, err := db.Exec(query, message)
	if err != nil {
		return fmt.Errorf("error inserting echo: %w", err)
	}
	fmt.Println("Echo added in table: ", message)
	return nil
}

/*func getTestoneList() ([]TestOne, error) {
	query := `SELECT id, one FROM testone`
	var rows []TestOne
	err := db.Select(&rows, query)
	if err != nil {
		return nil, fmt.Errorf("error selecting rows in testone: %w", err)
	}
	return rows, nil
}*/

/*
func addRowTestone(id int, one string) error {
	query := `INSERT INTO testone (id, one) VALUES ($1, $2)`
	_, err := db.Exec(query, id, one)
	if err != nil {
		return fmt.Errorf("error inserting row: %w", err)
	}
	return nil
}*/

/*
func getTable() ([]string, error) {
	query := `SELECT tablename FROM pg_catalog.pg_tables WHERE table_schema = 'public';`
	var tables []string
	err := db.Select(&tables, query)
	if err != nil {
		log.Fatalf("Error selected tables list: %v", err)
	}
	return tables, nil
}*/

/*func GetDB() *sqlx.DB {
	return db
}*/
