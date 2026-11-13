# bittype/ACX-010: An agnostic forwarding mechanism is owned by Horos

**Status.** Identified on 13 November 2026. Implementation correction pending.

**Scope.** The function-to-interface mechanism in
`transform/syntax/horos/mapping/go.Library[T]`, inspected on main at
`e731e506da15c375bd74bd87b813a367355595bf`. The owner clarified the placement
requirement after a review treated successful independent registration as
sufficient justification for leaving this helper with Horos.
[Inspected implementation](https://github.com/Bitspark/bittype/blob/e731e506da15c375bd74bd87b813a367355595bf/transform/syntax/horos/mapping/go/mapping.go)

## Rules involved

[AC-11](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-11)
requires an agnostic mechanism to live in the agnostic layer from its first
implementation. This adapter's complete mechanism checks whether a function
is nil and forwards its arguments and result. That responsibility and its
failure rule are independent of Horos. Its declaration fixes Horos types,
placing the mechanism inside one language integration.

[AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)
provides the semantic ownership principle.
[AC-03](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-03)
and [AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07)
explain the assessment error. Open registration establishes an extension
property. Mechanism placement requires its own assessment.
[AC-09](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-09)
and [AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
govern the intended interface and correction order.

## Intended boundary

The agnostic generation contract owns the mechanics by which a description
function satisfies its public mapping interface. Language and mapping owners
supply the concrete source type, target type, interpretation reader and mapping
description. A library owns its supported profile and native representation
meaning.

This follows the [agnostic versus language-specific matrix](../abstraction-levels.md#agnostic-kernel-and-language-specific-rules)
and [kernel-owned composition](../kernel-composition.md). The invariant mechanism
has one agnostic implementation, instantiated with each contribution's types.

## Actual implementation and consequence

The public SPI is already generic. The mapping owner implements it and the
generation registry calls `Describe` at registration.

```go
// Implemented in contracts/generation/go.
type MappingLibrary[F, T, R any] interface {
    Describe(from, to representation.Language) (
        MappingDescription[F, T, R], error,
    )
}
```

The complete forwarding helper below is implemented in the Horos mapping
package. In this excerpt, `g` denotes the generation contract, `h` the Horos
syntax vocabulary, and `source` its generation interpretation interface.

```go
type Library[T any] func(
    from, to representation.Language,
) (g.MappingDescription[h.Language, T, source.Reader], error)

func (p Library[T]) Describe(
    from, to representation.Language,
) (g.MappingDescription[h.Language, T, source.Reader], error) {
    if p == nil {
        return g.MappingDescription[h.Language, T, source.Reader]{},
            g.ErrInvalid
    }
    return p(from, to)
}
```

The helper interprets no descriptor, binder, kind constraint or native form.
Substituting the source and reader types leaves its operation and nil failure
identical. Its fixed parameters prevent another source from instantiating this
same helper for its own types. The generic public SPI allows a direct
implementation, but repeating that forwarding mechanism in every language
would perpetuate the ownership error.

Horos-specific interpretation still belongs to Horos. The concrete list,
record and other mappings still own their native representation choices.
Their meaningful work does not move merely because this forwarding helper does.

## What worked and what it did not establish

The generic interface accepts independent mappings. The common adapter now
returns the mapping's own description and imposes no profile whitelist.
The profile policy correction remains recorded in [bittype/ACX-008](0008-structural-profile-gate.md).

Those properties explain why generation can work while the remaining helper
has the wrong owner. The present finding comes from source inspection and the
owner's explicit clarification. It identifies a mechanism already expressible
through the shared typed contract. Qualification of the placement correction
remains pending.

## Required correction

1. Derive the function-adaptation responsibility from the general generation
   contract. If retained, implement it once in the agnostic layer with source,
   target and reader types supplied by its users.
2. Migrate the language-owned mappings to that mechanism and remove the
   Horos-specific forwarding implementation. Preserve the mappings' independent
   descriptions and profile ownership.
3. Preserve the nil-function failure, language and marker checks, description
   capture and scope guarantees. The registry retains selection and composition.
4. Exercise the shared helper with distinct source types and inspect the actual
   dependency direction. Verify that each language only supplies its own
   variation and the mechanism has one implementation.

## Correction and qualification status

This change documents the requirement and the confirmed case. The executable
helper and its callers still need correction. The package moves in
[bittype/ACX-009](0009-package-boundaries.md) and the profile correction in ACX-008
remain qualified for their recorded scopes. This case owns the additional
agnostic-placement finding.

## Assessment lesson

The review asked whether a new contribution could enter without changing the
kernel and answered yes. It then used that result to excuse the placement of
an invariant helper inside one source language. The owner clarified that
agnostic placement is independently required, including when only one language
currently uses the mechanism.

Assess the knowledge and guarantees of the mechanism as well as the openness
of the public SPI. Fixed language types, a small body and the possibility of
writing another implementation cannot justify per-language copies of
language-independent work. Code similarity alone is insufficient evidence
that two pieces of work share that responsibility.
