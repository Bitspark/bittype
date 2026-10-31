# AST generation and pipeline design

**Direction: 31 October 2026.** Generators and formatters will be implemented
in Go. The present slice implements only the generic AST data types and typed
stage contract. No parser, language grammar, generator, formatter, composition
executor or command is implemented.

The canonical model remains a [proposal](model.md). Its syntax and semantic
decisions must be settled before a language implementation claims conformance.

## Target pipeline

```text
canonical source
    -> parse, resolve and check once
    -> AST[Canonical]
         -> AST[TypeScript] -> TypeScript source
         -> AST[Go]         -> Go source
         -> AST[Rust]       -> Rust source
```

Each generator transforms a checked canonical AST into a target-language AST.
It chooses native type representations, names, imports, generic syntax and
recursive indirection according to an explicit target configuration.
It does not emit source strings.

Each formatter consumes only its target-language AST and prints source text.
It owns escaping, lexical spelling, layout and precedence. It does not resolve
canonical declarations or decide how a canonical nullable field should be
represented. Those choices must already be explicit in the target AST.

Parsing, reference resolution and declaration checking belong to the canonical
frontend. File IO, dependency retrieval, configuration and output paths belong
to the application. A formatter is usable with a target AST built by another
producer, without a canonical input or generator dependency.

The output scope is native types and interfaces, including generic declarations
and applications. Executable bodies, value validators, serialization, callable
implementations, adapters and transport behavior are outside this milestone.

## Go layout

There is one root Go module, `github.com/Bitspark/bittype`. Components are
separate packages; separate modules are unnecessary for this initial foundation.
The final path component denotes the implementation language, not the target.
For example, `syntax/rust/go` is Go code describing Rust syntax.

Only the entries marked **present** exist in the initial implementation.
Other paths are the target layout, not placeholder packages.

```text
go.mod                                      present
ast/go/
  ast.go, node.go, origin.go                 present: data declarations
  walk.go                                   planned: generic traversal
pipeline/go/
  stage.go                                  present: typed contract
  compose.go                                planned: typed composition
syntax/
  canonical/go/{nodes,grammar,builders}.go   planned
  typescript/go/{nodes,grammar,builders}.go  planned
  golang/go/{nodes,grammar,builders}.go      planned
  rust/go/{nodes,grammar,builders}.go        planned
frontend/canonical/go/
  parse.go, resolve.go, check.go             planned
generate/
  typescript/go/{generate,types,names,imports}.go
  golang/go/{generate,types,names,imports}.go
  rust/go/{generate,types,names,imports}.go   all planned
format/
  typescript/go/format.go                    planned
  golang/go/{format,stdast}.go               planned
  rust/go/format.go                          planned
cmd/bittype/go/
  main.go, compose.go, files.go              planned
integration/go/
  contracts_test.go, dependencies_test.go    present
  testdata/*.go                             present: API type-check fixtures
testdata/
  canonical/, target-asts/, source-text/     planned: language fixtures
docs/
  model.md, ast.md, generation.md            present
```

The target grammars cover the syntax needed for native type declarations.
They are not full programming-language parsers. A concrete unsupported type
case should identify a missing grammatical form; it must not lead to a raw
code-string escape hatch.

## Typed stage contract

The [pipeline package](../pipeline/go/stage.go) defines:

```go
type Stage[I, O any] func(I) (O, error)
```

A generator will satisfy a stage from one AST type to another; a formatter
will satisfy a stage from a target AST to source text. For example, after the
language markers exist, the type relationships are:

```text
Stage[AST[Canonical], AST[TypeScript]]
Stage[AST[TypeScript], string]
```

`Stage` imports nothing and does not know that ASTs or languages exist.
The function type is the interface between producers and consumers; no
additional object interface or function adapter is needed. Concrete stage
factories may later close over explicit target configuration.

A stage treats input and reachable storage as read-only. It owns mutable
output storage it creates. On error, callers must discard the output.
The zero/nil stage means no configured stage and must not be invoked.
These behavioral requirements are documented contracts, not enforcement
provided by the function type.

The following signature describes future composition, **not a present API**:

```go
func Then[A, B, C any](first Stage[A, B], next Stage[B, C]) Stage[A, C]
```

It would connect stages only when the intermediate type agrees and propagate
the first failure. There is no `Then` implementation or stub in this slice.
Scheduling, caching, cancellation and diagnostic aggregation are also deferred.

## Composition and dependencies

The application root will explicitly compose selected concrete stages.
It resolves input once and passes the resulting canonical AST to independently
selected target generators. Each formatter receives the corresponding target
AST. Shared inputs remain read-only by contract; outputs have separate storage.
Whether targets run sequentially or concurrently is an application decision,
not a requirement imposed by the AST model.

Arrows in the table mean a direct package dependency is allowed:

| Component | Dependencies |
| --- | --- |
| Generic AST | Standard library only; currently no imports |
| Generic pipeline | Standard library only; currently no imports, including no AST dependency |
| Each syntax package | Generic AST |
| Canonical frontend | Canonical syntax, generic AST |
| Each generator | Canonical syntax, its own target syntax, generic AST; pipeline only when exposing a stage |
| Each formatter | Its own target syntax, generic AST; pipeline only when exposing a stage |
| Application composition root | Explicitly selected frontend, generators, formatters and pipeline |
| Integration tests | The APIs under test and native compiler/checker tooling |

Standard-library utilities are permitted in every layer. There are no external
Go module dependencies in this slice and no dependencies on another Bitspark
implementation repository. The AST and pipeline packages do not import each
other. Syntax packages do not import generators or formatters. A formatter does
not import the canonical grammar or its generator.

The planned Go formatter can adapt its language AST to the standard
[`go/ast`](https://pkg.go.dev/go/ast) representation and use
[`go/format`](https://pkg.go.dev/go/format).
That adaptation belongs to the Go formatter, not to the generic AST.
TypeScript and Rust printers are Go implementations too; Node and Rust tooling
are needed for later output qualification, not as generation-time dependencies.

No global backend discovery, hidden default registrations or contract catalogue
is introduced. The generic AST and pipeline must never acquire Function,
Stream, Cell, Wire or runtime imports, name switches, special kinds or structural
fingerprint matching. Adding an ordinary library contract must not change these
foundations. Later library-owned realizations require a separate scoped design;
this slice supplies no binding extension registry.

## Evidence and next increments

The present Go module targets Go 1.27.0 language features; CI uses the pinned
Go 1.27.3 toolchain. The [Go workflow](../.github/workflows/go.yml) checks the
exact committed revision on a GitHub runner: formatting, tests and vet.
[Documentation checks](../.github/workflows/docs.yml) also run remotely.
No compilation or formatting workload is run on the workstation.

Public API fixtures establish that an independent consumer can declare a
language without central registration and that incompatible language/stage
types are rejected. A production dependency-closure check guards the empty
foundation dependency boundary. There are no grammar or generation correctness
claims yet.

Subsequent, separately authorized increments are:

1. Specify the canonical grammar and independent semantic fixtures.
2. Implement language-owned AST vocabularies, checks and builders.
3. Implement the canonical frontend and AST-to-AST generators.
4. Implement independent target formatters and application composition.
5. Compile the generated TypeScript, Go and Rust examples with pinned native
   toolchains on shared workers or GitHub runners.

Generator tests will compare target ASTs. Formatter tests will accept
independently authored target ASTs. End-to-end tests will compile and use the
same semantic examples, including open/applied pagination and presence/null
distinctions. Compiler success alone will not prove all canonical value
invariants.
