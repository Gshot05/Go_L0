package main

import (
	"database/sql"
	"fmt"
	"log"
	"myapp/internal/handler"
	"myapp/internal/service"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

var (
	db *sql.DB
	nc *nats.Conn
)

func main() {
	//Подключение к базе данных
	connStr := "user=goadmin password=P@$$vv0RD dbname=GoDataBase sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Подключение к натс
	natsURL := "localhost"
	log.Println("Подключение NATS...")
	nc, err = nats.Connect(natsURL)
	if err != nil {
		log.Println("ошибка подключения к NATS")
		log.Fatal(err)
	}
	log.Println("Подключение к NATS успешно")
	defer nc.Close()

	databaseService := service.NewDatabaseService(db, nc)
	databaseHandler := handler.NewDatabaseHandler(databaseService)

	http.HandleFunc("/", databaseHandler.Index)
	http.HandleFunc("/get-info", databaseHandler.GetInfo)
	http.HandleFunc("/add-info", databaseHandler.AddData)

	fmt.Println("Запуск сервера...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
