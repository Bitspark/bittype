package semantic

import (
	"errors"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"slices"
)

var ErrEvidenceUnavailable = errors.New("required semantic evidence is not established")

// EvidenceRequirement is a consumer-owned compatibility requirement, not a
// request to run a provider. Semantics additionally requires a published graph.
type EvidenceRequirement struct {
	Operation     semantics.Key
	Profile       string
	Configuration []Setting
	Reports       []Requirement
	Semantics     bool
}

func (r EvidenceRequirement) Clone() EvidenceRequirement {
	r.Configuration = cloneSettings(r.Configuration)
	r.Reports = slices.Clone(r.Reports)
	return r
}

// ValidateEvidenceRequirements applies one shared metadata budget to the set.
func ValidateEvidenceRequirements(requirements []EvidenceRequirement) error {
	limits := DefaultLimits()
	if len(requirements) > limits.MaxInputs {
		return ErrLimit
	}
	budget := newMetadataBudget(limits)
	for _, r := range requirements {
		if err := validateDependencySpecification(r.Operation, r.Profile, r.Configuration, r.Reports, &budget); err != nil {
			return err
		}
	}
	return nil
}

// AdmitEvidence shares semantic dependency admission, then requires complete
// reports and any requested model. It never invokes a resolver or analyzer.
func AdmitEvidence(source Snapshot, requirements []EvidenceRequirement, results []Result) ([]Result, error) {
	if !source.Valid() || len(requirements) != len(results) {
		return nil, ErrInvalidRequest
	}
	if err := ValidateEvidenceRequirements(requirements); err != nil {
		return nil, err
	}
	for i, r := range requirements {
		dependency, err := admitDependency(source.state, Dependency{Result: results[i], Operation: r.Operation, Profile: r.Profile, Configuration: r.Configuration, Reports: r.Reports})
		if err != nil {
			return nil, err
		}
		if !dependency.reportsEstablished() || r.Semantics && !results[i].HasSemantics() {
			return nil, ErrEvidenceUnavailable
		}
	}
	return slices.Clone(results), nil
}

// Inputs returns immutable envelopes in a separately owned slice.
func (s Snapshot) Inputs() []representation.AST {
	if !s.Valid() {
		return nil
	}
	return slices.Clone(s.state.trees)
}

// EntryTree is the captured tree containing the explicit entry node.
func (s Snapshot) EntryTree() (representation.AST, error) {
	if !s.Valid() {
		return representation.AST{}, ErrInvalidRequest
	}
	for _, input := range s.state.trees {
		info, err := input.Inspect()
		if err != nil {
			return representation.AST{}, err
		}
		if info.ID == s.state.entry.Tree {
			return input, nil
		}
	}
	return representation.AST{}, ErrInvalidRequest
}

// Failure retains semantic reports/diagnostics when a higher-level operation
// cannot publish an admitted result. It never fabricates a partial Source.
type Failure struct{ Result Result }

func (e *Failure) Error() string {
	diagnostics := e.Result.Diagnostics()
	if len(diagnostics) > 0 {
		return "semantics unavailable: " + diagnostics[0].Message
	}
	return ErrEvidenceUnavailable.Error()
}
func (e *Failure) Unwrap() error { return ErrEvidenceUnavailable }
