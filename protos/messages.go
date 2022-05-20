package protos

// Повідомлення з плат
type BoardMessage struct {

	// Тип борда {rs485 BoardTypes_t}
	BoardType int

	// Дані від датчика
	Message []byte
}

// Повідомлення с материнської плати
type MainboardMessage struct {
	// Датчики - 1, реле/входи - 2, вольтаж - 3
	Type int

	// Саме повідомлення
	Message string
}

// Канал отримання команд для материнської плати
type MainboardCommandChannel chan []byte
