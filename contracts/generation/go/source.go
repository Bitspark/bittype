package generation

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Captured is the evidence context of one generation request: the admitted
// request (snapshot, entry tree and the exact ordered evidence) and the typed
// semantic graph that the designated evidence published. Only Plan creates it,
// by admitting its input against the source's own contract, and Plan opens that
// source over exactly this value. The kernel walks this graph and no other.
// Two captures are different evidence even when they share one graph.
type Captured[F any] struct{ state *captured[F] }

type captured[F any] struct {
	request transform.Request[F]
	graph   semantics.Semantics[F]
}

// admitInput admits input against a source contract. The contract must designate
// exactly one evidence requirement that publishes a semantic graph; the kernel
// unpacks that result's graph itself instead of accepting one from a source.
func admitInput[F any](language representation.Language, contract transform.Contract, input transform.Input) (Captured[F], error) {
	index := -1
	for i, requirement := range contract.Evidence {
		if requirement.Semantics {
			if index >= 0 {
				return Captured[F]{}, fmt.Errorf("%w: several evidence requirements publish a graph", ErrInvalid)
			}
			index = i
		}
	}
	if index < 0 {
		return Captured[F]{}, fmt.Errorf("%w: no evidence requirement publishes a graph", ErrInvalid)
	}
	request, err := transform.AdmitRequest[F](language, contract, input)
	if err != nil {
		return Captured[F]{}, err
	}
	result, ok := request.Evidence(index)
	if !ok {
		return Captured[F]{}, semantic.ErrEvidenceUnavailable
	}
	envelope, ok := result.Semantics()
	if !ok {
		return Captured[F]{}, semantic.ErrEvidenceUnavailable
	}
	graph, err := representation.UnpackSemantics[F](request.Snapshot().Language(), envelope)
	if err != nil {
		return Captured[F]{}, err
	}
	return Captured[F]{&captured[F]{request: request, graph: graph}}, nil
}

func (c Captured[F]) Valid() bool { return c.state != nil && c.state.graph.Valid() }

// Request is the admitted transform request; its evidence is the exact admitted set.
func (c Captured[F]) Request() transform.Request[F] {
	if c.state == nil {
		return transform.Request[F]{}
	}
	return c.state.request
}
func (c Captured[F]) Graph() semantics.Semantics[F] {
	if c.state == nil {
		return semantics.Semantics[F]{}
	}
	return c.state.graph
}

// Verdict is the outcome of a source judgment over the captured evidence.
type Verdict uint8

const (
	// Admitted: the evidence establishes that the subject may be represented.
	Admitted Verdict = iota + 1
	// Unsupported: valid source that this generation integration cannot represent.
	Unsupported
	// Rejected: the evidence refutes the subject's admissibility.
	Rejected
	// Undetermined: the evidence does not establish the judgment either way.
	Undetermined
)

// Judgment is a source-owned conclusion. Diagnostic optionally explains a
// non-admitted verdict at a source site. The kernel never interprets it, and a
// diagnostic cannot add an error class: only the verdict classifies.
type Judgment struct {
	Verdict    Verdict
	Diagnostic error
}

// Interpretation is a source integration's reading of one captured graph. The
// kernel calls it with handles it issued for that graph and performs every
// structural step itself. An interpretation uses only the captured evidence:
// it never resolves, analyzes or traverses on the kernel's behalf.
type Interpretation[F any] interface {
	// Declaration classifies a declaration entity the walk reached. ok is false
	// for an entity this source does not generate (for example a kind or module).
	// Every selector the interpretation returns is qualified by the capture's
	// source language. A selector in another language or profile is outside the
	// run: planning fails with ErrScope before any exact or default rule is
	// chosen.
	Declaration(context.Context, Declaration[F]) (Selector, bool, error)
	// Expression returns the selector of a provider-expression occurrence.
	Expression(context.Context, Use[F]) (Selector, error)
	// JudgeDeclaration judges a declaration before its body is planned.
	JudgeDeclaration(context.Context, Declaration[F]) (Judgment, error)
	// JudgeUse judges one occurrence under the assumptions of the declaration
	// it occurs in. Before representing an occurrence the kernel judges it and
	// every occurrence in its subtree, including omitted actuals and uncaptured
	// operands. The judgment of a use covers that use: a callee declaration is
	// classified only when a use of it is represented, and judged only when it
	// is generated.
	JudgeUse(context.Context, Use[F]) (Judgment, error)
}

// Source is what a provider opens over one capture: the kernel-facing
// interpretation and the reader that mapping libraries use to interpret source
// payloads. R reads kernel handles; it supplies no navigation. Only Plan opens
// a source, over the capture it walks. Interpreting that capture is the
// provider's responsibility: the kernel cannot check where a provider's
// answers come from.
type Source[F, R any] struct {
	Interpretation Interpretation[F]
	Reader         R
	// Roles is the issuer of the role selectors that the reader hands to
	// rules (prototype, bittype#196, row 12). A source without field roles
	// leaves it zero.
	Roles RoleIssuer[F]
	// Semantics names the source interpretation that the checked references
	// use: the fifth version axis (row 16). The plan records it.
	Semantics string
}

// bittype:dimension Syntax
type SourceProvider[F, R any] interface {
	// Contract declares the required evidence; exactly one requirement
	// publishes the semantic graph.
	Contract() transform.Contract
	// Open interprets the capture Plan admitted against Contract, and only that
	// capture. A source resolves kernel handles against it (see
	// Declaration.In), so it can refuse handles from a run over other evidence.
	Open(context.Context, Captured[F]) (Source[F, R], error)
}

// SiteError is a kernel failure at a source syntax site, for diagnostics.
type SiteError struct {
	Err     error
	Site    ast.NodeRef
	Message string
}

func (e *SiteError) Error() string {
	if e.Site.Tree == "" {
		return e.Message + ": " + e.Err.Error()
	}
	return fmt.Sprintf("%s#%d: %s: %v", e.Site.Tree, e.Site.Node, e.Message, e.Err)
}
func (e *SiteError) Unwrap() error { return e.Err }

// explanation carries a source diagnostic. It does not unwrap, so a diagnostic
// cannot add a failure class to the verdict; As still reaches the diagnostic
// and any site it carries.
type explanation struct{ err error }

func (e explanation) Error() string      { return e.err.Error() }
func (e explanation) As(target any) bool { return errors.As(e.err, target) }

// verdict converts a judgment into the kernel's failure meaning.
func verdict(j Judgment) error {
	wrap := func(kind error) error {
		if j.Diagnostic == nil {
			return kind
		}
		return fmt.Errorf("%w: %w", kind, explanation{j.Diagnostic})
	}
	switch j.Verdict {
	case Admitted:
		return nil
	case Unsupported:
		return wrap(ErrUnsupported)
	case Rejected:
		return wrap(ErrRejected)
	case Undetermined:
		return wrap(semantic.ErrEvidenceUnavailable)
	default:
		return fmt.Errorf("%w: invalid judgment", ErrInvalid)
	}
}

func cloneOperands[F any](o []Operand[F]) []Operand[F] {
	out := slices.Clone(o)
	for i := range out {
		out[i].Locator = slices.Clone(out[i].Locator)
	}
	return out
}
