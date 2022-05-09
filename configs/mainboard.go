package configs

const Rs232Adr = "/dev/ttyS2"
const Rs232Baud = 115200
const Rs232Timeout = 300
const Rs232StopByte = 0x0D

const Rs485Adr = "/dev/ttyUSB0"
const Rs485Baud = 115200
const Rs485Timeout = 600
const Rs485StartByte = 0xFA
const Rs485StopByte = 0xFE

func Rs232Config() (string, int, int, byte) {
	return Rs232Adr, Rs232Baud, Rs232Timeout, Rs232StopByte
}

func Rs485Config() (string, int, int) {
	return Rs485Adr, Rs485Baud, Rs485Timeout
}
