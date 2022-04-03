package parser

import "github.com/arizonahanson/oryx/pkg/ast"

// cast to []interface{}
func slice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

// build []ast.Any from first, rest=[[_, next], ...]
func join(first, rest interface{}, index int) []ast.Any {
	if first == nil {
		return []ast.Any{}
	}
	more := slice(rest)
	result := make([]ast.Any, len(more)+1)
	result[0] = first.(ast.Any)
	for i, group := range more {
		next := slice(group)[index]
		result[i+1] = next.(ast.Any)
	}
	return result
}

func merge(first, rest interface{}, keyIndex int, valueIndex int) map[ast.String]ast.Any {
	pair := slice(first)
	if pair == nil {
		return map[ast.String]ast.Any{}
	}
	more := slice(rest)
	result := make(map[ast.String]ast.Any, len(more)+1)
	// assign helper
	assign := func(keyval []interface{}, keyN int, valN int) {
		key := keyval[keyN].(ast.String)
		result[key] = keyval[valN].(ast.Any)
	}
	// assign pairs
	assign(pair, keyIndex, valueIndex)
	for _, group := range more {
		pair := slice(group)
		assign(pair, keyIndex+1, valueIndex+1)
	}
	return result
}

func pos(p position) *ast.Position {
	return &ast.Position{Row: int64(p.line), Column: int64(p.col), Offset: int64(p.offset)}
}
