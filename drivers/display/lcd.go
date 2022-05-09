package display

import (
	"fmt"
	"strings"
	"time"

	"periph.io/x/conn/v3/i2c"
	// "github.com/d2r2/go-i2c"
)

const (
	// Commands
	CMD_Clear_Display        = 0x01
	CMD_Return_Home          = 0x02
	CMD_Entry_Mode           = 0x04
	CMD_Display_Control      = 0x08
	CMD_Cursor_Display_Shift = 0x10
	CMD_Function_Set         = 0x20
	CMD_CGRAM_Set            = 0x40
	CMD_DDRAM_Set            = 0x80

	// Options
	OPT_Increment = 0x02 // CMD_Entry_Mode
	OPT_Decrement = 0x00
	// OPT_Display_Shift  = 0x01 // CMD_Entry_Mode
	OPT_Enable_Display = 0x04 // CMD_Display_Control
	OPT_Enable_Cursor  = 0x02 // CMD_Display_Control
	OPT_Enable_Blink   = 0x01 // CMD_Display_Control
	OPT_Display_Shift  = 0x08 // CMD_Cursor_Display_Shift
	OPT_Shift_Right    = 0x04 // CMD_Cursor_Display_Shift 0 = Left
	OPT_8Bit_Mode      = 0x10
	OPT_4Bit_Mode      = 0x00
	OPT_2_Lines        = 0x08 // CMD_Function_Set 0 = 1 line
	OPT_1_Lines        = 0x00
	OPT_5x10_Dots      = 0x04 // CMD_Function_Set 0 = 5x7 dots
	OPT_5x8_Dots       = 0x00
)

const (
	PIN_BACKLIGHT byte = 0x08
	PIN_EN        byte = 0x04 // Enable bit
	PIN_RW        byte = 0x02 // Read/Write bit
	PIN_RS        byte = 0x01 // Register select bit
)

type LcdType int

const (
	LCD_UNKNOWN LcdType = iota
	LCD_16x2
	LCD_20x4
)

type ShowOptions int

const (
	SHOW_NO_OPTIONS ShowOptions = 0
	SHOW_LINE_1                 = 1 << iota
	SHOW_LINE_2
	SHOW_LINE_3
	SHOW_LINE_4
	SHOW_ELIPSE_IF_NOT_FIT
	SHOW_BLANK_PADDING
)

type Lcd struct {
	i2c       *i2c.Dev
	backlight bool
	lcdType   LcdType
}

func NewLcd(i2c *i2c.Dev, lcdType LcdType) (*Lcd, error) {
	lcd := &Lcd{i2c: i2c, backlight: false, lcdType: lcdType}
	initByteSeq := []byte{
		0x03, 0x03, 0x03, // base initialization
		0x02, // setting up 4-bit transfer mode
		CMD_Function_Set | OPT_2_Lines | OPT_5x8_Dots | OPT_4Bit_Mode,
		CMD_Display_Control | OPT_Enable_Display,
		CMD_Entry_Mode | OPT_Increment,
	}
	for _, b := range initByteSeq {
		err := lcd.writeByte(b, 0)
		if err != nil {
			return nil, err
		}
	}
	err := lcd.Clear()
	if err != nil {
		return nil, err
	}
	err = lcd.Home()
	if err != nil {
		return nil, err
	}
	return lcd, nil
}

type rawData struct {
	Data  byte
	Delay time.Duration
}

func (lcd *Lcd) writeRawDataSeq(seq []rawData) error {
	for _, item := range seq {
		_, err := lcd.i2c.Write([]byte{item.Data})
		if err != nil {
			return err
		}
		time.Sleep(item.Delay)
	}
	return nil
}

func (lcd *Lcd) writeDataWithStrobe(data byte) error {
	if lcd.backlight {
		data |= PIN_BACKLIGHT
	}
	seq := []rawData{
		{data, 0},                               // send data
		{data | PIN_EN, 200 * time.Microsecond}, // set strobe
		{data, 30 * time.Microsecond},           // reset strobe
	}
	return lcd.writeRawDataSeq(seq)
}

func (lcd *Lcd) writeByte(data byte, controlPins byte) error {
	err := lcd.writeDataWithStrobe(data&0xF0 | controlPins)
	if err != nil {
		return err
	}
	err = lcd.writeDataWithStrobe((data<<4)&0xF0 | controlPins)
	if err != nil {
		return err
	}
	return nil
}

func (lcd *Lcd) getLineRange(options ShowOptions) (startLine, endLine int) {
	var lines [4]bool
	lines[0] = options&SHOW_LINE_1 != 0
	lines[1] = options&SHOW_LINE_2 != 0
	lines[2] = options&SHOW_LINE_3 != 0
	lines[3] = options&SHOW_LINE_4 != 0
	startLine = -1
	for i := 0; i < len(lines); i++ {
		if lines[i] {
			startLine = i
			break
		}
	}
	endLine = -1
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] {
			endLine = i
			break
		}
	}
	return startLine, endLine
}

func (lcd *Lcd) splitText(text string, options ShowOptions) []string {
	var lines []string
	startLine, endLine := lcd.getLineRange(options)
	w, _ := lcd.getSize()
	if w != -1 && startLine != -1 && endLine != -1 {
		for i := 0; i <= endLine-startLine; i++ {
			if len(text) == 0 {
				break
			}
			j := w
			if j > len(text) {
				j = len(text)
			}
			lines = append(lines, text[:j])
			text = text[j:]
		}
		if len(text) > 0 {
			if options&SHOW_ELIPSE_IF_NOT_FIT != 0 {
				j := len(lines) - 1
				lines[j] = lines[j][:len(lines[j])-1] + "~"
			}
		} else {
			if options&SHOW_BLANK_PADDING != 0 {
				j := len(lines) - 1
				lines[j] = lines[j] + strings.Repeat(" ", w-len(lines[j]))
				for k := j + 1; k <= endLine-startLine; k++ {
					lines = append(lines, strings.Repeat(" ", w))
				}
			}

		}
	} else if len(text) > 0 {
		lines = append(lines, text)
	}
	return lines
}

func (lcd *Lcd) ShowMessage(text string, options ShowOptions) error {
	lines := lcd.splitText(text, options)
	startLine, endLine := lcd.getLineRange(options)
	i := 0
	for {
		if startLine != -1 && endLine != -1 {
			err := lcd.SetPosition(i+startLine, 0)
			if err != nil {
				return err
			}
		}
		line := lines[i]
		for _, c := range line {
			err := lcd.writeByte(byte(c), PIN_RS)
			if err != nil {
				return err
			}
		}
		if i == len(lines)-1 {
			break
		}
		i++
	}
	return nil
}

func (lcd *Lcd) TestWriteCGRam() error {
	err := lcd.writeByte(CMD_CGRAM_Set, 0)
	if err != nil {
		return err
	}
	var a byte = 0x55
	for i := 0; i < 80; i++ {
		err := lcd.writeByte(a, PIN_RS)
		if err != nil {
			return err
		}
		a = a ^ 0xFF
	}
	return nil
}

func (lcd *Lcd) BacklightOn() error {
	lcd.backlight = true
	err := lcd.writeByte(0x00, 0)
	if err != nil {
		return err
	}
	return nil
}

func (lcd *Lcd) BacklightOff() error {
	lcd.backlight = false
	err := lcd.writeByte(0x00, 0)
	if err != nil {
		return err
	}
	return nil
}

func (lcd *Lcd) Clear() error {
	err := lcd.writeByte(CMD_Clear_Display, 0)
	return err
}

func (lcd *Lcd) Home() error {
	err := lcd.writeByte(CMD_Return_Home, 0)
	time.Sleep(3 * time.Millisecond)
	return err
}

func (lcd *Lcd) getSize() (width, height int) {
	switch lcd.lcdType {
	case LCD_16x2:
		return 16, 2
	case LCD_20x4:
		return 20, 4
	default:
		return -1, -1
	}
}

func (lcd *Lcd) SetPosition(line, pos int) error {
	w, h := lcd.getSize()
	if w != -1 && (pos < 0 || pos > w-1) {
		return fmt.Errorf("cursor position %d "+
			"must be within the range [0..%d]", pos, w-1)
	}
	if h != -1 && (line < 0 || line > h-1) {
		return fmt.Errorf("cursor line %d "+
			"must be within the range [0..%d]", line, h-1)
	}
	lineOffset := []byte{0x00, 0x40, 0x14, 0x54}
	var b byte = CMD_DDRAM_Set + lineOffset[line] + byte(pos)
	err := lcd.writeByte(b, 0)
	return err
}

func (lcd *Lcd) Write(buf []byte) (int, error) {
	for i, c := range buf {
		err := lcd.writeByte(c, PIN_RS)
		if err != nil {
			return i, err
		}
	}
	return len(buf), nil
}

func (lcd *Lcd) Command(cmd byte) error {
	err := lcd.writeByte(cmd, 0)
	return err
}
