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

func TestCallingFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let fivePlusTen = fn() {5 + 10;};
			fivePlusTen();
			`,
			expected: 15,
		},
		{
			input: `
			let one = fn() {1;}
			let two = fn() {2;}
			one() + two()
			`,
			expected: 3,
		},
		{
			input: `
			let a = fn() {1;}
			let b = fn() {a () + 1}
			let c = fn() {b () + 1}
			c()
			`,
			expected: 3,
		},
		{
			input: `
			let earlyExit = fn() {return 99;100;}
			earlyExit()
			`,
			expected: 99,
		},
		{
			input: `
			let earlyExit = fn() {return 99;return 100;}
			earlyExit()
			`,
			expected: 99,
		},
		{
			input: `
			let noReturn = fn() {}
			noReturn()
			`,
			expected: Null,
		},
		{
			input: `
			let noReturn = fn() {};
			let noReturn2 = fn() {noReturn()};
			noReturn2();
			`,
			expected: Null,
		},
		{
			input: `
			let return1 = fn() {1};
			let returnReturn1 = fn() {return1};
			returnReturn1()();
			`,
			expected: 1,
		},
		{
			input: `
			let returnsOneReturner = fn() {
				let returnsOne = fn() {1;};
				returnsOne
			}
			returnsOneReturner()();
			`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let one = fn() {let one =  1; one};
			one();`,
			expected: 1,
		},
		{
			input: `
			let oneAndTwo = fn() {
				let one =1 ;
				let two = 2;
				one + two;
			}
			oneAndTwo();`,
			expected: 3,
		},
		{
			input: `
			let oneAndTwo = fn() {
				let one =1 ;
				let two = 2;
				one + two;
			}
			let threeAndFour = fn() {
				let three = 3;
				let four = 4;
				three + four;
			}
			oneAndTwo() + threeAndFour();`,
			expected: 10,
		},
		{
			input: `
			let firstFooBar = fn() { let foobar = 50; foobar; }
			let secondFoobar = fn() { let foobar = 100; foobar; }
			firstFooBar() + secondFoobar();`,
			expected: 150,
		},
		{
			input: `
			let globalSeed = 50;
			let minusOne = fn() { let num = 1; globalSeed - num }
			let minusTwo = fn() { let num = 2; globalSeed - num }
			minusOne() + minusTwo();`,
			expected: 97,
		},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

func TestCallingFunctionsWithBindingsAndArgs(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let a = fn(b) {b}
			a(4);`,
			expected: 4,
		},
		{
			input: `
			let sum = fn(a, b) {a + b;}
			sum(1, 2);`,
			expected: 3,
		},
		{
			input: `
			let sum = fn(a, b) {
				let c = a +b;
				return c;
			}
			sum(1, 2);`,
			expected: 3,
		},
		{
			input: `
			let sum = fn(a, b) {
				let c = a +b;
				return c;
			}
			sum(1, 2) + sum(3, 4);`,
			expected: 10,
		},
		{
			input: `
			let sum = fn(a, b) {
				let c = a +b;
				return c;
			}
			let outer = fn() {
				return sum(1, 2) + sum(3, 4);
			}
			outer();`,
			expected: 10,
		},
		{
			input: `
			let globalNum = 10;
			let sum = fn(a, b) {
				let c = a +b;
				return c + globalNum;
			}

			let outer = fn() {
				sum(1,2) + sum(3,4) +globalNum;
			}

			outer() + globalNum;
			`,
			expected: 50,
		},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `fn() {1;}(1);`,
			expected: `wrong number of arguments: want=0, got=1`,
		},
		{
			input:    `fn(a) {a}();`,
			expected: `wrong number of arguments: want=1, got=0`,
		},
		{
			input:    `fn(a, b) {a+b}(1);`,
			expected: `wrong number of arguments: want=2, got=1`,
		},
	}

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		require.Nil(t, err, "compiler error")

		vm := New(comp.Bytecode())
		err = vm.Run()
		require.NotNil(t, err, "expected VM error but resulted in none.")

		require.Equal(t, tt.expected, err.Error(), "wrong vm error")
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, &object.Error{
			Message: "argument to `len` not supported, got INTEGER",
		}},
		{`len("one", "two")`, &object.Error{
			Message: "wrong number of arguments. got=2, want=1",
		}},
		{`len([1,2,3])`, 3},
		{`len([])`, 0},
		{`puts("hello", "world!")`, Null},
		{`first([1,2,3])`, 1},
		{`first([])`, Null},
		{`first(1)`, &object.Error{
			Message: "argument to `first` must be ARRAY, got INTEGER",
		}},
		{`last([1,2,3])`, 3},
		{`last([])`, Null},
		{`last(1)`, &object.Error{
			Message: "argument to `last` must be ARRAY, got INTEGER",
		}},
		{`rest([1,2,3])`, []int{2, 3}},
		{`rest([])`, Null},
		{`push([], 1)`, []int{1}},
		{`push(1,1)`, &object.Error{
			Message: "argument to `push` must be ARRAY, got INTEGER",
		}},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let newClosure = fn(a) {
				fn() {a;}
			}
			let closure = newClosure(99);
			closure();
			`,
			expected: 99,
		},
		{
			input: `
			let newAdder = fn(a, b) {
				fn(c) { a+ b + c};
			}
			let adder = newAdder(1,2);
			adder(8)
			`,
			expected: 11,
		},
		{
			input: `
			let newAdder = fn(a, b) {
				let c = a + b;
				fn(d) { c+d};
			}
			let adder = newAdder(1,2);
			adder(8)
			`,
			expected: 11,
		},
		{
			input: `
			let newAdderOuter = fn(a, b) {
				let c = a + b;
				fn(d) {
					let e = d + c;
					fn(f) { e + f; };
				}
			}
			let newAdderInner = newAdderOuter(1, 2);
			let adder = newAdderInner(3);
			adder(8);
			`,
			expected: 14,
		},
		{
			input: `
			let a = 1;
			let newAdderOuter = fn(b) {
				fn(c) {
					fn(d) {a + b + c + d};
				}
			}
			let newAdderInner = newAdderOuter(2);
			let adder = newAdderInner(3);
			adder(8);
			`,
			expected: 14,
		},
		{
			input: `
			let newClosure = fn(a, b) {
				let one = fn() {a;}
				let two = fn() {b;}
				fn() { one() + two() }
			}
			let closure = newClosure(9, 90);
			closure();
			`,
			expected: 99,
		},
	}

	for _, tt := range tests {
		runVmTest(t, tt)
	}
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let countDown = fn(x) {
				if (x == 0) {
					return 0;
				}
				return countDown(x - 1);
			}
			countDown(1);
			`,
			expected: 0,
		},
		{
			input: `
			let countDown = fn(x) {
				if (x == 0) {
					return 0;
				}
				return countDown(x - 1);
			}
			let wrapper = fn() {
				countDown(1);
			}
			wrapper()
			`,
			expected: 0,
		},
		{
			input: `
			let wrapper = fn() {
				let countDown = fn(x) {
					if (x == 0) {
						return 0;
					}
					return countDown(x - 1);
				}

				countDown(1);
			}
			wrapper();
			`,
			expected: 0,
		},
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
	case *object.Error:
		errObj, ok := actual.(*object.Error)
		require.True(t, ok)
		require.Equal(t, expected.Message, errObj.Message)
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
