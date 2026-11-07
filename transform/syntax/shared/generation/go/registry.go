package generation

import (
	"context"
	"reflect"
	"slices"
	"sync"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// References identify registrations in one registry, not names or claimed IDs.
type SourceRef struct {
	owner *Registry
	id    uint64
}
type TargetRef struct {
	owner *Registry
	id    uint64
}
type MappingRef struct {
	owner *Registry
	id    uint64
}
type BindingRef struct {
	Provider MappingRef
	ID       string
}
type Replacement struct {
	Selector       g.Selector
	Expected, With BindingRef
}
type Selection struct {
	Source   SourceRef
	Target   TargetRef
	Mappings []MappingRef
	Replace  []Replacement
}
type Request struct {
	Input     transform.Input
	To        representation.Language
	Roots     []ast.NodeRef
	Output    ast.TreeID
	Selection Selection
	Limits    g.Limits
}

// Registry is additive, safe for concurrent registration and invocation, and
// must not be copied after use. Callbacks never run under its lock. Private
// erased storage is accessible only through typed adapters checked before use.
type Registry struct {
	mu       sync.RWMutex
	next     uint64
	sources  map[uint64]sourceEntry
	targets  map[uint64]targetEntry
	mappings map[uint64]mappingEntry
}
type sourceEntry struct {
	language representation.Language
	api      reflect.Type
	contract transform.Contract
	provider any
}
type targetEntry struct {
	from, to representation.Language
	api      reflect.Type
	provider any
}
type mappingEntry struct {
	from, to             representation.Language
	sourceAPI, targetAPI reflect.Type
	provider             any
	adapter              func(Request, sourceEntry, targetEntry, []mappingEntry) (transform.ASTTransformer, error)
}

func RegisterSource[F any, S g.SourceSession[F]](r *Registry, language representation.Language, p g.SourceProvider[F, S]) (SourceRef, error) {
	if r == nil || representation.IsNil(p) || representation.CheckMarker[F](language) != nil {
		return SourceRef{}, g.ErrInvalid
	}
	contract := p.Contract()
	if err := semantic.ValidateEvidenceRequirements(contract.Evidence); err != nil {
		return SourceRef{}, err
	}
	contract.Evidence = slices.Clone(contract.Evidence)
	for i := range contract.Evidence {
		contract.Evidence[i] = contract.Evidence[i].Clone()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.sources == nil {
		r.sources = map[uint64]sourceEntry{}
	}
	r.next++
	r.sources[r.next] = sourceEntry{language, reflect.TypeFor[S](), contract, p}
	return SourceRef{r, r.next}, nil
}
func RegisterTarget[F, T any, N g.Native[F, T]](r *Registry, from, to representation.Language, p g.TargetProvider[F, T, N]) (TargetRef, error) {
	if r == nil || representation.IsNil(p) || representation.CheckMarker[F](from) != nil || representation.CheckMarker[T](to) != nil {
		return TargetRef{}, g.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.targets == nil {
		r.targets = map[uint64]targetEntry{}
	}
	r.next++
	r.targets[r.next] = targetEntry{from, to, reflect.TypeFor[N](), p}
	return TargetRef{r, r.next}, nil
}
func RegisterMapping[F, T any, S g.SourceSession[F], N g.Native[F, T]](r *Registry, from, to representation.Language, p g.MappingProvider[F, T, S, N]) (MappingRef, error) {
	if r == nil || representation.IsNil(p) || representation.CheckMarker[F](from) != nil || representation.CheckMarker[T](to) != nil {
		return MappingRef{}, g.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.mappings == nil {
		r.mappings = map[uint64]mappingEntry{}
	}
	r.next++
	r.mappings[r.next] = mappingEntry{from, to, reflect.TypeFor[S](), reflect.TypeFor[N](), p, adapter[F, T, S, N]}
	return MappingRef{r, r.next}, nil
}

// Transformer captures a checked explicit selection. It is bound to the exact
// source snapshot; registering it as a reusable route would not authorize other
// snapshots. Generate is the normal one-shot entry.
func (r *Registry) Transformer(request Request) (transform.ASTTransformer, error) {
	if r == nil || !request.Input.Source.Valid() || request.To.Validate() != nil || !request.Limits.Valid() || request.Output == "" || len(request.Roots) == 0 || len(request.Selection.Mappings) == 0 {
		return nil, g.ErrInvalid
	}
	if len(request.Roots) > request.Limits.References || len(request.Selection.Mappings) > request.Limits.Plans || len(request.Selection.Replace) > request.Limits.Plans {
		return nil, g.ErrLimit
	}
	request.Input.Evidence = slices.Clone(request.Input.Evidence)
	request.Roots = slices.Clone(request.Roots)
	request.Selection.Mappings = slices.Clone(request.Selection.Mappings)
	request.Selection.Replace = slices.Clone(request.Selection.Replace)
	trees := map[ast.TreeID]int{}
	for _, tree := range request.Input.Source.Inputs() {
		info, err := tree.Inspect()
		if err != nil {
			return nil, err
		}
		if info.ID == request.Output {
			return nil, g.ErrScope
		}
		trees[info.ID] = info.Nodes
	}
	seenRoots := map[ast.NodeRef]bool{}
	for _, root := range request.Roots {
		if root.Node == 0 || int(root.Node) > trees[root.Tree] || seenRoots[root] {
			return nil, g.ErrScope
		}
		seenRoots[root] = true
	}
	sref, tref := request.Selection.Source, request.Selection.Target
	r.mu.RLock()
	src, sok := r.sources[sref.id]
	target, tok := r.targets[tref.id]
	var mappings []mappingEntry
	seen := map[MappingRef]bool{}
	valid := sref.owner == r && tref.owner == r && sok && tok
	for _, ref := range request.Selection.Mappings {
		entry, exists := r.mappings[ref.id]
		if ref.owner != r || !exists || seen[ref] {
			valid = false
		}
		seen[ref] = true
		mappings = append(mappings, entry)
	}
	r.mu.RUnlock()
	if !valid {
		return nil, g.ErrScope
	}
	from := request.Input.Source.Language()
	if !from.Compatible(src.language) || !from.Compatible(target.from) || !request.To.Compatible(target.to) {
		return nil, g.ErrScope
	}
	for _, m := range mappings {
		if !from.Compatible(m.from) || !request.To.Compatible(m.to) || src.api != m.sourceAPI || target.api != m.targetAPI {
			return nil, g.ErrScope
		}
	}
	for _, replacement := range request.Selection.Replace {
		if !replacement.Selector.Valid() || !from.Compatible(replacement.Selector.Language()) || !seen[replacement.Expected.Provider] || !seen[replacement.With.Provider] || replacement.Expected.ID == "" || replacement.With.ID == "" || replacement.Expected == replacement.With {
			return nil, g.ErrInvalid
		}
	}
	// All adapters have the same checked F/T/S/N instantiation. Choosing one here
	// chooses no representation and captures no mapping-provider policy.
	t, err := mappings[0].adapter(request, src, target, mappings)
	if err != nil {
		return nil, err
	}
	return withDiagnostics(t, from, request.To), nil
}
func (r *Registry) Generate(ctx context.Context, request Request) (representation.ASTSet, error) {
	if ctx == nil {
		return representation.ASTSet{}, g.ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return representation.ASTSet{}, err
	}
	t, err := r.Transformer(request)
	if err != nil {
		return representation.ASTSet{}, err
	}
	return t.Transform(ctx, request.Input)
}
