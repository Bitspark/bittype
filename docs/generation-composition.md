# Automatic generation composition

The kernel selects registered descriptions, closes their dependencies into a
semantic plan, and constructs that plan's native ASTs. Consumers register
libraries once and request a source capture, roots, target and generation
profile. [Registration SPI](kernel-composition.md#registration-and-request) ·
[Mapping SPI](generation-spi.md)

## Ownership and dependencies

| Owner | Responsibility |
| --- | --- |
| Source interpretation and reader | Classify declarations, name descriptors and judge uses from supplied evidence; interpret payloads for its mapping libraries |
| Kernel walk | Admit the graph, resolve roots, open declarations, follow references and applications, check child captures and binder scope |
| Mapping library | Declare exact source patterns, chosen native forms, declaration/helper recipes and generic conventions |
| Kernel planner | Selection, replacement, scoped substitution, dependency closure, limits and whole-request failure |
| Native model | Native form identity and validity, without AST syntax |
| Native target and shared constructor | Realize those forms with fresh syntax, local binder links and requested exposures |
| Target Files formatter | Placement, module syntax and checking actual bindings |
| Verification composition | Separate judgments over captured source and target artifacts |

## Construction and selection

<!-- bittype:overview scope=construction-and-selection part=features -->
A Selector combines a source language/profile/Go marker with an opaque key.
Only the source owner interprets that key. The registry freezes descriptions
at registration. The planner selects exact identities and explicit replacements
before expanding a declaration body. Builtin and independent library mappings
use the same public contract.
<!-- bittype:overview:end -->

<!-- bittype:overview scope=construction-and-selection part=exclusions -->
A missing or ambiguous demanded representation fails. Exact mapping errors do
not invoke a structural fallback. Source admission remains required for all
actuals, including those the native convention omits. Mapping callbacks run
only during planning, and cannot access native construction state.
<!-- bittype:overview:end -->

[Registry implementation](../transform/syntax/shared/generation/go/registry.go) ·
[Guarded planner](../contracts/generation/go/planning.go)

## Demands, generics and final output

Requested roots need declaration recipes; an external use rule cannot silently
drop a requested publication. Dependencies and helpers are emitted once, without
automatically becoming public. A selected external Cell representation can keep
its Horos Wire body opaque. The source analyzer still admits that original body.

One GenericConvention relates source binders, native parameters and complete
application arguments. The planner uses that convention for every ordinary
application. It does not specialize Page from its first use. The independent
Go fixture adds a native parameter computed from the complete mapped actual and
composes it through nested applications.

GenerationPlan contains no source reader or mapping expansion. Its Run constructs
fresh syntax from captured target expressions. Native symbols and fragments are
bound to that run. Scope, foreign-reference, cancellation and budget failures
produce no output. Current native profiles reject cyclic definition dependencies.
[Plan/result contracts](generation-spi.md#plan-and-output)

## Extension examples and qualification

The [planning tests](../integration/go/generation_planning_test.go) demonstrate
that planning performs no AST construction, repeated runs do not repeat mapping,
registration changes cannot alter a captured plan, and foreign slots or child
captures fail. An independently supplied two-argument Pair reverses its native
arguments inside Page. The same TypeScript target and shared planner handle it.

The [Horos composition tests](../integration/go/horos_generation_composition_test.go)
replace list with ReadonlyArray and select a package-owned Cell. The
[independent source fixture](../integration/go/generation_independent_test.go)
uses another source language, an extra native parameter and a native-only helper.

The existing four-language pagination, compiler and
[generated Go verification](../integration/go/horos_generation_verification_test.go)
tests remain. The [TypeScript placement tests](../integration/go/typescript_placement_test.go)
reuse one output with Person separated or colocated and a conflicting local Cell.
They check both local references and the actual external package identity.

Generation does not establish native kind enforcement or source/native
equivalence. Those retain their separate profiles and verification contracts.
C++ includes and Rust visibility remain the documented
[open Files boundary](architecture-cases/0004-generation-file-syntax.md).
