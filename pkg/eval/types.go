package eval

import "github.com/arizonahanson/oryx/pkg/ast"

type Future func() (ast.Any, error)

func (future Future) String() string {
	return "?← " // should not happen
}

func (future Future) GoString() string {
	return "?← " // should not happen
}

func (future Future) Equal(any ast.Any) bool {
	return false // not comparable
}
