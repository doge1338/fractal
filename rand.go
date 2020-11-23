package main

import (
	"sync/atomic"
	"time"
)

// xorshift random

var randState = uint64(time.Now().UnixNano())

func RandUint64() uint64 {
	oldState := atomic.LoadUint64(&randState)
	newState := ((oldState ^ (oldState << 13)) ^ (oldState >> 7)) ^ (oldState << 17)
	atomic.StoreUint64(&randState, newState)
	return newState
}

func RandFloat64() float64 {
	return float64(RandUint64()/2) / (1 << 63)
}
