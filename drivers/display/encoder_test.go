package display

import (
	"fmt"
	"log"
	"testing"
	"time"

	// "github.com/d2r2/go-i2c"
	"github.com/mclaut/ec11"

	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

func TestEncoder(t *testing.T) {
	fmt.Println("Started")
	e, err := ec11.New(15, 363, 14)
	if err != nil {
		t.Error(err)
	}
	ch := e.Start()
	for i := range ch {
		fmt.Println(i)
	}
}

func TestDisplay(t *testing.T) {
	if _, err := driverreg.Init(); err != nil {
		log.Fatal(err)
	}

	b, err := i2creg.Open("0")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	// Dev is a valid conn.Conn.
	dev := &i2c.Dev{Addr: 0x27, Bus: b}

	if err != nil {
		log.Fatal(err)
	}
	// defer i2c.Close()
	lcd, err := NewLcd(dev, LCD_20x4)
	if err != nil {
		log.Fatal(err)
	}
	err = lcd.BacklightOn()
	if err != nil {
		log.Fatal(err)
	}
	err = lcd.ShowMessage("--=! Let's rock !=--", SHOW_LINE_1)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second)
	err = lcd.ShowMessage("Welcome to RPi dude!", SHOW_LINE_2)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second)
	err = lcd.ShowMessage("Welcome to RPi dude!", SHOW_LINE_3)
	if err != nil {
		log.Fatal(err)
	}

	// time.Sleep(5 * time.Second)
	// err = lcd.BacklightOff()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
