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
| [bittype/ACX-006: Common semantic traversal delegated to source integrations](0006-generation-source-traversal.md) | Traversal corrected and independently verified at `6132170` on 9 November 2026 after #191 was rejected and #192 corrected it. Residual common name/site access through source readers identified on 13 November 2026; correction pending | AC-02, AC-06, AC-01, AC-05, AC-07, AC-09, AC-11 |
| [bittype/ACX-007: Matrix dependencies keyed by display names](0007-matrix-contract-identity.md) | Corrected, including unnamed argument identity | AC-05, AC-07 |
| [bittype/ACX-008: The shared horos mapping policy admitted one profile](0008-structural-profile-gate.md) | Corrected, with limits, state of 12 November 2026: a mapping states its profile in its own description and its constructor takes none; no shared code lists, writes or refuses a profile; the identities are `structural-representation/1` and `structural-enforcement/1`. A first form kept a refusal of one name in the adapter and was replaced before merging. Two independent reviews read the correction, at `fd75abe` and at `d61bc5d`. The registry does not judge a stated profile | AC-03, AC-02 |
| [bittype/ACX-009: Package boundaries do not follow contributions](0009-package-boundaries.md) | Corrected, with limits, state of 12 November 2026: the native and common horos mapping packages, the verification laws, the Go and TypeScript query services, the Go, Rust, TypeScript and horos identities, the horos names that a writer and a reader each spelled, the tag prefixes of canonical horos values, the horos syntax node kinds declared twice, the TypeScript graph names read by spelling, the symbol-named import check and the import rules by role. No confirmed instance remains. Two independent reviews read the moves of 11 November 2026 and the first commit of the Go move, pull request #246. None read its second commit or pull requests #247, #248, #250, #254 and #255. The limits name what is still not examined, which [#256](https://github.com/Bitspark/bittype/issues/256) tracks. Two candidates are not confirmed and are not part of the case. A third, the `Self` rule, was answered in principle on 12 November 2026 and is pending with the Go naming instance of bittype/ACX-004 | AC-10, AC-07 |
| [bittype/ACX-010: An agnostic forwarding mechanism is owned by Horos](0010-agnostic-mapping-adapter.md) | Identified on 13 November 2026. The function-to-interface adapter belongs in the agnostic generation layer, with source, target and reader supplied by its users. Implementation correction pending | AC-11 |
| [bittype/ACX-011: Common Files semantic orchestration is implemented per language](0011-files-semantic-orchestration.md) | Identified on 13 November 2026. Horos and TypeScript separately assemble the same capture/resolve/analyze/evidence sequence. Language organization and meaning checks stay with their owners. Correction pending | AC-02, AC-11 |
| [bittype/ACX-012: Verification-law role classification fixes one language pair](0012-verification-role-language-pair.md) | Identified on 13 November 2026. The shared architecture test recognizes only the literal horos-to-go law path. Correction pending; production law composition is a separate question | AC-03, AC-11 |
| [bittype/ACX-014: The bounded writer for a child process had one copy per caller](0014-bounded-process-writer.md) | Corrected, with limits, state of 13 November 2026: `internal/processoutput` holds the one writer. The Rust and TypeScript compiler bridges and the loader of the reporting tool use it, and no copy is left. The loader's budget, which its copy did not enforce, holds. The child process of the tests, which the correction first copied twice more, has one owner as well, `internal/processoutput/processoutputtest`. An independent review read the first form at `c8c3509`, and its findings are corrected. It read the second form at `9ff77f3`, and its verdict is to merge | AC-11, AC-05, AC-07 |
| [bittype/ACX-015: Models hold operations and import operations](0015-models-hold-operations.md) | Identified on 13 November 2026, after the owner directed that the root splits into models and operations. By the design's table 23 of 129 production packages hold a model-side part and an operation, and they make 17 imports of an operation. Correction pending, designed in [models and operations](../models-and-operations.md). For the four intended-forms packages the correction is made on the branch of the migration of #196, and `main` shows the misplacement until that branch merges, or until the split is made on `main` if its pull request closes unmerged. The platform's rules page names the case as the example of AC-12, which covers the grammar tables, the unpublished keys of the Go bridge and descriptions under a stage's name. Containment and the import rule on units are the owner's direction for bittype | AC-01, AC-02, AC-06, AC-11, AC-12 |

The number 0013 is taken by a case record on the open branch
`fix/definition-is-body-310` and is not on `main` yet.

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
