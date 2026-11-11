# Documentation map

This map identifies the owner of each current contract. A historical proposal or
consultation does not override a current contract guide. Status below describes
the implemented scope, not complete language coverage.

| Question | Owning documentation | Status |
| --- | --- | --- |
| What belongs to Bittype? | [Charter](../CHARTER.md), [model boundary](model.md), [dependencies](dependencies.md) | Current scope and dependency rules |
| Where does code belong? | [Layout](../LAYOUT.md) | Current stage/dimension/scope/component hierarchy |
| What exists at each abstraction level? | [Contexts and typed relations](abstraction-levels.md) | Conceptual model, with a current-carrier correspondence table |
| How do AST, CST, Files and Source relate? | [AST](ast.md), [text stages](parsing-formatting.md), [Files/Source](source.md) | Implemented carriers and role contracts; coverage varies by language |
| What does the semantic graph establish? | [Semantics](semantics.md), [resolution](resolution.md), [analysis](analysis.md) | Implemented contracts; language-owned judgments |
| What is canonical Horos input? | [Model](model.md), [Horos syntax](horos-syntax.md), [kinds](kinds.md), [frontend](horos.md) | Revision-4 frontend using public Horos 0.10.2 |
| How is generation composed? | [Overview](generation.md), [registration](kernel-composition.md), [SPI](generation-spi.md) | Declarative planning and separate native construction |
| Who chooses native representations? | [Type-owned mappings](type-generation-plugins.md), [generic elaboration](generic-elaboration.md), [representation profile](semantic-generation.md) | Four native pagination mappings; enforcing profile unsupported |
| How does a consumer generate Page/Cell? | [Runnable fixture walkthrough](page-cell-example.md), [planner](generation-planner.md) | Existing tested TypeScript composition |
| When are files and imports chosen? | [Generation/Files boundary](source.md#from-generated-declarations-to-files), [TypeScript imports](typescript-type-imports.md) | Generated Files handoff implemented for TypeScript |
| How should other targets place declarations? | [Go/Rust/C++ placement](native-file-placement.md) | Proposal; no native Files formatter for these targets |
| Which native semantic facts exist? | [Native semantics](native-semantics.md) | Bounded Go, TypeScript and Rust; C++ deferred |
| What does verification accept? | [Joint verification](verification.md), [verification architecture](type-system-architecture.md), [Logos contract map](logos-verification.md) | Public association, typed interpretations and scoped Horos/Go laws |
| What remains for native constraints? | [GEN-02](generation-gaps.md#gen-02), [constraints](native-kind-constraints.md), [Go admission](go-native-admission.md) | Proposed encodings; Go compiler fixtures are not production emission |
| What broader reasoning is proposed? | [Type-system model](type-system-model.md), [preservation reasoning](logos-preservation.md) | Language-owned theories; no new Bittype proof foundation |
| What belongs above static types? | [Binding generation](binding-generation.md), [representation families](representation-families.md), [runtime values](runtime-values.md), [optional extensions](hydrated-types.md), [concept inventory](wire-extension-preservation.md) | Accepted ownership of binding generation (research 0014) and the family contract (research 0015); the selected plan, families and adapters are not implemented |
| What is implemented and what has evidence? | [Provider matrix](providers.md), [language overview](language-support.md), [capability reports](capability-reports.md) | Generated source-derived reports; execution evidence is separate |
| What architectural corrections remain? | [Generation gaps](generation-gaps.md), [case index](architecture-cases/README.md) | GEN-01's correction is independently verified; Rust, C++ and Go Files ownership remains open (ACX-004); the structural profile gate is open (ACX-008); GEN-02 is an unsupported guarantee |

The exact implemented SPI belongs in its contract guide, with its source link.
Conceptual guides explain responsibilities and refer to that definition rather
than maintaining a second interface. Examples must identify omitted setup or
point to their complete runnable fixture. A proposal must state what remains
unimplemented and must not advertise a capability merely because its syntax exists.

## Research and historical decisions

[Research records](../research-docs/README.md) retain the questions, original
submissions, expert advice and dated evaluations. Their code baselines and
recommendation dispositions describe those consultations. Later decisions can
replace their suggested interfaces while preserving their substantive concerns.
The current guides above own today's contracts. Original submissions and advice
are evidence, not active specifications to rewrite as if they described today.

Update a contract's guide, its affected examples/status statements and the
[changelog](../CHANGELOG.md) together. Update the owning architectural case when
its correction or qualification changes. This is documentation maintenance, not
a new approval, provenance or heuristic checking workflow.
