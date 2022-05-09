package serial

import (
	"github.com/goburrow/serial"
	"io"
	"sync"
	"time"
)

var (
	DefaultReadTimeout int = 30 // in Milliseconds
	DefaultDataBits        = 8
	DefaultParity          = "N"
	DefaultStopBits        = 1
)

const (
	TrimSymbols = " \n\r\t"
)

type SerialPort struct {
	port io.ReadWriteCloser
	mtx  sync.Locker
}

func NewPort(port string, baud int, timeout int) (*SerialPort, error) {
	var s = new(SerialPort)
	var err error
	config := &serial.Config{
		Address:  port,
		BaudRate: baud,
		DataBits: DefaultDataBits,
		StopBits: DefaultStopBits,
		Parity:   DefaultParity,
		Timeout:  time.Duration(timeout) * time.Millisecond,
		RS485:    serial.RS485Config{},
	}

	s.port, err = serial.Open(config)

	s.mtx = new(sync.Mutex)

	return s, err
}

func (s *SerialPort) Lock() {
	s.mtx.Lock()
}

func (s *SerialPort) Unlock() {
	s.mtx.Unlock()
}

func (s *SerialPort) Close() error {
	return s.port.Close()
}

func (s *SerialPort) Write(b []byte) (int, error) {
	n, err := s.port.Write(b)
	return n, err
}

func (s *SerialPort) Read() ([]byte, error) {
	b := make([]byte, 1024)
	n, err := s.port.Read(b)
	c:=make([]byte,n)
	copy(c,b)
	return c, err
}
