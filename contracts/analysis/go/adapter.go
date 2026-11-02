package analysis

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Provider supplies language rules. It must be safe for concurrent calls when
// the application calls concurrently. Cancellation inside a synchronous analyzer
// requires provider cooperation; the adapter checks before and after invocation.
type Provider[L any] interface {
	language.LanguageProvider
	Analyze(context.Context, Request[L]) (Result, error)
}

// Analyzer is an opaque heterogeneous capability. Only NewAnalyzer constructs
// one, so dynamic registration cannot bypass typed capture/publication checks.
type Analyzer struct {
	syntax representation.Language
	contract Contract
	run func(context.Context, Input) (Result, error)
}

func NewAnalyzer[L any](syntax representation.Language, contract Contract, provider Provider[L]) (Analyzer, error) {
	if representation.IsNil(provider) { return Analyzer{}, ErrInvalidAnalyzer }
	if err := representation.CheckMarker[L](syntax); err != nil { return Analyzer{}, err }
	if err := validateContract(contract); err != nil { return Analyzer{}, err }
	id, name := provider.ID(), provider.Name()
	if id != syntax.ID() || strings.TrimSpace(name) == "" || !utf8.ValidString(name) {
		return Analyzer{}, fmt.Errorf("%w: provider identity does not match syntax", ErrInvalidAnalyzer)
	}
	contract = cloneContract(contract)
	result := Analyzer{syntax: syntax, contract: contract}
	result.run = func(ctx context.Context, input Input) (Result, error) {
		if ctx == nil { return Result{}, ErrInvalidRequest }
		if err := ctx.Err(); err != nil { return Result{}, err }
		request, err := captureRequest[L](ctx, syntax, contract, input)
		if err != nil { return Result{}, err }
		output, err := provider.Analyze(ctx, request)
		if err != nil { return Result{}, err }
		if err := ctx.Err(); err != nil { return Result{}, err }
		if output.state == nil { return Result{}, ErrInvalidResult }
		if output.state.owner != request.state.capture { return Result{}, ErrStaleResult }
		return output, nil
	}
	return result, nil
}

func (a Analyzer) ID() language.ID { return a.syntax.ID() }
func (a Analyzer) Name() string { return a.syntax.Name() }
func (a Analyzer) Language() representation.Language { return a.syntax }
func (a Analyzer) Contract() Contract { return cloneContract(a.contract) }
func (a Analyzer) Validate() error {
	if a.run == nil { return ErrInvalidAnalyzer }
	if err := a.syntax.Validate(); err != nil { return err }
	return validateContract(a.contract)
}
func (a Analyzer) Analyze(ctx context.Context, input Input) (Result, error) {
	if err := a.Validate(); err != nil { return Result{}, err }
	return a.run(ctx, input)
}

func captureRequest[L any](ctx context.Context, syntax representation.Language, contract Contract, input Input) (Request[L], error) {
	bad := func(reason string) (Request[L], error) {
		return Request[L]{}, fmt.Errorf("%w: %s", ErrInvalidRequest, reason)
	}
	limits := input.Limits
	if limits == (Limits{}) { limits = DefaultLimits() }
	if !limits.Valid() { return bad("invalid limits") }
	if len(input.Trees) == 0 { return bad("no input trees") }
	if len(input.Trees) > limits.MaxInputs || len(input.Requested) > limits.MaxRequested || len(input.Configuration) > limits.MaxConfigurationEntries {
		return Request[L]{}, ErrLimit
	}
	c := &capture{entry: input.Entry, contract: cloneContract(contract), limits: limits, catalogue: make(map[ast.TreeID]ast.NodeID)}
	bytes := limits.MaxConfigurationBytes
	if !take(&bytes, len(input.Entry.Tree)) { return Request[L]{}, ErrLimit }
	requested := input.Requested
	if len(requested) == 0 {
		requested = make([]semantics.Key, len(contract.Obligations))
		for i, obligation := range contract.Obligations { requested[i] = obligation.Key }
	}
	if len(requested) > limits.MaxRequested { return Request[L]{}, ErrLimit }
	seen := make(map[semantics.Key]bool)
	for _, key := range requested {
		if !key.Valid() || seen[key] { return bad("invalid or duplicate requested obligation") }
		if !keyBytes(&bytes, key) { return Request[L]{}, ErrLimit }
		seen[key] = true
	}
	clear(seen)
	for _, setting := range input.Configuration {
		if !setting.Key.Valid() || seen[setting.Key] { return bad("invalid or duplicate configuration key") }
		if !keyBytes(&bytes, setting.Key) || !take(&bytes, len(setting.Value)) { return Request[L]{}, ErrLimit }
		seen[setting.Key] = true
	}
	c.requested = slices.Clone(requested)
	c.configuration = cloneSettings(input.Configuration)
	state := &requestState[L]{capture: c, inputs: make([]ast.AST[L], 0, len(input.Trees))}
	nodes, inputBytes := limits.MaxInputNodes, limits.MaxInputBytes
	for _, envelope := range input.Trees {
		if err := ctx.Err(); err != nil { return Request[L]{}, err }
		if err := syntax.Check(envelope); err != nil { return Request[L]{}, err }
		value, err := representation.Unpack[L](syntax, envelope)
		if err != nil { return Request[L]{}, err }
		if _, found := c.catalogue[value.ID]; found { return bad("duplicate input tree ID") }
		if !take(&nodes, len(value.Nodes)) || !boundTree(value, &inputBytes) { return Request[L]{}, ErrLimit }
		c.catalogue[value.ID] = ast.NodeID(len(value.Nodes))
		state.inputs = append(state.inputs, value)
	}
	if !c.contains(c.entry) { return bad("entry is outside captured inputs") }
	return Request[L]{state: state}, nil
}

func boundTree[L any](input ast.AST[L], bytes *int) bool {
	if !take(bytes, len(input.ID)) { return false }
	for _, node := range input.Nodes {
		if !take(bytes, len(node.Kind)) { return false }
		for _, attribute := range node.Attributes {
			if !take(bytes, len(attribute.Name)) || !take(bytes, len(attribute.Text)) { return false }
		}
		for _, child := range node.Children { if !take(bytes, len(child.Role)) { return false } }
		for _, reference := range node.References {
			if !take(bytes, len(reference.Role)) || !take(bytes, len(reference.Target.Tree)) { return false }
		}
		for _, span := range node.Origin.Spans { if !take(bytes, len(span.Source)) { return false } }
		for _, origin := range node.Origin.Nodes { if !take(bytes, len(origin.Tree)) { return false } }
	}
	return true
}
