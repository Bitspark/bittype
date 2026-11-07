# bittype/ACX-005 — Existing callbacks prescribed the redesigned interface

**Scope and status.** The 8 November 2026 correction concerns the generation
planner proposal. The revised mapper SPI and its qualification are pending
implementation.

**Decision that violated AC-09.** The original proposal required the executor
to use the five existing generation operations. It also gave a mapper two
surfaces: a declarative description of the intended output and an `Open`
method that returned executable contributions. The rationale for retaining
the callbacks was that they already carried binder and scope guarantees.
Those guarantees are real requirements. That reasoning did not independently
justify the same public methods or a second authority over the representation.
[Original mapper sketch](https://github.com/Bitspark/bittype/blob/dd10caa/docs/generation-planner.md#what-a-mapper-declares) ·
[Original execution proposal](https://github.com/Bitspark/bittype/blob/dd10caa/docs/generation-planner.md#execution-guarantees) ·
[AC-09](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-09).

**Concrete consequence.** A rule could describe an array of the mapped element
while its executable contribution constructed a readonly array of that same
element. Both could name the same binding and child. Checking those identities
alone would leave the planner unable to guarantee that execution realized the
representation it showed the caller.

**Correction in the proposal.** Libraries supply mapping descriptions. The
planner fixes the selected expressions, owned definitions, binder conventions
and public destinations. The target builder realizes those decisions under the
scope and identity guarantees of the run. An inline type needs a use expression,
an owned declaration needs a body, and public exposure follows the caller's
request. Reuse of code can be assessed after this interface is selected.
[Revised mapper design](../generation-planner.md#what-a-mapper-declares) ·
[Execution and qualification](../generation-planner.md#execution-guarantees).
