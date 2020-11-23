package main

import (
	"time"
)

// xorshift random

var randState = uint64(time.Now().UnixNano())

func RandUint64() uint64 {
	randState = ((randState ^ (randState << 13)) ^ (randState >> 7)) ^ (randState << 17)
	return randState
}

func RandFloat64() float64 {
	return float64(RandUint64() / 2) / (1 << 63)
}
