package handler

import (
	"fmt"
	"html/template"
	"log"
	"myapp/internal/service"
	"net/http"
	"strconv"
)

type DatabaseHandler struct {
	databaseService *service.DatabaseService
}

func NewDatabaseHandler(databaseService *service.DatabaseService) *DatabaseHandler {
	return &DatabaseHandler{
		databaseService: databaseService,
	}
}

func (h *DatabaseHandler) Index(w http.ResponseWriter, r *http.Request) {
	// Парсим шаблон
	tmpl, err := template.ParseFiles("../../web/index.html")
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

func (h *DatabaseHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	numberStr := r.FormValue("number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	jsonData, err := h.databaseService.GetInfo(number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Парсим шаблон
	tmpl, err := template.ParseFiles("../../web/index.html")
	if err != nil {
		http.Error(w, "Oшибка парсинга шаблона", http.StatusInternalServerError)
		log.Println("Oшибка парсинга шаблона:", err)
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

func (h *DatabaseHandler) AddData(w http.ResponseWriter, r *http.Request) {
	// Получаем данные из формы
	jsonData := r.FormValue("jsonData")

	// Вызываем метод для добавления данных
	id, err := h.databaseService.AddData(jsonData)
	if err != nil {
		http.Error(w, "Ошибка при добавлении данных", http.StatusInternalServerError)
		log.Printf("Ошибка при добавлении данных: %v", err)
		return
	}

	// Выводим какую-то информацию о добавленных данных
	fmt.Fprintf(w, "Добавлены новые данные. ID: %d", id)
}
