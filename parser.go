package flatjson

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type errUnexpected struct {
	tok      token
	lit      string
	expected []token
}

func newErrUnexpected(tok token, lit string, expected ...token) error {
	return &errUnexpected{
		tok:      tok,
		lit:      lit,
		expected: expected,
	}
}

func (e errUnexpected) Error() string {
	return fmt.Sprintf("expected %v, found %s", e.expected, e.tok.formatLiteral(e.lit))
}

type parser struct {
	s   *scanner
	buf struct {
		tok token
		lit string
		n   int
	}
}

type assignment struct {
	identifier string
	properties []interface{}
	value      interface{}
}

func newParser(r io.Reader) *parser {
	return &parser{
		s: newScanner(r),
	}
}

func (p *parser) parseAssignment() (*assignment, error) {
	tok, lit := p.scanIgnoreWhitespaceAndComments()
	if tok != tokenIdentifier {
		p.unscan()
		return nil, newErrUnexpected(tok, lit, tokenIdentifier)
	}
	identifier := lit
	var properties []interface{}
FOR:
	for {
		tok, lit := p.scanIgnoreWhitespaceAndComments()
		switch {
		case tok == token('='):
			p.unscan()
			break FOR
		case tok == token('.') || tok == token('['):
			p.unscan()
			property, err := p.parsePropertyAccess()
			if err != nil {
				return nil, err
			}
			properties = append(properties, property)
		default:
			return nil, newErrUnexpected(tok, lit, token('='), token('.'), token('['))
		}
	}
	tok, lit = p.scanIgnoreWhitespaceAndComments()
	if tok != token('=') {
		return nil, newErrUnexpected(tok, lit, token('='))
	}
	value, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	tok, lit = p.scanIgnoreWhitespaceAndComments()
	if tok != token(';') {
		return nil, newErrUnexpected(tok, lit, token(';'))
	}
	return &assignment{
		identifier: identifier,
		properties: properties,
		value:      value,
	}, nil
}

func (p *parser) parseAssignments() ([]*assignment, error) {
	var assignments []*assignment
FOR:
	for {
		tok, lit := p.scanIgnoreWhitespaceAndComments()
		switch tok {
		case tokenEOF:
			break FOR
		case tokenIdentifier:
			p.unscan()
			assignment, err := p.parseAssignment()
			if err != nil {
				return nil, err
			}
			assignments = append(assignments, assignment)
		default:
			return nil, newErrUnexpected(tok, lit, tokenEOF, tokenIdentifier)
		}
	}
	return assignments, nil
}

func (p *parser) parsePropertyAccess() (interface{}, error) {
	tok, lit := p.scanIgnoreWhitespaceAndComments()
	switch tok {
	case token('.'):
		switch tok, lit := p.scanIgnoreWhitespaceAndComments(); tok {
		case tokenIdentifier:
			return lit, nil
		default:
			return nil, newErrUnexpected(tok, lit, tokenIdentifier)
		}
	case token('['):
		var property interface{}
		switch tok, lit := p.scanIgnoreWhitespaceAndComments(); tok {
		case tokenNumber:
			property64, _ := strconv.ParseUint(lit, 10, 64)
			property = int(property64)
		case tokenString:
			property = lit
		default:
			return nil, newErrUnexpected(tok, lit, tokenNumber, tokenString)
		}
		if tok, lit := p.scanIgnoreWhitespaceAndComments(); tok != token(']') {
			return nil, newErrUnexpected(tok, lit, token(']'))
		}
		return property, nil
	default:
		return nil, newErrUnexpected(tok, lit, token('.'), token('['))
	}
}

func (p *parser) parseValue() (interface{}, error) {
	tok, lit := p.scanIgnoreWhitespaceAndComments()
	switch tok {
	case tokenNumber:
		return json.Number(lit), nil
	case tokenString:
		return lit, nil
	case tokenFalse:
		return false, nil
	case tokenTrue:
		return true, nil
	case tokenNull:
		return nil, nil
	case '[':
		tok, lit := p.scanIgnoreWhitespaceAndComments()
		if tok != token(']') {
			return nil, newErrUnexpected(tok, lit, token(']'))
		}
		return []interface{}{}, nil
	case '{':
		tok, lit := p.scanIgnoreWhitespaceAndComments()
		if tok != token('}') {
			return nil, newErrUnexpected(tok, lit, token('}'))
		}
		return make(map[string]interface{}), nil
	default:
		return nil, newErrUnexpected(tok, lit, tokenNumber, tokenString, tokenFalse, tokenTrue, tokenNull, token('['), token('{'))
	}
}

func (p *parser) scan() (token, string) {
	if p.buf.n > 0 {
		p.buf.n = 0
	} else {
		p.buf.tok, p.buf.lit = p.s.scan()
	}
	return p.buf.tok, p.buf.lit
}

func (p *parser) scanIgnoreWhitespaceAndComments() (token, string) {
	for {
		if tok, lit := p.scan(); tok != tokenWhitespace && tok != tokenComment {
			return tok, lit
		}
	}
}

func (p *parser) unscan() {
	p.buf.n = 1
}
