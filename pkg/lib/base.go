package lib

import (
	"fmt"

	"github.com/arizonahanson/oryx/pkg/ast"
	"github.com/arizonahanson/oryx/pkg/eval"
)

var BaseLib = map[string]eval.FuncType{
	"&&":     _and,
	"and":    _and,
	"||":     _or,
	"or":     _or,
	"==":     _equalQ,
	"equal?": _equalQ,
	/*
		==
		!=
		<
		<=
		>
		>=
		+
		-
		*
		/
		%
	*/
}

func BaseEnv(outer *eval.Env) *eval.Env {
	env := eval.NewEnv(outer)
	for key, fn := range BaseLib {
		env.SetFunc(key, fn)
	}
	return env
}

func exactLen(exp ast.Expr, n int) error {
	if len(exp) != n {
		return fmt.Errorf("%#v: wanted %d arg(s), got %d", exp[0], n-1, len(exp)-1)
	}
	return nil
}

func minLen(exp ast.Expr, n int) error {
	if len(exp) < n {
		return fmt.Errorf("%#v: wanted at least %d arg(s), got %d", exp[0], n-1, len(exp)-1)
	}
	return nil
}

func _equalQ(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	if err := minLen(exp, 3); err != nil {
		return ast.Null{}, err
	}
	val1, err := eval.Eval(exp[1], env)
	if err != nil {
		return ast.Null{}, err
	}
	for _, item := range exp[2:] {
		val2, err := eval.Eval(item, env)
		if err != nil {
			return ast.Null{}, err
		}
		if !val1.Equal(val2) {
			return ast.Boolean(false), nil
		}
	}
	return ast.Boolean(true), nil
}

func _and(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	if len(exp) == 1 {
		return ast.Boolean(true), nil
	}
	for _, item := range exp[1 : len(exp)-1] {
		val, err := eval.Eval(item, env)
		if err != nil {
			return ast.Null{}, err
		}
		if (val.Equal(ast.Boolean(false)) || val.Equal(ast.Null{})) {
			return val, nil
		}
	}
	return eval.EvalFuture(exp[len(exp)-1], env), nil
}

func _or(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	if len(exp) == 1 {
		return ast.Boolean(false), nil
	}
	for _, item := range exp[1 : len(exp)-1] {
		val, err := eval.Eval(item, env)
		if err != nil {
			return ast.Null{}, err
		}
		if (!val.Equal(ast.Boolean(false)) && !val.Equal(ast.Null{})) {
			return val, nil
		}
	}
	return eval.EvalFuture(exp[len(exp)-1], env), nil
}
