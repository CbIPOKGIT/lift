package logger

import (
	"errors"
	"strconv"
	"strings"

	"github.com/CbIPOKGIT/lift/pkg/conv"
)

const WRONG_MESSAGE_FORMAT = "wrong message format"

type MainBoardSensorsStatus struct {
	Inited bool

	// Статус двері ліфтової
	DoorMachine int

	// Тампер
	CoverDevice int
}

type MainBoardVoltageStatus struct {
	Inited bool

	// Напруга блока живлення
	NetVoltage int

	// Напруга аккумулятора
	AccumVolatege int

	// Напруга живлення міні ПК
	MiniPcVoltage int

	// Працюємо від мережі
	PowerSupplyNetwork int
}

//Декодер материнської плати
type MainBoard struct {
	// Стани сенсорів
	MainBoardSensorsStatus

	// Стан по живленню
	MainBoardVoltageStatus

	// Канал передачі даних на сервер
	MessageHandler MessageHandler
}

// Обробляємо статус сенсорів плати
func (m *MainBoard) SetSensorMessage(message string) {
	status, err := m.decodeSensorMessage(message)
	if err != nil {
		// Щось будемо робити
		return
	}

	if !m.MainBoardSensorsStatus.Inited || status.DoorMachine != m.DoorMachine {
		m.DoorMachine = status.DoorMachine
		m.MessageHandler("machine_room_door", m.DoorMachine, "")
	}

	if !m.MainBoardSensorsStatus.Inited || status.CoverDevice != m.CoverDevice {
		m.CoverDevice = status.CoverDevice
		m.MessageHandler("cover_device", m.CoverDevice, "")
	}

	m.MainBoardSensorsStatus.Inited = true
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

	if binary.GetBitBoolValue(3) {
		status.DoorMachine = 1
	} else {
		status.DoorMachine = 0
	}

	if binary.GetBitBoolValue(4) {
		status.CoverDevice = 1
	} else {
		status.CoverDevice = 0
	}

	return status, nil
}

// Обробляємо статус реле плати
func (m *MainBoard) SetRelayMessage(message string) {

}

// Обробляємо дані напруги
func (m *MainBoard) SetVoltageMethod(message string) {
	status, err := m.decodeVoltageMessage(message)
	if err != nil {
		// Щось будемо робити
		return
	}

	if !m.MainBoardVoltageStatus.Inited || m.MainBoardVoltageStatus.PowerSupplyNetwork != status.PowerSupplyNetwork {
		m.MessageHandler("controller_power", status.PowerSupplyNetwork, "")
	}

	m.MainBoardVoltageStatus = *status

	m.MainBoardVoltageStatus.Inited = true
}

// Обробляємо статус сенсорів плати
func (MainBoard) decodeVoltageMessage(message string) (*MainBoardVoltageStatus, error) {
	parts := strings.Split(message, "=")
	if len(parts) != 2 {
		return nil, errors.New(WRONG_MESSAGE_FORMAT)
	}

	parts = strings.Split(parts[1], ",")
	if len(parts) != 3 {
		return nil, errors.New(WRONG_MESSAGE_FORMAT)
	}

	nVoltage, errN := strconv.Atoi(parts[0])
	aVoltage, errA := strconv.Atoi(parts[1])
	mVoltage, errM := strconv.Atoi(parts[2])

	if errA != nil || errN != nil || errM != nil {
		return nil, errors.New(WRONG_MESSAGE_FORMAT)
	}

	status := MainBoardVoltageStatus{
		AccumVolatege: aVoltage,
		NetVoltage:    nVoltage,
		MiniPcVoltage: mVoltage,
	}

	if status.NetVoltage > 0 {
		status.PowerSupplyNetwork = 1
	} else {
		status.PowerSupplyNetwork = 0
	}

	return &status, nil
}
