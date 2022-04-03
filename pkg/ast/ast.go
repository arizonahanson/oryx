package ast

import (
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// useful numbers
var Zero = NewNumber(0)
var One = NewNumber(1)

// number from int
func NewNumber(num int64) Number {
	return Number(decimal.NewFromInt(num))
}

// number to decimal impl
func (val Number) Decimal() decimal.Decimal {
	return decimal.Decimal(val)
}

// parse quoted, escaped string
func (val String) Unquote() (String, error) {
	str := strings.Replace(val.Val, "\\/", "/", -1)
	str, err := strconv.Unquote(str)
	return String{str}, err
}

// parse number
func (val String) Number() (Number, error) {
	dec, err := decimal.NewFromString(val.Val)
	return Number(dec), err
}

func NewSymbol(sym string, pos *Position) Symbol {
	return Symbol{Val: sym, Pos: pos}
}
