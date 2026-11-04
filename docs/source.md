# Syntax and Files

**4 November 2026. Owner-selected symmetry; proposed Go APIs only.** This replaces
the separate file-layout proposal. No carrier, provider, registration or generator
implementation is added by this documentation change. The previous proposal is
removed; there is no second active design or compatibility path.

## 1. One structure at two levels

```text
Syntax : CST    : AST
Files  : Source : Semantics
```

| | Syntax | Files |
| --- | --- | --- |
| External form | Text | File set |
| Concrete representation | CST | Source |
| Abstract representation | AST | Semantics |
| Parse | Text -> CST | Files -> Source |
| Format | AST -> CST | Semantics -> Source |
| Print | CST -> Text | Source -> Files |
| Abstract view | CST exposes AST | Source exposes Semantics |

Use the names **Source** and **Semantics**, without adding a Model suffix on only
one side. Both rows describe representations. Syntax and Files are processing
dimensions, not language IDs. There is no separate layout carrier, provider,
registration or pipeline stage. File placement, grouping, module organization and
concrete import routes belong to Source; their construction belongs to Files
formatting. Files printing emits the already chosen files.

Files means an in-memory collection of paths and text. Reading directories,
fetching dependencies, selecting a build's inputs and writing destinations remain
application IO. An archive, editor buffer or native cube can supply the same
collection. No provider implicitly reads a path or fetches an import.

## 2. Revised Go carriers

These sketches show the intended public shapes. Supporting private storage,
admission, limits and typed adapters are specified below, not implemented here.
The implemented [semantic graph](semantics.md) is currently called
`semantics.Model`. The proposed replacement name is `semantics.Semantics[L]`.
That is a clean future rename and typed language association of the existing fact
graph, not a new type algebra or an alias retained at the old name. Current API
references elsewhere continue to describe implemented source until that migration.

The existing Syntax carriers stay:

```go
// trees/ast/go
type AST[L any] struct {
    ID    TreeID
    Root  NodeID
    Nodes []Node[L]
}

// trees/cst/go — existing carrier, fields abbreviated only for readability.
type CST[L any] struct {
    Abstract ast.AST[L]
    Root     NodeID
    Nodes    []Node[L]
}

func (CST[L]) ASTView() ast.AST[L]
```

Files is the external representation of the other row:

```go
// Proposed files/go
type ID string
type Path string

type File struct {
    ID   ID
    Path Path
    Text string
}

type Files struct {
    Entries []File
}
```

File IDs are local inventory identities, not paths or semantic declaration IDs.
Paths are relative slash-separated names. IDs and paths are unique; no empty,
absolute, drive/UNC, dot, dot-dot, backslash or NUL path is admitted. A file cannot
also be an ancestor directory of another file. The application selects the root
and enforces host-specific path equivalence, reserved names, existing-content and
symlink rules before IO. File selection/order and IDs must be deterministic under
the selected policy, not dependent on parallel parser completion order.

Semantics retains the existing provider-owned graph and queries:

```go
// Proposed revision of semantics/go; reuses Entity, Term, Data and their rules.
type Semantics[L any] struct { state *state }

type View[L any] interface {
    SemanticsView() Semantics[L]
}

func (Semantics[L]) SemanticsView() Semantics[L]
func (Semantics[L]) Data() Data
func (Semantics[L]) Entity(EntityID) (Entity, bool)
func (Semantics[L]) Term(TermID) (Term, bool)
```

`L` associates an admitted graph with a language marker in the same way as the
syntax carriers. The configured publication adapter also checks language ID,
profile and actual Go marker. Generic records remain provider-owned facts: there
is no universal notion of a type, package or behavioral contract added here.
Existing resolution, expression, binding and dependency queries migrate with the
graph. IDs remain local to one graph, and accessors return copies. Constructing a
well-formed graph alone does not prove language validity; accepted analysis
reports remain required. Do not manufacture a marker conversion by casting facts.

Source contains concrete organization and the abstract view it realizes:

```go
// Proposed source/go
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
    Kind semantics.Key
    From files.ID
    To   files.ID
    At   *ast.NodeRef
}

// Sealed by the current Files operation's publisher; zero value is invalid.
type Source[L any] struct { state *state[L] }

// Construction input, never itself an admitted Source.
type Publication[L any] struct {
    Snapshot  semantic.Snapshot
    Analysis  semantic.Result
    Units     []Unit[L]
    Groups    []Group
    Relations []Relation
}

func (Source[L]) Entry() ast.NodeRef
func (Source[L]) Units() []Unit[L]
func (Source[L]) Groups() []Group
func (Source[L]) Relations() []Relation
func (Source[L]) Snapshot() semantic.Snapshot
func (Source[L]) Evidence() semantic.Result
func (Source[L]) SemanticsView() semantics.Semantics[L]
```

The view accessor returns the graph from the retained Analysis result; it performs
no resolution, analysis or formatting. Both Semantics and Source implement
`semantics.View[L]`, just as AST and CST implement `ast.View[L]`. Publication is
untrusted construction storage, like a raw CST record. Source is its sealed
admitted value: the graph, successful receipt and exact capture are retained
together. Unit/CST and organization accessors return deep copies; evidence and
snapshot accessors return immutable handles. SemanticsView alone does not prove
this association.

Publication uses the existing semantic evidence admission against Snapshot and
the configured language/profile/report requirements. Analysis must belong to
that exact snapshot and establish the required graph. Every unit's full AST must
match the captured AST, not just tree/node IDs; its ID/path/SourceID must match
captured Files context, with complete one-to-one unit coverage. Concrete text must
correspond to the captured file contents under the selected Syntax printer.
CST/AST consistency remains a language-owned admission obligation, not a semantic
claim inferred by a generic structural checker. Groups/relations and entry must
agree with this same captured context.

The publisher derives the graph from Analysis instead of accepting an independent
replacement. It atomically bounds, checks and copies construction storage. There
is no public unchecked Source constructor or mutable Abstract field. Changing a
path, AST, CST or grouping requires a new publication with appropriate capture
and evidence; later mutation cannot change the Source. Each Files invocation has
a fresh private publication identity, and its adapter rejects Source values
returned from an earlier invocation.

Each unit is one independently parsed/printed file and owns its CST, whose AST
view supplies the tree ID. There is no duplicate tree-ID field or mutable
declaration-to-file map. Declaration placement follows its current node reference
to the unit containing that tree. A file can contain several Horos modules;
several Go files can belong to one package. The shared carrier infers neither
from filenames. A language profile owns Group/Relation key meanings. Local
relations must identify supplied units; external artifacts remain explicit
dependency context rather than invented local files.

Source must contain an explicit valid Entry, unique file/source/tree IDs, complete
unit ownership and existing group/relation sites. Relation sites belong to their
From unit. Language-specific grouping and import assertions must agree with the
contained syntax. Cycles are allowed structurally and interpreted by the provider.
No core switch interprets these keys or names Function, Stream, Cell or a language.

## 3. Parse, format and print provider interfaces

Use the same three role names at both levels. The proposed contracts are generic
in their input and output, with no combined provider and no separate Files-layout
role. These replace the existing single-language-parameter text interfaces in a
future coordinated migration; this document changes no current Go interface.

```go
// contracts/parse/go
type Provider[I, O any] interface {
    language.Provider
    Parse(context.Context, I) (O, error)
}

// contracts/format/go
type Provider[I, O any] interface {
    language.Provider
    Format(context.Context, I) (O, error)
}

// contracts/print/go
type Provider[I, O any] interface {
    language.Provider
    Print(context.Context, I) (O, error)
}
```

Each provider still has ID and Name. The different type arguments describe the
two levels; they do not register a language or erase its marker. The precise
supported instantiations are:

| Role | Syntax | Files |
| --- | --- | --- |
| parse.Provider | `[parse.Input, cst.CST[L]]` | `[parse.FilesRequest[L], source.Source[L]]` |
| format.Provider | `[ast.View[L], cst.CST[L]]` | `[format.FilesRequest[L], source.Source[L]]` |
| print.Provider | `[cst.CST[L], string]` | `[source.Source[L], files.Files]` |

`parse.Input` remains the existing SourceID/TreeID/Text record. Files-level request
records carry the explicit context needed to establish/use Semantics:

```go
// contracts/parse/go
type FilesInput struct {
    Files         files.Files
    Entry         files.ID
    Configuration []semantic.Setting
    Limits        FilesLimits
}

// contracts/format/go
type FilesInput[L any] struct {
    Abstract      semantics.Semantics[L]
    Snapshot      semantic.Snapshot
    Evidence      semantic.Result
    Configuration []semantic.Setting
    Limits        FilesLimits
}
```

Each stage owns a captured request supplied to its provider:

```go
// contracts/parse/go
type FilesRequest[L any] struct { state *requestState[L] }
func (FilesRequest[L]) Input() FilesInput
func (FilesRequest[L]) Publish(source.Publication[L]) (source.Source[L], error)

// contracts/format/go
type FilesRequest[L any] struct { state *requestState[L] }
func (FilesRequest[L]) Input() FilesInput[L]
func (FilesRequest[L]) Publish(source.Publication[L]) (source.Source[L], error)
```

The application supplies construction Inputs to configured adapters. The adapters
capture them and create FilesRequest; Input accessors return copies. Providers
use their request's Publish method and cannot construct another invocation's
private identity. Parse publication must match captured input Files, entry and
configuration. Format publication retains fresh output capture/evidence and must
establish the configured preservation obligations against its admitted input.
Both callbacks return only a Source issued by their own request. Syntax's existing
snapshot rules continue to apply.

The entry file is explicit. Its language profile defines the syntax entry within
that file, including selecting a Horos module when there are multiple sections.
Requested configuration supplies that selection; no provider chooses the first
alphabetical file as an implicit program root. The formatter's Evidence must
establish the required analysis profile/reports for exactly Snapshot and Abstract, including
its language and graph identity. Equal bytes, entity numbers or names do not
establish that association. The provider declares the requirements it needs;
missing/incompatible/incomplete evidence fails before formatting.

Parse-to-Format therefore carries the whole established result:

```go
parsed, err := filesParser.Parse(ctx, parse.FilesInput{ /* supplied files */ })
// Handle err. filesParser is the adapter that constructs FilesRequest.
formatted, err := filesFormatter.Format(ctx, format.FilesInput[L]{
    Abstract: parsed.SemanticsView(),
    Snapshot: parsed.Snapshot(),
    Evidence: parsed.Evidence(),
    Configuration: organization,
})
// Handle err before Files Print. The adapter supplies a new FilesRequest.
```

No caller must re-analyze just to recover evidence discarded by parsing. The same
composition works with a formatted Source: its accessors return its output
capture/receipt, never those of the preceding Source. A graph-only value does not
satisfy Files formatting's evidence requirements.

`FilesLimits` is a stage-owned record whose implementation must bound files,
total text/path bytes, CST/AST nodes, links, graph records, depth and total work,
plus the already established provider configuration/diagnostic budgets. Limits
are checked before copying/indexing input and before output publication. The
default values and any separately bounded substage budgets must be specified
and tested during implementation. Cancellation or error produces no usable
complete output; diagnostics retain file and syntax locations.

There is no inference of valid registrations from I/O shapes. Typed registration
adapters admit exactly the listed pairings, language/profile and actual L marker.
The engine registers parser, formatter and printer independently for each
dimension, for example RegisterSyntaxParser and RegisterFilesParser, with matching
Format/Print registrations. Registering any one capability installs no others.
There is no automatic fallback from Files to Syntax when semantic work is missing.

The input/output generics do not replace the existing explicit resolution and
analysis interfaces. Files providers are visible compositions of explicitly
supplied lower-stage and semantic capabilities. Applications construct them with
those dependencies; the engine does not discover analyzers or run them during an
unrelated syntax operation. A formatter's use of validation is part of its Files
contract, not a hidden engine side effect. Calls run outside registry locks.

## 4. Files parsing: one source base, many text parses

```text
Files
  -> Syntax Parse for each file
  -> CSTs and their AST views, retaining separate file identities
  -> capture one explicit same-language/build context
  -> Resolve over the complete supplied context
  -> Analyze using that exact resolution result
  -> publish Source, including the established Semantics
```

The application discovers/selects and reads the files before calling Files Parse.
It supplies external dependencies and build/module configuration explicitly.
Native build selection and preprocessing are not invented by the generic kernel.
Mixed-language repositories are partitioned into explicitly selected language/
build contexts; one Source has one language and semantic profile.

The provider invokes the existing Syntax parser once per file. It supplies unique
SourceID and TreeID values and retains the association to file ID/path. No file
text is concatenated with another. Parsing can run with bounded concurrency when
the supplied providers support it; semantic dependencies do not impose parser
order. Unresolved names are permitted in the per-file ASTs until resolution.

File identity/path can affect resolution. Extend semantic capture with immutable
Files context and file-to-tree/source associations, rather than placing paths in
the semantic graph. The proposed request-context records are:

```go
// files/go — input context only, not another abstract/concrete representation.
type SyntaxLink struct {
    File   ID
    Tree   ast.TreeID
    Source ast.SourceID
}

type Context struct {
    Files Files
    Links []SyntaxLink
}

// contracts/semantic/go — proposed extension of existing capture Input.
type Input struct {
    Entry ast.NodeRef
    Trees []representation.AST
    Files *files.Context
    Limits Limits
}
```

Capture checks and owns this correspondence, including the bounded text inventory;
provider metadata/configuration checks cover the selected environment. Changing
files, paths, dependency selection, entry or build configuration requires a new
capture/invocation. Request accessors copy mutable storage. Evidence from another
capture cannot be reused merely because tree IDs match. Files Context is supplied
before analysis and does not claim to expose Semantics; it is not an alternative
Source type or a renamed layout subsystem.

Files Parse succeeds with Source only after the provider's required semantic
obligations complete successfully and every semantic syntax site belongs to the
captured units. The existing distinction between complete validity, rejection,
missing context and exhausted resources is preserved in diagnostics. A nil Go
error from an analyzer alone is insufficient. An editor may retain individual
successful CSTs after a project failure, but they do not constitute an admitted
Source with established Semantics.

The Horos gap remains concrete: its current resolver accepts one Program tree.
The future profile must decode several per-file Programs into an explicit module
pool and preserve every original tree-qualified site. Duplicate canonical pins
and contradicted module-name/pin assertions fail deterministically; distinct
versions may coexist where Horos permits them. Imports select exact pins, not
filesystem paths. Files cannot implicitly contribute fragments to one canonical
module. Unused modules stay outside the entry's reachable semantic closure under
the owning profile. No general assembly-provider interface is introduced.

Native semantic resolvers/analyzers remain unimplemented at this baseline. Their
existing per-file parsers cover type-declaration subsets, not arbitrary source
repositories. Files parsing cannot claim completion by silently skipping
unsupported executable syntax, unavailable preprocessing or missing dependencies.

## 5. Files formatting: realize Semantics as Source

```text
Established Semantics + explicit representation/organization configuration
  -> choose units, paths, package/module organization and concrete import routes
  -> construct each unit's same-language AST with complete current references
  -> Syntax Format each AST into its CST
  -> validate the complete concrete realization and its semantic correspondence
  -> publish Source with its abstract Semantics view
```

This is same-language realization. It does not choose a different language or
turn Horos Semantics into Go Semantics. Language transformation remains an
explicit separate operation. The formatter may preserve open generics and choose
different file arrangements, but it cannot silently change type meaning,
visibility, binding or representation precision. Target-introduced dependencies,
helpers and external artifacts participate in file organization, not just the
canonical input's dependency edges.

The semantic profile must contain enough established information to construct its
supported source forms. Missing information is an explicit unsupported case;
formatting may not guess it from names, reconstruct lost native meaning, read
ambient files, or use raw-source escape payloads. The initial source profiles and
the facts they require must be qualified independently for each language.

The existing graph includes current syntax sites as well as semantic meaning.
Formatting new units therefore requires an explicit mapping to their new sites.
Preserve the meanings/identities under a checked entity/term correspondence;
numeric graph IDs are not persistent identities. The default qualification path
uses the explicitly supplied resolver and analyzer to establish facts for the
new captured units and checks the provider's declared preservation obligations.
That fresh capture and successful result are retained in the newly published
Source, whose SemanticsView exposes the result's graph. A future direct
publication path needs an equally explicit, tested correspondence contract;
it is not an automatic shortcut.
The generic kernel only checks correspondence integrity and provenance, not the
truth of language-specific equivalence claims.

Do not copy an old analysis receipt onto moved syntax. Current references must
point into the new units. No output byte positions are fabricated before the
syntax formatter/printer has established the corresponding text. Files formatting
chooses source organization; Syntax formatting chooses concrete spelling. Both
operate without filesystem IO. This design requires no historical provenance
archive or retention workflow.

For pagination, a TypeScript Source may put Page in pagination.ts and Person /
PeoplePage in people.ts with the configured type import. A Go Source may put the
same declarations in pagination.go and people.go within one package, with no
import between the two files. The semantic application PeoplePage = Page<Person>
must be preserved. These are future qualification cases, not current generator
output options. Rust module scaffolding and C++ header syntax require explicitly
supported syntax/profile extensions before those Files formatters can claim them.

## 6. Files printing and round trips

Files Print receives an admitted Source, invokes the explicitly selected Syntax
printer on each stored CST and returns Files using the already recorded paths.
It chooses no new filenames, imports, groups, bindings or whitespace. It neither
resolves nor analyzes. Filesystem staging, overwrite policy, commit and rollback
belong to the application after the complete file set has been produced.

For admitted source input, Parse then Print retains each original file path and
text, including trivia. Files are matched by ID/path, not parser scheduling order.
Formats may choose a new organization; Format then Parse must recover the intended
Semantics under the explicit correspondence, not necessarily the original file
boundaries or historical node IDs. Source.SemanticsView and CST.ASTView are stored
views, not reparsing operations. Cross-file validity, semantic preservation,
lossless text round trips and native compiler checks remain distinct evidence.

## 7. Multi-unit AST transformations

Retain one ASTTransformer abstraction and the existing explicit route registry.
The proposed generalization is to a collection of target ASTs:

```go
// trees/representation/go
type ASTSet struct {
    Entry ast.NodeRef
    Trees []AST
}

// contracts/transform/go
type ASTTransformer interface {
    FromID() language.ID
    ToID() language.ID
    Transform(context.Context, Input) (representation.ASTSet, error)
}
```

The entry is explicit; the collection has unique tree IDs, one language/profile/
marker, admitted trees and checked internal reference integrity. A single-unit
mapping returns a one-element collection. No parallel old return type, placement
field or compatibility alias remains after that future migration.

An ASTSet is syntax, not Source. A transformer must not fabricate Semantics merely
by packaging target ASTs. To enter the Files pipeline, qualify and explicitly
establish target-language Semantics, then use its Files formatter. Current native
pagination mappings only establish the documented target syntax/representation
guarantees; native semantic providers are a prerequisite for this proposed full
pipeline. The existing explicitly selected Syntax Format/Print operations remain
useful on their own; the Files API never silently downgrades to them.

Each route leg captures the complete target collection afresh, keeps its explicit
entry and discards the preceding leg's semantic receipts. Requirements for new
analysis are satisfied through explicit application composition, not hidden engine
invocation. No file is silently dropped or transformed alone to fake a project.
File-relative imports require explicit matching input context when their meaning
is checked; an unlocated fragment cannot claim a complete Source.

## 8. Packages and migration

The stage-first hierarchy remains. The planned additions/changes are:

```text
files/go/                       new: Files, IDs/paths, captured input associations
source/go/                      new: sealed Source, publication records, evidence/view
semantics/go/                   revise: Semantics[L] and View[L], existing fact graph
trees/representation/go/        extend: ASTSet and typed Semantics admission; no Source import
contracts/{parse,format,print}/go/ revise: the same generic role contracts
contracts/semantic/go/          extend: immutable Files context and typed graph admission
contracts/transform/go/         revise: collection output
engine/go/                      extend: independent Syntax/Files role registrations
parse/<language>/go/            existing Syntax parsing
format/<language>/go/           existing Syntax formatting
print/go/                       existing shared Syntax printing
parse/files/<language>/go/      new: explicit Files parsing composition
format/files/<language>/go/     new: language-owned Files realization
print/files/go/                 new: shared Files emission with selected Syntax printer
provider/<language>/go/         explicit composition of separately registered stages
integration/go/                both-level and cross-package qualification
```

The extra `files` directory is a processing dimension under a stage, not a
language called Files. Syntax implementations keep their established paths;
the generic contracts expose both dimensions with the same role names. Naming
does not install providers. Architecture discovery/capability tests must account
for this dimension explicitly rather than misclassifying its directory as a new
language or allowing arbitrary dependency directions.

Files depends on AST addresses only for the input associations. Semantics remains
provider-owned facts with no Files or concrete-language dependency. Source depends
on Files, CST/AST, Semantics and the shared semantic capture/result contract needed
to retain evidence, not the engine or concrete providers. Shared semantic capture
may depend on Files context, never Source; representation also must not import
Source. Source is its own sealed typed envelope, so no reverse Source-envelope
import is required there. This keeps the dependency graph acyclic.
The role contracts/adapters depend on typed carriers and
the shared evidence boundary, not a concrete language. Concrete Files providers
compose their language's explicitly supplied syntax/resolution/analysis stages.
The engine remains language-agnostic, with no contract-library name switches.

Implementation would migrate the existing semantic graph name and all public
consumers together, introduce the two-level carriers/admission, generalize the
three role interfaces and collection transform boundary, then qualify concrete
Files compositions. Current Model references are not compatibility APIs added by
this proposal; they describe today's unmodified implementation. The superseded
layout proposal's packages are not created, renamed into aliases or registered.

The owner’s later 4 November standing rule adds no backup, checksum, validation,
review or provenance workflow as an implicit task or delivery prerequisite. This
proposal defines the requested parsing, semantic-analysis and representation
behavior; it does not establish an additional delivery-gate checklist.

## 9. Current versus proposed

Implemented today: five Syntax parser/formatter/printer compositions, AST/CST
carriers, the generic semantic fact graph, explicit semantic capture and Horos
resolution/analysis, and single-unit native pagination transformations. The
[provider matrix](providers.md) continues to report those actual capabilities.

Proposed here: Files/Source carriers, typed Semantics naming, Files-level
parse/format/print, the generic two-argument role interfaces, multi-file Horos
resolution and multi-unit transformer output. None is claimed implemented by
this document. Native semantic completeness, runtime adapters, serialization,
behavior and package/deployment work are outside this documentation change.
