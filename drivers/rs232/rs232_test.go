package rs232

import (
	"fmt"
	"testing"
	"time"
)

var port *RS232

func _TestNewPort(t *testing.T) {
	var err error
	port, err = NewPort("/dev/ttyS2", 115200, 30, 0x0D)
	if err != nil {
		t.Error(err)
	}
}

func _TestRS232_DoRequest(t *testing.T) {
	var b []byte
	var err error

	time.Sleep(time.Second * 10)
	ts := time.Now()
	b, err = port.DoRequest("ATS?")
	td := time.Since(ts)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(b), "\t--- \t", td)

}
