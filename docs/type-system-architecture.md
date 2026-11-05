# Type-system reasoning: ownership and implementation design

**Evaluated design, 5 November 2026. Documentation only.**
[Research 0010](../research-docs/0010-type-system-model-architecture.md) applies the
expert review to the [type-system model](type-system-model.md),
[symbolic pagination proof](logos-preservation.md) and
[generic generation design](generic-elaboration.md).

Use Logos's semantic-package direction directly. Language integrations supply
their meanings and judgments; mapping owners supply selected representation
relationships. Bittype associates the requested claim and evidence with the
actual generation operation. It does not gain a second logical foundation or a
universal type-system interpreter.

The first implementation should exercise two exact statements: TextBox admission
and result representation, and Page's conditional generic preservation. A narrow
Logos-owned typed checker is the recommended prerequisite for **production
symbolic certification**. This document supplies neither that checker nor new
providers, dependencies, runtime bindings or public Go APIs.

## Ownership and dependencies

| Owner | Responsibility |
| --- | --- |
| Horos | Canonical declaration/kind semantics, original argument contracts, authoritative query procedures and checked public views |
| Native language integration | Language-specific semantic objects, judgments, configuration, artifact interpretation and supported coverage |
| Logos | Proof-language contexts, binding, substitution, evidence checking and justified semantic-package composition |
| Type-owned mapping integration | The selected static representation relation, native family/helper choices and precisely stated mapping laws |
| Bittype integration | Captured source/output association, actual mapping selection, the independently constructed expected statement and associated outcomes |
| Caller composition | Explicitly select and invoke the required stages, checking route and accepted assumptions |

Logos already has a [semantic-package proposal](https://github.com/Bitspark/logos/blob/d0df8bbb6914e8002a441cb79241da8c9a3db7b0/docs/architecture/decisions/0038-typed-judgments-and-semantic-packages.md)
and [typed statement design](https://github.com/Bitspark/logos/blob/d0df8bbb6914e8002a441cb79241da8c9a3db7b0/docs/design/0014-typed-package-production.md).
They cover signatures/contexts, meaning, judgments, evidence, adequacy,
composition and operational services. Those concepts are the destination for
the seven requirements in the model document, not a new Bittype-owned framework.
The contracts remain proposed; their richer checker is test-private Rust, not a
stable Go package or deployed service that Bittype can use today.

The dependency direction is:

```text
explicit caller composition
  -> existing Bittype syntax/semantic/generation stages
  -> optional generation-verification composition
       -> explicit source and native integrations
       -> selected mapping laws
       -> admitted Logos checking adapter

source integration -> authoritative Horos API
language and mapping theories -> Logos statement/checking contracts
```

Horos core does not import Bittype or acquire a mandatory Logos dependency.
Logos does not import Bittype or learn Horos, Go, Page, Data, Function, Stream or
Cell as kernel meanings. The general Bittype engine and carriers do not import
concrete theories. Declarative theory and mapping assets remain with their
semantic owners; a directory containing a rule does not establish its truth.

## Four assurance boundaries

Keep four different claims visible:

1. **Object-language meaning:** what a Horos kind, Go constraint, Rust trait
   obligation or TypeScript assignment actually means.
2. **Foundation soundness:** why the admitted proof rules justify conclusions.
3. **Evidence validity:** whether this evidence proves this exact statement under
   the stated dependencies and assumptions.
4. **Model adequacy:** whether the modeled objects and judgments describe the
   actual artifacts and promised native profile.

Proof checking establishes the third relative to the first two. It cannot correct
an extractor that omitted a field or confused a Go alias with a defined type.
There are three correspondence boundaries to justify: authoritative Horos to
source semantic objects, actual generated artifacts to native semantic objects,
and modeled native judgments to the promised specification/compiler profile.
Each needs its own argument or explicitly accepted assumption.

For a promise about actual compiler-accepted applications, the crucial direction is:

```text
actual native acceptance within the declared domain
  => modeled native acceptance
  => source admission and the selected representation premises
```

Positive examples that compile do not establish this coverage direction. Also
distinguish compiler acceptance from specification well-formedness: some C++
ill-formed programs require no diagnostic. A claim must say which notion it
uses; compiler success cannot silently substitute for specification validity.

## Binding and the expected statement

An invocation binds the actual source capture, actual target capture, relevant
dependencies/environment, selected mappings, intended interpretation and requested
claim. The claim fixes its argument domain, checking point, permitted consumer
operations and extension policy before evidence is considered.

The expected statement is constructed by the selected integration from those
inputs, independently of candidate evidence. A proof cannot choose a weaker
statement, omit a generated helper or redefine the interpretation to make itself
valid. Module/declaration identities, original parameter slots, actual target
fields and mapping dependencies must remain associated.

Three identities are related, not collapsed:

- Horos's real declaration parameter.
- The native declaration's actual formal parameter.
- The variable in the proof context.

Their correspondence connects the owning languages' opening/application operations
to proof-language substitution. Same spelling or a common Binder struct does not
establish it. Extending a logical context with an assumption is different from
extending a native environment with a specialization, trait implementation or
merged declaration; the latter can change the judgment.

Reuse existing snapshots and canonical identities. A newly captured equal-looking
tree is not the same snapshot. Identical canonical modules admitted in different
Horos Programs do not authorize handle exchange. Replay resolves the supplied
canonical inputs and reconstructs local handles.

These are functional statement inputs and semantic assumption dependencies. This
design requires no new universal digest, persistent manifest, checksum,
administrative audit or provenance workflow. A content identity does not prove
truth. An AST claim covers its captured AST; a claim about emitted files must also
state or establish formatter/printer correspondence with the actual files.

## Explicit composition

Keep the existing ASTTransformer boundary initially. An explicitly selected
composition performs:

1. Admit the supplied source and semantic evidence using existing contracts.
2. Fix actual mapping choices, including helper/import dependencies.
3. Invoke the configured transformation.
4. Capture its actual output collection and relevant helpers.
5. Explicitly invoke any required native resolution/analysis stages that exist
   for the chosen fragment; missing capability is not filled by hidden analysis.
6. Build the expected statement and ask the selected checker to assess evidence.
7. Publish the outcome associated with this invocation and its assumptions.

Evidence about one generation leg is not automatically target semantic evidence
or evidence about the next leg. The current engine deliberately discards a
previous leg's evidence when capturing the next input. This remains correct.

An illustrative **private integration shape**, not a public API, is:

```go
type CheckRequest struct {
    Source    semantic.Snapshot
    Output    semantic.Snapshot
    Evidence  []semantic.Result
    Binding   GenerationBindingRef
    Expected  StatementRef
    Candidate EvidenceRef
    Policy    AssumptionPolicyRef
}
```

The Ref types above stand for bound integration-owned references, not defined Go
types, string tokens accepted on trust, or a new serialized schema. Statement
and proof representation follow the admitted Logos contract; Bittype does not
copy its proposed expression calculus. The two worked statements should exercise
this seam before a public `contracts/verification` API is designed.

### Proposed package roles

Preserve the [stage-first layout](../LAYOUT.md) and one root Go module. These are
future role locations, not directories or provider capabilities created now:

```text
contracts/verification/go/                 eventual shared association boundary
verification/files/go/                     explicit whole-claim orchestration
verification/files/source/horos/go/         authoritative source facts to statements
verification/files/target/golang/go/        actual native facts to model inputs
verification/files/checker/logos/go/        selected Logos checker adapter

resolve/files/golang/go/                    native resolution when required
analyze/files/golang/go/                    explicitly supported native judgments
integration/go/                            TextBox/Page composition qualification
```

Files is the dimension because a claim can cover multiple generated units,
dependencies and a module environment. Source, target and checker are
collaborators within that role, like the existing proposed generation
collaborators; they are not new carrier dimensions or hidden providers. Native
semantic extraction belongs to its native semantic stages, not another analyzer
buried in generation. Later languages follow the same matrix when implemented;
no empty stubs, registry cells or compatibility aliases are added.

Begin with a private composition seam; promote the shared contract only after
both examples identify its necessary shape. Generic orchestration depends only
on carriers and explicit collaborators. Concrete adapters may depend on their
own language/foundation contracts under explicit dependency rules.

An optional Go package still affects its module's requirements if it imports
Logos. Package isolation is not module isolation. A separate external integration
module would need a concrete dependency justification and an explicit architecture
decision; this consultation does not split Bittype or create another repository.
The existing type-owned integration-module proposal remains independent of this
shared orchestration layout.

## Outcomes and trust

Reuse the existing separation of completion and verdict, with distinct named
obligations where needed. Do not give a single Rejected result two meanings.

| Observation | Meaning |
| --- | --- |
| Valid proof under accepted assumptions | Requested claim established relative to that basis |
| Checked admissible counterexample | Requested claim refuted within the stated domain |
| Wrong or malformed evidence | Evidence rejected; the claim remains undetermined |
| Horos Unproved or exhausted checking | Incomplete; no negative language judgment follows |
| Unsupported construct, operation or checker | Unsupported, with no stronger result |
| Compiler/model disagreement | Adequacy of the affected profile is unresolved |

Current Logos LeafVerifier and RuleAdmissibility hooks are binary. If a bridge
uses them, outcome-bearing checking occurs before a leaf admits the exact
checker-issued result for immutable matching inputs. A producer's verified flag
is not such a result. Leaf rejection does not establish source nonmembership.
The existing Theory/GroundDecider route already distinguishes noncompletion;
do not flatten it into the binary-hook limitation.

Trusted Horos procedure results and independently checked composition proofs
have different roles. Importing the former as a precise assumption does not
independently verify Horos. A pure Theory callback is not inherently a stronger
proof than a leaf: strength depends on the actual checking contract and basis.

## The Logos checking prerequisite

Recommend a narrow Logos-owned **typed, check-only milestone**, admitted through
Logos's own contract and implementation process. Bittype documentation cannot
activate proposed Logos semantics. The needed capabilities are:

- Checked signatures and ordered contexts.
- The dependent-function/pair fragment needed for quantification, implication
  and paired results, plus the justified equality/conversion operations used.
- Abstraction, assumption discharge, capture-avoiding substitution and theorem application.
- Checking against the independently supplied exact expected statement.
- Explicit assumed laws, procedure dependencies and their reachable use.
- Bounded checking with invalid, incomplete and unsupported outcomes kept distinct.

General proof search, automatic package discovery, complete native theories and
the full Logos roadmap are not prerequisites for this small checker. A scripted
proof producer can supply the first evidence. Supported claims and admitted
package imports must still be explicit; a manifest or registration is not proof.

A thin Go process client is a candidate if the admitted checker remains a Rust
executable. Its versioned interface and resource contract must then be agreed
with Logos. A supported Go checker API could use the same ownership boundary.
The process protocol and deployment are not selected or implemented here. The
owner's shared/CI workload rule governs development executions; it does not imply
a permanent remote service or a new publishing/artifact workflow.

Today's ground Rule/Builtin/Leaf DAG can transport a theorem established by an
admitted external checker. It does not itself implement universal introduction.
RuleAdmissibility receives no full instantiated derivation and cannot check
scope discharge merely through its current signature.

A separate consumer certificate calculus is a possible alternative only when
an explicit delivery requirement justifies it. It would need a certificate
language, binder/substitution rules, dependency control, resource behavior and a
soundness argument for every lifted step. It is not selected now and must never
be an automatic fallback when the typed route is unsupported.

## Statement one: TextBox

With the ordinary prelude imported as h:

```text
type TextBox<T: h.Text> = { value: T; };
type Label = TextBox<text>;
```

Fix the intended associations text -> Go string and bool -> Go bool. Fix a native
environment E and a domain D_E before assessing either candidate. A useful first
fragment has closed concrete string/bool types, finite alias chains and
nonrecursive defined types over them; no symbolic forwarding or arbitrary
extension is silently included.

```go
// Candidate A
type TextBox[T any] struct { Value T }

// Candidate B, assessed separately
type TextBox[T interface{ string }] struct { Value T }
```

For A, a counterexample needs all three established premises: Go admits bool,
the fixed interpretation maps it to canonical bool, and Horos refutes
`bool : h.Text`. An unproved source query or unmapped native type is insufficient.

For the declared Go fragment, specify whole-domain characterizations:

```text
Accept_any(t)        iff t in D_E
Accept_exact(t)      iff t in D_E and Identical_E(t, string)
Accept_underlying(t) iff t in D_E and Identical_E(Underlying_E(t), string)
```

The exact constraint admits string and aliases of that identity; a defined
NamedString is different. `~string` admits a larger domain and needs an applicable
interpretation/representation relation for those types. It is not inherently
unsound: missing mapping support is not established source nonmembership.
[Go's generics description](https://go.dev/blog/intro-generics) explains this
exact/underlying distinction.

Admission still leaves a separate result-representation obligation for the actual
Value field, its required presence and binder substitution. Relative completeness
is assessed over the independently promised source domain, not every possible
member of the ordinary Text kind.

## Statement two: symbolic Page

```text
type Page<T> = { items: list<T>; next?: text; };
type DataPage<D: h.Data>: h.Data = Page<D>;
type Person = { name: text; active: bool; };
type Nested = DataPage<DataPage<Person>>;
```

Let H_C and G_C be the supported source/native proper-type domains for the bound
configuration C, K the particular resolved ordinary kind, P_H the actual source
Page family and P_G its actual generated Go family. The expected statement is:

```text
forall T in H_C, N in G_C:
  Member(T,K) -> Rep_C(T,N) ->
    (Member(P_H(T),K) and Rep_C(P_H(T),P_G(N)))
```

Formation requirements belong in those domains or explicit premises. K is fixed
by this statement; this is not a theorem over every kind named K. T and N are
arbitrary qualifying types, not entries in a finite catalogue.

### The authoritative source guarantee

Horos already supports the required constructor query. In public API shape:

```go
page, err := program.Value(pageDeclarationID)
// Handle err. canonicalArrow is the ordinary resolved kind expression K => K.
arrow, err := program.CompileKind(moduleScope, canonicalArrow)
// Handle err; only a Proven verdict supports the requested guarantee.
verdict, err := program.Satisfies(page, arrow)
```

This is an explanatory sequence, not production code or an executed example.
Horos checks the requested input kind against every original bound, applies the
existing constructor under real rigid parameters, then checks its retained body
against the requested result kind. Its [kind specification](https://github.com/Bitspark/horos/blob/ad7c6f493d726a82e445875d42af319adce0e3ce/docs/spec/kinds.md)
explicitly supports querying the same Page body for Data => Data and other
particular contracts.

An existing admitted DataPage declaration and its conditional Membership facts
provide another route. Public declaration/expression accessors already expose
its body, resolved callee, parameters, contracts and assumptions. The bridge must
check that actual association; it must not invent a general definitional-equality
proof or equate distinct declaration IDs by unfolding.

Neither route requires a mandatory helper declaration, a copied Page, template
lifting or a new inspection API. The **missing** integration translates the
authoritative result into a precisely scoped source-only schema in the selected
evidence foundation. All original assumptions, declarations and dependencies stay
visible. The current verdict/handle is not itself a portable proof, and several
closed Satisfies calls cannot establish a universal schema.

Initially this is a trusted Horos procedure result. The composition checker proves
the cross-language theorem relative to that source basis; it does not independently
prove Horos correct. A source schema is not an axiom asserting the final native
preservation conclusion.

### Representation and selected laws

For the first slice, Rep_C means the selected **static type/schema association**,
including relevant native well-formedness. It does not assert an adapter, runtime
value bijection or conformance of arbitrary native values.

The proof connects actual named families, opening/substitution, all paired fields,
the selected list/slice law and the optional helper. The existing example is:

```go
type BittypeOptional[V any] struct { Present bool; Value V }
type Page[N any] struct {
    Items []N
    Next BittypeOptional[string]
}
```

The names alone establish nothing. Do not replace a native named type with its
anonymous body using fabricated native type equality. The mapping's representation
relation and named-family law must justify the association while preserving native
identity distinctions.

Each law has an exact proposition and checked evidence or an assumption accepted
by the requester. Builtins use the same path. An imported Page's list, optional
and helper choices are fixed dependencies of its theorem instance; an element
override does not silently replace them. Source membership evidence and
configuration-specific representation evidence have separate reuse conditions.

A later value-denotation relation could be many-to-one, for example for inactive
optional payloads, nil/empty slices or distinct native names. That requires its
own stated laws. No such runtime relation follows from this static theorem.

### Checked abstraction and nested use

The proof shape is a real abstraction checked by the foundation:

```text
lambda T, N, member_T, rep_T_N.
  pair(
    source_page_guarantee(T, member_T),
    representation_of_actual_named_families(T, N, rep_T_N, selected_laws)
  )
```

Every operation has a justified type. The checker verifies scope, sorted variables,
explicit assumptions, lawful imported results, capture-avoiding substitution and
the exact final quantified type. A fresh ground atom, a spelling test or a
producer's parametric Boolean is not generalization. With checked abstraction
and uniform imported laws, this proof needs no additional global parametricity
theorem. Proof-language functions do not add source lambdas or a privileged
Function contract to Horos/Bittype.

Apply the theorem to established Person membership and representation. Its pair
of results supplies membership and representation for the inner Page; those
results become the premises for the outer application. The
[Introduce/Apply composition](generic-elaboration.md#one-generic-convention-for-formals-and-applications)
retains real binder association and whole-application result support.

An actual DataPage use still obeys that declaration's original argument contract
and its justified connection to its body. Ignored parameters remain constrained.
A weaker requested result does not relax source admission.

### Native admission remains independent

An enforcing target must account for the complete admitted application. In the
[Go certificate convention](go-native-admission.md#guarantee-and-scope), its
arguments are `(S, N, E)`: S is a native source-identity token, N the selected
native payload type, and E a native certificate type connecting that exact pair.
E is not a runtime evidence value or, by itself, a checked Logos proof.

Fix the source world, resolved kind K, selected mappings, issuing package and
consumer operations in configuration C. The admission obligation is:

```text
for every (S,N,E) in the declared native argument domain:
  NativeAccepts_C(S,N,E) ->
    exists T:
      Denotes_C(S,T) and Member(T,K) and Rep_C(T,N)
```

Denotes_C fixes the intended canonical subject of S independently of the proof;
it cannot select whichever T makes the conclusion convenient. The existence
obligation includes establishing that interpretation for every accepted tuple.
No inverse from N to T is assumed: different source identities can share a native
payload. The admission argument must justify why the actual E is an authorized
leaf or constrained composition for this S/N pair under C, including inference,
aliases and the promised checking point.

This tuple is a selected Go convention, not a mandatory carrier or parameter
scheme in the agnostic kernel or other languages. Their models must likewise
include all inputs relevant to admission. Compose the admission theorem with
Page's conditional theorem using its supported source/native domains and
formation premises for the end-to-end guarantee. Page[N any] alone
supplies no native Data gate, and admission followed by alias erasure does not
retain source identity in the resulting payload.

The Go proposal's PageEvidence requires evidence for its input and supplies
evidence for the resulting Page application, enabling nested use. Its compiler
fixtures exercise the native encoding; they do not establish the source and
representation laws authorizing issuance or the general admission theorem.

An unmapped native argument exposes missing coverage, not source nonmembership.
A Data enforcement counterexample requires an already fixed interpretation that
Horos actually refutes against Data. Bool is in Data, so TextBox's bool
counterexample cannot be reused unchanged for DataPage.

#### Go marker-method scope extension

The [Go encoding](go-native-admission.md#encoding-the-admission-relation) needs
empty private marker-method declarations on its certificate types. Producing
those method bodies is an explicit proposed extension beyond the original
types/interfaces-only generation scope. The methods satisfy the native admission
constraint; no certificate value needs to be passed or method called.

The selected Go convention owns these declarations alongside its certificate
types. Future emission must use Go-owned AST/CST tooling and include the actual
methods and issuing package in the claim's artifact association. This does not
introduce Data, Page or certificate semantics into the general generator or
Logos kernel. It does not authorize protocol operations, runtime adapters,
conversion or validation behavior.

These methods are implemented only in independently authored compiler fixtures.
Documenting the extension does not implement or activate production method
emission; the current types-only generators and Unsupported enforcement profile
remain unchanged.

### Recursion and later certificates

Keep recursive-kind reasoning inside the authoritative Horos procedure initially.
This preserves its greatest-fixed-point semantics while allowing a finite
composition proof relative to the source guarantee. Includes can return Unproved
without an admitted source counterexample; do not convert that into negation.

Independent membership certificates are deferred until Horos defines the actual
certificate judgment, guards, local obligations, admitted domain and checker.
A finite post-fixed-point relation may be useful, but pair revisitation alone is
insufficient. Inclusion's cuts and larger observation universe require separate
coverage and adequacy; a membership technique does not automatically solve it.

## Other native languages

| Language | First abstraction stress case | Required distinction |
| --- | --- | --- |
| TypeScript | A formed generic or conditional type normalizes to never; bottom remains assignable | Formation, normalization, assignment and inhabitation are different judgments |
| Go | A private witness is inferred or forwarded into another generic | Nameability, obtainability and admission are different properties |
| Rust | Alias use versus demanded associated-type normalization, including the effective Sized gate | Alias bounds, trait obligations, normalization and checking point differ |
| C++ | Change a supporting specialization or constraint-evaluation environment | Extension policy, visibility and semantic context are inputs to the judgment |

The [TypeScript compatibility reference](https://www.typescriptlang.org/docs/handbook/type-compatibility.html)
distinguishes assignment and subtype compatibility and documents bottom
assignability. The [Rust alias-bound lint](https://doc.rust-lang.org/rustc/lints/listing/warn-by-default.html#type-alias-bounds)
describes bounds not enforced at ordinary alias use sites, with feature-dependent
alternatives. The [C++ constraint rules](https://eel.is/c++draft/temp.constr.atomic)
include inconsistent satisfaction cases that require no diagnostic. Bind the
actual compiler/version/options/edition/features when qualifying a claim.

These cases add language-owned judgments and coverage, not cases in the generic
engine. The current AST subsets do not represent all these constructs. Extend
the owning syntax/semantic provider only when the next admitted fragment requires
it; a table of future cases is not an implemented native analyzer.

## Implementation and qualification sequence

After implementation is authorized:

1. Define the two exact statements, domains, actual artifact associations and
   evidence dependencies above. Exercise a private seam before a public SPI.
2. Establish the concrete TextBox bool counterexample through existing ground
   checking capabilities and authoritative source refutation. No sampling theorem.
3. Model the selected Go fragment's exact/underlying acceptance characterizations
   and field representation. Retain the example's explicit Go 1.25 baseline with
   matching compiler/specification inputs; the moving live specification now
   describes Go 1.27. This document ran neither toolchain.
4. Obtain the Logos-owned typed check-only contract and admitted implementation.
   In parallel, prepare the precise Horos schema bridge and language/mapping
   packages. A research/conformance example is not production certification.
5. Supply a scripted Page proof, check its actual quantified statement, then
   instantiate it for Person and again for the nested result. No general search
   engine or finite argument whitelist is required.
6. Qualify source projection, target projection, native-model adequacy and any
   emission promise. Then exercise the other language distinctions and promote
   only the shared API shown necessary by both examples.

The future regression set includes:

| Change or input | Required observation |
| --- | --- |
| TextBox[bool] under any with the fixed interpretation | Checked counterexample to TextBox preservation |
| NamedString under the exact string term | Native rejection in the admitted fragment |
| Change exact string to ~string | Earlier evidence loses association; reassess the larger domain |
| Remove a field, alter optionality or change a helper | Representation evidence rejected |
| Same-spelling foreign binder or escaping assumption | Binding/evidence rejection |
| Generalize a fact depending on one chosen rigid atom | Evidence rejection |
| Drop an ignored parameter's original bound | Source admission still required |
| Change a fixed imported mapping between nested uses | Old representation evidence cannot be reused |
| Foreign Horos handle, wrong snapshot or weaker expected statement | Input/evidence rejection |
| Exhaust a source/checker budget | Incomplete, never logical falsehood |
| Justify a recursive cycle solely by revisitation | Certificate rejected |
| Actual native acceptance outside modeled acceptance | Adequacy defect; no unqualified result |

These are acceptance requirements, not passed tests. Relevant automated tests
and native typechecking accompany implementation using maintained shared/CI
execution. All Logos workloads remain off the workstation. No extra administrative
validation system or speculative performance target is introduced.

## Unsupported conditions and forbidden shortcuts

Production symbolic certification remains unavailable until the required checking
contract/executable and source/native bridges are admitted and qualified. Missing
abstraction/substitution checking, absent actual-output correspondence, lost
original constraints, unsupported recursive obligations, uncovered native
operations or unresolved adequacy defects prevent the stronger result.

Do not add a Bittype-owned foundation, universal type enum, closed judgment enum,
Data/Type case or named-contract switch. Do not treat mapping selection as proof,
invalid evidence as refutation, an assumed final theorem as derivation, a fresh
atom as universal introduction, or a private constructor name as native issuance.
Do not hide resolution, carry evidence across incompatible outputs/configurations,
or fall back silently to erased constraints. Current representation generation
and its Unsupported enforcement profile remain unchanged.
