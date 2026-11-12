// Package source retains concrete files together with their established semantics.
// Source has no filesystem operations and no knowledge of language rules.
package source

import (
	"context"
	"errors"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"reflect"
	"slices"
	"sync/atomic"
)

type Unit[L any] struct {
	File   files.ID
	Path   files.Path
	Source ast.SourceID
	Syntax cst.CST[L]
}
type GroupID string
type Group struct {
	ID      GroupID
	Kind    semantics.Key
	Members []files.ID
	Sites   []ast.NodeRef
}
type Relation struct {
	Kind     semantics.Key
	From, To files.ID
	At       *ast.NodeRef
}
type Publication[L any] struct {
	Snapshot  semantic.Snapshot
	Analysis  semantic.Result
	Units     []Unit[L]
	Groups    []Group
	Relations []Relation
}

// Source retains concrete units and their established semantics.
//
// bittype:dimension Files
type Source[L any] struct{ state *state[L] }
type state[L any] struct {
	owner       *publicationIdentity
	publication Publication[L]
	abstract    semantics.Semantics[L]
}

// A Files stage refuses an input in one of four categories, which a caller
// tests with errors.Is: invalid input is ErrInvalid, unsupported coverage is
// ErrUnsupported, missing evidence is semantic.ErrEvidenceUnavailable and an
// exceeded limit is ErrLimit. The category of a refusal is part of the Files
// contract. A refusal that carries another package's value joins it to its
// category, so errors.Is holds for both. docs/source.md lists the refusals
// that carry a category.
var ErrInvalid = errors.New("invalid Source")
var ErrStale = errors.New("Source belongs to another Files invocation")
var ErrLimit = errors.New("Source resource limit")

// ErrUnsupported refuses an input that is well formed and that this profile
// does not cover yet. A later change may turn an ErrUnsupported refusal into
// success, and it may neither admit an input refused with ErrInvalid nor move a
// refusal to another category.
var ErrUnsupported = errors.New("Files input outside supported coverage")

// Limits bound concrete storage and metadata in addition to semantic capture limits.
type Limits struct {
	Files                                           files.Limits
	MaxNodes, MaxLinks, MaxBytes, MaxDepth, MaxWork int
}

func DefaultLimits() Limits {
	return Limits{files.DefaultLimits(), 400000, 1600000, 16 << 20, 512, 4000000}
}
func (l Limits) Valid() bool {
	return l.Files.Valid() && l.MaxNodes > 0 && l.MaxLinks > 0 && l.MaxBytes > 0 && l.MaxDepth > 0 && l.MaxDepth <= 512 && l.MaxWork > 0
}

// Policy is an explicit language-owned admission contract. Check establishes
// CST/AST agreement, grouping and the operation's preservation obligations.
// Print supplies the configured Syntax printer; neither callback is discovered.
type Policy[L any] struct {
	Evidence       semantic.EvidenceRequirement
	SemanticLimits semantic.Limits
	Limits         Limits
	Print          func(context.Context, cst.CST[L]) (string, error)
	Check          func(context.Context, Publication[L]) error
}
type Publisher[L any] struct{ state *publisherState[L] }
type publicationIdentity struct{ marker byte }

type publisherState[L any] struct {
	identity   *publicationIdentity
	ctx        context.Context
	definition representation.Language
	policy     Policy[L]
	active     atomic.Bool
}

func NewPublisher[L any](ctx context.Context, definition representation.Language, policy Policy[L]) (Publisher[L], error) {
	if policy.SemanticLimits == (semantic.Limits{}) {
		policy.SemanticLimits = semantic.DefaultLimits()
	}
	if !policy.SemanticLimits.Valid() {
		return Publisher[L]{}, ErrInvalid
	}
	if ctx == nil || policy.Print == nil || policy.Check == nil || !policy.Limits.Valid() || !policy.Evidence.Semantics {
		return Publisher[L]{}, ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return Publisher[L]{}, err
	}
	if err := representation.CheckMarker[L](definition); err != nil {
		return Publisher[L]{}, err
	}
	if err := semantic.ValidateEvidenceRequirements([]semantic.EvidenceRequirement{policy.Evidence}); err != nil {
		return Publisher[L]{}, err
	}
	policy.Evidence = policy.Evidence.Clone()
	p := &publisherState[L]{ctx: ctx, definition: definition, policy: policy, identity: &publicationIdentity{}}
	p.active.Store(true)
	return Publisher[L]{p}, nil
}
func (p Publisher[L]) Close() {
	if p.state != nil {
		p.state.active.Store(false)
	}
}
func (p Publisher[L]) Accept(s Source[L]) error {
	if p.state == nil || !p.state.active.Load() {
		return ErrInvalid
	}
	if err := p.state.ctx.Err(); err != nil {
		return err
	}
	if !s.Valid() || s.state.owner != p.state.identity {
		return ErrStale
	}
	return nil
}
func (s Source[L]) Valid() bool { return s.state != nil }
func (s Source[L]) Entry() ast.NodeRef {
	if !s.Valid() {
		return ast.NodeRef{}
	}
	return s.state.publication.Snapshot.Entry()
}
func (s Source[L]) Snapshot() semantic.Snapshot {
	if !s.Valid() {
		return semantic.Snapshot{}
	}
	return s.state.publication.Snapshot
}
func (s Source[L]) Evidence() semantic.Result {
	if !s.Valid() {
		return semantic.Result{}
	}
	return s.state.publication.Analysis
}
func (s Source[L]) SemanticsView() semantics.Semantics[L] {
	if !s.Valid() {
		return semantics.Semantics[L]{}
	}
	return s.state.abstract
}
func (s Source[L]) Units() []Unit[L] {
	if !s.Valid() {
		return nil
	}
	units := slices.Clone(s.state.publication.Units)
	for i := range units {
		packed, _ := representation.PackCST(s.Snapshot().Language(), units[i].Syntax)
		units[i].Syntax, _ = representation.UnpackCST[L](s.Snapshot().Language(), packed)
	}
	return units
}
func cloneGroups(groups []Group) []Group {
	out := slices.Clone(groups)
	for i := range out {
		out[i].Members = slices.Clone(out[i].Members)
		out[i].Sites = slices.Clone(out[i].Sites)
	}
	return out
}
func cloneRelations(relations []Relation) []Relation {
	out := slices.Clone(relations)
	for i := range out {
		if out[i].At != nil {
			at := *out[i].At
			out[i].At = &at
		}
	}
	return out
}
func (s Source[L]) Groups() []Group {
	if !s.Valid() {
		return nil
	}
	return cloneGroups(s.state.publication.Groups)
}
func (s Source[L]) Relations() []Relation {
	if !s.Valid() {
		return nil
	}
	return cloneRelations(s.state.publication.Relations)
}

func (p Publisher[L]) Publish(input Publication[L]) (Source[L], error) {
	if p.state == nil || !p.state.active.Load() {
		return Source[L]{}, ErrInvalid
	}
	q := p.state
	if err := q.ctx.Err(); err != nil {
		return Source[L]{}, err
	}
	if !q.definition.Compatible(input.Snapshot.Language()) {
		return Source[L]{}, representation.ErrLanguageMismatch
	}
	if _, err := semantic.AdmitEvidence(input.Snapshot, []semantic.EvidenceRequirement{q.policy.Evidence}, []semantic.Result{input.Analysis}); err != nil {
		return Source[L]{}, err
	}
	if err := input.Snapshot.CheckLimits(q.policy.SemanticLimits); err != nil {
		return Source[L]{}, err
	}
	graph, ok := input.Analysis.Semantics()
	if !ok {
		return Source[L]{}, semantic.ErrEvidenceUnavailable
	}
	if err := graph.CheckLimits(q.policy.SemanticLimits.Semantics); err != nil {
		return Source[L]{}, err
	}
	abstract, err := representation.UnpackSemantics[L](q.definition, graph)
	if err != nil {
		return Source[L]{}, err
	}
	contextFiles := input.Snapshot.Files()
	if contextFiles == nil || len(input.Units) != len(contextFiles.Files.Entries) {
		return Source[L]{}, ErrInvalid
	}
	limits := q.policy.Limits
	if _, err := contextFiles.Files.Capture(limits.Files); err != nil {
		return Source[L]{}, err
	}
	if len(input.Units) > limits.Files.MaxFiles || len(input.Groups) > limits.MaxLinks || len(input.Relations) > limits.MaxLinks {
		return Source[L]{}, ErrLimit
	}
	catalogue := map[ast.TreeID]representation.AST{}
	for _, tree := range input.Snapshot.Inputs() {
		info, _ := tree.Inspect()
		catalogue[info.ID] = tree
	}
	filesByID := map[files.ID]files.File{}
	links := map[files.ID]files.SyntaxLink{}
	for _, file := range contextFiles.Files.Entries {
		filesByID[file.ID] = file
	}
	for _, link := range contextFiles.Links {
		links[link.File] = link
	}
	units, trees := map[files.ID]bool{}, map[ast.TreeID]files.ID{}
	nodes, edges, bytes, work := limits.MaxNodes, limits.MaxLinks, limits.MaxBytes, limits.MaxWork
	take := func(left *int, n int) bool {
		if n < 0 || n > *left {
			return false
		}
		*left -= n
		return true
	}
	owned := input
	owned.Units = make([]Unit[L], len(input.Units))
	for i, unit := range input.Units {
		file, exists := filesByID[unit.File]
		link := links[unit.File]
		if !exists || units[unit.File] || unit.Path != file.Path || unit.Source != link.Source || unit.Syntax.Abstract.ID != link.Tree {
			return Source[L]{}, ErrInvalid
		}
		if !take(&nodes, len(unit.Syntax.Nodes)) || !take(&work, len(unit.Syntax.Nodes)) {
			return Source[L]{}, ErrLimit
		}
		for _, node := range unit.Syntax.Nodes {
			if !take(&edges, len(node.Children)+len(node.Syntax)+len(node.Origin.Spans)+len(node.Origin.Nodes)) || !take(&bytes, len(node.Kind)) {
				return Source[L]{}, ErrLimit
			}
			for _, span := range node.Origin.Spans {
				if !take(&bytes, len(span.Source)) {
					return Source[L]{}, ErrLimit
				}
			}
			for _, origin := range node.Origin.Nodes {
				if !take(&bytes, len(origin.Tree)) {
					return Source[L]{}, ErrLimit
				}
			}
			if node.Text != nil && !take(&bytes, len(*node.Text)) {
				return Source[L]{}, ErrLimit
			}
		}
		packed, err := representation.PackCST(q.definition, unit.Syntax)
		if err != nil {
			return Source[L]{}, err
		}
		concrete, err := representation.UnpackCST[L](q.definition, packed)
		if err != nil {
			return Source[L]{}, err
		}
		if !representation.SameAST(q.definition, catalogue[link.Tree], concrete) {
			return Source[L]{}, ErrInvalid
		}
		// Structural admission above rules out cycles. Bound traversal depth/work too.
		type visit struct {
			id    cst.NodeID
			depth int
		}
		stack := []visit{{concrete.Root, 1}}
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if v.depth > limits.MaxDepth || !take(&work, 1) {
				return Source[L]{}, ErrLimit
			}
			for _, child := range concrete.Nodes[v.id-1].Children {
				stack = append(stack, visit{child, v.depth + 1})
			}
		}
		text, err := q.policy.Print(q.ctx, concrete)
		if err != nil {
			return Source[L]{}, err
		}
		if text != file.Text {
			return Source[L]{}, ErrInvalid
		}
		owned.Units[i] = Unit[L]{unit.File, unit.Path, unit.Source, concrete}
		units[unit.File] = true
		trees[link.Tree] = unit.File
	}
	contains := func(site ast.NodeRef) bool {
		tree, exists := catalogue[site.Tree]
		if !exists {
			return false
		}
		info, _ := tree.Inspect()
		return site.Node > 0 && int(site.Node) <= info.Nodes
	}
	groups := map[GroupID]bool{}
	for _, group := range input.Groups {
		if group.ID == "" || groups[group.ID] || !group.Kind.Valid() || len(group.Members) == 0 {
			return Source[L]{}, ErrInvalid
		}
		if !take(&edges, len(group.Members)+len(group.Sites)) || !take(&bytes, len(group.ID)+len(group.Kind.Owner)+len(group.Kind.Name)) {
			return Source[L]{}, ErrLimit
		}
		members := map[files.ID]bool{}
		for _, file := range group.Members {
			if !units[file] || members[file] {
				return Source[L]{}, ErrInvalid
			}
			members[file] = true
			if !take(&bytes, len(file)) {
				return Source[L]{}, ErrLimit
			}
		}
		for _, site := range group.Sites {
			if !contains(site) || !members[trees[site.Tree]] {
				return Source[L]{}, ErrInvalid
			}
			if !take(&bytes, len(site.Tree)) {
				return Source[L]{}, ErrLimit
			}
		}
		groups[group.ID] = true
	}
	for _, relation := range input.Relations {
		if !relation.Kind.Valid() || !units[relation.From] || !units[relation.To] {
			return Source[L]{}, ErrInvalid
		}
		if !take(&bytes, len(relation.Kind.Owner)+len(relation.Kind.Name)+len(relation.From)+len(relation.To)) {
			return Source[L]{}, ErrLimit
		}
		if relation.At != nil {
			if !contains(*relation.At) || trees[relation.At.Tree] != relation.From {
				return Source[L]{}, ErrInvalid
			}
			if !take(&bytes, len(relation.At.Tree)) {
				return Source[L]{}, ErrLimit
			}
		}
	}
	owned.Groups, owned.Relations = cloneGroups(input.Groups), cloneRelations(input.Relations)
	// Call the language guard on a separate owned copy. It cannot mutate publication.
	sealed := Source[L]{&state[L]{q.identity, owned, abstract}}
	guardInput := Publication[L]{owned.Snapshot, owned.Analysis, sealed.Units(), sealed.Groups(), sealed.Relations()}
	if err := q.policy.Check(q.ctx, guardInput); err != nil {
		return Source[L]{}, err
	}
	if err := q.ctx.Err(); err != nil {
		return Source[L]{}, err
	}
	if !q.active.Load() {
		return Source[L]{}, ErrStale
	}
	return sealed, nil
}

// Envelope is a sealed heterogeneous Source, separate from tree representation
// to keep semantic capture independent of concrete Source.
type sourceEvidence interface {
	Snapshot() semantic.Snapshot
	Evidence() semantic.Result
}

type Envelope struct {
	definition representation.Language
	value      sourceEvidence
}

func Pack[L any](definition representation.Language, value Source[L]) (Envelope, error) {
	if err := representation.CheckMarker[L](definition); err != nil {
		return Envelope{}, err
	}
	if !value.Valid() || !definition.Compatible(value.Snapshot().Language()) {
		return Envelope{}, ErrInvalid
	}
	return Envelope{definition, value}, nil
}
func Unpack[L any](definition representation.Language, value Envelope) (Source[L], error) {
	if err := representation.CheckMarker[L](definition); err != nil {
		return Source[L]{}, err
	}
	if !definition.Compatible(value.definition) {
		return Source[L]{}, representation.ErrLanguageMismatch
	}
	result, ok := value.value.(Source[L])
	if !ok || !result.Valid() {
		return Source[L]{}, ErrInvalid
	}
	return result, nil
}
func (e Envelope) Language() representation.Language { return e.definition }
func (e Envelope) Valid() bool                       { return e.value != nil }

// EqualFiles compares captured physical context, never semantic equivalence.
func EqualFiles(a, b files.Files) bool {
	if len(a.Entries) != len(b.Entries) {
		return false
	}
	byID := map[files.ID]files.File{}
	for _, f := range a.Entries {
		byID[f.ID] = f
	}
	for _, f := range b.Entries {
		if actual, ok := byID[f.ID]; !ok || !reflect.DeepEqual(actual, f) {
			return false
		}
		delete(byID, f.ID)
	}
	return true
}

func (e Envelope) Snapshot() semantic.Snapshot {
	if !e.Valid() {
		return semantic.Snapshot{}
	}
	return e.value.Snapshot()
}
func (e Envelope) Evidence() semantic.Result {
	if !e.Valid() {
		return semantic.Result{}
	}
	return e.value.Evidence()
}
func (e Envelope) SemanticsView() representation.Semantics {
	graph, _ := e.Evidence().Semantics()
	return graph
}
