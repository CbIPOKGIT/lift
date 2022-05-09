package board

import (
	"github.com/CbIPOKGIT/lift/drivers/rs485"
)

type Board struct {
	Id           uint8
	CpuId        uint64
	ReadInterval uint16
	Status       uint8
	Name         string
	CurrentData  string
	BoardType    rs485.BoardTypes_t
	Port         *rs485.RS485
}

func New() (*Board, error) {
	return nil, nil
}
