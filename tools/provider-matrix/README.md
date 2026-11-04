# Provider matrix generation

`go run ./tools/provider-matrix/go` updates only the `providers` and `transformers`
generated regions in [docs/providers.md](../../docs/providers.md). The tool is
development tooling in the root Go module; runtime packages do not import it.
It has no third-party dependencies and executes no provider constructors or methods.

## Source discovery

The Go command supplies package files and dependency export data for the active
Go build context. The tool reads repository source with `go/parser` and `go/types`.
Test packages and infrastructure under contracts, trees, engine, internal,
integration and tools are not implementation candidates. The contract packages
are still read to discover the interfaces being implemented.

- Columns come from concrete types under `syntax/` implementing `language.Provider`.
  ID and Name come from their constant-return methods; no language catalogue exists.
- Rows come from interfaces named `Provider` under `contracts/` that embed language
  identity. Go method sets and inferred generic arguments establish the implemented
  interface. Carrier types distinguish Syntax and Files, independent of package aliases.
- Concrete providers are associated through their language identity. Reusable
  providers are instantiated for each language marker that satisfies their Go constraint.
  Public constructor signatures supply their required provider collaborators.
  Availability is the least fixed point of those requirements; circular dependencies
  cannot create availability. Infrastructure adapters never supply implementations.
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

Implementation links use Go token positions, ignoring `//line` remapping. Simple
single-return delegation is followed to its entry point; dynamic dispatch is not.
Distinct implementations remain separate links in a cell. Shared implementations
repeat the same link across languages. Direct transformer edges are separate from
engine routes and profile-specific applicability.

Output is deterministic, with languages sorted by ID. There are no dates, revision
records or manually maintained line numbers. Links target the module's GitHub
`main` branch; `-source-url` overrides the prefix. `-root` selects the repository
and `-output` the Markdown file. Authored prose and newline style are preserved.

## Automatic updates

[Provider matrices](../../.github/workflows/provider-matrices.yml) runs on GitHub
after relevant changes land on `main`, or by workflow dispatch. It reads latest
main, generates the tables, and, when they change, creates or updates the bot-owned
`automation/provider-matrices` branch and pull request, then squash-merges that
docs-only update. No manual table edit or merge is needed. The job is serialized
and uses `GITHUB_TOKEN`, with write credentials available only during publication.

The bot commits only `docs/providers.md`, using `[skip ci]`. Token-authored changes
do not recursively start another update. Unchanged output creates no commit. This
is output generation, not a validation gate; existing checks remain separate.
Portable generation runs on the shared build service or GitHub runner.

Source changes can briefly precede their generated documentation while the job
runs. Publication failures remain workflow failures, never a claim that the
repository was updated merely because a local table was produced.
