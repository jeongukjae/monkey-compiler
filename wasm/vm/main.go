package main

import (
	"fmt"
	"strings"
	"syscall/js"
	"time"

	"monkey/compiler"
	"monkey/lexer"
	"monkey/parser"
	"monkey/vm"
)

func main() {
	fmt.Println("Initializing wasm")

	js.Global().Set("compileAndRun", js.FuncOf(func(this js.Value, s []js.Value) interface{} {
		if len(s) == 0 {
			return js.ValueOf("")
		}

		input := strings.Trim(s[0].String(), " ")
		if input == "" {
			return js.ValueOf("")
		}

		result := map[string]interface{}{
			"ErrorString":            "",
			"ElapsedTimeCompilation": -1,
			"ElapsedTimeVMInit":      -1,
			"ElapsedTimeRuntime":     -1,
			"Result":                 "",
		}

		start := time.Now()
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			result["ErrorString"] = fmt.Sprintf("compiler error: %s", err)
			return js.ValueOf(result)
		}
		result["ElapsedTimeCompilation"] = time.Since(start).Milliseconds()

		start = time.Now()
		machine := vm.New(comp.Bytecode())
		result["ElapsedTimeVMInit"] = time.Since(start).Milliseconds()

		start = time.Now()
		err = machine.Run()
		if err != nil {
			result["ErrorString"] = fmt.Sprintf("vm error: %s", err)
			return js.ValueOf(result)
		}
		result["ElapsedTimeRuntime"] = time.Since(start).Milliseconds()
		result["Result"] = machine.LastPoppedStackElement().Inspect()

		return js.ValueOf(result)

	}))

	c := make(chan bool)
	<-c
}
