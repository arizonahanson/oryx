package parser

import (
	"strconv"
	"strings"

	"github.com/arizonahanson/oryx/pkg/ast"
	"github.com/shopspring/decimal"
)

// parse quoted, escaped string
func NewStringFromString(val string) (ast.String, error) {
	str := strings.Replace(val, "\\/", "/", -1)
	str, err := strconv.Unquote(str)
	return ast.String{Val: str}, err
}

// parse number
func NewNumberFromString(num string) (ast.Number, error) {
	dec, err := decimal.NewFromString(num)
	return ast.Number(dec), err
}
