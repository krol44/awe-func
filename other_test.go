package aweFunc

import (
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	t.Parallel()

	tt := []struct {
		St string
	}{{St: "1"}, {St: "2"}, {St: "3"}, {St: "4"}}
	pp := PrettyPrint(tt)
	if len(pp) != 78 {
		t.Fail()
	}
}
