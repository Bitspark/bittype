package source

import (
	"context"
	"strings"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/horos/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func Contract() transform.Contract {
	return transform.Contract{Evidence: []semantic.EvidenceRequirement{{Operation: analysis.Kind(), Profile: n.AnalysisProfile, Reports: n.AnalysisContract().Semantic.SemanticsRequirements, Semantics: true}}}
}

// interpretation is the horos reading of one captured graph for the kernel.
// It classifies declarations, names descriptors and judges proper sorts from
// the analysis facts. It never walks the graph on the kernel's behalf.
type interpretation struct {
	captured   g.Captured[h.Language] // handles resolve against this capture only
	graph      semantics.Semantics[h.Language]
	language   representation.Language
	data       semantics.Data
	sorts      map[semantics.Subject]semantics.TermID
	results    map[semantics.EntityID]semantics.TermID
	identities map[semantics.EntityID]semantics.TermID
	sites      map[semantics.TermID]ast.NodeRef
	rootModule ast.NodeRef
}

func open(ctx context.Context, c g.Captured[h.Language]) (*interpretation, error) {
	if ctx == nil || !c.Valid() {
		return nil, g.ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	request := c.Request()
	input := request.Source()
	if request.Entry() != (ast.NodeRef{Tree: input.ID, Node: input.Root}) {
		return nil, semantic.ErrInvalidRequest
	}
	s := &interpretation{captured: c, graph: c.Graph(), language: request.Snapshot().Language(), data: c.Graph().Data(), sorts: map[semantics.Subject]semantics.TermID{}, results: map[semantics.EntityID]semantics.TermID{}, identities: map[semantics.EntityID]semantics.TermID{}, sites: map[semantics.TermID]ast.NodeRef{}}
	for _, ref := range input.Nodes[input.Root-1].References {
		if ref.Role == h.RoleRoot {
			s.rootModule = ref.Target
		}
	}
	if s.rootModule.Tree == "" {
		return nil, semantic.ErrInvalidRequest
	}
	for _, e := range s.data.Expressions {
		if e.At.Role == n.ExpressionRole() {
			s.sites[e.Term] = e.At.Site
		}
	}
	// horos publishes at most one value for each of these relations per subject.
	for _, f := range s.data.Facts {
		switch f.Relation {
		case n.SortRelation():
			if s.sorts[f.Subject] != 0 {
				return nil, semantic.ErrInvalidRequest
			}
			s.sorts[f.Subject] = f.Term
		case n.ResultSortRelation():
			if f.Subject.Entity == 0 || s.results[f.Subject.Entity] != 0 {
				return nil, semantic.ErrInvalidRequest
			}
			s.results[f.Subject.Entity] = f.Term
		case n.IdentityRelation():
			if f.Subject.Entity == 0 || s.identities[f.Subject.Entity] != 0 {
				return nil, semantic.ErrInvalidRequest
			}
			s.identities[f.Subject.Entity] = f.Term
		}
	}
	return s, nil
}

// definition is a horos type definition the kernel reached, or ok=false for any
// other entity (a kind definition, module or parameter).
func (s *interpretation) definition(d g.Declaration[h.Language]) (semantics.Entity, bool, error) {
	id, ok := d.In(s.captured)
	if !ok {
		return semantics.Entity{}, false, g.ErrScope
	}
	e := s.data.Entities[id-1]
	if e.Kind != n.TypeDefinitionKind() {
		return semantics.Entity{}, false, nil
	}
	if e.Owner == 0 || e.Body == 0 || len(e.Declarations) != 1 {
		return semantics.Entity{}, false, semantic.ErrInvalidRequest
	}
	return e, true, nil
}

func (s *interpretation) Declaration(ctx context.Context, d g.Declaration[h.Language]) (g.Selector, bool, error) {
	if err := ctx.Err(); err != nil {
		return g.Selector{}, false, err
	}
	e, ok, err := s.definition(d)
	if err != nil || !ok {
		return g.Selector{}, false, err
	}
	pin, err := s.modulePin(e.Owner)
	if err != nil {
		return g.Selector{}, false, err
	}
	key, err := Named(s.language, pin, e.Name)
	return key, true, err
}

func (s *interpretation) Expression(ctx context.Context, u g.Use[h.Language]) (g.Selector, error) {
	if err := ctx.Err(); err != nil {
		return g.Selector{}, err
	}
	d, ok, err := s.descriptor(u)
	if err != nil {
		return g.Selector{}, err
	}
	if !ok {
		return g.Selector{}, semantic.ErrInvalidRequest
	}
	return Builtin(s.language, d.Constructor, d.Discriminator)
}

// JudgeDeclaration requires a proper result sort and proper parameter sorts:
// this integration has no default representation for type constructors.
func (s *interpretation) JudgeDeclaration(ctx context.Context, d g.Declaration[h.Language]) (g.Judgment, error) {
	if err := ctx.Err(); err != nil {
		return g.Judgment{}, err
	}
	id, _ := d.In(s.captured)
	e, ok, err := s.definition(d)
	if err != nil {
		return g.Judgment{}, err
	}
	if !ok {
		return g.Judgment{Verdict: g.Unsupported}, nil
	}
	site := s.sites[e.Body]
	if j := s.proper(s.results[id], site); j.Verdict != g.Admitted {
		return j, nil
	}
	for _, p := range e.Parameters {
		if j := s.proper(s.sorts[semantics.Subject{Entity: p.Entity}], site); j.Verdict != g.Admitted {
			return j, nil
		}
	}
	return g.Judgment{Verdict: g.Admitted}, nil
}

// JudgeUse requires a proper sort for the occurrence's term. The successful
// horos analysis (program-validity Holds) already established the whole
// occurrence, including every actual and operand, under its declaration's
// original assumptions; this adds the generation restriction to proper types.
func (s *interpretation) JudgeUse(ctx context.Context, u g.Use[h.Language]) (g.Judgment, error) {
	if err := ctx.Err(); err != nil {
		return g.Judgment{}, err
	}
	term, ok := u.In(s.captured)
	if !ok {
		return g.Judgment{}, g.ErrScope
	}
	return s.proper(s.sorts[semantics.Subject{Term: term}], s.sites[term]), nil
}

// proper: a missing sort fact is uncovered, not malformed; an arrow sort is a
// known result that this integration cannot represent.
func (s *interpretation) proper(sort semantics.TermID, site ast.NodeRef) g.Judgment {
	if sort == 0 || int(sort) > len(s.data.Terms) {
		// Not an unsupported form: the evidence simply does not cover it.
		return g.Judgment{Verdict: g.Undetermined, Diagnostic: &g.SiteError{Err: semantic.ErrEvidenceUnavailable, Site: site, Message: "no established sort"}}
	}
	t := s.data.Terms[sort-1]
	if t.Kind == semantics.ProviderExpression && t.Operator == n.ProperSort() {
		return g.Judgment{Verdict: g.Admitted}
	}
	return g.Judgment{Verdict: g.Unsupported, Diagnostic: s.diagnostic(site, "constructor-valued types require an explicit representation")}
}

func (s *interpretation) diagnostic(site ast.NodeRef, message string) error {
	return &text.Diagnostic{Code: text.Unsupported, Message: message, Node: &site}
}

func (s *interpretation) descriptor(u g.Use[h.Language]) (Descriptor, bool, error) {
	id, ok := u.In(s.captured)
	if !ok {
		return Descriptor{}, false, g.ErrScope
	}
	t := s.data.Terms[id-1]
	d := Descriptor{Origin: s.sites[id]}
	if t.Kind != semantics.ProviderExpression {
		return d, false, nil
	}
	if t.Operator.Owner != "horos" || !strings.HasPrefix(t.Operator.Name, "descriptor/") {
		return d, false, semantic.ErrInvalidRequest
	}
	d.Constructor = strings.TrimPrefix(t.Operator.Name, "descriptor/")
	if d.Constructor == "scalar" || d.Constructor == "record" {
		p, err := s.payload(id)
		if err != nil {
			return d, false, err
		}
		if p.Len() < 2 {
			return d, false, semantic.ErrInvalidRequest
		}
		d.Discriminator = atom(p.At(1))
	}
	return d, true, nil
}
