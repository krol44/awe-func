package aweFunc

import (
	"fmt"
	"testing"
)

func TestRandInt(t *testing.T) {
	t.Parallel()

	i := RandInt(1, 10)
	if fmt.Sprintf("%T", i) != "int" || i <= 0 {
		t.Fail()
	}
}

func TestUniqueID(t *testing.T) {
	t.Parallel()

	s := UniqueID("test")
	if len(s) != 20 {
		t.Fail()
	}
}
