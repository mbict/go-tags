package tags_test

import (
	"github.com/mbict/go-tags"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type ParserSuite struct{}

var _ = Suite(&ParserSuite{})

var goodPathTests = []struct {
	Description string
	Tag         string
	Expected    []tags.Param
}{
	{
		Description: "empty tag",
		Tag:         ``,
		Expected:    []tags.Param{},
	}, {
		Description: "ignored tag",
		Tag:         `-`,
		Expected:    nil,
	}, {
		Description: "only keyword",
		Tag:         `max`,
		Expected:    []tags.Param{tags.Param{Name: "max", Args: nil}},
	}, {
		Description: "only keyword with whitespace",
		Tag:         `    max     `,
		Expected:    []tags.Param{tags.Param{Name: "max", Args: nil}},
	}, {
		Description: "keyword with parentesis but no params",
		Tag:         `max()`,
		Expected:    []tags.Param{tags.Param{Name: "max", Args: nil}},
	}, {
		Description: "keyword with parentesis but no params filled with whitespace",
		Tag:         `max(      )`,
		Expected:    []tags.Param{tags.Param{Name: "max", Args: nil}},
	}, {
		Description: "keyword with parentesis with close semicolon",
		Tag:         `max();`,
		Expected:    []tags.Param{tags.Param{Name: "max", Args: nil}},
	}, {
		Description: "keyword with 1 param",
		Tag:         `max(123)`,
		Expected:    []tags.Param{tags.Param{Name: "max", Args: []string{"123"}}},
	}, {
		Description: "keyword with 1 param is enclosed with whitespace",
		Tag:         `max(    123   )`,
		Expected:    []tags.Param{tags.Param{Name: "max", Args: []string{"123"}}},
	}, {
		Description: "keyword with 3 params",
		Tag:         `max(123,456,7)`,
		Expected:    []tags.Param{tags.Param{Name: "max", Args: []string{"123", "456", "7"}}},
	}, {
		Description: "keyword with 1 param enclosed with quotes and has whitespace",
		Tag:         `max("a a")`,
		Expected:    []tags.Param{tags.Param{Name: "max", Args: []string{"a a"}}},
	}, {
		Description: "keyword with 3 params enclosed with quotes and has whitespace",
		Tag:         `max("a","bc d"," ef gg ")`,
		Expected:    []tags.Param{tags.Param{Name: "max", Args: []string{"a", "bc d", " ef gg "}}},
	}, {
		Description: "keyword with 1 params enclosed with quotes and a escaped quote",
		Tag:         `max("a\"a")`,
		Expected:    []tags.Param{tags.Param{Name: "max", Args: []string{"a\"a"}}},
	}, {
		Description: "4 keywords with various methods separated with semicolon",
		Tag:         `max;min();in("a","b","c");between(1,2,3)`,
		Expected: []tags.Param{
			tags.Param{Name: "max"},
			tags.Param{Name: "min"},
			tags.Param{Name: "in", Args: []string{"a", "b", "c"}},
			tags.Param{Name: "between", Args: []string{"1", "2", "3"}},
		},
	}, {
		Description: "4 keywords with various methods separated with comma",
		Tag:         `max,min(),in("a","b","c"),between(1,2,3)`,
		Expected: []tags.Param{
			tags.Param{Name: "max"},
			tags.Param{Name: "min"},
			tags.Param{Name: "in", Args: []string{"a", "b", "c"}},
			tags.Param{Name: "between", Args: []string{"1", "2", "3"}},
		},
	}, {
		Description: "4 keywords with various methods separated with commas and semicolons",
		Tag:         `max;min(),in("a","b","c");between(1,2,3)`,
		Expected: []tags.Param{
			tags.Param{Name: "max"},
			tags.Param{Name: "min"},
			tags.Param{Name: "in", Args: []string{"a", "b", "c"}},
			tags.Param{Name: "between", Args: []string{"1", "2", "3"}},
		},
	}, {
		Description: "escape test",
		Tag:         `in("\\","\"","\d\"\t")`,
		Expected: []tags.Param{
			tags.Param{Name: "in", Args: []string{`\`, `"`, `\d"\t`}},
		},
	},
}

func (ms *ParserSuite) TestParse(c *C) {
	for _, test := range goodPathTests {
		result, err := tags.Parse(test.Tag)

		c.Assert(err, IsNil, Commentf("test failed gives error for `%s`", test.Description))
		c.Assert(result, DeepEquals, test.Expected, Commentf("test failed for `%s`", test.Description))
	}
}

func (ms *ParserSuite) TestParseMap(c *C) {
	for _, test := range goodPathTests {
		result, err := tags.ParseMap(test.Tag)

		c.Assert(err, IsNil, Commentf("test failed gives error for `%s`", test.Description))
		if test.Expected == nil {
			c.Assert(result, IsNil, Commentf("test failed for `%s`", test.Description))
		} else {
			c.Assert(result, HasLen, len(test.Expected), Commentf("test failed for `%s`", test.Description))
			for _, param := range test.Expected {
				v, hasKey := result[param.Name]
				c.Assert(hasKey, Equals, true, Commentf("test failed for `%s`", test.Description))
				c.Assert(v, DeepEquals, param.Args, Commentf("test failed for `%s`", test.Description))
			}
		}

	}
}

func (ms *ParserSuite) TestCustomParser(c *C) {
	parser := tags.NewParser("-")
	for _, test := range goodPathTests {
		result, err := parser.Parse(test.Tag)

		c.Assert(err, IsNil, Commentf("test failed gives error for `%s`", test.Description))
		c.Assert(result, DeepEquals, test.Expected, Commentf("test failed for `%s`", test.Description))
	}
}

var badPathTests = []struct {
	Description string
	Tag         string
	Expected    string
}{
	{
		Description: "EOF too early",
		Tag:         "min(",
		Expected:    "Unexpected EOF",
	}, {
		Description: "EOF too early",
		Tag:         "min( 1234",
		Expected:    "Unexpected EOF",
	}, {
		Description: "Not closed quoted string",
		Tag:         `min("abbcbcb);test`,
		Expected:    "Unexpected EOF",
	}, {
		Description: "Closing a keyword with parenthesis without a opening",
		Tag:         "min)",
		Expected:    "Unexpected token `\\)` expected a delimter \\; or \\,",
	}, {
		Description: "Unexpected token when the current keyword is not closed with a delimiter",
		Tag:         "min()max",
		Expected:    "Unexpected token `max` expected a delimter \\; or \\,",
	}, {
		Description: "Unexpected token when the current opened parenthesis is not closed",
		Tag:         "min(;",
		Expected:    "Unexpected token `\\;` expected a \\, or \\)",
	},
}

func (ms *ParserSuite) TestParseFailingPatterns(c *C) {
	for _, test := range badPathTests {
		result, err := tags.Parse(test.Tag)

		c.Assert(result, IsNil, Commentf("test failed should not give results for test `%s`", test.Description))
		c.Assert(err, ErrorMatches, test.Expected)
	}
}

func (ms *ParserSuite) TestParseMapFailingPatterns(c *C) {
	for _, test := range badPathTests {
		result, err := tags.ParseMap(test.Tag)

		c.Assert(result, IsNil, Commentf("test failed should not give results for test `%s`", test.Description))
		c.Assert(err, ErrorMatches, test.Expected)
	}
}
