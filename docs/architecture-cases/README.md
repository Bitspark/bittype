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
| [bittype/ACX-004: Generation emits file syntax before placement](0004-generation-file-syntax.md) | TypeScript corrected; C++, Rust and Go (package clause, and since 12 November 2026 the case of every name) remain open. Order of correction recorded 11 November 2026: Go and Rust visibility move first, C++ is deferred until its semantic provider | AC-02, AC-06 |
| [bittype/ACX-005: Existing callbacks prescribed the redesigned interface](0005-generation-interface-design.md) | Corrected and qualified by declarative planning and native construction | AC-09 |
| [bittype/ACX-006: Common semantic traversal delegated to source integrations](0006-generation-source-traversal.md) | Corrected and independently verified: #191 was rejected on 8 November 2026; #192, correcting exact evidence and selector identity, was verified at `6132170` on 9 November 2026 | AC-02, AC-06, AC-01, AC-05, AC-07, AC-09 |
| [bittype/ACX-007: Matrix dependencies keyed by display names](0007-matrix-contract-identity.md) | Corrected, including unnamed argument identity | AC-05, AC-07 |
| [bittype/ACX-008: The shared horos mapping policy admitted one profile](0008-structural-profile-gate.md) | Corrected, with limits, state of 12 November 2026: a mapping states its profile in its own description and its constructor takes none; no shared code lists, writes or refuses a profile; the identities are `structural-representation/1` and `structural-enforcement/1`. A first form kept a refusal of one name in the adapter and was replaced before merging. Two independent reviews read the correction, at `fd75abe` and at `d61bc5d`. The registry does not judge a stated profile | AC-03, AC-02 |
| [bittype/ACX-009: Package boundaries do not follow contributions](0009-package-boundaries.md) | Corrected, with limits, state of 12 November 2026: the native and common horos mapping packages, the verification laws, the Go and TypeScript query services, the Go, Rust, TypeScript and horos identities, the horos names that a writer and a reader each spelled, the tag prefixes of canonical horos values, the horos syntax node kinds declared twice, the TypeScript graph names read by spelling, the symbol-named import check and the import rules by role. No confirmed instance remains. Two independent reviews read the moves of 11 November 2026 and the first commit of the Go move, pull request #246. None read its second commit or pull requests #247, #248, #250, #254 and #255. The limits name what is still not examined, which [#256](https://github.com/Bitspark/bittype/issues/256) tracks. Two candidates are not confirmed and are not part of the case. A third, the `Self` rule, was answered in principle on 12 November 2026 and is pending with the Go naming instance of bittype/ACX-004 | AC-10, AC-07 |

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
