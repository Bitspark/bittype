# Detailed capability reports

**Proposal, 6 November 2026. No implementation in this change.** The owner
requested a way to explain the extent of each implemented capability behind the
[generated matrices](providers.md), without putting a coverage essay in each
cell or inventing completeness percentages. This document specifies proposed
source annotations, extraction rules and report navigation. The annotations,
extractors, report files and CI result integration described below do not yet
exist. Existing `bittype:dimension` and `bittype:shared` annotations are unchanged.

## Purpose and meaning

The current [matrix tool](../tools/provider-matrix/README.md) discovers provider
implementations and typed transformer endpoints. Existence is distinct from
feature breadth, exercised examples and architectural correctness.

A detailed report combines four separately labelled sources:

| Information | Owner and permitted interpretation |
| --- | --- |
| Declared scope and exclusions | Authored provider/profile documentation; defines the intended boundary and known limitations |
| Extracted definitions | Owning implementation's structured declarations; establishes what is written in those declarations |
| Associated cases | Explicit links from tests to operations and scenarios; establishes that a case has been identified |
| Executed results | Existing test runner output for those cases and their execution context; establishes the observed result of the named tests |

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
effort or fixed feature denominator is proposed. An omitted test is unassessed,
not necessarily unsupported. Known exclusions remain explicit authored facts.

## Ownership and extension boundary

The existing source discovery continues to own language/contract association,
operation identity and direct transformer endpoints. Annotations cannot create a
matrix cell, register a runtime provider or override those discoveries.

The reporting mechanism owns annotation syntax, resolved symbol association,
bounded inputs, result structure, source links, diagnostics and deterministic
rendering. Independently supplied extractors own interpretation of particular
structured definition formats. The command's composition selects extractors by
the resolved Go declaration/type identities they handle. It contains no language
or feature catalogue. The initial grammar extractor uses exactly the same
contract as an external extractor.

The extraction contract must be narrow enough to preserve these guarantees:

- A request identifies the annotated operation, the selected definition symbol,
  its package and type, the applicable profile context, and a bounded read-only
  source/type view. It does not grant arbitrary repository or environment access.
- A result is an inventory associated with that request. It contains feature
  entries, structural child/attribute constraints, links to opaque custom checks
  and explicit unsupported extraction reasons. It cannot create unrelated
  providers, claim a test passed or certify semantic validity.
- Feature keys are scoped to the owning definition/provider context. Source
  locations must belong to the supplied source view. Host association checks
  prevent an extractor from attaching another operation's definitions or tests.
- Selection is explicit and unique. Missing extractors produce an unavailable
  inventory with a reason; overlapping selections are configuration errors.
  Neither invokes an ordered fallback or guesses from package/feature names.
- Cancellation, bounded traversal and failure produce explicit outcomes. No
  partial extraction is labelled a complete inventory.

This is a development-tool contract, not a new provider SPI in `contracts/`, a
runtime registry, an AST node annotation scheme or a Logos proof contract. An
extractor is trusted tooling code supplied by the composition; a read-only API
does not claim to sandbox arbitrary Go code. The source being inspected is data:
the tool never runs provider constructors, initializers or predicate functions.

For a new structured definition format, its owner supplies an extractor through
this tooling contract and the consumer selects it. Shared discovery, renderer
and unrelated extractors do not change. Another language reusing the existing
grammar format needs only its own definitions and annotations. A genuinely new
kind of report fact may require a deliberate contract extension; arbitrary JSON,
callbacks or renderer snippets are not an escape hatch.

## Annotation format

The first proposal adds three directives in Go declaration doc comments. The
following is annotation grammar, not implemented Go tooling:

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
sets. Their origins stay separate. Duplicate roots or conflicting definitions
of a feature within the same effective operation/profile are errors unless the
owning format defines an explicit composition that its extractor represents.
There is no last-writer-wins rule.

### `scope`: existing authored documentation

Attach a scope link to an operation/delegate or a typed transformer factory:

```go
// bittype:scope docs/native-semantics.md#first-typescript-profile profile="captured-types-1"
```

Document references are relative to the annotated source module's repository
root, with an optional Markdown fragment. They must resolve inside that supplied
root to an existing document/anchor. Absolute paths and path escapes are errors.
An independently supplied library resolves its documentation against its own
supplied source root, not Bittype's root. The composition supplies the publication
URL base; source annotations do not choose an arbitrary publishing destination.

Scope links are additive. The report links to the authored sections; it does not
parse prose into automatic promises, duplicate exclusions into a manually
maintained status catalogue or treat a linked proposal as delivered support.
No scope annotation means scope documentation is unassociated, not unrestricted.

### `case`: a named test association

Attach a case directive to the actual Go test function:

```go
// bittype:case target=parse.Parse feature=mapped-type case=key-remapping expects=accept
func TestMappedKeyRemapping(t *testing.T) {
    // Assert acceptance and the expected key-remapping structure.
}
```

`expects` is `accept` or `reject`: it describes the operation outcome asserted
by the case, not the result of executing the test. Rejection cases must assert
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

Source extraction discovers associations even when no tests have run. Execution
status can be attached only from explicitly supplied existing test-runner output
for the inspected source and relevant configuration. A global green job does
not establish that a named case ran. A result for different source, an unknown
profile or another compiler configuration is not silently reused as a current
qualification. No extra checksum, history store or provenance workflow is needed;
use the context already available to the existing CI run, or omit the result.

Keep case execution states separate from expected operation outcomes: passed,
failed, skipped and not observed. A failed test may represent an implementation
or test defect; it does not automatically mean a valid feature is unsupported.
Skipped and unobserved cases are never counted as passed. Missing results leave
the source-derived report usable with associations only. Stale or mismatched
supplied results produce a visible diagnostic and no attributed success.

The documentation generator does not launch compilers/tests or download arbitrary
results. Existing remote test execution remains responsible for qualification.
This proposal adds no architectural score, administrative gate or release step.

## Reports and links from both matrices

Proposed output is one generated `docs/provider-capabilities.md` file containing
provider sections followed by transformer sections. It is not created by this
documentation change. The existing `docs/providers.md` remains the compact entry
point. The current source links stay; generation would add a `details` link beside
each implemented operation/factory, for example:

```text
Provider cell:     Check · details
Transformer cell: New · details
Unavailable cell: /
```

Each link resolves to the report for that specific implementation, not just its
language. Multiple implementations in a cell each retain their own links. Section
anchors derive deterministically from the full discovered identity tuple
(including profile when present), using escaped identity text and deterministic
collision disambiguation, not line numbers or checksums. A symbol/profile change
regenerates both the anchor and incoming links together. Generated links and
sections must agree before publication; a malformed annotation must not publish
a half-updated matrix. Existing authored prose remains outside generated regions.

A provider report contains:

1. Language, dimension, role, implementation link and explicit/unspecified profile.
2. Links to the owning scope and exclusions.
3. Extracted inventory, its completeness/limitations as an extraction, and links
   to custom predicates. Extraction completeness is not language completeness.
4. Named case associations with their expected operation outcomes and any supplied
   execution results. Definitions without associated cases are visibly unassessed.
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

## Placement and implementation sequence

The proposed tooling layout is:

```text
tools/provider-matrix/
  go/                    existing command; discovery and explicit composition
  capability/go/         proposed tooling-only records and extraction contract
  extract/grammar/go/    proposed extractor for the current grammar format
docs/
  providers.md           existing generated matrices with proposed detail links
  provider-capabilities.md  proposed generated reports
  capability-reports.md  this authored design
```

Independent extractors may live in their owners' tooling packages. Runtime
packages never import the reporting packages. The command may select concrete
extractors; shared report orchestration depends only on their contract. Exact Go
API details belong to implementation review, but the ownership, association and
failure guarantees above are requirements from the first slice.

Implementation can proceed in bounded steps:

1. Add annotation parsing/resolution and associated source records. Keep existing
   matrix availability discovery unchanged. Test malformed and ambiguous metadata.
2. Add the grammar extractor through the tooling contract and a materially
   different independently supplied extractor fixture, for example declarative
   Go struct definitions instead of the current helper/map grammar encoding.
   Compose both through the same contract, including invalid associations, with
   shared orchestration fixed. Another language using the same rule table alone
   would not demonstrate this extension boundary.
3. Generate inventory/case sections and links from both matrices. Initially show
   cases as defined, with execution not observed. Missing extractors remain explicit.
4. Optionally consume named results already produced by existing remote tests.
   Extend the existing documentation update path when this is implemented; no
   separate reporting service or status ledger is proposed.

Relevant future tests cover symbol/profile isolation, custom predicates, partial
extraction, bounds, duplicate definitions, accept/reject distinctions, missing and
skipped results, generic-provider association, transformer directions, deterministic
anchors and unchanged authored prose. They qualify these mechanics; they do not
replace architectural reasoning or establish whole-language coverage.

Until implementation lands, the matrices retain their existing source-only links
and meaning. This proposal changes no executable annotations, provider methods,
grammars, generation mappings, semantic models, workflows or runtime behavior.
