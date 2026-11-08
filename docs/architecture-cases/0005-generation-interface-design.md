# bittype/ACX-005 — Existing callbacks prescribed the redesigned interface

**Scope and status.** The 8 November 2026 correction concerns the generation
planner design. The declarative mapping SPI and the separation of planning from
AST construction are implemented and qualified. Source-reader narrowing remains
a separate proposed change.

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

**Implemented correction.** Libraries supply mapping descriptions. The
planner fixes the selected expressions, owned definitions, binder conventions
and public destinations. Construction receives NativeDeclaration records without the source binding or
application templates used to reach them. The target builder realizes those decisions under the
scope and identity guarantees of the run. An inline type needs a use expression,
an owned declaration needs a body, and public exposure follows the caller's
request. Reuse of code can be assessed after this interface is selected.
[Revised mapper design](../generation-planner.md#what-a-mapper-declares) ·
[Execution and qualification](../generation-planner.md#execution-guarantees).


**Qualification.** The Page<Cell> plan is inspectable before any Build call.
Repeated runs do not recapture source uses or reexecute mapping expansion.
Changing an expansion's mutable flag after planning cannot replace the array
form it selected. A separate Pair library reverses two native arguments inside
Page without kernel or target changes. An independent non-Horos source exercises
an additional native parameter and a helper under the same convention.
[Plan tests](../../integration/go/generation_planning_test.go) ·
[Independent source](../../integration/go/generation_independent_test.go)

**Limits.** The native target remains the trusted interpreter of its own forms.
This establishes a single mapping authority and separation of stages, not a
cross-language proof of every native lowering. Kind enforcement and the third
source-traversal SPI change are not claimed by this correction. The old
procedural SPI, second mapping authority and operation tickets were removed.
