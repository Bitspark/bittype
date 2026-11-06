# Provider matrices and capability reports

`go run ./tools/provider-matrix/go` updates the `providers` and `transformers`
generated regions in [docs/providers.md](../../docs/providers.md) and renders
[docs/provider-capabilities.md](../../docs/provider-capabilities.md) and the compact
[language overview](../../docs/language-support.md). The committed reports contains source-derived identity, scope links, structural inventories and
named case associations. It contains no current-pass labels or completeness score.

This is development tooling in the root Go module. Runtime and language-provider
packages do not import it, and it executes no inspected initializers, provider
constructors, methods or predicates. The public host accepts independently
supplied extractors. The command explicitly composes the current grammar-format
extractor; source discovery and dimension metadata remain the authority for
matrix availability.

The [capability report design](../../docs/capability-reports.md) defines `grammar`,
`scope` and `case` annotations and separates declarations, authored expectations,
execution observations and their applicability. Reports do not certify whole
languages, arbitrary configurations or architectural correctness.

## Source discovery

The Go command supplies package files and dependency export data for the active
Go build context. The host selects one explicit module in the existing Bittype
repository layout and reads its production source with `go/parser` and `go/types`.
Dependency export data provides type identity, not ambient dependency source,
annotations or documentation. Cross-module source discovery is not implemented.
Active `TestGoFiles` and `XTestGoFiles` form a separate declaration index for
annotations; their bodies are not interpreted and test-only providers cannot
create matrix cells.
Test packages and infrastructure under contracts, trees, engine, internal,
integration and tools are not implementation candidates. The contract packages
are still read to discover the interfaces being implemented.

- Columns come from concrete types under `syntax/` implementing `language.Provider`.
  ID and Name come from their constant-return methods; no language catalogue exists.
- Rows come from interfaces named `Provider` under `contracts/` that embed language
  identity. Go method sets and inferred generic arguments establish the implemented
  interface. Declaration-owned dimension comments classify their operation boundaries;
  no package or type name determines the dimension.
- Concrete providers are associated through their language identity. Reusable
  providers are instantiated for each language marker that satisfies their Go constraint.
  Public constructor signatures supply their required provider collaborators.
  Availability is the least fixed point of those requirements; circular dependencies
  cannot create availability. Infrastructure adapters never supply implementations.
- Dependency matching uses `go/types.Identical`, including unnamed arguments with
  package-owned private fields or methods. Row grouping replaces only the selected
  language marker with a common Go type parameter and compares the resulting Go
  types. Operation/call/context association also retains Go identity. Type strings
  are display labels, not equality keys; a rendered row or anchor collision fails
  rather than merging unequal contributions.
- Operation input carriers are not required providers. In particular, `Source[L]`
  as an input does not make printing depend on parsing or semantic analysis.
- Generic providers currently have one language parameter and a public constructor
  returning that provider. Other generic construction shapes require discovery
  support, not guessed cells or handwritten capability flags. Configuration remains
  caller-supplied; this documents implementations, not arbitrary factory execution,
  configuration validity, semantic correctness or complete language coverage.
- Transformer factories supply endpoints through actual
  `transform.NewTransformer[From, To]` calls. Both markers must be discovered
  languages. Unrecoverable endpoints stop generation rather than omitting an edge.

## Dimension metadata

The doc comment on a defined boundary type can contain one of these directives:

```go
// bittype:dimension Syntax
// bittype:dimension Files
// bittype:shared
```

These are three alternative examples, not directives to combine on one declaration.
Dimension labels are source-owned identifiers; the generator has no list of their
names. A label starts with a letter and continues with letters, digits, hyphens or
underscores. `Shared`, `Unclassified` and `Mixed` are reserved output states.
Annotations may appear in line or block doc comments. Annotate individual types
inside grouped declarations. Aliases inherit the target's metadata and cannot
declare another dimension.

The foundational AST/View/CST and text-input types declare `Syntax`.
Files/Source/Semantics and semantic requests/results declare `Files`.
`language.Provider` explicitly declares itself shared. These comments add no
runtime methods, imports, registry or registration effects.

For each instantiated provider contract, the generator collects metadata from its
own declaration and the input/output types of its methods, including inherited
operations. It follows aliases, pointer/collection elements, callback signatures,
and type arguments of unannotated generic wrappers. An annotated carrier owns its
boundary, so traversal stops there. Struct fields, underlying storage, helper
methods and generic constraints do not classify an operation: a semantic result
containing AST references is still in its declared dimension.

One distinct dimension produces that label. Multiple dimensions produce a sorted
`Mixed: A + B` label, retaining the conflict. With none, an explicitly shared
contract produces `Shared`; otherwise the result is `Unclassified`. Embedding a
shared identity interface or accepting a shared input does not turn missing
operation metadata into shared classification. Malformed or duplicate directives
stop generation with a source location, before any document is rewritten.

Resolution and analysis consequently belong to `Files` through their semantic
boundary types. Their independent use on captured ASTs requires no physical files
and does not change this architectural grouping.

## Rendering

Implementation links use physical Go token positions, ignoring `//line`
remapping. Statically discovered wrapper/delegate relationships support source
navigation; dynamic dispatch is not followed. Canonical operation and receiver,
contract arguments, marker and endpoints remain distinct before rendering.
A delegate's case or result does not qualify its wrapper. Distinct contributions
remain separate links and detailed sections even when one matrix cell aggregates
them. Shared implementations repeat across applicable language instantiations.
Direct transformer edges remain separate from engine routes and applicability.

Output is deterministic for the selected source/build context, with languages
sorted by ID and rows by dimension label then contract. There is no special
ordering table for known dimensions. Detail anchors encode canonical identity
tuples, including absent versus explicit profiles, without hashes or line-number
identity. Static reports have no timestamps or inspected HEAD records. Links
target the module's GitHub `main` branch; `-source-url` overrides the prefix.
`-root` selects the module root, `-output` the matrix Markdown file, `-details`
the static details path and `-overview` the compact language view. Repeatable
`-build-flag` values select supported Go build options. `-timeout` supplies a
cooperative command deadline. Authored matrix prose and its newline style are preserved.

All three outputs are rendered and validated before writes begin. A malformed
annotation or extraction failure prevents publication. OS writes are not a
transaction: the caller must stop publication on a write failure. The existing
updater publishes the completed three documents together in one Git commit.

The [compact-view contract](../../docs/capability-reports.md#compact-language-overview)
defines exact scope excerpt markers, feature inventories and test deduplication.
Excerpts live in the existing owning scope documents. The renderer has no list
of languages, supported features or exclusion keywords. Missing excerpts remain
explicit, and test associations never become passing results by inference.

## Public composition

| Package | Responsibility |
| --- | --- |
| [report/go](report/go/report.go) | Public source-discovery, association, validation, rendering and write host |
| [capability/go](capability/go/model.go) | Definition selectors, immutable source view and fixed structural inventory schema |
| [extract/grammar/go](extract/grammar/go/extract.go) | Interpretation of the exact owning Rules type and helpers |
| [observation/go](observation/go/types.go) | Invocation-scoped Go JSON import and exact-context test association |
| `go/` | Command selecting the extractors and explicit optional observation inputs |

The host imports no concrete extractor. A consuming tool can choose its own
bindings using the same public runner:

```go
documents, err := report.Generate(ctx, report.Config{
    Root: root,
    Extractors: []report.Binding{
        {Format: grammar.Format, Extractor: grammar.Extractor{}},
    },
})
if err != nil {
    return err
}
return report.Write(documents)
```

Here `report` and `grammar` refer to the public packages in the table. `Generate`
returns rendered documents without mutating files; `Write` reports failures and
does not commit or publish. An independent struct-field extractor fixture uses
this same runner without changing discovery or rendering.

Extractors receive only a definition, source view and limits. The host owns the
format/root/provider/profile association and validates returned physical spans
and structural records. Known, partial, unknown and opaque fields stay distinct;
known unrestricted values and unbounded cardinalities are explicit, not unknown.
No match is unavailable extraction; multiple matches are configuration errors.
An arbitrary payload or renderer callback cannot extend the result schema.

Source loading, metadata, traversal, queries and accepted output are bounded.
These are cooperative guarantees over trusted tooling, not a hostile-extractor
sandbox. In-process Go checking and arbitrary extractor code are not forcibly
preemptible.

## Exact-invocation observations

The default generation path never fetches test results. Existing remote jobs can
supply captures explicitly:

```sh
go run ./tools/provider-matrix/go -results invocation.json -context source-context.json -execution-output report.md -execution-overview overview.md
```

`invocation.json` uses [observation.Invocation](observation/go/types.go) and sits
beside its `stdout.jsonl` and `stderr.log`. It records the actual command,
allowlisted environment, run/job/attempt, process completion and source selection.
`source-context.json` independently describes the inspected source/module/harness
and Go build context. Its build flags replace `-build-flag`; supplying both is an
error. `-results` can be repeated without merging the invocation identities.
The execution report is an invocation artifact; it does not replace the
committed static report or mutate it into a latest-success ledger.

The importer retains tests, packages, build failures, stderr, repeated
occurrences, exact subtest names and incomplete captures. Association requires
matching complete, unmodified source selections and exact top-level runner names.
Internal/external declarations with a duplicate runner name remain ambiguous.
Package passes and parent/subtest outcomes are not substitutes for a missing
named observation. Skip and incomplete are never pass; an early pass does not
erase a later process failure.

Profile labels and `expects` labels remain authored assertions, not facts
inferred from Go JSON or the error value. An absent profile is unspecified, not
a wildcard. Observations from different contracts/native/installed-consumer
invocations retain their own contexts. Missing or mismatched context is visible
without attributing success. The report command does not run tests or weaken
their original exit status.

## Automatic updates

[Provider matrices](../../.github/workflows/provider-matrices.yml) runs on GitHub
after relevant changes land on `main`, or by workflow dispatch. A merged-PR event
also triggers generation when `[skip ci]` suppresses the ordinary push workflow.
Only merged PRs qualify, and the job always reads `main`, never an unmerged head.
Direct pushes with skip instructions require explicit workflow dispatch, following
[GitHub's event behavior](https://docs.github.com/en/actions/how-tos/manage-workflow-runs/skip-workflow-runs).
It reads latest main, generates all three static documents and, when they change,
creates or updates the bot-owned
`automation/provider-matrices` branch and pull request, then squash-merges that
docs-only update. No manual table edit or merge is needed. The job is serialized
and uses `GITHUB_TOKEN`, with write credentials available only during publication.

The bot commits only `docs/providers.md`, `docs/provider-capabilities.md` and `docs/language-support.md`,
using `[skip ci]`. Token-authored changes do not recursively start another update.
Unchanged output creates no commit. This
is output generation, not a validation gate; existing checks remain separate.

Source changes can briefly precede their generated documentation while the job
runs. Publication failures remain workflow failures, never a claim that the
repository was updated merely because a local table was produced.
