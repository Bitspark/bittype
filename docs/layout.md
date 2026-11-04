# Source layout and multi-unit pipelines

**4 November 2026. Proposed Go contracts; documentation only.** The owner asked
for a source-file model alongside syntax and semantics, a concrete independent
`layout.Provider`, composition with multi-unit AST transformations, and an
explanation of parsing a source base through per-text parsers. None of the new
types, methods, registrations or package paths below is implemented by this
change. The implemented baseline is the [single-unit pagination mapping](semantic-generation.md).

## Ownership and current gaps

Semantic identity, target-language syntax, physical file placement and filesystem
IO are separate responsibilities:

| Responsibility | Owner |
| --- | --- |
| Declaration identity, generic binding and type meaning | Language resolver/analyzer; Horos owns canonical rules |
| Packages, modules, imports, qualified references and visibility expressed in source | Language AST and language-owned processing stages |
| Which syntax tree is written to which relative path; grouping and links between source units | The proposed layout model and language-specific layout provider |
| Enumerating and reading files, selecting a build's source set and dependencies, writing outputs | Consuming application |

File placement can affect language meaning. Moving a declaration across a package
or module boundary can change lookup, visibility and import cycles. A layout
provider must either construct the appropriate native syntax under its declared
profile or reject the move. The generic model does not define a universal
package, namespace or import system.

The current transformer returns one `representation.AST`. The initial native
profile bundles the reachable Horos module closure into one source unit. Go has
a configured package name; C++ has a configured namespace. `ast.TreeID` and
`ast.SourceID` are opaque addresses, not paths. Parsers already take text rather
than filenames. Generic semantic capture already accepts several same-language
AST envelopes, but the current Horos resolver profile requires one Program tree.
Native semantic resolvers/analyzers are not implemented. These facts are the
starting point, not claims that a multi-file pipeline already works.

This proposal uses a layout stage after semantic-to-native mapping and before
formatting. Mapping produces logical target syntax with explicit declaration
references. Layout chooses file boundaries and finalizes imports and qualified
references. This refines the earlier informal phrase "plan before imports are
finalized": it does not require a new language-neutral intermediate type algebra
or a second implementation of native type mapping.

## 1. Shared layout model

The proposed `layout/go` package uses existing AST addresses and
`semantics.Key` as an opaque owner/name vocabulary key. It does not store terms,
validity judgments, generic substitutions or Horos constructors. All identifiers
below are local to one captured layout and its associated AST collection.

```go
package layout // proposed: github.com/Bitspark/bittype/layout/go

import (
    "github.com/Bitspark/bittype/semantics/go"
    "github.com/Bitspark/bittype/trees/ast/go"
)

type UnitID string
type GroupID string
type RelativePath string

// One independently parsed/printed text resource in the initial profile.
type File struct {
    Unit   UnitID
    Path   RelativePath
    Source ast.SourceID
    Tree   ast.TreeID
}

// Provider-owned organizational grouping. Kind does not invoke a provider.
// Sites identify syntax anchors when the organization is expressed in source.
type Group struct {
    ID      GroupID
    Kind    semantics.Key
    Members []UnitID
    Sites   []ast.NodeRef
}

// A provider-owned relation between supplied units, such as an import.
// At identifies the syntax that expresses it, when present.
type Relation struct {
    Kind semantics.Key
    From UnitID
    To   UnitID
    At   *ast.NodeRef
}

// Mutable construction records; not a validated or immutable model.
type Data struct {
    Files     []File
    Groups    []Group
    Relations []Relation
}

// New validates, bounds and copies Data; accessors return independent copies.
type Model struct { state *modelState }

type Limits struct {
    MaxFiles       int
    MaxGroups      int
    MaxRelations   int
    MaxLinks       int
    MaxPathBytes   int
    MaxTotalBytes  int
}

func New(Data, Limits) (Model, error)
func (Model) Data() Data
func (Model) File(UnitID) (File, bool)
func (Model) FileForTree(ast.TreeID) (File, bool)
```

These are signature sketches, not a compilable implementation. `modelState` is
private implementation storage. A later implementation must define positive
default limits and enforce them before copying or indexing caller data.

The initial physical-file profile has one AST tree per file. A Horos Program
tree can still contain several module declarations; a file is not a canonical
module. Several files can belong to one language package. The generic model
does not infer either relationship from a directory or a filename.

Declaration placement is derived: a current `ast.NodeRef` identifies its tree,
and `FileForTree` identifies its file. There is no duplicate mutable
declaration-to-file table. When a layout operation moves declarations, its new
trees contain their new current addresses, while `Origin.Nodes` retains the old
addresses. Live references must target the new addresses; historical origins
cannot substitute for current references. One source declaration may contribute
to several output nodes without conflating their identities.

`Groups` are organization assertions, not a second set of semantic module
entities. For example, a Go-specific profile can record a package grouping and
point to its package clauses. It checks their agreement. The generic carrier
only checks identifiers, membership, sites and budgets. It does not interpret
the key, invent a package name or claim the grouping is semantically valid.
Likewise, unit relations describe source organization; semantic dependency
edges remain in the semantic model. External package artifacts stay in explicit
language configuration/dependency inputs, not fictitious local files. A local
relation's endpoints must both be present.

Paths use relative slash-separated components, with no empty, `.` or `..`
components, absolute/drive/UNC forms, backslashes or NULs. File and unit IDs,
tree IDs, source IDs and exact paths are unique in this initial profile. A file
cannot also be an ancestor directory of another file. Group IDs are unique;
members and sites must exist and must not repeat within a group. Relation sites
must belong to their `From` unit. Duplicate relation records are rejected;
cycles are permitted structurally and interpreted by the owning provider.

`New` can check record-internal consistency. Binding the model to a collection
also checks that every listed tree and node exists and that every physical tree
has exactly one file. The application chooses the actual destination root and
checks host path equivalence, reserved names, symlink containment and collisions
with existing content before writing. Target deployment/host path rules are
explicit policy, never inferred from the workstation running the planner.

## 2. Collections and immutable context

Keep `AST[L]` unchanged. Add one construction carrier beside the existing opaque
tree envelopes in `trees/representation/go`:

```go
package representation

type ASTSet struct {
    Entry ast.NodeRef
    Trees []AST
}
```

`Entry` is explicit, exists in the set and has the role required by the selected
language profile. It is not implicitly the first file. The provider determines
whether it selects a root module, a package anchor or another supported entry.
All trees in one set have distinct tree IDs and the same registered language,
syntax profile and actual Go marker. Empty, mixed or malformed sets fail. The
initial API remains a single-language operation; the application partitions a
mixed-language source base into explicit language/build contexts.

`ASTSet` is caller-owned construction storage. Its tree envelopes already own
their storage, but its slice is still mutable. Stage adapters capture that slice
and enforce collection limits and each tree's syntax admission. No new generic
"project" AST node or fake directory grammar is added to every language.

Extend existing semantic capture input with optional layout context:

```go
package semantic

type Input struct {
    Entry  ast.NodeRef
    Trees  []representation.AST
    Layout *layout.Model // proposed addition; copied handle to immutable data
    Limits Limits
}
```

The captured layout is available to a resolver/analyzer through a new
`Request.Layout() (layout.Model, bool)` accessor. It is request context, not
part of `semantics.Model`. Capture validates its association with the captured
trees; a provider that needs physical layout requires it in its profile.
Pure logical AST transformations and current callers can omit it.

The snapshot owns the association among entry, trees and layout. Changing a path,
grouping, source text, selected dependencies or entry requires a fresh capture.
An equal-looking layout or repeated tree IDs cannot reuse evidence from another
capture. Additional language build/module configuration remains explicit and
is included in the provider's existing configuration-matching requirements.
No resolver consults an ambient directory or package environment.

File `Source` identifies the text resource associated with that unit. Node
origins can refer to older source resources and trees after generation. They are
not required to equal the new output file's Source ID; fabricating output byte
spans before printing is forbidden. Parsed CSTs and source texts remain in the
application's source inventory for exact reprinting and source diagnostics.

## 3. Independent layout.Provider

The proposed provider operates on same-language ASTs. It returns a complete,
in-memory output plan: finalized native ASTs plus their file layout. It does not
write files. `Plan` can split/merge units, add native imports/module declarations,
adjust qualification and remap current references. It does not choose a new
data representation, expand generic contracts or change canonical type meaning.

```go
package layout // proposed: github.com/Bitspark/bittype/contracts/layout/go

import (
    "context"

    "github.com/Bitspark/bittype/contracts/language/go"
    "github.com/Bitspark/bittype/contracts/semantic/go"
    model "github.com/Bitspark/bittype/layout/go"
    "github.com/Bitspark/bittype/trees/ast/go"
    "github.com/Bitspark/bittype/trees/representation/go"
)

type Provider[L any] interface {
    language.Provider
    Plan(context.Context, Request[L]) (Result, error)
}

type Contract struct {
    Profile  string
    Evidence []semantic.EvidenceRequirement
}

type Input struct {
    Source   semantic.Snapshot
    Evidence []semantic.Result
    Options  Options
}

type Options struct {
    Configuration []semantic.Setting
    Limits        Limits
}

type Limits struct {
    MaxTrees      int
    MaxNodes      int
    MaxLinks      int
    MaxBytes      int
    MaxDepth      int
    MaxWork       int
    Layout        model.Limits
}

type Request[L any] struct { state *requestState[L] }

func (Request[L]) Entry() ast.NodeRef
func (Request[L]) Inputs() []ast.AST[L]
func (Request[L]) InputLayout() (model.Model, bool)
func (Request[L]) Evidence(int) (semantic.Result, bool)
func (Request[L]) Configuration() []semantic.Setting
func (Request[L]) Limits() Limits
func (Request[L]) Publish(Publication[L]) (Result, error)

type Publication[L any] struct {
    Entry  ast.NodeRef
    Trees  []ast.AST[L]
    Layout model.Data
}

type Result struct { state *resultState }

func (Result) Syntax() representation.ASTSet
func (Result) Layout() model.Model
```

The sketch follows existing capture/publication rather than introducing a second
semantic fact system. Results belong to their exact invocation; a provider cannot
return a result retained from another call. Publication checks/copies the complete
tree collection and layout, remapped references, entry, limits and language
syntax. Accessors return independent slices or immutable envelopes. Errors and
cancellation yield no usable partial plan. Language rules and target-specific
relocation validity remain the provider's responsibility; publication is not
a universal proof of semantic equivalence.

The typed registration adapter would be `layout.NewPlanner[L](definition,
contract, provider)`. It validates ID, profile and Go marker and uses the existing
`semantic.AdmitEvidence` boundary. `Engine.RegisterLayout(planner)` and
`Engine.PlanLayout(ctx, input)` would register and invoke this capability
independently. Registration installs no parser, resolver, analyzer or transformer.
Calls execute outside registry locks. Providers/configuration are explicitly
selected by the application; the engine has no native language list or default
placement policy.

Two different admission profiles must be distinguished. A generated-syntax
profile can accept an exact checked transformer output whose declaration links
are already explicit, validating all links it uses and treating supported native
intrinsics/external references according to its documented configuration. It
rejects unresolved or contradictory references. Reorganizing arbitrary parsed
source can require native resolution/analysis evidence. It cannot pretend that
source spellings alone establish bindings. Native analyzers do not exist today,
so the latter is a later capability. Horos analysis cannot be reused as native
analysis after a transformation.

The profile's closed set of admitted constructs must specify relocation rules
for visibility, imports, target dependencies, cycles and declarations that must
stay together. Unknown rules or missing dependency context cause rejection.
Native helpers and imports introduced by representation mapping participate in
planning; planning from canonical dependency edges alone would miss them.

## 4. Multi-unit AST transformations

Replace the singular return type at the existing transformation boundary when
this design is implemented. Retain one ASTTransformer abstraction and route
registry; do not introduce a competing legacy/single-file transformer API.

```go
package transform

type Output struct {
    Syntax representation.ASTSet
    Layout *layout.Model // optional; nil means logical units without placement
}

type ASTTransformer interface {
    FromID() language.ID
    ToID() language.ID
    Transform(context.Context, Input) (Output, error)
}
```

`Input` still contains the existing captured `Source` and explicit `Evidence`.
The typed mapping adapter publishes/checks every target tree, the explicit entry
and any supplied layout. A mapping that emits one tree returns a one-element set.
No compatibility alias or old parallel execution route remains after migration.

`Engine.TransformTo` returns `transform.Output`. Existing path selection remains
about registered language edges. Each subsequent edge receives a fresh capture
of the complete checked output set, its entry and any output layout, never the
previous edge's semantic evidence. An edge requiring new semantic analysis still
requires a separate explicit resolve/analyze call. The engine never picks one
tree and discards the others, executes each file as an independent transformation,
or guesses an entry.

Logical outputs may have complete type syntax and resolved declaration links
before physical imports are added. Their layout profile must support that form.
A mapping which already emits file-relative imports must supply matching output
placement; it cannot produce an unlocated file-relative import and claim a
finished physical program. Target-aware relocation handles existing supplied
placement, or rejects it when it cannot preserve meaning.

The public composition would be:

```go
// Proposed API flow; handle every error before continuing.
mapped, err := engine.TransformTo(ctx, canonicalInput, targetID)

target, err := engine.Capture(ctx, semantic.Input{
    Entry:  mapped.Syntax.Entry,
    Trees:  mapped.Syntax.Trees,
    Layout: mapped.Layout,
})

planned, err := engine.PlanLayout(ctx, layoutcontract.Input{
    Source: target,
    Options: layoutcontract.Options{Configuration: placementConfiguration},
})

// For each planned unit: look up its AST, format it, print its CST,
// and associate the resulting bytes with that unit's RelativePath.
// Only the application commits the complete file set to its chosen root.
```

Here `layoutcontract` denotes `contracts/layout/go`; `layout` in the earlier
transform.Output sketch denotes the shared model package. The initial generated
layout profile requires explicit target links, not a native analysis result.
Profiles that require semantics add their compatible results to `Input.Evidence`.

The layout stage is itself an AST-to-AST operation within one language, with
extra physical organization output. It gets its own capability because callers
can arrange already generated or, with appropriate evidence, parsed source
without changing language. It is not a same-language edge in route search and
cannot be used as a hidden retry after a failed transformation.

After layout, all native syntax is final. Formatters only construct CSTs; printers
only emit their leaves. They never choose paths, split files, resolve imports or
move declarations. A whole-output buffer/staging area lets the application
finish formatting all units before changing any destination. Filesystem commit,
overwrite policy and rollback belong to that application, not the planner.

## 5. Pagination in two layouts

Start with canonical modules containing `Page<T>`, `Person` and
`PeoplePage = Page<Person>`. Mapping retains those declaration identities in
origins and makes current target references explicit. It can initially bundle
all three declarations as it does today. The target layout provider receives that
syntax plus an explicit placement configuration and produces, for example:

| Target | Unit path | Contents and organization |
| --- | --- | --- |
| TypeScript | `pagination.ts` | exported Page |
| TypeScript | `people.ts` | Person, PeoplePage and a type import of Page under the selected module-resolution policy |
| Go | `pagination.go` | Page in package example |
| Go | `people.go` | Person and PeoplePage in the same package example; no import between these two files |
| Rust | `src/pagination.rs` | Page with the visibility selected by the crate policy |
| Rust | `src/people.rs` | Person, PeoplePage and explicit access to Page |
| Rust | `src/lib.rs` | module declarations connecting the generated units |
| C++ | `include/example/pagination.hpp` | Page and its required includes/namespace declarations |
| C++ | `include/example/people.hpp` | Person, PeoplePage and the selected include route to Page |

These are proposed qualification cases, not current output options. Language
AST vocabulary/grammar additions needed for module scaffolding or header
protection must be explicit and qualified; no raw-text escape node is allowed.
Build/package manifests and compiler preprocessing are outside this source-only
layout profile. Locked external artifacts and module-resolution choices are
supplied by consuming build configuration.

Moving PeoplePage must not duplicate Page's generic body, specialize it silently,
merge unrelated same-spelled binders or change Data obligations. Layout acts on
the selected native representation; it does not strengthen the native value
conformance guarantees documented by the mapping.

## 6. Parsing a whole source base

The parser remains exactly a text parser:

```go
type Provider[L any] interface {
    language.Provider
    Parse(parse.Input) (cst.CST[L], error)
}
```

The application supplies a bounded inventory. An illustrative application-owned
record is:

```go
type SourceFile struct {
    Unit     layout.UnitID
    Path     layout.RelativePath
    Language representation.Language
    Input    parse.Input // explicit SourceID, TreeID and Text
}
```

This is not a proposal for another parser interface. The inventory records the
text bytes actually read, with unique source/tree/unit IDs and relative paths.
The physical root, file selection and retrieval policy stay outside the engine.
An editor buffer, an in-memory collection, an archive or a native cube can supply
the same records without pretending to be a filesystem.

The flow is:

```text
application-selected files and explicit build/dependency configuration
    -> read/capture SourceFile records
    -> Parse(file.Input) once per file, using its registered language
    -> retain each CST, its AST view and source/path association
    -> group the supplied trees into explicit language/build contexts
    -> Capture(entry, all context trees, input layout)
    -> Resolve across that context
    -> Analyze with that resolution result
```

The application can parallelize parsing with bounded concurrency when providers
are concurrency-safe. Deterministic IDs and inventory ordering do not depend on
completion order. Parsing one file requires no other file's semantic result, so
cross-file/cyclic references do not force parser order. Parse failure is reported
at that file; a project-generation operation produces no usable complete result
after any required file fails. An editor can retain individual successful CSTs,
but must label the project incomplete and cannot use them as complete analysis.

Do not concatenate all files into one source string. That loses per-file grammar
boundaries, source spans, module context and independent lossless reprinting.
Do not turn an import into an instruction for a parser or resolver to open an
ambient file or fetch a package. Missing dependencies remain missing until the
application supplies them under its explicit retrieval policy and captures a new
context. A successful per-file parse says nothing about cross-file validity.

Native source selection can depend on language/build configuration. The
application selects and records that configuration and the resulting input set;
the appropriate resolver interprets it. Parsing an arbitrary existing repository
is not currently supported: the parsers admit the documented type-declaration
subsets, and preprocessing/executable syntax may require future profiles.
Unsupported files cannot be silently skipped to claim a complete source base.

### Horos's concrete assembly gap

Today each parsed file gives a Program AST. A Program can already contain
multiple module sections, but the resolver accepts only one Program tree. The
proposed multi-file resolver profile would decode each supplied Program into an
explicit module pool, choose the entry module from the explicit Entry, and call
Horos's authoritative resolution with that pool. It preserves the original
tree-qualified sites; it does not concatenate source or fabricate an aggregate
tree that loses provenance.

The new profile must define duplicates deterministically: reject duplicate
supplied canonical module pins and contradictory pin/name assertions, rather
than letting file iteration overwrite a module. Different pinned versions may
coexist where Horos permits them; imports select their exact pins. File paths
never replace Horos canonical pins or module names. Two files cannot implicitly
contribute fragments to one canonical module; that would need a separately
specified authoring rule. Unused supplied modules remain outside the entry's
reachable semantic closure, while supplied assertions still receive the
validation required by the Horos profile.

This extension lives in the Horos resolver/private bridge and its documented
profile, using the existing collection in semantic.Input. No general
assembly-provider SPI is needed for this concrete case. It is future work in
this proposal. Native resolvers will similarly interpret their explicit context
under their own rules when implemented; a batch of parsed native files is not
already a resolved native project.

## 7. Package layout and dependency direction

Proposed additions and changes, with no empty packages created now:

```text
layout/go/                         new: physical model, copies and structural checks
contracts/layout/go/               new: independent Provider, request/result, adapter
layout/<language>/go/              new: language-owned arrangement and import rules
trees/representation/go/           extend: ASTSet construction carrier
contracts/semantic/go/             extend: captured optional layout context
contracts/transform/go/            change: multi-unit Output, typed output admission
engine/go/                        extend: independent layout registration/execution;
                                            complete collection route handling
resolve/horos/go/                  extend: multi-Program input profile
internal/horosmodel/               extend: pool assembly with current per-tree sites
integration/go/                   extend: cross-package qualification
docs/layout.md                    this proposal
```

`layout/go` depends on AST addresses and the existing owner/name key type, not on
semantic requests, an engine, native syntax or IO. `representation.ASTSet` knows
neither layout nor semantics. Shared semantic capture may depend on the layout
model; the model never depends on capture. The layout contract depends on the
model, representation, semantic evidence/capture and identity. Concrete layout
providers depend on their own target syntax and contract. The engine depends
only on those agnostic contracts. `provider/<language>/go` may compose the
optional stage; explicit registration remains mandatory.

The new capability is optional and has a tested capability profile per language.
The existing provider matrix continues to describe implemented source; this
proposal does not mark any layout or project-parsing capability implemented.
No central switch chooses Go packages, Rust modules, TypeScript imports, C++
headers, or Function/Stream/Cell behavior. Libraries continue to own contracts.

## 8. Implementation sequence and acceptance

Only this document and navigation/status corrections are authorized in the
current change. A subsequent implementation would proceed in bounded stages:

1. Add the shared model/ASTSet and exact-context capture, plus foreign-language
   fixtures proving no core language list. Qualify ownership, duplicate IDs,
   paths, missing sites, resource limits and stale/mutated context rejection.
2. Replace the single-output transformation boundary and migrate every mapping,
   engine route, consumer and test together. A one-unit result remains ordinary
   collection data, not an old API alias. Qualify complete-set capture, explicit
   entries, cross-tree links and evidence isolation between route legs.
3. Add the layout contract and independent engine registration. Qualify marker/
   profile/evidence checks, callback locking, cancellation, deterministic output,
   remapping, immutable publication and empty results on failure.
4. Implement explicitly bounded generated-layout profiles and pagination fixtures
   for each target. Check imports, target-introduced dependencies, visibility,
   helpers, collisions, invalid cycles and origins. Compile the whole generated
   file set and an independent consuming package through maintained remote gates.
5. Extend Horos resolution to supplied per-file Programs; qualify imported,
   nested and Data-bound pagination across files, missing/conflicting/duplicate
   modules, unchanged pin identity after physical relocation, and diagnostic
   paths/spans. Native semantic project support remains a separate increment.

Source round trips preserve each file's original text through its CST. Formatted
output must parse per unit and reconstruct the intended context; text need not
recover historical origins or external semantic IDs. Whole-project compilation,
cross-file binding checks and round trips are separate evidence. No passing
compiler run claims runtime value conformance, adapter behavior or deployment.
