# Changelog

## Package boundaries — 11 November 2026

- Add [package boundaries](docs/package-boundaries.md): the difference between parts that are wrongly grouped in one package and contributions that are wrongly coupled, the question that separates them, what to look for, the findings of a scan at `3bf45c3`, the proposed import rules by role and the order of work ([#218](https://github.com/Bitspark/bittype/issues/218)). Nothing is implemented.
- Add [bittype/ACX-009](docs/architecture-cases/0009-package-boundaries.md): the mapping packages, the verification laws and two analysis provider packages group separate contributions, contract identities and horos names are published by provider packages, and the import check cannot see either ([#219](https://github.com/Bitspark/bittype/issues/219), [#220](https://github.com/Bitspark/bittype/issues/220), [#221](https://github.com/Bitspark/bittype/issues/221), [#222](https://github.com/Bitspark/bittype/issues/222)).
- Record the owner's request of 11 November in `AGENTS.md`, and note in `LAYOUT.md`, [generation](docs/generation.md) and [representation families](docs/representation-families.md) what the current check does not establish and that the mapping split precedes the families.

## Order of correction and the profile gate — 11 November 2026

- Record the owner's 11 November direction in [ACX-004](docs/architecture-cases/0004-generation-file-syntax.md): each misplaced placement decision moves to its Files formatter first and coverage follows additively. The Go package clause ([#161](https://github.com/Bitspark/bittype/issues/161)) and Rust visibility ([#159](https://github.com/Bitspark/bittype/issues/159)) can start; C++ ([#160](https://github.com/Bitspark/bittype/issues/160)) is deferred until its semantic provider, with the reason. The Rust external path is recorded as a further pending instance. No correction is implemented.
- Add [bittype/ACX-008](docs/architecture-cases/0008-structural-profile-gate.md): the shared horos mapping policy admits exactly one profile name, and that name is taken from the first fixture ([#216](https://github.com/Bitspark/bittype/issues/216)). Pending.
- Note in [representation families](docs/representation-families.md) the two changes that precede the additive work, and that ontos delivered the expanded-size measure in v0.17.0.

## Chain issues made startable — 10 November 2026

- Record two further pending instances in ACX-004: the Go package clause and the C++ namespace are fixed by target configuration during generation, while the placement proposal takes each declaration's native scope from the consumer's placement ([#161](https://github.com/Bitspark/bittype/issues/161), [#160](https://github.com/Bitspark/bittype/issues/160)). No correction is implemented.
- Point the delivery tables of [representation families](docs/representation-families.md) and [binding generation](docs/binding-generation.md) at the per-family issues split out of #165, the libraries/service and Content tracking issues, and the ordered plan in [#206](https://github.com/Bitspark/bittype/issues/206), which replaces #5.

## Representation families — 9 November 2026

- Record research 0015's accepted design ([representation families](docs/representation-families.md)). Each family is a representation contract, instantiated once per use, whose native type and plan entry are projections of one selected instance. The plan is a graph with a shared envelope and family-owned payloads declared in horos, five version axes and explicit generic requirements; it stays unstable until six vertical cases pass. Correct 0014's union rule: the obstruction is information loss, not overlap. Nothing is implemented yet.

## GEN-01 verified — 9 November 2026

- Record the independent verification of the GEN-01 correction: #192 at
  `6132170`, board record 000005. ACX-006 and GEN-01 read verified. The
  provider-integrity and single-site diagnostic limits remain; #193 tracks
  certifying that a provider's judgments rest on the admitted evidence.

## Binding generation ownership — 9 November 2026

- Record the accepted ownership of binding generation for horos-declared contracts
  (research 0014, [binding generation](docs/binding-generation.md)). bittype stays
  the representation authority and will publish a typed selected representation
  plan. Per-target binding backends compose paired native and adapter contributions
  from it. The charter gains the plan interface only; nothing is implemented yet.

## Generation, semantics and documentation — 8 November 2026

Changes since the previously recorded module-encoding entry:

- Add bounded Go, TypeScript and Rust semantic integrations with explicit compiler
  configuration and capture-bound facts; C++ semantics remain deferred.
- Implement public joint verification association, direct Logos checking and
  library-owned Horos/Go projections/laws, including independent reversed Pair
  composition and checked nested theorem reuse.
- Replace per-pair generation assembly and procedural mapping operations with
  independently registered descriptions, kernel-owned planning, one generic
  convention and separate native AST construction for all four targets.
- Carry generation.Output into TypeScript Files; decide imports/exports after
  placement and check external symbols, public roots and conflicting local names.
- Record open source traversal ownership (GEN-01), native enforcement (GEN-02)
  and Rust/C++ Files ownership (ACX-004).
- Correct GEN-01 (ACX-006, research 0013): the kernel admits the evidence graph and
  walks declarations, references, applications and children itself, with
  run-scoped occurrence handles and binder-keyed actuals. Sources supply
  classification, descriptor selectors and judgments (admitted, unsupported,
  rejected, undetermined) plus a reader for their libraries. Every actual is
  judged, including omitted ones. The Open/Admit/CheckChild/Catalog reader is
  removed without aliases; horos and the independent source, which now publishes
  its own graph, drive the same walk. GEN-02 and ACX-004 remain unfinished.
- Correct the GEN-01 delivery after independent verification rejected it:
  `PlanningInput` takes the source language, input and provider, and `Plan`
  admits the evidence and opens the source itself (`AdmitInput` is no longer
  exported). Handles resolve with `In(Captured)` against the exact capture. A
  source selector outside the capture's source language fails as `ErrScope`. A
  term recorded at several syntax sites is reported at the occurrence's
  declaration. A use whose callee is not generated or whose callee selector is
  foreign is reported as that use. `Plan` checks mapping
  descriptions before it admits evidence, as the registry does. ACX-006 and
  GEN-01 stay open until the correction is verified.
- Reconcile active architecture/API guides, Horos graph examples and dependency
  versions; add a documentation map, runnable Page/Cell fixture walkthrough,
  native placement proposal, exact verification contribution reference and
  current-carrier context table. Repair section links and clarify historical
  consultation authority. Move the maintained matrix identity case into Bittype.


## Horos module encoding — 5 November 2026

- Add the general `horos-module` command: one Horos source module to canonical
  Ontos bytes, sharing the existing parser and module lowering with source pins.
- Add a manual generation-only workflow that builds the encoder and lowers a
  selected repository source on a GitHub runner. No semantic resolution or
  prelude-specific logic is added to the command or agnostic kernel.
- Pin Horos Go 0.10.2, whose complete standard kinds module is authored in Horos
  source and loaded from the generated canonical asset. Keep Ontos at 0.16.0.

## Unified Horos declarations — 4 November 2026

- Replace separate Horos and Horos/Wire frontends with one revision-2 Horos AST
  and provider family. Remove old packages/profiles without compatibility aliases.
- Preserve explicit Type/Data parameter bounds in syntax and provider-owned
  semantic terms. Check spelling structurally; resolve/analysis delegate meaning
  to Horos and retain shared immutable evidence.
- Pin verified public Horos 0.9.0 with Ontos 0.16.0, removing the development
  pseudo-version and retaining the no-replacement installed-consumer gate.
- Migrate canonical fixtures and the independently installed pagination consumer;
  qualify one Page, DataPage aliases, data/mixed specialization and open forwarding.
- Update the provider matrix and ownership/model guides. No generators or runtime
  bindings are introduced; the agnostic kernel and native language stages are unchanged.


## Unreleased

### 4 November 2026

- Remove the Horos-to-HorosWire template-lifting transformer, its tests, exclusive
  indexing helpers and architecture allowances. Author optional pagination directly
  and use the same generic with data and Wire arguments. Preserve closed ground
  embedding and specialization; remove all compatibility paths. Update active
  design documents and align the public Horos dependency to 0.8.0.

- Remove `text.Provider` and the engine's combined provider API without aliases.
  Register parser, checker, formatter and printer adapters independently against
  language identity/profile metadata, in any order. Pack/PackCST validate generic
  carrier structure; explicit checkers preserve syntax admission for Check,
  semantic Capture and every transformation endpoint. Add partial-registration,
  ownership, preflight, concurrency, removed-API and installed-consumer coverage;
  update current documentation and the provider matrix.

- Move parsing, formatting and printing contracts to `parse.Provider`,
  `format.Provider` and `print.Provider`; move source input to `parse.Input`.
  `text.Provider` composes these independent contracts with `check.Provider`.
  Update consumers, independent positive/negative fixtures, dependency checks,
  documentation and the six-language provider matrix. Remove the old names
  without compatibility aliases.

### 3 November 2026

- Move AST syntax admission to `check.Provider` in `contracts/check/go`,
  independent of text and CST contracts. `representation.NewLanguage` accepts
  this interface; `text.Provider` composes it with parsing, formatting and
  printing. Update Horos and native providers, independent consumers, dependency
  checks and the provider matrix. Remove `text.Checker` without an alias.

### 2 November 2026

- Standardize the public provider interfaces as `language.Provider`,
  `resolution.Provider`, `analysis.Provider` and `text.Provider`. Individual
  text capabilities are `text.Parser`, `text.Formatter`, `text.Checker` and
  `text.Printer`. Update all current consumers, fixtures and documentation;
  remove the previous interface names without compatibility aliases.

### 31 October 2026

- Adopt the stage-first hierarchy: trees/{ast,cst,representation},
  contracts/{language,pipeline,text,transform} and uniform native stage paths.
  Move generic AST checks to check/ast. Remove superseded public paths and the
  engine's moved CST exports rather than retaining compatibility aliases.
- Consolidate checked AST/CST snapshots and descriptors in representation;
  preserve sealed payloads and construction-time marker checks across the split.
- Discover the complete native stage matrix in tests, enforce per-package
  production dependencies, and share provider conformance/round-trip/fuzz cases.

- Implement ASTTransformer, typed adapters and owned AST/CST registration
  envelopes with explicit language profiles and marker checks.
- Implement an empty-by-default syntax engine with explicit provider/edge
  registration, unique shortest paths, explicit routes, checked execution and
  existing native text composition. Add independent graph-oracle, ownership,
  negative public API and concurrent-registry tests, including a remote race gate.
- Keep canonical AST and actual language mappings deferred; record the delivered
  frame in [docs/engine.md](docs/engine.md).

- Add the generic CST and borrowed AST-view interface. Native parsers retain
  source tokens/trivia, formatters construct CST directly, and the shared printer
  emits supplied leaves without knowing languages or formatting conventions.
- Implement native TypeScript, Go, Rust and C++ parsing, formatting and syntax
  checks, with independent AST goldens, source round trips, malformed-input
  regressions, fuzzing and remote native compiler qualification.
- Replace the identity interface with LanguageProvider (ID/Name), embedded by
  every specialized text provider; compose providers explicitly. That native-text
  slice added no alias, registry, canonical AST or generator.
- Document the pipeline and ownership in
  [docs/parsing-formatting.md](docs/parsing-formatting.md).

- Add TypeScript, Go, Rust and C++ AST vocabularies over the unchanged generic
  carrier, with independent target consumers and language-separation checks.
- Add an agnostic language identity interface and provider ID/Name methods
  for the owner's intended central registry. That earlier slice added no
  registry, canonical AST or native text stages.
- Document node shapes, native syntax limits and the updated dependency/layout
  boundary in [docs/syntax.md](docs/syntax.md). Generators remain unimplemented;
  native formatting is delivered by the later text slice above.
- Document the Go AST-to-AST generation pipeline, target layout, composition
  strategy and dependency boundaries in [docs/generation.md](docs/generation.md).
- Add generic AST carriers, language-owned node kinds, containment/reference
  links and provenance, plus the independent typed `Stage[I, O]` contract.
- Add remote Go formatting, type-separation and dependency tests, and vet.
  That initial foundation included no language vocabularies or stage executors.
- Add the proposed canonical value-type model and a types-only implementation
  plan for TypeScript, Go and Rust in [docs/model.md](docs/model.md).
- Replace stale bootstrap guidance in the README, charter and contributor
  instructions with the current scope. Function, Stream and Cell stay
  library-owned; the model proposal is not a released language.
- Add remote documentation checks for local links, JSON examples and whitespace.

## Historical bootstrap — 26 September 2026

- Established the repository charter and component-first source layout.
- The original broader scope is retained in Git history. It is superseded for
  current work by the 31 October direction.
- Unmerged research at
  [`91363a5`](https://github.com/Bitspark/bittype/tree/91363a569c62899132ec365103abde7e41c760d9/research-docs)
  remains historical design evidence, not the current model specification.
