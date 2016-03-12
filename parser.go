package tags

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type Param struct {
	Name   string
	Values []string
}

type Parser interface {
	Parse(tag string) ([]Param, error)
}

type parser struct {
	ignoreKeyword string
}

var defaultParser = NewParser("-")

// NewParser creates a new parser with a custom ignoreKeywords
func NewParser(ignoreKeyword string) Parser {
	return &parser{
		ignoreKeyword: ignoreKeyword,
	}
}

// Parse parses the provided tag and returns a slice with the found parameter (keywords/value) pairs.
// If the provided tag is empty or has the ignore identifier "-" it will return a nil slice
// The default default tag parser is used.
func Parse(tag string) ([]Param, error) {
	return defaultParser.Parse(tag)
}

// Parse parses the provided tag and returns a slice with the found parameter (keywords/value) pairs.
// If the provided tag is empty or has the ignore identifier (custom defined) it will return a nil slice
func (p *parser) Parse(tag string) ([]Param, error) {
	tag = strings.TrimSpace(tag)
	if len(tag) == 0 || tag == p.ignoreKeyword {
		return nil, nil
	}

	var params []Param = nil
	var s scanner.Scanner
	s.Init(strings.NewReader(tag))
	s.Error = func(_ *scanner.Scanner, _ string) {}

	var param *Param = nil
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if param != nil && s.TokenText() == "(" {
			//scan enclosure
			for tok != scanner.EOF {
				tok = s.Scan()
				if s.TokenText() == ")" {
					//end of enclosure
					break
				} else if tok == scanner.EOF {
					return nil, errors.New("Unexpected EOF")
				} else if s.TokenText() == "," {
					//ignore it
					continue
				} else if s.TokenText() == ";" {
					return nil, fmt.Errorf("Unexpected token `%s` expected a , or )", s.TokenText())
				} else {
					value, err := strconv.Unquote(s.TokenText())
					if err != nil {
						value = s.TokenText()
					}
					param.Values = append(param.Values, value)
				}
			}
		} else if s.TokenText() == ";" || s.TokenText() == "," {
			params = append(params, *param)
			param = nil
			continue
		} else if param == nil {
			param = &Param{
				Name:   s.TokenText(),
				Values: nil,
			}
		} else {
			return nil, fmt.Errorf("Unexpected token `%s` expected a delimter ; or ,", s.TokenText())
		}
	}

	if param != nil {
		return append(params, *param), nil
	}

	return params, nil
}
