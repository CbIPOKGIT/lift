package display

import (
	"errors"
	"fmt"
	"github.com/CbIPOKGIT/lift/drivers/nanopi"
	"github.com/CbIPOKGIT/lift/internal/mainboard"
	"time"

	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

const numRows = 3
const menuStr = "prev     ok     next"
const pageSearch = 1
const pageClear = 2

type page [numRows + 1]string

var page1 = page{"Sys info", " ", " ", menuStr}
var page2 = page{"Search board", " ", " ", menuStr}
var page3 = page{"Stop all boards", " ", " ", menuStr}
var staticPages = []page{page1, page2, page3}

var numStatPages = len(staticPages)

type Display struct {
	mb          *mainboard.MainBoard
	CurrentPage uint8
	lcd         *Lcd
	offRender   *chan bool
	pages       []page
	rows        [numRows + 1]string
}

func New(mb *mainboard.MainBoard) *Display {
	if _, err := driverreg.Init(); err != nil {
		fmt.Println("driverreg", err)
	}

	b, err := i2creg.Open("0")
	if err != nil {
		fmt.Println("driverreg", err)
	}

	// Dev is a valid conn.Conn.
	dev := &i2c.Dev{Addr: 0x27, Bus: b}

	lcd, err := NewLcd(dev, LCD_20x4)
	if err != nil {
		fmt.Println("NewLcd", err)
	}
	err = lcd.BacklightOn()
	if err != nil {
		fmt.Println("BacklightOn", err)
	}
	fmt.Println("Init display")
	offRender := make(chan bool, 2)
	return &Display{
		mb:          mb,
		CurrentPage: 0,
		lcd:         lcd,
		offRender:   &offRender,
		pages:       staticPages,
	}
}

func (d *Display) AddPage(title string) {
	d.pages = append(d.pages, page{title, " ", " ", menuStr})
}

func (d *Display) RemovePage(pos uint8) {
	d.pages = append(d.pages[:pos], d.pages[pos+1:]...)
}

func (d *Display) Exec(page int8) string {
	switch page {
	case 0:
		return nanopi.GetMainInfo()
	//case 1:
	//	board, err := d.mb.SearcBrd()
	//	if err != nil {
	//		return "Error search"
	//	}
	//	title := board.BoardType.String() + ":" + strconv.Itoa(int(board.Id))
	//	fmt.Println("Added page", title)
	//	d.AddPage(title)
	//	return "Added: " + title
	//case 2:
	//	var num int
	//
	//	for _, v := range d.mb.GetActiveBoards() {
	//		d.mb.RemoveBoard(v.Id)
	//		d.RemovePage(uint8(numStatPages))
	//		num++
	//	}
	//	return "Removed: " + strconv.Itoa(int(num)) + " boards"
	default:
		//idS := strings.Split(d.pages[page][0], ":")
		//id, _ := strconv.Atoi(idS[1])
		//return "Data:" + d.mb.BoardData[id]
		return "foo"
	}
}

func (d *Display) RenderData() {
	maxLen := 20
	var cp int8 = int8(d.CurrentPage)

	fmt.Println("started render")
	timer := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-*d.offRender:
			fmt.Println("off render")
			return
		case <-timer.C:
			strData := d.Exec(cp)
			if len(strData) > maxLen {
				d.ShowData(strData[0:maxLen], 1)
				d.ShowData(strData[maxLen+1:], 2)
			} else {
				d.ShowData(strData, 1)
			}
			if cp == pageSearch || cp == pageClear {
				return
			}
		}
	}
}

func (d *Display) ChangePage(dir int8) {

	var cp int8 = int8(d.CurrentPage)
	numPages := int8(len(d.pages) - 1)
	if dir == 0 {
		go d.RenderData()
		return
	}
	if len(*d.offRender) == 0 {
		fmt.Println("terminate old render")
		*d.offRender <- true
	}

	cp += dir

	switch {
	case cp > numPages:
		cp = 0
	case cp < 0:
		cp = numPages
	}
	d.CurrentPage = uint8(cp)
	d.RenderPage(cp)
}

func (d *Display) RenderPage(cp int8) {
	for y := 0; y <= numRows; y++ {
		d.ShowData(d.pages[cp][y], uint8(y))
	}
}

func (d *Display) ShowData(str string, pos uint8) error {
	if pos > numRows {
		return nil
	}
	if str != d.rows[pos] {
		d.rows[pos] = str
		d.PrintMessage(pos, str)
		return nil
	}

	return errors.New("same data")
}

func (d *Display) PrintMessage(x uint8, message string) error {
	var pos ShowOptions = SHOW_LINE_1
	switch x {
	case 1:
		pos = SHOW_LINE_2
	case 2:
		pos = SHOW_LINE_3
	case 3:
		pos = SHOW_LINE_4
	}
	err := d.lcd.ShowMessage("                    ", pos)
	if err != nil {
		fmt.Println("ShowMessage empty", err)
		return err
	}
	time.Sleep(time.Millisecond)

	err = d.lcd.ShowMessage(message, pos)
	if err != nil {
		fmt.Println("ShowMessage", err)
		return err
	}

	return nil
}
