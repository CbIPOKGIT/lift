package rs485

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/CbIPOKGIT/lift/drivers/crc"
	"github.com/CbIPOKGIT/lift/drivers/serial"
)

const (
	DefaultSearchHash       string = "$*%&#^@!"
	DefaultSearchAddr       uint8  = 255
	DefaultRequestStartByte byte   = 0xFB
	DefaultRequestStopByte  byte   = 0xFE
)

// Boards codes
type BoardTypes_t byte

type boardMaps_t map[BoardTypes_t]string

const (
	OptoBoard   BoardTypes_t = 0x10
	SensorBoard BoardTypes_t = 0x20
	MotorBoard  BoardTypes_t = 0x30
	DoorBoard   BoardTypes_t = 0x00
	Unknown     BoardTypes_t = 0x01
)

var boardsMap = boardMaps_t{
	OptoBoard:   "OptoBoard",
	SensorBoard: "SensorBoard",
	MotorBoard:  "MotorBoard",
	DoorBoard:   "DoorBoard",
}

func (b BoardTypes_t) String() string {
	if s, ok := boardsMap[b]; ok {
		return s
	}

	return "not found"
}

var (
	DefaultReadTimeout          int = 60  // in milliseconds
	DefaultSearchInterval       int = 600 // in milliseconds
	DefaultSearchProcessTimeout int = 60  // in seconds
)

// var errors
var (
	ErrTimeout            = errors.New("read timeout")
	ErrEmptyArgument      = errors.New("empty argument not allowed")
	ErrReceivedCRC        = errors.New("response have wrong crc")
	ErrEmptyAnswer        = errors.New("empty Answer")
	ErrSearchHexDecode    = errors.New("decode byte to hex error")
	ErrSearchWrongAnswer1 = errors.New("wrong search's Answer(=)")
	ErrSearchWrongAnswer2 = errors.New("wrong search's Answer(,)")
	ErrSearchParseCPUid   = errors.New("wrong CPU id")
)

type RS485 struct {
	port      Porter
	stopByte  byte
	startByte byte
}

type SearchResult struct {
	Addr  byte
	Cpuid uint64
	Hash  bool
	Err   error
}

func NewPort(portName string, baud int, timeout int) (*RS485, error) {
	// test comment
	var err error
	r := new(RS485)
	r.port, err = serial.NewPort(portName, baud, timeout)
	return r, err
}

func (r *RS485) Close() error {
	return r.port.Close()
}

func (r *RS485) SetStartByte(b byte) {
	r.startByte = b
}

func (r *RS485) SetStopByte(b byte) {
	r.stopByte = b
}

func (r *RS485) DoRequest(addr byte, mess string) Response {
	ts := time.Now()
	rsp := Response{}
	var err error

	// create Request
	request := r.createRequest(addr, mess)

	r.port.Lock()
	defer r.port.Unlock()

	// make Request
	err = r.write(request)
	if err != nil {
		rsp.Err = err
		rsp.Duration = time.Since(ts)
		return rsp
	}

	// read response
	rawAnswer, err := r.read()
	if err != nil {
		rsp.Duration = time.Since(ts)
		rsp.Err = err
		return rsp
	}

	rawAnswer, err = stripStopStartSymbols(rawAnswer)
	if err == ErrEmptyArgument {
		rsp.Err = ErrEmptyAnswer
		return rsp
	}

	// convert message
	//fmt.Printf("%x\n", rawAnswer)
	btAnswer, err := hex.DecodeString(string(rawAnswer[1:]))
	//fmt.Printf(">>>%s - %v:%v \n", btAnswer, err, rsp.Duration)
	if err != nil {
		rsp.Err = err
		rsp.Duration = time.Since(ts)
		return rsp
	}

	// check for crc
	tmpAnswer := crc.CRC16X25(btAnswer[:len(btAnswer)-2])
	//fmt.Printf("%x = %x\n", btAnswer[len(btAnswer)-2:], tmpAnswer)
	if !bytes.Equal(btAnswer[len(btAnswer)-2:], tmpAnswer) {
		rsp.Err = ErrReceivedCRC
		rsp.Duration = time.Since(ts)
		return rsp
	}

	rsp.Response = btAnswer[:len(btAnswer)-2]
	rsp.Duration = time.Since(ts)
	return rsp
}

func (r *RS485) Search() SearchResult {
	rsp := SearchResult{}
	rawRequest := r.createRequest(DefaultSearchAddr, "ATSEARCH")

	tmrStop := time.NewTimer(time.Second * time.Duration(DefaultSearchProcessTimeout))
	tmr := time.NewTicker(time.Millisecond * time.Duration(DefaultSearchInterval))
	defer tmr.Stop()
	defer tmrStop.Stop()
Loop:
	for {
		select {
		case <-tmr.C:
			// send broadcast message
			rsp.Err = r.write(rawRequest)
			if rsp.Err != nil {
				return rsp
			}

			// listen Answer if device button is pressed
			var rawAnswer []byte
			rawAnswer, rsp.Err = r.read()
			if rsp.Err != nil {
				if rsp.Err.Error() == "serial: timeout" {
					continue
				}
				return rsp
			}

			if len(rawAnswer) > 0 {
				rsp = r.parseSearch(rawAnswer)
				if rsp.Err != nil {
					return rsp
				}

				return rsp
			}

		case <-tmrStop.C:
			fmt.Println("STOP")
			break Loop

		default:
			time.Sleep(time.Millisecond * 5)
		}

	}

	return rsp
}

func (r *RS485) GetBoardType(addr byte) (BoardTypes_t, error) {
	response := r.DoRequest(addr, "ATCPUID?")
	log.Println("Response")
	log.Println(response.Response)
	log.Println(response.Err)
	if response.Err != nil {
		return 0, response.Err
	}
	b := string(response.Response[len(response.Response)-2:])
	log.Println("Board type from response")
	log.Println(b)
	d, err := strconv.Atoi(b)
	if err != nil {
		return 0, err
	}

	return BoardTypes_t(d), nil
}

func (r *RS485) write(b []byte) error {
	_, err := r.port.Write(b)
	return err
}

// function drop start byte and finest byte
func (r *RS485) read() ([]byte, error) {
	res := make([]byte, 0, 1024)
	var b []byte
	var err error
Loop:
	for {
		select {
		case <-time.After(time.Millisecond * time.Duration(DefaultReadTimeout)):
			return nil, ErrTimeout
		default:
			b, err = r.port.Read()
			if err != nil {
				return nil, err
			}
			res = append(res, b...)
			if res[0] == r.startByte && res[len(res)-1] == r.stopByte {
				break Loop
			}
		}
	}

	return res, nil
}

func stripStopStartSymbols(b []byte) ([]byte, error) {
	if b == nil {
		return nil, ErrEmptyArgument
	}

	return b[1 : len(b)-1], nil
}

func (r *RS485) createRequest(addr byte, message string) []byte {
	req := make([]byte, 0, (len(message)+2)*2+3)
	req = append(req, DefaultRequestStartByte, addr)
	req = append(req, hex.EncodeToString(append([]byte(message), crc.CRC16X25([]byte(message))...))...)
	req = append(req, DefaultRequestStopByte)
	return req
}

func (r *RS485) parseSearch(b []byte) SearchResult {
	var err error
	res := SearchResult{
		Addr:  0,
		Cpuid: 0,
		Hash:  false,
		Err:   nil,
	}

	b, err = stripStopStartSymbols(b)
	if err != nil {
		res.Err = ErrEmptyAnswer
	}

	res.Addr = b[0]

	b = stripCRC(stripAddr(b))

	b, err = hex.DecodeString(string(b))
	if err != nil {
		res.Err = err
		return res
	}
	s := string(b)

	tmp0 := strings.Split(s, "=")
	if len(tmp0) != 2 {
		res.Err = ErrSearchWrongAnswer1
		return res
	}

	tmp1 := strings.Split(tmp0[1], ",")
	if len(tmp1) != 2 {
		res.Err = ErrSearchWrongAnswer2
		return res
	}

	res.Cpuid, err = strconv.ParseUint(tmp1[0], 10, 64)
	if err != nil {
		res.Err = ErrSearchParseCPUid
		return res
	}

	// get Hash from received data
	flag := true
	hashR := make([]byte, 0, 32)
	for _, j := range []byte(tmp1[1]) {
		if j == '"' {
			if flag {
				flag = !flag
				continue
			} else {
				break
			}
		}
		hashR = append(hashR, j)
	}

	hashR, err = hex.DecodeString(string(hashR))
	if err != nil {
		res.Err = ErrSearchHexDecode
		return res
	}

	res.Hash = checkCPUSHA256(strings.Trim(tmp1[0], " \"\t\n\r,."), string(hashR))

	return res
}

func checkCPUSHA256(cpuid, sha string) bool {
	b := []byte(cpuid)

	var b1 [32]byte
	for i, j := range []byte(sha) {
		b1[i] = j
	}
	b = append(b, []byte(DefaultSearchHash)...)
	s := sha256.Sum256(b)
	return s == b1
}

func stripCRC(b []byte) []byte {
	return b[:len(b)-4]
}

func stripAddr(b []byte) []byte {
	return b[1:]
}

func (r *RS485) SetAddr(addr byte, cpuid uint64) error {
	cpuS := strconv.FormatUint(cpuid, 10)
	addrS := strconv.FormatUint(uint64(addr), 10)
	tmp := sha256.Sum256([]byte(cpuS + addrS + DefaultSearchHash))
	sha := make([]byte, 32)
	for i, j := range tmp {
		sha[i] = j
	}

	shaS := hex.EncodeToString(sha)
	mess := "ATSETADDR=" + cpuS + "," + addrS + ",\"" + shaS + "\""

	resp := r.DoRequest(255, mess)

	fmt.Println(resp)
	if resp.Err != nil {
		fmt.Println(resp.Err.Error())
	}
	return resp.Err
}
