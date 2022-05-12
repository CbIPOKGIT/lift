package models

const (
	TABLE = "boards"
)

type Board struct {
	Id      uint8  `xorm:"pk autoincr"`
	CpuId   uint64 `xorm:"cpu_id not null"`
	ReadInt uint16 `xorm:"read_interval not null"`
	Status  uint8  `xorm:"status default '0'"`
	Name    string `xorm:"name default ''"`
	Current string `xorm:"current_data default ''"`
	Type    byte   `xorm:"board_type not null"`
}

type Boards []Board

func (Board) TableName() string {
	return TABLE
}

func (Boards) TableName() string {
	return TABLE
}

func (b *Boards) All() error {
	con, err := CreateConnection()
	if err != nil {
		return err
	}
	defer con.Close()

	return con.Find(b)
}
