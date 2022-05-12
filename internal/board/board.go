package board

import (
	"log"
	"time"

	"github.com/CbIPOKGIT/lift/drivers/rs485"
	"github.com/CbIPOKGIT/lift/protos"
)

// Зчитуємо дані з плати
func (board *Board) ReadData(messenger protos.BoardSpeaker) {
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
			board.Lock()

			response := board.Port.DoRequest(board.Id, command)

			board.Unlock()

			if response.Err != nil {
				log.Println("Error while response")
				log.Println(response.Err)
				continue
			}

			if board.CurrentData != string(response.Response) {
				message := protos.BoardMessage{
					Message:   response.Response,
					BoardType: int(board.BoardType),
				}

				messenger.ReciveFromBoard(&message)

				board.CurrentData = string(response.Response)
			}

		}
	}
}
