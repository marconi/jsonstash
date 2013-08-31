package utils

import (
	"encoding/json"
	"fmt"
)

func CheckJSON(data string) bool {
	var v interface{}
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		fmt.Println("Invalid JSON: ", err)
		return false
	}
	return true
}

func ToJSON(data []string) string {
	out := ""
	for index, d := range data {
		out += d
		if (index + 1) != len(data) {
			out += ","
		}
	}
	return fmt.Sprintf("[%s]", out)
}
