package service

import (
	"database/sql"
	"fmt"

	"github.com/nats-io/nats.go"
)

type Cache struct {
	OrderUID string `json:"order_uid"`
}

type DatabaseService struct {
	db    *sql.DB
	nc    *nats.Conn
	cache map[int]Cache
}

func NewDatabaseService(db *sql.DB, nc *nats.Conn) *DatabaseService {
	return &DatabaseService{
		db:    db,
		nc:    nc,
		cache: make(map[int]Cache),
	}
}

func (s *DatabaseService) GetInfo(number int) (string, error) {
	var jsonData string
	err := s.db.QueryRow("SELECT Name_Json_Info FROM Json_Info WHERE ID_Json_Info = $1", number).Scan(&jsonData)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса к базе данных")
	}

	message := fmt.Sprintf("Выполнен запрос к бд с номером id: %d", number)
	s.nc.Publish("log", []byte(message))

	return jsonData, nil
}

func (s *DatabaseService) AddAndCacheData(orderUID string) (int, string, error) {
	// Вставляем новые данные в базу и получаем возвращенный ID
	var newID int
	err := s.db.QueryRow("INSERT INTO Json_Info (Name_Json_Info) VALUES ($1) RETURNING ID_Json_Info", orderUID).Scan(&newID)
	if err != nil {
		return 0, "", fmt.Errorf("ошибка при добавлении данных в базу: %v", err)
	}

	// Сохраняем данные в кэше
	newCacheData := Cache{OrderUID: orderUID}
	s.cache[newID] = newCacheData

	// Отправляем сообщение в NATS
	message := fmt.Sprintf("Добавлены новые данные: %s", orderUID)
	s.nc.Publish("log", []byte(message))

	return newID, orderUID, nil
}
