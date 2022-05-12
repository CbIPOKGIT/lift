package test

import (
	"log"
	"time"

	"github.com/CbIPOKGIT/lift/drivers/rs485"
	"github.com/CbIPOKGIT/lift/internal/mainboard"
)

func Test(mb *mainboard.MainBoard) {
	boardIdConv := 0

	boardId := uint8(boardIdConv)
	boardId++

	for i := 0; i < 3; i++ {
		log.Println("Board id ", boardId)
		cpuId, errSr := searchBoard(mb.P485)
		if errSr == nil {
			log.Println("CPU id - ", cpuId)
		} else {
			log.Println("Error")
		}
		boardId++
		time.Sleep(time.Second * 5)
	}

	log.Fatal("Stop")
}

func searchBoard(port *rs485.RS485) (uint64, error) {

	rsp := port.Search()
	log.Println("Search resp")
	log.Println(rsp)
	if rsp.Err != nil {
		return 0, rsp.Err
	}

	return rsp.Cpuid, nil
}
