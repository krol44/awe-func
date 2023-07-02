package aweFunc

import (
	"fmt"
	"math/rand"
	"time"
)

// RandInt return rand int
func RandInt(min int, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(max-min+1) + min
}

// UniqueID return unique string id with prefix, `test-64a0643d4e76891`
// subject to collision
func UniqueID(prefix string) string {
	now := time.Now()
	sec := now.Unix()
	use := now.UnixNano() % 0x100000
	return fmt.Sprintf("%s-%d%08x%05x", prefix, RandInt(10, 99), sec, use)
}
