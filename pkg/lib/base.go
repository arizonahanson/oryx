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
	"+":      _add,
	"add":    _add,
	"-":      _sub,
	"sub":    _sub,
	"*":      _mul,
	"mul":    _mul,
	"/":      _div,
	"div":    _div,
	"quo":    _quo,
	"rem":    _rem,
	"!":      _not,
	"not":    _not,
	":=":     _defE,
	"def!":   _defE,
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

func oneArg(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	if err := exactLen(exp, 2); err != nil {
		return ast.Null{}, err
	}
	val, err := eval.Eval(exp[1], env)
	if err != nil {
		return ast.Null{}, err
	}
	return val, nil
}

func _defE(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	if err := exactLen(exp, 3); err != nil {
		return ast.Null{}, err
	}
	switch sym := exp[1].(type) {
	default:
		return ast.Null{}, fmt.Errorf("called with non-symbol %#v", exp[1])
	case ast.Symbol:
		env.Set(sym, eval.FutureEval(exp[2], env))
		return ast.Null{}, nil
	}
}

func _add(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	res := ast.Zero.Decimal()
	for _, item := range exp[1:] {
		val, err := evalNumber(item, env)
		if err != nil {
			return ast.Null{}, err
		}
		res = res.Add(val.Decimal())
	}
	return ast.Number(res), nil
}

func _sub(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	res := ast.Zero.Decimal()
	for i, item := range exp[1:] {
		val, err := evalNumber(item, env)
		if err != nil {
			return ast.Null{}, err
		}
		if i == 0 && len(exp) > 2 {
			res = val.Decimal()
		} else {
			res = res.Sub(val.Decimal())
		}
	}
	return ast.Number(res), nil
}

func _mul(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	res := ast.One.Decimal()
	for _, item := range exp[1:] {
		val, err := evalNumber(item, env)
		if err != nil {
			return ast.Null{}, err
		}
		res = res.Mul(val.Decimal())
	}
	return ast.Number(res), nil
}

func _quo(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	if err := exactLen(exp, 4); err != nil {
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
	val3, err := evalNumber(exp[3], env)
	if err != nil {
		return ast.Null{}, err
	}
	q, _ := val1.Decimal().QuoRem(val2.Decimal(), int32(val3.Decimal().IntPart()))
	return ast.Number(q), nil
}

func _rem(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	if err := exactLen(exp, 4); err != nil {
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
	val3, err := evalNumber(exp[3], env)
	if err != nil {
		return ast.Null{}, err
	}
	_, r := val1.Decimal().QuoRem(val2.Decimal(), int32(val3.Decimal().IntPart()))
	return ast.Number(r), nil
}

func _div(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	res := ast.One.Decimal()
	for i, item := range exp[1:] {
		val, err := evalNumber(item, env)
		if err != nil {
			return ast.Null{}, err
		}
		if i == 0 && len(exp) > 2 {
			res = val.Decimal()
		} else {
			res = res.Div(val.Decimal())
		}
	}
	return ast.Number(res), nil
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
	return eval.FutureEval(exp[len(exp)-1], env), nil
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
	return eval.FutureEval(exp[len(exp)-1], env), nil
}

func _not(exp ast.Expr, env *eval.Env) (ast.Any, error) {
	val, err := oneArg(exp, env)
	if err != nil {
		return ast.Null{}, err
	}
	return ast.Boolean(val == ast.Boolean(false) || val == ast.Null{}), nil
}
