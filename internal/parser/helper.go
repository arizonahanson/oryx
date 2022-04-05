package parser

import (
	"github.com/arizonahanson/oryx/pkg/ast"
)

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

func swap(first, rest interface{}, opIndex int, rightIndex int) []ast.Any {
	// return left if no right
	more := slice(rest)
	left := first.([]ast.Any)
	if len(more) == 0 || more[0] == nil {
		return left
	}
	// iterate right side
	result := make([]ast.Any, len(more)+2)
	i := 0
	for _, group := range more {
		frag := slice(group)
		op := frag[opIndex].(ast.Symbol)
		right := frag[rightIndex].([]ast.Any)
		if i == 0 {
			result[0] = op
			if len(left) > 1 {
				result[1] = ast.Expr(left)
			} else {
				result[1] = left[0]
			}
		}
		if !op.Equal(result[0]) {
			// new op
			newexpr := make([]ast.Any, len(result)-i)
			newexpr[0] = op
			newexpr[1] = ast.Expr(result[:i+2])
			result = newexpr
			i = 0
		}
		if len(right) > 1 {
			result[i+2] = ast.Expr(right)
		} else {
			result[i+2] = right[0]
		}
		i++
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

func symbol(c *current) (ast.Symbol, error) {
	return ast.Symbol{Val: string(c.text), Pos: pos(c.pos)}, nil
}

func pos(p position) *ast.Position {
	return &ast.Position{Row: int64(p.line), Column: int64(p.col), Offset: int64(p.offset)}
}
