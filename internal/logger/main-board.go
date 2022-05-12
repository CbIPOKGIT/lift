package logger

import (
	"errors"
	"strings"

	"github.com/CbIPOKGIT/lift/pkg/conv"
)

const WRONG_MESSAGE_FORMAT = "wrong message format"

type MainBoardSensorsStatus struct {
	Door bool
}

//Декодер материнської плати
type MainBoard struct {
	// Щойно ініціалізований. Щоб на початку ми обов'язково відправили початкові статуси
	JustInited bool

	// Стани сенсорів
	MainBoardSensorsStatus

	// Канал передачі даних на сервер
	ToServer MessageToServerChannel
}

// Обробляємо статус сенсорів плати
func (m *MainBoard) SetSensorMessage(message string) {
	status, err := m.decodeSensorMessage(message)
	if err != nil {
		// Щось будемо робити
	}

	if m.JustInited || status.Door != m.Door {
		m.Door = status.Door
		m.sendMessageToServer("lift_door", m.Door, "")
	}

	m.JustInited = false
}

// Обробляємо статус реле плати
func (m *MainBoard) SetRelayMessage(message string) {

}

// Обробляємо статус сенсорів плати
func (MainBoard) decodeSensorMessage(message string) (*MainBoardSensorsStatus, error) {
	parts := strings.Split(message, "=")
	if len(parts) != 2 {
		return nil, errors.New(WRONG_MESSAGE_FORMAT)
	}

	binary := conv.NewBinary()
	binary.SetString(parts[1])

	status := new(MainBoardSensorsStatus)

	status.Door = binary.GetBitUnsafe(3)

	return status, nil
}

// Форматуємо і відправляємо повідомлення на сервер
func (m *MainBoard) sendMessageToServer(command string, status bool, message string) {
	data := &MessageToServer{
		Command: command,
		Message: message,
	}

	if status {
		data.Status = "1"
	} else {
		data.Status = "0"
	}

	m.ToServer <- data
}
