package board

import (
	"log"
	"time"

	"github.com/CbIPOKGIT/lift/drivers/rs485"
)

func (board *Board) ReadData() {
	timer := time.NewTicker(time.Millisecond * time.Duration(board.ReadInterval))
	defer timer.Stop()

	// Команда, за допомогою якої виконується опитування
	var command string

	if board.BoardType == rs485.MotorBoard {
		command = "ATSPEED?"
	} else {
		command = "ATS?"
	}

	for {
		select {
		case <-timer.C:
			log.Println(board.BoardType)
			response := board.Port.DoRequest(board.Id, command)

			if response.Err == nil {
				data := string(response.Response)
				log.Println(data)
			} else {
				log.Println("Error while response")
				log.Println(response.Err)
			}
		}
	}

	// defer t1.Stop()

	// configs.Loger("started reader for ", board.Id)
	// port := board.Port

	// cmdReq := "ATS?"
	// if board.BoardType == rs485.MotorBoard {
	// 	cmdReq = "ATSPEED?"
	// }

	// for {
	// 	select {
	// 	case <-*board.terminate:
	// 		configs.Loger("terminate board", board.Id)
	// 		return
	// 	case message := <-*board.TriggerChan:
	// 		fmt.Println("trigger", message)
	// 		*board.LogerCh <- configs.Event{
	// 			BoardId:    board.Id,
	// 			EventsType: "Trigger",
	// 			Data:       message,
	// 			IsChange:   true,
	// 		}

	// 	case <-t1.C:
	// 		board.isBusy.Lock()
	// 		rsp2 := port.DoRequest(board.Id, cmdReq)
	// 		board.isBusy.Unlock()
	// 		if rsp2.Err == nil {
	// 			currData := string(rsp2.Response)
	// 			isNewData := board.checkData(currData)
	// 			ev := configs.Event{
	// 				BoardId:    board.Id,
	// 				EventsType: board.GetType(),
	// 				Name:       board.Name,
	// 				Refresh:    board.ReadInterval,
	// 				Data:       currData,
	// 				Status:     board.Status,
	// 				IsChange:   isNewData,
	// 				Chanels:    board.Chanels,
	// 			}
	// 			*board.LogerCh <- ev
	// 		}
	// 	}
	// }

}
