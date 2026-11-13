package semantic

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type snapshotState struct {
	syntax    representation.Language
	files     *files.Context
	entry     ast.NodeRef
	trees     []representation.AST
	catalogue map[ast.TreeID]ast.NodeID
}

// Snapshot identifies one admitted input collection. Identical addresses or
// bytes captured again are a different snapshot. Storage is private and immutable.
type Snapshot struct{ state *snapshotState }

func (s Snapshot) Valid() bool { return s.state != nil }
func (s Snapshot) LanguageID() language.ID {
	if !s.Valid() {
		return ""
	}
	return s.state.syntax.ID()
}

// Entry is the node that the capture names as its entry. Every capture has
// one, and this contract gives it no meaning: it is an address. It names no
// root of a native program, such as the packages of a Go build, the root file
// of a Rust crate, a translation-unit context of C++ or the input files of a
// TypeScript compiler. It limits no check: a provider answers for every
// captured tree, whichever holds the entry. A provider may give the entry a
// meaning of its own, as the horos providers do for the module that a request
// enters through.
func (s Snapshot) Entry() ast.NodeRef {
	if !s.Valid() {
		return ast.NodeRef{}
	}
	return s.state.entry
}

func normalizedLimits(limits Limits) (Limits, error) {
	if limits == (Limits{}) {
		limits = DefaultLimits()
	}
	if !limits.Valid() {
		return Limits{}, ErrInvalidRequest
	}
	return limits, nil
}

func boundInputs(trees []representation.AST, limits Limits) (map[ast.TreeID]ast.NodeID, error) {
	if len(trees) == 0 {
		return nil, ErrInvalidRequest
	}
	if len(trees) > limits.MaxInputs {
		return nil, ErrLimit
	}
	catalogue := make(map[ast.TreeID]ast.NodeID)
	nodes, entries, textBytes := limits.MaxInputNodes, limits.MaxInputEntries, limits.MaxInputBytes
	for _, tree := range trees {
		info, err := tree.Inspect()
		if err != nil {
			return nil, err
		}
		if _, found := catalogue[info.ID]; found {
			return nil, ErrInvalidRequest
		}
		if !take(&nodes, info.Nodes) || !take(&entries, info.Entries) || !take(&textBytes, info.TextBytes) {
			return nil, ErrLimit
		}
		catalogue[info.ID] = ast.NodeID(info.Nodes)
	}
	for _, tree := range trees {
		if err := tree.CheckReferences(catalogue); err != nil {
			return nil, err
		}
	}
	return catalogue, nil
}

// Capture bounds and retains immutable envelopes before any checker callback.
func Capture(ctx context.Context, checker check.Checker, input Input) (Snapshot, error) {
	if ctx == nil {
		return Snapshot{}, ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return Snapshot{}, err
	}
	if err := checker.Validate(); err != nil {
		return Snapshot{}, err
	}
	limits, err := normalizedLimits(input.Limits)
	if err != nil {
		return Snapshot{}, err
	}
	if len(input.Trees) > limits.MaxInputs || len(input.Entry.Tree) > limits.MaxInputBytes {
		return Snapshot{}, ErrLimit
	}
	trees := slices.Clone(input.Trees)
	catalogue, err := boundInputs(trees, limits)
	if err != nil {
		return Snapshot{}, err
	}
	if input.Entry.Tree == "" || input.Entry.Node == 0 || input.Entry.Node > catalogue[input.Entry.Tree] {
		return Snapshot{}, ErrInvalidRequest
	}
	var contextFiles *files.Context
	if input.Files != nil {
		captured, err := input.Files.Capture(files.Limits{MaxFiles: limits.MaxInputs, MaxBytes: limits.MaxInputBytes, MaxLinks: limits.MaxInputs}, catalogue)
		if err != nil {
			return Snapshot{}, err
		}
		contextFiles = &captured
	}
	for _, tree := range trees {
		if err := ctx.Err(); err != nil {
			return Snapshot{}, err
		}
		if err := checker.Check(tree); err != nil {
			return Snapshot{}, err
		}
	}
	if err := ctx.Err(); err != nil {
		return Snapshot{}, err
	}

	return Snapshot{state: &snapshotState{syntax: checker.Language(), entry: input.Entry, trees: trees, catalogue: catalogue, files: contextFiles}}, nil
}

// Dependency is an explicit input receipt and the consumer's exact requirements.
// Configuration compares keys and bytes, never inferred semantic equivalence.
type Dependency struct {
	Result        Result
	Operation     semantics.Key
	Profile       string
	Configuration []Setting
	Reports       []Requirement
}

func (d Dependency) reportsEstablished() bool {
	for _, required := range d.Reports {
		report, found := d.Result.Report(required.Obligation)
		if !found || report.Status != Complete || report.Verdict != required.Verdict {
			return false
		}
	}
	return d.Result.Valid()
}

type metadataBudget struct{ requested, settings, bytes int }

func newMetadataBudget(limits Limits) metadataBudget {
	return metadataBudget{requested: limits.MaxRequested, settings: limits.MaxConfigurationEntries, bytes: limits.MaxConfigurationBytes}
}

func validateSettings(settings []Setting, budget *metadataBudget) error {
	if !take(&budget.settings, len(settings)) {
		return ErrLimit
	}
	seen := make(map[semantics.Key]bool)
	for _, setting := range settings {
		if !setting.Key.Valid() || seen[setting.Key] {
			return ErrInvalidRequest
		}
		if !keyBytes(&budget.bytes, setting.Key) || !take(&budget.bytes, len(setting.Value)) {
			return ErrLimit
		}
		seen[setting.Key] = true
	}
	return nil
}

// ValidateDependencySpecification validates consumer-owned metadata before copying.
func ValidateDependencySpecification(operation semantics.Key, profile string, settings []Setting, reports []Requirement) error {
	budget := newMetadataBudget(DefaultLimits())
	return validateDependencySpecification(operation, profile, settings, reports, &budget)
}

func validateDependencySpecification(operation semantics.Key, profile string, settings []Setting, reports []Requirement, budget *metadataBudget) error {
	if !keyBytes(&budget.bytes, operation) || !take(&budget.bytes, len(profile)) || !take(&budget.requested, len(reports)) {
		return ErrLimit
	}
	if !operation.Valid() || strings.TrimSpace(profile) == "" || !utf8.ValidString(profile) {
		return ErrInvalidOperation
	}
	if err := validateSettings(settings, budget); err != nil {
		return err
	}
	seen := make(map[semantics.Key]bool)
	for _, report := range reports {
		if !keyBytes(&budget.bytes, report.Obligation) {
			return ErrLimit
		}
		if !report.Obligation.Valid() || seen[report.Obligation] || report.Verdict > Rejected {
			return ErrInvalidOperation
		}
		seen[report.Obligation] = true
	}
	return nil
}

func admitDependency(source *snapshotState, d Dependency) (Dependency, error) {
	if !d.Result.Valid() || d.Result.state.owner.source != source || d.Result.state.owner.operation != d.Operation || d.Result.state.owner.contract.Profile != d.Profile {
		return Dependency{}, ErrInvalidRequest
	}
	actual := d.Result.state.owner.configuration
	if len(actual) != len(d.Configuration) {
		return Dependency{}, ErrInvalidRequest
	}
	for _, expected := range d.Configuration {
		found := false
		for _, setting := range actual {
			if expected.Key == setting.Key && bytes.Equal(expected.Value, setting.Value) {
				found = true
				break
			}
		}
		if !found {
			return Dependency{}, ErrInvalidRequest
		}
	}
	for _, required := range d.Reports {
		advertised := false
		for _, obligation := range d.Result.state.owner.contract.Obligations {
			advertised = advertised || obligation.Key == required.Obligation
		}
		if !advertised {
			return Dependency{}, ErrInvalidRequest
		}
	}
	d.Configuration = cloneSettings(d.Configuration)
	d.Reports = slices.Clone(d.Reports)
	return d, nil
}

// Operation is the typed invocation bridge shared by specialized capabilities.
// Its kind/profile are explicit compatibility promises, not semantic discovery.
type Operation struct {
	syntax   representation.Language
	contract Contract
	kind     semantics.Key
	run      func(context.Context, Snapshot, Options, []Dependency) (Result, error)
}

func NewOperation[L any](syntax representation.Language, contract Contract, kind semantics.Key, provider language.Provider, run func(context.Context, Request[L]) (Result, error)) (Operation, error) {
	if representation.IsNil(provider) || run == nil || !kind.Valid() {
		return Operation{}, ErrInvalidOperation
	}
	if err := representation.CheckMarker[L](syntax); err != nil {
		return Operation{}, err
	}
	if err := contract.Validate(); err != nil {
		return Operation{}, err
	}
	if provider.ID() != syntax.ID() || strings.TrimSpace(provider.Name()) == "" || !utf8.ValidString(provider.Name()) {
		return Operation{}, ErrInvalidOperation
	}
	contract = contract.Clone()
	op := Operation{syntax: syntax, contract: contract, kind: kind}
	op.run = func(ctx context.Context, source Snapshot, options Options, dependencies []Dependency) (Result, error) {
		if ctx == nil || !source.Valid() {
			return Result{}, ErrInvalidRequest
		}
		if err := ctx.Err(); err != nil {
			return Result{}, err
		}
		if !syntax.Compatible(source.state.syntax) {
			return Result{}, representation.ErrLanguageMismatch
		}
		limits, err := normalizedLimits(options.Limits)
		if err != nil {
			return Result{}, err
		}
		if err := source.CheckLimits(limits); err != nil {
			return Result{}, err
		}
		if len(options.Requested) > limits.MaxRequested || len(dependencies) > limits.MaxInputs {
			return Result{}, ErrLimit
		}
		requested := options.Requested
		if len(requested) == 0 {
			if len(contract.Obligations) > limits.MaxRequested {
				return Result{}, ErrLimit
			}
			requested = make([]semantics.Key, 0, len(contract.Obligations))
			for _, obligation := range contract.Obligations {
				requested = append(requested, obligation.Key)
			}
		}
		budget := newMetadataBudget(limits)
		if !take(&budget.requested, len(requested)) {
			return Result{}, ErrLimit
		}
		if err := validateSettings(options.Configuration, &budget); err != nil {
			return Result{}, err
		}
		seen := make(map[semantics.Key]bool)
		for _, key := range requested {
			if !key.Valid() || seen[key] {
				return Result{}, ErrInvalidRequest
			}
			if !keyBytes(&budget.bytes, key) {
				return Result{}, ErrLimit
			}
			seen[key] = true
		}
		// Bound the entire invocation before copying caller/dependency metadata.
		for _, dependency := range dependencies {
			if err := validateDependencySpecification(dependency.Operation, dependency.Profile, dependency.Configuration, dependency.Reports, &budget); err != nil {
				return Result{}, err
			}
		}
		c := &capture{source: source.state, operation: kind, contract: contract, requested: slices.Clone(requested), configuration: cloneSettings(options.Configuration), limits: limits}
		for _, dependency := range dependencies {
			admitted, err := admitDependency(source.state, dependency)
			if err != nil {
				return Result{}, err
			}
			c.dependencies = append(c.dependencies, admitted)
		}
		state := &requestState[L]{capture: c}
		for _, envelope := range source.state.trees {
			if err := ctx.Err(); err != nil {
				return Result{}, err
			}
			tree, err := representation.Unpack[L](syntax, envelope)
			if err != nil {
				return Result{}, err
			}
			state.inputs = append(state.inputs, tree)
		}
		result, err := run(ctx, Request[L]{state: state})
		if err != nil {
			return Result{}, err
		}
		if err := ctx.Err(); err != nil {
			return Result{}, err
		}
		if !result.Valid() {
			return Result{}, ErrInvalidResult
		}
		if result.state.owner != c {
			return Result{}, ErrStaleResult
		}
		return result, nil
	}
	return op, nil
}

func (o Operation) Language() representation.Language { return o.syntax }
func (o Operation) Contract() Contract                { return o.contract.Clone() }
func (o Operation) Validate() error {
	if o.run == nil {
		return ErrInvalidOperation
	}
	return nil
}
func (o Operation) Run(ctx context.Context, source Snapshot, options Options, dependencies []Dependency) (Result, error) {
	if err := o.Validate(); err != nil {
		return Result{}, err
	}
	return o.run(ctx, source, options, dependencies)
}

func (r Request[L]) Dependency(index int) (Result, bool) {
	if r.state == nil || index < 0 || index >= len(r.state.capture.dependencies) {
		return Result{}, false
	}
	return r.state.capture.dependencies[index].Result, true
}

func (r Result) Operation() semantics.Key {
	if !r.Valid() {
		return semantics.Key{}
	}
	return r.state.owner.operation
}
func (r Result) Snapshot() Snapshot {
	if !r.Valid() {
		return Snapshot{}
	}
	return Snapshot{state: r.state.owner.source}
}

// Files returns the captured correspondence; nil means syntax-only input.
func (s Snapshot) Files() *files.Context {
	if !s.Valid() || s.state.files == nil {
		return nil
	}
	c := s.state.files.Clone()
	return &c
}
func (s Snapshot) Language() representation.Language {
	if !s.Valid() {
		return representation.Language{}
	}
	return s.state.syntax
}
func (s Snapshot) Same(other Snapshot) bool { return s.Valid() && s.state == other.state }
func (r Request[L]) Files() *files.Context {
	if r.state == nil {
		return nil
	}
	return (Snapshot{state: r.state.capture.source}).Files()
}

// CaptureSettings bounds and copies caller-owned configuration without interpreting it.
func CaptureSettings(settings []Setting, limits Limits) ([]Setting, error) {
	normalized, err := normalizedLimits(limits)
	if err != nil {
		return nil, err
	}
	budget := newMetadataBudget(normalized)
	if err := validateSettings(settings, &budget); err != nil {
		return nil, err
	}
	return cloneSettings(settings), nil
}

// CheckLimits bounds all retained input context without running a provider.
func (s Snapshot) CheckLimits(limits Limits) error {
	if !s.Valid() {
		return ErrInvalidRequest
	}
	l, err := normalizedLimits(limits)
	if err != nil {
		return err
	}
	if _, err := boundInputs(s.state.trees, l); err != nil {
		return err
	}
	if s.state.files != nil {
		if _, err := s.state.files.Capture(files.Limits{MaxFiles: l.MaxInputs, MaxBytes: l.MaxInputBytes, MaxLinks: l.MaxInputs}, s.state.catalogue); err != nil {
			return err
		}
	}
	return nil
}
