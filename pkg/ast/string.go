package ast

import (
	"strconv"
	"strings"
)

// parse quoted, escaped string
func NewStringFromString(val string) (String, error) {
	str := strings.Replace(val, "\\/", "/", -1)
	str, err := strconv.Unquote(str)
	return String{Val: str}, err
}
