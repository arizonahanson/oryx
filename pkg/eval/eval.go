package eval

import "github.com/arizonahanson/oryx/pkg/ast"

func EvalBytes(bytes []byte, env *Env) (ast.Any, error) {
	arg, err := Parse(bytes)
	if err != nil {
		return ast.Null{}, err
	}
	return Eval(arg, env)
}

func EvalFile(filename string, env *Env) (ast.Any, error) {
	arg, err := ParseFile(filename)
	if err != nil {
		return ast.Null{}, err
	}
	return Eval(arg, env)
}

// eager evaluation
func Eval(any ast.Any, env *Env) (ast.Any, error) {
	return FutureEval(any, env).Get()
}

// lazy evaluation
func FutureEval(any ast.Any, env *Env) Future {
	return func() (val ast.Any, err error) {
		return eval(any, env)
	}
}

func eval(any ast.Any, env *Env) (val ast.Any, err error) {
	val = any
	switch arg := val.(type) {
	default:
		break
	case ast.Symbol:
		// symbol
		val, err = env.Get(arg)
	case ast.Array:
		// array
		res := make(ast.Array, len(arg))
		for i, item := range arg {
			res[i], err = Eval(item, env)
			if err != nil {
				return
			}
		}
		val = res
	case ast.Map:
		// map
		res := make(ast.Map, len(arg))
		for key, item := range arg {
			res[key], err = Eval(item, env)
			if err != nil {
				return
			}
		}
		val = res
	case ast.Expr:
		// expression
		if len(arg) == 0 {
			// eval to null
			return ast.Null{}, nil
		}
		// pre-eval first item
		var first ast.Any
		first, err = Eval(arg[0], env)
		if err != nil {
			return
		}
		switch fn := first.(type) {
		case Func:
			// eval to function-call future
			val = fn.Future(arg, env)
		default:
			// eval to array
			res := make(ast.Array, len(arg))
			for i, item := range arg {
				if i == 0 {
					res[0] = first
					continue
				}
				res[i], err = Eval(item, env)
				if err != nil {
					return
				}
			}
			val = res
		}
	}
	return
}
