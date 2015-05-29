//line expr.y:17
package main

import __yyfmt__ "fmt"

//line expr.y:18
import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"unicode/utf8"
)

var result_value int
var result_error string

//line expr.y:33
type exprSymType struct {
	yys int
	num int
}

const NUM = 57346

var exprToknames = []string{
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'('",
	"')'",
	"NUM",
}
var exprStatenames = []string{}

const exprEofCode = 1
const exprErrCode = 2
const exprMaxDepth = 200

//line expr.y:92

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
L:
	for {
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

//line yacctab:1
var exprExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const exprNprod = 13
const exprPrivate = 57344

var exprTokenNames []string
var exprStates []string

const exprLast = 23

var exprAct = []int{

	7, 4, 5, 2, 21, 9, 6, 8, 12, 13,
	9, 1, 8, 16, 3, 19, 20, 17, 18, 14,
	15, 10, 11,
}
var exprPact = []int{

	-3, -1000, -1000, 17, -3, -3, 13, -1000, -1000, -3,
	2, 2, -1000, -1000, 2, 2, -5, 13, 13, -1000,
	-1000, -1000,
}
var exprPgo = []int{

	0, 3, 14, 6, 0, 11,
}
var exprR1 = []int{

	0, 5, 1, 1, 1, 2, 2, 2, 3, 3,
	3, 4, 4,
}
var exprR2 = []int{

	0, 1, 1, 2, 2, 1, 3, 3, 1, 3,
	3, 1, 3,
}
var exprChk = []int{

	-1000, -5, -1, -2, 4, 5, -3, -4, 10, 8,
	4, 5, -1, -1, 6, 7, -1, -3, -3, -4,
	-4, 9,
}
var exprDef = []int{

	0, -2, 1, 2, 0, 0, 5, 8, 11, 0,
	0, 0, 3, 4, 0, 0, 0, 6, 7, 9,
	10, 12,
}
var exprTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	8, 9, 6, 4, 3, 5, 3, 7,
}
var exprTok2 = []int{

	2, 3, 10,
}
var exprTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var exprDebug = 0

type exprLexer interface {
	Lex(lval *exprSymType) int
	Error(s string)
}

const exprFlag = -1000

func exprTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(exprToknames) {
		if exprToknames[c-4] != "" {
			return exprToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func exprStatname(s int) string {
	if s >= 0 && s < len(exprStatenames) {
		if exprStatenames[s] != "" {
			return exprStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func exprlex1(lex exprLexer, lval *exprSymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = exprTok1[0]
		goto out
	}
	if char < len(exprTok1) {
		c = exprTok1[char]
		goto out
	}
	if char >= exprPrivate {
		if char < exprPrivate+len(exprTok2) {
			c = exprTok2[char-exprPrivate]
			goto out
		}
	}
	for i := 0; i < len(exprTok3); i += 2 {
		c = exprTok3[i+0]
		if c == char {
			c = exprTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = exprTok2[1] /* unknown char */
	}
	if exprDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", exprTokname(c), uint(char))
	}
	return c
}

func exprParse(exprlex exprLexer) int {
	var exprn int
	var exprlval exprSymType
	var exprVAL exprSymType
	exprS := make([]exprSymType, exprMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	exprstate := 0
	exprchar := -1
	exprp := -1
	goto exprstack

ret0:
	return 0

ret1:
	return 1

exprstack:
	/* put a state and value onto the stack */
	if exprDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", exprTokname(exprchar), exprStatname(exprstate))
	}

	exprp++
	if exprp >= len(exprS) {
		nyys := make([]exprSymType, len(exprS)*2)
		copy(nyys, exprS)
		exprS = nyys
	}
	exprS[exprp] = exprVAL
	exprS[exprp].yys = exprstate

exprnewstate:
	exprn = exprPact[exprstate]
	if exprn <= exprFlag {
		goto exprdefault /* simple state */
	}
	if exprchar < 0 {
		exprchar = exprlex1(exprlex, &exprlval)
	}
	exprn += exprchar
	if exprn < 0 || exprn >= exprLast {
		goto exprdefault
	}
	exprn = exprAct[exprn]
	if exprChk[exprn] == exprchar { /* valid shift */
		exprchar = -1
		exprVAL = exprlval
		exprstate = exprn
		if Errflag > 0 {
			Errflag--
		}
		goto exprstack
	}

exprdefault:
	/* default state action */
	exprn = exprDef[exprstate]
	if exprn == -2 {
		if exprchar < 0 {
			exprchar = exprlex1(exprlex, &exprlval)
		}

		/* look through exception table */
		xi := 0
		for {
			if exprExca[xi+0] == -1 && exprExca[xi+1] == exprstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			exprn = exprExca[xi+0]
			if exprn < 0 || exprn == exprchar {
				break
			}
		}
		exprn = exprExca[xi+1]
		if exprn < 0 {
			goto ret0
		}
	}
	if exprn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			exprlex.Error("syntax error")
			Nerrs++
			if exprDebug >= 1 {
				__yyfmt__.Printf("%s", exprStatname(exprstate))
				__yyfmt__.Printf(" saw %s\n", exprTokname(exprchar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for exprp >= 0 {
				exprn = exprPact[exprS[exprp].yys] + exprErrCode
				if exprn >= 0 && exprn < exprLast {
					exprstate = exprAct[exprn] /* simulate a shift of "error" */
					if exprChk[exprstate] == exprErrCode {
						goto exprstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if exprDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", exprS[exprp].yys)
				}
				exprp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if exprDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", exprTokname(exprchar))
			}
			if exprchar == exprEofCode {
				goto ret1
			}
			exprchar = -1
			goto exprnewstate /* try again in the same state */
		}
	}

	/* reduction by production exprn */
	if exprDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", exprn, exprStatname(exprstate))
	}

	exprnt := exprn
	exprpt := exprp
	_ = exprpt // guard against "declared and not used"

	exprp -= exprR2[exprn]
	exprVAL = exprS[exprp+1]

	/* consult goto table to find next state */
	exprn = exprR1[exprn]
	exprg := exprPgo[exprn]
	exprj := exprg + exprS[exprp].yys + 1

	if exprj >= exprLast {
		exprstate = exprAct[exprg]
	} else {
		exprstate = exprAct[exprj]
		if exprChk[exprstate] != -exprn {
			exprstate = exprAct[exprg]
		}
	}
	// dummy call; replaced with literal code
	switch exprnt {

	case 1:
		//line expr.y:47
		{
			result_value = exprS[exprpt-0].num
		}
	case 2:
		exprVAL.num = exprS[exprpt-0].num
	case 3:
		//line expr.y:54
		{
			exprVAL.num = exprS[exprpt-0].num
		}
	case 4:
		//line expr.y:58
		{
			exprVAL.num = -exprS[exprpt-0].num
		}
	case 5:
		exprVAL.num = exprS[exprpt-0].num
	case 6:
		//line expr.y:65
		{
			exprVAL.num = exprS[exprpt-2].num + exprS[exprpt-0].num
		}
	case 7:
		//line expr.y:69
		{
			exprVAL.num = exprS[exprpt-2].num - exprS[exprpt-0].num
		}
	case 8:
		exprVAL.num = exprS[exprpt-0].num
	case 9:
		//line expr.y:76
		{
			exprVAL.num = exprS[exprpt-2].num * exprS[exprpt-0].num
		}
	case 10:
		//line expr.y:80
		{
			exprVAL.num = exprS[exprpt-2].num / exprS[exprpt-0].num
		}
	case 11:
		exprVAL.num = exprS[exprpt-0].num
	case 12:
		//line expr.y:87
		{
			exprVAL.num = exprS[exprpt-1].num
		}
	}
	goto exprstack /* stack new state and value */
}
