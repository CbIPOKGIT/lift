package mainboard

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/CbIPOKGIT/lift/configs"
	"github.com/CbIPOKGIT/lift/drivers/rs485"
	"github.com/CbIPOKGIT/lift/internal/board"
)

func (mb *MainBoard) Disconnect() {
	mb.P232.Close()
	mb.P485.Close()
}

// Підключення бордів
func (mb *MainBoard) LoadBoards() {

	mb.Boards = make([]*board.Board, 0, 255)

	if boards, err := mb.LoadBoardsFromStorage(); err == nil && len(*boards) > 0 {
		mb.Boards = *boards
	} else {
		for {
			if _, err := mb.SearcBoard(); err == nil {
				time.Sleep(time.Second)
			} else {
				mb.Storage.Set("boards", mb.Boards)
				return
			}
		}
	}
}

// Завантажуємо данні про boards зі сховища
func (mb *MainBoard) LoadBoardsFromStorage() (*Boards, error) {
	data, _ := mb.Storage.Get("boards")

	boards := make(Boards, 0)

	if err := json.Unmarshal([]byte(data), &boards); err == nil {
		for i := 0; i < len(boards); i++ {
			boards[i].Port = mb.P485
		}
		return &boards, nil
	} else {
		return nil, err
	}
}

// Виконуємо команду борда
// Поки що on/off
func (mb *MainBoard) GetData(command string) (string, error) {
	resp, err := mb.P232.DoRequest(configs.TranslateCommand(command))
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// Пошук плат
func (mb *MainBoard) SearcBoard() (*board.Board, error) {
	boardID := len(mb.Boards) + 1
	log.Println("Searching board")

	// mb.IsBusy.Lock()
	cpuId, errSearch := searchBoard(mb.P485)
	if errSearch != nil {
		log.Println("Error")
		log.Println(errSearch)
		return nil, errSearch
	}

	bdata := board.BoardData{Id: uint8(boardID), CpuId: cpuId, ReadInterval: configs.BOARD_READ_INTERVAL}
	if newBoard, err := board.New(mb.P485, bdata); err == nil {
		if errAddBoard := mb.AddBoard(newBoard); errAddBoard == nil {
			return newBoard, nil
		} else {
			return nil, errAddBoard
		}

	} else {
		return nil, err
	}
	// go curBoard.ReadData()
}

func (mb *MainBoard) AddBoard(board *board.Board) error {
	for _, b := range mb.Boards {
		if board.CpuId == b.CpuId {
			return errors.New("Board already connected")
		}
	}
	mb.Boards = append(mb.Boards, board)
	return nil
}

func searchBoard(port *rs485.RS485) (uint64, error) {

	rsp := port.Search()
	if rsp.Err != nil {
		return 0, rsp.Err
	}

	return rsp.Cpuid, nil
}
