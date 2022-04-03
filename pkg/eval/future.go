package eval

import (
	"github.com/arizonahanson/oryx/pkg/ast"
)

// trampoline to resolve futures
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

// resolve futures asynchronously and return new future
func (future Future) Go() (await Future) {
	channel := make(chan Future, 1)
	// resolve
	async := func() {
		val, err := future.Get()
		channel <- func() (ast.Any, error) {
			return val, err
		}
	}
	// await future
	await = func() (ast.Any, error) {
		return <-channel, nil
	}
	go async()
	return
}
