package eval

import (
	"io"

	"github.com/arizonahanson/oryx/internal/parser"
	"github.com/arizonahanson/oryx/pkg/ast"
)

// parse a slice of bytes as an ast
func Parse(in []byte) (ast.Any, error) {
	val, err := parser.Parse("parse", in)
	if err != nil {
		return ast.Null{}, err
	}
	return val.(ast.Any), nil
}

// parse a file as an ast
func ParseFile(filename string) (ast.Any, error) {
	val, err := parser.ParseFile(filename)
	if err != nil {
		return ast.Null{}, err
	}
	return val.(ast.Any), nil
}

// parse reader output as an ast
func ParseReader(read io.Reader) (ast.Any, error) {
	val, err := parser.ParseReader("read", read)
	if err != nil {
		return ast.Null{}, err
	}
	return val.(ast.Any), nil
}
