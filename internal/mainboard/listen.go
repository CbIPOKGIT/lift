package mainboard

import (
	"time"

	"github.com/CbIPOKGIT/lift/protos"
)

func (mb *MainBoard) Listen(bus protos.GlobalBus) {

	// Зчитуємо дані з приєднаних плат
	for _, board := range mb.Boards {
		go board.ReadData(bus)
	}

	// Зчитуємо статус материнської плати
	go mb.ListenMBStatus(bus)

	// Виконуємо команди з сервера
	go mb.ListenServerCommands(bus.GetCommandsBus())
}

// Зчитуємо показники з материнської плати
func (mb *MainBoard) ListenMBStatus(messenger protos.MBSpeaker) {
	ticker := time.NewTicker(time.Duration(mb.ReadInterval) * time.Millisecond)

	for {
		select {

		//
		case <-ticker.C:
			status, err := mb.GetData("status_sensors")

			if err != nil {
				// Логируем ошибку
				continue
			}

			if status != mb.StatusSensors {
				mb.StatusSensors = status
				messenger.ReciveFromMainboard(&protos.MainboardMessage{
					Sensor:  true,
					Message: status,
				})
			}
		}
	}
}

// Очікуємо та виконуємо команди з сервера
func (mb *MainBoard) ListenServerCommands(bus protos.MainboardCommandChannel) {
	for {
		select {
		case command := <-bus:
			mb.GetData(string(command))
		}
	}
}
