# Research records and current contracts

Consultations record a question against a specific baseline, independent advice
and a dated evaluation. Original `.submitted*` and `.advice*` files are historical
evidence. A `.verify-N.md` file is the report of a verification agent on a
set of recommendations: a proposal, kept as it was written, and no decision. An `applied` status records that consultation's delivery; it does not
make its suggested interface today's specification. The [documentation map](../docs/README.md)
identifies the current owners. Evaluated decisions are retained without rewriting
what the expert or submitting author said at the time.

| Consultation | Current documentation and subsequent changes |
| --- | --- |
| [0002: semantic analysis](0002-semantic-analysis-kernel.md) | [Semantics](../docs/semantics.md), [analysis](../docs/analysis.md); Horos revision 4 uses callee references and ordinary reference terms. |
| [0003: optional Wire extensions](0003-optional-wire-extensions.md) | [Concept inventory](../docs/wire-extension-preservation.md), [unified declarations](../docs/unified-declarations.md); template lifting and the second frontend were removed. |
| [0005: general kinds](0005-general-kinds.md) | [Kinds](../docs/kinds.md); ordinary kinds and the source-authored prelude own their meaning. |
| [0006: prelude source](0006-prelude-source.md) | [Dependencies](../docs/dependencies.md), [module encoder](../tools/horos-module/README.md). |
| [0007: type-owned mappings](0007-type-owned-generator-plugins.md) | [Mapping design](../docs/type-generation-plugins.md), [current SPI](../docs/generation-spi.md); procedural operations and per-pair assembly are superseded. |
| [0008: generic protocol](0008-generic-generation-protocol.md) | [Generic elaboration](../docs/generic-elaboration.md); one declarative convention replaces operation tickets. |
| [0009: native constraints](0009-native-kind-constraint-soundness.md) | [Native constraint proposal](../docs/native-kind-constraints.md), [GEN-02](../docs/generation-gaps.md#gen-02); native enforcement is still absent. |
| [0010: reasoning architecture](0010-type-system-model-architecture.md) | [Verification architecture](../docs/type-system-architecture.md), [joint contract](../docs/verification.md); public association and composable interpretations are implemented for the bounded route. |
| [0011: public generation SPI](0011-public-generation-spi.md) | [Declarative SPI](../docs/generation-spi.md), [ACX-005](../docs/architecture-cases/0005-generation-interface-design.md); the delivered procedural surface was subsequently removed. |
| [0012: capability reports](0012-capability-reports.md) | [Capability reporting](../docs/capability-reports.md), [matrix identity correction](../docs/architecture-cases/0007-matrix-contract-identity.md). |
| [0013: source traversal](0013-generation-source-traversal.md) | [Semantic traversal and source interpretation](../docs/generation-spi.md#semantic-traversal-and-source-interpretation), [GEN-01](../docs/generation-gaps.md#gen-01), [ACX-006](../docs/architecture-cases/0006-generation-source-traversal.md); kernel-owned traversal replaced the 0011 Open/Admit reader. |
| [0014: binding generation](0014-binding-generation.md) | [Binding generation](../docs/binding-generation.md); bittype publishes a selected representation plan and per-target backends compose paired library contributions. Accepted direction; not implemented. |
| [0015: representation families](0015-representation-families.md) | [Representation families](../docs/representation-families.md); each family is one representation contract instantiated once per use, delivered with its plan payload, conversion and qualification. Accepted design; not implemented. |
| [0016: native Files placement](0016-native-files-placement.md) | [Source and Files](../docs/source.md), [bittype/ACX-004](../docs/architecture-cases/0004-generation-file-syntax.md); evaluated on 12 November 2026 and in delivery: language-owned placement types with one shared admission, a typed placement bound to its output, naming intent realized by Files, fragment roots, recorded destinations. The work items are #271 to #278, #161 and #159. |
| [0017: the contribution contract](0017-contribution-contract.md) | [Binding generation](../docs/binding-generation.md), [representation families](../docs/representation-families.md); evaluated on 13 November 2026 and in delivery: one decision record per selection with one round of requests before a rule decides, a registered projector, publisher-owned contribution identities, use entries and representation nodes, open generics, a typed facet contract, and a prototype before the migration. The work items are #196, #195, #207 and #197. The design page is [the contribution contract](../docs/contribution-contract.md). The reports of its three verification agents are [1](0017-contribution-contract.verify-1.md), [2](0017-contribution-contract.verify-2.md) and [3](0017-contribution-contract.verify-3.md). |

There is no implication of a missing current feature merely because a historical
submission names a different API. Conversely, later working examples do not close
an explicitly open architectural defect or prove native constraint preservation.
Follow [generation gaps](../docs/generation-gaps.md) and the
[owning case index](../docs/architecture-cases/README.md) for current status.
