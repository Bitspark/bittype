// Package lex supplies configurable lexical mechanics. Languages supply their
// punctuation and comment/literal conventions; it has no language registry.
package lex

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/text/go"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token struct {
	Class, Text string
	Start, End  int
}
type Options struct {
	Punctuation     []string
	NestedComments  bool
	Lifetimes       bool
	RawIdentifiers  bool
	BacktickStrings bool
	Directives      []string
}

const MaxBytes = 4 << 20
const MaxTokens = 200000

func Scan(input text.Input, options Options) ([]Token, error) {
	bad := func(code text.Code, at, end int, message string) ([]Token, error) {
		return nil, &text.Diagnostic{Code: code, Message: message, Span: &ast.Span{Source: input.SourceID, Start: uint64(at), End: uint64(end)}}
	}
	if input.SourceID == "" || input.TreeID == "" {
		return bad(text.SyntaxError, 0, 0, "source and tree IDs are required")
	}
	if len(input.Text) > MaxBytes {
		return bad(text.LimitExceeded, 0, 0, "input exceeds 4 MiB")
	}
	if !utf8.ValidString(input.Text) {
		return bad(text.SyntaxError, 0, len(input.Text), "source must be UTF-8")
	}
	s := input.Text
	var out []Token
	for i := 0; i < len(s); {
		start := i
		class := "punctuation"
		r, size := utf8.DecodeRuneInString(s[i:])
		switch {
		case unicode.IsSpace(r) || (i == 0 && r == '\ufeff'):
			class = "whitespace"
			i += size
			for i < len(s) {
				r, size = utf8.DecodeRuneInString(s[i:])
				if !unicode.IsSpace(r) {
					break
				}
				i += size
			}
		case strings.HasPrefix(s[i:], "//"):
			class = "comment"
			for i < len(s) && s[i] != '\n' && s[i] != '\r' {
				i++
			}
			for _, directive := range options.Directives {
				if strings.HasPrefix(s[start:i], directive) {
					return bad(text.Unsupported, start, i, "directive comment is outside the supported syntax")
				}
			}
		case strings.HasPrefix(s[i:], "/*"):
			class = "comment"
			i += 2
			depth := 1
			for i < len(s) && depth > 0 {
				if strings.HasPrefix(s[i:], "*/") {
					depth--
					i += 2
				} else if options.NestedComments && strings.HasPrefix(s[i:], "/*") {
					depth++
					i += 2
				} else {
					_, n := utf8.DecodeRuneInString(s[i:])
					i += n
				}
			}
			if depth != 0 {
				return bad(text.SyntaxError, start, i, "unterminated comment")
			}
		case options.Lifetimes && r == '\'' && i+1 < len(s):
			j := i + 1
			for j < len(s) {
				rr, n := utf8.DecodeRuneInString(s[j:])
				if rr != '_' && !unicode.IsLetter(rr) && !unicode.IsDigit(rr) {
					break
				}
				j += n
			}
			if j == i+1 || (j < len(s) && s[j] == '\'') {
				return bad(text.Unsupported, start, j, "character values are outside this type syntax")
			}
			class = "lifetime"
			i = j
		case r == '"' || r == '\'' || (r == 96 && options.BacktickStrings):
			class = "string"
			quote := byte(r)
			i += size
			closed := false
			for i < len(s) {
				if s[i] == quote {
					i++
					closed = true
					break
				}
				if s[i] == '\\' && quote != 96 {
					i++
					if i == len(s) {
						break
					}
					_, n := utf8.DecodeRuneInString(s[i:])
					i += n
					continue
				}
				if (s[i] == '\n' || s[i] == '\r') && quote != 96 {
					return bad(text.SyntaxError, start, i, "newline in quoted token")
				}
				_, n := utf8.DecodeRuneInString(s[i:])
				i += n
			}
			if !closed {
				return bad(text.SyntaxError, start, i, "unterminated quoted token")
			}
		case unicode.IsLetter(r) || r == '_' || r == '$':
			class = "identifier"
			i += size
			if options.RawIdentifiers && r == 'r' && i < len(s) && s[i] == '#' {
				i++
			}
			for i < len(s) {
				rr, n := utf8.DecodeRuneInString(s[i:])
				if !unicode.IsLetter(rr) && !unicode.IsDigit(rr) && rr != '_' && rr != '$' {
					break
				}
				i += n
			}
		case r >= '0' && r <= '9':
			class = "number"
			i += size
			for i < len(s) {
				c := s[i]
				if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c == '.' {
					i++
					continue
				}
				if (c == '+' || c == '-') && i > start && (s[i-1] == 'e' || s[i-1] == 'E') {
					i++
					continue
				}
				break
			}
		default:
			i += size
			for _, punct := range options.Punctuation {
				if strings.HasPrefix(s[start:], punct) {
					i = start + len(punct)
					break
				}
			}
		}
		out = append(out, Token{class, s[start:i], start, i})
		if len(out) > MaxTokens {
			return bad(text.LimitExceeded, start, i, "too many tokens")
		}
	}
	return out, nil
}
