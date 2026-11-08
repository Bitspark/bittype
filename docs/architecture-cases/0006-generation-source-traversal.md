# bittype/ACX-006 — Common semantic traversal delegated to source integrations

**Status: open, 8 November 2026.** The corrective source-reader SPI is not
finalized or implemented. [GEN-01](../generation-gaps.md#gen-01) owns the work
description and completion criteria.

## Actual ownership failure

The public generation Source/Catalog contract asks each source integration to
supply navigation of declarations, binders, references and application arguments
as well as interpretation of source meanings. Horos's `Open` and `Admit` perform
both roles. The shared Semantics model already carries these structural facts
and exposes a direct-application binding query.

The concrete consequence is that an independent source has to recreate a
navigation adapter for facts that the common semantic context already owns.
The kernel receives an additional source-specific interpretation of common
structure instead of traversing the established graph once. The independent
fixture demonstrates composability through the present interface, but still
reimplements that responsibility.
[Source contract](../../contracts/generation/go/source.go) ·
[Horos implementation](../../transform/syntax/horos/source/go/source.go) ·
[Shared semantic queries](../../semantics/go/model.go) ·
[Independent reader](../../integration/go/testdata/generationfixture/source.go)

This finding does not assert a Horos dispatch table in the kernel or absent
source checking. It concerns where common traversal lives. Separating semantic
planning from AST construction in PR #155 did not repair this different boundary.

## Rules involved

- [AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02):
  common structural navigation belongs to the shared mechanism; source
  interpretation belongs to its language owner.
- [AC-06](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-06):
  successful generation through the first reader does not redefine the intended
  kernel/language division or make its correction optional.
- [AC-01](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-01):
  describe this remaining architectural defect explicitly instead of claiming
  that passing composition tests complete the architecture.

## Required correction and qualification

Move traversal of established graph structure into the kernel and design narrow
source interpretation contracts. Preserve snapshot/evidence association, real
binder identities, scope, complete argument admission and failure meaning.
A structural binding query does not establish that an application is allowed;
Horos still owns its sort/kind judgments. Do not transplant all of `Admit` into
the kernel or retain incumbent signatures as a design constraint.

Qualify the same traversal with Horos and an independent source, valid nested
applications and negative scope/evidence cases. Remove the superseded reader
surface and its duplicated traversal. Follow AC-09 when selecting the replacement
and AC-05 when preserving guarantees; neither unrestricted callbacks nor a
package move is a correction. Update this same case when that work is delivered.

GEN-02's unsupported native enforcement guarantee is separate. Declared
representation-profile erasure is not evidence that GEN-01 is acceptable, and
fixing GEN-01 does not implement native constraint enforcement.
