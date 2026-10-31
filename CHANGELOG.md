# Changelog

## Unreleased

### 31 October 2026

- Document the Go AST-to-AST generation pipeline, target layout, composition
  strategy and dependency boundaries in [docs/generation.md](docs/generation.md).
- Add generic AST carriers, language-owned node kinds, containment/reference
  links and provenance, plus the independent typed `Stage[I, O]` contract.
- Add remote Go formatting, type-separation and dependency tests, and vet.
  No language grammars, generators, formatters or stage executors are included.
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
