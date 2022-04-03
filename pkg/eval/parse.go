package eval

import (
	"io"

	"github.com/arizonahanson/oryx/internal/parser"
	"github.com/arizonahanson/oryx/pkg/ast"
)

func Parse(in []byte) (ast.Any, error) {
	val, err := parser.Parse("parse", in)
	if err != nil {
		return ast.Null{}, err
	}
	return val.(ast.Any), nil
}

func ParseFile(filename string) (ast.Any, error) {
	val, err := parser.ParseFile(filename)
	if err != nil {
		return ast.Null{}, err
	}
	return val.(ast.Any), nil
}

func ParseReader(read io.Reader) (ast.Any, error) {
	val, err := parser.ParseReader("read", read)
	if err != nil {
		return ast.Null{}, err
	}
	return val.(ast.Any), nil
}
