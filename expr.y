// Modifications copyright 2015 Evan Miller

// Original copyright statement:
//
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is an example of a goyacc program.
// To build it:
// go tool yacc -p "expr" expr.y (produces y.go)
// go build -o expr y.go
// expr
// > <type an expression>

%{

package main

import (
	"bytes"
    "errors"
	"fmt"
	"log"
	"unicode/utf8"
)

var result_value int
var result_error string

%}

%union {
	num int
}

%type	<num>	expr expr1 expr2 expr3

%token '+' '-' '*' '/' '(' ')'

%token	<num>	NUM

%%

top:
	expr
	{
        result_value = $1
	}

expr:
	expr1
|	'+' expr
	{
		$$ = $2
	}
|	'-' expr
	{
		$$ = -$2
	}

expr1:
	expr2
|	expr1 '+' expr2
	{
		$$ = $1 + $3
	}
|	expr1 '-' expr2
	{
		$$ = $1 - $3
	}

expr2:
	expr3
|	expr2 '*' expr3
	{
		$$ = $1 * $3
	}
|	expr2 '/' expr3
	{
		$$ = $1 / $3
	}

expr3:
	NUM
|	'(' expr ')'
	{
		$$ = $2
	}


%%

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type exprLex struct {
	line []byte
	peek rune
}

// The parser calls this method to get each new token.  This
// implementation returns operators and NUM.
func (x *exprLex) Lex(yylval *exprSymType) int {
	for {
		c := x.next()
		switch c {
		case eof:
			return eof
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return x.num(c, yylval)
		case '+', '-', '*', '/', '(', ')':
			return int(c)

		case ' ', '\t', '\n', '\r':
		default:
			log.Printf("unrecognized character %q", c)
		}
	}
}

// Lex a number.
func (x *exprLex) num(c rune, yylval *exprSymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			log.Fatalf("WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	add(&b, c)
	L: for {
		c = x.next()
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'x':
			add(&b, c)
		default:
			break L
		}
	}
	if c != eof {
		x.peek = c
	}
    value := 0
	if n, _ := fmt.Sscanf(b.String(), "%v", &value); n != 1 {
		log.Printf("bad number %q", b.String())
		return eof
	}
    yylval.num = value
	return NUM
}

// Return the next rune for the lexer.
func (x *exprLex) next() rune {
	if x.peek != eof {
		r := x.peek
		x.peek = eof
		return r
	}
	if len(x.line) == 0 {
		return eof
	}
	c, size := utf8.DecodeRune(x.line)
	x.line = x.line[size:]
	if c == utf8.RuneError && size == 1 {
		log.Print("invalid utf8")
		return x.next()
	}
	return c
}

// The parser calls this method on a parse error.
func (x *exprLex) Error(s string) {
    result_error = s
}

func evaluateExpression(line string) (int, error) {
    result_value = -1
    result_error = ""
    exprParse(&exprLex{line: []byte(line)})
    if result_error == "" {
        return result_value, nil
    }
    return result_value, errors.New(result_error)
}
