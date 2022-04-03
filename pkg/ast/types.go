package ast

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// type:any
type Any interface {
	String() string
	GoString() string
	Equal(arg Any) bool
}

// type:null
type Null struct{}

func (val Null) String() string {
	return "null"
}

func (val Null) GoString() string {
	return val.String()
}

func (val Null) Equal(arg Any) bool {
	switch arg.(type) {
	default:
		return false
	case Null:
		return true
	}
}

// type:boolean
type Boolean bool

func (val Boolean) String() string {
	return val.GoString()
}

func (val Boolean) GoString() string {
	return fmt.Sprintf("%t", val)
}

func (val Boolean) Equal(arg Any) bool {
	switch arg.(type) {
	default:
		return false
	case Boolean:
		return val == arg
	}
}

// type:number
type Number decimal.Decimal

func (val Number) String() string {
	return val.GoString()
}

func (val Number) GoString() string {
	return val.Decimal().String()
}

func (val Number) Equal(arg Any) bool {
	switch num := arg.(type) {
	default:
		return false
	case Number:
		return val.Decimal().Equal(num.Decimal())
	}
}

// type:string
type String struct {
	Val string
}

func (val String) String() string {
	return val.Val
}

func (val String) GoString() string {
	return fmt.Sprintf("%#v", val.Val)
}

func (val String) Equal(arg Any) bool {
	switch arg.(type) {
	default:
		return false
	case String:
		return val == arg
	}
}

// type:array
type Array []Any

func (val Array) String() string {
	return val.GoString()
}

func (val Array) GoString() string {
	res := make([]string, len(val))
	for i, item := range val {
		res[i] = item.GoString()
	}
	return "[" + strings.Join(res, " ") + "]"
}

func (val Array) Equal(arg Any) bool {
	switch val2 := arg.(type) {
	default:
		return false
	case Array:
		if len(val) != len(val2) {
			return false
		}
		for i, item := range val {
			item2 := val2[i]
			if !item.Equal(item2) {
				return false
			}
		}
		return true
	}
}

// type:map
type Map map[String]Any

func (val Map) String() string {
	return val.GoString()
}

func (val Map) GoString() string {
	i := 0
	res := make([]string, len(val))
	for key, item := range val {
		res[i] = key.GoString() + ":" + item.GoString()
		i++
	}
	return "{" + strings.Join(res, " ") + "}"
}

func (val Map) Equal(arg Any) bool {
	switch val2 := arg.(type) {
	default:
		return false
	case Map:
		if len(val) != len(val2) {
			return false
		}
		for key, item := range val {
			item2, ok := val2[key]
			if !ok || !item.Equal(item2) {
				return false
			}
		}
		return true
	}
}

// type:symbol
type Symbol struct {
	Val string
	Pos *Position
}

type Position struct {
	Row, Column, Offset int64
}

func (val Symbol) String() string {
	return val.Val
}

func (val Symbol) GoString() string {
	if val.Pos != nil {
		return fmt.Sprintf("%s<%d,%d;%d>", val.Val, val.Pos.Row, val.Pos.Column, val.Pos.Offset)
	}
	return fmt.Sprintf("%s<?>", val.Val)
}

func (val Symbol) Equal(arg Any) bool {
	switch sym := arg.(type) {
	default:
		return false
	case Symbol:
		return val == sym
	}
}

// type:expression
type Expr []Any

func (val Expr) String() string {
	return val.GoString()
}

func (val Expr) GoString() string {
	res := make([]string, len(val))
	for i, item := range val {
		res[i] = item.GoString()
	}
	return "(" + strings.Join(res, " ") + ")"
}

func (val Expr) Equal(arg Any) bool {
	switch val2 := arg.(type) {
	default:
		return false
	case Expr:
		if len(val) != len(val2) {
			return false
		}
		for i, item := range val {
			item2 := val2[i]
			if !item.Equal(item2) {
				return false
			}
		}
		return true
	}
}

// type:operator
type Operator struct {
	Val string
}

func (val Operator) String() string {
	return val.GoString()
}

func (val Operator) GoString() string {
	return val.Val
}

func (val Operator) Equal(arg Any) bool {
	switch op := arg.(type) {
	default:
		return false
	case Operator:
		return val.Val == op.Val
	}
}
