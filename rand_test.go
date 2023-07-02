package aweFunc

import (
	"fmt"
	"testing"
)

func TestRandInt(t *testing.T) {
	i := RandInt(1, 10)
	if fmt.Sprintf("%T", i) != "int" || i <= 0 {
		t.Fail()
	}
}

func TestUniqueID(t *testing.T) {
	s := UniqueID("test")
	if len(s) != 20 {
		t.Fail()
	}
}
