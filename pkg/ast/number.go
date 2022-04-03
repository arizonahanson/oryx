package ast

import (
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
