# Public generation operations

This implements the
[type-owned generation design](type-generation-plugins.md) and
[generic elaboration protocol](generic-elaboration.md). Its scope is public
operation contracts and the machinery needed to exercise them safely with
explicitly selected collaborators. [Automatic composition](generation-composition.md)
now selects mappings, publishes roots and schedules definitions; all four native
generators use this SPI. Native constraint enforcement and runtime bindings are separate work.

The [expert review](../research-docs/0011-public-generation-spi.md) addresses
construction, ownership and generic scope. The existing
[joint verification contract](verification.md) concerns completed captured
source/target artifacts. A generation product is not a final-output proof.

The public package is `github.com/Bitspark/bittype/contracts/generation/go`.
It contains both the operation contracts and their guarded lifecycle adapter,
following the existing semantic operation pattern. Keeping issuance with the
opaque handles avoids an unrestricted public constructor between packages.
The stage-owned `generation/syntax/go` composer uses this boundary for selection
and scheduling. Session also supports explicit driving; no forwarding alias or
second issuance path is introduced.

## Construction and driving

The consumer supplies `Config[F,T]`: a `Source[F]`, `Native[F,T]`, target language
descriptor, receipt requirements, preselected `Plan[F,T]` records and an explicit
selection callback over already admitted source uses. `Begin` captures plan
records and source-owner identities, admits the exact supplied snapshot/receipts,
and calls the configured source's `CheckInput`. It runs no semantic provider.
`transform.Request.Snapshot()` preserves this original capture even when the
transformer requires no receipts; recapturing equal syntax is a different input.

The sequence is explicit:

```go
run, err := generation.Begin[F, T](ctx, input, config, limits)
// Handle every error; any failed operation prevents successful Finish.
product, err := run.Use(use)
publication, err := run.Publish(declaration)
callee, err := run.Acquire(declaration)
definition, err := run.Define(callee)
output, err := run.Finish()
```

These are alternative demands in a larger consumer workflow, not a requirement
that every mapping implement every operation. Acquire is idempotent for the
same selected declaration within the invocation. A local owned reservation must
be defined exactly once before Finish. Publication is separate and may precede
definition; each requested root must be exposed through its publisher. Only
Finish returns a usable complete AST set. It fails if required definitions are
unfinished, a callback failed, cancellation occurred or a budget was exceeded.

`Plan.Header` supplies a target layout when no source convention is needed.
Source formals must be empty on that route. `AcquireHelper(key)` and
`UseRequest.RequireHelper(key)` address explicitly supplied native-only plans.
Their source declaration/environment are zero; a helper cannot inherit a
caller's source binders. Its target-owned header may still have explicit native
parameters. Helper keys are local to the captured configuration and are not
synthetic source identities.

A Session permits nested callbacks on one call stack. Copies share its entire
private invocation state; consumers serialize calls across all copies. Finish
freezes the obligation ledger before target assembly and rejects reentrant
operations. Recursive publication of the same root is also rejected. Separate
sessions can run concurrently when their collaborators support
that use. The host copies its owned slices and records; collaborators must keep
their private records and policy stable. Copying a Go interface cannot freeze a
mutable implementation.

## Typed addresses and authenticated products

`Space[K]` belongs to an integration storing immutable records. Its `Ref[K]`
contains a nonempty private owner identity and a slot; `Owner[K]` permits ownership
checks without issuing addresses. Roles distinguish declarations, real binders,
source uses/environments, target destinations/layouts/fragments and associations.
These are record locators, not a second semantic graph or a type algebra. A
nonzero reference alone proves neither existence nor applicability.

The source supplies read-only owners for its roles. The host checks these exact
owners and invokes `Source.Authenticate` for existing records. `Open` supplies
canonical binders and assumptions; `Admit` supplies a complete applicable use.
An occurrence and its effective subject remain separate, so two occurrences of
the same binder can reuse its product without pretending their AST addresses
are equal. Existing Snapshot, semantic Result and AST envelopes remain direct
inputs rather than being translated into new record spaces.

Mappings issue tokens into their own immutable association stores. A request's
`Return(token)` calls the authority fixed by its plan, then checks the returned
binding and asks the configured target to admit the actual payload fragment.
The public `Captured` construction record is only a callback result: there is no
public `AcceptCaptured` or free source/payload/support product constructor.
The candidate cannot select its own validator. `Product.Binding`, `Payload` and
`Token` expose immutable identities, not a mutable proof tree or a delayed reader
of whatever scope happens to be current.

`Binding` retains run, configuration, issuer, effective source subject,
assumption environment, native frame and applicable signature. A separate
callback ticket controls issuance/return. Formal products survive Introduce
into Define, and child products survive into parent Apply, while old requests
cannot issue and old publications/definitions cannot satisfy new callbacks.
The first slice requires exact frame equality. Even closed native syntax may
carry assumptions from another frame; lack of free native symbols is insufficient.

Each library can own a distinct association store. A convention consumes an
accepted argument through explicitly configured compatible premise readers.
There is no central product variant or string-keyed evidence dictionary to edit
when adding a library. Missing compatible support is an error.

## Responsibility boundaries

| Owner | Responsibility |
| --- | --- |
| Source integration | Interpret supplied semantic receipts and explicit query results; establish applicability, effective arguments and original assumptions |
| Selected mapping | Choose the representation and establish its association with the effective source use |
| Generic convention | Introduce native formals and elaborate complete applications using one retained convention |
| Target integration | Own native grammar, symbols, constraints, imports and assembly |
| Shared generation host | Associate requests, products and scopes; preserve lifecycle, limits and whole-invocation failure |
| Final verification composition | Check claims about the actual assembled target using its captured semantics and the selected Logos contract |

The shared code interprets no Horos constructor, kind or named library contract.
Function, Stream and Cell have the same extension path as any other library.
Builtin mappings also belong to supplied integrations. None of these operations
implements invocation, streaming, hydration or application-value conversion.

## Five operations

**Use** obtains the selected representation of an admitted source use.
**Publish** exposes a requested native declaration. **Define** fills a locally
owned reservation. These are independently demanded capabilities: an inline
mapping can supply Use without Publish or Define, and an external generic need
not define its callee locally. A requested publication cannot silently disappear.

**Introduce** establishes the native formal layout and products for the real
binders of one canonical open declaration. **Apply** uses that retained convention
in the caller's environment and returns the complete application's representation
and support. Native formal count need not equal source binder count. An extra
native descriptor or witness parameter creates no additional source binder.

These names are scoped to generation. Generation Publish does not publish a
semantic verdict, and generation Define is not the Logos logical-definition
constructor.

## Source and representation environments

The source substitution environment maps real source binders to effective source
arguments. The representation environment maps those binders to scoped native
products. These are different environments: names and positional argument indexes
cannot replace binder identity or retained assumptions.

For `Outer<U> = Page<U>`, the application consumes Outer.U's product. Page's own
formal cannot escape its definition environment. For
`Nested<U> = Page<Page<U>>`, the outer application consumes support for the inner
Page result; support for U alone does not describe that result. Admission also
covers unused arguments, as in `Ignore<T: K> = text`.

A helper has its own scope and explicit formals. Keeping an immutable handle to a
caller's binder does not make capturing that binder legal. Reusable syntax must
materialize fresh containment nodes; shared symbols remain references.

## What successful generation establishes

The host enforces association and lifecycle rules. Semantic correctness of a
mapping, applicability of a source use and adequacy of a native representation
remain explicit integration obligations. Trusted Go collaborators are not a
sandbox. Generic typing and struct copies do not establish immutable ownership.

Generation support may describe a symbolic representation under assumptions.
After assembly, final verification must capture the actual target, project its
native semantics and prepare the requested claim. An earlier receipt about a
different target cannot be relabelled onto newly emitted syntax. The existing
Horos/Go verification composition and public Logos checker remain separately
invoked; generation performs no hidden resolution, analysis or proof search.

## Qualification cases

The first slice must be usable from independent public consumers, including an
installed consumer with no private imports or local module replacements. Its
fixtures distinguish host checks from the supplied integrations' semantics.

| Case | Required behavior |
| --- | --- |
| Page and a concrete Person application | One retained open layout; the application uses its selected argument product |
| Outer and renamed binders | Apply retains the caller's actual binder and assumptions |
| Nested application | The outer use consumes the inner result association |
| More native formals than source binders | Exact source-binder coverage without positional one-to-one native matching |
| Nongeneric definition | Publish and Define work without an unnecessary generic capability |
| External generic | Apply works without owning or defining the external callee |
| Inline representation with a helper | Inline and owned effects can compose |
| Wrong support, foreign scope or changed layout | The invocation fails; no successful partial output |
| Missing requested capability | Explicit failure, with no fallback mapping or omitted root |
| Cancellation or exhausted limits | No successful completion or output publication |

Native compilation qualifies emitted syntax in the tested domain. It does not
establish source applicability, universal native admission, runtime conformance
or the correctness of every possible integration.
