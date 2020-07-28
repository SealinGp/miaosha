package main

import (
	"encoding/json"
	"fmt"
)

type a struct {
	B b  `json:"b"`
	C c  `json:"c"`
}
type b struct {
	B1 string `json:"b_1"`
}
type c struct {
	C1 string `json:"c_1"`
}

//{"b":{"b1":"xxx"}},"c":{"c1":"xxx"}}
func main() {
	str := `{"b":{"b_1":"xxx"},"c":{"c_1":"xxx"}}`
	m := a{}
	json.Unmarshal([]byte(str),&m)
	fmt.Println(m)
}