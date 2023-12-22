package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

// Определение структуры Cache
type Cache struct {
	OrderUID string `json:"order_uid"`
}

// Объявление глобальных переменных
var (
	db        *sql.DB
	cacheJson map[string]Cache
	nc        *nats.Conn
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

	fmt.Println("Запуск сервера...")

	//Подключение к натс
	natsURL := "localhost:"
	log.Println("Подключение NATS...")
	nc, err = nats.Connect(natsURL)
	if err != nil {
		log.Println("Ошибка подключения к NATS")
		log.Fatal(err)
	}
	log.Println("Подключение к NATS успешно")
	defer nc.Close()

	// Инициализация map
	cacheJson = make(map[string]Cache)

	// Добавим обработчик HTTP-запросов для веб-интерфейса
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/get-info", handleGetInfo)

	// Запуск сервера на порту 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// Парсим шаблон
	tmpl, err := template.ParseFiles("../static/index.html")
	if err != nil {
		http.Error(w, "Ошибка парсинга шаблона", http.StatusInternalServerError)
		log.Println("Ошибка парсинга шаблона:", err)
		return
	}

	// Отображаем шаблон
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
		log.Println("Ошибка отображения шаблона:", err)
		return
	}
}

func handleGetInfo(w http.ResponseWriter, r *http.Request) {
	// Получаем значение числа из формы
	numberStr := r.FormValue("number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	// Выполняем запрос к базе данных
	var jsonData string
	err = db.QueryRow("SELECT Name_Json_Info FROM Json_Info WHERE ID_Json_Info = $1", number).Scan(&jsonData)
	if err != nil {
		http.Error(w, "Ошибка запроса к базе данных", http.StatusInternalServerError)
		log.Printf("Ошибка запроса к базе данных: %v", err)
		return
	}

	message := fmt.Sprintf("Выполнен запрос к бд с номером id: %d", number)
	nc.Publish("log", []byte(message))

	// Сохраняем полученную информацию в мапе cacheJson
	var cacheData Cache
	if err := json.Unmarshal([]byte(jsonData), &cacheData); err != nil {
		http.Error(w, "Ошибка декондинга JSON", http.StatusInternalServerError)
		log.Printf("Ошибка декондинга JSON: %v", err)
		return
	}

	cacheJson[numberStr] = cacheData

	// Парсим шаблон
	tmpl, err := template.ParseFiles("../static/index.html")
	if err != nil {
		http.Error(w, "Ошибка парсинга шаблона", http.StatusInternalServerError)
		log.Println("Ошибка парсинга шаблона:", err)
		return
	}

	// Отображаем шаблон с информацией
	data := struct {
		Info   string
		Number int
	}{
		Info:   jsonData,
		Number: number,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
		log.Println("Ошибка отображения шаблона:", err)
		return
	}
}
