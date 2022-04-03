package ast

import (
	"strconv"
	"strings"
)

// parse quoted string with escape-sequences
func NewStringFromString(val string) (String, error) {
	str := strings.Replace(val, "\\/", "/", -1)
	str, err := strconv.Unquote(str)
	return String{Val: str}, err
}
