package logger

import (
	"strconv"
	"strings"
	"sync"
)

const (
	COMMAND_LIFT_SPEED = "lift_speed"
)

type BoardInterface interface {
	SetData([]byte)
	SetMessageHandler(MessageHandler)
}

func NewBoard(boardType int, handler MessageHandler) BoardInterface {
	var board BoardInterface

	if boardType == 48 {
		board = new(MotorBoard)
	}

	if board == nil {
		return nil
	}
	board.SetMessageHandler(handler)
	return board
}

//--------------------------------- Декодер MOTOR BOARD ---------------------------------
type MotorBoard struct {
	sync.Mutex
	Speed          int
	MessageHandler MessageHandler
}

// Передаємо канал
func (m *MotorBoard) SetMessageHandler(mh MessageHandler) {
	m.MessageHandler = mh
}

// Розпарсюємо та записуємо дані
func (m *MotorBoard) SetData(data []byte) {
	m.Lock()
	defer m.Unlock()

	parts := strings.Split(string(data), "=")
	if len(parts) != 2 {
		return
	}

	speed, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	m.Speed = speed

	m.sendMessage()
}

// Відправляєм повідомлення
func (m *MotorBoard) sendMessage() {
	var status int
	if m.Speed == 0 {
		status = 0
	} else if m.Speed <= 140 && m.Speed >= -100 {
		status = 1
	} else {
		status = 2
	}

	m.MessageHandler(COMMAND_LIFT_SPEED, status, "")
}
