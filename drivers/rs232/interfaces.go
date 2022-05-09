package rs232

type Porter interface {
	Lock()
	Unlock()
	Close() error
	Write([]byte) (int, error)
	Read() ([]byte, error)
}
