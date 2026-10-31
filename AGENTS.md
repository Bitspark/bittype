# Working here as an agent

Read [CHARTER.md](CHARTER.md), [docs/model.md](docs/model.md) and
[LAYOUT.md](LAYOUT.md) before changing this tree.

## Current scope

The owner selected canonical definitions as input and TypeScript, Go and Rust
types as outputs. The model document is a proposal; distinguish proposed rules,
implemented behavior, published revisions and delivered consumer uptake.

- Keep the model independent of Wire, storage formats and target languages.
- Function, Stream, Cell and future behavioral contracts are library-owned.
  Do not add their names as core kinds, switches, imports, registrations,
  fingerprints or hidden defaults.
- Preserve declaration origins, generic slots and applications. Equal-looking
  structures do not authorize substituting one library contract for another.
- This milestone generates types only. Do not add adapters, callable behavior,
  runtime value validators, serialization or networking.
- Do not turn TypeScript declarations into the authoritative input again.
- A general generation kernel, if needed, stays in a separate module and does
  not import the declaration language or contract libraries.
- Use the current model proposal as the working document. Prior research and
  Nightseam are evidence to inspect, not specifications to copy or restore as
  competing active documents.
- Settle normalization and identity content before fixing an encoding. Preserve
  the distinction between declaration identity, value shape, generated-code
  compatibility and behavior.

## Delivery

Use a task branch or worktree and a pull request. Squash a green change onto
`main`; keep the designated primary checkout on `main`.

For documentation, check relative Markdown links, JSON example syntax and
`git diff --check`, including staged additions. The documentation workflow runs
these checks on the exact committed revision on a GitHub runner. Portable
compilation and later language checks also run on the shared build service or
GitHub runners, using the maintained repository gate when one exists.

Keep historical commits, research evidence, license notices and unrelated work.
Do not add private checkout dependencies, local coordination state or credentials.
Repository names are lowercase; the company name is Bitspark.
