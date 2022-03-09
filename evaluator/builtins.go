package evaluator

import (
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":   object.GetBuiltinByName("len"),
	"type":  object.GetBuiltinByName("type"),
	"puts":  object.GetBuiltinByName("puts"),
	"first": object.GetBuiltinByName("first"),
	"last":  object.GetBuiltinByName("last"),
	"rest":  object.GetBuiltinByName("rest"),
	"push":  object.GetBuiltinByName("push"),
}
