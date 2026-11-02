# Contract types above the Horos data model

**Owner decision, 2 November 2026: Ontos and Horos stay data-only; Bittype adds
the Wire type; hydration stays above Horos.** The type language below is a
proposal, with no implementation. It replaces the earlier proposal to put
`horos/wire` in Horos declaration revision 2. That format was never implemented
and is not retained as an alias or migration input.

The current [Horos frontend](horos.md), declaration revision 1, native text
stages and agnostic engine remain unchanged. This document specifies a separate
Bittype-owned contract-language layer; it does not revive the superseded
independent Bittype data-type algebra.

The owner's 3 November follow-up is specified in
[template lifting and native representations](type-representations.md): partial
lifting rules, generic pagination generation and library-owned target recipes
for TypeScript, Go, Rust and C++. Those are also documentation proposals only.

## 1. Ownership and value domains

| Component | Responsibility |
| --- | --- |
| Deixis | Generic tree structure over opaque contents |
| Ontos | Ground data values, structural identity and encoding |
| Horos | Data descriptors, data declarations, generic binding and conformance |
| Bittype contract-language provider | Wire message types and composition of data types with capability types |
| Bittype agnostic kernel | Generic AST/CST, semantic facts, registration and transformations without concrete type semantics |
| Bitwire / bitruntime | Capability contracts and hydration realization above Horos |
| Contract libraries | Function, Stream, Cell declarations, native representations and behavioral laws |

[Deixis permits opaque capability slots](https://github.com/Bitspark/deixis/blob/e0b0410d6f76d795992cc0a5792c507a19ef05ee/docs/SLOTS.md)
without requiring them to be portable data. The
[Ontos core](https://github.com/Bitspark/ontos/blob/a3f1d081f2e82f1e0ec39688c4060b5606b16c63/docs/spec/ontos-core.md)
instead defines finite, immutable, context-free ground values. Horos retains
[that data-only boundary](https://github.com/Bitspark/horos/blob/main/docs/spec/hydrated-types.md).

```text
GroundValue   = ByteAtom | Tuple(GroundValue*)
HydratedValue = ByteAtom | HydratedTuple(HydratedValue*) | Wire
```

The second domain already belongs to the
[Bitwire hydrated contract](https://github.com/Bitspark/bitwire/blob/e2d2d641744a49f1d8098b6ea5f2724454d72b96/docs/wire/hydrated-native.md)
and its runtime realization. Ground values embed unchanged. A Wire is an opaque
live sending capability, not an Ontos atom or an arbitrary object with a send
method. Hydration conveys recognized capabilities through exports and proxies.
Bittype describes these values; it does not hydrate them or implement a new
reference protocol. Receiving, closure and resource lifetime remain runtime and
endpoint-owner concerns.

Every type description below is itself ground Ontos data. It contains syntax,
not a live endpoint. That makes it encodable, not a Horos descriptor: the
language interpreting the encoded description still matters. Loading or
analyzing it performs no sends, exports, binding or endpoint allocation.

## 2. Exact proposed Wire expression and ground embedding

Use the ground value abbreviations `A`, `T`, `I`, `L`, `M`, `H` and `Pin` from
[Horos's notation](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md#1-representation-and-notation).
They construct raw UTF-8 atoms, canonical text/integers/lists/maps, Horos tags
and exact SHA-256 document pins. Add one Bittype-owned syntax abbreviation:

```text
C(tag, arguments...) = Tuple(Atom(UTF8("bittype/contract/" + tag)), arguments...)

Wire(E) = C("wire", E)
        = Tuple(Atom(UTF8("bittype/contract/wire")), E)
Data(G) = C("data", G)
```

These are proposed exact value equations, not implemented builders or a readable
source grammar. Wire takes exactly one type operand E, at tuple path `[1]`,
declaring its message type. It conveys sending authority only. There is no
implicit message type, direction flag, result slot or behavioral-contract name.

G must be a complete, well-formed, closed Horos ground descriptor. `Data(G)`
admits only ground values conforming to G. The Bittype provider delegates G's
meaning to Horos; no second scalar, refinement, equality or ground collection
model is introduced. In particular:

- `Data(H("any"))` excludes every value containing a Wire.
- `Data(H("not", G))` takes the complement within ground Ontos values only.
- `Data(H("exact", value))` keeps the entire literal ground and opaque to
  contract-parameter substitution, even if it resembles a Wire expression.

G is a separate language boundary. A contract parameter or application cannot
be inserted into its descriptor positions and handed to Horos as a descriptor.
Existing Horos declaration tooling can specialize data applications to G; an
AST transformation must retain the original declaration, arguments and origins
alongside that result. A future direct import bridge is separate scoped work.

## 3. Structural composition and generic declarations

The initial proposed contract expression grammar is deliberately bounded:

```text
E = C("data", G)
  | C("wire", E)
  | C("tuple", E...)
  | C("compound", A(label), E...)
  | C("list-of", E)
  | C("record", A("closed"), C("field", T(name), A(presence), E)...)
  | C("union", E, E, ...)
  | C("param", I(index))
  | C("apply", target, L(E...))

target = C("local", T(name)) | C("imported", T(alias), T(name))
presence = "req" | "opt"
```

Tuple and compound accept zero or more children, list-of exactly one, and union
at least two alternatives. Records accept zero or more distinct named fields;
the mode is explicitly closed. The provider owns these lifted structural forms.
They reuse the existing tuple, labelled compound, list and text-keyed map
layouts, allowing hydrated children where Horos permits only data:

```text
tuple(v0, ...)       = HydratedTuple(v0, ...)
compound(label, ...) = HydratedTuple(A(label), ...)
list(v0, ...)        = HydratedTuple(A("list"), v0, ...)
record(k0: v0, ...)  = HydratedTuple(A("map"), HydratedTuple(T(k0), v0), ...)
```

Record keys remain ground canonical text, unique and ordered by their Ontos
codec bytes; required/optional membership follows the Horos record convention.
Only their field values are lifted. Union admits values of any listed branch.
A tuple with no live leaf retains the ordinary ground representation. This is
an explicit Bittype interpretation above Horos, not live values inserted into
Ontos containers or a widened Horos checker. Wire equality, ordering, hashing
and serialization are not introduced.

Other ground forms, including maps, sets, open records, complements and guarded
recursion, remain available inside Data(G) with their unchanged Horos meaning.
They are not implicitly lifted. This first contract proposal adds no hydrated
top/complement, live map keys/set elements, recursive named application or
session-type recursion. Unsupported lifts must report that limitation. Wire
message operands are schema children, but a live Wire remains a value leaf;
schema traversal never follows a receiver or performs interaction.

Proposed contract modules have their own tag and revision:

```text
C("module", I(1), T(moduleName), importsMap, declarationsMap)
C("import", T(moduleName), Pin(importedContractModule))
C("type", L(T(parameterName), ...), E)
```

Imports and declarations use canonical text-keyed maps. This proposal adopts
the existing [Horos declaration conventions](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md)
for identifier spelling, exact document pins, ordered unique parameter names,
zero-based lexical slots, local/imported targets, explicit argument arity and
acyclic import/application graphs. Bittype owns their interpretation for these
contract modules. All imported modules in this initial envelope are contract
modules; a Horos module cannot be relabelled or imported through this form.
The envelope is not Horos revision 1 or 2, and current tools accept none of it.

Each parameter ranges over contract expressions E. Substitution visits only
the declared contract-type positions; it never scans arbitrary Ontos literals
or the contents of Data(G) for matching tags. Resolution preserves declaration
identity, lexical binding and each application even when specialization could
produce equal shapes. Specialization retains the contract-language profile and
Data boundaries; absence of a Wire tag does not authorize conversion to Horos.

### Pagination and lifting

Horos's `Page<Person>` is still a data application and can be specialized to G
and embedded with Data(G). Its T slot cannot accept `Wire<Message>`.

A Bittype-owned lifting transformer can be designed to reuse the Page template
by mapping its record/list structure into the contract forms above and preserving
its T binder. The [bounded lifting rules](type-representations.md#structural-lifting-rules)
specify this proposed transformation. Conceptually, its separately identified
result is:

```text
Page<T> = record closed {
  items: list-of(T)
  next?: Data(scalar(utf8-text))
}
```

Here Page names the derived contract declaration, so `Page<Wire<Message>>` is
meaningful in the upper language. Users do not hand-author another Page, and a
`LiftedPage` spelling is not required. The transform records the supported forms,
input/output profiles, original declaration and provenance. Derived contract
modules have their own identities and exact import pins; they are not the
original Horos declarations with widened semantics. The build selects lifting
explicitly, with no fallback or blanket reinterpretation of Horos constructors.
Configuration syntax and implementation remain deferred; this is explanatory
notation. Native generation can retain this generic and compose its arguments
through the [selected target representations](type-representations.md#3-pagination-through-the-type-pipeline).

## 4. Library-owned Function, Cell and Stream

These complete proposed module equations illustrate composition. They specify
no deployed package or protocol compatibility. The abbreviated applications are:

```text
P(i)                 = C("param", I(i))
Use(alias, name, ...) = C("apply", C("imported", T(alias), T(name)), L(...))
Local(name, ...)      = C("apply", C("local", T(name)), L(...))

WireLibrary = C("module", I(1), T("example.wire"), M{}, M{
  "Wire": C("type", L(T("M")), C("wire", P(0)))
})

FunctionLibrary = C("module", I(1), T("example.contracts.function"), M{
  "wire": C("import", T("example.wire"), Pin(WireLibrary))
}, M{
  "Unit": C("type", L(), Data(H("tuple"))),
  "Function": C("type", L(T("A"), T("B")),
    Use("wire", "Wire", C("tuple", P(0), Use("wire", "Wire", P(1)))))
})

CellLibrary = C("module", I(1), T("example.contracts.cell"), M{
  "fn": C("import", T("example.contracts.function"), Pin(FunctionLibrary))
}, M{
  "Cell": C("type", L(T("T")), C("record", A("closed"),
    C("field", T("get"), A("req"),
      Use("fn", "Function", Use("fn", "Unit"), P(0))),
    C("field", T("set"), A("req"),
      Use("fn", "Function", P(0), Use("fn", "Unit")))))
})

StreamLibrary = C("module", I(1), T("example.contracts.stream"), M{
  "fn": C("import", T("example.contracts.function"), Pin(FunctionLibrary))
}, M{
  "Step": C("type", L(T("T")),
    C("union", C("compound", A("item"), P(0)), C("compound", A("end")))),
  "Stream": C("type", L(T("T")), C("record", A("closed"),
    C("field", T("next"), A("req"),
      Use("fn", "Function", Use("fn", "Unit"), Local("Step", P(0)))),
    C("field", T("cancel"), A("req"),
      Use("fn", "Function", Use("fn", "Unit"), Use("fn", "Unit")))))
})
```

Each Pin is the exact 32-byte SHA-256 atom of the preceding expanded module's
canonical Ontos bytes, not a placeholder string or version lookup. Module names
are illustrative library identities, never built-in registrations. Their meaning
in explanatory notation is:

```text
Unit = Data(tuple())
Function<A,B> = Wire<tuple(A, Wire<B>)>
Cell<T> = record closed { get: Function<Unit,T>, set: Function<T,Unit> }
Step<T> = union(compound("item", T), compound("end"))
Stream<T> = record closed {
  next: Function<Unit,Step<T>>
  cancel: Function<Unit,Unit>
}
```

Function is an illustrative success-only request/reply representation. Its
library owns fresh reply endpoints, matching, one-reply discipline, failure,
completion and resource management. Unit is a bare empty tuple, not an implicit
native void. Cell's library owns state and mutation guarantees. The pull Stream
library owns demand, order, concurrent-next policy and terminal/cancellation
races. These behaviors are not consequences of the type shapes. Production
bindings need their own full protocols and qualification.

A future generator preserves resolved library identity through imports, aliases
and applications. The [native representation proposal](type-representations.md#4-ordinary-libraries-and-optional-native-representations)
specifies library-owned type recipes selected by exact identity and configured
on the concrete language transformers. Ordinary structural compositions need no
custom recipe. Matching an expanded shape or a name such as Cell cannot select
a representation or binding. Type recipes are separate from future executable
adapters; neither is implemented. Adding another contract changes the library
and consuming composition, not the generator core, runtime, Wire or hydration.

## 5. Provider, pipeline and dependencies

The contract language must have its own Go language marker, language ID and
qualified syntax/resolution/analysis profiles, distinct from `horos` and its
implemented `horos-declarations-1/ast-2`, `horos-declarations-1/resolution-1`
and `horos-declarations-1/analysis-2` profiles.
Exact registered identifiers are assigned with implementation. An engine instance
registers the chosen providers explicitly; no name switch or default is added.

The [stage-first layout](../LAYOUT.md) applies. Proposed paths use `contract`
as the stage-directory name, not a new generic kernel package:

| Proposed component | Responsibility |
| --- | --- |
| `syntax/contract/go` | Contract-language vocabulary and identity on the existing AST carrier |
| `check/contract/go` | Syntax-only admission; unresolved names can reach analysis |
| `resolve/contract/go`, `internal/contract` | Module/binder resolution, explicit ground embeddings and declared operand positions |
| `analyze/contract/go` | Provider-owned semantics adapted to shared immutable facts and reports |
| `provider/contract/go` | Explicit AST-only capability composition |
| `transforms/contract-to-<native>/go` | Future native type AST mappings with explicit representation profiles |
| `transforms/horos-to-contract/go` | Future, separately qualified data embedding or structural lifting |

These are proposed paths only; no packages or APIs are added now. Existing
`syntax/horos`, `resolve/horos`, `check/horos`, `analyze/horos` and
`provider/horos` remain data-only. Horos-to-native mappings remain independently
useful for data generation; a capability layer is not a mandatory detour.

The generic semantic model can carry provider-owned expressions with located
operands and ordinary Reference/Application facts. The contract provider owns
the Wire message operand and preserves each Data boundary; the kernel checks
graph integrity, storage and publication without decoding those semantics.
No universal semantic Wire term, structural fingerprint or core import follows.

The contract resolver/analyzer may delegate ground semantics to Horos through
reviewed, narrowly qualified package dependencies. Current architecture checks
permit those external model dependencies only in the Horos frontend; the future
provider needs its own explicit allowance and tested capability profile.
Generic trees, `contracts/*`, semantic facts, engine, printer and native text
stages remain independent of Horos, Wire and behavioral libraries. The existing
`contracts/` directory contains agnostic SPIs, not this proposed language.

Static tooling has no dependency on a live runtime. Behavioral library
implementations use existing Wire/hydration; consumers compose selected libraries.
The dependency never points from Horos to Bittype or from hydration to Function,
Stream or Cell. Existing formatters/printer consume future target ASTs normally.

## 6. Qualification and forbidden shortcuts

This task changes documentation only. Current Horos APIs and the Bittype Horos
provider reject the proposed contract format. No generator, runtime binding,
parser, validator, production dependency, conformance vector or release is added.

Before implementation is claimed, independent fixtures must establish exact
namespace/revision/arity checks, ground Data boundaries, generic/literal scope,
resolved library identity, profile-preserving specialization, supported lifting
and separate syntax/semantic outcomes. Include wire-free `any`/complement cases
to prove that domain selection is not inferred by searching for a Wire tag.
The data-only Page and its explicitly lifted counterpart require separate
identity and meaning checks. Native compilation and runtime protocol tests are
separate evidence, not substitutes for those checks.

Recognizing a Wire proves only runtime-admitted sending capability. `Wire<E>`
declares a message contract; observing the leaf cannot prove that a remote
receiver follows E. Future bindings may validate actual messages. Analysis
reports cannot establish remote execution, one reply, cancellation, variance or
protocol equivalence. Type analysis itself must never send, bind or export.

Forbidden: adding Wire to Horos or the agnostic kernel; silently reinterpreting
ground descriptors or generic arguments; embedding live handles in Ontos data;
core Function/Stream/Cell kinds, imports, switches or registrations; behavior
chosen by structural fingerprints, including generated central dispatch tables;
and a replacement hydration/export engine. The new ownership boundary belongs
in a concrete language provider, not in those general layers.
