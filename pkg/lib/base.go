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
	"!=":     _nequalQ,
	"<":      _ltQ,
	"lt?":    _ltQ,
	"<=":     _lteqQ,
	"lteq?":  _lteqQ,
	">":      _gtQ,
	"gt?":    _gtQ,
	">=":     _gteqQ,
	"gteq?":  _gteqQ,
	/*
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

func evalNumber(exp ast.Any, env *eval.Env) (*ast.Number, error) {
	val, err := eval.Eval(exp, env)
	if err != nil {
		return nil, err
	}
	switch num := val.(type) {
	default:
		return nil, fmt.Errorf("called with non-number %#v", val)
	case ast.Number:
		return &num, nil
	}
}

func _ltQ(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	if err := exactLen(exp, 3); err != nil {
		return ast.Null{}, err
	}
	val1, err := evalNumber(exp[1], env)
	if err != nil {
		return ast.Null{}, err
	}
	val2, err := evalNumber(exp[2], env)
	if err != nil {
		return ast.Null{}, err
	}
	return ast.Boolean(val1.Decimal().LessThan(val2.Decimal())), nil
}

func _lteqQ(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	if err := exactLen(exp, 3); err != nil {
		return ast.Null{}, err
	}
	val1, err := evalNumber(exp[1], env)
	if err != nil {
		return ast.Null{}, err
	}
	val2, err := evalNumber(exp[2], env)
	if err != nil {
		return ast.Null{}, err
	}
	return ast.Boolean(val1.Decimal().LessThanOrEqual(val2.Decimal())), nil
}

func _gtQ(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	if err := exactLen(exp, 3); err != nil {
		return ast.Null{}, err
	}
	val1, err := evalNumber(exp[1], env)
	if err != nil {
		return ast.Null{}, err
	}
	val2, err := evalNumber(exp[2], env)
	if err != nil {
		return ast.Null{}, err
	}
	return ast.Boolean(val1.Decimal().GreaterThan(val2.Decimal())), nil
}

func _gteqQ(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	if err := exactLen(exp, 3); err != nil {
		return ast.Null{}, err
	}
	val1, err := evalNumber(exp[1], env)
	if err != nil {
		return ast.Null{}, err
	}
	val2, err := evalNumber(exp[2], env)
	if err != nil {
		return ast.Null{}, err
	}
	return ast.Boolean(val1.Decimal().GreaterThanOrEqual(val2.Decimal())), nil
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

func _nequalQ(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	if err := exactLen(exp, 3); err != nil {
		return ast.Null{}, err
	}
	val1, err := eval.Eval(exp[1], env)
	if err != nil {
		return ast.Null{}, err
	}
	val2, err := eval.Eval(exp[2], env)
	if err != nil {
		return ast.Null{}, err
	}
	return ast.Boolean(!val1.Equal(val2)), nil
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
