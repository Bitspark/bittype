# Advice for consult
**Research doc:** doc_0eea29f2-5250-4c7d-89ce-608d96fd1b2f.md
**Document:** doc_0eea29f2-5250-4c7d-89ce-608d96fd1b2f
**Advisor:** ChatGPT, model not read from the answer page
**Session:** nsess_1fd56f60-dd60-4eaa-99f5-3679ba3a32b5
**Run:** run_f7cfb3b5-7cf2-40d8-8884-ce150a7adcf7
**Chat:** https://chatgpt.com/c/6af07002-0f54-8390-bc6c-189427382e33
**Verify:** nightfall consult verify run_f7cfb3b5-7cf2-40d8-8884-ce150a7adcf7 --output <this file>

# Expert response: Kernel-owned semantic traversal and narrow source interpretation

Research 0013 — Design assessment and recommendations

## Executive assessment

I agree with the central architectural diagnosis: the generation kernel should own traversal of the shared semantic graph, while source integrations should own the interpretation of language-specific constructs and the judgments that determine whether those constructs are admissible.

I would adopt the proposed direction, but make three important refinements:

1. The kernel should own structural traversal, not semantic validity. It should determine what entity, term, binder, or child is structurally connected to another. It should never infer that the resulting structure is valid in the source language.
2. Admission should be an explicit, evidence-backed source judgment over a complete use. It must cover every actual, including arguments omitted by a native representation, without making the kernel responsible for interpreting source constraints.
3. The fundamental object exposed to mapping libraries should be a kernel-issued, run-scoped occurrence handle, not a source-created navigation record or a raw graph ID. This is important because a graph term and a use of that term are not necessarily the same thing.

My recommended architecture is therefore:

Shared semantic graph

Entities · Binders · Terms · Applications · Operands · Occurrence associations

Generation kernel — structural authority

Root resolution · Declaration opening · Reference traversal · Application binding · Child validation · Scope · Limits · Planning

Source integration

Selectors · Constructor interpretation · Payload interpretation · Applicability judgments · Evidence

Mapping libraries

Native representation · Declaration recipes · Operand selection · Captures · Native conventions

Closed native program

Existing target construction and validation

The key distinction is that the kernel controls how the graph is traversed, while the source controls what the traversed structures mean.

This is consistent with AC-02, AC-05 and AC-09, and it addresses the duplication identified in the current implementation. In particular, the existing `Open` and `Admit` methods combine responsibilities that should not belong to the same owner.

research(20261108-124839).md

## 1. Where is the principled line?

### 1.1 My proposed ownership boundary

I would draw the boundary as follows.

| Operation                                                     | Owner                                       | Reason                                       |
| ------------------------------------------------------------- | ------------------------------------------- | -------------------------------------------- |
| Resolve a syntax occurrence to its graph association          | Kernel                                      | Structural association already exists        |
| Identify the declaration entity at a root                     | Kernel                                      | Graph traversal                              |
| Decide whether the entity is generatable                      | Source                                      | Language-specific declaration classification |
| Retrieve declaration body and real binders                    | Kernel                                      | Direct graph structure                       |
| Follow a reference to its target                              | Kernel                                      | Structural relation                          |
| Distinguish reference, application and provider expression    | Kernel                                      | Graph variant inspection                     |
| Pair application arguments with real binders                  | Kernel                                      | Structural correspondence                    |
| Determine whether a term is a child of another                | Kernel                                      | Graph containment                            |
| Interpret provider-expression constructors                    | Source                                      | Language semantics                           |
| Decode provider-expression operands                           | Source                                      | Provider-defined encoding                    |
| Determine stable declaration identity                         | Source                                      | Language-specific identity                   |
| Determine whether a use satisfies its assumptions             | Source                                      | Semantic judgment                            |
| Determine whether a type is supported by a generation profile | Source, through its interpretation contract | Language-specific admissibility              |
| Select native representation                                  | Mapping library                             | Representation ownership                     |
| Validate native output                                        | Target                                      | Native-language semantics                    |

There is one qualification: structural traversal and structural admission are not quite the same thing.

The kernel should own the structural facts and verify that every handle belongs to the admitted graph and the current generation context. It should not decide that an entity is a valid horos type definition merely because its graph `Kind` happens to contain a particular key.

That classification belongs to horos.

### 1.2 The three term variants provide a clean boundary

The existing graph has exactly three term variants:

- `Reference`
- `Application`
- `ProviderExpression`

This is a particularly useful division.

Reference

The kernel follows `Target`.

It can determine whether the target is a parameter entity or a declaration entity, and whether that parameter belongs to the declaration whose body is currently being traversed.

It must not inspect `Entity.Kind` or `Entity.Name` to decide what that declaration means.

Application

The kernel follows `Callee` and `Arguments`.

When the callee is a direct reference to a declaration, it can use `DirectParameterBindings` to associate actuals with real binders.

It must not infer saturation or source-language applicability from that correspondence.

ProviderExpression

The kernel recognizes that this is a provider-owned expression and follows its graph-level operand relations.

It does not interpret `Operator`, `Encoding`, `Data` or `Locator`.

For example, the kernel may establish that an operand is structurally a child of a record expression. Only horos can explain that locator `[k+2, 3]` denotes the type of the kth record field.

That is the right division.

research(20261108-124839).md

### 1.3 Scope: direct binder ownership, not an inferred Owner-chain rule

I recommend that the kernel use the set of real binders belonging to the declaration currently being opened.

For the running example:

```
type Page<T> = { items: list<T>; next?: text; };
```

When traversing `Page`:

- `Page.T` is an admissible formal binder.
- A reference to `Page.T` resolves to its native binding.
- A reference to a parameter belonging to another declaration is a scope violation.

I would not make the kernel infer lexical visibility by walking an arbitrary `Owner` chain.

The graph's `Owner` relation is structural, but the semantics of lexical scope, enclosing declarations and inherited assumptions can differ between source languages.

For horos, the direct-parameter rule is sufficient. If a future language supports nested generic scopes or inherited parameters, its integration should explicitly establish the appropriate scope judgment.

In other words:

- The kernel owns the mechanics of checking binder membership.
- The source owns the language's rules for which binders are applicable.
- An arbitrary ownership chain must not silently become a language-independent scoping rule.

### 1.4 Does the earlier expert advice support this?

Yes.

The distinction between correspondence and applicability is exactly the right one.

`DirectParameterBindings` establishes:

> This argument structurally corresponds to this parameter.

It does not establish:

> This argument is a valid substitution under the source language's rules.

The kernel can safely use the former to plan native arguments. The source must establish the latter through previously supplied evidence.

The earlier advice should not be interpreted as prohibiting the kernel from following graph structure. It prohibits the kernel from becoming a semantic evaluator.

### 1.5 Comparable architectural patterns

Three useful precedents are:

- LLVM IR: A common intermediate representation owns its instruction and operand relationships, while language front ends own source-language interpretation and lowering decisions.
- MLIR dialects: Shared infrastructure owns operation traversal and generic IR mechanics, while dialect-specific interfaces supply operation-specific semantics, verification and transformations.
- Roslyn semantic models: Syntax and structural relationships are distinct from semantic judgments such as symbol binding and type information.

The relevant lesson is not that bittype should copy any of these systems wholesale. It is that a shared representation should own the invariant mechanics of navigating that representation.

A language hook should supply the meaning of a node, not reimplement the infrastructure's ability to find that node's children.

## 2. How should admission judgments be shaped?

This is the most important part of the redesign after traversal itself.

The current interface makes it too easy to confuse four fundamentally different conditions:

1. No structural view exists.
2. A semantic judgment has rejected the input.
3. A required fact is unknown or uncovered.
4. The request or graph is malformed.

These must remain distinct.

### 2.1 Use whole-use judgments

I recommend a source-owned operation conceptually equivalent to:

```
JudgeUse(
    context.Context,
    EvidenceContext,
    DeclarationContext,
    UseOccurrence,
) (Judgment, error)
```

The judgment applies to the complete source use, including all its actuals.

It must be evaluated under the assumptions of the declaration in which the use occurs.

For example:

```
type Page<T> = { items: list<T>; next?: text; };
type People = Page<Cell<Person>>;
```

The application `Page<Cell<Person>>` is not admitted merely because the outer application is structurally understandable.

The source must establish that the complete application, including `Cell<Person>`, satisfies its language-specific requirements.

This also prevents a mapping library from bypassing a source constraint by omitting a native argument.

### 2.2 Separate judgment from traversal

I would use two independent operations.

Kernel operation: structural inspection

Determines:

- The term variant.
- The referenced entity, if any.
- The application callee and arguments.
- The binder-to-actual associations.
- The structural children.
- The current declaration and environment.

Source operation: semantic judgment

Determines:

- Whether the use is admissible.
- Whether its assumptions hold.
- Whether its language-specific constraints hold.
- Whether the available evidence establishes the required conclusion.

Neither operation should attempt to do the other's job.

### 2.3 Recommended outcome model

I would use a typed outcome rather than making every non-success condition an error.

| Outcome                               | Owner                       | Recommended representation                 |
| ------------------------------------- | --------------------------- | ------------------------------------------ |
| Structural view available             | Kernel                      | Successful structural view                 |
| No structural view                    | Kernel                      | Explicit value outcome                     |
| Judgment holds                        | Source                      | Judgment value                             |
| Judgment rejected                     | Source                      | Judgment value with diagnostic information |
| Judgment undetermined                 | Source                      | Judgment value                             |
| Required evidence unavailable         | Evidence/admission boundary | Typed error                                |
| Malformed graph or mismatched request | Kernel/contract boundary    | Typed error                                |
| Foreign handle or wrong scope         | Kernel                      | Typed error                                |
| Resource limit exhausted              | Kernel                      | Typed error                                |

An illustrative API:

```
type JudgmentStatus uint8

const (
    JudgmentHolds JudgmentStatus = iota
    JudgmentRejected
    JudgmentUndetermined
)

type Judgment struct {
    Status JudgmentStatus
    Report ReportReference
}

type StructuralStatus uint8

const (
    StructuralAvailable StructuralStatus = iota
    StructuralUnavailable
)

type StructuralResult struct {
    Status StructuralStatus
    View   StructuralView
    Reason StructuralReason
}
```

The exact Go representation can differ. The important point is that a structural limitation must not be mistaken for a semantic rejection.

For example:

- A parameter-headed application may be structurally unsupported by the generation kernel.
- An arrow-sorted type may be semantically valid horos but unsupported by a generation profile.
- An absent sort fact means that the required judgment has not been established.
- A graph address outside the admitted graph is malformed input or a scope violation, depending on the contract boundary.

These are not interchangeable.

research(20261108-124839).md

### 2.4 What should happen with omitted actuals?

The kernel must admit every actual before it permits a mapping to represent the parent application.

This is a mandatory invariant.

Consider a native convention that omits a phantom source parameter. The mapping library must not be able to make an otherwise invalid source use acceptable simply by never capturing that argument.

My recommendation:

1. The kernel obtains all structural actuals.
2. It asks the source to judge the complete parent use.
3. It independently ensures that each actual occurrence is admitted, even if the native convention will omit it.
4. If an actual is subsequently represented, its children are traversed through the normal kernel mechanism.
5. No mapping callback can bypass these steps.

I would distinguish admission of an actual occurrence from recursively generating its representation. An omitted actual need not produce native output, but it cannot escape source admission.

The source's judgment may validate the complete semantic application using its existing analysis evidence. The kernel must not rerun that analysis.

### 2.5 What about ignored provider-expression operands?

I would apply a similar principle, with one important distinction.

An operand ignored by a mapping is not necessarily equivalent to an omitted generic actual.

A generic actual participates in the application's parameter substitution. A provider-expression operand is defined by the source's constructor semantics.

Therefore:

- The source judgment must establish the admissibility of the complete provider expression.
- The kernel must validate any operand that a mapping captures as a structural child.
- An ignored operand must not be silently assumed to be semantically irrelevant.
- Whether a particular operand must independently satisfy a generation-specific restriction should be explicit in the source's generation contract.

I would not automatically recursively generate every operand of every provider expression. That would make generation eager and could impose unintended requirements on source constructs.

However, a source integration must not be able to use an ignored operand to evade a required applicability judgment.

### 2.6 A crucial distinction: evidence coverage versus evidence truth

The existing evidence model has a useful property:

`Complete/Undetermined` is not the same as `Complete/Rejected`, and neither is the same as an absent report.

Preserve that distinction.

I would make the source judgment consume only admitted reports and evidence references. The source must not fill an evidence gap by invoking an analyzer.

For horos:

- Missing sort fact: uncovered evidence, not necessarily malformed graph.
- Arrow sort: a known result that the generation profile does not support.
- Proper sort: a successful generation-specific restriction.
- Failed source membership: a rejected semantic judgment.

This is especially important because horos's analysis can accept constructor-valued declarations that generation intentionally does not support.

The analysis and generation layers must remain independent.

## 3. What is a subject, and how should it be identified?

I think the current design exposes an important conceptual problem: it treats a graph term, an occurrence, and an effective subject as if they were variations of the same object.

They are not.

### 3.1 Three distinct concepts

I would explicitly distinguish:

| Concept           | Meaning                                                             | Owner              |
| ----------------- | ------------------------------------------------------------------- | ------------------ |
| Graph term        | A node in the immutable semantic graph                              | Semantic graph     |
| Occurrence        | A use of a term in a particular declaration and environment         | Generation kernel  |
| Effective subject | The source-certified semantic subject associated with an occurrence | Source integration |

This distinction matters because the graph permits shared terms.

Imagine the same term being reached from two different declaration bodies. The graph term ID is identical, but the two occurrences can have different environments and different applicability judgments.

Consequently:

\\[ \text{GraphTermID} \neq \text{OccurrenceIdentity} \\]

Nor should two occurrences be considered equivalent merely because the source reports the same effective subject.

### 3.2 Kernel-issued occurrence handles

I recommend replacing source-issued `Use` addresses with opaque handles minted by the kernel.

Conceptually:

```
type Occurrence[F any] struct {
    // private:
    run       *runIdentity
    graph     *graphIdentity
    term      semantics.TermID
    declaration semantics.EntityID
    environment environmentIdentity
    path      occurrencePath
}
```

This is illustrative, not a proposed public struct. Its fields should remain inaccessible to mapping libraries.

The handle should identify:

- The generation run.
- The exact admitted semantic graph.
- The underlying term.
- The declaration under which the term is reached.
- The environment in which it is interpreted.
- Its structural occurrence.

This solves several problems simultaneously.

Foreign-reference rejection

A handle from another run, graph, or source session is rejected. Merely possessing a valid graph-local ID is insufficient.

Scope preservation

A term reached under another declaration's assumptions cannot be passed off as the current declaration's occurrence.

Shared graph terms

Two paths to the same graph term can produce distinct occurrences.

Capture identity

A mapping callback can capture a precise occurrence, not just an address that happens to resolve to a term.

### 3.3 What about effective subjects?

I would retain the concept only if there is an actual semantic consumer for it.

It should not be used as the kernel's traversal identity.

In particular:

- The kernel must not deduplicate occurrences using effective-subject equality.
- The kernel must not infer that two uses have the same assumptions because their effective subjects match.
- Memoization must include the relevant environment and generation context.
- Recursive-type handling must not depend on an equivalence relation that has not been defined.

If future features require semantic equivalence, let the source provide an explicit equivalence operation with documented guarantees.

Do not introduce equivalence into the structural handle model prematurely.

### 3.4 Callback lifetimes

The existing bounded-lifetime guarantee is correct and should be preserved.

I would implement a callback-scoped capability model:

```
type Expansion[F, T any] interface {
    Current() (OccurrenceView[F], error)
    Mapped(Occurrence[F]) (Expression[T], error)
}
```

The actual implementation would keep the capability associated with an active callback invocation.

When the callback returns:

- All accessors become invalid.
- All capture operations become invalid.
- Retaining a handle for later use fails.
- A callback cannot use a previously obtained handle in another expansion.

Go cannot fully enforce borrowing or linear ownership statically. Therefore, runtime capability validation is appropriate.

One additional refinement: callbacks should not receive the underlying graph object or a graph-wide reader. They should receive only the capabilities required for their operation.

This substantially reduces the possibility of bypassing the kernel.

### 3.5 Should the kernel reject shared terms?

Not categorically.

The graph explicitly permits a containment DAG, so a shared term is not inherently malformed.

I would distinguish:

- A term reached in multiple environments: distinct occurrences.
- A term shared within the same environment and structural context: potentially reusable.
- A term whose sharing violates a source-language capture or ownership rule: source rejection.
- A foreign or forged occurrence handle: kernel rejection.

The current horos integration rejects a term shared by two declaration bodies during its indexing pass. That may be a legitimate horos-specific constraint, but it should not automatically become a universal graph restriction.

The source should decide whether such sharing is semantically admissible. The kernel should preserve occurrence-specific scope regardless.

research(20261108-124839).md

## 4. What should mapping libraries see?

This is where the redesign must be especially careful.

The mapping libraries are independent extension points. They must not be forced to reconstruct structural relationships, but they must retain enough source-specific information to interpret their constructs.

I recommend a two-part interface:

1. A kernel-owned structural view.
2. A narrow, source-owned interpretation reader.

### 4.1 The structural view

A mapping callback should receive something conceptually like:

```
type UseView[F any] struct {
    Occurrence Occurrence[F]

    Kind UseKind

    // For a reference or direct application:
    Declaration DeclarationView[F]

    // For an application:
    Actuals []ActualView[F]

    // For a parameter reference:
    Bound Binder[F]
}

type ActualView[F any] struct {
    Binder Binder[F]
    Use    Occurrence[F]
}

type DeclarationView[F any] struct {
    Declaration Declaration[F]
    Formals     []FormalView[F]
    Body        Occurrence[F]
}
```

These records are produced by the kernel, not the source.

The mapping library should never have to ask:

- Which declaration does this reference target?
- Which binder corresponds to this argument?
- Is this occurrence a child of the current expression?
- Is this binder in the declaration currently being expanded?

Those are kernel questions.

### 4.2 Keep a narrow source reader

The horos mapping libraries still need to understand things that the kernel cannot.

For example:

- Record field names.
- Optional versus required fields.
- Descriptor constructors.
- Descriptor discriminators.
- Source syntax locations.
- Declaration names and module pins.
- Source-specific unsupported diagnostics.

I would retain a source-specific reader for those purposes, but redesign it so it cannot supply an alternative traversal.

For example:

```
type Reader interface {
    DeclarationInfo(Declaration[h.Language]) (
        DeclarationInfo, error)

    ParameterInfo(Binder[h.Language]) (
        ParameterInfo, error)

    ReadDescriptor(Occurrence[h.Language]) (
        DescriptorView, error)

    ReadRecord(Occurrence[h.Language]) (
        RecordView, error)
}
```

The important difference is that the reader consumes kernel-issued handles.

It does not return independently constructed `UseData` records containing callees, environments and actuals.

A `RecordView` may describe fields and their source locations, but each field's type should be identified by an occurrence handle that the kernel can validate.

### 4.3 Remove `CheckChild`

I would remove `Source.CheckChild` entirely.

The kernel already knows the structural child relations.

For a captured child, the kernel can verify:

1. The parent handle belongs to the current callback.
2. The child handle belongs to the same run and graph.
3. The child is an actual structural child of the parent.
4. The child is being used in the permitted declaration and environment.
5. The capture is within the relevant work limits.

The source reader should not be responsible for this.

This is one of the clearest opportunities to eliminate duplicated logic.

### 4.4 Actuals should be identified by real binder, not position

My recommendation is strongly in favor of binder-keyed actuals.

Consider:

```
type Cell<V> = wire<V>;
```

The semantic relationship is:

\\[ \text{Cell.V} \longleftrightarrow \text{Person} \\]

It is not fundamentally:

\\[ \text{Argument}[0] \longleftrightarrow \text{Person} \\]

Position is only a convenient encoding of the relationship in the current graph.

The mapping library should receive:

```
ActualView{
    Binder: cellV,
    Use:    personOccurrence,
}
```

This has several benefits:

- It preserves the identity of the real source binder.
- It avoids confusing source parameter order with native parameter order.
- It makes reordering explicit.
- It supports conventions that omit or add native parameters.
- It avoids accidental positional assumptions in independently authored mapping libraries.

A mapping library can then explicitly select the actual associated with the binder it wants.

For a single-parameter declaration, a helper can make this convenient without changing the underlying semantics.

Position may still be exposed as an ordered iteration convenience, because the source declaration has an established parameter order. But it should not be the identity of the binding.

### 4.5 The native generic convention should remain separate

I would preserve the conceptual distinction already present in `GenericConvention`.

There are three separate mappings:

1. Source binder to native declaration parameter.
2. Source actual to source binder.
3. Native declaration parameter to native application argument.

For example:

```
Source:
    Page<T>

Native declaration:
    Page<U>

Binding:
    T -> U

Source application:
    Page<Cell<Person>>

Structural binding:
    T -> Cell<Person>

Native application:
    Page<Cell<Person>>
```

The kernel composes these relationships but does not assume they have identical positions.

This is particularly important when a native mapping introduces helper parameters or eliminates phantom parameters.

### 4.6 What happens to the Cell mapping?

The Cell mapping becomes simpler conceptually.

Instead of reading `use.Actuals[0]`, it selects the actual associated with Cell's real binder `V`.

It then captures that occurrence using the kernel's capture mechanism.

The Cell library still owns the native representation:

```
model.External(
    "@bitspark/cell",
    "Cell",
    g.Mapped(slot),
)
```

The Cell body remains unopened.

That last property is important: the kernel should follow the application structure to its actual, but selecting an external representation must not force expansion of the declaration's source body.

research(20261108-124839).md

## 5. What else should be watched?

I see several additional risks beyond the five questions.

### 5.1 Evidence must be inseparable from the graph

This is a critical issue.

The document correctly notes that the graph itself does not have a pointer back to its snapshot. Its relationship to the snapshot exists through the result that published it.

I would make the evidence context an immutable kernel-owned object created during admission.

It should contain:

- The exact snapshot identity.
- The exact admitted result identities.
- The typed semantic graph.
- The evidence requirements.
- The established reports.
- The provider and language profile.

The source must not be allowed to substitute a graph obtained from a different analysis result.

Ideally, the kernel should unpack the graph itself from the admitted evidence, rather than accepting an arbitrary graph supplied by the source integration.

This is directly aligned with the candidate split's first point.

### 5.2 The graph needs more structural query support

The current graph offers point queries and `DirectParameterBindings`, but lacks some useful structural operations.

I would add language-neutral queries where they represent universal graph mechanics:

| Query                              | Purpose                        |
| ---------------------------------- | ------------------------------ |
| Entity declarations at syntax site | Root resolution                |
| Direct entity parameters           | Declaration opening            |
| Entity body                        | Declaration traversal          |
| Term variant and structural edges  | Term traversal                 |
| Direct structural children         | Capture validation             |
| Direct parameter bindings          | Application substitution       |
| Parent/child membership            | Scope and capture verification |

These should be graph operations, not horos helpers.

However, I would resist adding a large general-purpose traversal framework to `semantics` merely to accommodate generation.

The semantic graph should expose primitive structural capabilities. The generation kernel should own generation-specific traversal policy.

That keeps the graph reusable by other consumers.

### 5.3 The four existing error classes need disciplined handling

I would preserve the distinction between:

- `ErrScope`
- `ErrUnsupported`
- `ErrInvalid`
- `ErrLimit`

But I would define their ownership more precisely.

| Error            | Recommended meaning                                                                |
| ---------------- | ---------------------------------------------------------------------------------- |
| `ErrScope`       | The kernel has detected an invalid structural relationship, environment or handle  |
| `ErrUnsupported` | Valid source structure cannot be represented under the selected generation profile |
| `ErrInvalid`     | A malformed contract, request or graph invariant                                   |
| `ErrLimit`       | A bounded generation resource has been exhausted                                   |

A rejected semantic judgment should normally carry source diagnostic information rather than being converted indiscriminately into `ErrInvalid`.

Likewise, an uncovered fact must not be translated into a semantic rejection.

### 5.4 Deterministic planning and memoization

The redesign should specify which objects may be memoized.

I recommend:

- Declarations: memoized by admitted declaration identity within a run.
- Occurrences: memoized only when occurrence identity and environment are equivalent.
- Semantic subjects: never assumed equivalent without an explicit source guarantee.
- Native expressions: memoized only where the mapping contract establishes that reuse is safe.

This becomes particularly important when recursive types are eventually introduced.

A declaration-level recursion detector is useful, but it should not confuse two different occurrences of the same graph term.

### 5.5 Recursive types need an explicit policy

Although recursive types are out of scope, the new architecture should not accidentally prevent them.

The kernel should distinguish:

- A declaration currently being expanded.
- A reference to a previously reserved declaration.
- A structural cycle in the term containment graph.
- A semantic recursive reference through declaration identity.

These are different situations.

The current graph rejects cyclic term containment but permits cyclic entity references. That is a useful starting point.

Future recursive generation should probably rely on declaration reservation and reference resolution, rather than trying to permit cyclic term containment.

### 5.6 Native enforcement of source bounds remains a separate problem

GEN-02 is correctly out of scope.

Nevertheless, the new interface should preserve enough information to address it later.

In particular, the source-owned contract should be able to expose the relevant source constraint in a structured, stable form, without forcing the kernel to understand horos kind semantics.

A future target-specific enforcement mechanism could then translate source constraints into native constraints where possible.

The kernel should not silently convert a proven source membership into an unconditional native type fact.

### 5.7 Mapping libraries need stronger capability boundaries

One risk I would address immediately is that the current mapping library is parameterized by the entire `SourceSession`.

That gives a mapping callback more authority than it needs.

I would split the capabilities:

- Declaration recipe: declaration metadata, formal binders, body occurrence.
- Descriptor mapping: descriptor interpretation and validated operand capture.
- Application mapping: declaration selector and binder-keyed actuals.
- Native construction: no source access.

A library should not receive the complete source session merely because it needs to read one record descriptor.

This reduces accidental coupling and makes the contract easier to test.

### 5.8 Test the invariants directly

The existing qualification gaps are significant enough that I would make negative tests part of the redesign acceptance criteria.

| Test                                                    | Expected result                    |
| ------------------------------------------------------- | ---------------------------------- |
| Graph from another snapshot                             | Rejected                           |
| Evidence from another snapshot                          | Rejected                           |
| Foreign declaration handle                              | Rejected                           |
| Foreign occurrence handle                               | Rejected                           |
| Wrong environment                                       | `ErrScope`                         |
| Parameter owned by another declaration                  | `ErrScope`                         |
| Capture of a non-child                                  | `ErrScope`                         |
| Capture after callback return                           | Rejected                           |
| Omitted actual with failed admission                    | Rejected                           |
| Omitted actual with missing evidence                    | Evidence failure                   |
| Missing sort fact                                       | Uncovered evidence                 |
| Arrow sort                                              | Unsupported                        |
| Parameter-headed application                            | No structural view / Unsupported   |
| Curried application                                     | No structural view / Unsupported   |
| Zero-parameter constructor alias applied to an argument | Unsupported, not malformed-request |
| Multiple mappings for one selector                      | `ErrAmbiguous`                     |
| Mapping failure after declaration reservation           | No usable output                   |
| Exhausted reference or depth limit                      | `ErrLimit`, no partial output      |
| Shared graph term in different environments             | Distinct occurrences               |

I would also test that the independent source publishes a semantic graph and exercises exactly the same kernel traversal as horos.

Otherwise the architectural improvement has not actually been demonstrated.

research(20261108-124839).md

## 6. Recommended implementation sequence

I would implement the redesign in six stages, with each stage establishing a testable architectural invariant.

| Stage | Work                         | Acceptance condition                                                |
| ----- | ---------------------------- | ------------------------------------------------------------------- |
| 1     | Evidence and graph admission | Kernel owns the exact admitted graph                                |
| 2     | Structural traversal         | Kernel opens declarations and follows references/applications       |
| 3     | Occurrence handles           | Scope, graph identity and callback lifetime enforced                |
| 4     | Source judgment contract     | Whole-use applicability and unknown/rejected distinctions preserved |
| 5     | Mapping API redesign         | Libraries consume kernel views and narrow source readers            |
| 6     | Migration and qualification  | Horos and independent source use identical traversal                |

I would migrate the horos mapping libraries only after the structural and judgment contracts are established.

Otherwise, the implementation risks reproducing the current interface under different names.

The final migration should remove:

- `Source.Open`
- `Source.Admit`
- `Source.CheckChild`
- Source-owned term-to-declaration traversal indexes
- Source-owned structural application pairing
- Mapping-library access to parallel navigation records

No compatibility aliases should remain, consistent with the stated completion criteria.

## 7. Final decisions on the five questions

For clarity, here are my direct answers.

| Question                | Recommendation                                                                                                                                                      |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1. Principled boundary  | Kernel owns all graph navigation and structural correspondence. Source owns interpretation, stable identity and semantic judgments.                                 |
| 2. Admission            | Source judges complete uses under their original assumptions. Kernel ensures all actuals are admitted, including omitted ones.                                      |
| 3. Subject identity     | Distinguish graph terms, occurrences and effective subjects. Use opaque, kernel-issued, run-scoped occurrence handles.                                              |
| 4. Mapping-library view | Provide kernel-owned declaration/use views and binder-keyed actuals, alongside a narrow source-specific reader.                                                     |
| 5. Additional risks     | Evidence provenance, shared-term environments, callback capabilities, memoization, deterministic failure, recursive types and future native constraint enforcement. |

## 8. Overall verdict

Recommend adopting the kernel-owned traversal direction

I would accept the candidate split, subject to the refinements above.

There is one architectural principle I would make explicit in the implementation specification:

The kernel must be the sole authority on structural navigation, but it must never become the authority on source-language meaning.

This gives each participant a clear responsibility:

- The semantic graph records structure and facts.
- The kernel follows structure, maintains scope, enforces resource limits and plans output.
- The source integration interprets language-specific constructs and establishes semantic admissibility from supplied evidence.
- Mapping libraries choose native representations without reconstructing source navigation.
- Targets validate and construct the closed native program.

The most important design choice is the occurrence handle. Without it, the redesign could move traversal into the kernel while still leaving scope, foreign-reference rejection and callback safety dependent on source-owned records. That would be an incomplete architectural correction.

I would also make binder-keyed actuals non-negotiable. Real binder identity is the semantic relationship; position is merely one possible representation of that relationship.

The decisive acceptance test should be `People = Page<Cell<Person>>`. The kernel must traverse the entire declaration/reference/application structure identically for horos and the independent test source, while the two integrations supply different interpretations and judgments. The Cell mapping must receive the actual associated with its real binder, capture Person, and avoid opening Cell's body.

If that works, alongside the negative scope, evidence, foreign-handle and omitted-actual tests, the project will have demonstrated the intended separation rather than simply refactored its existing interfaces.

My final recommendation: redesign the contract from the graph and its ownership guarantees, not from `Open` and `Admit`. Make structural navigation a kernel capability, semantic admission a source capability, and occurrence identity a kernel-enforced invariant.