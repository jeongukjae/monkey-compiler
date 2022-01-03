package main

import (
	"fmt"
	"monkey/repl"
	"strings"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)

	in := make(chan string)
	out := make(chan string)

	fmt.Println("Initializing wasm")
	go repl.StartChannel(in, out)

	js.Global().Set("writeCommand", js.FuncOf(func(this js.Value, s []js.Value) interface{} {
		if len(s) == 0 {
			return js.ValueOf("")
		}

		command := strings.Trim(s[0].String(), " ")
		if command == "" {
			return js.ValueOf("")
		}

		in <- command
		output := <-out

		return js.ValueOf(output)
	}))

	<-c
}
