# Advice for consult
**Research doc:** doc_e667d0e2-7926-4302-8252-ac6670b1af61.md
**Document:** doc_e667d0e2-7926-4302-8252-ac6670b1af61
**Advisor:** ChatGPT, model "6 Pro" as labelled on the answer page
**Session:** nsess_f4bf5e68-6d60-4145-9e35-a45c7fbaf842
**Run:** run_a409da7b-a708-41bc-bba3-bd0b011d33ff
**Chat:** https://chatgpt.com/c/6ae89698-5720-83ed-ab63-cea18b8e6996
**Verify:** nightfall consult verify run_a409da7b-a708-41bc-bba3-bd0b011d33ff --output <this file>

## Recommendation

I would proceed with **option 3: a shared graph of established semantic facts, provider-owned expression data, and explicit coverage**. The proposed architectural boundary is sound, but the first implementation should be narrower and more precise than the illustrative model.

The kernel should establish structural integrity, ownership, and association with a particular captured request. Horos should establish the correctness of resolution, parameter binding, descriptor interpretation, and semantic verdicts. Those responsibilities need separate contracts.

Before freezing the API, I would make six changes:

1. Distinguish **references, applications, and provider expressions** explicitly.
2. Preserve descriptor data through **bounded, immutable, provider-encoded payloads**, alongside exposed semantic operands.
3. Add explicit **occurrence-to-term associations and declaration-dependency facts**.
4. Replace the caller-supplied snapshot string with **immutable request handles and request-bound results**.
5. Separate **analysis completion, semantic verdicts, and graph availability**.
6. Complete the Horos admission migration, including any **Horos-owned unresolved syntax traversal** needed to admit inputs before resolution.

This review treats the revisions and code excerpts in your attachment as the implementation baseline. I have not independently audited those repository commits. The external references below support the language counterexamples and architectural comparisons.

## 1. The smallest useful shared semantic graph

### Start from the four consumer queries

The four stated queries justify a relatively small core:

| Component | Minimum responsibility |
|---|---|
| Entities | Distinguish declarations, their parameters, and any necessary module or builtin entities. |
| Terms | Represent references, applications, and provider-defined expressions. |
| Syntax associations | Connect an occurrence to an established entity, and an expression occurrence to its term. |
| Dependencies | Record directed, labeled declaration dependencies. |
| Analysis reports | State which obligations were completed, over which domain, with which verdicts. |
| Diagnostics | Explain rejection, missing context, unsupported work, and resource exhaustion. |

**Scopes, scope links, and `Members` are not necessary to answer those four queries.** I would defer them unless an immediate consumer needs scope navigation or resolution explanations. Established resolution can be inspected without exporting a generic lookup algorithm.

This is consistent with the useful separation in existing compiler APIs: Go’s `types.Info` distinguishes definitions, uses, expression types, scopes, and instantiations; TypeScript’s documented checker API distinguishes syntax nodes, declared entities represented by symbols, and types. Those distinctions are useful precedents without requiring their entire APIs to become the shared contract. [Go Packages](https://pkg.go.dev/go/types%40go1.27.3?utm_source=chatgpt.com)

### Keep entities distinct from their bodies

An entity should retain:

- A snapshot-local identity.
- A provider-owned kind.
- A display name.
- Its defining owner, where applicable.
- Current declaration sites.
- Its ordered parameters.
- An optional body term.

I would rename `Entity.Meaning` to **`Body`**, or define it explicitly as the checked declaration expression. “Meaning” invites consumers to assume expansion, normalization, or complete denotation—none of which this first model establishes.

Parameter ownership must be explicit or mechanically derivable from a uniquely owning declaration. Display names must play no role in substitution.

Allow entities without declaration sites for builtins and explicitly synthesized entities. Do not invent an AST declaration for a builtin merely to satisfy a nonempty-site invariant.

Keep language-specific distinctions in provider vocabulary. For example, a Go alias and a Go defined type can have different provider-owned entity kinds. The kernel need not understand their consequences.

### Use three term variants

The essential distinction is:

```go
// Conceptual shapes, not a proposal for exported mutable storage.

Reference {
    Target EntityID
}

Application {
    Target    EntityID
    Arguments []TermID
}

ProviderExpression {
    Operator Key
    Encoding Key
    Data     ImmutableBytes
    Operands []Operand
}
```

A `Reference` identifies an entity. An `Application` applies a declaration to arguments. A `ProviderExpression` preserves a language operation or descriptor together with its data and semantic children.

This resolves several problems at once:

- `Reference(Person)` differs from `Application(Person, [])`.
- An open parameter is a reference to an established parameter entity.
- Operators such as Horos `list-of` or TypeScript indexed access need not become artificial declaration entities.
- A constant can be a provider expression with data and no semantic children.
- The model does not assume that every argument or operand is a type.

For the initial `Application` contract, I recommend **fully resolved positional instantiations with effective arguments**. Horos’s current explicit application model fits that contract.

Under that restriction, expose parameter bindings by zipping the target declaration’s ordered parameter identities with its effective arguments. There is no need to store an independently mutable `Substitution` map.

If both representations are stored, agreement must be an invariant. Otherwise an application can simultaneously claim:

- Its first argument is `Person`.
- Its first parameter is bound to something else.

A future provider may need to distinguish written arguments, inferred arguments, defaults, and effective parameter bindings. That is a concrete extension point. It should not be hidden inside an initially ambiguous `Arguments` field.

### Preserve descriptor data without inventing a common type algebra

For Horos, my preferred first implementation is:

1. Preserve the canonical encoding of `Expression.Value()`.
2. Identify its format with a versioned provider-owned encoding key.
3. Expose semantic child terms at exactly the paths supplied by `Expression.Operands()`.
4. Use `Expression.Constructor()` to identify the provider operation.

This preserves:

- Scalar distinctions.
- Exact literal values.
- Record labels.
- Required versus optional fields.
- Open versus closed record mode.
- Refinement operands.
- Other literal metadata not represented by type-argument lists.

The shared kernel copies and bounds those bytes. It does not decode Ontos or interpret Horos constructors.

That choice still provides useful agnostic inspection: declaration identities, occurrence bindings, applications, parameter bindings, and dependencies remain explicit. Descriptor interpretation remains provider-owned, as required.

The guarantee should be written narrowly:

> The graph preserves provider expression data and exposes the semantic operands identified by the provider.

It does **not** establish a language-neutral interpretation of every field, label, constant, or refinement.

If a later consumer needs “enumerate fields and optionality across languages,” add a specific projection capability with its own coverage contract. Do not silently derive that capability from preserved bytes.

For the adverse literal example, an `exact` value containing a tuple tagged like `horos/decl/param` remains literal data. The adapter follows Horos’s operand paths; it never discovers semantic occurrences by scanning arbitrary tuple tags.

Also distinguish **representation equality from semantic equivalence**. Equal payload bytes establish equality of that encoded representation. They do not establish type equivalence, alias identity, or equivalence under different declaration contexts.

I prefer this payload mechanism to source-fragment references as the sole fidelity mechanism. Source references can be useful provenance, but they couple interpretation to a particular syntax snapshot and do not naturally cover provider-produced constants or elaborated expressions.

### Make the pagination mapping precise

The mapping should retain at least these separate entities:

| Entity | Required representation |
|---|---|
| `pagination.Page` | Declaration with a body and one ordered parameter. |
| `Page.T` | Parameter entity owned by `Page`. |
| `people.Person` | Separate declaration entity. |
| `people.PeoplePage` | Separate declaration whose body applies `Page`. |

There is a subtlety in the supplied Horos model: the shorthand `Person` may itself be encoded as an application with zero arguments. Consequently, the body may be:

```text
PeoplePage.Body =
    Application(Page, [
        Application(Person, [])
    ])
```

Its binding view is:

```text
Page.T → Application(Person, [])
```

Do not normalize the inner zero-argument application into a reference merely because the explanatory notation omitted angle brackets.

For an open caller such as `Outer<T> = Page<T>`, the argument term refers to **`Outer.T`**, while the application binding’s parameter identity is **`Page.T`**. Traversal of application arguments must remain in the caller’s binding environment.

A practical mapping can allocate entities in a first pass using Horos’s canonical declaration and parameter identities, then translate expressions in a second pass. This avoids name-based identity and handles forward references cleanly.

No expansion follows from the presence of a binding. `PeoplePage` remains inspectable as a declaration, and specialization remains an explicit Horos operation.

### Add syntax associations and dependencies explicitly

The proposal has bindings but lacks a direct expression-occurrence-to-term association. Add one so a consumer can locate an application in syntax and inspect its term without searching every declaration body.

Use role-qualified occurrences when one AST node can contain several semantic positions. A suitable conceptual occurrence address is:

```text
(current NodeRef, provider-owned role, optional ordinal)
```

A binding at that occurrence records an established target. Multiple records must mean multiple role-distinguished facts, never an unlabeled set of possible candidates.

Add explicit dependency records:

```text
(from entity, to entity, relation key, optional occurrence)
```

For Horos, define a precise relation such as “declaration application dependency.” Include applications nested inside arguments and descriptor operands. Under that definition, `PeoplePage` directly depends on both `Page` and `Person`.

Do not call this “all semantic dependencies.” Other providers may need distinct relations for imports, constraints, defaults, or elaboration.

An absent edge only supports a negative inference when coverage promises complete enumeration of that exact relation over the relevant domain.

### Native counterexamples the representation must survive

These examples constrain the representation without requiring native semantic analyzers now:

| Example | Information that must survive | Work that can remain unsupported |
|---|---|---|
| Go `type ID int` versus `type ID = int` | Separate declaration identity and declaration form. Go distinguishes a new defined type from an alias. | Identity, assignability, and method-set reasoning. |
| Go `[3]T` versus `[4]T` | Element term and length data; array length participates in the type. | General constant evaluation. |
| Rust `Container<'a, T, 3>` | Distinct lifetime, type, and constant parameter sorts and argument data. | Lifetime inference, borrow checking, and trait solving. |
| TypeScript `T[K]` and `keyof T` | Provider operators and their operands. | Evaluation of those operators. |
| C++ `T*`, `T&`, and `T&&` | Distinct provider operations. | Conversion, reference binding, and overload rules. |

The Go specification supplies the alias/definition and array distinctions. Rust explicitly has lifetime, type, and constant generic parameters. TypeScript documents indexed access as a type operation, including combinations with `keyof`. The C++ working draft distinguishes pointers, lvalue references, and rvalue references. [The Go Programming Language](https://go.dev/ref/spec?utm_source=chatgpt.com)

Accordingly, retain provider-owned parameter sorts, and avoid interpreting every term as a type. Rust associated-type projections can likewise remain provider expressions until an actual provider supports resolving them. [The Rust Reference](https://doc.rust-lang.org/reference/items/associated-items.html?utm_source=chatgpt.com)

TypeScript mapped types remain outside the admitted profile described in the note. They are a useful future test for local binders, rather than a reason to introduce a full binder calculus now.

## 2. Request context, registration, ownership, and provenance

### Use three separate immutable boundaries

Keep these distinct:

1. The admitted syntax envelope.
2. The captured analysis request.
3. The published analysis result.

The typed adapter should create the latter two boundaries. A schematic interface is:

```go
type Provider[L any] interface {
    language.LanguageProvider
    Analyze(context.Context, Request[L]) (Result, error)
}

type Request[L any] struct {
    state *requestState[L] // private, immutable captured state
}

type Result struct {
    state *resultState // private, immutable published state
}

func NewAnalyzer[L any](
    syntax representation.Language,
    contract AnalysisContract,
    provider Provider[L],
) (Analyzer, error)
```

Prefer an opaque heterogeneous `Analyzer` handle constructed through `NewAnalyzer[L]`. That makes the typed adapter an enforced boundary rather than an optional convention.

Separate **syntax profile** from **analysis profile**:

- The syntax profile identifies the admitted input contract.
- The analysis profile identifies supported facts, obligations, options, and vocabulary.

Neither should be inferred from compiler versions or structural hashes.

### Capture an entry and an input catalogue

A request should contain an explicit entry `NodeRef` and an owned catalogue of input trees.

Do not impose one module per tree. The existing Horos projection can contain a module closure in one tree; the request design should accommodate that.

For version one:

- Require all input envelopes to match the registered ID, syntax profile, and actual Go marker.
- Reject duplicate `TreeID`s within the request.
- Capture the entry, inputs, dependency context, requested analyses, effective configuration, and limits before invoking the provider.
- Build the node catalogue from those captured inputs.
- Validate result sites against that catalogue.
- Keep historical `Origin.Nodes` outside the current semantic input closure.

The typed adapter can build the catalogue after `Unpack[L]`; a new public marker-agnostic AST catalogue is not necessary initially.

Two details deserve explicit tests:

**First, AST node zero may be valid.** Do not generalize the semantic graph’s zero-as-absent convention to `ast.NodeID`. Optional locations need an explicit representation of absence.

**Second, `TreeID` is an address within this captured request.** Its reuse in another request must confer no compatibility or freshness.

The initial provider should not acquire hidden dependencies from the network or filesystem. If dependency acquisition is added later, materialize that context before analysis or explicitly extend the provenance contract.

### Bind results when they are constructed

Allocate a fresh internal identity for each invocation and associate it with the captured request state. Request-created builders and their finished results carry that identity.

On return, compare the carried identity with the invocation’s identity.

**Do not attach the current identity to an arbitrary graph after the provider returns.** That would accept an accidentally retained result from an earlier invocation whenever its references happened to pass structural validation.

A request-created builder is a practical mechanism:

```text
request.NewBuilder()
    → add facts, reports, and diagnostics
    → finish immutable result
    → adapter verifies invocation association and integrity
```

The limitation must remain explicit: a provider can still copy stale facts into a fresh builder. A token establishes which request a publication is associated with; it cannot prove that the facts were computed correctly.

Hashes do not solve that stronger problem either. Correctness remains a trusted-provider obligation.

### Make ownership safe before publication

Deep copying only after a callback returns is insufficient if the provider retains and concurrently mutates the returned slices: the copy itself can race with those writes. Go requires synchronization for concurrent access involving mutation. [The Go Programming Language](https://go.dev/ref/mem?utm_source=chatgpt.com)

Use:

- Private request storage, with AST accessors returning independent copies.
- Builders that copy mutable inputs when accepting them.
- A defined builder concurrency rule.
- A finalization operation that publishes private storage.
- Result accessors that return immutable handles or recursively owned copies.

Audit nested storage: operand paths, literal bytes, diagnostics, parameters, coverage subjects, and source locations.

Providers must not mutate arguments concurrently while a builder is copying them. Likewise, callers must synchronize mutation of raw ASTs during admission. An ownership boundary cannot retroactively make an ingress race safe.

### State what the kernel can actually verify

| Kernel or typed-adapter guarantee | Provider promise |
|---|---|
| Marker, language, and profile compatibility | The analysis implements the declared language contract. |
| Result belongs to this invocation | Facts were computed from the captured request. |
| IDs and current syntax sites exist | The stated declaration is the correct resolved target. |
| Term variants and parameter ownership are structurally consistent | Argument sorts and binding rules are correct. |
| Payloads and graph size respect limits | Payloads and operand locators faithfully preserve meaning. |
| Reports are structurally consistent with requested obligations | Claimed coverage is complete and verdicts are justified. |
| Published storage has independent ownership | Supplied annotations were checked against authoritative results. |

The kernel cannot prove complete resolution or dependency enumeration because it does not know which language-specific positions require those facts.

For a bounded first graph, I would require acyclic **term containment**, while allowing entity references and dependency relations to form cycles. Recursive semantics can be represented through entity references or provider-owned recursive constructs. A generic acyclicity rule over the whole graph would reject legitimate language representations.

### Attach analyzers explicitly and release locks before callbacks

`RegisterAnalyzer` should attach a capability to an existing language entry and reject:

- Missing language registration.
- Invalid or typed-nil providers.
- Marker, ID, or syntax-profile mismatch.
- Duplicate analyzer attachment.

Capture provider metadata during adapter construction. Registration should not invoke arbitrary provider callbacks while holding the registry lock.

Execution should snapshot the analyzer and descriptors under the engine lock, release it, then perform admission checks and analysis. Require concurrent-call safety from the provider instance, or introduce an explicit per-invocation factory. Silently serializing a provider behind a mutex can break reentrant calls.

Finally, use one publication rule:

> A non-nil Go error means no usable result, regardless of any accompanying output.

Apply it to cancellation, incompatible requests, stale results, malformed provider output, and internal failures.

## 3. Horos migration, failures, coverage, and diagnostics

### The admission change must be real

The supplied `Pack → Check → bridge.Read → decl.Resolve` chain must change. Adding an analyzer above it cannot admit unresolved inputs.

Introduce an explicit new syntax profile, for example `horos-declarations-1/ast-2`, and a separately versioned analysis contract.

The split should be:

| Syntax admission | Semantic analysis |
|---|---|
| AST structure and permitted node shapes | Pinned dependency closure |
| Decodable declaration representation | Import and declaration resolution |
| Well-formed target syntax | Parameter scope and range |
| Allowed annotation shapes | Application arity |
| Syntax resource limits | Descriptor rules and recursion restrictions |
| Preservation of unresolved spellings | Correctness of supplied binding annotations |

A syntactically well-formed unknown declaration or out-of-range parameter use should be packable under the new profile and rejected by analysis.

Missing optional binding annotations should be accepted. A present annotation is an assertion to validate:

- Correct annotation: accepted.
- Missing optional annotation: accepted.
- Wrong but structurally valid target: semantic rejection.
- Malformed reference or forbidden annotation shape: admission rejection.

Do not repair supplied annotations.

The old `equivalent` procedure cannot continue requiring every expected resolved reference. Preserve appropriate syntax-shape checks, then separately verify every supplied semantic annotation against authoritative resolution.

Also audit `Model(view)` and callers that currently rely on `Pack` implying a checked `Program`. They must consume successful analysis or invoke an explicitly semantic operation. Otherwise the old admission guarantee survives indirectly.

### A potentially necessary Horos change is missing from the plan

The document correctly identifies the circular admission path, but there is another dependency: **`Expression.Operands()` is currently exposed by a checked Program.**

If unresolved AST construction needs to distinguish descriptor type positions from literal positions, those positions must be available before successful resolution.

There are two acceptable implementation paths:

- Reuse an existing authoritative raw syntax projection, if it already preserves those boundaries.
- Factor a narrowly scoped unresolved decoder or walker inside Horos.

Bittype must not recreate a descriptor-constructor switch or scan arbitrary tuple tags to bridge this gap. That would introduce the literal-capture bug before the analyzer runs.

This is worth resolving before estimating the migration. A structured-error extension alone may not be sufficient.

### Preserve the resolver’s all-or-nothing behavior

The smallest sound adapter is:

1. Decode the admitted unresolved input and supplied context.
2. Invoke Horos’s resolver.
3. On success, allocate entities from canonical module, declaration, and parameter identities.
4. Translate the three checked expression kinds.
5. Record bindings, expression associations, and direct dependencies.
6. Validate supplied annotations.
7. Publish only after the complete graph passes integrity and output-budget checks.

The translation should switch on the three declaration expression kinds. Descriptor traversal should use `Operands()`, without another constructor interpreter.

Never publish `descriptorAny()` as a term meaning. It is an internal validation skeleton, according to the supplied resolver excerpts.

Do not manufacture partial checked facts from a failed `Resolve`.

### Separate completion from verdict and graph availability

Replace `Completed []Key` with explicit reports for declared obligations.

For the first Horos contract, useful obligations are:

- Entry-program validity.
- Declaration and parameter indexing.
- Occurrence binding enumeration.
- Application and parameter-binding enumeration.
- Direct declaration-dependency enumeration.

Give every obligation a precise domain. Version one can use a fixed domain such as the entry module and its resolved transitive import closure. State whether unused supplied modules are merely available context or also part of the analysis obligation.

Use separate dimensions:

- **Status:** complete, incomplete, unsupported.
- **Verdict:** undetermined, holds, rejected.
- **Model availability:** complete graph present or absent.

A completed enumeration need not have a truth verdict. A completed validity decision can have either `Holds` or `Rejected`.

An early counterexample can establish that the submitted program is invalid while leaving binding and dependency enumeration incomplete. Consequently:

> `Complete(program-validity) + Rejected` does not imply `Complete(binding-enumeration)`.

Require every requested obligation to receive a report. An omitted required obligation must not silently count as completed.

For version one, I would prohibit usable `Conditional` judgments. Add them when assumptions have a structured representation and a consumer can reason about their discharge.

Also ensure that failure judgments can refer to the request or an input occurrence. They must not require `TermID` operands from a semantic model that does not exist after rejection.

### Use a precise failure matrix

| Situation | Result classification | Model |
|---|---|---|
| Resolution and complete export succeed | Validity holds; promised enumerations complete | Present |
| Definite bad declaration, arity, scope, or prohibited cycle | Validity rejected; exhaustive enumeration not claimed | Absent |
| Required pinned module is missing | Incomplete context; verdict undetermined | Absent |
| Name is absent from the correct, verified supplied module | Semantic rejection | Absent |
| Supplied dependency has a wrong pin or expected name | Submitted context rejected | Absent |
| Valid analysis budget is exhausted | Incomplete due to resources | Absent |
| Resolver succeeds but graph export exceeds its budget | Export incomplete; any retained validity verdict must be separately identified | Absent |
| Requested relation is unsupported | Unsupported report for that relation | May coexist with a complete graph for supported obligations |
| Options or limits are invalid | Request/configuration error | No usable output |
| Cancellation, stale result, corrupt graph, or internal failure | Go error | No usable output |

A wrong supplied dependency deserves particular care: it proves a problem with the submitted context. It does not prove that the root source is invalid under its proper dependency context. Make the rejected subject explicit.

Likewise, distinguish an admission resource limit from malformed syntax. Neither should become a semantic `Rejected` verdict merely because admission returned an error.

### Extend Horos’s structured failures at the authoritative sites

The adapter needs stable distinctions for at least:

- Invalid configuration.
- Missing dependency.
- Invalid or mismatched dependency context.
- Definite source semantic failure.
- Resource exhaustion.

More specific reasons—unknown alias, unknown declaration, arity mismatch, prohibited cycle—can accompany those classes.

Include relevant structured details and the failure subject. Populate module, declaration, and path at the point where Horos knows them.

Do not infer these distinctions by parsing messages or redoing lookup in bittype.

If the old coarse error API must temporarily remain the only source, ambiguous failures should produce no definitive semantic conclusion. That fallback is conservative, but it would not fully qualify the requested contract.

Diagnostic locations should map through the current input’s decoding index. When no precise mapping exists, return a module-level or unlocated diagnostic. Historical origins and fabricated root locations are poor substitutes.

For cancellation, checking before invocation and publication is straightforward. Prompt interruption requires cooperation from the underlying resolver; the wrapper cannot forcibly preempt synchronous trusted Go code.

## 4. Bounded implementation and qualification

### Implement in six bounded changes

| Change | Deliverable and exit condition |
|---|---|
| **1. Specify the contract** | Define the four consumer queries, term variants, annotation policy, coverage domain, failure matrix, and profiles. |
| **2. Make minimal Horos changes** | Add structured failure distinctions and any required unresolved traversal, preserving semantic ownership and all-or-nothing resolution. |
| **3. Build the generic boundary** | Implement immutable graphs, request-bound results, catalogue validation, budgets, and typed adapter construction. Qualify them with independent providers. |
| **4. Migrate Horos admission** | Remove semantic resolution from syntax `Check`, change the profile, and migrate dependent APIs and tests together. |
| **5. Implement the Horos projection** | Map checked identities and operand paths; verify annotations; expose the four queries through a generic-only consumer. |
| **6. Integrate through maintained gates** | Pin the Horos revision and run exact committed source through the prescribed remote architecture, fixture, race, and fuzz gates. |

That sequence respects the note’s prohibition on local portable compilation, tests, and formatting. No performance claim or copying optimization is needed.

Keep `semantics/go` independent of Horos, Ontos, and native language implementations. Put request/provenance wrappers in `contracts/analysis/go`, language adaptation in `analyze/horos/go`, and explicit composition in the provider.

### Make the acceptance tests discriminate between plausible wrong designs

| Test family | Required observation |
|---|---|
| **Admission separation** | An unresolved unknown name packs successfully, then analysis rejects it. |
| **Annotation semantics** | Missing optional annotation succeeds; wrong supplied target fails; input remains unchanged. |
| **Context classification** | Missing pin, wrong supplied pin, and absent name in a verified module produce distinct outcomes. |
| **Binder identity** | Pagination retains all four identities; same-spelled parameters and declarations do not capture one another. |
| **Term fidelity** | Zero-argument application differs from reference; exact values, scalar kinds, and record metadata survive. |
| **Literal boundaries** | Declaration-looking tuples inside literal positions produce no captured parameter or application edges. |
| **Coverage honesty** | Early rejection does not advertise complete binding or dependency enumeration. |
| **Exhaustion** | Resolver and exporter limits produce incompleteness, never semantic rejection. |
| **Provenance** | Returning request A’s result for request B fails, including when `TreeID` and graph shape are unchanged. |
| **Changed context** | Changing only dependencies, effective options, or limits cannot reuse the prior result. |
| **Ownership** | Mutating retained builder inputs or accessor-returned copies cannot change the published snapshot. |
| **Boundary errors** | Wrong marker/profile, duplicate tree IDs, historical-only sites, dangling IDs, and output-plus-error are rejected appropriately. |
| **Concurrency** | Concurrent and reentrant callbacks neither deadlock nor share mutable result storage. |
| **Extensibility** | An unfamiliar provider registers and answers the structural queries without kernel edits or Horos imports. |

Include valid open declarations and guarded recursive ground descriptors. Successful analysis must not imply specialization, and groundness must not imply the absence of permitted recursion.

The independent consumer should demonstrate the actual acceptance criteria: resolve an occurrence, inspect an application and its parameter bindings, traverse the specified dependency relation, and enforce required coverage and verdicts.

### A second production analyzer is unnecessary now

I would require a **test-only non-Horos provider before freezing the SPI**.

It should exercise:

- Alias versus defined-declaration metadata.
- Literal data such as different array lengths.
- Lifetime, type, and constant parameter sorts.
- Provider-specific operators.
- Builtins without declaration sites.
- Unknown language markers and vocabulary.

That qualifies structural extensibility. It does not establish actual Go, Rust, TypeScript, or C++ semantic support.

I would not make a production native analyzer a prerequisite for the first Horos release. A tiny real Go adapter becomes justified if a concrete interface decision remains unresolved and existing frontend/source mapping plus `go/types` can settle it cheaply. Otherwise, the trigger should be the first consumer requiring real native semantic answers.

## 5. What I would revise or defer

The most consequential architectural opportunity is to make **consumer obligations** the organizing principle of the API. A consumer should ask whether the facts it needs were established, rather than infer success from `err == nil`, a nonempty graph, or the absence of diagnostics.

For example, a transformation might require both:

- Complete application-binding enumeration.
- A `Holds` verdict for entry-program validity.

An inspection UI might accept a completed rejection with diagnostics and no graph. The same result contract can serve both without ambiguity.

I would also explicitly document three independent identities:

1. **Syntax address:** `TreeID` plus `NodeID`, interpreted within captured inputs.
2. **Snapshot-local semantic identity:** entity and term IDs.
3. **Provider identity:** for example, Horos declaration identity containing a module pin.

A module pin identifies an exact canonical document. It does not replace request provenance or establish semantic equivalence. Keeping these distinctions explicit prevents future caching and transformation work from building on the wrong identity boundary.

The remaining deferrals should have concrete triggers:

| Deferred capability | Trigger for adding it |
|---|---|
| Scopes and lookup explanations | A consumer needs visibility navigation or an explanation of resolution. |
| Common field/member projection | A consumer needs cross-language field enumeration or optionality. |
| Conditional judgments | A real analysis produces represented assumptions that consumers can inspect or discharge. |
| Inference, defaults, packs, and higher kinds | A supported provider needs to distinguish written syntax from effective bindings. |
| Partial checked graphs | Horos or another provider can return sound partial evidence with precise coverage. |
| Native equivalence, conversion, or assignability | A relation-specific provider and consumer require those answers. |
| Persistent semantic identity and caching | Measured reuse needs justify them, with an explicit context and equivalence contract. |
| Cross-language normalization | A deliberately scoped transformer defines and qualifies its preservation guarantees. |

Function, Stream, and Cell should pass through every layer as ordinary library declarations. Their names should introduce no privileged branches in the graph, engine, or generic validation.

I would consider the first release qualified when the four generic consumer queries work against a complete Horos export, the failure matrix is observable through structured results, unresolved inputs can actually reach the analyzer, and an independent provider passes the same ownership and integrity boundary. That is a useful production contract with clear limits—and enough room for later language support without committing the kernel to Horos’s type algebra.