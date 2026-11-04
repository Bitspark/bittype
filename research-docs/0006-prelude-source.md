# Research: Authoring the Horos prelude in Horos

**ID:** 0006
**Date:** 4 November 2026
**Status:** evaluated — design consultation complete; implementation not requested
**Owner:** codex-hydrated-type-design
**Scope:** Architectural explanation and recommendation before implementation.
**Run-ID:** run_c97ae94d-1d80-4ea5-b039-a7835103cc6d
**Document-ID:** doc_2e194bc2-3e54-47f8-bc08-1dbc58b1371d
**Consultation key:** bittype-research-0006-f4f86a35-a5df-4149-8fc5-df68af9a314c

## Evaluated recommendation

The [expert advice](0006-prelude-source.advice.md) agrees that the standard kinds
should be maintained in Horos source. The Go construction was an authoring and
packaging shortcut, not a requirement of the kind model. Its output already
uses the ordinary canonical module loader and general kind calculus.

Adopt this design for a subsequent implementation:

```text
Horos: authoritative horos.kinds.horos
  -> Bittype: ordinary Horos parser
  -> Bittype: existing syntax-to-canonical module encoder
  -> Horos: generated canonical Ontos module asset
  -> Horos: ordinary Ontos decoding and module admission
```

The smallest public tool is an ordinary single-module encoder in Bittype. It
must not know the name `horos.kinds`, invoke semantic resolution, or introduce a
second parser. A small private helper can expose the existing `reader.module`
conversion to both the command and `SourcePins`. A separate public Go API is not
needed until a library caller actually requires it.

Keep the complete maintained source and generated asset together under
`horos/core/go/declarations/standard/`, with proposed names
`horos.kinds.horos` and `horos.kinds.ontos`. Horos can embed the asset and retain
the useful `Prelude`, `PreludePin` and `PreludeImport` accessors. Remove
`buildPrelude` and the helpers used only to hand-author its definitions. The
compiled asset is the runtime distribution form; only the Horos source is edited.

There are three separate relationships:

| Relationship | Meaning |
| --- | --- |
| Selected formation dependency | Primitive rules use the configured standard module's Data definition. |
| Explicit source import | A user can refer to names such as h.Data. |
| Supplied source tree | Tools have actual syntax sites and text locations. |

Selecting the formation dependency remains legitimate language composition. The
kind evaluator obtains Data's meaning from the ordinary declaration; it does not
recognize the name. An unrelated local or imported Data cannot replace that
binding. Changing the selected standard Data definition can affect formation
even in programs that never explicitly import it. Importing another version of
the module does not change the resolver's selected formation dependency.

Source authorship does not guarantee source availability in every invocation.
Retain the three narrow, exact-module-identity accommodations for an absent
prelude source tree. When callers supply the actual source through ordinary
Syntax/Files inputs, use its real locations. Generated AST projection must not
pretend to recover original source offsets. General canonical-only dependencies
are a separate capability, not a reason to relax all missing-site errors now.

### Bootstrap and update sequence

1. Add the general source-to-canonical command to Bittype against the already
   available Horos release. Obtain that executable before preparing the new
   Horos package.
2. Author the complete standard module in Horos and use that available executable
   to produce its canonical asset. This is syntax lowering, with no call to
   `resolve.Model` or admission against the installed prelude.
3. Replace Go-authored bodies with embedded-asset decoding and ordinary loading.
   Ship both maintained source and generated asset in the Horos Go module.
4. Update Bittype's Horos dependency afterward. Downstream Horos builds consume
   the packaged asset and do not invoke Bittype.

The preparation executable may contain or initialize the older Horos prelude
because of its existing package dependencies. The relevant boundary is that
lowering the new module never resolves its definitions through that prelude.
Using an already available executable also avoids a workspace redirect that
would build the tool against an incomplete target missing its embedded asset.

### Disposition of the advice

The research owner evaluated the recommendations against the source paths in
this brief. This is architectural analysis, with no validation/check workflow.
The following are design dispositions; none claims an implementation change.

| Recommendation | Verdict | Reason and proposed action |
| --- | --- | --- |
| Author all standard kind definitions in Horos | HOLDS | Existing kind syntax expresses the complete module. Move all exports, not just the illustrative subset. |
| Keep primitives and ordinary formation-reference wiring in Horos | HOLDS | The calculus supplies meaning to source definitions; grammar.go selects an ordinary scoped Data reference without a Data-specific evaluator. Preserve that selection. |
| Add a general single-module encoding command | HOLDS | SourcePins already uses reader.module before resolution. Expose that conversion internally and encode the resulting value; keep the command ignorant of prelude names. |
| Use a public Go lowering API | DEFERRED | Add it when a concrete library caller needs canonical module values from admitted ASTs; the command alone meets this use case. |
| Prepare with an available Bittype executable, then package the asset | HOLDS | Keeps runtime dependencies one-way and avoids building the incomplete Horos target to create its own required asset. |
| Remove authored Go bodies while retaining useful accessors | HOLDS | Prelude/PreludePin/PreludeImport remain ordinary access conveniences; buildPrelude and definition-only helpers are the superseded implementation. |
| Preserve canonical identity semantics | HOLDS | Existing pins address canonical values. Equivalent kind meaning alone does not imply unchanged module identity; no new identity scheme is needed. |
| Retain exact-prelude source accommodations for canonical-only use | HOLDS | Changing authoring does not supply a source tree to every session. Keep the absent-source case narrow; supplied source follows ordinary correspondence. |
| Expose embedded source text for standalone navigation | DEFERRED | Add only when a concrete standalone integration requires navigation without a caller-supplied source file. Reuse the maintained source, with no duplicate literal. |
| Generalize canonical-only dependencies and source maps | DEFERRED | Generalize dependency inputs when a non-prelude canonical-only dependency is required. Add source-map machinery only for a consumer that cannot obtain needed original locations from supplied source. |
| Move the parser below Horos, add configurable prelude replacement, or require manual standard-library assembly everywhere | REJECTED for this change | These expand ownership or public configuration without addressing a missing capability in this small source-authoring correction. Reconsider only for a separate concrete requirement. |
| Document the selected Data's formation role and language coverage limits | HOLDS | A library-defined contract can still be part of the selected language semantics. A portable asset does not implement revision-4 declarations in Rust or TypeScript. |

No production code, package versions or release state changed in this
consultation. Implementation remains a separate next step.

## Question

The owner asks: “Why are these definitions constructed in Go? Shouldn't the
prelude be written directly in Horos and loaded through the same machinery as
any other module?” The definitions are Scalar, List, Record and Data, which the
recent general-kind implementation represents as ordinary canonical declarations
but currently authors with Go constructor calls.

Recommend the smallest coherent authoring, compilation, distribution and loading
design. Account for the real bootstrap dependencies, and distinguish a standard
library dependency from a special semantic evaluator. An explanation that only
moves the Go table to another package will not answer the owner’s concern.

This is a design consultation. No runtime or compiler implementation is requested
in this turn. The owner prohibits backup, checksum, validation/check and unsolicited
provenance workflows; do not propose those as delivery prerequisites. Semantic
admission performed by the language remains product behavior. Module identities
described below are existing language semantics, not a proposed verification scheme.

## Context

### System and terminology

Bitspark has three relevant independently released repositories:

| Repository | Responsibility |
| --- | --- |
| ontos | Immutable Atom/Tuple values, data embeddings, canonical encoding and structural identity. |
| horos | Data descriptors and value conformance; a Go declaration layer for modules, generic types, applications and general structural kinds. |
| bittype | Go AST/CST infrastructure, Horos readable syntax, resolution/analysis providers, and native type transformations. |

An AST is an abstract syntax tree. A CST additionally retains concrete source
tokens/trivia. Bittype's shared kernel is language-agnostic: language providers own
their vocabulary, semantic interpretation and stage composition. A canonical
Horos module is an Ontos value; it is not Go source or a Bittype AST. Readable
`.horos` text is an authoring spelling for that canonical model.

A *sort* distinguishes proper types (`P`) from constructors (`(P, P) -> P`, for
example). A *kind* constrains type descriptions of a sort. Scalar, List and Record
classify description shapes; Data is a recursively defined kind of descriptions
that contain only ground-data forms. `top(P)` admits all proper types, including
Wire types. Data is not a reserved name or a dedicated proof procedure.

Structural kinds match normalized type observations: aliases/applications are
evaluated; designated type children are traversed; literal payloads stay opaque;
ground `rec`/`ref` binders become graph backedges. A description rooted at a union
is classified by its union shape, not by pretending it has each member’s shape.

There is one revision-4 declaration language. Ground Ontos/Horos values and ground
conformance remain data-only. Optional Ontos/wire owns mixed runtime values;
Bitwire/Bitruntime own hydration. Function, Stream and Cell remain library-owned
contracts. None is a special case in the kind engine, Bittype kernel or runtime.

The current Go Horos module uses Ontos v0.16.0. Bittype depends on Horos's Go module
v0.10.1 and Ontos v0.16.0. Horos does not depend on Bittype. Other Horos language
implementations currently support ground conformance, not this declaration layer.

### Existing canonical module and kind representation

In the following notation, `D(tag, ...)` means an Ontos tuple headed by the atom
`horos/decl/<tag>`. `text`, `integer`, `list` and `map` are existing Ontos data
embeddings. Module identity is the digest of its canonical Ontos encoding.
Imports carry the expected module name plus that identity. Different local aliases
do not make different imported declarations; same-spelled definitions in another
module are not the standard definitions.

```text
Module = D("module", integer(4), text(moduleName), importsMap, declarationsMap)
Import = D("import", text(expectedModuleName), identityAtom)

Sort = D("sort-proper")
     | D("sort-arrow", list(Sort...), Sort)

KindParameter = D("kind-parameter", text(name), Sort)
KindDefinition = D("kind", list(KindParameter...), Sort, Kind)

Kind = D("top", Sort) | D("bottom", Sort)
     | D("kind-param", integer(slot))
     | D("kind-local", text(name))
     | D("kind-imported", text(alias), text(name))
     | D("kind-apply", namedKindReference, list(Kind...))
     | D("or", Sort, list(Kind...)) | D("and", Sort, list(Kind...))
     | D("node", text(form), selectorMap, Sequence)
     | D("nu", list(Kind...), integer(entry))
     | D("recursive-var", integer(depth), integer(slot))
     | D("arrow", list(Kind...), Kind)

Sequence = D("epsilon") | D("item", Kind)
         | D("concat", list(Sequence...))
         | D("alternative", list(Sequence...))
         | D("star", Sequence)
```

`epsilon` means zero type children; `item(K)` one child satisfying K;
`star(item(K))` any finite sequence of them. `nu` is a guarded greatest fixed
point: recursive references must descend through a node's type child. Named kind
definitions otherwise form a dependency DAG. The kind engine implements these
constructs generically; it does not recognize the names Scalar, Record or Data.

### Actual Go authoring

Horos `core/go/declarations/prelude.go:12-27` currently contains:

```go
var standardPrelude = buildPrelude()
var standardPreludePin = Digest(sha256.Sum256(codec.Encode(standardPrelude)))

func Prelude() core.Value { return standardPrelude }
func PreludePin() Digest { return standardPreludePin }

func PreludeImport() core.Value {
    pin := PreludePin()
    return preludeForm("import", preludeText("horos.kinds"), core.NewAtom(pin[:]))
}
```

`core.Value` is an Ontos value; `Digest` is the existing module identity type.
`preludeForm` constructs a tuple with the `horos/decl/` tag, and helpers encode
ordinary text, integer, list and map values. `buildPrelude` at lines 72-145 uses
these helpers to author the library. These excerpts show representative bodies:

```go
define("Type", preludeForm("top", proper))
define("Int", selected("scalar", "selector", "int", epsilon))
define("Text", selected("scalar", "selector", "utf8-text", epsilon))
define("Bool", selected("scalar", "selector", "bool", epsilon))
define("Scalar", union(
    preludeReference("Int"), preludeReference("Text"), preludeReference("Bool"),
    leaf("text-len"), leaf("key-grammar"), leaf("text-non-blank"), leaf("lowercase-uuid"),
))
// The unary-family loop includes List -> list-of and Wire -> wire.
define(family[0], preludeNode(family[1], empty, item(parameter(0))), "K")
define("Record", preludeNode("record", empty, many(parameter(0))), "K")
// dataBody is an ordinary recursive union of the ground families below.
define("Data", preludeForm("nu", preludeList(dataBody), preludeInteger(0)))
return preludeForm("module", preludeInteger(4), preludeText("horos.kinds"),
    empty, preludeMap(entries...))
```

`define` produces a normal KindDefinition, `proper` is the proper sort,
`selected` adds a selector map to a node pattern, `union` creates `or`,
`leaf(form)` is a node with an empty selector map and epsilon children,
`many(K)` is `star(item(K))`, and `parameter(0)` refers to the first kind parameter.
There is no separate prelude declaration representation.

### The same definitions in the implemented readable syntax

This is a transcription of the current definitions into already-supported source
syntax, not an existing maintained source file. The whole module would also
declare the other referenced families listed in Data.

```text
module horos.kinds;

kind Type = top(P);
kind Int  = node("scalar", {"selector": "int"}, epsilon);
kind Text = node("scalar", {"selector": "utf8-text"}, epsilon);
kind Bool = node("scalar", {"selector": "bool"}, epsilon);

kind Scalar = Int | Text | Bool
  | node("text-len", {}, epsilon)
  | node("key-grammar", {}, epsilon)
  | node("text-non-blank", {}, epsilon)
  | node("lowercase-uuid", {}, epsilon);

kind List<K> = node("list-of", {}, item(K));
kind Record<K> = node("record", {}, star(item(K)));

kind Data = nu {
  D = Any | AtomFamily | Scalar | Exact
    | Tuple<D> | Compound<D> | Seq<D>
    | List<D> | Set<D> | Map<D, D>
    | ClosedRecord<D> | OpenRecord<D>
    | UnionOf<D> | AndOf<D> | Not<D>;
} in D;
```

All kind parameters/result sorts default to P. `Any` above matches the ground
`any` descriptor; it is not kind top. `Exact` treats its payload as opaque.
Record observations include optional fields; open records also contribute an
implicit `any` child. Data excludes Wire structurally. The module has only kind
definitions and no imports, so its own kind admission requires the primitive
shape vocabulary and general kind calculus, not Data-based type formation.

Consumers explicitly import the standard module:

```text
module example.pagination;
import h from horos.kinds pin "sha256:<actual-64-hex-module-identity>";

type Page<T> = { items: list<T>; next?: text; };
type DataPage<D: h.Data>: h.Data = Page<D>;
type Person = { name: text; active: bool; };
type People = DataPage<Person>;
```

The identity placeholder is explanatory, not executable source. There are no
implicit Type/Data source names. A local kind named Data does not change primitive
formation. An omitted type parameter or result contract elaborates to `top(P)`.

### Actual canonical loading and formation configuration

The public API accepts canonical values, not text. Horos
`core/go/declarations/model.go:91-100` delegates as follows (wrapper omitted):

```go
func ResolveWithLimits(root core.Value, dependencies []core.Value, limits Limits) (*Program, error) {
    p, err := model.Resolve(root, dependencies, limits, grammar())
    // Wrap or return the result/error.
}
```

`core/go/declarations/grammar.go:12-24` constructs this internal grammar:

```go
return model.Grammar{
    DeclarationPrefix: "horos/decl/",
    DescriptorPrefix:  "horos/",
    Any:               descriptorAny,
    Positions:         descriptorPositions,
    WF:                descriptorWF,
    GroundWF:          groundWF,
    Metadata:          descriptorMetadata,
    Schema:            descriptorSchema(),
    Prelude:           Prelude(),
    OperandKind:       descriptorOperandKind,
    GroundKind:        preludeReference("Data"),
}
```

The grammar fixes primitive arities, type-operand positions, literal metadata and
recursion rules. It supplies the ordinary prelude plus expressions for primitive
operand contracts. The kind engine is not given a new Data-specific judgment.
`descriptorOperandKind` at lines 105-117 does nevertheless name Data explicitly:

```go
func descriptorOperandKind(constructor string, args []core.Value) core.Value {
    switch constructor {
    case "seq", "set-of", "map-of", "and", "not", "rec":
        return preludeReference("Data")
    case "record":
        if len(args) > 0 {
            if mode, ok := args[0].(core.Atom); ok && string(mode.Bytes()) == "open" {
                return preludeReference("Data")
            }
        }
    }
    return preludeForm("top", preludeForm("sort-proper"))
}
```

These choices preserve existing ground/mixed-type admission. Tuple, compound,
list, closed-record, union and Wire message operands may be proper types; the
listed forms require the selected ordinary Data kind. Existing ground
well-formedness, literal checks and guarded recursion remain authoritative.

The internal resolver in `core/go/internal/declmodel/resolve.go:16-44`:

1. Places root, supplied dependencies and `grammar.Prelude` in one module pool.
2. Loads the prelude and root using the same `load` closure. It parses module
   values, follows pinned imports, checks expected names and detects cycles.
3. Reads kind expressions and ordinary references for all loaded declarations.
4. Admits kind definitions through the general kind calculus.
5. Resolves `GroundKind` in the prelude's scope and then checks type declarations.

Representative actual statements, expanded for readability:

```go
all := append([]core.Value{root}, dependencies...)
all = append(all, grammar.Prelude)
// pool and rootPin/preludePin are formed from these values.
if err := load(preludePin, 0); err != nil { return nil, err }
if err := load(rootPin, 0); err != nil { return nil, err }

preludeOwner := &Declaration{id: DeclarationID{Module: p.prelude}}
g, err := r.readKind(grammar.GroundKind, preludeOwner, nil, nil, 0)
// Instantiate and admit g using the general checker.
```

`internal/declmodel/evaluate.go:120-124` resolves each OperandKind in the same
prelude scope. A consumer's `Data` definition cannot substitute for it. The
prelude is always in the semantic closure even when the user has no explicit
import; that availability does not put its names into the user’s lexical scope.

Thus authoring and dependency injection are special, but parsing canonical module
values, binding kinds, membership and inclusion already use ordinary machinery.
Please do not assume there is a completely separate “prelude evaluator.”

### Readable syntax and the real dependency cycle

Bittype `parse/horos/go/parse.go:58-104` exposes:

```go
func Parse(ctx context.Context, input parsecontract.Input) (cst.CST[s.Language], error)
```

It parses one or more ordinary module sections into CST/AST, admits syntax,
and calls `horosmodel.SourcePins(tree)` to fill module identities. It does no
filesystem/network IO and does not semantically resolve names. `kinds.go` in the
same package implements node patterns, sequence patterns, `nu`, kind applications,
sorts, and Boolean/arrow expressions.

Bittype `internal/horosmodel/source.go:14-30` already encodes module syntax before
semantic resolution:

```go
func SourcePins[L any](tree ast.AST[L]) (map[ast.NodeID]string, error) {
    r := reader[L]{tree: tree, declarationPrefix: "horos/decl/", descriptorPrefix: "horos/"}
    out := make(map[ast.NodeID]string)
    for _, child := range tree.Nodes[tree.Root-1].Children {
        value, err := r.module(child.Node)
        if err != nil { return nil, err }
        pin, err := decl.Pin(value)
        if err != nil { return nil, text.AtNode(tree, child.Node, err.Error()) }
        out[child.Node] = pin.String()
    }
    return out, nil
}
```

`reader.module` in `internal/horosmodel/read.go:187-225` encodes import and
declaration maps, then produces exactly:

```go
return tagged(r.declarationPrefix, "module", data.EncodeInt(big.NewInt(4)),
    data.MustEncodeText(r.attr(id, s.AttributeName)), importsValue, declarationsValue), nil
```

This syntax-only encoder is private. There is no public API/tool that returns
canonical module values from admitted source syntax without semantic resolution.
Public `resolve.Model` currently resolves and returns a Horos Program, whose
modules expose `Value()`. Using that semantic path just to produce a new prelude
would unnecessarily involve the already-installed prelude. Exposing the existing
syntax-to-canonical conversion is a smaller possibility than writing a parser.

Current package dependency:

```text
Bittype parser -> Bittype private Horos model helper -> Horos declarations -> Ontos
```

Having Horos declarations import that parser to load source at runtime creates a
Go import cycle. The semantic prelude itself does not have a corresponding cycle:
it contains only kinds, whose primitive forms/sorts are already intrinsic.

### Source locations and Files composition

Bittype has a text pipeline (`text -> CST/AST`) and a Files pipeline that captures
supplied source files and typed semantic evidence. Resolution/analysis run
explicitly; providers do not silently retrieve imports. Multiple source trees can
supply module dependencies without concatenating source strings.

Because Horos always injects a prelude that may not have a supplied syntax tree,
current Bittype has exact-prelude-identity accommodations:

- `internal/horosmodel/read.go:540`: permits the prelude to be absent from source
  correspondence, while ordinary supplied modules must have syntax sites.
- `internal/horosmodel/export.go:179`: permits those prelude entities to have no
  source site. Otherwise an entity without a syntax site is an error.
- `internal/horosfiles/stages.go:208`: omits a physical file-import relation if
  the prelude source file was not supplied.

Supplying the prelude as a real source dependency already follows the ordinary
path and gives actual sites. Projection can also materialize an AST from the
canonical prelude, but those generated nodes do not establish original text
locations. The question includes how much of this boundary should change now;
please distinguish an authoring change from a broader compiled-module facility.

## Options considered

1. **Horos source plus canonical distribution asset.** Author the prelude once
   as `.horos` in Horos; use Bittype's parser and existing syntax encoder in a
   preparation/publication tool; package the resulting canonical Ontos module
   with Horos. Runtime decodes it and uses ordinary module admission. Remove
   `buildPrelude` and authored Go definitions entirely. The generated asset is
   the runtime distribution form, not another editable source. This is the
   leading candidate, but the compiler stage, bootstrap recipe and source
   availability policy need a precise boundary.
2. **Move/factor the existing readable frontend below Horos.** Could permit
   runtime source loading without a cycle, but changes established ownership and
   stage structure. We see no demonstrated need for this just to ship a prelude.
3. **Application supplies all prelude dependencies explicitly.** Could make even
   library availability ordinary and remove prelude-only source exceptions, but
   canonical-only Horos users still need a standard formation dependency and a
   defined convenience composition. Is this desirable now or separate scope?
4. **Private prelude parser, hand-coded second definition, or semantic bypass.**
   Not acceptable: duplicate grammar/meaning, authored fallback or privileged
   admission merely hides the problem. The owner requires clean replacements.

For candidate 1 the desired flow is:

```text
authoritative horos.kinds.horos
  -> existing Bittype Horos parser
  -> existing syntax-to-canonical module conversion
  -> canonical Ontos module asset shipped with Horos
  -> ordinary Ontos decoder
  -> ordinary Horos module/kind admission
```

Using an already available Bittype tool to prepare an asset is a build-tool
dependency, not a runtime Horos-to-Bittype import. Avoid requiring a package to
compile itself through its newly generated prelude before it can run. Preserve
ordinary declaration identity and limits; do not add a schema registry, new
identity mechanism, fallback decoder or Data-specific membership rule.

## Prior consultation and outside references

[0005-general-kinds.advice.md](0005-general-kinds.advice.md) recommended an ordinary
prelude and distinguished primitive shape vocabulary, kind definitions and
formation contracts. Its three separately pinned public artifacts were adapted
to internal responsibilities of the existing revision; no new signature or
profile identity system was added. It did not settle source authoring or the
specific Bittype-to-Horos compilation dependency, so this question is narrower
than rerunning the general-kind design.

Go's [embed package](https://pkg.go.dev/embed) can place generated asset bytes in
the runtime package without runtime filesystem access; embed patterns are scoped
to that Go module. Rust's [core prelude](https://doc.rust-lang.org/core/prelude/index.html)
is a useful example of a language/library prelude whose availability policy is
distinct from its library contents. These are packaging comparisons, not a claim
that their semantics match Horos or that Horos needs implicit source imports.

## Questions for the expert

1. What should be intrinsic in this architecture, and what belongs in ordinary
   Horos source? In particular, how should primitive formation refer to the
   standard Data definition while Scalar/List/Record/Data remain ordinary kinds?
2. Recommend a concrete single-source compilation and distribution arrangement
   that uses the existing parser/encoder, preserves Horos's one-way runtime
   dependency, and avoids resolving the new prelude through an old one. Include
   the minimal public boundary and a coherent bootstrap/update sequence.
3. How should ordinary dependency loading and source locations work for the
   prelude? Distinguish the minimum change needed now from any justified general
   compiled-module support, and address the current exact-prelude exceptions.
4. What alternative framing, overlooked failure mode, or simplification would
   you recommend? Please challenge the leading option where appropriate and
   keep the result proportionate to this small kind-only standard module.
