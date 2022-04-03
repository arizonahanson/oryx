package eval

import "github.com/arizonahanson/oryx/pkg/ast"

// type:future
type Future func() (ast.Any, error)

func (future Future) String() string {
	return future.GoString()
}

func (future Future) GoString() string {
	return "???"
}

func (future Future) Equal(any ast.Any) bool {
	// not comparable
	return false
}

// type:function
type Func func(exp ast.Expr, env *Env) (ast.Any, error)

func (fn Func) String() string {
	// TODO
	return "func"
}

func (fn Func) GoString() string {
	// TODO
	return "func"
}

func (fn Func) Equal(any ast.Any) bool {
	// not comparable
	return false
}
