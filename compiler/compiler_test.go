package compiler

import (
	"monkey/ast"
	"monkey/code"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"

	"github.com/stretchr/testify/require"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
			},
		},
	}

	for _, tt := range tests {
		runCompilerTests(t, tt)
	}
}

// Helper functions
func runCompilerTests(t *testing.T, tt compilerTestCase) {
	t.Helper()

	program := parse(tt.input)
	compiler := New()
	err := compiler.Compile(program)
	require.Nil(t, err)

	bytecode := compiler.Bytecode()
	testInstructions(t, tt.expectedInstructions, bytecode.Instructions)
	testConstants(t, tt.expectedConstants, bytecode.Constants)
}

func testInstructions(
	t *testing.T,
	expected []code.Instructions,
	actual code.Instructions,
) {
	concatenated := concatInstruction(expected)

	require.Equal(t, len(concatenated), len(actual), "wrong instructions length")
	require.ElementsMatch(t, concatenated, actual, "wrong instruction")
}

func testConstants(
	t *testing.T,
	expected []interface{},
	actual []object.Object,
) {
	require.Equal(t, len(expected), len(actual), "wrong number of constatns")

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			testIntegerObject(t, int64(constant), actual[i])
		}
	}
}

func testIntegerObject(t *testing.T, expected int64, actual object.Object) {
	result, ok := actual.(*object.Integer)
	require.True(t, ok)
	require.Equal(t, expected, result.Value, "object has wrong value")
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func concatInstruction(s []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, ins := range s {
		out = append(out, ins...)
	}

	return out
}
