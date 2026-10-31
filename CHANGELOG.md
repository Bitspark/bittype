# Changelog

## Unreleased

### 31 October 2026

- Add the generic CST and borrowed AST-view interface. Native parsers retain
  source tokens/trivia, formatters construct CST directly, and the shared printer
  emits supplied leaves without knowing languages or formatting conventions.
- Implement native TypeScript, Go, Rust and C++ parsing, formatting and syntax
  checks, with independent AST goldens, source round trips, malformed-input
  regressions, fuzzing and remote native compiler qualification.
- Replace the identity interface with LanguageProvider (ID/Name), embedded by
  every specialized text provider; compose providers explicitly. No alias,
  registry implementation, canonical AST or generator is added.
- Document the pipeline and ownership in
  [docs/parsing-formatting.md](docs/parsing-formatting.md).

- Add TypeScript, Go, Rust and C++ AST vocabularies over the unchanged generic
  carrier, with independent target consumers and language-separation checks.
- Add an agnostic language identity interface and provider ID/Name methods
  for the owner's intended central registry. Registry implementation remains
  deferred; that earlier slice added no canonical AST or native text stages.
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
