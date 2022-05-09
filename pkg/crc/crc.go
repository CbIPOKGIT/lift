package crc

import "github.com/sigurn/crc16"

func AddCRC16X25(data []byte) []byte {
	t := crc16.MakeTable(crc16.CRC16_X_25)
	cs := crc16.Checksum(data, t)
	data = append(data, uint8(cs>>8), uint8(cs&0xff))
	return data
}

func CRC16X25(data []byte) []byte {
	cs := crc16.Checksum(data, crc16.MakeTable(crc16.CRC16_X_25))
	return []byte{uint8(cs >> 8), uint8(cs & 0xff)}
}
