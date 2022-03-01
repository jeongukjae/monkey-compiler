package vm

import (
	"testing"

	"github.com/stretchr/testify/require"

	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerVm(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 > 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"false == true", false},
		{"false != true", true},
		{"true != false", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) {5;})", true},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) {10} ", 10},
		{"if (true) {10} else {20}", 10},
		{"if (false) {10} else {20}", 20},
		{"if (1) {10}", 10},
		{"if (1 < 2) {10}", 10},
		{"if (1 < 2) {10} else {20}", 10},
		{"if (1 > 2) {10} else {20}", 20},
		{"if (1 > 2) {10}", Null},
		{"if (false) {10}", Null},
		{"if ((if(false) {true})) {10} else {20}", 20},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two  = 2; one + two", 3},
		{"let one = 1; let two  = one + one ; one + two", 3},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" +"key"`, "monkey"},
		{`"mon" +"key" + "banana"`, "monkeybanana"},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1,2,3]", []int{1, 2, 3}},
		{"[1 +2, 3*4, 5 +6]", []int{3, 12, 11}},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"{}", map[object.HashKey]int64{}},
		{
			"{1:2,2:3}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			},
		},
		{
			"{1 + 1:2*2, 3+3:4*4}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1,2,3][1]", 2},
		{"[1,2,3][0 + 2]", 3},
		{"[[1,1,1]][0][0]", 1},
		{"[][0]", Null},
		{"[1,2,3][99]", Null},
		{"[1][-1]", Null},
		{"{1:1,2:2}[1]", 1},
		{"{1:1,2:2}[2]", 2},
		{"{1:1,2:2}[3]", Null},
		{"{}[1]", Null},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

// helper functions
func runVmTest(t *testing.T, tt vmTestCase) {
	t.Helper()

	program := parse(tt.input)

	comp := compiler.New()
	err := comp.Compile(program)
	require.Nil(t, err, "compiler error")

	vm := New(comp.Bytecode())
	err = vm.Run()
	require.Nil(t, err, "vm error")

	stackElem := vm.LastPoppedStackElement()
	testExpectObject(t, tt.expected, stackElem)
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testExpectObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		testIntegerObject(t, int64(expected), actual)
	case bool:
		testBooleanObject(t, bool(expected), actual)
	case *object.Null:
		require.Equal(t, Null, actual, "object is not Null")
	case string:
		testStringObject(t, expected, actual)
	case []int:
		testIntegerArrayObject(t, expected, actual)
	case map[object.HashKey]int64:
		testIntegerHashObject(t, expected, actual)
	}
}

func testIntegerObject(t *testing.T, expected int64, actual object.Object) {
	result, ok := actual.(*object.Integer)
	require.True(t, ok)
	require.Equal(t, expected, result.Value, "object has wrong value")
}

func testBooleanObject(t *testing.T, expected bool, actual object.Object) {
	result, ok := actual.(*object.Boolean)
	require.True(t, ok)
	require.Equal(t, expected, result.Value, "object has wrong value")
}

func testStringObject(t *testing.T, expected string, actual object.Object) {
	result, ok := actual.(*object.String)
	require.True(t, ok)
	require.Equal(t, expected, result.Value, "object has wrong value")
}

func testIntegerArrayObject(t *testing.T, expected []int, actual object.Object) {
	result, ok := actual.(*object.Array)
	require.True(t, ok)
	require.Equal(t, len(expected), len(result.Elements), "wrong number of elements")
	for i, expectedElement := range expected {
		testIntegerObject(t, int64(expectedElement), result.Elements[i])
	}
}
func testIntegerHashObject(t *testing.T, expected map[object.HashKey]int64, actual object.Object) {
	result, ok := actual.(*object.Hash)
	require.True(t, ok)
	require.Equal(t, len(expected), len(result.Pairs), "wrong number of elements")
	for key, value := range expected {
		pair, ok := result.Pairs[key]
		require.True(t, ok, "no pair for given key")
		testIntegerObject(t, value, pair.Value)
	}
}
