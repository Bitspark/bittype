# Automatic generation composition

The engine assembles a generation request from supplied source, target and mapping registrations. See the [kernel construction SPI and pagination walkthrough](kernel-composition.md) and [generation operation SPI](generation-spi.md).

## Ownership and dependencies

| Package | Responsibility |
| --- | --- |
| contracts/generation/go | Construction and operation contracts, scoped issuance and product association |
| transform/syntax/shared/generation/go | Registration, checked selection, exact replacement and definition scheduling |
| transform/syntax/horos/source/go | Read the supplied semantic receipt and authenticate Horos declarations, binders and operands |
| transform/syntax/shared/target/go | Symbolic native fragments, layouts and bounded AST assembly |
| transform/syntax/&lt;native&gt;/target/go | Native construction, scopes and target syntax |
| transform/syntax/horos/mapping/go | Horos structural declaration policy and associations |
| transform/syntax/horos-to-&lt;native&gt;/mapping/go | Selectable Horos representations for each target |

## Construction and selection

Consumers supply providers once and select their registrations per request. The source resolves explicit declaration sites. Each mapping opens against the same source and native instances. The kernel collects the results, selects exact bindings, publishes the requested roots and schedules acquired definitions. The source authenticates applications including actual arguments omitted from the native representation. See the [registration code](../transform/syntax/shared/generation/go/registry.go) and [composition code](../transform/syntax/shared/generation/go/adapter.go).

<!-- bittype:overview scope=construction-and-selection part=features -->
A Selector combines a compatible source language/profile/Go marker with an
opaque key. Display names and graph-local IDs are not selector identity. The
Horos adapter constructs keys for exact module-pin/declaration pairs, builtin
constructors with literal discriminators, and bound variables. Only that adapter
interprets Horos identity. Native choices never follow a name such as Cell or
Data, a kind classification, or an expanded structural fingerprint.
<!-- bittype:overview:end -->

<!-- bittype:overview scope=construction-and-selection part=exclusions -->
A replacement names the selector, the expected provider and binding ID, and the new provider and binding ID. The kernel rejects ambiguous, stale, repeated or unused replacements and conflicting helpers. It applies the selected structural policy only where the source declares eligibility. A failed mapping ends the invocation.
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

After root publication, the kernel drains Session.PendingDefinitions until every
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
syntax or failed operations return no usable output. The request supplies explicit SPI limits
in [contracts.go](../contracts/generation/go/contracts.go); consumers can select stricter limits. No cache or live collaborator state crosses invocations.

## Extension examples and qualification

The [Horos composition tests](../integration/go/horos_generation_composition_test.go)
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

Bittype's integration suite retains four-target pagination, imported identity,
source bound rejection, naming limits, native compilation, installed-consumer
qualification and [verification of actual generated Go output](../integration/go/horos_generation_verification_test.go),
including nested Page reuse through the landed Logos checker. Generation products
are not those verification results. Source analysis, native analysis, mapping
meaning and final-output checking remain separately owned and invoked.
