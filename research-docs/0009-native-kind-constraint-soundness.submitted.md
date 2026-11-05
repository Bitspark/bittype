# Research: Native kind constraint soundness across representations and packages

**ID:** 0009
**Date:** 5 November 2026
**Status:** reviewed
**Reviewed:** https://github.com/Bitspark/bittype/issues/85#issuecomment-6664161148
**Owner:** codex-hydrated-type-design
**Issue:** https://github.com/Bitspark/bittype/issues/85
**Companion consultation:** 0008, https://github.com/Bitspark/bittype/issues/84
**Inspected baseline:** bittype f2179283150b1e7344844fc97503ca26fe40795c

## Question

What precisely can it mean to preserve a canonical type's generic kind constraints
in TypeScript, Go, Rust and C++, when consumers can choose native representations
and generate packages separately? What is the smallest sound, useful encoding
domain, and is the proposed source-witness/membership-evidence approach the right
one at all?

This is a conceptual, documentation-only consultation before implementation. We
want counterexamples and alternatives, not confirmation of the sketches. No
runtime bindings, adapters, value validators, generated calls or stream behavior
are in scope. No compiler runs or implementation are requested from this
consultation. Source snippets below are inspected code; proposed target snippets
are not claimed to compile or enforce their intended contract.

The paired consultation 0008 owns the agnostic generation protocol: elaborating
one source binder/application into target parameters/arguments, coherent
Use/Publish/Define, and declaration ownership. This consultation owns the
semantic obligations that such a protocol must carry and the real limits of
candidate target encodings. Please challenge the overall conceptual approach
even when that affects the companion's protocol. The two independent answers
will be reconciled before any implementation.

## Context

### Vocabulary and system boundary

There are two independent Go modules:

- **horos** owns canonical type declarations, module imports, generic application,
  kinds and semantic judgments. Ground types describe inert **ontos** data
  values. The declaration language also admits a `wire<T>` type describing a
  communication capability. Ground value conformance stays data-only.
- **bittype** owns generic AST/CST carriers, provider contracts, an explicitly
  configured syntax engine, Horos frontend tooling and native type generation.
  AST means abstract syntax tree; CST means a concrete syntax tree carrying
  source spelling. Language-specific Go packages own their tree vocabulary.

Here **ground** names the inert ontos value/descriptor domain. It is not a synonym
for every closed proper type: **closed** means no free generic parameters, and
`wire<text>` is closed and proper but is outside ground Horos. Extraction to a
ground descriptor additionally requires a closed proper expression, established
data membership, materialization and ordinary ground descriptor well-formedness.

The generator is implemented in Go and produces target ASTs, then target CSTs,
then source text. TypeScript, Go, Rust and C++ are output languages. Generation
does not generate operational behavior. A library may declare Function, Stream
or Cell and supply a native type mapping, but Wire, hydration, the engine and
the generic generation coordinator must not learn those names or semantics.

An **open** declaration still has generic parameters. A **proper type** is a
type-level value describing runtime values. A **constructor** accepts type-level
arguments and returns a proper type or another constructor. A **kind** describes
admitted proper types or constructors, not runtime inhabitants. A **sort** is
the coarse shape `P` (proper) or an arrow of sorts. It is not a kind constraint.
An **association** below is a selected canonical-to-native type mapping, not a
runtime conversion. A **profile** states the output's promised domain and API.

An analysis **receipt** binds provider-owned facts to an exact captured input.
The agnostic semantic graph stores those facts but does not interpret Horos
kinds. A **source witness** is a proposed compile-time identity for a canonical
type expression. **Membership evidence** is a proposed native encoding of an
already-established source judgment and its assumptions; it is not a proof that
arbitrary payload values are valid. These are design terms, not existing APIs.

### Existing direction and unresolved delta

Prior research 0005 selected general structural kinds. Research 0006 moved the
standard prelude into ordinary Horos source. Research 0007 selected type-owned
generation mappings, exact consumer overrides and separate Use/Publish/Define.
We are not reopening those ownership choices as a routine rewrite. We are asking
what they imply for the unresolved native enforcement design.

The main documents are `docs/native-kind-constraints.md` and
`docs/type-generation-plugins.md`. Their relevant substance is reproduced here;
the expert does not need repository access. Main includes their alignment and
composition clarification from PRs 80-83. Neither plugin SPI nor constraint
encodings are implemented yet. The existing pagination generators use a private
source-specific dialect interface, which the plugin design will replace.

Current native generation supports text, bool, lists, named closed records,
optional fields and ordinary proper generic parameters/applications. It retains
open generic declarations. It rejects unsupported exports rather than silently
omitting them. Integer domains, Wire payload mappings, residual constructor
parameters/results, open records and unused proper parameters are outside this
implemented representation slice. We need an explicit useful initial enforcement
domain, not a claim that all Horos types already have native mappings.
The current native generator also rejects recursive declarations. That limitation
is separate from Horos's existing support for recursive kind judgments.

### Source semantics: what a kind actually classifies

Horos revision 4 uses finite simple sorts:

```text
Sort = P | (Sort, ...) -> Sort

KindDefinition = parameters with sorts + result sort + kind expression
TypeDefinition = parameters with kind contracts + result kind + type expression
```

Kinds have ordinary references and applications, top/bottom at a sort,
union/intersection, primitive node patterns with selected metadata and child
sequence patterns, guarded greatest-fixed-point recursion, and constructor
arrow contracts. Kind unions range over descriptions; they are not value unions.
The source language's primitive observation vocabulary is finite. There are no
arbitrary callback predicates, subtree-equality conditions or semantic inference
from every possible runtime inhabitant hidden in the kind model.

Representative standard definitions, schematically:

```text
List<K>       = node(list-of, {}, item(K))
Record<K>     = node(record, {}, star(item(K)))
ClosedRecord<K> = node(record, {mode: closed}, star(item(K)))
Map<K,V>      = node(map-of, {}, concat(item(K), item(V)))
Wire<K>       = node(wire, {}, item(K))
UnionOf<K>    = node(union, {}, star(item(K)))

Data = nu D.
    Any | AtomFamily | Scalar | Exact
  | Tuple<D> | Compound<D> | Seq<D>
  | List<D> | Set<D> | Map<D,D>
  | ClosedRecord<D> | OpenRecord<D>
  | UnionOf<D> | AndOf<D> | Not<D>
```

`node` describes an admitted primitive's form, metadata and the whole ordered
sequence of its designated type children. `item(K)` consumes one child;
`star(item(K))` admits any finite number. `nu` denotes a greatest fixed point.
Every recursive cycle descends through a node's type child. Ground `rec`/`ref`
become lexical backedges in a finite observation graph; its unfolding can be
infinite. An `exact` literal has no type children even if its data resembles a
type declaration. Optional record fields participate. Open records contribute
an implicit ground `any` child for unspecified fields.
The uppercase prelude kind `Any` recognizes the one ground `any` descriptor form;
it is not kind top. Its occurrence in Data therefore does not admit every type.

All prelude names, including Data, are ordinary imported kind definitions. A
renamed definition follows its resolved identity; a user's unrelated local
kind named Data has no privileged interpretation. `any` is a ground data type,
not kind top. Kind bottom has no types; an empty value union is still a proper
type description. A record with a Wire field is outside Data even if some other
field makes its runtime values impossible.

Aliases/applications normalize to structural observations for membership, while
declaration identities remain available for mapping selection. Two source
descriptions can have different kind membership and the same selected native
payload. For example, the kind `h.Text | h.Bool` contains a text description or
a boolean description; it does not contain a type-level union description whose
root is `union`. A native `string | boolean` bound does not preserve that fact.

Known constructors retain checked bodies and lexical environments. To establish
`F : A -> B`, where A and B are input and result kinds, Horos first establishes
that every type admitted by A satisfies F's original parameter contract
(`A <= original-domain`), then establishes B under fresh rigid input assumptions. A narrower query
cannot repair an invalid open declaration. Abstract constructors expose only
their established guarantees. Ignored arguments must still satisfy their
contracts. Results distinguish Proven, Refuted, Unproved, Invalid and Exhausted;
only established facts may justify output. A missing fact is not a negative
fact, and neither missing nor exhausted analysis admits a permissive encoding.

### Shared examples

These source excerpts assume explicit imports of `horos.kinds` as `h` and the
library module supplying `Cell` as `cell`, each with its ordinary module pin.
The syntax is Horos readable source; `nu` and kind annotations are compile-time.

```text
type Page<T> = {
  items: list<T>;
  next?: text;
};
type DataPage<D: h.Data>: h.Data = Page<D>;
type Person = { name: text; active: bool; };
type People = DataPage<Person>;
type EditablePeople = Page<cell.Cell<Person>>;

kind Tree = nu {
  X = h.Text | h.Bool | h.List<X> | h.ClosedRecord<X>;
} in X;
type TreePage<T: Tree>: Tree = Page<T>;
type TreePeople = TreePage<Person>;

type Ignore<T: Tree> = text;
```

The same Page body can establish `Data -> Data` and `Tree -> Tree`; there is no
copied DataPage implementation. Tree is intentionally smaller than Data: it
does not admit union-root or Wire-root descriptions. The optional `next` field
is text and remains relevant. Ignore's ignored argument is still constrained.
Cell is an ordinary library declaration; if its body includes Wire it is not
Data. Unconstrained Page may still compose with a supported library Cell
representation. Choosing a function/channel/handle native type supplies no call
or streaming semantics.

Required scenarios, not executed test results:

| Scenario | Source fact / question |
| --- | --- |
| `TreePage<Person>` | Person's closed record of text/bool belongs to Tree. |
| `TreePage<list<Person>>` | A generic list membership implication must retain its prerequisites. |
| `TreePage<wire<text>>` | Outside Tree; rejection is already required in canonical source. |
| TreePage applied to a union description over text/bool | Outside Tree despite a plausible native payload union. |
| `Ignore<wire<text>>` | Invalid argument even though the output text itself belongs to Tree. |
| `Page<cell.Cell<Person>>` | Uses Cell's selected library mapping; no Data obligation is invented. |
| `DataPage<cell.Cell<Person>>` with a Wire-containing Cell | Rejected regardless of the native Cell representation. |
| Rename/import Tree differently | Same semantics, without a named-kind branch. |
| Override every list mapping with a generic Collection | Open Page and every application must use that same representation family. |
| Override only native Person to an external type | Source membership may be unchanged; selected native association changes. |

### Candidate guarantee: three different relations

The proposed document currently says an admitted native application must retain:

```text
S : K       established source membership, under its bound assumptions
S -> N      the selected native representation association
```

It must not infer either solely from native assignability. This formulation may
still be too vague: what exactly ranges over S and N for a native caller? Is
soundness enough, or should every supported source application also be expressible
in the native API? We want a precise contract with any necessary trust assumptions.

A third, independent proposition is that a particular native value denotes a
valid Horos value. That is outside this static type-only stage. For example, a
Go string type does not promise valid UTF-8 bytes, and a native record type need
not enforce exact closed-record value semantics. Neither a witness nor a
membership token should accidentally be advertised as a value validator.

Each profile needs to specify independent dimensions together:

| Dimension | Choices to distinguish |
| --- | --- |
| Rejection timing | Type formation, normalization, or later value use; an uninhabitable result is not automatically a formation error. |
| Argument domain | Generated named applications, generated witness families, or independently authored native payload types. |
| Extension model | One closed generation unit, separately generated trusted modules, or unrestricted native customization. |
| Native API promise | Alias interchangeability, nominal wrapper identity, retained source association, or admission followed by erasure. |

These dimensions are not synonyms or one ordered scale. Please recommend which
combinations are useful and realistic for the initial profiles.

### Options and constraints

Considered approaches are deliberately still open:

1. Direct payload bounds on a precisely supported domain.
2. Closed generated witness/evidence domains, with named concrete applications.
3. Open compositional witness families with a documented extension trust model.
4. Source-only admission and representation generation, explicitly weaker than
   native enforcement; this is what the current representation profile does.
5. Another design that gives users a simpler or stronger useful guarantee.

Fixed constraints: no name-sensitive Data/Tree/Function/Stream/Cell cases;
Horos owns kinds and semantic judgments; the generic kernel and native syntax
backend do not learn them; source and target choices remain explicit; no runtime
registry, reflection-based membership, validator or operational behavior; no
silent erasure or fallback when an enforcing profile is selected; no hidden
monomorphization of an open declaration. Unsupported representations and unknown
proof obligations remain explicit failures with source locations.

Flexible choices include supported fragment, target language baseline, native
generic arity, wrappers versus aliases where the profile permits them, bounded
closed-world output versus independent extension, and the shape of compile-time
support artifacts. Please do not treat the sketches as requirements.

### Existing code and available semantic information

All snippets in this subsection are from the inspected bittype baseline. Lines
are source anchors; omitted code is identified. This is enough to distinguish
the actual current information from future witnesses.

`semantics/go/model.go`, lines 19-72 and 103-119 (selected declarations):

```go
type EntityID uint32
type TermID uint32
type Key struct{ Owner, Name string }

type Parameter struct {
    Entity EntityID
    Sort   Key // coarse sort, not a structured kind
}

type TermKind uint8
const (
    Reference TermKind = iota + 1
    Application
    ProviderExpression
)
type Term struct {
    Kind      TermKind
    Target    EntityID
    Callee    TermID
    Arguments []TermID
    Operator  Key
    Encoding  Key
    Data      []byte
    Operands  []Operand
}
type Operand struct {
    Locator []uint32
    Term    TermID
}
type Subject struct {
    Entity EntityID
    Term   TermID
}
type Fact struct {
    Subject  Subject
    Relation Key
    Term     TermID
}
```

IDs are local graph addresses, not durable identities. A Subject selects one
entity or term. A ProviderExpression contains source-owned encoded data plus
explicit operand links. `Semantics[L]` owns an immutable graph for language L;
the kernel only admits graph structure and receipt association. Horos exports
separate canonical module/declaration/parameter identities, result and parameter
contracts, and memberships. Native stages do not interpret those encodings.

`internal/horosmodel/facts.go`, lines 50-70 (current export logic, abbreviated
only to focus on the retained conclusion and assumptions):

```go
// Each membership is one established source guarantee with its original rigid
// assumptions. Query results never mutate this graph or imply a principal kind.
func (e *exporter[L]) memberships(
    subject semantics.Subject, memberships []decl.Membership, depth int,
) error {
    for _, membership := range memberships {
        kind, err := e.kind(membership.Kind(), ast.NodeRef{}, 0, false, depth+1)
        if err != nil { return err }
        operands := []semantics.Operand{{Locator: []uint32{0}, Term: kind}}
        for i, assumption := range membership.Assumptions() {
            reference, err := e.allocateTerm(depth+2)
            if err != nil { return err }
            if err := e.fillTerm(reference, semantics.Term{
                Kind: semantics.Reference,
                Target: e.parameters[assumption.Parameter()],
            }); err != nil { return err }
            contract, err := e.kind(assumption.Kind(), ast.NodeRef{}, 0, false, depth+2)
            if err != nil { return err }
            value, err := e.provider("membership/assumption", core.NewAtom(nil),
                []semantics.Operand{
                    {Locator: []uint32{0}, Term: reference},
                    {Locator: []uint32{1}, Term: contract},
                }, depth+1)
            if err != nil { return err }
            operands = append(operands, semantics.Operand{
                Locator: []uint32{1, uint32(i)}, Term: value,
            })
        }
        value, err := e.provider("membership/proven", core.NewAtom(nil), operands, depth)
        if err != nil { return err }
        if err := e.fact(subject, e.frontend.MembershipRelation(), value); err != nil {
            return err
        }
    }
    return nil
}
```

Here `decl.Membership` and its bound kind/parameter handles come from Horos's
checked Program. `e.kind` projects a kind graph, `e.provider` encodes a
provider-owned term, `e.fact` adds the associated fact, and allocation helpers
enforce graph resource bounds. `core.NewAtom(nil)` is an inert empty payload;
the linked operands carry this expression's relevant semantics. This does not
create a native witness or independently establish the membership again.

`internal/horosgen/generate.go`, lines 30-38, 110-123 and 174-200:

```go
const RepresentationProfile = "pagination-representation/1"
const EnforcementProfile = "pagination-enforcement/1"

// Current Dialect contains other native construction methods as well.
type Dialect[L any] interface {
    Parameter(*Arena[L], string, ast.NodeRef) ast.NodeID
    Reference(*Arena[L], string, ast.NodeID, []ast.NodeID, ast.NodeRef) ast.NodeID
    Text(*Arena[L], ast.NodeRef) ast.NodeID
    Boolean(*Arena[L], ast.NodeRef) ast.NodeID
    List(*Arena[L], ast.NodeID, ast.NodeRef) ast.NodeID
    Record(*Arena[L], ast.NodeID, string, []ast.NodeID, []Field)
    Alias(*Arena[L], ast.NodeID, string, []ast.NodeID, ast.NodeID)
    // ReservedNames, Finish, FieldName omitted.
}
```

Arena is target AST allocation storage; NodeID/NodeRef identify syntax nodes;
Field contains a mapped member's name, optionality, type node and origin. At
the start of Generate, the enforcement branch currently returns:

```go
case EnforcementProfile:
    site := request.Entry()
    return fail(&text.Diagnostic{
        Code: text.Unsupported,
        Message: "native kind enforcement requires an explicit target representation",
        Node: &site,
    })
```

The representation branch requires the analysis receipt and semantic model.
Its parameter method receives only the parameter name and origin. Thus a
source-bound parameter currently becomes a bare TypeScript generic / Go `any`
parameter / ordinary Rust or C++ generic. This profile expressly permits
constraint erasure; we must not describe it as enforcing kinds. No enforcement
code is hidden elsewhere in this interface.

The actual call in `internal/horosgen/generate.go:420` and TypeScript emitter in
`transforms/horos-to-typescript/go/transform.go:47-49` make that erasure concrete:

```go
g.targets[p.Entity] = g.dialect.Parameter(&g.arena, name, pe.Declarations[0])

func (dialect) Parameter(a *g.Arena[s.Language], name string, o ast.NodeRef) ast.NodeID {
    return a.Add(s.KindTypeParameter, g.Attr("name", name), o)
}
```

`p.Entity` is the real source binder; `pe` is its semantic entity. The emitter's
`g` import is the private generator package and `s` is TypeScript syntax. The
target parameter node has a name and origin, with no constraint child.

### Proposed plugin composition boundary

One agnostic coordinator composes a source adapter, selected type-owned mappings
and a target syntax backend. Builtin and library mappings use the same supplied
mechanism. Exact declaration identity selects a library mapping; a source-owned
builtin identity selects a builtin mapping. Selection precedes body expansion.
Consumer overrides name the exact binding they replace; no priority search,
most-specific-kind lookup or error-triggered fallback is allowed.

The current proposed operations are:

```go
// F and T are source and target language marker types.
type Plugin[F, T any] interface {
    ID() PluginID
    Use(context.Context, UseRequest[F, T]) (Type[T], error)
}
type Publisher[F, T any] interface {
    Publish(context.Context, PublishRequest[F, T]) (Symbol[T], error)
}
type Definer[F, T any] interface {
    Define(context.Context, DefineRequest[F, T]) (Decl[T], error)
}
type BinderLink[F any] struct {
    Source      Ref[F]
    TargetIndex int
}
```

Use produces a proper target type expression. Publish exposes an explicitly
requested open declaration. Define fills a previously reserved declaration or
helper; it runs with a fresh formal environment. Request values carry the source
reference/declaration or definition specification and the invocation context.
Type, Symbol and Decl are owned target handles, not raw source strings.
`Context.Lower` recursively uses the selected mapping of a child. Ref includes
the source's effective lexical/substitution/assumption environment. Parameter
binding uses real binder identity, not spelling. The source adapter exposes the
explicitly supplied semantic model and receipts; plugins never run analysis or
resolution implicitly.

BinderLink currently pairs one real source binder to one target slot. It does
not solve witness/native/evidence tuple expansion or a `Native<S>` projection.
The companion consultation owns that protocol. This consultation should state
the relations its design must preserve without turning the generic kernel into
a Horos kind evaluator or a catalogue of target-language exceptions.

Generic coherence already requires that open publication and concrete use agree
under represented substitutions. Mapping every `list<T>` to `Collection<T>` can
compose through one open Page. Mapping only `list<Byte>` to ByteBuffer while
the already-published Page keeps Collection needs a separate explicit
specialization design; it is not a permitted hidden optimization. Bounds,
witnesses and evidence arguments must participate in the same coherence rule.


### Candidate target mechanisms and their open problems

Direct native bounds are allowed only where a defined representation-domain rule
preserves the needed source distinction. Two source descriptions mapped to one
native type but with different memberships defeat a general payload-only bound.
The alternative proposal carries generated source witnesses and membership
evidence. It should not require handwritten lists or grants for new source
types: ordinary source analysis and generation supply any needed declarations.

TypeScript sketch:

```typescript
export type TreePage<S extends TreeMember> = Page<Native<S>>;
export type People = TreePage<PersonType>;
```

TreeMember/Native/PersonType are proposed support types, not implemented APIs.
We have not established an enforcement construction. Structural brands/private
symbols alone may not keep type-only evidence invariant or unmanufacturable.
A union of source witnesses is not one source union descriptor. Unintended
conditional distribution can turn `Page<string> | Page<boolean>` into
`Page<string | boolean>`, admitting mixed collections. Mapping an invalid type
to `never` does not necessarily reject forming the application. The domain must
account for `any`, `never`, intersections, widening, recursive type limits and
the normal assertion/unsafe escape hatches without claiming a security boundary.

Two concrete challenges for the expert: `never` ordinarily satisfies an upper
bound, while `Page<never>` can still contain an empty items array; rejecting only
inhabited item values does not reject that type application. Also, if evidence's
representation slot is covariant, evidence associated with Person may widen to
evidence for `unknown`, potentially admitting unrelated page contents. An
invariance marker must work under the actual target comparison rules, not just
look invariant in a generated declaration.

An illustrative variance counterexample, without assertions or `any`:

```typescript
declare const proofKey: unique symbol;
declare const personKey: unique symbol;
type PersonType = typeof personKey;
type Person = { name: string };
type Evidence<S, N> = { readonly [proofKey]: S; readonly native: N };
type PersonInTree = Evidence<PersonType, Person>;
type UnsafeTreePage<S, N, P extends Evidence<S, N>> = { items: N[] };
type Widened = UnsafeTreePage<PersonType, unknown, PersonInTree>;
const unrelated: Widened = { items: [123] };
```

This is a challenge to a candidate encoding, not a proposed generated body.
Function-property versus method variance and compiler strictness options matter.
Likewise, a private class constructor need not prevent naming an unrestricted
generic evidence class in a type argument: this is a type-only admission problem.

Go sketch (generic alias support is a target-version condition):

```go
type TreePage[S, T any, P TreeEvidence[S, T]] = Page[T]
type People = TreePage[PersonType, Person, PersonInTree]

type evidence[S, T, K any] struct{}
type TreeEvidence[S, T any] interface { evidence[S, T, treeKind] }
type PersonInTree = evidence[PersonType, Person, treeKind]
```

The evidence constructor is package-private. The public alias is emitted only
for an established tuple; the constraint is the exact named instantiation, not
an empty interface or a `~struct{}` approximation. It is intended to tie S, T and
K even when unused in the value representation. Generic evidence families must
retain their prerequisites. This is currently a single-package candidate.
Different packages' private constructors are not interchangeable. Making one
unrestricted public constructor available seems to defeat the relationship. Go
also does not provide an arbitrary recursive predicate on all fields of an
unknown struct, and direct type sets distinguish `string` from `~string`.
Source and kind tokens must have genuinely distinct native identity: aliases
of the same `struct{}` do not create different tokens merely by being named.

```go
type TextType = struct{}
type WireTextType = struct{} // Identical to TextType: not a new identity.
type TextInTree = evidence[TextType, string, treeKind]
type Wrong = TreePage[WireTextType, string, TextInTree]
```

Go's constrained generic alias may reject a bad instantiation while its result
is still identical to `Page[T]`; the source/evidence parameters need not remain
observable in that resulting type. Admission at an alias boundary and recovering
source evidence from an arbitrary payload later are different guarantees.

Rust sketch:

```rust
pub struct TreePage<S: TreeMember> {
    pub page: Page<S::Repr>,
}
pub type People = TreePage<PersonType>;
```

TreeMember would associate Repr and encode membership in a controlled generated
domain. Appropriate sealing and ownership would need to prevent merely asserting
unestablished facts. Ordinary Rust alias bounds are not generally enforced at
use sites; a wrapper changes nominality, access syntax and conversion needs.
We must not silently exchange source-alias interchangeability for a wrapper.
Trait coherence/orphan rules, unions of kind alternatives, separately generated
crates and recursive/coinductive kinds all need an explicit account. Treating
native trait recursion as the Horos greatest fixed point would be unjustified.
The ordinary-alias warning does not settle every alias-based alternative:
a fully qualified associated-type projection might impose obligations at a
different stage. Please assess its promise and rejection timing separately
before concluding that a nominal wrapper is always required.

C++ sketch, targeting C++20 unless another standard is explicitly selected:

```cpp
template<class S>
    requires tree_member<S>
using TreePage = Page<typename S::representation>;
using People = TreePage<PersonType>;
```

Here generated traits/concepts express established membership over a stated
descriptor domain. Unknown source witnesses have no positive membership by
default. Native aliases erase distinctions; separate descriptor tokens may be
needed. Open user specializations cannot automatically count as Horos-established
facts. Partial-specialization coherence, independent compilation and compiler
limits also matter. A concept over native layout or trivial copyability cannot
stand in for source Data membership.
Concepts themselves cannot be specialized, but a supporting class or variable
template can be. Private placement alone is not necessarily a sealing mechanism:
explicit-specialization rules have their own access-checking provisions. Please
state which admitted consumer operations could forge a positive relationship.

For all languages, a finite list of generated witnesses must be described as
finite; it is not a solution to arbitrary open recursive membership. A sealed
generated closure might be useful initially, but consumers need to combine
independently generated libraries. How much can this be extended without making
native consumers implement a second Horos checker or permitting handwritten
semantic grants? Trusted generator output, safe native source, and hostile code
with assertions or explicit customization are different trust domains.
Already-proved recursive source membership might be exported as one concrete
certificate without encoding a greatest-fixed-point solver in the target. This
is distinct from admitting previously unknown recursive source families through
open native composition. Please assess these cases separately.

### Overrides, aliases and package composition

Suppose package A publishes `Page<T>` and Person plus their association/evidence.
Package B imports A's canonical declarations and chooses `Person -> external.User`
instead of A's Person representation. The source fact Person : Tree may still
hold, but evidence tied to the old native type is not enough. Similarly, selecting
an external Page implementation may hide all structural fields while its source
admission and generic obligations still matter.

Which support artifact is reusable: a source membership theorem, a source
witness, a representation projection, a complete native generic declaration, or
none? Which are configuration-relative? Equality of canonical names, plugin IDs
or native spelling alone is not proof of output compatibility. A new structural
mapping must not silently rewrite an imported public generic header. We want a
minimal composition rule and explicit failure boundary, not a new global
registry, universal digest or provenance subsystem.
Within one selected configuration, membership in two different kinds must still
use the same source-to-native association. A design where each membership trait
independently selects a different Repr could break that relationship. Conversely,
two consumer configurations may legitimately choose different representations
for the same source; one globally fixed PersonType::Repr cannot simply be
overwritten by a later consumer.

Generic families need to preserve ignored-argument constraints, conditional
assumptions and result contracts even when the payload alias reduces to a
simpler type. Native higher-order encodings may be deferred, but their absence
must not invalidate promises made for ordinary proper generic parameters.


## Questions for the expert

1. **What should the preservation contract promise?** Give a precise distinction
   between source membership, selected representation and native value validity;
   choose useful admission/enforcement levels and trust domains. Explain whether
   soundness, completeness over a represented domain, formation rejection and
   separate compilation can all be promised together, or where they conflict.
2. **Which encodings are sound and useful in each target?** For TypeScript, Go,
   Rust and C++, assess direct bounds and the witness/evidence sketches, identify
   concrete counterexamples, and propose a minimal workable initial domain or
   explicit impossibility boundary. Include native unions/bottom/widening,
   ignored generic arguments, aliases versus wrappers, and recursive kinds.
3. **How should independently generated modules and overrides compose?** Explain
   ownership and identity of source witnesses, source membership, native
   associations and exported evidence; address sealing versus extension, Rust
   coherence, Go private evidence, and C++ customization without manual semantic
   grants. State what a representation override preserves or invalidates.
4. **What semantic obligations must the agnostic generation protocol carry?**
   Use Page, DataPage, TreePage, Ignore and Page<Cell<Person>> to describe the
   required relationship between open definitions and concrete applications,
   conditional/result evidence and target argument expansion. State semantic
   laws and counterexamples; leave concrete Go contracts, header elaboration and
   argument-synthesis mechanics to the companion consultation.
5. **What are we missing, and what would you design differently?** Challenge the
   witness approach and even the usefulness of a universal native-enforcement
   goal. Recommend a staged conceptual plan, decisive counterexamples and honest
   unsupported boundaries before implementation.

## References and scope of evidence

Repository sources: `semantics/go/model.go`, `internal/horosmodel/facts.go`,
`internal/horosgen/generate.go`, `docs/kinds.md`, `docs/model.md`,
`docs/type-generation-plugins.md`, `docs/native-kind-constraints.md`; research
0005-0007. Source inspection only; no target examples were compiled for this brief.

Primary language references for the limitations to examine:

- [TypeScript generic constraints](https://www.typescriptlang.org/docs/handbook/2/generics.html)
- [TypeScript conditional type distribution](https://www.typescriptlang.org/docs/handbook/2/conditional-types.html)
- [TypeScript type compatibility](https://www.typescriptlang.org/docs/handbook/type-compatibility.html)
- [TypeScript strict function variance](https://www.typescriptlang.org/tsconfig/strictFunctionTypes.html)
- [Go general interfaces](https://go.dev/ref/spec#General_interfaces)
- [Go type identity](https://go.dev/ref/spec#Type_identity)
- [Go generic instantiation](https://go.dev/ref/spec#Instantiations)
- [Rust alias-bound limitation](https://doc.rust-lang.org/rustc/lints/listing/warn-by-default.html#type-alias-bounds)
- [Rust associated types](https://doc.rust-lang.org/reference/items/associated-items.html#associated-types)
- [Rust implementation coherence](https://doc.rust-lang.org/reference/items/implementations.html#trait-implementation-coherence)
- [C++ template constraints](https://eel.is/c++draft/temp.constr)
- [C++ specialization access rules](https://eel.is/c++draft/temp.spec.general#6)
