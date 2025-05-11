package main

import (
	"context"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/philipslstwoyears/calculator-go/internal/agent"
	calc "github.com/philipslstwoyears/calculator-go/internal/calculator"
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"github.com/philipslstwoyears/calculator-go/internal/server"
	"github.com/philipslstwoyears/calculator-go/internal/storage"
	"log"
	_ "modernc.org/sqlite"
	"time"
)

func main() {
	err := godotenv.Load("cmd/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	data := storage.New(db)
	ch := make(chan dto.Expression)
	workers := calc.New(data, ch)
	agent := agent.New(data, ch)
	go func() {
		if err := agent.RunServer(); err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(time.Second)
	serv, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	if err := workers.Start(); err != nil {
		log.Fatal(err)
	}
	if err := serv.RunServer(); err != nil {
		log.Fatal(err)
	}
}

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "store.db")
	if err != nil {
		return nil, err
	}

	err = db.PingContext(context.Background())
	if err != nil {
		return nil, err
	}
	err = CreateTables(context.Background(), db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTables(ctx context.Context, db *sql.DB) error {
	const (
		usersTable = `
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			login TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);`

		createUniqueLoginIndex = `
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_login ON users(login);`

		expressionsTable = `
	CREATE TABLE IF NOT EXISTS expressions (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		expression TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		status TEXT NOT NULL,
		result FLOAT,
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);`
	)

	// Создаём таблицы
	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, createUniqueLoginIndex); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, expressionsTable); err != nil {
		return err
	}

	return nil
}
