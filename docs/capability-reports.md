# Detailed capability reports

The [ACX-007 identity correction](https://github.com/Bitspark/platform/blob/main/docs/concepts/architecture-cases/0007-bittype-matrix-contract-identity.md)
records the existing matrix bug found during implementation: short display names
could satisfy an unrelated constructor dependency. The initial report implementation
landed in [PR126](https://github.com/Bitspark/bittype/pull/126), but independent review
found that even package-qualified type strings lose private member ownership in
unnamed arguments. Availability, contribution deduplication and case-context
association now use `go/types.Identical`. Row grouping substitutes only the selected
language marker with a shared type parameter before that comparison, preserving
the other Go identities, private members and structural properties. Display labels
and anchors never establish type equivalence; ambiguous rows or anchors diagnose.

The [generated matrices](providers.md) link each discovered contribution to its
[detailed report](provider-capabilities.md). Source annotations, format-owned
extractors and explicit test observations explain the extent of an implementation
without inventing completeness percentages. The committed report is source-only;
execution observations belong to reports for specific test invocations.

This design incorporates the [evaluated expert recommendations](../research-docs/0012-capability-reports.md#applied).
The [tool guide](../tools/provider-matrix/README.md) covers commands and public
composition. Existing `bittype:dimension` and `bittype:shared` annotations retain
their separate attachment and interpretation rules.

## Purpose and meaning

The current [matrix tool](../tools/provider-matrix/README.md) discovers provider
implementations and typed transformer endpoints. Existence is distinct from
feature breadth, exercised examples and architectural correctness.

A detailed report keeps these facts separate:

| Information | Owner and permitted interpretation |
| --- | --- |
| Contribution identity | Go discovery identifies the language, instantiated contract, operation or direct transformer factory |
| Declared scope and exclusions | Authored provider/profile documentation; defines the intended boundary and known limitations |
| Extracted definitions | Owning implementation's structured declarations; establishes what is written in those declarations |
| Associated cases | Explicit links from tests to operations and scenarios; establishes that a case has been identified |
| Executed results | Existing test runner output for those cases and their execution context; establishes the observed result of the named tests |
| Applicability | Host-owned association determines whether an observation matches the supplied source declarations and execution context |

A rule-table entry does not prove that the parser, formatter or analyzer handles
it. A passing case does not prove all combinations of a feature. Neither a
passing matrix nor a report establishes semantic ownership, additive composition
or architectural conformance. Those require reasoning about responsibilities
and actual change propagation under the owner's
[architectural-correctness rule](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md).
An architectural violation remains a defect even when its examples pass; do not
convert it to a percentage or merely describe it as missing feature coverage.

Scope is relative to a language revision, provider profile and supplied
configuration. Types-only tooling is not measured against every executable
construct in a complete language. No overall percentage, estimate of remaining
effort or fixed feature denominator is used. An omitted test is unassessed within
the declared scope, not necessarily unsupported. Known exclusions remain
explicit authored facts. An unavailable extractor, an unknown structural field,
a skipped test and an unobserved test are distinct facts.

## Ownership and extension boundary

The existing source discovery continues to own language/contract association,
operation identity and direct transformer endpoints. Annotations cannot create a
matrix cell, register a runtime provider or override those discoveries.

The public host in `tools/provider-matrix/report/go` owns annotation syntax,
resolved symbol association, bounded inputs, result validation, source links,
diagnostics and deterministic rendering. Independently supplied extractors own
interpretation of particular structured definition formats. The command composes
`report.Binding` values pairing format identities with extractors selected by
exact definition symbol or named-type origin. True Go aliases retain their
original identity; a newly defined lookalike type does not match. The host has no
language or feature catalogue and does not import the concrete grammar extractor.

The extraction contract must be narrow enough to preserve these guarantees:

- A request contains only the selected definition symbol, a captured source view
  and limits. Provider, profile and test destinations remain host-owned; the
  extractor cannot redirect the result to another contribution.
- The fixed result schema contains structural features, ordered roles,
  cardinalities, allowed kinds, attributes, reference-role allowance and located
  extraction limitations. It cannot create providers, claim test outcomes or
  certify semantic validity.
- Feature identity is `(definition root, format identity, feature key)`. Equal
  spellings from different roots or formats do not merge. Source locations must
  belong to the supplied capture and are validated before host association.
- Selection is explicit and unique. Missing extractors produce an unavailable
  inventory with a reason; overlapping selections are configuration errors.
  Neither invokes an ordered fallback or guesses from package/feature names.
- Loading, source queries, traversal and accepted results are bounded. Unknown
  syntax may yield an explicitly partial inventory; exhausting a resource limit
  fails the operation instead of publishing a complete-looking subset.

The public contract is definition-centered:

```go
type Extractor interface {
    Selectors() []Selector
    Extract(context.Context, DefinitionRequest) (Inventory, error)
}

type DefinitionRequest struct {
    Root   Symbol
    Source SourceView
    Limits Limits
}
```

`SourceView` returns immutable captured source strings and copied records for
declarations, exact expression spans, canonical symbols, named-type origins and
compile-time constants. Extractors may parse private copies. No mutable host AST
or `go/types` graph is exposed. Locations use physical byte offsets; `//line`
directives do not redirect links. Foreign files and spans are rejected.

Each invocation selects one explicit source module in the existing Bittype
repository layout. Dependency export data supplies type identity, not dependency
source, annotations or documentation. Active `TestGoFiles` and `XTestGoFiles`
have a separate declaration index and cannot introduce production providers.
The public runner does not imply arbitrary multi-module or repository-layout
support. There is no ambient documentation search. The Go command follows the
caller's module and toolchain configuration and may download missing dependencies;
their export data does not expand the selected source/document boundary.

This is a development-tool contract, not a new provider SPI in `contracts/`, a
runtime registry, an AST node annotation scheme or a Logos proof contract. An
extractor is trusted tooling code supplied by the composition. Loading bounds
metadata, file count and captured source bytes; extraction bounds queries,
traversal, records and accepted output. These are cooperative guarantees, not a
hostile-code sandbox: arbitrary Go extractor code and an active `go/types` call
cannot be forcibly preempted by an interface. The caller's context bounds
cooperative operations; existing CI jobs also impose an outer job deadline. The tool never runs
inspected provider constructors, initializers or predicate functions.

For a new structured definition format, its owner supplies an extractor through
this tooling contract and the consumer selects it. Shared discovery, renderer
and unrelated extractors do not change. Another language reusing the existing
grammar format needs only its own definitions and annotations. A genuinely new
kind of report fact may require a deliberate contract extension; arbitrary JSON,
callbacks or renderer snippets are not an escape hatch.

## Annotation format

Three directives in Go declaration doc comments associate existing source:

```text
// bittype:grammar <symbol-ref> [profile=<value>]
// bittype:scope <document-ref> [profile=<value>]
// bittype:case target=<symbol-ref> feature=<key> case=<key> expects=<outcome> [profile=<value>] [language=<symbol-ref>]
```

There is one directive per line. A token without whitespace may be unquoted;
values containing whitespace use double-quoted Go string-literal syntax. No
shell expansion, interpolation or execution occurs. Named arguments can be
ordered freely, occur once and must be known to that directive. Required values
must be nonempty. The first slice recognizes standalone `//` doc-comment lines;
block-comment directives, trailing statement comments and Markdown/code examples
are not annotation sites. A malformed intended directive is a located diagnostic,
not silently ignored metadata.

`profile` is an optional exact profile label, using an existing profile's value.
It does not invent, install or verify that profile. Omission means the annotation
has not specified a profile; it never means every possible caller configuration.
Different explicit profiles remain distinct, and an unqualified record is not
silently promoted into a qualified one. The annotated test must actually select
the stated profile; this is an assertion about the test, not a way to configure it.

### Symbol and feature references

A symbol reference resolves to an actual Go declaration through the loaded
package and type information:

- `rules` names a declaration in the comment's package.
- `parse.Parse` uses that file's existing import alias.
- `Provider.Check` names a method in the current package;
  `parse.Provider.Parse` uses an imported receiver type.
- `github.com/Bitspark/bittype/parse/syntax/typescript/go#Parse` is the fully
  qualified form. A method uses `import/path#Receiver.Method`.

The qualified form needs no otherwise-unused Go import. Pointer receivers and
generic declarations resolve to their declared symbols; actual language/type
instantiations remain those discovered for the matrix. Unresolved or ambiguous
references are errors. The extractor does not search for a similarly named
function or infer an operation from a path spelling. Sources outside the
explicitly supplied package/dependency set are not fetched implicitly.

Feature and case keys are nonempty UTF-8 values, scoped to their owner. Examples
use lowercase words and hyphens by convention; extraction preserves existing
keys rather than requiring owners to rename them. Quote values when necessary
and escape them when rendering. For the grammar extractor, the feature key is the
existing node-kind value, such as `mapped-type`; it is not duplicated in a separate
language catalogue. Cases refine that feature, for example `key-remapping` or
`nested-conditionals`. Features need not be AST nodes: semantic or transformation
tests may use owner-defined keys such as `generic-binding` under their operation.
Such a key has no structural-definition evidence until an extractor supplies it.
Equal keys in different languages or provider contexts imply no shared semantics.

### `grammar`: an implementation-owned inventory

Attach this directive to the operation or its explicitly discovered delegate:

```go
// bittype:grammar rules
// bittype:scope docs/syntax.md#typescript
func Check(view ast.View[s.Language]) error {
    // Existing implementation uses rules.
}
```

The symbol must identify a structured definition handled by the selected
extractor. The directive claims that this definition belongs to the operation;
it is not proof of execution reachability. The tool records that association,
checks that the operation belongs to a discovered provider/delegation chain and
does not perform speculative whole-program call analysis.

Multiple roots may be declared when the operation composes several definition
sets. Their root/format namespaces stay separate, including equally spelled
feature keys. Duplicate root/profile associations or duplicate feature keys
inside one extracted root are errors. There is no last-writer-wins rule.

### `scope`: existing authored documentation

Attach a scope link to an operation/delegate or a typed transformer factory:

```go
// bittype:scope docs/native-semantics.md#first-typescript-profile
```

Document references are relative to the annotated source module's repository
root, with an optional Markdown fragment. They must resolve inside that supplied
root to an existing document/anchor. Absolute paths, path escapes and symlink
escapes are errors. Only the explicitly selected source module supplies
documents; dependency export data does not enable dependency documentation.
The composition supplies the publication URL base; annotations do not choose an
arbitrary publishing destination.

Scope links are additive. The report links to the authored sections; it does not
parse prose into automatic promises, duplicate exclusions into a manually
maintained status catalogue or treat a linked proposal as delivered support.
No scope annotation means scope documentation is unassociated, not unrestricted.
Identity documentation belongs on `ID`/`Name` methods, not a provider type
declaration.

### `case`: a named test association

Attach a case directive to the actual Go test function:

```go
// bittype:case target=parse.Parse feature=mapped-type case=key-remapping expects=accept
func TestMappedKeyRemapping(t *testing.T) {
    // Assert acceptance and the expected key-remapping structure.
}
```

`expects` is a nonempty owner-defined label, such as `accept`,
`located-rejection`, `round-trip-preserved` or `no-partial-model`. It describes
the outcome asserted by the case, not the result of executing the test. The
host does not interpret those labels as a closed support taxonomy. Rejection
cases must assert
the relevant diagnostic/category in their test code; merely observing some
failure does not establish the desired boundary. A successful rejection test is
shown as an expected-rejection case, never as successful support for the input.
Acceptance/rejection follows the owning operation's contract. For semantic
providers, a rejected validity report and an operation error remain different;
the reporter never derives this distinction merely from a Go error value.

The first slice binds each annotation to its package and enclosing named test
function. Multiple directives may identify the separate associations a test
actually exercises. Duplicate association keys on that test are errors. A
table-driven function may describe an explicitly named family of cases, but the
report does not invent individual subtest results or multiply a parent pass into
independent qualifications. Precise dynamic subtest metadata is deferred.

`target` must resolve to a discovered operation/delegate or a typed transformer
factory. Test packages are read only for associations; their implementations do
not become production matrix cells. Generic targets that cannot be associated
with one language require `language` naming an actual discovered Go marker. A
language annotation must agree with the target's typed instantiation. Transformer
directions come from their factories; no source/target names are inferred from
test names. Ambiguous endpoint configurations remain unresolved instead of
being credited to every edge.

Case links are exact. Testing parsing does not qualify checking, formatting,
semantic analysis or transformation. An end-to-end test may explicitly associate
several operations, but it must exercise and assert the claimed behavior of each.
The tool cannot determine the adequacy of those assertions by scanning the test.

Static wrapper/delegate relationships support navigation only. A delegate's case
or observation does not assert that a wrapper executed, preserved its arguments
or selected the same configuration. Each annotation keeps its exact target.
Generic-composer tests do not qualify native transformer factories they never
call. Profile labels remain authored assertions; Go test events cannot observe
which profile a test selected internally.

## Extraction rules

### Existing discovery stays authoritative

Continue discovering language markers, instantiated role contracts and transformer
endpoints through Go types. Existing dimension annotations still classify the
rows. Keep separate identities for language, dimension, instantiated contract,
implementation/delegate and explicit profile. A displayed cell may aggregate
several implementations; its detailed report must keep them distinguishable.

Annotations associate records with these discovered identities. A missing
implementation cannot become implemented through a scope link or case annotation.
The shared CST printer may have carrier-level cases without a language grammar
inventory. This is a legitimate report, not a reason to synthesize grammar facts.

### Initial extractor: the existing grammar format

The initial extractor owns the meaning of
[`internal/grammar.Rules`](../internal/grammar/rules.go). It recognizes its exact
Go type and helper declarations, not any map named `rules`. Within the supplied
source view it may read constant expressions, identifiers, composite literals
and bounded acyclic references to definition slices/maps. It recognizes the
owning format's `One`, `Opt`, `Many`, `Attributes`, `Attr`, `Set` and corresponding
literal forms by resolved symbol identity. It does not evaluate arbitrary Go.
The inventory describes declared source definitions; it does not establish that
mutable tables retain their initial values during execution or that every
declared rule is reachable. The report must not strengthen that observation.

For each declared node kind it emits:

- Ordered child roles, minimum/maximum cardinality and declared allowed kinds.
- Attribute names, required/optional flags and explicit finite allowed values
  when a `Set` or equivalent declarative form supplies them.
- Declared reference allowance and links to the defining source.
- Opaque predicate links where an attribute uses a custom function.

Knowledge is recorded for each field, rather than assigning one support flag to
the whole feature:

| State | Meaning |
| --- | --- |
| Known | The declaration supplies the value, including known empty, false or zero |
| Partial | A collection has a known subset and an explicit unresolved portion |
| Unknown | The field could not be decoded; no invented empty value replaces it |
| Opaque | An attribute predicate is identified but its behavior is not interpreted |

Unbounded cardinality and unrestricted attribute values have explicit flags;
they are not unknown values. Every uncertain field has a reason and physical
source location. An undecodable entry key remains an unresolved portion of the
inventory. Its decoded/partial extent concerns the selected declaration and the
fixed structural schema, not feature or language completeness.

The owning helper syntax is interpreted only in that extractor. The reporting
host does not know TypeScript, `mapped-type` or the names of grammar helpers.
An unsupported initializer, cyclic definition or exhausted traversal budget
does not result in an invented empty set or a complete inventory. Report the
specific extraction limitation; independently extracted entries, if displayed,
are explicitly partial.

For example, the current TypeScript mapped-type rule can yield this inventory:

```text
mapped-type
  children, in order:
    type-parameter   exactly one type-parameter
    name-type        zero or one declared type expression
    type             exactly one declared type expression
  attributes:
    readonly-modifier   optional; add | remove
    optional-modifier   optional; add | remove
```

The inventory is structural. Procedural restrictions in `Check`, such as `infer`
placement or prohibiting defaults on mapped parameters, do not automatically
follow from it. Their documentation and named cases remain separate. Do not
scan switches and infer support from a case label, count lines of code or use
word heuristics to turn custom checks into semantic guarantees.

### Test results remain a separate input

The committed report is source-only: it records associations even when no tests
have run, with no mutable latest-pass labels, timestamps or inspected HEAD value.
Existing remote jobs publish invocation-specific observations as artifacts or
job summaries. A documentation-only commit does not relabel an earlier test
result as evidence for the new revision.

The observation package imports an explicit invocation capture with Go JSON
stdout and stderr. It preserves exact emitted test names, package events, build
failures, repeated occurrences, selection arguments, process exit and capture
completeness. Events are observations first, including failed or unassociated
executions. Package success, output text or cached-looking messages cannot
supply a missing named execution. An early test pass does not erase a later
package, build or invocation failure.

Attribution requires an independently supplied matching source context:
revision, module and harness identities, selected modules/replacements and Go
build selection. Unknown, modified, incomplete or unequal inputs prevent
attributed success and retain their reasons. Preparation changes are disclosed,
not hidden behind a recorded HEAD. An independently installed consumer harness
retains its own context instead of donating results to repository tests.

The join uses exact tested-package/top-level-test-name pairs. Internal and
external test declarations sharing a runner name remain ambiguous. Child events
retain their emitted names without becoming authored cases; parent outcomes do
not become independent subtest qualifications. Repeated invocations and repeated
occurrences remain separate.

Invocation arguments and allowlisted environment describe the supplied execution
configuration. Go test JSON does not establish which provider profile or mapping
a test selected internally; profile annotations remain authored assertions.
Contracts and native-compiler jobs keep their separate observations and context.

Keep case execution states separate from expected operation outcomes: passed,
failed, skipped, incomplete and not observed. A failed test may represent an
implementation or test defect; it does not automatically mean a valid feature is unsupported.
Skipped and unobserved cases are never counted as passed. Missing results leave
the source-derived report usable with associations only. Stale or mismatched
supplied results produce a visible diagnostic and no attributed success.

The documentation generator does not run qualification tests or native-language
compiler checks, or download arbitrary results. Its Go loader may prepare
dependency export data. This adds no
architectural score, administrative gate or release step.

For captured results, generate a separate execution report with:

```sh
go run ./tools/provider-matrix/go -results invocation.json -context source-context.json -execution-output report.md
```

The command does not launch tests, fetch a latest successful run or create a
result ledger. The existing test workflow owns execution and its exit status.

## Reports and links from both matrices

The generated `docs/provider-capabilities.md` contains provider and transformer
details. The generated [language overview](language-support.md) is the compact
support-depth entry point; [the matrices](providers.md) show availability. Source links stay;
generation adds a `details` link beside each implemented operation/factory:

```text
Provider cell:     Check · details
Transformer cell: New · details
Unavailable cell: /
```

Each link resolves to the report for that specific implementation, not just its
language. Multiple implementations in a cell each retain their own links. Section
anchors encode the full discovered identity tuple injectively. Attached definition
and scope anchors also encode whether a profile is absent or explicitly named.
They do not depend on line
numbers, feature counts, timestamps, hashes or inspection revisions. A
symbol/profile change regenerates both the anchor and incoming links together.
Generated links and sections must agree before publication; a malformed annotation must not publish
a half-updated matrix. All three documents are fully rendered and validated before
writing. File writes are not a filesystem transaction: a write failure stops
publication, and the existing updater commits the successful three documents together.
Existing authored prose remains outside generated regions.

### Language overview matrices

`docs/language-support.md` groups the same captured reporting facts by discovered
language. It shows a dimension × provider-contract matrix for each language and a
separate source-language × target-language matrix of direct transformers, alongside
owner-written features/exclusions and expandable inventories and exact test sources.
The README and both matrices link to it. No language, stage or feature catalogue
is added to the renderer. Each selected definition is extracted once, then shared
by the detailed and language views, including independently supplied formats.

**Owner direction, 6 November 2026:** preserve conceptual matrix axes. Do not
replace stage dimensions or transformation endpoints with lists in a combined
cell, or suppress empty rows/columns to make the overview smaller. The transformer
matrix and each provider matrix remain visible; expandable details hold their
supporting evidence.

Provider rows are discovered dimension labels. Columns are the actual originating
Go provider declarations, compared by Go type identity. Instantiations of one
provider in different dimensions share a column, while each cell retains every
exact instantiated operation and its evidence link. Same-spelled unrelated
contracts receive qualified column labels; they are never merged. All discovered
axes appear for every language, including completely empty rows or columns.
Language identity appears as metadata above its matrix, separate from pipeline
operations. There is no fixed list of dimensions, roles or languages.

Transformer rows and columns both contain every discovered language. Only
explicitly discovered direct edges populate cells. Reverse, composed and diagonal
identity routes are not inferred. Both documentation views use the same direct-edge
renderer. `/` means no discovered implementation for that pairing, not a claim
that the pairing is inherently impossible. Populated cells link separately to each
implementation and its exact capability details, including scope/profile/test
associations. Language-wide test totals remain in the evidence section; they do
not serve as per-cell qualification.

Existing `bittype:scope` documents own the prose selected for the overview:

```markdown
## Type fragment

<!-- bittype:overview scope=type-fragment part=features -->
Aliases and generic applications have structured syntax tooling.
<!-- bittype:overview:end -->

<!-- bittype:overview scope=type-fragment part=exclusions -->
Executable bodies remain outside this profile.
<!-- bittype:overview:end -->
```

The marker's `scope` is an existing heading or explicit anchor in that document;
`-` means the document itself. Selection requires the exact fragment of the
operation's `bittype:scope` annotation. A file-wide scope does not select its
children. Profile absence remains distinct from a named profile. Several passages
can select the same scope, in document order. The owning doc remains the source
of truth; these passages are authored claims, not inferred or verified support.
They cannot add a provider, transformer edge or test association.

Passages contain paragraphs with emphasis, single-backtick code spans and simple
inline Markdown links, at most 4 KiB each. Ordinary
inline links are rebased to the owning source document. Block markup, reference
links, images and raw HTML are outside this small excerpt format. Fenced/indented
examples outside excerpts are ignored. Invalid markers, missing anchors, empty
or oversized excerpts stop generation before publication. Scope files retain the
existing module containment and 2 MiB read bounds. An absent feature/exclusion
passage is reported as unreported, never interpreted as complete support or no
exclusions. No expectation label or feature spelling is classified by keywords.

Scope passages explicitly associated with multiple languages appear once in a
shared expandable section. Scopes without excerpts remain listed inside the
evidence section; they are not silently treated as fully supported.

Repeated wrappers do not inflate case or test counts. Structural inventories may
be folded by exact definition/format/profile for readability; all operation and
target associations remain linked separately. Counts describe associated test
declarations and authored cases, not execution, percentage coverage or separate
execution of every case annotation. A delegate's evidence remains its own.

The committed overview has no execution outcomes. The existing CI capture steps
also emit `overview.md` beside each detailed invocation `report.md`. Outcomes
remain a language × outcome matrix: pass, fail, skip and incomplete occupy
separate columns; any additional supplied outcome retains its own column. `/`
means no matching observation, while zero means matching observations exist but
none has that outcome. This uses
the same exact source/build/harness association as the detailed report, counts
only matched test occurrences, retains repeated pass/fail outcomes and the
process exit, and credits no language for a mismatched invocation. Full
diagnostics and unassociated observations remain in that invocation's detailed
report and raw capture. Counts across languages are not additive: one declared
test can explicitly cover several languages.

A provider report contains:

1. Language, dimension, role, implementation link and explicit/unspecified profile.
2. Links to the owning scope and exclusions.
3. Extracted inventory, its completeness/limitations as an extraction, and links
   to custom predicates. Extraction completeness is not language completeness.
4. Exact named case associations and their authored expected outcomes. Execution
   reports add invocation observations separately. Definitions without associated
   cases are visibly unassessed.
5. Missing associations or unavailable inventory/results, without treating absence
   as either unsupported behavior or successful coverage.

An illustrative feature section would look like this; it is not a current result:

| Feature | Extracted definition | Associated case | Expected outcome | Execution |
| --- | --- | --- | --- | --- |
| `mapped-type` | Owning rule link | `key-remapping` | Accept | Not observed |
| `mapped-type` | Owning rule link | `extra-members` | Reject | Not observed |

A transformer report additionally identifies the exact source/target endpoints,
factory and selected profile/composition where supplied. It links to the owning
mapping scope and cases for that direction. Target AST vocabulary does not imply
source representability. A factory's availability does not establish every
possible configured mapping, native constraint preservation or runtime behavior.
Library-owned mappings keep their own scope/cases; the generic report host does
not learn their type names or semantics. Reverse and composed paths receive no
credit from a forward edge's cases.

The generated file is a view of definitions and associations, not a new source
of provider support flags. Exclusions remain in the owning scope documents. No
manual edits to cells, feature rows or observed-result labels are required.

## Packages and composition

The tooling lives in the root Go module:

```text
tools/provider-matrix/
  go/                    command and explicit extractor/result composition
  report/go/             public discovery, association, rendering and write host
  capability/go/         public definition/source/structural extraction contract
  extract/grammar/go/    owning Rules-format extractor
  observation/go/        invocation decoding and exact-context association
docs/
  providers.md           compact generated matrices within authored prose
  provider-capabilities.md  generated static details
  language-support.md       generated compact language overview
  capability-reports.md  this authored design
```

Independent extractors live in their owners' tooling packages. Runtime and
language-provider packages never import the reporting packages. The command
selects concrete extractors; the report host depends only on their contract and
the observation model. A materially different struct-field extractor fixture
uses the same public runner without changes to the host or renderer. A second
language using the same grammar table alone would not establish that boundary.

`report.Generate` returns fully rendered documents without writing them.
`report.Write` writes those three documents and reports write failures; it does not commit or
publish. External tooling consumers can compose another extractor over the
runner's explicitly selected module snapshot. See the
[public composition example](../tools/provider-matrix/README.md#public-composition).

The annotation rollout covers native and Horos grammar roots, operation scope
links and bounded cases for syntax, formatting, pagination and native semantic
subsets. It is not an exhaustive case catalogue. Transformer factories have scope
links; generic composer tests are not attributed to factory wrappers they do not
execute. Source observations and mechanical tests never establish whole-language
coverage or replace architectural reasoning.

Stronger process isolation, explicit dependency-source selection and precise
dynamic subtest metadata are deferred until a supported consumer needs them.
They are not implicit guarantees of the public host. Missing evidence remains
visible without a new reporting service, ledger or approval workflow.
