package main

import (
	"DB-project/cmd/database"
	"DB-project/cmd/server"
	"log"
)

func main() {
	// Инициализация базы данных
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Запуск веб-сервера
	server.StartServer(db)
}
