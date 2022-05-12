package conv

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const (
	BinLength     int    = 8
	ErrOutOfRange string = "ErrOutOfRange¬"
)

type Binary [BinLength]bool

var table []int

func init() {
	table = make([]int, 0, 32)
	for i := 0; i < BinLength; i++ {
		table = append(table, pow2x(i))
	}
}

func NewBinary() Binary {
	return Binary{}
}

func (b Binary) String() string {
	var s string
	for i := BinLength - 1; i >= 0; i-- {
		if b[i] {
			s += "1"
		} else {
			s += "0"
		}
	}
	return s
}

func (b *Binary) SetInt(i int) {
	*b = int2Bin(i)
}

func (b Binary) GetInt() int {
	return bin2Int(b)
}

func (b Binary) GetIntString() string {
	return strconv.Itoa(b.GetInt())
}

func (b *Binary) SetString(s string) {
	i, _ := parseResp(s)
	*b = int2Bin(i)
}

func (b Binary) GetBit(num int) (bool, error) {
	err := checkRange(num)
	if err != nil {
		return false, errors.New("out of range")
	}

	return b[num], nil
}

func (b Binary) GetBitInt64(num int) (int64, error) {
	if checkRange(num) != nil {
		return 0, errors.New(ErrOutOfRange)
	}
	if b[num] {
		return int64(1), nil
	}
	return int64(0), nil
}

func (b *Binary) SetBit(num int, value bool) error {
	err := checkRange(num)
	if err != nil {
		return err
	}

	b[num] = value

	return nil
}

func bin2Int(b Binary) int {
	var s int

	for i := 0; i < BinLength; i++ {
		if b[i] {
			s += table[i]
		}
	}

	return s
}

func int2Bin(p int) [BinLength]bool {
	ar := [BinLength]bool{}

	for i := BinLength - 1; i >= 0; i-- {
		m := table[i]
		if p/m == 1 {
			ar[i] = true
			p = p - m
		} else {
			ar[i] = false
		}
	}
	return ar
}

func pow2x(x int) int {
	return int(math.Pow(float64(2), float64(x)))
}

func checkRange(x int) error {
	if x < 0 || x >= BinLength {
		return errors.New("bit index out of range")
	}
	return nil
}

func parseResp(str string) (int, error) {
	s := strings.Split(str, "=")
	if len(s) != 2 {
		return 0, errors.New("error split")
	}
	i, err := strconv.Atoi(s[1])
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (b Binary) GetBitUnsafe(num int) bool {
	err := checkRange(num)
	if err != nil {
		return false
	}

	return b[num]
}
