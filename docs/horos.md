# Horos declaration AST and semantics

**Revision 4 implemented in Go source, 4 November 2026.** One Horos frontend
projects general kind definitions and type definitions, including constructor
results. It consumes public Horos Go 0.10.2. Horos owns the
[canonical format](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md).
Bittype owns its [readable spelling](horos-syntax.md), AST projection and provider
composition. The [kind design](kinds.md) supplies the semantic overview.

## Packages and public helpers

- `syntax/horos/go` defines the language marker, AST/CST vocabulary and identity.
  It depends on generic carriers and identity, not Horos's semantic implementation.
- `check/syntax/horos/provider/go` admits syntax: node forms, roles, cardinalities and lexical
  spellings. Unresolved names and semantically invalid contracts can reach
  resolution; syntax admission does not establish membership or module meaning.
- `resolve/files/horos/provider/go` delegates meaning to Horos and projects the checked result.
  Its public helpers remain `Resolve(root, dependencies, treeID)`,
  `Project(program, treeID)` and `Model(view)`. `Model` is an explicitly semantic
  single-view convenience, not a generic carrier operation.
- The registered `ResolveReferences` provider accepts the captured collection and
  explicit entry. Its private `internal/horosmodel.DecodeCollection` path retains
  several Program trees and tree-qualified sites. It does not concatenate source
  or expose a second public model wrapper.
- `analyze/files/horos/provider/go` consumes the explicit compatible resolution result and
  preserves its established typed Semantics and required reports.
- `parse/syntax/horos/provider/go` builds a lossless CST with AST view; `format/syntax/horos/provider/go` builds
  CST directly from an AST view. The shared printer emits leaves. Neither text
  stage resolves or applies type-level values.
- The consumer registers Horos syntax, checking, resolution and analysis through their own stage providers. The language ID remains `horos`, with display name
  `Horos`. Implementing a capability installs nothing in an engine.

The new profiles are:

| Boundary | Profile |
| --- | --- |
| AST syntax | `horos-declarations-4/ast-1` |
| Resolution | `horos-declarations-4/resolution-1` |
| Analysis | `horos-declarations-4/analysis-1` |

Applications register the language descriptor and selected check/resolve/analyze
and Syntax stages explicitly. The existing symmetric role contracts remain
`parse.Provider[parse.Input, cst.CST[Language]]`,
`format.Provider[ast.View[Language], cst.CST[Language]]` and
`print.Provider[cst.CST[Language], string]`, with context passed to each operation.
[Files composition](source.md) uses the corresponding Files contracts. Old
profile readers and aliases are removed rather than negotiated implicitly.

Only the concrete Horos stages, their private model bridge and source-specific
transformation reader consume Horos/Ontos dependencies. Generic trees, contracts,
engine and printer, and native Syntax stages, remain independent. Model dependency
alignment belongs to the [dependency guide](dependencies.md).

## Syntax vocabulary

<!-- bittype:overview scope=syntax-vocabulary part=features -->
The canonical module envelope remains revision, name, imports and declarations,
with revision 4. There is no schema-digest field. The AST retains Program roots,
module names/pins, imports, ordered parameters and literal Atom/Tuple fragments.
Module declarations now distinguish kind and type definitions.
<!-- bittype:overview:end -->

| Construct | AST vocabulary and purpose |
| --- | --- |
| Type definition | `KindTypeDefinition` (`type-definition`), with parameter children, result kind and body |
| Kind definition | `KindKindDefinition` (`kind-definition`), with parameter children, result sort and body |
| Type binder | `KindTypeParameter`, name/index and `RoleContract` kind |
| Kind binder | `KindKindParameter`, name/index and `RoleSort` sort |
| Sort | `KindProperSort` or `KindArrowSort` |
| Top/bottom | `KindTop`, `KindBottom`, explicitly at a sort |
| Kind binder use/reference/application | `KindKindParameterReference`, `KindKindReference`, `KindKindApplication` |
| Boolean/constructor kinds | `KindUnion`, `KindIntersection`, `KindArrow` |
| Structural kind | `KindNodePattern`, with selected metadata and child-sequence pattern |
| Recursive kind | `KindRecursiveKind`; `RoleEquation` children contain `KindEquation` bodies |
| Recursive use | `KindRecursiveVariable`, with lexical depth and index |
| Child sequence | `KindSequenceEpsilon`, `KindSequenceItem`, `KindSequenceConcat`, `KindSequenceAlternative`, `KindSequenceStar` |
| Type reference | `KindReference`, separate from an application |
| Type application | `KindApplication`, with `RoleCallee` and one or more `RoleArgument` children |

Both definition forms retain `RoleParameter`, `RoleResult` and `RoleBody`.
A recursive group's entry is its `AttributeIndex`; equation indices determine
canonical slots. Optional equation names support readable output, while lexical
depth/slot determines binding. The formatter freshens names only when needed to
preserve that binding.

<!-- bittype:overview scope=syntax-vocabulary part=exclusions -->
The old generic Declaration/Bound nodes and RoleBound are removed. Type/Data
names are ordinary references and never choose an AST variant. Primitive top
and sort syntax supplies the unconstrained proper default. Descriptor operands
still preserve the exact canonical value layout. Only Horos-designated type
positions contain type expressions; labels, scalar selectors, field wrappers and
`exact` payloads retain their literal roles.
<!-- bittype:overview:end -->

Parameter uses, local/imported references and expression callees retain separate
bindings. A constructor-valued declaration may have no declared parameters, so
its semantic sort cannot be derived from `Entity.Parameters` or AST parameter
count. Supplied binding annotations assert a target; resolution compares them
with Horos's established binding and does not repair them by mutating the AST.

## Multi-Program and standard-module behavior

The selected root and its transitive import closure determine the semantic
domain. Multiple supplied Program trees retain distinct tree-qualified addresses;
filenames and storage order do not establish visibility. Unrelated supplied
modules are not silently promoted into the selected closure. The existing Files
profile keeps its complete supplied-file admission rules.

The exact bundled standard module is identified by public `resolve.PreludePin()`.
A source module still writes an ordinary import to that identity. When the
bundled prelude has no supplied AST, authoritative kind entities may be
source-less. The bridge invents neither AST locations nor a Files relation to a
nonexistent prelude file. Supplied prelude source corresponds normally.
`Project` can materialize prelude syntax as part of a standalone projection;
Files formatting preserves the supplied file set rather than adding that file
implicitly.

Canonical identity is separate from local entity IDs and AST addresses. Origins
remain tooling metadata when available; they do not authorize library behavior
or stand in for semantic bindings. Literal payloads remain opaque throughout
canonical reconstruction and projection.

## Typed Semantics and provider-owned facts

The existing `Semantics[horos.Language]` and sealed `representation.Semantics`
envelope remain unchanged. Unpacking requires the language ID, explicit profile
and actual Go marker. There is no `Semantics[any]` conversion. The generic graph
stores structured provider terms and multivalued subject relations without
interpreting Horos kinds.

| Public relation helper | Horos-owned relation name | Meaning |
| --- | --- | --- |
| `SortRelation()` | `sort` | Sort of the subject's type-level value |
| `ResultSortRelation()` | `result-sort` | Result sort of a definition |
| `DeclaredKindRelation()` | `declared-kind` | The authored kind contract |
| `ResultKindRelation()` | `result-kind` | The declared result kind |
| `SymbolicBodyRelation()` | `symbolic-body` | Checked source body retaining binding structure |
| `MembershipRelation()` | `membership` | An individually established membership with assumptions |
| `IdentityRelation()` | `canonical-identity` | Canonical declaration/parameter identity |

Each relation has owner `horos`. Entities distinguish `type-definition`,
`kind-definition`, `type-parameter` and `kind-parameter`. `Entity.Body` retains
the checked symbolic source body. Kind expressions and sorts remain structured
provider-owned terms; the kernel acquires no Data, Scalar, List, Record or Wire
classification rule.

Membership records come from authoritative source semantics and retain exact
ParameterID assumptions. Missing membership facts assert nothing; several facts
may concern one subject. Arbitrary later queries do not mutate a captured model.
There is no principal-kind fact, and completeness does not mean enumeration of
all true memberships. In particular, a narrower result query must not become an
unconditional membership on an open body with weaker original assumptions.

Generic queries remain available without importing Horos: `Resolve`, `TermAt`,
`DirectParameterBindings` and `DependenciesFrom`. DirectParameterBindings pairs
actual formals only for a direct declaration reference with matching declared
arity. It does not fabricate parameters for a constructor-valued result or infer
bindings for an indirect callee. Dependency relations cover named semantic
references, including unapplied or ignored references, rather than applications
alone.

## Reports and resource outcomes

Required reports include program validity, declaration enumeration, occurrence
bindings, applications, reference dependencies and type-fact coverage. The
resolver also establishes reference resolution. The provider vocabulary uses
`reference-dependencies` for dependency coverage, `declaration-reference` for an
edge, and `type-facts` for required fact coverage. Analysis preserves the resolver's
required reports through its explicit receipt.

Successful validity/reference reports establish Holds. Completed enumeration
reports describe coverage without asserting a new language proposition. Missing
context, unproved required obligations and exhausted operations cannot publish
a usable partial semantic model. Export exhaustion can preserve an earlier
validity report while leaving graph publication incomplete. Unknown requested
obligations remain Unsupported rather than silently counted as complete.

`resolve.LimitsSetting()` retains provider-owned configuration for
`[MaxDepth, MaxNodes, MaxBytes, MaxModules]`. Absence uses Horos defaults; malformed
configuration is a request/configuration error. The kernel bounds and copies the
configuration without interpreting it. Canonical, syntax/carrier, graph and
expanded-output budgets remain distinct. Horos now charges general kind graphs,
sequence automata, environments and query work through its bounded operations.

## Horos value/query API

The projected Program supplies the general sealed API:

```text
Value(declarationID)                 -> TypeValue
Proper(closedDescriptor)            -> TypeValue
Apply(callee, []TypeValue)           -> TypeValue
CompileKind(moduleDigest, kindExpr) -> KindHandle
Satisfies(value, kind)               -> Verdict
Includes(actual, expected)           -> Verdict
Descriptor(value)                   -> materialized descriptor
GroundDescriptor(value)             -> admitted ground descriptor
```

These operations also return errors according to the public Horos API. Verdicts
are Proven, Refuted and Unproved; Invalid and Exhausted are distinct error
outcomes. Handles belong to one checked Program. `Apply` may return a constructor;
`Descriptor` and `GroundDescriptor` require the relevant closed proper form.
GroundDescriptor additionally establishes ordinary ground-kind membership and
unchanged ground well-formedness. Specialize/DataDescriptor compatibility wrappers
are not retained.

The frontend exports source semantics, not runtime values. Native pagination is
a separate [representation mapping](semantic-generation.md). Function, Stream
and Cell remain ordinary libraries; no query result selects their behavior,
transport, hydration or native representation.
