package compiler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
	}

	global := NewSymbolTable()

	a := global.Define("a")
	require.Equal(t, expected["a"], a)

	b := global.Define("b")
	require.Equal(t, expected["b"], b)
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
