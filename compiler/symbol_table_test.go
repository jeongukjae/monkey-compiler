package compiler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
		"c": {Name: "c", Scope: LocalScope, Index: 0},
		"d": {Name: "d", Scope: LocalScope, Index: 1},
		"e": {Name: "e", Scope: LocalScope, Index: 0},
		"f": {Name: "f", Scope: LocalScope, Index: 1},
	}

	global := NewSymbolTable()

	a := global.Define("a")
	require.Equal(t, expected["a"], a)

	b := global.Define("b")
	require.Equal(t, expected["b"], b)

	firstLocal := NewEnclosedSymbolTable(global)
	c := firstLocal.Define("c")
	require.Equal(t, expected["c"], c)

	d := firstLocal.Define("d")
	require.Equal(t, expected["d"], d)

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	e := secondLocal.Define("e")
	require.Equal(t, expected["e"], e)

	f := secondLocal.Define("f")
	require.Equal(t, expected["f"], f)
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		require.True(t, ok, "name %s not resolvable", sym.Name)
		require.Equal(t, sym, result)
	}
}

func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
		{Name: "c", Scope: LocalScope, Index: 0},
		{Name: "d", Scope: LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := local.Resolve(sym.Name)
		require.True(t, ok, "name", sym.Name, "is not resolvable")
		require.Equal(t, sym, result)
	}
}

func TestResolveNestedLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	expected := []struct {
		table          *SymbolTable
		exepctedSymbol []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			secondLocal,
			[]Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range expected {
		for _, sym := range tt.exepctedSymbol {
			result, ok := tt.table.Resolve(sym.Name)
			require.True(t, ok, "name", sym.Name, "is not resolvable")
			require.Equal(t, sym, result)
		}
	}
}
