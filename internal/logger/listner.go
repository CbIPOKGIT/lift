package logger

import "log"

func (l *Logger) Listen() {
	for {
		select {

		// Слухаємо статуси mainboard і
		case message := <-l.BoardsChannel:
			log.Println("Send data to server")
			log.Println(message)
			l.ServerSocket.Write(message)

			// l.CheckTriggers()

		case message := <-l.ServerSocket.Recive():
			l.MainboardCommandsChannel <- message
		}
	}
}
