package aweFunc

import (
	"testing"
)

func TestSliceChunk(t *testing.T) {
	t.Parallel()

	type TyTemp struct {
		St string
	}
	tt := []TyTemp{{St: "1"}, {St: "2"}, {St: "3"}, {St: "4"}}

	sc := SliceChunk(tt, 2)

	if len(sc) != 2 {
		t.Fail()
	}

	if sc[1][1].St != "4" {
		t.Fail()
	}
}

func TestSliceReverse(t *testing.T) {
	t.Parallel()

	tt := []int{1, 2}
	sr := SliceReverse(tt)
	if sr[0] != 2 {
		t.Fail()
	}
}
