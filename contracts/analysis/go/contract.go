// Package analysis defines explicit, typed semantic-analysis capabilities.
// It owns request/result publication, never a language's semantic rules.
package analysis

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var (
	ErrInvalidAnalyzer = errors.New("invalid semantic analyzer")
	ErrInvalidRequest = errors.New("invalid analysis request")
	ErrInvalidResult = errors.New("invalid analysis result")
	ErrStaleResult = errors.New("analysis result belongs to another request")
	ErrLimit = errors.New("analysis resource limit")
)

// Obligation names an analysis and its provider-defined domain. The domain is
// interpreted relative to the captured entry and inputs, not fetched by the core.
type Obligation struct {
	Key semantics.Key
	Domain semantics.Key
}

// Requirement requires Complete and this exact verdict before a model may be
// published. Undetermined is useful for complete enumerations without a verdict.
type Requirement struct {
	Obligation semantics.Key
	Verdict Verdict
}

type Contract struct {
	Profile string
	Obligations []Obligation
	ModelRequirements []Requirement
}

type Completion uint8

const (
	Unreported Completion = iota
	Complete
	Incomplete
	Unsupported
)

type Verdict uint8

const (
	Undetermined Verdict = iota
	Holds
	Rejected
)

// Report separates completion from truth. Unsupported and Incomplete require
// Undetermined; a completed validity decision may hold or reject its subject.
type Report struct {
	Obligation semantics.Key
	Status Completion
	Verdict Verdict
}

// Diagnostic may name a current input node. Subject and Code are provider-owned
// vocabulary, allowing input-context failures to differ from source rejection.
type Diagnostic struct {
	Code semantics.Key
	Subject semantics.Key
	Message string
	At *ast.NodeRef
}

// Setting is uninterpreted configuration. Its owner defines encoding and meaning.
type Setting struct {
	Key semantics.Key
	Value []byte
}

// Limits bound captures and publication. A wholly zero value in Input selects
// DefaultLimits; partly specified or otherwise invalid limits are rejected.
type Limits struct {
	MaxInputs int
	MaxInputNodes int
	MaxInputEntries int
	MaxInputBytes int
	MaxRequested int
	MaxConfigurationEntries int
	MaxConfigurationBytes int
	MaxDiagnostics int
	MaxDiagnosticBytes int
	Model semantics.Limits
}

func DefaultLimits() Limits {
	return Limits{
		MaxInputs: 256, MaxInputNodes: 200000, MaxInputEntries: 1600000, MaxInputBytes: 16 << 20,
		MaxRequested: 256, MaxConfigurationEntries: 256, MaxConfigurationBytes: 1 << 20,
		MaxDiagnostics: 256, MaxDiagnosticBytes: 1 << 20, Model: semantics.DefaultLimits(),
	}
}

func (l Limits) Valid() bool {
	return l.MaxInputs > 0 && l.MaxInputNodes > 0 && l.MaxInputEntries > 0 && l.MaxInputBytes > 0 &&
		l.MaxRequested > 0 && l.MaxConfigurationEntries > 0 && l.MaxConfigurationBytes > 0 &&
		l.MaxDiagnostics > 0 && l.MaxDiagnosticBytes > 0 && l.Model.Valid()
}

// Input is caller construction storage. Analyze captures its storage before
// calling a provider. Callers must not mutate it concurrently with capture.
// Every tree must match the analyzer's syntax descriptor. Empty Requested selects
// all advertised obligations; unknown requested keys receive Unsupported reports.
type Input struct {
	Entry ast.NodeRef
	Trees []representation.AST
	Requested []semantics.Key
	Configuration []Setting
	Limits Limits
}

// LanguageID selects the envelope's registered language. Analyze still checks
// every input's ID, profile and actual Go marker before invoking the provider.
func (i Input) LanguageID() language.ID {
	if len(i.Trees) == 0 {
		return ""
	}
	return i.Trees[0].LanguageID()
}

// Publication is provider construction storage. Request.Publish validates and
// copies it during the callback, before returning an immutable Result. Producers
// must not mutate these records concurrently with publication.
type Publication struct {
	Facts *semantics.Data
	Reports []Report
	Diagnostics []Diagnostic
}

func cloneContract(c Contract) Contract {
	c.Obligations = slices.Clone(c.Obligations)
	c.ModelRequirements = slices.Clone(c.ModelRequirements)
	return c
}

func validateContract(c Contract) error {
	bad := func(reason string) error { return fmt.Errorf("%w: %s", ErrInvalidAnalyzer, reason) }
	if strings.TrimSpace(c.Profile) == "" || !utf8.ValidString(c.Profile) || len(c.Obligations) == 0 || len(c.Obligations) > 256 || len(c.ModelRequirements) > 256 {
		return bad("invalid analysis contract")
	}
	bytes := 1 << 20
	if !take(&bytes, len(c.Profile)) {
		return bad("contract metadata exceeds limit")
	}
	known := make(map[semantics.Key]bool)
	for _, obligation := range c.Obligations {
		if !obligation.Key.Valid() || !obligation.Domain.Valid() || known[obligation.Key] || !keyBytes(&bytes, obligation.Key) || !keyBytes(&bytes, obligation.Domain) {
			return bad("invalid or duplicate obligation")
		}
		known[obligation.Key] = true
	}
	seen := make(map[semantics.Key]bool)
	for _, requirement := range c.ModelRequirements {
		if !known[requirement.Obligation] || seen[requirement.Obligation] || requirement.Verdict > Rejected {
			return bad("invalid model requirement")
		}
		seen[requirement.Obligation] = true
	}
	return nil
}

func take(left *int, n int) bool {
	if n < 0 || n > *left {
		return false
	}
	*left -= n
	return true
}

func keyBytes(left *int, key semantics.Key) bool {
	return take(left, len(key.Owner)) && take(left, len(key.Name))
}
