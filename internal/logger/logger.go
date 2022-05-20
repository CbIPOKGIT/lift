package logger

import (
	"log"
	"strconv"
	"sync"

	wsconnections "github.com/CbIPOKGIT/lift/pkg/ws-connections"
	"github.com/CbIPOKGIT/lift/protos"
)

// Функція, що відправляє дані на сервер
type MessageHandler func(string, int, string)

// Структура повідомлення від датчика до сервера
type MessageToServer struct {
	Command string `json:"cmd"`
	Status  string `json:"status"`
	Message string `json:"msg"`
}

// Канал для доставки повідомлення від інтерфейса плати
type MessageToServerChannel chan *MessageToServer

// Список плат, з яких приходять сигнали
type BoardsMapa struct {
	sync.Mutex
	Mapa map[int]BoardInterface
}

// Теоритично це буде наш основний посередник
//
// Приймає дані від сенсорів
type Logger struct {
	// Стан материнської плати
	MainBoard *MainBoard

	// Список активних плат
	Boards BoardsMapa

	// Канал для отримання розпарсених данних та відправки їх на сервер
	BoardsChannel MessageToServerChannel

	// Сокет зв'язку з сервером
	ServerSocket *wsconnections.WebsocketClient

	// Канал передачі команд материнській платі
	MainboardCommandsChannel protos.MainboardCommandChannel
}

func New() *Logger {
	logger := new(Logger)

	logger.Boards.Mapa = make(map[int]BoardInterface)

	logger.BoardsChannel = make(MessageToServerChannel, 255)

	logger.MainBoard = new(MainBoard)

	// Канал передачі показників на сервер
	logger.MainBoard.MessageHandler = logger.sendMessageToServer

	// Створюємо канал для передачі команд на материнську плату
	logger.MainboardCommandsChannel = make(protos.MainboardCommandChannel)

	logger.CreateServerSocket()

	go logger.Listen()

	return logger
}

func (l *Logger) CreateServerSocket() {
	l.ServerSocket = wsconnections.New()

	if err := l.ServerSocket.Connect(); err != nil {
		log.Println("Error connect to server")
	} else {
		log.Println("Successfull connected to server")
	}
}

// Віддаємо канал команд материнській платі
func (l *Logger) GetCommandsBus() protos.MainboardCommandChannel {
	return l.MainboardCommandsChannel
}

// Форматуємо і відправляємо повідомлення на сервер
func (l *Logger) sendMessageToServer(command string, status int, message string) {
	data := &MessageToServer{
		Command: command,
		Message: message,
		Status:  strconv.Itoa(status),
	}

	l.BoardsChannel <- data
}
