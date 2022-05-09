package mainboard

import (
	"github.com/CbIPOKGIT/lift/configs"
	"github.com/CbIPOKGIT/lift/drivers/rs232"
	"github.com/CbIPOKGIT/lift/drivers/rs485"
	"github.com/CbIPOKGIT/lift/internal/board"
	"github.com/CbIPOKGIT/lift/internal/storage"
)

type Boards []*board.Board

type MainBoard struct {
	P232 *rs232.RS232
	P485 *rs485.RS485

	// Список підключених плат
	Boards Boards

	// Локальне сховище для локальних данних
	Storage storage.Storage
}

func New() (*MainBoard, error) {
	mainBoard := new(MainBoard)

	// Підключаємся до порта 232
	if err := mainBoard.createPort232(); err != nil {
		return nil, err
	}

	// Підключаємся до порта 485
	if err := mainBoard.createPort485(); err != nil {
		return nil, err
	}

	// Підключаємо сховище
	mainBoard.Storage = storage.New()

	// Підключаємо борди
	mainBoard.LoadBoards()

	return mainBoard, nil
}

func (mb *MainBoard) createPort232() error {
	port, err := rs232.NewPort(configs.Rs232Config())
	if err != nil {
		return err
	}

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

	port.SetStartByte(configs.Rs485StartByte)
	port.SetStopByte(configs.Rs485StopByte)

	mb.P485 = port
	return nil
}
