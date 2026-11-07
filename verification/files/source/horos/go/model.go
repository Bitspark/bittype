// Package horos supplies explicit, artifact-associated authoritative Horos queries.
// These are projection inputs, not proof outcomes or a second kind evaluator.
package horos

import (
	"context"
	"errors"

	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	bridge "github.com/Bitspark/bittype/internal/horosmodel"
	resolve "github.com/Bitspark/bittype/resolve/files/horos/provider/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

var ErrForeign = errors.New("foreign or missing Horos projection handle")

var frontend = bridge.Frontend[s.Language]{
	DeclarationPrefix: "horos/decl/", DescriptorPrefix: "horos/", Owner: "horos", ResolutionProfile: resolve.Profile,
	ResolveModel: func(root core.Value, deps []core.Value, limits decl.Limits) (bridge.Program, error) {
		return decl.ResolveWithLimits(root, deps, limits)
	},
	PinModel: decl.PinWithLimits,
}

type state struct {
	snapshot semantic.Snapshot
	receipt  semantic.Result
	program  *decl.Program
	decoded  *bridge.Decoded[s.Language]
}
type Model struct{ state *state }
type Declaration struct {
	owner *state
	value *decl.Declaration
}
type Value struct {
	owner *state
	value decl.TypeValue
}
type Kind struct {
	owner  *state
	value  decl.KindHandle
	module decl.Digest
}
type Judgment struct {
	owner   *state
	value   Value
	kind    Kind
	verdict decl.Verdict
}

func Evidence() semantic.EvidenceRequirement {
	return semantic.EvidenceRequirement{Operation: resolution.Kind(), Profile: resolve.Profile,
		Reports: resolve.Contract().Semantic.SemanticsRequirements, Semantics: true}
}

// Inspect explicitly reconstructs the owner-local Program from the exact captured
// collection, after admission of its resolution receipt. Prepare never calls it.
func Inspect(ctx context.Context, snapshot semantic.Snapshot, receipt semantic.Result) (Model, error) {
	if ctx == nil {
		return Model{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return Model{}, err
	}
	if _, err := semantic.AdmitEvidence(snapshot, []semantic.EvidenceRequirement{Evidence()}, []semantic.Result{receipt}); err != nil {
		return Model{}, err
	}
	if err := snapshot.CheckLimits(semantic.DefaultLimits()); err != nil {
		return Model{}, err
	}
	var trees []ast.AST[s.Language]
	for _, input := range snapshot.Inputs() {
		tree, err := representation.Unpack[s.Language](snapshot.Language(), input)
		if err != nil {
			return Model{}, err
		}
		trees = append(trees, tree)
	}
	decoded, err := frontend.DecodeCollection(trees, snapshot.Entry(), decl.DefaultLimits())
	if err != nil {
		return Model{}, err
	}
	program, err := decoded.Resolve(decl.DefaultLimits())
	if err != nil {
		return Model{}, err
	}
	if err := ctx.Err(); err != nil {
		return Model{}, err
	}
	return Model{&state{snapshot, receipt, program.(*decl.Program), decoded}}, nil
}

func (m Model) Valid() bool           { return m.state != nil }
func (m Model) Same(other Model) bool { return m.Valid() && m.state == other.state }
func (m Model) Snapshot() semantic.Snapshot {
	if !m.Valid() {
		return semantic.Snapshot{}
	}
	return m.state.snapshot
}
func (m Model) Receipt() semantic.Result {
	if !m.Valid() {
		return semantic.Result{}
	}
	return m.state.receipt
}
func (m Model) Root() decl.Digest {
	if !m.Valid() {
		return decl.Digest{}
	}
	return m.state.program.Root()
}
func (m Model) Declaration(id decl.DeclarationID) (Declaration, error) {
	if !m.Valid() {
		return Declaration{}, ErrForeign
	}
	d, ok := m.state.program.Declaration(id)
	if !ok {
		return Declaration{}, ErrForeign
	}
	return Declaration{m.state, d}, nil
}
func (d Declaration) Valid() bool  { return d.owner != nil && d.value != nil }
func (d Declaration) Model() Model { return Model{d.owner} }
func (d Declaration) Definition() *decl.Declaration {
	if !d.Valid() {
		return nil
	}
	copy := *d.value
	return &copy
}
func (d Declaration) Site() ast.NodeRef {
	if !d.Valid() {
		return ast.NodeRef{}
	}
	return d.owner.decoded.DeclarationSite(d.value.ID())
}
func (d Declaration) ParameterSite(index int) ast.NodeRef {
	if !d.Valid() {
		return ast.NodeRef{}
	}
	return d.owner.decoded.ParameterSite(decl.ParameterID{Declaration: d.value.ID(), Index: index})
}
func (d Declaration) ParameterKind(index int) (Kind, error) {
	if !d.Valid() || index < 0 || index >= len(d.value.Parameters()) {
		return Kind{}, ErrForeign
	}
	return d.Model().CompileKind(d.value.ID().Module, d.value.Parameters()[index].Kind.Value())
}
func (d Declaration) Value() (Value, error) {
	if !d.Valid() {
		return Value{}, ErrForeign
	}
	v, err := d.owner.program.Value(d.value.ID())
	if err != nil {
		return Value{}, err
	}
	return Value{d.owner, v}, nil
}
func (m Model) Proper(descriptor core.Value) (Value, error) {
	if !m.Valid() {
		return Value{}, ErrForeign
	}
	v, err := m.state.program.Proper(descriptor)
	if err != nil {
		return Value{}, err
	}
	return Value{m.state, v}, nil
}
func (m Model) CompileKind(module decl.Digest, expression core.Value) (Kind, error) {
	if !m.Valid() {
		return Kind{}, ErrForeign
	}
	k, err := m.state.program.CompileKind(module, expression)
	if err != nil {
		return Kind{}, err
	}
	return Kind{m.state, k, module}, nil
}

// Arrow uses the same canonical Horos kind syntax as any other request.
func (k Kind) Arrow(module decl.Digest) (Kind, error) {
	if k.owner == nil || module != k.module {
		return Kind{}, ErrForeign
	}
	v := k.value.Expression().Value()
	return (Model{k.owner}).CompileKind(module, core.NewTuple(core.NewAtom([]byte("horos/decl/arrow")), data.EncodeList([]core.Value{v}), v))
}
func (k Kind) Model() Model                  { return Model{k.owner} }
func (k Kind) Module() decl.Digest           { return k.module }
func (k Kind) Expression() *decl.Kind        { return k.value.Expression() }
func (v Value) Model() Model                 { return Model{v.owner} }
func (v Value) Expression() *decl.Expression { return v.value.Expression() }
func (v Value) Same(other Value) bool {
	return v.owner != nil && v.owner == other.owner && v.value == other.value
}
func (k Kind) Same(other Kind) bool {
	return k.owner != nil && k.owner == other.owner && k.value == other.value
}
func (m Model) Apply(callee Value, arguments ...Value) (Value, error) {
	if !m.Valid() || callee.owner != m.state {
		return Value{}, ErrForeign
	}
	args := make([]decl.TypeValue, len(arguments))
	for i, v := range arguments {
		if v.owner != m.state {
			return Value{}, ErrForeign
		}
		args[i] = v.value
	}
	v, err := m.state.program.Apply(callee.value, args)
	if err != nil {
		return Value{}, err
	}
	return Value{m.state, v}, nil
}
func (m Model) Satisfies(ctx context.Context, value Value, kind Kind) (Judgment, error) {
	if ctx == nil || !m.Valid() || value.owner != m.state || kind.owner != m.state {
		return Judgment{}, ErrForeign
	}
	if err := ctx.Err(); err != nil {
		return Judgment{}, err
	}
	v, err := m.state.program.Satisfies(value.value, kind.value)
	if err != nil {
		return Judgment{}, err
	}
	if err := ctx.Err(); err != nil {
		return Judgment{}, err
	}
	return Judgment{m.state, value, kind, v}, nil
}
func (j Judgment) Model() Model          { return Model{j.owner} }
func (j Judgment) Value() Value          { return j.value }
func (j Judgment) Kind() Kind            { return j.kind }
func (j Judgment) Verdict() decl.Verdict { return j.verdict }

func (v Value) Descriptor() (core.Value, error) {
	if v.owner == nil {
		return nil, ErrForeign
	}
	return v.owner.program.Descriptor(v.value)
}

type Field struct {
	Name     string
	Optional bool
	Type     *decl.Expression
	Site     ast.NodeRef
}

// Fields retains every field of an actual closed record, in canonical order.
// Other descriptors are outside this particular static mapping's source view.
func (d Declaration) Fields() ([]Field, error) {
	if !d.Valid() || d.value.Form() != decl.TypeDefinition {
		return nil, ErrForeign
	}
	e := d.value.Expression()
	if e == nil || e.Kind() != decl.Descriptor || e.Constructor() != "record" {
		return nil, errors.New("source body is not a closed record")
	}
	t, ok := e.Value().(core.Tuple)
	if !ok || t.Len() < 2 {
		return nil, ErrForeign
	}
	mode, ok := t.At(1).(core.Atom)
	if !ok || string(mode.Bytes()) != "closed" {
		return nil, errors.New("open records are outside this projection")
	}
	operands := e.Operands()
	fields := make([]Field, 0, t.Len()-2)
	for i := 2; i < t.Len(); i++ {
		field, ok := t.At(i).(core.Tuple)
		if !ok || field.Len() != 4 {
			return nil, ErrForeign
		}
		name, err := data.ReadText(field.At(1))
		if err != nil {
			return nil, err
		}
		flag, ok := field.At(2).(core.Atom)
		if !ok {
			return nil, ErrForeign
		}
		var ty *decl.Expression
		for _, op := range operands {
			path := op.Path()
			if len(path) == 2 && path[0] == i && path[1] == 3 {
				ty = op.Expression()
				break
			}
		}
		if ty == nil {
			return nil, ErrForeign
		}
		fields = append(fields, Field{Name: name, Optional: string(flag.Bytes()) == "opt", Type: ty, Site: d.owner.decoded.OperandSite(d.owner.decoded.BodySite(d.value.ID()), []int{i, 3})})
	}
	return fields, nil
}
