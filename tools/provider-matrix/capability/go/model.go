package capability

import (
	"context"
	"fmt"
)

// Selector specifies exactly one original definition or named-type origin.
// Selection is host-owned; no match is absence and multiple matches are errors.
type Selector struct {
	Definition *Symbol
	NamedType  *Symbol
}

// Extractor interprets a selected declaration. It receives no provider/profile
// identities and cannot redirect its result to another operation. The host binds
// each extractor to a format identity in consuming configuration.
type Extractor interface {
	Selectors() []Selector
	Extract(context.Context, DefinitionRequest) (Inventory, error)
}

type DefinitionRequest struct {
	Root   Symbol
	Source SourceView
	Limits Limits
}

// Limits bound cooperative traversal and the accepted result. They do not
// sandbox an in-process extractor. Source loading and command deadlines belong
// to the host. Exhausting a limit is an error, not a partial successful report.
type Limits struct {
	MaxDepth           int
	MaxExpansions      int
	MaxRecords         int
	MaxDiagnosticBytes int
}

func DefaultLimits() Limits {
	return Limits{MaxDepth: 64, MaxExpansions: 100000, MaxRecords: 50000, MaxDiagnosticBytes: 4096}
}

// EffectiveLimits fills omitted limits and rejects negative values.
func EffectiveLimits(l Limits) (Limits, error) {
	d := DefaultLimits()
	fields := []*int{&l.MaxDepth, &l.MaxExpansions, &l.MaxRecords, &l.MaxDiagnosticBytes}
	defaults := []int{d.MaxDepth, d.MaxExpansions, d.MaxRecords, d.MaxDiagnosticBytes}
	for i, field := range fields {
		if *field < 0 {
			return Limits{}, fmt.Errorf("negative extraction limit")
		}
		if *field == 0 {
			*field = defaults[i]
		}
	}
	return l, nil
}

type Knowledge string

const (
	Known   Knowledge = "known"
	Partial Knowledge = "partial"
	Unknown Knowledge = "unknown"
	Opaque  Knowledge = "opaque"
)

// Field makes uncertainty local. Known includes known empty/false/zero. Partial
// carries a known subset of a collection; Unknown carries no value. Opaque is
// reserved for an attribute predicate whose behavior was not interpreted.
// Any state other than Known requires a reason and a physical source location.
type Field[T any] struct {
	State  Knowledge
	Value  T
	At     Span
	Reason string
}

type UpperBound struct {
	Unbounded bool
	Value     int
}

// ValueSet distinguishes an absent value predicate from a finite accepted set.
// Known empty Values with Unrestricted false means no accepted string values.
type ValueSet struct {
	Unrestricted bool
	Values       []string
}

type Role struct {
	Name  string
	At    Span
	Min   Field[int]
	Max   Field[UpperBound]
	Kinds Field[[]string]
}

type Attribute struct {
	Name     string
	At       Span
	Required Field[bool]
	Values   Field[ValueSet]
}

// Feature is an interpreted source entry, not a statement about all runtime
// values or reachable behavior. Its identity is (definition root, format, Key),
// with the first two components attached by the host. Roles preserve order.
// References lists allowed reference-role names, not resolution semantics.
type Feature struct {
	Key        string
	At         Span
	Roles      Field[[]Role]
	Attributes Field[[]Attribute]
	References Field[[]string]
}

type Limitation struct {
	At      Span
	Message string
}

type ExtractionExtent string

const (
	Decoded       ExtractionExtent = "decoded"
	PartlyDecoded ExtractionExtent = "partial"
)

// Extent concerns only the selected declaration and this structural schema.
// Decoded requires every emitted field to be known and no unresolved portion.
// Neither extent asserts runtime stability, semantic completeness or coverage.
type Inventory struct {
	Extent      ExtractionExtent
	Features    []Feature
	Limitations []Limitation
}
