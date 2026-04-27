package output

import (
	"encoding/json"
	"fmt"
	"os"
)

type response struct {
	OK    bool   `json:"ok"`
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func Success(data any) {
	b, _ := json.Marshal(response{OK: true, Data: data})
	fmt.Println(string(b))
}

func Failure(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	b, _ := json.Marshal(response{OK: false, Error: err.Error()})
	fmt.Println(string(b))
	os.Exit(1)
}
