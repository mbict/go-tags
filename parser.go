package tags

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

// Param stores a keyword and the arguments extracted from a parsed tag
type Param struct {
	Name string
	Args []string
}

// Parser interface only exports parse interface
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
func (p *parser) Parse(tag string) (params []Param, err error) {
	tag = strings.TrimSpace(tag)
	if len(tag) == 0 || tag == p.ignoreKeyword {
		return nil, nil
	}

	var s scanner.Scanner
	s.Init(strings.NewReader(tag))
	s.Error = func(_ *scanner.Scanner, _ string) {}

	var param *Param
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if param != nil && s.TokenText() == "(" {
			args, err := parseArguments(&s)
			if err != nil {
				return nil, err
			}
			param.Args = args

		} else if s.TokenText() == ";" || s.TokenText() == "," {
			params = append(params, *param)
			param = nil
		} else if param == nil {
			param = &Param{
				Name: s.TokenText(),
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

func parseArguments(s *scanner.Scanner) ([]string, error) {
	var args []string
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {

		if s.TokenText() == ")" {
			//end of enclosure
			break
		} else if tok == scanner.EOF {
			return nil, errors.New("Unexpected EOF")
		} else if s.TokenText() == "," {
			//ignore it
		} else if s.TokenText() == ";" {
			return nil, fmt.Errorf("Unexpected token `%s` expected a , or )", s.TokenText())
		} else {
			if value, err := strconv.Unquote(s.TokenText()); err == nil {
				args = append(args, value)
			} else {
				args = append(args, s.TokenText())
			}
		}
	}
	return args, nil
}
