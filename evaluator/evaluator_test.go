package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 -10", 10},
		{"2 * 2 * 2* 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"50 / 2 * 2 + 10", 60},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, tt.expected, evaluated)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world!";`, "hello world!"},
		{`"hello" + " " + "world!";`, "hello world!"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, tt.expected, evaluated)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 > 1", false},
		{"1 < 1", false},
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
		{`"hello" + "world" == "helloworld"`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, tt.expected, evaluated)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, tt.expected, evaluated)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) {10}", 10},
		{"if (false) {10}", nil},
		{"if (1) {10}", 10},
		{"if (1 < 2) {10}", 10},
		{"if (1 > 2) {10}", nil},
		{"if (1 > 2) {10} else {20}", 20},
		{"if (1 < 2) {10} else {20}", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, int64(integer), evaluated)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 5; 9;", 5},
		{"return 2 * 5;9;", 10},
		{"9; return 2 * 5;9;", 10},
		{`
		if (10 > 1) {
			if (10  >1 ) {
				return 10;
			}
			return 5;
		}
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, tt.expected, evaluated)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"- true", "unknown operator: -BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 5){ true + false;}", "unknown operator: BOOLEAN + BOOLEAN"},
		{`
		if (10 > 5){
			if (true) {
				return true + false;
			}
			return 1;
		}`, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{`{"name": "Monkey"}[fn(x) {x}]`, "unhashable as hash key: FUNCTION"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		require.True(t, ok, "no error object returned")
		require.Equal(t, tt.expectedMessage, errObj.Message, "wrong error message")
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b= a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, tt.expected, evaluated)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) {x + 2;};"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	require.True(t, ok, "object is not function")
	require.Equal(t, 1, len(fn.Parameters), "function has wrong parameters")
	require.Equal(t, "x", fn.Parameters[0].String(), "wrong parameter")
	require.Equal(t, "(x + 2)", fn.Body.String(), "wrong body")
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) {x;} identity(5);", 5},
		{"let identity = fn(x) { return x;} identity(5);", 5},
		{"let double = fn(x) { return x * 2;} double(5);", 10},
		{"let add = fn(x, y) { x + y;} add(5, 5);", 10},
		{"let add = fn(x, y) { x + y;} add(5 + 5, add(5, 5));", 20},
		{"fn(x) {x;}(5)", 5},
		{"(fn(x) {x;})(5)", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, tt.expected, evaluated)
	}
}

func TestClosure(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		return fn(y) {
			return x + y;
		}
	}

	let addTwo = newAdder(2);
	addTwo(2);
	`

	evaluated := testEval(input)
	testIntegerObject(t, 4, evaluated)
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(2)`, "argument to len not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, int64(expected), evaluated)
		case string:
			errObj, ok := evaluated.(*object.Error)
			require.True(t, ok, "object is not error")
			require.Equal(t, expected, errObj.Message, "wrong error message")
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)

	require.True(t, ok, "object is not array, %s", result)
	require.Equal(t, 3, len(result.Elements), "array has wrong num of elements, %s", result)

	testIntegerObject(t, 1, result.Elements[0])
	testIntegerObject(t, 4, result.Elements[1])
	testIntegerObject(t, 6, result.Elements[2])
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`[1,2,3][0]`, 1},
		{`[1,2,3][1]`, 2},
		{`[1,2,3][2]`, 3},
		{`let i = 0; [1][i]`, 1},
		{`[1,2,3][1+1]`, 3},
		{`let myArray = [1,2,3]; myArray[2];`, 3},
		{`let myArray = [1,2,3]; myArray[0] + myArray[1] + myArray[2]`, 6},
		{`let myArray = [1,2,3]; let i = myArray[0]; myArray[i]`, 2},
		{`[1,2,3][3]`, nil},
		{`[1,2,3][-1]`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, int64(expected), evaluated)
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiteral(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 -9,
		two: 1 +1,
		"thr" + "ee": 6 / 2,
		4 : 4,
		true: 5,
		false: 6
	}
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)

	require.True(t, ok, "object is not Hash, %s", evaluated)
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		(&object.Boolean{Value: true}).HashKey():   5,
		(&object.Boolean{Value: false}).HashKey():  6,
	}

	require.Equal(t, len(expected), len(result.Pairs))
	for key, value := range expected {
		pair, ok := result.Pairs[key]
		require.True(t, ok)
		testIntegerObject(t, value, pair.Value)
	}
}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["bar"]`,
			nil,
		},
		{
			`{5:5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, int64(integer), evaluated)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

//
// Helper functions
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, expected int64, actual object.Object) {
	result, ok := actual.(*object.Integer)
	require.True(t, ok, "object is not integer, %s", actual)
	require.Equal(t, expected, result.Value, "object has wrong value")
}

func testStringObject(t *testing.T, expected string, actual object.Object) {
	result, ok := actual.(*object.String)
	require.True(t, ok, "object is not string, %s", actual)
	require.Equal(t, expected, result.Value, "object has wrong value")
}

func testBooleanObject(t *testing.T, expected bool, actual object.Object) {
	result, ok := actual.(*object.Boolean)
	require.True(t, ok, "object is not boolean, got %s", actual)
	require.Equal(t, expected, result.Value, "object has wrong value")
}

func testNullObject(t *testing.T, actual object.Object) {
	require.Equal(t, NULL, actual, "object is not NULL")
}
