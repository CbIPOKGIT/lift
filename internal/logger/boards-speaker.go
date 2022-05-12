package logger

import (
	"github.com/CbIPOKGIT/lift/protos"
)

// Додаємо новий інтерфейс датчика
func (l *Logger) AddBoard(boardType int) BoardInterface {
	board := NewBoard(boardType, l.BoardsChannel)

	l.Boards.Lock()
	defer l.Boards.Unlock()

	l.Boards.Mapa[boardType] = board

	return board
}

// Отримуємо дані з датчика та записуємо в відповідний інтерфейс
func (l *Logger) ReciveFromBoard(message *protos.BoardMessage) {
	board, has := l.Boards.Mapa[message.BoardType]
	if !has {
		board = l.AddBoard(message.BoardType)
	}

	if board == nil {
		return
	}

	go board.SetData(message.Message)
}

func (l *Logger) ReciveFromMainboard(data *protos.MainboardMessage) {

	handlers := map[bool]func(string){
		true:  l.MainBoard.SetSensorMessage,
		false: l.MainBoard.SetRelayMessage,
	}

	go handlers[data.Sensor](data.Message)
}
