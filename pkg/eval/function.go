package eval

import "github.com/arizonahanson/oryx/pkg/ast"

func (fn Func) Future(exp ast.Expr, env *Env) Future {
	return func() (ast.Any, error) {
		return fn(exp, env)
	}
}
