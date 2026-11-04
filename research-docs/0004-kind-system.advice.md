# Advice for consult
**Research doc:** doc_efdb21fa-1de0-443d-974d-788e2b2871a5.md
**Document:** doc_efdb21fa-1de0-443d-974d-788e2b2871a5
**Advisor:** ChatGPT, model "6 Pro" as labelled on the answer page
**Session:** nsess_680318ab-bee4-4ce3-9a26-af21b616067b
**Run:** run_03dd46c7-4356-4db4-92cf-2e2fd3f98645
**Chat:** https://chatgpt.com/c/6aeb4438-2b8c-83ed-96dd-281f755e8511
**Verify:** nightfall consult verify run_03dd46c7-4356-4db4-92cf-2e2fd3f98645 --output <this file>

## Recommendation

Implement **explicit, saturated higher-order kinds with a separate, compositional Data proof system**. Preserve actual constructor identities and symbolic result summaries when constructors pass through parameters.

The indispensable distinction is between:

- **Admission:** what arguments a constructor promises to accept.
- **Result evidence:** what can be proved about its result, first under abstract parameter contracts and later with known arguments.

This supports every required `Wrapped`, `Compose`, `Page`, `Ignore`, and `DataWrapped` judgment without requiring partial application, type lambdas, or kind polymorphism. The following is a proposed replacement against the supplied implementation snapshot.

## 1. Core calculus

### Kinds and interfaces

Use finite, explicitly annotated kinds:

\[
K ::= \mathsf{Type}
\mid (K_1,\ldots,K_n)\Rightarrow\mathsf{Type}
\qquad n\geq1.
\]

These are **n-ary, fully saturated constructors**, not curried arrows. Domains may contain further constructor kinds. Every application and declaration body produces a proper type.

Separately define parameter interfaces:

\[
S ::= \mathsf{Type}\mid\mathsf{Data}
\mid (S_1,\ldots,S_n)\Rightarrow R,
\qquad R\in\{\mathsf{Type},\mathsf{Data}\}.
\]

Erasing `Data` to `Type`, recursively, gives an interface’s shape kind. Thus `Data` is a refinement of proper types; it does not introduce another universe or a value type. Canonical descriptions being Ontos data does not establish that the types they describe satisfy `Data`.

Examples:

```text
Page:
    kind       (Type) => Type
    admission  (Type) => Type

DataPage:
    kind       (Type) => Type
    admission  (Data) => Type
    result     provably Data

Wrapped:
    kind       ((Type) => Type, Type) => Type
```

The unary source spelling `Type -> Type` can remain convenient. For multiple inputs, specify tuple-domain notation; do not silently interpret full application as currying.

### Judgments and application

Use three distinct judgments:

```text
Γ ⊢ e : K       e is a valid expression of shape kind K
Γ ⊢ e ⇐ S       e satisfies interface S
Γ ⊢ Data(e)     e has an established Data proof
```

The application rule is:

\[
\frac{
c\text{ has admission }(S_1,\ldots,S_n)\Rightarrow R
\quad
\Gamma\vdash a_i\Leftarrow S_i\ \text{for every }i
}{
\Gamma\vdash c\langle a_1,\ldots,a_n\rangle:\mathsf{Type}
}
\]

Require exact arity, including every argument even if unused. Derive `Data` separately from the declared result guarantee or the known constructor’s checked summary.

References and variables obey these rules:

- A parameter has its declared interface and erased kind.
- A declaration with parameters denotes a constructor when referenced unapplied.
- A declaration without parameters denotes a proper type.
- Every descriptor child designated as a type position must have kind `Type`.
- Every declaration body must have kind `Type`; a declared Data result additionally requires a proof under its parameter interfaces.

Consequently, `Wrapped<Person,Person>` fails at the first argument; `wire<Page>` fails because `Page` is a constructor; `Wrapped<Page,Person>` has the correct shapes.

Parameter annotations are explicit. An omitted ordinary annotation elaborates to `Type`; infer checked result summaries, not missing higher-kinded parameter annotations.

### First implementation boundary

Include constructor references, constructor parameters, higher-order input kinds, parameter callees, nested saturated applications, and verified Data interfaces.

Defer:

- Partial application.
- Constructor-valued declaration results.
- Bare constructor aliases.
- Source type lambdas.
- Kind variables and rank-1 kind polymorphism.

Explicitly reject:

```text
Return<F: Type -> Type> = F;    // body has constructor kind
P = Page;                      // constructor-valued alias
Compose<Page>;                 // insufficient arguments
Return<Page><Person>;          // Return is not a valid declaration
```

An explicitly parameterized forwarding declaration such as `P<T> = Page<T>` remains an ordinary distinct declaration. The implementation must not manufacture these as lifting machinery.

Higher-kinded parameters do not inherently require kind polymorphism, and saturation is a separate language decision; GHC’s documentation makes both distinctions explicit. [ghc.gitlab.haskell.org](https://ghc.gitlab.haskell.org/ghc/doc/users_guide/exts/poly_kinds.html?utm_source=chatgpt.com)

## 2. Data precision and constructor compatibility

### Structural interface compatibility

Let \(S\preceq T\) mean that anything satisfying interface \(S\) can satisfy interface \(T\).

Besides reflexivity:

\[
\mathsf{Data}\preceq\mathsf{Type}.
\]

For matching arities:

\[
\frac{
B_i\preceq A_i\text{ for every }i
\quad R_A\preceq R_B
}{
(A_1,\ldots,A_n)\Rightarrow R_A
\preceq
(B_1,\ldots,B_n)\Rightarrow R_B
}.
\]

Inputs are contravariant; results are covariant. This concerns constructor substitutability, not native container variance. FreeST documents the same arrow polarity, although its other semantics are outside this proposal. [FreeST](https://freest-lang.github.io/tutorial/kinds-and-types/?utm_source=chatgpt.com)

Therefore:

| Interface comparison | Result |
|---|---|
| `Type -> Data` satisfies `Data -> Type` | Yes |
| `Data -> Data` satisfies `Type -> Type` | No: insufficient input admission |
| `Type -> Type` satisfies `Data -> Data` | No: insufficient result guarantee |

The last row is why arrow compatibility alone cannot establish Page’s preservation property.

### A checked declaration can establish stronger result interfaces

Suppose a named declaration `d` has original input interfaces \(B_1,\ldots,B_n\). To check it against an expected interface \((A_1,\ldots,A_n)\Rightarrow R\):

1. Require matching shape and arity.
2. Require \(A_i\preceq B_i\) for every input.
3. Introduce fresh rigid arguments \(x_i:A_i\).
4. If \(R=\mathsf{Data}\), evaluate `d`’s already-checked symbolic result summary on those arguments and require a Data proof.

This verifies a universally usable interface using symbolic inputs, not sampled concrete instantiations.

For an **abstract constructor parameter**, use only its declared interface and structural compatibility. Its future implementation is unavailable while checking the generic body.

Crucially, step 2 preserves declared admission even when the body ignores an argument:

```text
type OnlyData<D: Data> = text;
Wrapped<OnlyData, Person>       // reject
```

`OnlyData` cannot become an arbitrary-Type constructor merely because its current result is constant.

### Data rules and summaries

Preserve the existing primitive rules, including their other well-formedness obligations:

| Expression category | Data rule |
|---|---|
| Data primitives, `any`, valid `exact` literals | Proved Data |
| `wire<M>` | Not Data; still validate `M` as a proper type |
| Tuple, compound, list, closed record, union | Conjunction of designated value-bearing child proofs |
| Existing data-only constructors | Require their child Data obligations before admitting the expression |
| Proper parameter `T:Type` | Symbolic requirement `Need(T)` |
| Proper parameter `D:Data` | Proved Data |
| Abstract constructor application | Retain a symbolic application; use its declared guarantee for the current proof |
| Known constructor application | Substitute checked arguments into its stored result summary |

A suitable internal summary DAG is:

```text
Q ::= true
    | false
    | Need(ParameterID)
    | And(Q...)
    | Probe(constructor, checked-expression-arguments)
```

`false` represents failure of the prescribed compositional Data criterion. An unresolved `Probe` represents information still dependent on a constructor. **Do not encode an abstract call as the existing `impossible` flag.**

`Probe` arguments must retain typed expression structure and constructor identities, not merely Data booleans.

For the supplied definitions:

```text
δ(Page<T>)          = δ(T)
δ(Ignore<T>)        = true
δ(Send<T>)          = false
δ(Wrapped<F,T>)     = Probe(F, [T])
δ(Compose<F,G,T>)   = Probe(F, [G<T>])
```

Evaluate these under an environment without overwriting the generic summary with a specialization’s answer.

### Required derivations

For `Page ⇐ Data -> Data`, introduce `x:Data`. Its original `Type` input accepts `x`; its summary becomes `δ(x)=true`. This is an additional checked interface of the **same Page**.

For `Ignore ⇐ Type -> Data`, introduce `x:Type`; the result summary is unconditionally true.

For `Send ⇐ Data -> Data`, the input check succeeds but the result remains `wire<x>`, so the guarantee fails.

The resulting application judgments are:

| Application | Judgment and reason |
|---|---|
| `Wrapped<Page,Person>` | Data: substitute Page, then prove `Data(Person)` |
| `Wrapped<Ignore,wire<Person>>` | Data: Ignore’s result ignores its admitted argument |
| `Wrapped<Send,Person>` | Valid Type, not Data |
| `Compose<Page,Page,Person>` | Data: apply Page’s conditional summary twice |
| `DataWrapped<Page,Person>` | Data: Page satisfies the required preservation interface |
| `DataWrapped<Send,Person>` | Reject the constructor argument |
| `DataWrapped<Ignore,wire<Person>>` | Reject the second argument despite Ignore’s pure result |
| `Wrapped<DataPage,Person>` | Reject: DataPage cannot accept arbitrary Type inputs |

Open-body checking must precede specialization:

```text
Bad<F: Type -> Type, D: Data> = DataPage<F<D>>;
```

Reject: `F` could be `Send`, so `Data(F<D>)` is not established.

```text
Good<F: Data -> Data, D: Data> = DataPage<F<D>>;
```

Accept: `D` satisfies the input contract and `F` guarantees a Data result.

The converse failure also matters:

```text
BadInput<F: Data -> Data> = F<wire<Person>>;
```

Reject even if every proposed caller supplies Page. A known caller cannot repair an invalid open definition.

### Deliberate limitation and alternatives

This first interface language cannot express both “accept every Type” and “preserve Data inputs” for the **same abstract parameter**. `F:Data -> Data` permits only Data calls; `F:Type -> Type` guarantees no Data result. Known Page supports both independently verified uses.

A later finite conditional guarantee or interface intersection can expose both properties. This extension does not require kind polymorphism.

Type/Data subkinds are another presentation of much of this small system, provided they retain body-derived preservation evidence. Bounded kind polymorphism could express reusable schemes, but adds quantified binders and instantiation without being necessary for these examples. Type lambdas offer greater uniformity at a larger binding and normalization boundary.

Do not cite richer higher-order subtyping undecidability as an objection to this finite structural comparison: the type-interval calculus in the cited paper includes substantially stronger mechanisms. [arxiv.org](https://arxiv.org/abs/2107.01883?utm_source=chatgpt.com)

## 3. Canonical representation and bounded checking

### Concrete revisioned grammar

A possible revision 3 encoding, using the document’s notation, is:

```text
Module =
  D("module", I(3), T(name), importsMap, declarationsMap)

Claim = A("type") | A("data")

Interface =
  D("proper", Claim)
  | D("constructor", L(Interface, ...), Claim)
    // nonempty input list; result always proper

Parameter =
  D("parameter", T(name), Interface)

Declaration =
  D("type", L(Parameter, ...), Claim, E)

E =
  Horos descriptor with E only in schema-designated type positions
  | D("param", I(slot))
  | D("local", T(name))
  | D("imported", T(alias), T(name))
  | D("apply", E, L(E, ...))
```

Store interfaces once and derive their erased kinds; avoid redundant independently authoritative kind fields. An omitted source result claim elaborates to `type`; stronger inferred summaries remain checked semantic facts.

An application has a genuine expression callee. In this first calculus, a callee that is itself an application fails kind checking because applications produce proper types. The representation remains principled without claiming constructor-valued results are supported.

Reject empty application lists: a zero-parameter declaration is referenced directly.

### Binding and literal opacity

Keep one ordered parameter scope per declaration:

```text
DeclarationID = (exact module digest, declaration name)
ParameterID   = (DeclarationID, slot)
```

Interface trees introduce no named binders or dependencies. Fresh symbolic arguments used for interface checking are internal proof variables, not fabricated declarations.

A named constructor retains its declaration identity and definition environment. Its internal local/imported references must never be re-resolved in the caller’s module. Substitute by binder identity, not spelling.

Traverse only Horos’s schema-designated expression positions. For example:

```text
H("exact", D("param", I(0)))
```

contains no generic occurrence. Apply the same rule to literal tuples resembling `apply`, constructor-reference, Wire, or interface tags. Preserve literal bytes, labels, presence markers, ordering, and metadata exactly.

Keep expression-form tags and semantic kinds as distinct Go types.

### Checking algorithm

Use these stages, returning no checked Program if any required stage fails:

1. **Admit and decode.** Validate revision, tuple shapes, interface structure, primitive payloads, sizes, and supplied import pins/names. Perform no implicit IO.
2. **Resolve.** Bind every semantic reference and parameter. Collect dependency edges from **all named references**, including unapplied constructors and ignored arguments.
3. **Reject declaration cycles.** Then check declarations in dependency order. Ground `rec`/`ref` remains a separate mechanism.
4. **Check every open body.** Validate every argument and primitive obligation under declared interfaces. Validate explicit result claims. Build shared symbolic summaries.
5. **Answer applications and interface queries.** Use structural interface comparison and bounded summary evaluation with actual argument environments.

The expanded dependency rule is essential:

```text
type IgnoreF<F: Type -> Type, T> = text;
type Loop<T> = IgnoreF<Loop,T>;
```

Reject this cycle even though `Loop` contributes nothing to the result.

Cache summary evaluations by expression/summary identity plus substitution environment. Cache interface checks by exact declaration identity and expected interface. Include relevant dependency/configuration identity; erased kinds or Data flags alone cannot distinguish Page, Ignore, and Send.

### Specialization and equality

Replace the untyped constructor-argument path with a deliberately typed API: validated proper-type arguments and checked constructor handles carrying declaration/dependency ownership.

For a closed application:

- Validate all arguments against original admission contracts.
- Substitute closures using definition scope and stable binder identities.
- Expand named applications while retaining provenance separately.
- Revalidate the closed result.
- Return a ground descriptor through `DataDescriptor` only after establishing Data.

Keep open declarations immutable. Internal summary closures do not require source lambdas.

Maintain three separate notions:

1. **Canonical identity:** exact module bytes and declaration names.
2. **Kind/interface comparison:** structural equality and the specified refinement relation.
3. **Expansion:** named substitution used for result proofs or descriptor materialization.

No general extensional constructor equality, eta equality, or value-type equivalence solver is required. Equal expanded descriptors do not permit merging declaration identities or substituting behavioral libraries. Renaming a binder may preserve meaning while changing canonical bytes and pins.

### Resource outcomes and migration

Retain the quoted limits—default depth 256, configured maximum 512, 100,000 nodes/work steps, 16 MiB, and 256 supplied modules—and extend accounting to interfaces, comparisons, environments, summaries, and normalization.

Maintain separate admission, semantic-work, and output budgets. Charge aggregate work across the operation, not a fresh allowance per declaration. Count emitted descriptor occurrences and encoded bytes even when their internal representation is shared.

For example:

```text
Dup<T> = tuple<T,T>;
D0 = text;
D1 = Dup<D0>;
...
Dn = Dup<D(n-1)>;
```

Data analysis can succeed cheaply while output materialization exceeds its budget.

Distinguish **valid**, **invalid**, **incomplete due to limits**, and **valid but unsupported by a mapping**. Exhaustion establishes neither invalidity nor Data. Failed specialization returns no usable partial result.

Release revision 3 as a clean replacement. Reissue migrated modules with new bytes and pins; repin dependents. Preserve revision 2 artifacts and evidence through their immutable releases. The replacement reader rejects revision 2 rather than silently translating it. Existing pagination source and native behavior should remain equivalent, while canonical hashes necessarily change.

## 4. Bittype facts and native generation

### Agnostic graph changes

Revise application structure to:

```text
Reference   { Target: EntityID }
Application { Callee: TermID, Arguments: []TermID }
```

Remove application validation against `Entities[Target].Parameters`. Validate callee/argument edges, active fields, ownership, containment acyclicity, and budgets. Horos decides applicability, arity, interfaces, and saturation.

Constructor parameters remain real parameter entities. Their kinds must not be simulated through fabricated declaration parameters.

Expose structured provider-owned facts for:

- Erased kind and declared interface.
- Established Data facts and conditional summaries.
- Canonical declaration/parameter identity.
- Source and reduction provenance.

A generic subject-to-fact-term relation is sufficient. Kind/interface trees can use `ProviderExpression` operators with explicit operand edges. An opaque sort key alone is insufficient; ordinary applications must remain visible graph structure.

The kernel acquires no Type/Data rules, primitive table, library recognition, or provider invocation.

Bittype’s readable syntax gains interface syntax, bare constructor references, and expression callees. Its private Horos bridge validates and projects authoritative results. Syntax admission must not perform semantic resolution.

Any explicit source reduction creates a new snapshot requiring fresh resolution and analysis. Provenance back to an old AST does not validate an old receipt against the new AST.

### Incremental native policy

First preserve the existing ordinary pagination profile in all four languages. Until additional mapping support is qualified, reject residual higher-kinded forms explicitly.

| Target | Direct policy for residual higher kinds | Qualified extension |
|---|---|---|
| **TypeScript** | Reject constructor parameters and `F<T>` through them. The general feature remains an open proposal. | An explicit alternative representation could be added separately; do not introduce it implicitly. [GitHub](https://github.com/microsoft/TypeScript/issues/1213?utm_source=chatgpt.com) |
| **Go** | Reject; generic types must be instantiated when used. | Use an explicitly selected source-specialization route. [The Go Programming Language](https://go.dev/ref/spec?utm_source=chatgpt.com) |
| **Rust** | Reject direct constructor-parameter mapping. | An opt-in family-marker/GAT representation requires traits, projections, bounds, and applicable `where` clauses in the AST and mapping. [The Rust Reference](https://doc.rust-lang.org/reference/items/associated-items.html?utm_source=chatgpt.com) |
| **C++** | Reject until nested template-parameter AST/checker/formatter support exists. | Then qualify fixed-arity template-template parameters over ordinary type slots, excluding packs/defaults/value parameters initially. Native matching remains a separate check. [eel.is](https://eel.is/c%2B%2Bdraft/temp.param?utm_source=chatgpt.com) |

A narrow C++ extension could emit:

```cpp
template<template<class> class F, class T>
struct Wrapped {
    F<T> value;
};
```

Do not generalize that qualification automatically to arbitrary nested higher-order interfaces.

### Explicit source specialization

Offer a separate, explicitly selected constructor-binding transformation with named roots and bindings. For example, binding `Wrapped`’s `F` to Page can produce:

```text
type WrappedPage<T> = { value: Page<T>; };
```

This retains an open ordinary generic and references the one Page declaration. It does not copy Page for each caller. This transformation is distinct from allowing partial application as a source expression.

The residual document must be checked anew and satisfy the target’s complete supported profile. Existing restrictions—including Wire, unnamed nested records, unsupported descriptors, and unused parameters—still apply.

Default emitters must not silently perform this reduction even for a closed `Wrapped<Page,Person>`. Nor may they silently drop unsupported exports: selecting output roots is an explicit policy.

Native compilation does not prove Horos Data conformance for arbitrary native instantiations. No universal native `any`, runtime adapter, raw-source escape, or library-name dispatch follows from this proposal.

## 5. Missing obligations and independent acceptance criteria

### Data preservation does not imply guardedness preservation

This is an important additional failure case:

```text
type Id<T> = T;
```

Id preserves Data, but substituting it into a putative recursive body `rec X. F<X>` may expose an unguarded reference. A `Data -> Data` guarantee cannot establish that an application supplies a recursion guard.

Preserve the existing requirement that application arguments be recursion-closed, retain the separate guardedness checker, and recheck expanded closed descriptors. Broader interaction between constructors and recursive references needs a separate guard-preservation judgment.

Also specify Data as the prescribed compositional certificate. Do not accidentally promise complete semantic reasoning about emptiness, intersections, or all extensionally data-only types.

### State the preservation obligations

Before freezing canonical bytes, make these specification obligations explicit:

1. **Substitution preservation:** substituting arguments satisfying declared interfaces into a checked body preserves its kind and well-formedness.
2. **Data preservation:** a Data proof for an open body remains valid under every such substitution.
3. **Specialization agreement:** for bounded closed examples, symbolic summary evaluation agrees with independently expanded descriptor classification.

These should guide implementation and qualification, rather than relying solely on successful example programs.

### Strong falsification suite

Author expected judgments independently of the production checker. Prioritize:

| Test family | Required evidence |
|---|---|
| **Precision and cache isolation** | All supplied Page/Ignore/Send/Wrapped/Compose cases; evaluate them in different orders to catch specialization facts contaminating generic summaries. |
| **Open validity and variance** | Both invalid open bodies above; `OnlyData` rejected as `Type -> Type`; higher-order polarity tests where a constructor demands stronger guarantees from its constructor argument than the expected interface permits. |
| **Unused arguments** | Reject `Ignore<DataPage<wire<Person>>>`, unbound references, and an unused constructor argument failing `Data -> Data`. |
| **Binding and identity** | Same-spelled parameters and constructors from distinct pinned modules remain distinct through nested applications; definition-scope imports remain unchanged. |
| **Opacity, cycles, recursion** | Literal lookalike tags remain byte-identical; ignored constructor-reference cycles reject; Data-preserving identity cannot manufacture guarded recursion. |
| **Explicit feature boundary** | Reject partial/over-application, constructor-valued bodies, bare constructor aliases, and unsupported kind syntax with stable diagnostics. |
| **Budgets and receipts** | Small shared DAGs with enormous expansions return incomplete/no output; changes to interfaces, constructor identities, pins, configuration, or intermediate ASTs invalidate stale evidence. |
| **Native qualification** | Pagination remains open and compiles in every pinned toolchain; unsupported higher kinds fail atomically; explicit residualization emits one Page; new C++ or Rust profiles receive separate positive and negative compiler fixtures. |

Build a small independent eager expander and primitive Data checker for bounded acyclic test programs. Compare it against the optimized summary implementation, while keeping hand-authored negative fixtures for invalid open definitions.

The central architectural change should therefore be a **constructor application calculus plus retained evidence**, owned by Horos. Bittype exposes that structure and consumes checked judgments; native mappings declare their representable subsets explicitly. This preserves the required precision without making general tools understand Function, Stream, Cell, or any other library’s identity or protocol.