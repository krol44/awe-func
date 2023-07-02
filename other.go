package aweFunc

import "encoding/json"

// PrettyPrint return pretty interface for log or debug
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
