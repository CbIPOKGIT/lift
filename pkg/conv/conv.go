package conv

import (
	"strconv"
	"strings"
)

type BinaryConverter struct {
	binary []string
}

func NewBinary() *BinaryConverter {
	converter := new(BinaryConverter)
	converter.binary = make([]string, 0)
	return converter
}

func (bc *BinaryConverter) SetString(str string) error {
	val, err := strconv.Atoi(str)
	if err != nil {
		return err
	}

	binary := strings.Split(strconv.FormatInt(int64(val), 2), "")
	for i := 0; i < len(binary)/2; i++ {
		binary[i], binary[len(binary)-1-i] = binary[len(binary)-1-i], binary[i]
	}

	bc.binary = binary
	return nil
}

func (bc *BinaryConverter) GetBitBoolValue(pos int) bool {
	if len(bc.binary) <= pos {
		return false
	}
	return bc.binary[pos] == "1"
}
