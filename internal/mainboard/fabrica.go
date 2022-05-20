package mainboard

import (
	"sync"

	"github.com/CbIPOKGIT/lift/configs"
	"github.com/CbIPOKGIT/lift/drivers/rs232"
	"github.com/CbIPOKGIT/lift/drivers/rs485"
	"github.com/CbIPOKGIT/lift/internal/board"
)

const MB_READ_INTERVAL = 200

type Boards []*board.Board

type MainBoard struct {
	sync.Mutex

	// Порти
	P232 *rs232.RS232
	P485 *rs485.RS485

	// Інтервал зчитування данних з материнської плати (мс)
	ReadInterval int

	// Поточний статус сенсорів плати
	StatusSensors string

	// Поточний статус реле плати
	StatusRelays string

	// Поточний статус напруги
	StatusVoltage string

	// Список підключених плат
	Boards Boards
}

func New() (*MainBoard, error) {
	mainBoard := new(MainBoard)

	// Встановлюємо значення інтервалу зчитування
	mainBoard.ReadInterval = MB_READ_INTERVAL

	// Підключаємся до порта 232
	if err := mainBoard.createPort232(); err != nil {
		return nil, err
	}

	// Підключаємся до порта 485
	if err := mainBoard.createPort485(); err != nil {
		return nil, err
	}

	// Підключаємо борди
	mainBoard.LoadBoards()

	return mainBoard, nil
}

func (mb *MainBoard) createPort232() error {
	port, err := rs232.NewPort(configs.Rs232Config())
	if err != nil {
		return err
	}

	// port.DoRequest("ATPING")

	port.DoRequest("ATLCDCLEAR")   //Треба переписати на окрему константу
	port.DoRequest("ATLCDLIGHTON") //але я ще не знаю what a fuck is this

	mb.P232 = port
	return nil
}

func (mb *MainBoard) createPort485() error {
	port, err := rs485.NewPort(configs.Rs485Config())
	if err != nil {
		return err
	}

	port.SetStartByte(0xFA)
	port.SetStopByte(0xFE)

	mb.P485 = port
	return nil
}
