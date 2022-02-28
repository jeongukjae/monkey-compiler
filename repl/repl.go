package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
)

const PROMPT = ">> "

const MONKEY_FACE = `
 .--.  .-"      "-.  .--.
/ .. \/  .-. .-.   \/ .. \
| | '|  /   Y   \  |' |  |
| \  \  \ 0 | 0 /  /  /  |
\ '- ,\.-"""""""-./, -' /
 ''-' /_   ^ ^   _\ '-''
     |  \._   _./  |
     \   \ '~' /   /
      '._ '-=-' _.'
         '-----'
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	inChan := make(chan string)
	outChan := make(chan string)

	go StartChannel(inChan, outChan)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		inChan <- scanner.Text()
		output := <-outChan
		io.WriteString(out, output)
	}
}

func StartChannel(in chan string, out chan string) {
	constants := []object.Object{}
	globals := make([]object.Object, 10)
	symbolTable := compiler.NewSymbolTable()

	for {
		line := <-in

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			out <- printParseErrors(p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			out <- fmt.Sprintf("Woops! Compilation failed:\n%s\n", err)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalState(code, globals)
		err = machine.Run()
		if err != nil {
			out <- fmt.Sprintf("Woops! Executing bytecode failed:\n%s\n", err)
		}

		stackTop := machine.LastPoppedStackElement()
		out <- stackTop.Inspect() + "\n"
	}
}

func printParseErrors(errors []string) string {
	var out bytes.Buffer
	out.WriteString(MONKEY_FACE)
	out.WriteString("Woops! We ran into some monkey business here!\n")
	out.WriteString("parser errors:\n")
	for _, msg := range errors {
		out.WriteString("\t" + msg + "\n")
	}
	return out.String()
}
