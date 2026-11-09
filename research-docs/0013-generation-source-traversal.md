# Research: Kernel-owned semantic traversal and narrow source interpretation in type generation

**ID:** 0013
**Date:** 8 November 2026
**Status:** applied
**Delivered:** https://github.com/Bitspark/bittype/pull/191, corrected in https://github.com/Bitspark/bittype/pull/192
**Run-ID:** run_f7cfb3b5-7cf2-40d8-8884-ce150a7adcf7
**Reviewed:** https://github.com/Bitspark/bittype/issues/158#issuecomment-6720989325
**Owner:** claude-bittype-gaps
**Issue:** https://github.com/Bitspark/bittype/issues/158
**Code baseline:** bittype 57187bcb0621c8da3ac1ef12e12887d4d079cd94; horos Go core 0.10.2

## Question

A type generator turns declarations written in a source language into type
declarations in native programming languages. Its generation kernel is meant to be
language-agnostic. Today the kernel asks each *source integration* (the plug-in for one
source language) for two different things through one interface:

- **Navigation** of an already established semantic graph: which declaration a use
  names, a declaration's parameters and body, which argument of an application belongs
  to which parameter, and whether one term is a child of another.
- **Interpretation** of source-specific meaning: stable declaration identity, the
  meaning of language-specific expression forms and their operands, and whether a use
  is admissible under the declaration's assumptions.

The navigation facts are already recorded in a shared, language-agnostic semantic graph.
A source integration that consumes that graph has to re-walk it to answer the kernel's
navigation questions. The project's architecture rules treat this as a defect: a shared
mechanism should own the invariant mechanics, and a language owner should contribute
only its meanings.

**The decision we need help with:** how should the boundary between the generation
kernel and a source integration be drawn, so that the kernel follows declaration,
reference and application structure through the shared graph, while source-owned
contracts keep interpretation and applicability judgments? The design must keep the
guarantees listed below, and it must say what mapping libraries, which today read the
source integration's records directly, see instead.

**How to read the code below.** The project rule for interface design (called AC-09)
gives the existing interface *zero weight*, including as a tiebreaker, and gives
migration effort zero weight as well. A retained element needs the same independent
justification as a new one. The code excerpts are evidence of who does what today. They
are not an interface to preserve.

The owner's intended division, for the running example `People = Page<Cell<Person>>`:

| Decision | Intended owner |
| --- | --- |
| Follow established references to Page, Cell and Person | Kernel, over the captured graph |
| Follow `Page.T` to `Cell<Person>` and `Cell.V` to Person | Kernel, over established application bindings |
| Interpret horos constructors, kind judgments and pinned identity | horos integration |
| Select a native representation of Cell | Cell's mapping library |
| Admit source evidence and preserve capture and scope association | Shared contract guarantees, with source-owned judgments |

## Context

### System overview

**bittype** is a Go library. It parses, checks and prints source code, and it generates
native type declarations for TypeScript, Go, Rust and C++. Its input language is
**horos**, a small declaration language for data types. A horos module declares named
type definitions, which may be generic. The running example spans two modules:

```horos
module cells;
type Cell<V> = wire<V>;      // wire: an opaque transport descriptor with no default native form

module consumer;
import cells from cells pin "sha256:<module pin>";
type Person = { name: text; };
type Page<T> = { items: list<T>; next?: text; };
type People = Page<cells.Cell<Person>>;
```

`text`, `bool`, `list<T>`, `wire<T>` and record braces are built-in *descriptors*.
Expected TypeScript output, with a mapping library that owns Cell and represents it by
an imported native `Cell`:

```typescript
import type { Cell } from "@bitspark/cell";
export interface Person { name: string; }
export interface Page<T> { items: T[]; next?: string; }
export type People = Page<Cell<Person>>;
```

Cell's horos body (`wire<V>`) is never generated. Its library supplies the
representation, but `Person` must still be admitted and generated.

**Terminology warning.** horos uses "kind" and "sort" the other way round from much of
the type-theory literature. A horos **kind** is a classification of types, close to a
type class or a bound: a parameter written `T: h.Text` must be a type that `h.Text`
admits, and `h` is the standard *prelude* module that defines kinds such as `h.Text`
and `h.Data`. A horos **sort** is what Haskell calls a kind: `proper` (an ordinary type,
`*`) or an `arrow` (a type constructor such as `* -> *`). Kinds appear only as bounds,
never as field types.

Independently owned parts cooperate:

| Part | Owns |
| --- | --- |
| horos **analysis provider** | Name resolution, sort and kind checking, and horos's semantic facts. It publishes them as one immutable semantic graph per captured input. |
| **Generation kernel** (Go package `generation`, `contracts/generation/go`) | Planning: selection of mappings, reservation of native declarations, substitution of native type arguments, scope checks, limits, and all-or-nothing output. It is meant to interpret no source or native rules. |
| **Source integration** (for horos: package `source`, `transform/syntax/horos/source/go`) | Today: navigation of the graph *and* horos interpretation. This is the subject of the document. |
| **Mapping libraries** | Declarative descriptions of native representations, one per source type family and target. Examples: horos lists to TypeScript arrays, horos records to TypeScript interfaces, and a library that owns `Cell<V>` and maps it to an imported native `Cell<T>`. |
| **Target** (one per native language) | Checks native expressions and declarations, and builds native syntax trees from a closed native program. |

**Scale and consumers.** The repository has five horos mapping libraries per target
(text, bool, list, record, declarations), four targets, and test-only libraries such as
the Cell mapping. Mapping libraries are meant to be written by the owners of type
families, so other people's code will implement this interface. Only horos publishes a
graph today. No second production source language is planned. An independent test
source (described below) stands in for one, so the cost of re-implementing navigation
is prospective, but it is the extension the architecture promises.

### Vocabulary

- **Generic parameters.** `F` is the source-language marker type, `T` the
  target-language marker type, and `S` the source integration's concrete session type.
  `L` is a language marker in general. Markers are compile-time types only.
- **Package names.** `generation` (alias `g`) is the kernel contract. `semantics` holds
  the graph. `semantic` (singular) holds snapshots, evidence and reports. `source` and
  `rules` are horos's integration and generic horos mapping rules. `model` is a target's
  native model. `representation`, `transform` and `ast` are shared plumbing for typed
  language envelopes, typed requests and syntax-node references.
- **Snapshot.** An immutable capture of input syntax trees. Capturing equal bytes twice
  yields two different snapshots: identity is by allocation, not content.
- **Result** or **evidence.** The immutable output of an analysis operation, bound to
  exactly one snapshot. It carries *reports* (named obligations, each with a completion
  status and a verdict), diagnostics and, optionally, a semantic graph.
- **Declaration, binder, use, environment.** A declaration is a named type definition.
  A binder is one of its real type parameters. A use is one type occurrence in a body.
  An environment is the set of binders and assumptions under which a use is
  interpreted. Today that is the enclosing declaration.
- **"Subject" has three senses.**
  - The graph's `semantics.Subject` names an entity or a term that a fact is about.
  - `UseData.Subject` is the source-certified *effective subject*. Two different
    occurrences may denote the same subject, and only the source may say so.
  - A *formal subject* is a synthetic use that stands for a binder, because a binder
    has no term of its own.
- **"Profile" has three senses.** A *generation profile* (for example
  `pagination-representation/1`) selects which mapping libraries apply. A *provider
  profile* names an analysis configuration inside an evidence requirement. A *language
  profile* qualifies a language identity, and therefore a selector.
- **Selector.** A stable identity string, owned by the source and qualified by
  language. Mapping libraries register rules against selectors. The kernel compares
  them and never interprets them. horos has three forms:
  - `declaration/<module pin>/<name>` for a named declaration. The *module pin* is the
    SHA-256 hash of the module's canonical encoding.
  - `builtin/<constructor>/<discriminator>`, length-prefixed, for a built-in
    descriptor. The discriminator distinguishes variants, for example `utf8-text` for a
    scalar or `closed` for a record.
  - `bound` for a use of a parameter.
- **Native program.** The closed list of planned native declarations that a target
  turns into syntax trees. Planning resolves every choice before construction (`Run`)
  starts. "A run" means one generation invocation.
- **Error outcomes.**

  | Error | Meaning |
  | --- | --- |
  | `ErrScope` | An address, record or relation outside its admitted scope: a foreign record, a wrong environment, a non-child. |
  | `ErrUnsupported` | Valid input outside supported coverage. Reported as a diagnostic with code Unsupported at a syntax site. |
  | `ErrInvalid` | A malformed request or description. |
  | `ErrLimit` | A work or size limit was exhausted. |
  | `ErrAmbiguous` | More than one mapping applies, or more than one recipe. |
  | `semantic.ErrInvalidRequest` | Malformed or mismatched semantic input. |
  | `semantic.ErrEvidenceUnavailable` | Required evidence is not established. |

- **Tracking identifiers.** GEN-01 (this gap), GEN-02 (missing native enforcement of
  source constraints) and ACX-006 (the architecture case record for this defect) are
  internal tracking identifiers.

### Guarantees that must survive

These come from earlier consultations and from the delivered interface. Each one is a
requirement on the replacement, independent of the current method names.

| Guarantee | Meaning |
| --- | --- |
| Capture identity | Planning consumes exactly the admitted snapshot and evidence. A recapture of equal bytes, or evidence from another snapshot, is rejected. |
| Real binders | A source binder is the declared parameter itself, never a name or position. Native conventions may reorder, omit or add native parameters, so binder *i* is not native parameter *i*. |
| Environments and assumptions | A use is interpreted under its own declaration's assumptions (for example `D: h.Data`). Native results built while planning a callee's body, which refer to the callee's own parameters, must not leak into a caller. An assumption must not become an unconditional fact. |
| Foreign-reference rejection | Records, handles or terms from another run, source or graph cannot be passed off as belonging to this run. Possessing an address is not admission. |
| Whole-use admission | A use is admitted as a whole, including every actual, before any mapping represents it. An actual that the native convention omits, such as a phantom parameter, still needs admission. Otherwise a rule that emits constant text could bypass a source bound by never traversing its argument. |
| No hidden analysis | Generation never runs a resolver or an analyzer. Missing evidence fails; it is never computed on demand. |
| Unknown is not false | An uncovered or absent fact is distinct from a failed judgment. |
| Bounded lifetimes | A mapping callback's access to the source and the kernel ends when the callback returns. Anything retained and used later fails. |
| Limits and whole-output failure | Work, depth, references, formals, reservations and output size are bounded. Any failure, cancellation or exhausted limit yields no usable output, and never triggers a weaker fallback mapping. |
| Selector ownership | Stable selector identity belongs to the source. The kernel compares selector identity, never spellings or graph-local IDs. |

### Why the current boundary is a defect, and what "done" means

Four project architecture rules apply, paraphrased:

- **AC-02.** Meanings stay with their semantic owners. A shared mechanism owns its
  invariant guarantees, and language owners contribute interpretations.
- **AC-05.** Composition preserves guarantees. Openness may not be bought with arbitrary
  payloads, unchecked claims, or callbacks that discard guarantees.
- **AC-06.** Limited coverage does not redefine the architecture. Working examples
  through today's interface do not make the wrong split acceptable.
- **AC-09.** As above: interfaces are derived from the domain, owners and guarantees.

Completion criteria set by the owner:

1. The kernel consumes declaration, reference and application structure through the
   shared semantic context. No source integration supplies a parallel traversal of the
   same facts.
2. Narrow source-owned contracts keep constructor and operand interpretation, stable
   selector identity and applicability judgments. Unknown or uncovered facts stay
   distinct from failed judgments, and there is no hidden analysis.
3. Every guarantee above stays enforced.
4. horos and an independent source exercise the same kernel traversal, including nested
   `Page<Cell<Person>>` and invalid scope and evidence cases.
5. The superseded reader interfaces and duplicated traversal are removed. No
   compatibility aliases remain.

### Constraints: what is fixed and what may change

| Area | Status |
| --- | --- |
| Generation contract (source interface, catalogue, mapping-library interface, expansion handles) | **May change freely.** This is what is being redesigned. |
| The horos source integration, its reading API, and the horos and native mapping libraries | May change, as consumers of the new contract |
| The independent test source | May change. It must publish a graph and supply only its own interpretation. |
| The graph data model (`semantics`) | Its meaning is fixed. Language-agnostic structural queries or indexes may be added if they interpret no provider vocabulary. |
| The horos analysis provider and its facts | Fixed vocabulary, owned by horos. Generation must not trigger more analysis. |
| Snapshots, evidence admission, reports | Fixed |
| Native program, targets, construction (`Run`), file placement | Fixed. This change should not touch them. |

**What the kernel may read and must not read** (our reading of the rules, for critique):

- **May read:** `Entity.Owner`, `Entity.Parameters` (as binder identities),
  `Entity.Body`, `Entity.Declarations` (syntax sites); `Term.Kind`, `Term.Target`,
  `Term.Callee`, `Term.Arguments`, and `Operand.Term` (as a child relation);
  `DirectParameterBindings`.
- **Must not read:** `Entity.Kind`, `Entity.Name`, `Parameter.Sort`, `Term.Operator`,
  `Term.Encoding`, `Term.Data`, the meaning of an operand `Locator`, and any fact
  relation or fact term.

**Out of scope:** native enforcement of source bounds in generated code (GEN-02), and
recursive types, which are rejected today. The design should not make either harder.

### The shared semantic graph

The graph is the shared context. It stores entities (declarations, parameters, modules
and so on), terms (the expression trees of bodies), occurrence associations (which term
or entity a syntax position denotes) and provider-owned facts. Its constructor validates
graph integrity only: addresses, parameter ownership, acyclic term containment. It does
not check language rules. IDs are one-based addresses local to one graph. Excerpt from
`semantics/go/model.go`; comments beginning with `//>` are annotations added for this
document:

```go
type EntityID uint32
type TermID uint32
type Key struct{ Owner, Name string } // provider-owned vocabulary, e.g. {horos, type-definition}

type Entity struct {
	Kind         Key
	Name         string
	Owner        EntityID      //> enclosing entity; for a parameter, its declaration
	Declarations []ast.NodeRef //> syntax sites that declare it (ast.NodeRef = tree ID + node ID)
	Parameters   []Parameter   //> real binders, in declaration order
	Body         TermID
}
type Parameter struct {
	Entity EntityID
	Sort   Key //> coarse provider-owned sort, e.g. {horos, proper} or {horos, constructor}
}

// Term has exactly one variant (Reference, Application or ProviderExpression).
// Application arguments are effective positional arguments, not a mixture of
// written, inferred, defaulted or packed arguments. Reference owns Target;
// Application owns Callee and Arguments. Applicability, arity and saturation
// belong to the provider, not graph admission. A provider expression's encoding
// defines its bytes and operand locators.
type Term struct {
	Kind      TermKind
	Target    EntityID  //> Reference
	Callee    TermID    //> Application
	Arguments []TermID  //> Application
	Operator  Key       //> ProviderExpression, e.g. {horos, descriptor/list-of}
	Encoding  Key       //> ProviderExpression payload encoding
	Data      []byte    //> ProviderExpression payload
	Operands  []Operand //> ProviderExpression children
}
type Operand struct {
	Locator []uint32 //> provider-defined path inside the payload
	Term    TermID
}
type Subject struct{ Entity EntityID; Term TermID }          //> exactly one is set
type Fact struct{ Subject Subject; Relation Key; Term TermID } //> several values per relation allowed
```

Graph admission enforces, among other things, that each parameter entity is listed by
exactly one entity and names it as its `Owner`, that a reference names an existing
entity, and that term containment (callee, arguments, operands) is acyclic. Entity
references may form cycles. Containment may form a DAG: one term may be shared.

The graph offers point queries: `Entity(id)`, `Term(id)`, `Facts(subject, relation)`,
`Resolve(occurrence)` and `TermAt(occurrence)` for syntax positions,
`DependenciesFrom(entity, relation)`, `Data()` (a full copy), and one structural pairing
query:

```go
// DirectParameterBindings pairs the real parameters of a direct Reference
// callee with an equally sized argument list. Parameter-headed applications,
// other callees and arity mismatches have no such view. This structural query
// does not establish applicability.
// A direct zero-parameter application returns an empty successful view.
func (m Semantics[L]) DirectParameterBindings(id TermID) ([]ArgumentBinding, bool)
type ArgumentBinding struct { Parameter EntityID; Argument TermID }
```

It applies to Application terms only. A plain `Reference` such as `Person` has no
view. There are no reverse or indexed queries, such as "terms in a declaration's body"
or "entity declared at a site". A consumer builds those indexes itself.

### Evidence and capture identity

The graph reaches generation inside evidence. A source integration declares the
evidence it needs:

```go
type EvidenceRequirement struct {
	Operation     semantics.Key // e.g. {bittype, analysis}
	Profile       string        // provider profile
	Configuration []Setting     // exact key/value settings
	Reports       []Requirement // each: obligation key + required verdict
	Semantics     bool          // a graph must be published
}
```

Evidence admission matches requirement *i* to result *i*. It rejects a result from a
different snapshot (pointer identity), any mismatch of operation, profile or
configuration, a required report that is not complete with exactly the required
verdict, and a missing graph when `Semantics` is set. Completion is Complete,
Incomplete or Unsupported. A verdict is Holds, Rejected or Undetermined.
"Complete/Undetermined" is a completed enumeration that asserts no truth value. The
horos source requires `program-validity` Complete/Holds, plus Complete/Undetermined
for its enumeration reports: declarations, occurrence bindings, applications,
reference dependencies and type facts.

A result returns its graph as an untyped envelope.
`representation.UnpackSemantics`, called with the snapshot language and the
envelope, produces the typed graph after checking the language identity, profile and marker. A graph has no pointer
back to its snapshot. The only binding between them is the result that published it.

### What horos records and what its analysis judges

Entity kinds emitted by horos:

| Entity kind | Owner | Parameters | Body |
| --- | --- | --- | --- |
| `{horos, module}` | none | none | none |
| `{horos, type-definition}` | its module | its type-parameter entities | the type expression |
| `{horos, kind-definition}` | its module | its kind-parameter entities | the kind expression |
| `{horos, type-parameter}` | its type-definition | none | none |
| `{horos, kind-parameter}` | its kind-definition | none | none |

Imported modules and the prelude are ordinary module entities in the same graph. An
import creates no entity: a reference to an imported declaration targets that
declaration's entity. horos has no nested generic declarations. Every type-definition
is owned by a module.

Type-definition bodies use the three term variants as follows:

| horos expression | Term |
| --- | --- |
| A parameter, `T` | `Reference` to the type-parameter entity |
| A named type with no arguments, `Person` | `Reference` to the type-definition. It is never a zero-argument application. |
| An application, `Page<Cell<Person>>` | `Application` with a callee term (normally a `Reference`) and argument terms |
| A descriptor: `list<T>`, `text`, `{ ... }`, `wire<V>` | `ProviderExpression`. Operator `{horos, descriptor/<constructor>}`, payload in the encoding `{ontos, codec-v1}` (a canonical binary value encoding), nested type terms as operands. |

Operand locators are paths into the encoded descriptor tuple, defined by horos. A list
element is at `[1]`. A record payload is the tuple `(horos/record, closed|open,
field…)`, where each field is the tuple `(horos/field, name, req|opt, type)`. The type
of zero-based field *k* is therefore at `[k+2, 3]`, and field names exist only in the
payload.

horos publishes several facts. Generation reads three of them:

- `sort`, on every type term, declaration and parameter. Its value is a sort term
  `{horos, sort/proper}` or `{horos, sort/arrow}`.
- `result-sort`, on declarations.
- `canonical-identity`, on modules. Its value is a term `{horos, identity/module}`
  holding the 32-byte module pin.

It also reads the occurrence associations, for diagnostic sites. horos additionally
publishes `declared-kind` (a parameter's bound, as a kind term) and `membership` (a
proven kind membership with its assumptions, on every type term), but generation does
not read them today. A missing fact asserts nothing.

The analysis judges the whole module closure before it reports `program-validity`
Complete/Holds. Under each declaration's own parameter assumptions:

- every application exactly saturates its callee and has sort-correct arguments;
- every argument provably satisfies the callee's declared parameter kind;
- every descriptor operand satisfies its formation contract, horos's well-formedness
  rule for that descriptor's operands;
- the body satisfies its declared result kind.

When a judgment fails, no graph is published. The analysis does **not** require types
to be proper, that is, not type constructors. The horos generation integration adds
that restriction in every generation profile today, because constructor-valued types
have no default native form. horos therefore accepts forms that generation must reject,
and today it rejects some of them with the wrong kind of error:

| horos-valid form | Graph shape | Today's generation outcome |
| --- | --- | --- |
| Constructor-valued declaration, or constructor-sorted parameter | sort `sort/arrow` | Unsupported diagnostic when the declaration is opened (correct) |
| Parameter-headed application `F<A>` in `type Later<F: ...> = F<text>` | Application whose callee references a parameter | Unsupported at open, because `F` has an arrow sort. A later scope-error branch is unreachable but would misreport the case. |
| Curried application `Later<Alias><text>` | Application whose callee is an application | Unsupported diagnostic (correct) |
| Applying a zero-parameter constructor alias: `type Alias: ... = ListType; type X = Alias<text>;` | One-argument application whose callee has no parameters | Invalid-request error (wrong: should be Unsupported) |

The misclassifications show why "no structural view", "rejected judgment", "uncovered
fact" and "malformed input" need distinct outcomes.

### The current generation contract

The public contract, from `contracts/generation/go/source.go`, `mapping.go`,
`reference.go` and `description.go`. Handles are owner-qualified addresses:

```go
// A Space issues addresses for one owner. Space.Ref(n) returns the address
// numbered n (n > 0); Space.Index(r) returns n only if r came from this Space.
// A Ref is an address, not a proof that a record exists.
type Space[K any] struct{ /* private owner identity */ }
type Ref[K any] struct{ /* private owner pointer, slot uint64 */ }
type Owner[K any] struct{ /* identity check only: Owns(r) */ }

type Declaration[F any] = Ref[DeclarationRole[F]]
type Binder[F any]      = Ref[BinderRole[F]]
type Use[F any]         = Ref[UseRole[F]]
type Environment[F any] = Ref[EnvironmentRole[F]]

type Formal[F any] struct{ Binder Binder[F]; Subject Use[F] } //> Subject: the formal subject
type Actual[F any] struct{ Binder Binder[F]; Use Use[F] }

// These are provider construction records. Only the configured source can admit them.
type OpenData[F any] struct {
	Declaration Declaration[F]
	Environment Environment[F]
	Body        Use[F]
	Formals     []Formal[F]
}
type UseData[F any] struct {
	Use         Use[F]          //> the occurrence
	Subject     Use[F]          //> the effective subject
	Environment Environment[F]
	Callee      Declaration[F]  //> set for a use of a declaration
	Bound       Binder[F]       //> set for a use of a parameter in scope
	Actuals     []Actual[F]     //> callee binders paired with argument uses
}

type Source[F any] interface {
	CheckInput(context.Context, semantic.Snapshot, []semantic.Result) error //> same snapshot and evidence as the source was opened with
	Owners() SourceOwners[F]          //> the four Owner values
	RootEnvironment() Environment[F]  //> no kernel or test caller
	// Authenticate checks that each address belongs to this source's spaces and names
	// an existing record. It does not establish applicability.
	Authenticate(context.Context, References[F]) error
	Open(context.Context, Declaration[F]) (OpenData[F], error)
	// Admit establishes the whole use's applicability under its original
	// assumptions using supplied evidence, including actuals the mapping may
	// never represent. Use identifies the occurrence; Subject is the
	// source-certified effective subject. Missing evidence must fail, never
	// trigger hidden analysis.
	Admit(context.Context, Environment[F], Use[F]) (UseData[F], error)
	// CheckChild authenticates a direct semantic operand or generic actual of parent.
	CheckChild(context.Context, Environment[F], Use[F], Use[F]) error
}
type Catalog[F any] interface {
	Declarations() []Declaration[F] //> no kernel caller
	DeclarationSelector(Declaration[F]) (Selector, bool, error) //> bool: has a selector
	UseSelector(UseData[F]) (Selector, error)
}
type SourceSession[F any] interface {
	Source[F]
	Catalog[F]
	ResolveRoots(context.Context, []ast.NodeRef) ([]Declaration[F], error) //> syntax sites to declarations
}
type SourceProvider[F any, S SourceSession[F]] interface {
	Contract() transform.Contract // {Evidence []EvidenceRequirement}
	Open(context.Context, transform.Request[F]) (S, error) // request: snapshot, entry tree, admitted evidence
}
```

The kernel never checks address ownership itself. It checks only that `Owners()` is
valid. Every rejection of foreign or forged addresses happens inside the source's
`Authenticate`.

Mapping libraries are typed by the session type `S`, so they can call the source's own
reading API as well as the contract:

```go
type Slot struct{ /* unique identity */ }              //> a named hole in an output expression
type Capture[F any] struct{ Slot Slot; Use Use[F] }

type Pattern[F any, S SourceSession[F]] struct {
	Selector Selector
	Slots    []Slot
	// Capture interprets the selected source meaning only; the kernel checks each child relation.
	Capture func(context.Context, S, UseData[F]) ([]Capture[F], error)
}
type BindingDescription[F, T any, S SourceSession[F]] struct { //> a rule for uses matching Selector
	ID     string
	Input  Pattern[F, S]
	Output Expression[T]                                                     //> fixed output, or
	Expand func(context.Context, Expansion[F, T, S]) (Expression[T], error) //> computed output
}
type DeclarationRecipe[F, T any, S SourceSession[F]] struct { //> how to emit a declaration
	ID       string
	Selector Selector //> zero: the default recipe for any declaration
	Expand   func(context.Context, DeclarationRequest[F, T, S]) (DeclarationDescription[T], error)
}
type DeclarationDescription[T any] struct {
	Name    string
	Form    Form[T]        //> target-owned declaration form, e.g. interface or alias
	Origin  ast.NodeRef
	Generic GenericConvention[T]
	Body    Expression[T]
}
// GenericConvention: Bindings[i] is what source binder i means inside the native
// declaration; Arguments[j] says how a use fills native parameter j from the source
// actuals. Parameters may be reordered, omitted or added.
type GenericConvention[T any] struct {
	Parameters []ParameterDescription // {Name, Origin, RequiredUse}; RequiredUse: must appear in the body
	Bindings   []Expression[T]
	Arguments  []Expression[T]
}
type MappingLibrary[F, T any, S SourceSession[F]] interface {
	Describe(from, to representation.Language) (MappingDescription[F, T, S], error) // {Profile, Bindings, Declarations, Helpers}
}
```

`Expression[T]` is a native expression tree with a few kernel-owned placeholders:

- `Mapped(slot)` is replaced by the planned representation of the use captured in that
  slot.
- `Parameter(i)` is native type parameter *i* of the declaration being planned.
- `Argument(i)` is the representation of source actual *i* at a use site.
- `Helper(key)` refers to a library-owned helper declaration.

While a callback runs, a library holds an `Expansion`. Its methods fail after the
callback returns:

```go
func (r Expansion[F, T, S]) Source() (S, error)       // the session
func (r Expansion[F, T, S]) Use() (UseData[F], error) // the admitted use being mapped
// Mapped records use as a capture under a fresh slot and returns the placeholder
// Mapped(slot); the planner resolves it after the callback. A declaration recipe may
// capture only its body; a binding rule only a child accepted by Source.CheckChild.
// Each call counts against the References limit.
func (r Expansion[F, T, S]) Mapped(use Use[F]) (Expression[T], error)
func (r DeclarationRequest[F, T, S]) Open() (OpenData[F], error) // DeclarationRequest embeds Expansion
```

### Today's planner on the running example

The planner declares each requested root and resolves uses depth-first. Each step is
tagged **N** (navigation of structure the graph already holds) or **I**
(interpretation). Graph terms: `u0 = App(u1=Ref(Page), [u2])`,
`u2 = App(u3=Ref(Cell), [u4])`, `u4 = Ref(Person)`. Page's body is a record `g0` whose
field types are `g1 = list<g2>` with `g2 = Ref(Page.T)`, and `g3 = text`.

1. `declare(People)`: `Authenticate`, then `DeclarationSelector(People)` (**I**:
   module pin + name). There is no exact recipe, so the default recipe applies.
2. `Open(People)` returns body `u0` and no formals (**N**). It checks that
   People's result sort is proper (**I**).
3. The recipe calls `Declaration(People)` for name, pin and root-module flag (**I**),
   reads `u0`'s constructor to pick the alias form (**I**), and captures `u0` as its
   body.
4. `Admit(People, u0)`:
   - checks that `u0` lies in People's body (**N**: a pre-built term-to-declaration
     map);
   - checks that `u0`'s sort is proper (**I**);
   - finds callee Page through the Reference callee and pairs `Page.T` with `u2` (**N**).
5. `Admit(People, u2)` is the whole-use admission of the actual (**N**+**I**).
6. `DeclarationSelector(Page)` (**I**) finds no exact binding, so the planner calls
   `declare(Page)`. `Open(Page)` gives body `g0` and formal `T` (**N**) and checks that
   both are proper (**I**).
7. The recipe calls `Parameter(T)` for the name (**I**) and records the convention
   `T ↦ Parameter(0)`.
8. `Admit(Page, g0)`, then `UseSelector(g0)` (**I**: `builtin/record/closed`), selects
   the record rule.
   - The record rule calls `IsBody` (**N**) and `Fields(g0)` (**I**: names, optional
     flags, and field-type uses found through payload locators).
   - It captures `g1` and `g3`. Each capture passes `CheckChild` (**N**: is it an operand
     of `g0`?).
9. `use(Page, g1)`: `UseSelector` selects the list rule. Its capture reads the
   operand at `[1]`, which is `g2` (**I**), and `CheckChild(g1, g2)` passes (**N**).
10. `use(Page, g2)`: `Admit` reports that `g2` references `Page.T`, owned by Page, so the
    use is bound (**N**). The planner substitutes `Parameter(0)`.
11. Back at `u0`, the convention needs argument 0, so the planner calls
    `use(People, u2)`. `Admit` finds callee Cell and pairs `Cell.V` with `u4` (**N**).
    `DeclarationSelector(Cell)` (**I**) matches the Cell library's exact binding.
12. The Cell library's capture reads `use.Actuals[0].Use`, which is `u4`. That is the
    pairing the source computed in step 11 (**N**). `CheckChild(u2, u4)` passes (**N**).
    Cell's body is never opened.
13. `use(People, u4)`: `Admit` reports a reference to Person with no actuals (**N**).
    There is no exact binding, so the planner calls `declare(Person)`, which works like
    steps 1 to 8.
14. The result is `People = Page<Cell<Person>>`.

Every **N** step is answered by the source from structure already in the graph.

### The planner's use resolution (condensed)

```go
func (p *planner[F, T, S]) declare(d Declaration[F], depth int) (uint32, error) {
	if id := p.declarations[d]; id != 0 {
		if p.active[id] { return 0, fmt.Errorf("%w: recursive native declaration", ErrUnsupported) }
		return id, nil // declarations are memoized; uses are not
	}
	// selector, then exactly one recipe (exact, else the default; several: ErrAmbiguous;
	// an exact use binding without a recipe: ErrUnsupported)
	open, err := p.source.Open(p.ctx, d)
	id, _ := p.reserve()
	p.declarations[d], p.opens[d] = id, open // recorded before Expand, so recursion is detected
	desc, err := recipe.Expand(ctx, request)
	if len(desc.Generic.Bindings) != len(open.Formals) { return 0, ErrScope }
	env := map[Binder[F]]Expression[T]{} // real binder -> its native meaning
	for i, f := range open.Formals { env[f.Binder] = /* desc.Generic.Bindings[i] */ }
	desc.Body, err = p.resolve(desc.Body, captures, open.Environment, env, id, depth+1)
	// RequiredUse parameters must occur in the body; the target checks the declaration
}

// resolve replaces each Mapped(slot) placeholder by p.use(capturedUse, ...) in the
// caller's environment. Operands no mapping captures are never visited.

func (p *planner[F, T, S]) use(use Use[F], environment Environment[F], env map[Binder[F]]Expression[T], scope uint32, depth int) (Expression[T], error) {
	data, err := p.source.Admit(p.ctx, environment, use)
	if data.Use != use || data.Environment != environment { return Expression[T]{}, ErrScope }
	if data.Bound.Valid() { return env[data.Bound], nil } // missing entry: ErrScope
	// Whole-source admission applies even to arguments a representation omits.
	// (One Admit per actual; the actual's own subtree is not admitted here.)
	for _, a := range data.Actuals { p.source.Admit(p.ctx, environment, a.Use) }
	selector := /* DeclarationSelector(callee) if a callee, else UseSelector(data) */
	switch choices := p.bindings[selector.Key()]; len(choices) {
	case 0: // fall through to the callee's declaration
	case 1: // exact binding: run Capture, CheckChild each capture, Expand, resolve
		return ...
	default:
		return Expression[T]{}, ErrAmbiguous // never a fallback
	}
	if !data.Callee.Valid() { return Expression[T]{}, ErrUnsupported }
	id, err := p.declare(data.Callee, depth+1)
	// pair actuals with the callee's real binders; map only the actuals the native
	// convention references; build native arguments from convention.Arguments
	return Expression[T]{kind: DeclarationExpression, declaration: id, args: args}, nil
}
```

Construction later receives only the closed native program. It has no source access
and no mapping callbacks.

### The horos source integration

The horos source is the concrete `S`. Its addresses encode graph positions:

| Address | Slot number |
| --- | --- |
| Declaration | the entity ID |
| Binder | the parameter's entity ID |
| Use | a term ID; for a formal subject, the term count plus the parameter's entity ID |
| Environment | 1 for the root; for a declaration, its entity ID plus 1 |

At open time, the horos source copies the graph. It indexes the three facts it reads
(rejecting a second value on one subject, which the graph itself would allow) and
pre-walks every type-definition body:

```go
for i, e := range facts.Entities {
	id := semantics.EntityID(i + 1)
	if e.Kind != (semantics.Key{Owner: "horos", Name: "type-definition"}) { continue }
	// register the declaration, its binders and one synthetic formal-subject use per binder
	var walk func(semantics.TermID, int) error
	walk = func(t semantics.TermID, depth int) error {
		if depth > 512 { return ErrLimit }
		if old := termOwners[t]; old != 0 {
			if old != id { return semantic.ErrInvalidRequest } // a term shared by two bodies fails the whole session
			return nil
		}
		termOwners[t] = id
		for _, a := range terms[t-1].Arguments { walk(a, depth+1) }
		for _, o := range terms[t-1].Operands  { walk(o.Term, depth+1) } // callees are not walked
		return nil
	}
	walk(e.Body, 0)
}
```

`Open` and `Admit` combine navigation with horos judgments. In the excerpt, `i`, `e`,
`ti` and `formal` are the slot numbers decoded from the addresses:

```go
func (s *Source) Open(ctx context.Context, d g.Declaration[h.Language]) (g.OpenData[h.Language], error) {
	e := s.facts.Entities[i-1]
	if err := s.proper(s.results[i], e.Body); err != nil { ... }            // I: result sort proper
	o := g.OpenData[h.Language]{Declaration: d, Environment: env(i+1), Body: s.use(e.Body)} // N
	for _, p := range e.Parameters {
		if err := s.proper(s.sorts[Subject{Entity: p.Entity}], e.Body); err != nil { ... } // I
		o.Formals = append(o.Formals, g.Formal[h.Language]{Binder: s.bound[p.Entity], Subject: s.formalSubjects[p.Entity]}) // N
	}
	return o, nil
}

func (s *Source) Admit(ctx context.Context, env g.Environment[h.Language], use g.Use[h.Language]) (g.UseData[h.Language], error) {
	u := g.UseData[h.Language]{Use: use, Subject: use, Environment: env}
	if formal != 0 { // the synthetic formal-subject use
		if s.facts.Entities[formal-1].Owner != EntityID(e-1) { return u, g.ErrScope } // N
		u.Bound = s.bound[formal]
		return u, nil
	}
	if s.termOwners[ti] != EntityID(e-1) { return u, g.ErrScope }               // N: scope
	if err := s.proper(s.sorts[Subject{Term: ti}], ti); err != nil { return u, err } // I
	t := s.facts.Terms[ti-1]
	if t.Kind == Reference || t.Kind == Application {                            // N
		target := t.Target
		if t.Kind == Application {
			callee := s.facts.Terms[t.Callee-1]
			if callee.Kind != Reference { return u, s.Unsupported(use, "computed constructor application") }
			target = callee.Target
		}
		if b := s.bound[target]; b.Valid() { // a parameter
			if t.Kind != Reference || s.facts.Entities[target-1].Owner != EntityID(e-1) { return u, g.ErrScope }
			u.Bound, u.Subject = b, s.formalSubjects[target]
			return u, nil
		}
		u.Callee = s.decls[target]
		parameters := s.facts.Entities[target-1].Parameters
		if len(parameters) != len(t.Arguments) { return u, semantic.ErrInvalidRequest }
		for i, p := range parameters { // pairing actuals with real binders
			u.Actuals = append(u.Actuals, g.Actual[h.Language]{Binder: s.bound[p.Entity], Use: s.use(t.Arguments[i])})
		}
	}
	// The exact successful Horos analysis admits the entire original occurrence,
	// including every actual, under this declaration's original assumptions.
	return u, nil
}

// proper (I): the sort fact must be the term {horos, sort/proper}; an absent fact is
// reported as semantic.ErrInvalidRequest, an arrow sort as an Unsupported diagnostic.

func (s *Source) CheckChild(ctx context.Context, env g.Environment[h.Language], parent, child g.Use[h.Language]) error {
	p, _ := s.Admit(ctx, env, parent)
	if _, err := s.Admit(ctx, env, child); err != nil { return err }
	for _, a := range p.Actuals { if a.Use == child { return nil } }  // N: a generic actual
	v, _ := s.Read(parent)
	for _, o := range v.Operands { if o.Use == child { return nil } } // N: a descriptor operand
	return g.ErrScope
}
```

`DeclarationSelector` reads the module's identity fact and builds
`declaration/<pin>/<name>` (**I**). `UseSelector` decodes the descriptor operator and
payload discriminator (**I**). Beyond the contract, the horos session exposes a reading
API for horos mapping libraries:

```go
type Reader interface {
	g.SourceSession[h.Language]
	Declaration(g.Declaration[h.Language]) (Declaration, error)       // name, module pin, root-module flag, site, body
	Parameter(g.Binder[h.Language]) (Parameter, error)               // name and site
	IsBody(g.Environment[h.Language], g.Use[h.Language]) bool        // is the use the declaration's body? (N)
	Read(g.Use[h.Language]) (Expression, error)                      // constructor, discriminator, site, operands
	Operand(g.Use[h.Language], ...uint32) (g.Use[h.Language], error) // the operand at a locator
	Fields(g.Use[h.Language]) ([]Field, error)                       // record fields: name, optional, type use, site
	Unsupported(g.Use[h.Language], string) error                     // a diagnostic at the use's site
}
```

### How mapping libraries use the source

The default horos declaration recipe, shared by all four targets:

```go
Expand: func(ctx context.Context, r g.DeclarationRequest[h.Language, T, source.Reader]) (g.DeclarationDescription[T], error) {
	src, _ := r.Source()
	open, _ := r.Open()
	decl, _ := src.Declaration(open.Declaration)
	name, _ := Exported(decl.Name) // a native-safe identifier
	if !decl.Root { name = "M" + decl.ModulePin[:16] + "_" + name } // a declaration from an imported module
	body, _ := src.Read(open.Body)
	form := alias
	if body.Constructor == "record" { form = record }
	desc := g.DeclarationDescription[T]{Name: name, Form: form, Origin: decl.Origin}
	for i, f := range open.Formals {
		parameter, _ := src.Parameter(f.Binder)
		pname, _ := Exported(parameter.Name)
		desc.Generic.Parameters = append(desc.Generic.Parameters, g.ParameterDescription{Name: pname, Origin: parameter.Origin, RequiredUse: true})
		desc.Generic.Bindings = append(desc.Generic.Bindings, g.Parameter[T](i))
		desc.Generic.Arguments = append(desc.Generic.Arguments, g.Argument[T](i))
	}
	desc.Body, _ = r.Mapped(open.Body)
	return desc, nil
}
```

The TypeScript record rule (Go, Rust and C++ are analogous):

```go
Expand: func(ctx context.Context, r g.Expansion[h.Language, s.Language, source.Reader]) (g.Expression[s.Language], error) {
	src, _ := r.Source()
	use, _ := r.Use()
	if !src.IsBody(use.Environment, use.Use) {
		return g.Expression[s.Language]{}, src.Unsupported(use.Use, "nested records require a named declaration in this profile")
	}
	view, _ := src.Read(use.Use)
	fields, _ := src.Fields(use.Use)
	var members []g.Expression[s.Language]
	for _, f := range fields {
		ty, _ := r.Mapped(f.Use)
		members = append(members, g.At(model.Field(f.Name, f.Optional, ty), f.Origin)) // At: attach a source site
	}
	return g.At(model.Record(members...), view.Origin), nil
}
```

The list rule's capture reads the operand at locator `[1]`:

```go
Capture: func(ctx context.Context, src source.Reader, u g.UseData[h.Language]) ([]g.Capture[h.Language], error) {
	operand, err := src.Operand(u.Use, 1)
	return []g.Capture[h.Language]{{Slot: slot, Use: operand}}, err
}
```

The Cell library is keyed by Cell's declaration selector. It captures the actual
through the source's pairing:

```go
key, _ := source.Named(from, cellsModulePin, "Cell")
pattern := g.Pattern[h.Language, source.Reader]{Selector: key, Slots: []g.Slot{slot},
	Capture: func(ctx context.Context, src source.Reader, use g.UseData[h.Language]) ([]g.Capture[h.Language], error) {
		if len(use.Actuals) != 1 { return nil, g.ErrScope }
		return []g.Capture[h.Language]{{Slot: slot, Use: use.Actuals[0].Use}}, nil // by position
	}}
// Output: model.External("@bitspark/cell", "Cell", g.Mapped(slot))
```

### The independent source

To show that the kernel is not horos-specific, the repository has an independent test
source with its own toy language and its own mapping library. It publishes **no**
graph and requires no evidence. Its `Open` and `Admit` return `OpenData` and `UseData`
records built by hand in tables. Its `CheckChild` accepts a declaration's formal
subjects as children of the body, because its Page rule captures the formal subject
directly. It also builds a second occurrence with the same effective subject, to show
that the kernel must not equate occurrences by handle. It proves composability through
today's interface, but it re-implements the same navigation responsibility in a
different form. Under the target design it has to publish a graph and supply only its
own interpretation.

### Existing qualification

Several guarantees have no negative test that reaches the planner:

- The omitted-actual test fails earlier, at evidence admission. The planner's admission
  of an omitted actual has never been seen to fail.
- Foreign or forged records are untested.
- A wrong environment is untested.
- A `Capture` that returns a non-child is untested. `Expansion.Mapped` with a non-child
  is tested.
- Recursion is untested.

The replacement should be tested against those cases directly.

### What we've considered

**Earlier consultations did not consider kernel-owned traversal.** The consultation of
6 November 2026 that produced `Open`/`Admit` framed its constraint as "without turning
the agnostic host [kernel] into a semantic evaluator". Its expert wrote: "The source
integration is responsible for the meaning of the substitution, arity, saturation,
applicability, and required judgments. The driver [kernel] does not infer any of those
from the shape of `semantics.Term`." The consultation of 5 November 2026 said a
structural binding query "establishes correspondence, not applicability". We read both
as compatible with a kernel that follows structural correspondence and asks the source
for every judgment. Question 1 asks whether that reading holds.

**A candidate split.** This is ours, for critique, not a decision. It refines the
owner's intended division at the top of this document:

| Concern | Candidate owner | Notes |
| --- | --- | --- |
| Take the graph from the one evidence entry whose requirement sets `Semantics`, and unpack it typed | Kernel | Rules out a source consulting a different graph than the admitted one. |
| Resolve a requested root syntax site to the entity declaring it (`Entity.Declarations`) | Kernel | The source still classifies whether that entity is a generatable declaration. |
| Open a declaration: body term and real binders (`Entity.Parameters`) | Kernel | |
| Walk a body. A Reference to the declaration's own parameter is a bound use; to another declaration's parameter, a scope error. A Reference to a declaration entity is a zero-argument use of it. An Application with a direct Reference callee is a use of that callee, paired by `DirectParameterBindings`. | Kernel | |
| An application with no structural view (computed callee, parameter callee, unpaired arity), or a bare Reference to a declaration that has parameters | Kernel reports "no structural view", not a failed judgment | Arity and saturation *meaning* stay with the source. The kernel only notes that it cannot build a complete native application from the structure. |
| Child relation: a captured child must be an argument or operand of the current term (or the body of the declaration being opened) | Kernel | Replaces `CheckChild`. |
| Classify an entity as a generatable declaration and give its selector | Source | horos: kind `{horos, type-definition}`, pin plus name. |
| Give a descriptor's selector (constructor and discriminator) | Source | |
| Judge each declaration, formal and subject admissible under its declaration's assumptions, using supplied evidence only | Source | horos: the proper-sort checks. |
| Read payloads and operand meanings for mapping libraries | Source (its reading API) | horos: record fields, the list element at `[1]`. |

**Open points within that split** (each mapped to a question):

1. **Grain and coverage of admission** (Question 2). Should the source judge per term,
   per whole use (a term with all its actuals), or per declaration plus per term? The
   horos analysis has already judged the entire body when it reports success, so
   horos's per-term check today mainly enforces the proper-sort restriction. The walk is
   *demand-driven*: operands no mapping captures are never visited. For an omitted
   actual, should the kernel ask for one judgment of the whole actual, or walk its
   subtree? And does an operand that a mapping ignores need admission, as an omitted
   actual does?
2. **Outcome shape** (Question 2). There are four outcomes to keep apart: no structural
   view, rejected judgment, uncovered fact, and malformed input. Which are values and
   which are errors, and which belong to the kernel? Today an absent sort fact surfaces
   as an invalid-request error, which conflates "uncovered" with "malformed".
3. **Handles** (Question 3). Mapping libraries and the reading API could exchange
   graph-local term IDs that the kernel validates as structural children of the current
   subject. The output could then contain only this graph's structure, but a foreign
   handle would not be detectable as foreign. Alternatively, the kernel could mint
   run-scoped, owner-qualified handles that the reading API resolves through the
   kernel.
4. **Environment and scope** (Questions 1 and 3). Is scope "parameters owned by the
   declaration being opened" or an `Owner`-chain rule? A term may be shared between two
   bodies, and the graph permits that. Should such a term be accepted, with its
   environment taken from the walk path, or rejected as horos does today?
5. **Occurrence versus effective subject** (Question 3). Today `UseData` carries both,
   but the planner never compares effective subjects. Is a separate effective subject
   still needed, for example for future memoization or recursion, and who would own the
   equivalence?
6. **The mapping library's view** (Question 4). `Capture` and `Expand` need the current
   subject, the declaration being opened, the formals, and, for a declaration-keyed
   rule such as `Cell<V>`, the actuals. Should actuals be identified by real binder or by
   position?

**Rejected directions:**

- Moving `Admit` wholesale into the kernel. It mixes navigation with evidence and
  assumption checks that the kernel must not perform.
- Keeping `Open`/`Admit` and adding a shared helper that sources may call. The kernel
  would still depend on a source-supplied traversal, and nothing would prevent
  divergent implementations.
- Callbacks that let the source walk and call back into the kernel. That inverts
  control but leaves the traversal source-owned and weakens the child and scope checks.
- Moving the current reader to another package. That changes its location, not its
  responsibility.

## Questions for the Expert

1. **Where is the principled line?** Given a shared graph with three term variants
   (reference, application, provider expression with provider-defined operand
   locators) and provider-owned facts, how would you divide the generation walk between
   a language-agnostic kernel and a source-owned interpretation contract? That
   includes who owns scope: parameter ownership by the opened declaration, or an
   `Owner` chain. Please draw your own line before weighing our candidate. If yours
   resembles today's `Open`/`Admit` shape, the independent reason is what we need.
   Does our reading of the earlier advice ("correspondence, not applicability") hold?
   Which patterns from comparable systems, such as shared intermediate representations
   with language hooks or multi-language code-intelligence graphs, keep "follow the
   structure" from turning into "judge the structure"?
2. **How should admission judgments be shaped?** How would you model the grain of the
   source's judgments (declaration, whole use, term) and the outcomes (no structural
   view, rejected, uncovered, malformed)? Which of those outcomes does the kernel own,
   and does each suit a value or an error? The aims are that actuals a native convention
   omits are still admitted, that unknown stays distinct from failure, and that the
   kernel never interprets what a judgment means.
3. **What is a subject, and how is it identified?** Is a subject a graph term, an
   occurrence (a term reached under an environment, given that terms may be shared), or
   an effective subject whose equivalence the source owns? What handle design would you
   use so that capture identity and foreign-reference rejection hold, mapping callbacks
   cannot use subjects after their scope ends, and the source no longer re-supplies
   navigation records?
4. **What should mapping libraries see?** What structural view (current subject,
   declaration, formals, actuals) should source-typed mapping libraries receive
   alongside the source's own reading API? For a declaration-keyed rule such as a
   library-owned `Cell<V>`, should actuals be identified by real binder or by position,
   and what does each choice cost?
5. **What else would you watch for?** Looking at the design as a whole, including later
   work on recursive types and on carrying source bounds into native output, what risks,
   missing concepts or simplifications do you see? That includes any guarantee or
   constraint above that you think is wrong or missing.

## Applied

**Evaluated against:** `0013-generation-source-traversal.submitted.md` (sha256 `695dae3d45d35c2ff65559ba5a7f028537239d03ecc6053a958442eb7d7b7d8c`), attempt 1, run `run_f7cfb3b5-7cf2-40d8-8884-ce150a7adcf7`. Code at bittype `57187bcb0621c8da3ac1ef12e12887d4d079cd94`; the research branch changes only `research-docs/`.
**Graded:** https://github.com/Bitspark/bittype/issues/158#issuecomment-6721775550 (research owner claude-bittype-gaps; bittype has no architecture seat)

The verdicts are checked against the current code. "Kernel" means `contracts/generation/go` together with the registry in `transform/syntax/shared/generation/go`.

| # | Recommendation (advice section) | Verdict | Action |
|---|---|---|---|
| 1 | The kernel owns structural traversal: root resolution, declaration opening (body and real binders), reference following, variant inspection, application pairing, and child validation. The source classifies generatable declarations. (1.1) | HOLDS | Replace `Open`, `Admit`, `CheckChild`, `ResolveRoots` and `Authenticate` with one kernel walk over the admitted graph. The source keeps a declaration classification that returns a selector. |
| 2 | The kernel reads only structure: never `Entity.Kind`, `Entity.Name`, `Operator`, `Encoding`, `Data` or a locator's meaning. (1.2) | HOLDS | The kernel walk reads `Owner`, `Parameters`, `Body`, `Declarations`, the term variant, `Target`, `Callee`, `Arguments` and operand terms. Operand locators pass through unread. |
| 3 | Scope is direct binder ownership by the declaration being opened, not an `Owner`-chain rule. A language with nested scopes contributes its own scope judgment. (1.3) | HOLDS | A reference to an own parameter is a bound use, and a reference to another entity's parameter is `ErrScope`. Nested generic scopes stay unsupported until a language needs them. |
| 4 | The source judges whole uses under the assumptions of the enclosing declaration. (2.1) | HOLDS | A source use judgment receives a kernel occurrence, which carries its declaration. horos judges the use term's sort; its `program-validity` evidence covers the whole application. |
| 5 | Keep structural inspection separate from semantic judgment. (2.2) | HOLDS | The kernel builds the structural view. The source receives handles and returns judgments; it never returns navigation records. |
| 6 | Use a typed outcome model: structural availability as a value; holds, rejected and undetermined as judgment values with diagnostics; evidence, malformed input, scope and limits as typed errors. (2.3) | NEEDS ADAPTATION | A judgment is a value with four verdicts: Admitted, Unsupported (valid source the profile cannot represent; the arrow-sort case of 2.6), Rejected and Undetermined, each with an optional source diagnostic. The kernel maps Unsupported to `ErrUnsupported`, Rejected to a new `ErrRejected` and Undetermined to `semantic.ErrEvidenceUnavailable`. A zero judgment is `ErrInvalid`. "No structural view" is internal to the kernel and fails as `ErrUnsupported` at the term's site: no mapping can represent a use the kernel cannot pair. |
| 7 | Admit every actual, including those the native convention omits. Admission is separate from representation, and no callback can bypass it. (2.4) | HOLDS | Before any binding or convention applies, the kernel judges the parent use and then every actual occurrence. Represented actuals are traversed normally. Judgments are memoized per occurrence. |
| 8 | Provider-expression operands: the source judges the whole expression, the kernel validates captured operands as children, there is no eager generation, and any per-operand requirement is explicit in the source contract. (2.5) | HOLDS | The kernel validates captures only. horos's judgment of a descriptor occurrence covers its operands, because the horos analysis checks every operand's formation contract and the generation restriction applies to every term the walk reaches. The horos integration docs state this. |
| 9 | horos outcomes: a missing sort fact is undetermined, an arrow sort is unsupported, a proper sort is admitted, a failed membership is rejected. Never fill a gap with analysis. (2.6) | HOLDS | Today's horos `proper` returns `semantic.ErrInvalidRequest` for a missing fact. It changes to Undetermined. Arrow sorts give Unsupported with the existing diagnostic. Membership facts are not read today, so no rejection path exists for horos yet. |
| 10 | Distinguish graph terms, occurrences and effective subjects. Use opaque, kernel-issued, run-scoped occurrence handles that identify the run, graph, term, declaration and path. (3.1–3.2) | HOLDS | `g.Use[F]` becomes a kernel handle for one occurrence (run, term, declaration, parent). The kernel mints the children and checks run and parent identity. Kernel handles for declarations and binders replace the source `Space`/`Ref` records. |
| 11 | Keep the effective subject only if something consumes it. Never deduplicate or memoize by subject equality. (3.3) | HOLDS | Remove `UseData.Subject` and the formal-subject uses: the planner never compares subjects. Memoization covers declarations by entity and judgments by occurrence only. |
| 12 | Callback-scoped capabilities that expire on return. Callbacks receive neither the graph nor a graph-wide reader. (3.4) | HOLDS | `Expansion` and `DeclarationRequest` keep their active flag. They expose the kernel view, the source's narrow reader and `Mapped`. The graph and the source interpretation are not reachable from a callback. |
| 13 | Do not reject shared terms categorically. Two environments give two occurrences; a source may reject sharing its language forbids. (3.5) | HOLDS | Occurrences are path-based. The horos pre-walk that rejected cross-body sharing is removed: the horos exporter allocates a fresh term for every expression (`internal/horosmodel/export.go`), so the check never applied. |
| 14 | Mapping libraries get a kernel structural view: the current occurrence, its kind, the callee declaration, binder-keyed actuals, the bound binder, and the declaration's formals and body. (4.1) | HOLDS | Add `UseView` (occurrence, kind, callee, actuals, bound, operands, body flag) and `DeclarationView` (declaration, formals, body). The kernel builds both. |
| 15 | Keep a narrow source reader that consumes kernel handles. Its record view identifies field types as kernel-validated occurrences, and it returns no `UseData`. (4.2) | HOLDS | Mapping libraries are typed by a source reader `R`, not the session. The horos reader resolves kernel handles against its own graph. Its field and element readers select field-type occurrences from the kernel view's operands by horos locator. |
| 16 | Remove `CheckChild`. The kernel verifies callback ownership, run and graph, the structural child relation, and limits. (4.3) | HOLDS | `Mapped` and captured uses must be children the kernel minted for the current occurrence, from this run. A declaration recipe may map only its body occurrence. |
| 17 | Binder-keyed actuals. Position is only an iteration convenience. (4.4) | HOLDS | `Actual{Binder, Use}` in callee parameter order, with a lookup by binder. The Cell rule selects by Cell's binder `V`. |
| 18 | Keep the three mappings of the generic convention separate. (4.5) | HOLDS | `GenericConvention` is unchanged. |
| 19 | The Cell rule selects by binder and never opens Cell's body. (4.6) | HOLDS | The Cell test mapping selects Cell's formal through the kernel view of the callee declaration. The existing capture-count assertions stay. |
| 20 | An immutable evidence context built at admission: snapshot, admitted results, typed graph. The kernel unpacks the graph itself. (5.1) | HOLDS | Add a kernel capture step. It admits the request, requires exactly one evidence requirement with `Semantics`, and unpacks that result's graph as `Semantics[F]`. The source provider is opened over the capture. |
| 21 | Add language-neutral graph queries that are universal mechanics, without a traversal framework. (5.2) | NEEDS ADAPTATION | The kernel builds its run index from the graph (declaration sites, parameter owners, term sites). It adds no traversal framework to `semantics`; existing point queries and `DirectParameterBindings` suffice. A shared query is added only where a second consumer needs it. |
| 22 | Disciplined error ownership: `ErrScope` (structure or handle), `ErrUnsupported` (profile), `ErrInvalid` (malformed), `ErrLimit`. A rejection carries the source diagnostic, and uncovered is never a rejection. (5.3) | HOLDS | As row 6. The horos `Alias<text>` case becomes Unsupported (no structural view), not `semantic.ErrInvalidRequest`. |
| 23 | Memoization policy: declarations by identity, occurrences only with equal environment, subjects never, native expressions only where safe. (5.4) | HOLDS | Declarations by entity, judgments by occurrence, no expression memoization. |
| 24 | Keep recursion distinctions explicit, and leave future recursion to declaration reservation. (5.5) | HOLDS | Declaration-level active-set detection is unchanged, and a recursion test is added. Recursive generation stays bittype#164. |
| 25 | Keep room to expose source constraints in structured form for GEN-02, and never turn a proven membership into an unconditional native fact. (5.6) | DEFERRED | Trigger: GEN-02 (bittype#162) work starts. The judgment and reader contracts are source-owned, so a constraint reader can be added without kernel change. |
| 26 | Split mapping capabilities: declaration recipes, descriptor rules and application rules get only what each needs, not the whole session. (5.7) | NEEDS ADAPTATION | Callbacks receive the kernel view for their case (`DeclarationView` or `UseView`) plus the source reader `R`, a pure interpretation interface over kernel handles. They no longer receive the session or its interpretation. A separate reader per rule kind would duplicate one horos reading surface without removing any authority. |
| 27 | Test the invariants directly; the negative-test table is acceptance. (5.8) | HOLDS | The implementation PR adds those cases for horos and for the independent source, which publishes a real graph through its own analysis. |
| 28 | Implement in stages and migrate mapping libraries after the contracts settle. Finally remove `Open`, `Admit`, `CheckChild`, the source traversal indexes, source pairing and mapping access to navigation records, with no aliases. (6) | NEEDS ADAPTATION | One implementation PR, staged commit by commit in the advised order. The removal list is adopted in full, together with `Authenticate`, `CheckInput`, `Owners`, `RootEnvironment`, `Catalog.Declarations` and `ResolveRoots`. |

### Delivery

Owner of every accepted row: `claude-bittype-gaps`. Every accepted row is delivered by [bittype#191](https://github.com/Bitspark/bittype/pull/191).

**Validation method.** Validation names the test that fails on the broken behaviour. "Red" names the mutant seeded into that kernel check, which turned the test red in a local run. "Green" means the test passes on the PR head in CI: [checks of #191](https://github.com/Bitspark/bittype/pull/191/checks).

**Adaptation after review.** Two independent code reviews of the delivery found gaps against rows 7 and 8, which say an omitted part must not evade admission.

- The first review found that omitted actuals skipped the structural scope check, and that content nested in them was not judged.
- The second found that a body the recipe never represents was not admitted.

The kernel now admits each occurrence with its whole subtree before representing it: actuals, nested content and uncaptured operands. It also admits every generated declaration's body before the recipe runs. Representation stays demand-driven.

The judgment of a use covers that use. A callee declaration is classified when one of its uses is represented, and judged when it is generated.

**Adaptation after verification.** Independent verification by `codex-hydrated-type-design` (board claim `bittype-source-traversal-20261108`, record 000003) rejected #191.

- **Row 20 was not fully delivered in #191.** Advice 5.1 places the provider inside the immutable, kernel-owned evidence context. #191 built the capture but let callers of the exported `Plan` pass an already opened source next to it. A source opened over other evidence of the same graph could therefore plan. Now `Plan` takes the input and the provider, admits the evidence and opens the provider itself. Handles resolve against the exact capture (`In(Captured)`), not the graph.
- **Selector identity.** Rule lookup compared selector keys and dropped the language. A selector the source returns must now be in the capture's source language and profile.
- **Diagnostic site.** A term recorded at several syntax sites is reported at the occurrence's declaration.

The correction is [#192](https://github.com/Bitspark/bittype/pull/192). Independent verification by `codex-hydrated-type-design` (record 000005) accepted it at `6132170` on 9 November 2026.

| Rows | Validation (red → green) |
| --- | --- |
| 1, 2, 5, 14 (kernel walk, structure-only reads, separate judgment, kernel views) | `TestKernelTraversalIndependentNestedLibraryCell` and `TestDeclarativeIndependentSourceAndExtraParameter`. The independent source publishes a graph through its own analysis and drives the same walk as horos. The horos pagination, composition and placement suites stay green. The architecture test now allows `contracts/generation/go` to import `semantics/go`, and only that. |
| 3, 13 (direct binder scope; shared terms as distinct occurrences) | `TestKernelScopeIsDirectBinderOwnership`, which covers Stray, Borrow through a shared term, Lend, and SameA/SameB. Red: mutant M2, which treats any binder as bound. |
| 4, 7, 8, 9 (whole-use judgment, omitted actuals, operands, horos outcomes) | `TestKernelAdmitsOmittedActuals`: Kept is admitted, Leak is rejected, Blind is uncovered, Deep has nested refuted content, Smuggle and SmuggleLoose pass a foreign binder and an unapplied generic as omitted actuals, Hidden has a foreign binder in an uncaptured operand, and Shadow has a foreign binder in a body its recipe never represents. Red: M1 (admission skips actuals), M1b (admission skips operands), M4 (undetermined counted as admitted) and M10 (an unmapped body is not admitted). `TestHorosUnsupportedFormsAtTheirSite` covers a constructor-sorted parameter. |
| 6, 22 (outcome model, error ownership) | `TestKernelAdmitsOmittedActuals` checks `ErrRejected` against `semantic.ErrEvidenceUnavailable`. `TestKernelStructuralLimits` and `TestHorosUnsupportedFormsAtTheirSite` check that an applied alias, a curried application and a constructor-sorted parameter are Unsupported diagnostics located at the form, not malformed input. Red: M11 (horos diagnostic becomes a plain string). `TestKernelMapsEveryVerdict` covers a rejected declaration (Refused) and a zero verdict (Void). `TestVerdictKeepsOneFailureClass` covers a diagnostic that wraps a kernel class, and a located diagnostic that keeps its site. Red: M5 (rejected reported as undetermined), M6 (structural application check removed), M7 (declaration judgment skipped), M8 (zero verdict admitted), M9 (a diagnostic may add a class) and M12 (the explanation wrapper unwraps). |
| 10, 12, 16 (kernel occurrence handles, callback lifetimes, `CheckChild` removed) | `TestKernelCapturesAreChildrenOfThisRun` covers non-child captures through `Mapped` and `Pattern.Capture`, a recipe mapping a non-body use, and a handle kept from an earlier run. `TestPlanningExpansionScopeAndNoFallback` (`stale`, `swallowed-scope-error`) covers lifetimes. Red: M3 (`Mapped` accepts non-children), M3b (pattern captures unchecked) and M3c (a recipe may map any valid use). |
| 11, 23 (no effective subject; memoization policy) | `UseData.Subject` and the formal-subject uses are removed. Judgments are memoized per occurrence, declarations per entity. `TestKernelScopeIsDirectBinderOwnership` shows that one shared term is valid in SameA and SameB but a scope error through Borrow. |
| 15, 26 (narrow reader, capability split) | Mapping libraries are typed by the reader `R` and receive views plus that reader. The horos reader and the fixture reader are unexported values. All horos target suites stay green on the reader. |
| 17, 19 (binder-keyed actuals; Cell body unopened) | `TestGenerationPlanReversedLibraryArguments` selects Pair's actuals by binder names `L` and `R`. `TestLibraryOwnedCellAvoidsBodyExpansion` and `TestGenerationPlanPageCellSnapshot` select by Cell's binder, with one capture call. |
| 18 (generic convention) | Unchanged. `TestDeclarativeIndependentSourceAndExtraParameter` covers the extra Storage parameter. |
| 20 (evidence context) | `TestKernelRootsAndEvidence` covers foreign and missing evidence (`semantic.ErrInvalidRequest`), plus roots that are not a generated declaration (`ErrScope`). `TestGenerationPlanPreflightAndSourceEvidence` stays green. After verification: `TestKernelPlansExactlyTheAdmittedEvidence` plans two inputs that share a graph and differ in one report, through the registry and the exported planner. It also checks that an interpretation opened over the other input cannot resolve the run's handles. Red: M13 (handles resolve by graph) and M17 (source opened over a second admission). `TestKernelSelectorsAreQualified` varies the profile, the identity and the Go marker. Red: M14, M14b, M19 (profile only) and M20 (identity only). `TestKernelLocatesSharedTermsAtTheirDeclaration`. Red: M15 (first site wins). The strengthened `TestHorosUnsupportedFormsAtTheirSite` is red on M16 (term sites never used). `TestKernelLocatesUseFailuresAtTheUse` records the use's own site. Red: M18 (callee failures at the callee), M21 and M22 (use failures lose their term). `TestKernelLocatesSharedTermsAtTheirDeclaration` also covers both record orders and a repeated site. Red: M23 (the last site wins) and M24 (any repeated record is shared). `TestKernelChecksDescriptionsBeforeEvidence`. Red: M25 (admission first). `TestHorosReaderRefusesForeignHandles`. Red: M26 and M27 (horos reader ignores `In`). |
| 21 (graph queries) | The kernel run index is built in `occurrence.go` (`newWalk`). The `semantics` package is unchanged. |
| 24 (recursion distinctions) | `TestKernelStructuralLimits` (Rec reports recursive native declaration). |
| 27 (negative tests as acceptance) | `integration/go/generation_traversal_test.go`. |
| 28 (staging and removal) | The commits of #191, in the advised order: contract and migration, qualification, then docs. The removal list is in the PR description. |
| 25 | DEFERRED. Trigger: GEN-02 (bittype#162) work starts. |
