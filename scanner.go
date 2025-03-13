package flatjson

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type token int

const (
	tokenIllegal token = iota
	tokenEOF
	tokenComment
	tokenWhitespace
	tokenIdentifier
	tokenNumber
	tokenString
	tokenFalse
	tokenNull
	tokenTrue

	eof = rune(0)
)

var tokenStrings = map[token]string{
	tokenIllegal:    "illegal",
	tokenEOF:        "eof",
	tokenComment:    "comment",
	tokenWhitespace: "whitespace",
	tokenIdentifier: "identifier",
	tokenNumber:     "number",
	tokenString:     "string",
	tokenFalse:      "false",
	tokenTrue:       "true",
	tokenNull:       "null",
}

func (t token) String() string {
	if s, ok := tokenStrings[t]; ok {
		return s
	}
	return string(rune(t))
}

func (t token) formatLiteral(lit string) string {
	switch t {
	case tokenIllegal, tokenWhitespace, tokenIdentifier, tokenNumber:
		return t.String() + "(" + strconv.QuoteToASCII(lit) + ")"
	case tokenEOF, tokenComment, tokenFalse, tokenTrue, tokenNull:
		return t.String()
	default:
		return "'" + strconv.QuoteRuneToASCII(rune(t)) + "'"
	}
}

// A scanner implements a scanner for a subset of JavaScript assignments. Its
// implementation is based on
// https://blog.gopheracademy.com/advent-2014/parsers-lexers/.
type scanner struct {
	r *bufio.Reader
}

func newScanner(r io.Reader) *scanner {
	return &scanner{
		r: bufio.NewReader(r),
	}
}

func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	switch {
	case errors.Is(err, io.EOF):
		return eof
	case err == nil:
		return ch
	default:
		return rune(-1)
	}
}

func (s *scanner) scan() (token, string) {
	ch := s.read()
	switch {
	case ch == eof:
		return tokenEOF, ""
	case isDigit(ch) || ch == '-':
		s.unread()
		return s.scanNumber()
	case isLetter(ch) || ch == '_':
		s.unread()
		return s.scanIdentifier()
	case isSpace(ch):
		s.unread()
		return s.scanWhitespace()
	case ch == '"':
		s.unread()
		return s.scanString()
	case ch == '.' || ch == ';' || ch == '=' || ch == '[' || ch == ']' || ch == '{' || ch == '}':
		return token(ch), string(ch)
	case ch == '/':
		s.unread()
		return s.scanComment()
	default:
		return tokenIllegal, string(ch)
	}
}

func (s *scanner) scanComment() (token, string) {
	var sb strings.Builder
	ch := s.read()
	sb.WriteRune(ch)
	ch = s.read()
	switch {
	case ch == '/':
		sb.WriteRune(ch)
		for {
			ch := s.read()
			switch {
			case ch == eof:
				return tokenComment, sb.String()
			case ch == '\n':
				sb.WriteRune(ch)
				return tokenComment, sb.String()
			default:
				sb.WriteRune(ch)
			}
		}
	case ch == '*':
		sb.WriteRune(ch)
		for {
			ch := s.read()
			switch {
			case ch == eof:
				return tokenIllegal, sb.String()
			case ch == '*':
				sb.WriteRune(ch)
				ch = s.read()
				switch {
				case ch == '/':
					sb.WriteRune(ch)
					return tokenComment, sb.String()
				default:
					s.unread()
				}
			default:
				sb.WriteRune(ch)
			}
		}
	default:
		return tokenIllegal, sb.String()
	}
}

func (s *scanner) scanIdentifier() (token, string) {
	var sb strings.Builder
	sb.WriteRune(s.read())
FOR:
	for {
		ch := s.read()
		switch {
		case ch == eof:
			break FOR
		case !isDigit(ch) && !isLetter(ch) && ch != '_':
			s.unread()
			break FOR
		default:
			sb.WriteRune(ch)
		}
	}
	switch s := sb.String(); s {
	case "false":
		return tokenFalse, s
	case "null":
		return tokenNull, s
	case "true":
		return tokenTrue, s
	default:
		return tokenIdentifier, s
	}
}

func (s *scanner) scanNumber() (token, string) {
	var sb strings.Builder
	ch := s.read()
	switch {
	case ch == '-':
		sb.WriteRune(ch)
	default:
		s.unread()
	}
	ch = s.read()
	switch {
	case ch == '0':
		sb.WriteRune(ch)
	case '1' <= ch && ch <= '9':
		sb.WriteRune(ch)
		for ch = s.read(); isDigit(ch); ch = s.read() {
			sb.WriteRune(ch)
		}
		s.unread()
	}
	ch = s.read()
	switch {
	case ch == '.':
		sb.WriteRune(ch)
		for ch = s.read(); isDigit(ch); ch = s.read() {
			sb.WriteRune(ch)
		}
		s.unread()
	default:
		s.unread()
	}
	ch = s.read()
	switch {
	case ch == 'e' || ch == 'E':
		sb.WriteRune(ch)
		ch = s.read()
		switch {
		case ch == '+' || ch == '-':
			sb.WriteRune(ch)
		default:
			s.unread()
		}
		ch := s.read()
		switch {
		case isDigit(ch):
			sb.WriteRune(ch)
		default:
			return tokenIllegal, ""
		}
		for ch = s.read(); isDigit(ch); ch = s.read() {
			sb.WriteRune(ch)
		}
		s.unread()
	default:
		s.unread()
	}
	return tokenNumber, sb.String()
}

func (s *scanner) scanString() (token, string) {
	var sb strings.Builder
	if ch := s.read(); ch != '"' {
		return tokenIllegal, string(ch)
	}
FOR:
	for {
		ch := s.read()
		switch {
		case ch == '"':
			break FOR
		case ch == '\\':
			switch ch := s.read(); ch {
			case '"', '\\', '/':
				sb.WriteRune(ch)
			case 'b':
				sb.WriteRune('\b')
			case 'f':
				sb.WriteRune('\f')
			case 'n':
				sb.WriteRune('\n')
			case 'r':
				sb.WriteRune('\r')
			case 't':
				sb.WriteRune('\t')
			case 'u':
				var r rune
				for range 4 {
					ch := s.read()
					switch {
					case isDigit(ch):
						r = (r << 4) | (ch - '0')
					case 'A' <= ch && ch <= 'F':
						r = (r << 4) | (ch - 'A' + 0xa)
					case 'a' <= ch && ch <= 'f':
						r = (r << 4) | (ch - 'f' + 0xa)
					default:
						return tokenIllegal, string(ch)
					}
				}
				sb.WriteRune(r)
			default:
				return tokenIllegal, string(ch)
			}
		default:
			sb.WriteRune(ch)
		}
	}
	return tokenString, sb.String()
}

func (s *scanner) scanWhitespace() (token, string) {
	var buf strings.Builder
	buf.WriteRune(s.read())
FOR:
	for {
		ch := s.read()
		switch {
		case ch == eof:
			break FOR
		case !isSpace(ch):
			s.unread()
			break FOR
		default:
			buf.WriteRune(ch)
		}
	}
	return tokenWhitespace, buf.String()
}

func (s *scanner) unread() {
	_ = s.r.UnreadRune()
}

func isDigit(r rune) bool  { return '0' <= r && r <= '9' }
func isLetter(r rune) bool { return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z') }
func isSpace(r rune) bool  { return r == '\t' || r == '\n' || r == '\r' || r == ' ' }
