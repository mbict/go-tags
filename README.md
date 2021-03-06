[![wercker status](https://app.wercker.com/status/252a744c512d2eb933d775039947622b/s "wercker status")](https://app.wercker.com/project/bykey/252a744c512d2eb933d775039947622b)
[![Build Status](https://travis-ci.org/mbict/go-tags.png?branch=master)](https://travis-ci.org/mbict/go-tags)
[![GoDoc](https://godoc.org/github.com/mbict/go-tags?status.png)](http://godoc.org/github.com/mbict/go-tags)
[![GoCover](http://gocover.io/_badge/github.com/mbict/go-tags)](http://gocover.io/github.com/mbict/go-tags)
[![GoReportCard](http://goreportcard.com/badge/mbict/go-tags)](http://goreportcard.com/report/mbict/go-tags)

Tags
====

Tags is a simple string parser used for extracting keyword/arguments pairs from structure tags.

Why
===
I was recreating this feature for a few libraries i'm maintaining and developing (such as validation, binding and a ORM). 
So i felt the need to create a common library to take care of this job and keeps me stopping repeating my self ;).

Examples
========
Tags can parse a few formats

#### Single keyword tag
```go
tags.Parse(`required`)

// outputs map
// []Param{ Param{ Name: "required", Args:nil } 
```

#### Multiple keywords
You can delimit keywords with a comma `,` or a semicolon `;`
```go
tags.Parse(`required;email`)
tags.Parse(`required,email`)

// outputs map
// []Param{ 
//    Param{ Name: "required", Args:nil },
//    Param{ Name: "email", Args:nil },
// }
```

#### Keywords with params
You can delimit keywords with a comma `,` or a semicolon `;`
```go
tags.Parse(`between(10,20);in("foo", "bar", "foo bar")`)

// outputs map
// []Param{ 
//    Param{ Name: "between", Args:[]string{"10", "20"} },
//    Param{ Name: "in", Args:[]string{"foo", "bar", "foo bar"} },
// }
```

#### Ignored keywords
You can delimit keywords with a comma `,` or a semicolon `;`
```go
tags.Parse(`-`)

// gives nil pointer back
```

Error handling
==============
When a malformed tag is provided the parse function will return a error with the reason why.
This could be a unexpected end of string (EOF) or badly closed statements.