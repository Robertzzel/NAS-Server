package Services

import (
	"NAS-Server-Web/shared"
	"NAS-Server-Web/shared/configurations"
	"NAS-Server-Web/shared/models"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net"
	"strconv"
)

type DatabaseService struct {
	conn *shared.MessageHandler
}

func NewDatabaseService() (*DatabaseService, error) {
	address := configurations.GetDatabaseHost() + ":" + configurations.GetDatabasePort()
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	mh := shared.NewMessageHandler(conn)
	return &DatabaseService{mh}, nil
}

func (db *DatabaseService) CheckUsernameAndPassword(username, password string) (bool, error) {
	request := models.NewRequestMessage(0, []string{username, password})
	_ = db.conn.Write(request.GetBytesData())

	rawMsg, err := db.conn.Read()
	if err != nil {
		return false, err
	}

	response := models.NewResponseMessageFromBytes(rawMsg)
	return response.Body[0] == 1, nil
}

func (db *DatabaseService) GetUserAllocatedMemory(username string) (int, error) {
	request := models.NewRequestMessage(1, []string{username})
	_ = db.conn.Write(request.GetBytesData())

	rawMsg, err := db.conn.Read()
	if err != nil {
		return 0, err
	}

	response := models.NewResponseMessageFromBytes(rawMsg)

	memory, err := strconv.Atoi(string(response.Body))
	if err != nil {
		return 0, err
	}

	return memory, nil
}

func (db *DatabaseService) AddUser(username, password string, memory int) bool {
	request := models.NewRequestMessage(2, []string{username, password, fmt.Sprint(memory)})
	_ = db.conn.Write(request.GetBytesData())

	rawMsg, err := db.conn.Read()
	if err != nil {
		return false
	}

	response := models.NewResponseMessageFromBytes(rawMsg)
	return response.Status == 0
}