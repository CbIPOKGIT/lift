package rs232

import (
	"fmt"
	"testing"
	"time"
)

var p *RS232
var stop = make(chan bool)

func TestSerialPort(t *testing.T) {
	var err error
	p, err = NewPort("/dev/ttyS2", 115200, 30, 0x0D)
	if err != nil {
		t.Error(err)
	}
}

func TestSerialPort_Cmd(t *testing.T) {
	defer p.Close()

	currStr := ""
	b, _ := p.DoRequest("ATLCDPRINT=0,1,\"" + currStr + "\"")
	fmt.Println(string(b))

	// time.Sleep(time.Millisecond * 500)
	// if true {
	// 	return
	// }

	t1 := time.NewTicker(time.Millisecond * 500)
	defer t1.Stop()

	for {
		select {
		case <-t1.C:
			ts := time.Now()
			b, err := p.DoRequest("ATDOORBUTTON?")
			strData := string(b)
			fmt.Println("ddd|", strData, "|")
			td := time.Since(ts)
			if err != nil || strData == "ATERROR  " {
				t.Error("Error", err)
				continue
			}
			fmt.Println(currStr, "\t--- \t", td)

			if currStr != strData {
				p.DoRequest("ATLCDCLEAR")
				currStr = strData

				ts2 := time.Now()
				b2, err2 := p.DoRequest("ATLCDPRINT=0,0,\"" + currStr + "\"")
				td2 := time.Since(ts2)
				if err2 != nil {
					t.Error(err)
					continue
				}
				fmt.Println(string(b2), "\t--- \t", td2)

			}

		case <-stop:
			return
		}
	}
}
