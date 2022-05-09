package serial

import (
	"runtime"
	"testing"
	"time"
)

var p *SerialPort
var stop = make(chan bool)

func TestNewSerialPort(t *testing.T) {
	var err error
	p, err = NewPort("/dev/ttyS2", 115200, DefaultReadTimeout)
	if err != nil {
		t.Error(err)
	}
	runtime.GOMAXPROCS(1)
}

func TestSerialPort_Cmd(t *testing.T) {
	defer p.Close()

	// go func() {
	// 	t1 := time.NewTicker(time.Millisecond * 501)
	// 	defer t1.Stop()

	// 	for true {
	// 		select {
	// 		case <-t1.C:
	// 			res, err := p.Request("ATPING")
	// 			if err != nil {
	// 				t.Error(err)
	// 			}
	// 			log.Println("ATPING -", res.Duration(), res.Response())
	// 		case <-stop:
	// 			return
	// 		default:
	// 			break
	// 		}
	// 	}
	// }()

	// query := []string{"ATS?", "ATO?", "ATO=0", "ATRFID?", "ATVOLTAGE?", "ATDOORBUTTON?", "ATW?", "ATSYSTIME?", "ATALL?"} //, "ATDOORLEDS=", "ATDOORBLINK="}
	for true {
		<-time.After(time.Millisecond * 100)
		// for _, cmd := range query {
		/*
			blink := rand.Int31n(255)
			if cmd == "ATDOORBLINK=" {
				cmd += strconv.Itoa(int(blink))
			}
			ledOn := rand.Int31n(255)
			if cmd == "ATDOORLEDS=" {
				cmd += strconv.Itoa(int(ledOn))
			}
		*/

		// res, err := p.Request(cmd)
		// if err != nil {
		// 	log.Println(cmd)
		// 	t.Fatal(err)
		// }
		// log.Println(cmd, res.Duration(), res.Response())

		// }
	}

	stop <- true
}
