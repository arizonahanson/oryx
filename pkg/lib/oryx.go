package lib

import (
	"errors"

	"github.com/arizonahanson/oryx/pkg/ast"
	"github.com/arizonahanson/oryx/pkg/eval"
)

func DoString(in string, env *eval.Env) (ast.Any, error) {
	baseEnv := BaseEnv(env)
	val, err := eval.EvalBytes([]byte(in), baseEnv)
	if err != nil {
		return ast.Null{}, err
	}
	return last(val)
}

func DoFile(filename string, env *eval.Env) (ast.Any, error) {
	baseEnv := BaseEnv(env)
	val, err := eval.EvalFile(filename, baseEnv)
	if err != nil {
		return ast.Null{}, err
	}
	return last(val)
}

func last(val ast.Any) (ast.Any, error) {
	switch seq := val.(type) {
	default:
		return val, nil
	case ast.Array:
		if len(seq) > 0 {
			return seq[len(seq)-1], nil
		}
		return ast.Null{}, errors.New("?empty")
	}
}
