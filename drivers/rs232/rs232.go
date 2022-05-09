package rs232

import (
	"time"

	"github.com/CbIPOKGIT/lift/drivers/serial"
)

type RS232 struct {
	port     Porter
	stopByte byte
	timeout  time.Duration
}

func NewPort(portName string, baud int, timeout int, stopByte byte) (*RS232, error) {
	var err error
	r := new(RS232)
	r.port, err = serial.NewPort(portName, baud, timeout)
	r.stopByte = stopByte
	r.timeout = time.Duration(timeout) * time.Millisecond
	return r, err
}

func (r *RS232) Close() error {
	return r.port.Close()
}

func (r *RS232) SetStopByte(b byte) {
	r.stopByte = b
}

func (r *RS232) DoRequest(command string) ([]byte, error) {
	var err error
	// b := make([]byte, 0, 512)
	res := make([]byte, 0, 512)

	tmr := time.NewTicker(r.timeout)
	defer tmr.Stop()

	cmdB := createCommand(command)

	r.port.Lock()
	defer r.port.Unlock()

	_, err = r.port.Write(cmdB)
	if err != nil {
		return nil, err
	}
Loop:
	for {
		select {
		case <-tmr.C:
			break Loop
		default:
			b, err := r.port.Read()
			if err != nil {
				return nil, err
			}

			for _, bt := range b {
				if bt == r.stopByte {
					break Loop
				}
				res = append(res, bt)
			}
		}
	}
	return res, nil
}

func createCommand(s string) []byte {
	// b := make([]byte, 0, len(s)+1)
	b := []byte(s)
	b = append(b, 0x0A)
	return b
}
