package eval

import (
	"fmt"

	"github.com/arizonahanson/oryx/pkg/ast"
)

// environment or scope for symbols
type Env struct {
	parent *Env
	data   map[string]ast.Any
}

// find an environment and value using a symbol
func (env *Env) find(symbol ast.Symbol) (*Env, ast.Any) {
	scope := env
	for {
		val, ok := scope.data[symbol.Val]
		if !ok {
			if scope.parent != nil {
				scope = scope.parent
				continue
			}
			return nil, ast.Null{}
		}
		return scope, val
	}
}

// get a value from the environment using a symbol
func (env *Env) Get(symbol ast.Symbol) (val ast.Any, err error) {
	scope, val := env.find(symbol)
	if scope == nil {
		err = fmt.Errorf("%#v: not found", symbol)
	}
	return
}

// set a value to the environment
func (env *Env) Set(symbol ast.Symbol, any ast.Any) (val ast.Any) {
	val = any
	switch future := any.(type) {
	default:
		break
	case Future:
		// memoize future resolution
		val = Future(func() (res ast.Any, err error) {
			res, err = future.Get()
			env.Set(symbol, res)
			return
		})
	}
	env.data[symbol.Val] = val
	return
}
