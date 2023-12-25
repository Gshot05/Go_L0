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

