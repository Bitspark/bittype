// Package capability defines tooling-only structural extraction. Format owners
// interpret source; the report host owns selection, association and presentation.
package capability

import (
	"fmt"
	"strings"
)

// Symbol identifies an original Go package declaration, or a method on a named
// receiver. Local names are not package declarations. An import alias is never
// part of this identity.
type Symbol struct {
	Package  string
	Receiver string
	Name     string
}

func (s Symbol) String() string {
	if s.Receiver != "" {
		return s.Package + "#" + s.Receiver + "." + s.Name
	}
	return s.Package + "#" + s.Name
}

// ParseSymbol accepts package#name and package#Receiver.method annotations.
func ParseSymbol(value string) (Symbol, error) {
	pkg, name, ok := strings.Cut(value, "#")
	if !ok || pkg == "" || name == "" || strings.Contains(name, "#") || strings.ContainsAny(value, " \t\r\n") {
		return Symbol{}, fmt.Errorf("invalid declaration identity %q", value)
	}
	s := Symbol{Package: pkg, Name: name}
	if receiver, method, found := strings.Cut(name, "."); found {
		if receiver == "" || method == "" || strings.Contains(method, ".") {
			return Symbol{}, fmt.Errorf("invalid method identity %q", value)
		}
		s.Receiver, s.Name = receiver, method
	}
	return s, nil
}

// FileRef is issued by one source view. A spelling from another view does not
// authorize a query; each implementation must check membership in its capture.
type FileRef string

// Span uses physical UTF-8 byte offsets, with an exclusive End. Source line
// directives do not alter it. A zero span is permitted only for optional fields.
type Span struct {
	File       FileRef
	Start, End int
}

type File struct {
	Ref  FileRef
	Path string
	Text string
}

// Declaration carries the full declaration and optional initializer/type spans.
// Value is absent for functions/types; Type is absent when not explicitly given.
type Declaration struct {
	Symbol Symbol
	Span   Span
	Value  Span
	Type   Span
}

// TypeFact follows true aliases to their named origin. A newly defined type has
// its own origin, even if its underlying structure resembles another type.
// Text is diagnostic spelling, never a selection key.
type TypeFact struct {
	Text      string
	Origin    *Symbol
	Arguments []TypeFact
}

// Constant carries Go compile-time values. String values are unquoted contents;
// other values use go/constant exact text. Kind names are bool, string, int,
// float and complex. No evaluated runtime value enters this contract.
type Constant struct {
	Kind  string
	Value string
}

type Facts struct {
	Symbol   *Symbol
	Callee   *Symbol
	Type     TypeFact
	Constant *Constant
}

// SourceView exposes immutable captured strings and copied facts. Implementations
// reject foreign files/spans and resolve declarations against this capture.
// Facts accepts exact expression spans only. A valid expression with no relevant
// fact yields a zero Facts value. No mutable compiler objects are shared.
type SourceView interface {
	File(FileRef) (File, error)
	Declaration(Symbol) (Declaration, error)
	Facts(Span) (Facts, error)
}

// ValidateSpan checks physical membership and range through the issuing view.
func ValidateSpan(source SourceView, span Span) error {
	if source == nil || span.File == "" || span.Start < 0 || span.End <= span.Start {
		return fmt.Errorf("invalid physical source span")
	}
	file, err := source.File(span.File)
	if err != nil {
		return err
	}
	if file.Ref != span.File || span.End > len(file.Text) {
		return fmt.Errorf("source span is outside its captured file")
	}
	return nil
}
