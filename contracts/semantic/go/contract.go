// Package semantic owns captured inputs and bounded semantic publication.
// Language rules and operation-specific guarantees belong to providers.
package semantic

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
	ErrInvalidOperation = errors.New("invalid semantic operation")
	ErrInvalidRequest  = errors.New("invalid semantic request")
	ErrInvalidResult   = errors.New("invalid semantic result")
	ErrStaleResult     = errors.New("semantic result belongs to another invocation")
	ErrLimit           = errors.New("semantic resource limit")
)

// Obligation names an assertion or enumeration and its provider-defined domain. The domain is
// interpreted relative to the captured entry and inputs, not fetched by the core.
type Obligation struct {
	Key    semantics.Key
	Domain semantics.Key
}

// Requirement requires Complete and this exact verdict before a model may be
// published. Undetermined is useful for complete enumerations without a verdict.
type Requirement struct {
	Obligation semantics.Key
	Verdict    Verdict
}

type Contract struct {
	Profile           string
	Obligations       []Obligation
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
	Status     Completion
	Verdict    Verdict
}

// Diagnostic may name a current input node. Subject and Code are provider-owned
// vocabulary, allowing input-context failures to differ from source rejection.
type Diagnostic struct {
	Code    semantics.Key
	Subject semantics.Key
	Message string
	At      *ast.NodeRef
}

// Setting is uninterpreted configuration. Its owner defines encoding and meaning.
type Setting struct {
	Key   semantics.Key
	Value []byte
}

// Limits bound capture, invocation and publication. Zero Input/Options limits
// select defaults; partly specified or otherwise invalid limits are rejected.
// MaxRequested counts requested keys plus all dependency report requirements.
// MaxConfigurationEntries counts all own and dependency settings together.
// MaxConfigurationBytes bounds all requested keys, settings, dependency operation
// keys/profiles and dependency report keys together, before copying metadata.
type Limits struct {
	MaxInputs               int
	MaxInputNodes           int
	MaxInputEntries         int
	MaxInputBytes           int
	MaxRequested            int
	MaxConfigurationEntries int
	MaxConfigurationBytes   int
	MaxDiagnostics          int
	MaxDiagnosticBytes      int
	Model                   semantics.Limits
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

// Input is caller construction storage. Capture admits its immutable envelopes
// before calling an operation. Every tree must match the syntax descriptor.
// Callers must not mutate construction storage concurrently with capture.
type Input struct {
	Entry         ast.NodeRef
	Trees         []representation.AST
	Limits        Limits
}

// Options are captured separately for each invocation on an immutable Snapshot.
// Empty Requested selects all advertised obligations; unknown keys receive
// Unsupported reports. Callers must not mutate options during invocation.
type Options struct {
	Requested     []semantics.Key
	Configuration []Setting
	Limits        Limits
}

// LanguageID selects the envelope's registered language. Capture checks every
// input's ID, profile and actual Go marker before admitting the snapshot.
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
	Facts       *semantics.Data
	// ReuseModelFrom must be an admitted dependency of this invocation. It is
	// mutually exclusive with Facts; the model retains its original identity.
	ReuseModelFrom *Result
	Reports     []Report
	Diagnostics []Diagnostic
}

func (c Contract) Clone() Contract {
	c.Obligations = slices.Clone(c.Obligations)
	c.ModelRequirements = slices.Clone(c.ModelRequirements)
	return c
}

func (c Contract) Validate() error {
	bad := func(reason string) error { return fmt.Errorf("%w: %s", ErrInvalidOperation, reason) }
	if strings.TrimSpace(c.Profile) == "" || !utf8.ValidString(c.Profile) || len(c.Obligations) == 0 || len(c.Obligations) > 256 || len(c.ModelRequirements) > 256 {
		return bad("invalid semantic contract")
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
