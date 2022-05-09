package rs485

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

var port *RS485

func TestNewPort(t *testing.T) {
	var err error
	port, err = NewPort("/dev/ttyUSB0", 115200, 60)
	if err != nil {
		t.Error(err)
	}
	port.SetStartByte(0xFA)
	port.SetStopByte(0xFE)

}

func TestRS485_Search(t *testing.T) {
	fmt.Println("Search")
	rsp := port.Search()
	fmt.Println(rsp)
	if rsp.Err != nil {
		t.Error(rsp.Err)
	}
}

func TestRS485_SetAddr(t *testing.T) {
	var addr uint8 = 139
	fmt.Println("Set addr", addr)
	rsp := port.Search()
	fmt.Println(rsp)
	if rsp.Err != nil {
		t.Error(rsp.Err)
	}

	err := port.SetAddr(addr, rsp.Cpuid)
	if err != nil {
		t.Error(err)
	}
}

func TestRS485_GetBoardType(t *testing.T) {
	b, err := port.GetBoardType(139)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%x\n", b)
	fmt.Println(b)
}

func TestRS485_Request(t *testing.T) {
	var i uint64
	defer port.Close()

	ch := make(chan bool, 4)
	wg := new(sync.WaitGroup)
	wg.Add(4)

	go func() {
		defer wg.Done()
		t := time.NewTicker(time.Second * 10)
		defer t.Stop()
		for range t.C {
			fmt.Println(strings.Repeat("STOP ", 20))
			ch <- true
			ch <- true
			ch <- true
			//ch <- true
			return
		}

	}()

	go func() {
		defer wg.Done()
		var maxDuration time.Duration
		t := time.NewTicker(time.Millisecond * 25)
		defer t.Stop()
		for {
			select {
			case b, ok := <-ch:
				if ok {
					ch <- true
				}
				if b {
					fmt.Println("goroutine 1 - stop")
					return
				}

			case <-t.C:
				rsp := port.DoRequest(129, "ATSPEED?")
				if rsp.Err != nil {
					fmt.Println(strings.Repeat("-", 60))
				}

				if maxDuration < rsp.Duration {
					maxDuration = rsp.Duration
				}
				fmt.Printf("%d: %s \t | \t %v \t %v \t %v\n", i, rsp.Response, rsp.Err, rsp.Duration, maxDuration)
				i++
			default:
				<-time.After(time.Millisecond * 5)
			}
		}
	}()

	go func() {
		defer wg.Done()
		var maxDuration time.Duration
		t := time.NewTicker(time.Millisecond * 101)
		defer t.Stop()
		for {
			select {
			case b, ok := <-ch:
				if ok {
					ch <- true
				}
				if b {
					fmt.Println("goroutine 2 - stop")
					return
				}
			case <-t.C:
				rsp := port.DoRequest(140, "ATS?")
				if rsp.Err != nil {
					fmt.Println(strings.Repeat("-", 60), rsp.Response)
				}
				if maxDuration < rsp.Duration {
					maxDuration = rsp.Duration
				}
				fmt.Printf("%d: %s \t | \t %v \t %v \t %v\n", i, rsp.Response, rsp.Err, rsp.Duration, maxDuration)
				i++
			default:
				<-time.After(time.Millisecond * 5)
			}
		}
	}()

	go func() {
		defer wg.Done()
		var maxDuration time.Duration
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			select {
			case b, ok := <-ch:
				if ok {
					ch <- true
				}
				if b {
					fmt.Println("goroutine 3 - stop")
					return
				}
			case <-t.C:
				rsp := port.DoRequest(144, "ATS?")
				if rsp.Err != nil {
					fmt.Printf("%d: %s \t | \t %v \t %v \t %v\n", i, rsp.Response, rsp.Err, rsp.Duration, maxDuration)
				}
				if maxDuration < rsp.Duration {
					maxDuration = rsp.Duration
				}
				fmt.Printf("%d: %s \t | \t %v \t %v \t %v\n", i, rsp.Response, rsp.Err, rsp.Duration, maxDuration)
				i++
			default:
				<-time.After(time.Millisecond * 5)
			}

		}
	}()

	wg.Wait()
}
