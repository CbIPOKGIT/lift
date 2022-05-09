package rs485

import "time"

type Response struct {
	Response []byte
	Duration time.Duration
	Err      error
}
