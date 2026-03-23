package main

import (
	"fmt"
	"go_final/pkg/db"
	"go_final/pkg/server"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env файл не найден")
	}
}

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	if err := db.Init(dbFile); err != nil {
		fmt.Println("Ошибка БД", err)
		return
	}

	server.StartServer()
}
