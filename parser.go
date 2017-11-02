package tags

import (
	"errors"
	"fmt"
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
	//Parse will parse a tag and returns the params in order
	Parse(tag string) ([]Param, error)
	//ParseMao will parse a tag and will return an unordered map
	ParseMap(tag string) (map[string][]string, error)
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
// If the provided tag has the ignore identifier ("-") it will return a nil slice
// If the provided tag is empty it will return a empty slice
func Parse(tag string) ([]Param, error) {
	return defaultParser.Parse(tag)
}

// ParseMap parses the provided tag and returns a map with the found parameters.
// If the provided tag has the ignore identifier ("-") it will return a nil map
// If the provided tag is empty it will return an empty map
func ParseMap(tag string) (map[string][]string, error) {
	return defaultParser.ParseMap(tag)
}

// Parse parses the provided tag and returns a slice with the found parameter (keywords/value) pairs.
// If the provided tag has the ignore identifier (custom defined) it will return a nil slice
// If the provided tag is empty it will return a empty slice
func (p *parser) Parse(tag string) (params []Param, err error) {
	tag = strings.TrimSpace(tag)
	if tag == p.ignoreKeyword {
		return nil, nil
	}

	params = []Param{}
	paramReadFunc := func(name string, args []string) {
		param := Param{
			Name: name,
			Args: args,
		}
		params = append(params, param)
	}
	if err := parseTag(tag, paramReadFunc); err != nil {
		return nil, err
	}
	return params, nil
}

// ParseMap parses the provided tag and returns a map with the found parameters.
// If the provided tag has the ignore identifier (custom defined) it will return a nil map
// If the provided tag is empty it will return an empty map
func (p *parser) ParseMap(tag string) (map[string][]string, error) {
	tag = strings.TrimSpace(tag)
	if tag == p.ignoreKeyword {
		return nil, nil
	}

	params := make(map[string][]string)
	paramReadFunc := func(name string, args []string) {
		params[name] = args
	}
	if err := parseTag(tag, paramReadFunc); err != nil {
		return nil, err
	}
	return params, nil
}

type paramReadFunc func(name string, args []string)

func parseTag(tag string, paramReadFn paramReadFunc) error {
	var s scanner.Scanner
	s.Init(strings.NewReader(tag))
	s.Error = func(_ *scanner.Scanner, _ string) {}

	var (
		err  error
		name string
		args []string
	)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if name != "" && s.TokenText() == "(" {
			args, err = parseArguments(&s)
			if err != nil {
				return err
			}
		} else if s.TokenText() == ";" || s.TokenText() == "," {
			paramReadFn(name, args)
			name = ""
			args = []string{}
		} else if name == "" {
			name = s.TokenText()
		} else {
			return fmt.Errorf("Unexpected token `%s` expected a delimter ; or ,", s.TokenText())
		}
	}

	if name != "" {
		paramReadFn(name, args)
	}
	return nil
}

func parseArguments(s *scanner.Scanner) ([]string, error) {
	var args []string
	for {
		tok := s.Scan()
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
			args = append(args, unescape(s.TokenText()))
		}

	}
	return args, nil
}

// unescape only removes the escape character from a
// - double escaped backslash (\\) to (\)
// - escaped quotes (\") to (")
func unescape(in string) string {

	in = strings.TrimSpace(in)
	l := len(in)
	if l < 2 {
		return in
	}

	if in[0] != '"' || in[l-1] != '"' {
		return in
	}

	in = strings.Replace(in[1:l-1], `\\`, `\`, -1)
	return strings.Replace(in, `\"`, `"`, -1)
}
