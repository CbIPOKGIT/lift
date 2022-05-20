package logger

import (
	"github.com/CbIPOKGIT/lift/protos"
)

// Додаємо новий інтерфейс датчика
func (l *Logger) AddBoard(boardType int) BoardInterface {
	board := NewBoard(boardType, l.sendMessageToServer)

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

	handlers := map[int]func(string){
		1: l.MainBoard.SetSensorMessage,
		2: l.MainBoard.SetRelayMessage,
		3: l.MainBoard.SetVoltageMethod,
	}

	go handlers[data.Type](data.Message)
}
