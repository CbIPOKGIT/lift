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
	SetChannel(MessageToServerChannel)
}

func NewBoard(boardType int, channel MessageToServerChannel) BoardInterface {
	var board BoardInterface

	if boardType == 48 {
		board = new(MotorBoard)
	}

	if board == nil {
		return nil
	}
	board.SetChannel(channel)
	return board
}

//--------------------------------- Декодер MOTOR BOARD ---------------------------------
type MotorBoard struct {
	sync.Mutex
	Speed    int
	ToServer MessageToServerChannel
}

// Передаємо канал
func (m *MotorBoard) SetChannel(channel MessageToServerChannel) {
	m.ToServer = channel
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
	message := &MessageToServer{
		Command: COMMAND_LIFT_SPEED,
	}

	if m.Speed == 0 {
		message.Status = "0"
	} else if m.Speed <= 140 && m.Speed >= -100 {
		message.Status = "1"
	} else {
		message.Status = "2"
	}

	m.ToServer <- message
}
