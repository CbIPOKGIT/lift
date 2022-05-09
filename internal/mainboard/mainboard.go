package mainboard

import "github.com/CbIPOKGIT/lift/configs"

func (mb *MainBoard) Disconnect() {
	mb.P232.Close()
	mb.P485.Close()
}

// Підключення бордів
func (mb *MainBoard) LoadBoards() error {
	return nil
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
