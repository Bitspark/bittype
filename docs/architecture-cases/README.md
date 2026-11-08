# Bittype architectural cases

Bittype owns the records of architectural violations in its designs and implementations.
Use the shared [numbered architecture rules](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#numbered-rules)
and the [case format](https://github.com/Bitspark/platform/blob/main/docs/concepts/architecture-cases/_template.md),
with repository-qualified references when linking across projects. The shared
rules do not transfer ownership of these cases to the platform repository.

| Case | Status | Primary violated rules |
| --- | --- | --- |
| [bittype/ACX-001: Closed verification contributions](0001-verification-composition.md) | Corrected for scoped static verification | AC-02, AC-03, AC-04 |
| [bittype/ACX-002: Generation assembly escaped the kernel](0002-generation-assembly.md) | Corrected and qualified | AC-02, AC-03, AC-05 |
| [bittype/ACX-003: Files stage coverage assumed Horos](0003-files-stage-coverage.md) | Corrected in stage coverage and test composition | AC-03, AC-06 |
| [bittype/ACX-004: Generation emits file syntax before placement](0004-generation-file-syntax.md) | TypeScript corrected; C++ and Rust remain open | AC-02, AC-06 |
| [bittype/ACX-005: Existing callbacks prescribed the redesigned interface](0005-generation-interface-design.md) | Corrected in the proposal. Implementation pending. | AC-09 |

ACX-001 continues the same incident first recorded as
[platform/ACX-001](https://github.com/Bitspark/platform/blob/main/docs/concepts/architecture-cases/0001-bittype-verification-composition.md).
Its detailed findings and correction status are maintained here. The platform
page is a cross-project reference, not a second maintained case.

Update a case as its understanding and correction develop. Explain the actual
ownership failure, applicable rules, source sections, required correction and
qualification limits. A proposal, passing example or package move does not close
an architectural defect. Update the same record when the repair is implemented
and qualified; keep partial repairs explicitly pending. Do not add a score,
keyword gate, source archive or separate evidence workflow.
