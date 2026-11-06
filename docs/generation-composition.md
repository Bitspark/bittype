# Automatic generation composition

Implemented in Go source, 6 November 2026. This is the production composition of
[Use, Publish, Define and Introduce/Apply](generation-spi.md). The four former
private pagination Dialects and `internal/horosgen` have been deleted. The native
representation profile remains types and interfaces only; native constraint
enforcement, runtime bindings and protocol behavior are separate work.

## Ownership and dependencies

| Component | Responsibility |
| --- | --- |
| `contracts/generation/go` | Typed operations, authenticated products, lexical scope and guarded invocation lifecycle |
| `generation/syntax/go` | Exact selection, guarded overrides, requested-root publication and definition scheduling |
| `generation/syntax/source/horos/go` | Read the supplied Horos semantic receipt; expose declaration identities, actual arguments, real binders and descriptor operands |
| `generation/syntax/target/go` | Invocation-local symbolic native fragments, reservations, immutable layouts and bounded assembly |
| `generation/syntax/target/{typescript,golang,rust,cpp}/go` | Native grammatical construction, symbols and output units; no canonical representation catalogue |
| `transforms/horos-to-*/go` | Thin endpoint-qualified ASTTransformer factories requiring a configured composition |
| Horos `tooling/bittype` module | Horos-owned builtin mappings, ordinary structural declaration policy and optional four-target default compositions |
| Other type libraries and consuming applications | Exact library mappings, native artifact assumptions and explicit consumer overrides |

The optional [Horos integration module](https://github.com/Bitspark/horos/tree/main/tooling/bittype)
depends on Bittype. Bittype and Horos core do not depend on that module. There are
no builtin-mapping imports in the composer or thin transformer factories. A
library can implement the public operations directly; using Horos's optional
mapping helpers is not required. Native backends have no Horos/Ontos dependency.
The [architecture tests](../integration/go/architecture_test.go) enforce these
import boundaries.

## Construction and selection

An invocation factory receives the existing exact `transform.Request[F]`. It
constructs a source view, native store and fresh mapping state, then returns:

```go
Composition[F, T]{
    Config:   generation.Config[F, T]{Source: source, Native: native,
                  Target: target, Evidence: requirements},
    Catalog:  source,
    Bindings: bindings,
    Roots:    requestedDeclarations,
    Helpers:  nativeOnlyPlans,
    Limits:   limits,
}
```

This is a field-shape illustration; concrete imports and setup are in the
[Horos defaults](https://github.com/Bitspark/horos/tree/main/tooling/bittype/transforms).
The composer supplies Config's Plans and Select; prepopulating either is invalid.
The Catalog reports declaration selectors and eligibility for an explicitly
supplied structural policy. Inventory does not lower a body. Source admission
still covers the complete original application, including arguments a mapping
may omit from its native output.

<!-- bittype:overview scope=construction-and-selection part=features -->
A Selector combines a compatible source language/profile/Go marker with an
opaque key. Display names and graph-local IDs are not selector identity. The
Horos adapter constructs keys for exact module-pin/declaration pairs, builtin
constructors with literal discriminators, and bound variables. Only that adapter
interprets Horos identity. Native choices never follow a name such as Cell or
Data, a kind classification, or an expanded structural fingerprint.
<!-- bittype:overview:end -->

<!-- bittype:overview scope=construction-and-selection part=exclusions -->
`NewBindings` captures an effective table. Duplicate exact bindings fail. An
override must name the existing Mapping.ID in ExpectedID and provide a different
new ID; missing, stale, repeated and foreign-profile overrides fail. Exact
bindings precede the explicitly supplied structural policy. Failure of a
selected mapping ends the invocation; no fallback or automatic specialization
is attempted. Factories create invocation-local collaborators; their Plan and
selection policy must remain stable and must not lower bodies speculatively.
<!-- bittype:overview:end -->

## Demands, generics and final output

The composer publishes every requested root. A mapping lacking Publish cannot
silently drop that root, even if it can represent a use inline. Imported source
declarations are emitted only when selected rules demand them. An exact external
or inline library mapping can avoid expanding an unsupported source body; source
semantic admission and identity are still required.

Acquisition reserves a destination and freezes its single Introduce layout.
Apply receives the caller's actual arguments and formal environment and returns
the complete result association. Its result is not the argument's association
relabelled. The same open Page is reused by Page<Person>, Page<Page<Person>> and
Page<Cell<Person>>. No first-use specialization changes Page's published header.

After root publication, Run drains Session.PendingDefinitions until every
acquired owned declaration or helper is defined. Defining one body may demand
more definitions. External and inline representations need no local definition.
Native-only helpers use explicit keys and native parameters, with no synthetic
source declaration or borrowed caller binder. Session's explicit driving API
remains available; PendingDefinitions is scheduling data, not an eager Define.

The native store authenticates fragments and symbols and materializes fresh AST
occurrences with real local references. Formal visibility follows lexical frame,
not the current callback's target signature. It captures parameter layouts once,
rejects native name collisions, orders demanded definitions by their references,
and runs the native syntax checker before returning a complete ASTSet. The
current native stores reject cyclic definition dependencies; they do not invent
indirection. External native syntax references remain external, with no fabricated
local declaration link. Integration contracts and native compiler qualification
own the external artifact's meaning.

Cancellation, missing representations, exhausted budgets, unsupported native
syntax or failed operations return no usable output. Defaults use the SPI limits
in [contracts.go](../contracts/generation/go/contracts.go); mapping factories can
select stricter limits. No cache or live collaborator state crosses invocations.

## Extension examples and qualification

The [Horos composition tests](https://github.com/Bitspark/horos/blob/main/tooling/bittype/integration/go/composition_test.go)
exercise replacing the exact list mapping with ReadonlyArray in TypeScript.
Page<T>'s body composes that choice under T, and nested uses preserve it without
regenerating Page. A deliberately failing replacement is called once and yields
no fallback output.

The same tests define an ordinary imported Cell<T> whose body is wire<T>. The
default profile cannot lower Wire. Adding one exact library-owned mapping renders
Cell<T> as LibraryCell<T> inside Page, without changing the composer, source
frontend, native backend, runtime or Wire. The fixture supplies the external
native declaration to the TypeScript compiler. This is type representation only;
it does not implement Cell behavior or install a privileged Cell package.

The [independent-source fixture](../integration/go/generation_composition_test.go)
uses the same automatic composer without Horos. It exercises a generic with
additional native formals, nested whole-result associations, external ownership
and a native helper acquired while defining another declaration. Existing SPI
negative tests retain scope, replay, incomplete-root and mutation checks.

Horos's integration suite retains four-target pagination, imported identity,
source bound rejection, naming limits, native compilation, installed-consumer
qualification and [verification of actual generated Go output](https://github.com/Bitspark/horos/blob/main/tooling/bittype/integration/go/verification_test.go),
including nested Page reuse through the landed Logos checker. Generation products
are not those verification results. Source analysis, native analysis, mapping
meaning and final-output checking remain separately owned and invoked.
