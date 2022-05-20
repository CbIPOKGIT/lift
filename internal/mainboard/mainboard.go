package mainboard

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/CbIPOKGIT/lift/configs"
	"github.com/CbIPOKGIT/lift/drivers/rs485"
	"github.com/CbIPOKGIT/lift/internal/board"
	"github.com/CbIPOKGIT/lift/models"
)

func (mb *MainBoard) Disconnect() {
	mb.P232.Close()
	mb.P485.Close()
}

// Підключення бордів
func (mb *MainBoard) LoadBoards() {

	mb.Boards = make([]*board.Board, 0, 255)

	if boards, err := mb.LoadBoardsFromDB(); err == nil && len(*boards) > 0 {
		mb.Boards = *boards
		for i := len(mb.Boards) - 1; i >= 0; i-- {
			board := mb.Boards[i]
			board.Test(mb.P485)
		}
		os.Exit(1)

	} else {
		for {
			index := 1
			if board, err := mb.SearcBoard(uint8(index)); err == nil {
				mb.Boards = append(mb.Boards, board)
				time.Sleep(time.Second)
			} else {
				mb.SaveBoardsToDB()
				return
			}
		}
	}
}

// Завантажуємо данні про boards зі сховища
func (mb *MainBoard) LoadBoardsFromDB() (*Boards, error) {
	dbBoards := make(models.Boards, 0)
	if err := dbBoards.All(); err != nil {
		return nil, err
	}

	boards := make(Boards, 0, len(dbBoards))

	for _, db := range dbBoards {
		board := new(board.Board)
		board.Id = db.Id
		board.CpuId = db.CpuId
		board.ReadInterval = db.ReadInt
		board.Status = db.Status
		board.Name = db.Name
		board.CurrentData = db.Current
		board.BoardType = rs485.BoardTypes_t(db.Type)
		board.Port = mb.P485

		boards = append(boards, board)
	}

	return &boards, nil
}

// Зберігаємо дані про плати в базу
func (mb *MainBoard) SaveBoardsToDB() {
	// dbBoards := make(models.Boards, 0, len(mb.Boards))

	// for _, board := range v {
	// 	var dbBoard models.Board

	// 	dbBoard.Id = dbBoard.Id
	// 	dbBoard.CpuId = dbBoard.CpuId
	// 	dbBoard.Id = dbBoard.Id
	// 	dbBoard.Id = dbBoard.Id
	// 	dbBoard.Id = dbBoard.Id
	// 	dbBoard.Id = dbBoard.Id
	// 	dbBoard.Id = dbBoard.Id

	// 	dbBoards = append(dbBoards, dbBoard)
	// }
}

// Виконуємо команду борда
// Поки що on/off
func (mb *MainBoard) GetData(command string) (string, error) {
	resp, err := mb.P232.DoRequest(configs.TranslateCommand(command))
	if err != nil {
		log.Println("Error request to p232", "Command - ", command)
		return "", err
	}
	return string(resp), nil
}

// Пошук плат
func (mb *MainBoard) SearcBoard(boardID uint8) (*board.Board, error) {
	log.Println("Searching board")

	mb.Lock()
	defer mb.Unlock()

	cpuId, errSearch := mb.searchBoard(mb.P485)
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

func (mb *MainBoard) searchBoard(port *rs485.RS485) (uint64, error) {

	rsp := port.Search()
	if rsp.Err != nil {
		return 0, rsp.Err
	}

	return rsp.Cpuid, nil
}
