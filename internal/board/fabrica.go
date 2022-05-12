package board

import (
	"log"
	"sync"

	"github.com/CbIPOKGIT/lift/drivers/rs485"
)

type BoardData struct {
	Id           uint8
	CpuId        uint64
	ReadInterval uint16
}

type Board struct {
	sync.Mutex
	Id           uint8
	CpuId        uint64
	ReadInterval uint16
	Status       uint8
	Name         string
	CurrentData  string
	BoardType    rs485.BoardTypes_t
	Port         *rs485.RS485
}

func New(port *rs485.RS485, data BoardData) (*Board, error) {
	board := &Board{
		Id:           data.Id,
		CpuId:        data.CpuId,
		ReadInterval: data.ReadInterval,
		Port:         port,
		Status:       0,
	}

	if err := port.SetAddr(board.Id, board.CpuId); err != nil {
		return nil, err
	}

	boardType, err := port.GetBoardType(board.Id)
	if err != nil {
		log.Println("Error get board type")
		log.Println(err)
		boardType = rs485.DoorBoard
	}

	board.BoardType = boardType
	return board, nil
}
