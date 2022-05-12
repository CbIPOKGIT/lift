package database

import (
	"log"

	"github.com/CbIPOKGIT/lift/models"
	"xorm.io/xorm"
)

func InitDB() error {
	con, err := models.CreateConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	return initTables(con)
}

func initTables(con *xorm.Engine) error {
	board := new(models.Board)
	if err := con.Sync2(board); err != nil {
		return err
	}
	// if count, err := con.Count(board); err == nil && count == 0 {
	// 	stor := storage.New()
	// 	data, _ := stor.Get("boards")
	// 	mbBoards := make(mainboard.Boards, 0)

	// 	json.Unmarshal([]byte(data), &mbBoards)

	// 	boards := make(models.Boards, 0, len(mbBoards))

	// 	for _, bData := range mbBoards {
	// 		var board models.Board

	// 		board.Id = bData.Id
	// 		board.CpuId = bData.CpuId
	// 		board.ReadInt = bData.ReadInterval
	// 		board.Status = bData.Status
	// 		board.Name = bData.Name
	// 		board.Current = bData.CurrentData
	// 		board.Type = byte(bData.BoardType)

	// 		boards = append(boards, board)
	// 	}

	// 	if _, err := con.Insert(&boards); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	return nil
}
