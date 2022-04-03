package ast

import (
	"github.com/shopspring/decimal"
)

var Zero = NewNumber(0)
var One = NewNumber(1)

// number from int
func NewNumber(num int64) Number {
	return Number(decimal.NewFromInt(num))
}

// parse number from string
func NewNumberFromString(num string) (Number, error) {
	dec, err := decimal.NewFromString(num)
	return Number(dec), err
}

// number to shopspring decimal
func (val Number) Decimal() decimal.Decimal {
	return decimal.Decimal(val)
}
