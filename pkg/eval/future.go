package eval

import (
	"fmt"

	"github.com/arizonahanson/oryx/pkg/ast"
)

// trampoline to resolve future values
func (future Future) Get() (val ast.Any, err error) {
	val, err = future()
	for {
		if err != nil {
			return
		}
		switch future := val.(type) {
		default:
			return
		case Future:
			val, err = future()
		}
	}
}

// resolve future asynchronously and return new future
func (future Future) Go() (await Future) {
	tunnel := make(chan Future, 1)
	// resolve
	async := func() {
		val, err := future.Get()
		tunnel <- func() (ast.Any, error) {
			return val, err
		}
	}
	// await future
	await = func() (ast.Any, error) {
		return <-tunnel, nil
	}
	go async()
	return
}

// trace errors mapped to source
func (future Future) Trace(exp ast.Expr) Future {
	return func() (val ast.Any, err error) {
		val, err = future.Get()
		if err != nil {
			err = fmt.Errorf("%#v: %s", exp[0], err)
		}
		return
	}
}
