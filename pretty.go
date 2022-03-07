package pp

import (
	"bytes"
	"fmt"
	"go/format"
	"go/scanner"
	"go/token"
)

func Sprint(v ...interface{}) string {
	var s string
	if len(v) == 1 {
		s = fmt.Sprintf("%#v", v[0])
	} else {
		s = fmt.Sprintf("%#v", v)
	}
	if debug {
		fmt.Println(s)
	}
	var sc scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(s))
	sc.Init(file, []byte(s), nil, 0)
	_, tok, lit := sc.Scan()
	p := newPrinter(len(s))
	for sc.ErrorCount == 0 && p.print(tok, lit) {
		_, tok, lit = sc.Scan()
	}
	if debug {
		fmt.Println(p.String())
	}
	return frmt(p.Bytes(), s)
}

// frmt attempts to format a or b. If both fails, b is returned as it was.
func frmt(a []byte, b string) string {
	fs, err := format.Source(a)
	if err != nil {
		if debug {
			fmt.Printf("error formatting reprinted source: %v\n", err)
		}
		fs, err = format.Source([]byte(b))
		if err != nil {
			fmt.Printf("error formatting %s source: %v\n", "%#v", err)
			return b
		}
	}
	return string(fs)
}

type printer struct {
	*bytes.Buffer
	scope []token.Token
	prev  token.Token
}

func newPrinter(sz int) *printer {
	p := &printer{
		Buffer: &bytes.Buffer{},
	}
	p.Grow(sz * 4 / 3)
	return p
}

func (p *printer) print(tok token.Token, lit string) bool {
	if tok == token.EOF || tok == token.ILLEGAL {
		return false
	}
	if debug {
		fmt.Printf("%v %v %v\t%-8q\t%q\n", tok.IsKeyword(), tok.IsLiteral(), tok.IsOperator(), tok.String(), lit)
	}
	p.trackScope(tok)
	p.printBefore(tok)

	// write the actual token
	if tok.IsOperator() || tok.IsKeyword() {
		p.WriteString(tok.String())
	} else {
		p.WriteString(lit)
	}

	p.printAfter(tok)
	p.prev = tok
	return true
}

func (p *printer) trackScope(tok token.Token) {
	switch tok {
	case token.LBRACE, token.LPAREN: // push {(
		p.scope = append(p.scope, tok)
	case token.RBRACE, token.RPAREN: // pop on })
		p.scope = p.scope[:len(p.scope)-1]
	}
}

func (p *printer) printBefore(tok token.Token) {
	// make } ",\n}" except when prev token was ident (struct def doesn't take comma) or {
	if tok == token.RBRACE && p.prev != token.IDENT && p.prev != token.LBRACE {
		p.WriteString(",\n")
	}
	// put space between struct field name and type
	if tok == token.IDENT && p.prev == token.IDENT {
		p.WriteByte(' ')
	}
}

func (p *printer) printAfter(tok token.Token) {
	// put a newline after } or , except inside PAREN () scope
	if tok == token.LBRACE || tok == token.COMMA {
		if len(p.scope) == 0 || p.scope[len(p.scope)-1] == token.LBRACE {
			p.WriteByte('\n')
		}
	}
}
