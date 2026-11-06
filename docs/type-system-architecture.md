# Type-system reasoning: ownership and implementation design

**Evaluated design, 5 November 2026. Documentation only.**
[Research 0010](../research-docs/0010-type-system-model-architecture.md) applies the
expert review to the [type-system model](type-system-model.md),
[symbolic pagination proof](logos-preservation.md) and
[generic generation design](generic-elaboration.md).
The owner's subsequent [Logos source inspection](logos-verification.md),
6 November, refines the verification SPI boundary below. Its
[typed check-only follow-up](logos-verification.md#concrete-typed-check-only-contract)
now supplies a concrete Go integration target, landed in Logos PR 3881. This is
source-based design work, not a new expert consultation or
implemented Bittype API.
The subsequent [joint request/result contract](verification.md) defines preparation,
checking, ownership, lifetimes and failure behavior, with TextBox/Page walkthroughs.
It remains a documentation proposal; public package promotion and implementation
are separate from those contract decisions.

Use Logos's semantic-package direction directly. Language integrations supply
their meanings and judgments; mapping owners supply selected representation
relationships. Bittype associates the requested claim and evidence with the
actual generation operation. It does not gain a second logical foundation or a
universal type-system interpreter.

The first implementation should exercise two exact statements: TextBox admission
and result representation, and Page's conditional generic preservation. A narrow
Logos-owned typed checker supplies the checking prerequisite for **production
symbolic certification** once its delivery is admitted. This document supplies
neither that checker nor new providers, dependencies, runtime bindings or public
Go APIs.

## Ownership and dependencies

| Owner | Responsibility |
| --- | --- |
| Horos | Canonical declaration/kind semantics, original argument contracts, authoritative query procedures and checked public views |
| Native language integration | Language-specific semantic objects, judgments, configuration, artifact interpretation and supported coverage |
| Logos | Ground derivation checking and optional search/certification; landed typed check-only API; broader semantic-package composition remains proposed |
| Type-owned mapping integration | The selected static representation relation, native family/helper choices and precisely stated mapping laws |
| Bittype integration | Captured source/output association, actual mapping selection, the independently constructed expected statement and associated outcomes |
| Caller composition | Explicitly select and invoke the required stages, checking route and accepted assumptions |

Logos already has a [semantic-package proposal](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/docs/architecture/decisions/0038-typed-judgments-and-semantic-packages.md)
and [typed statement design](https://github.com/Bitspark/logos/blob/d94e8dc996ebba99e2bd26cde518df4918e2354f/docs/design/0014-typed-package-production.md).
They cover signatures/contexts, meaning, judgments, evidence, adequacy,
composition and operational services. Those concepts are the destination for
the seven requirements in the model document, not a new Bittype-owned framework.
The broader contracts remain proposed. The subsequent narrow typed check-only
delivery exposes a Go `Statement`/`Expr`/`Policy`/`Limits`/`Outcome` boundary and
the required checking operations without delivering the full package model.
Its [reviewed versus delivered status](logos-verification.md#concrete-typed-check-only-contract)
must remain explicit; no Bittype integration is implemented here.

The existing Go ground checker is usable independently: `contract.Check`
checks supplied positive derivations, and `certify.CheckWithDerivedContext`
binds closed-world checking to the caller's live inputs. Neither requires
Bittype to implement an EngineProvider. The
[production request trace and direct route](logos-verification.md#direct-checking-needs-no-search-provider)
separate these existing APIs from the typed proposal and identify the source,
native and artifact-association work Bittype still needs.

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

1. Admit the supplied source and semantic reports using existing contracts.
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

The private integration binds the following inputs; this is a responsibility
map, not a new Go struct or serialized schema:

| Input | Construction and ownership |
| --- | --- |
| Source and actual output captures | Existing Bittype snapshots retained separately by the proposed joint association |
| Source/native semantic results | Compatible provider-issued reports for those captures; not automatically proof evidence |
| Selected mappings, helpers and environment | Caller configuration and owning language/mapping integrations |
| Independently expected claim | Owning integration constructs the appropriate Logos ground claim or admitted typed statement from the bound artifacts |
| Candidate derivation/evidence | Supplied separately, in the selected Logos checking contract |
| Checking basis and limits | Explicit Logos profile or admitted typed receiver bindings/policy; logical assumptions remain distinct from permitted dependencies and budgets |

Claim and evidence representations belong to the selected Logos contract.
The [Logos contract map](logos-verification.md#ownership-map-for-the-verification-spi)
identifies the existing and proposed owners of those inputs. Bittype specifies
their artifact association and integration, without copying Logos's expression
calculus. The [joint contract walkthroughs](verification.md#textbox-walkthrough)
exercise the boundary for TextBox and Page without claiming completed integration.

### Proposed package roles

Preserve the [stage-first layout](../LAYOUT.md) and one root Go module. These are
future role locations, not directories or provider capabilities created now:

```text
contracts/verification/go/                 eventual request/result association, if needed
verification/files/go/                     explicit whole-claim orchestration
verification/files/source/horos/go/         authoritative source facts to statements
verification/files/target/golang/go/        actual native facts to model inputs
verification/files/checker/logos/go/        direct use/adaptation of selected Logos API

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
both examples identify its necessary shape. Contract walkthroughs can establish
that shape before both examples execute; production guarantees still require
their implementation and qualification. Generic orchestration depends only
on carriers and explicit collaborators. Concrete adapters may depend on their
own language/foundation contracts under explicit dependency rules.
The checker collaborator need not become a separate public provider SPI. Direct
composition with a supported Logos API is sufficient where it meets the claim.
No new verifier registry, foundation catalogue or automatic discovery is selected.

An optional Go package still affects its module's requirements if it imports
Logos. Package isolation is not module isolation. A separate external integration
module would need a concrete dependency justification and an explicit architecture
decision; this consultation does not split Bittype or create another repository.
The existing type-owned integration-module proposal remains independent of this
shared orchestration layout.

## Remaining public verification SPI decisions

**Refined after source inspection, 6 November 2026. Documentation only.** The
owner requested that the [Logos findings](logos-verification.md) narrow this
proposal. The four concerns below distinguish existing Logos capabilities,
Logos production gaps and Bittype's association responsibilities. The
[joint request/result proposal](verification.md) now defines the association,
lifetimes, preparation/checking split and failure behavior. The remaining public
decision is which demonstrated agnostic surface to promote from the first
private composition; it is not a missing proof calculus or undefined request
semantics. Implementation and qualification remain outstanding. Semantic
ownership, explicit composition and stage-first placement are unchanged.

The fixed boundaries are: languages own judgments, Logos owns proof machinery,
the expected statement is constructed independently of candidate evidence, and
verification binds actual source/output/mapping inputs. ASTTransformer continues
to return syntax. Resolution and analysis are explicit stages; the generic kernel
does not acquire named kinds, contract cases or hidden semantic operations.

### 1. Bind the owning claim and evidence contracts

**Concrete Logos contracts:** the ground route has `GroundLiteral`, `Derivation`,
`CheckProfile` and direct checking. The reviewed typed delivery has
`typed.Statement`, scoped contexts, `typed.Expr` and dependency-sensitive
`typed.Check`. Its delivery status is recorded in the
[contract map](logos-verification.md#concrete-typed-check-only-contract). The
broader typed-package proposal is a separate contract, not a required dependency
of this first integration.

**Contract decision:** the receiver-selected integrations supply input-associated
projections and procedure results. Preparation constructs the expected statement
without a candidate, then seals it with the actual artifacts, interpretation,
policy and limits. Owner-local references retain their models; fresh canonical
resolution requires reconstruction and fresh preparation/checking. See
[construction and preparation](verification.md#request-construction-record).
No Bittype proof-term algebra or serialization is required. Arbitrary bytes,
unqualified names and a producer-selected conclusion do not establish compatibility.

**Integration dependency:** admit the landed typed Go implementation in Bittype.
Claim/evidence syntax and the in-process call no longer need an invented Bittype
equivalent. The language projections and defined association still need
implementation; public promotion does not authorize another checking framework.

### 2. Immutable requests spanning source and target

Existing [semantic requests and results](../contracts/semantic/go/request.go)
are bound to one captured language context. A preservation request relates source
and target snapshots, their respective semantic results, selected mappings,
helpers, environment and the expected claim. The [joint contract](verification.md)
retains separate snapshot/receipt handles, admits each participant's evidence
under its own requirements, owns mutable construction collections and issues
results only through its checking path. Attaching one side's request identity
to the other side's result does not establish the association.

This is the main new Bittype-owned contract. The typed route retains the logical
statement, candidate and policy in its acceptance result; it does not capture
Bittype's artifacts or source locations. Keep logical assumptions, the receiver's
permitted dependencies and checking budgets distinct. Prevent evidence for another output,
mapping or context from being published as this invocation's result, while
preserving legitimate reuse under explicitly compatible dependencies. These are
functional inputs, not a new digest or provenance workflow.

Existing semantic results remain provider-issued reports. `Result.Valid()`
does not establish their semantic truth, and a `Holds` report is not automatically
a Logos proof. Admit authoritative procedure results under explicit assumptions
or provide checkable evidence; retain Unproved and exhaustion distinctly.

### 3. Provider compatibility and checking capability

**Concrete Logos contracts:** the ground route explicitly binds a `CheckProfile`
and consumer context. The typed route has exact foundation/checker identities,
receiver-owned `Policy` and explicit `Limits`. Caller composition selects the
actual installed checking dependency; an identity string does not select or
authenticate an executable. Search has a separate EngineProvider SPI.

**Bittype direction:** use the direct typed Go route for the first TextBox/Page
composition and bind its required inputs. The
[placement proposal](verification.md#placement-qualification-and-remaining-work)
keeps Logos-specific prepared requests/results in the concrete integration and
promotes only demonstrated agnostic association contracts. A direct check does
not require an EngineProvider or a second verifier registry. Do not invent an
independent compatibility descriptor duplicating Logos's contract. The generic
boundary must reject unsupported pairings without
interpreting Horos or Go judgments or trying a weaker checker.

Language-specific resolution and analysis retain their existing ID/Name-based
interfaces. A language ID alone cannot identify a cross-language checking
contract. No public checker-provider interface, automatic discovery or fallback
is selected merely because the eventual association boundary is called an SPI.

### 4. Result types and failure boundaries

**Concrete Logos contracts:** the ground `Check` returns Accept/Reject with codes,
including limit codes. The typed `Outcome` separates accepted, invalid, incomplete
and unsupported statuses with typed codes and logical locations. Its acceptance
retains the exact statement/candidate/policy and actual dependency closure.
Collection accessors are defensive; the Go acceptance is a value handle, not a
caller-mutable wrapper retained by the outcome. Use the selected API's actual
contract rather than flattening these results or duplicating its outcome types.

**Contract decision:** a result retains its prepared request, exact candidate and
actual Logos outcome. Only the direct checking path seals that result; it cannot
retag another outcome. Participant-qualified diagnostic sites retain Logos's
logical code/location when no source mapping exists. Preserve evidence
validity, the established judgment and execution completion separately. Reuse
the existing contracts; do not invent another proof outcome algebra or flatten
the result to a Boolean.

A valid proof establishes its requested statement relative to admitted
assumptions. A checked counterexample can refute a preservation claim; malformed
evidence cannot. Unsupported capability, exhaustion and cancellation do not
establish logical falsehood. The [failure contract](verification.md#result-and-failure-contract)
assigns malformed associations, unavailable preparation inputs, contract
violations and cancellation to errors without a joint result. Actual Logos
outcomes, including invalid candidates, are retained in a joint result. This
does not discard the underlying semantic reports or confuse incomplete source
evidence with refutation.
Logos's locations identify phases and declaration/context positions; Bittype
supplies the mapping to actual source/output sites. The synchronous checker has
no cancellation API, so cancellation of surrounding extraction or orchestration
must remain a separate operation-completion condition.

### Joint contract and implementation cases

The [joint contract](verification.md) specifies the construction record, sealed
request/result and concrete Logos bindings. Its walkthroughs identify which
source/native/mapping inputs are still missing implementations. Keep proof
production separate: checking receives an independently fixed expected statement
and supplied evidence, without hidden search or missing analysis. A new public
checker abstraction is not needed for this direct-call contract.

Use the two examples to review the candidate SPI before public promotion:

| Contract exercise | Required distinction |
| --- | --- |
| Concrete TextBox counterexample | Source nonmembership plus native acceptance and fixed representation can refute the claim; invalid evidence cannot |
| Symbolic Page and nested application | Actual binders, scoped assumptions, imported laws and instantiated result evidence remain associated |
| Changed output/helper, foreign snapshot or incompatible checker | Require fresh association; reject a stale result or unsupported binding. Candidate syntax may still pass a fresh check if the new meaning and laws justify it |
| Unsupported or exhausted checker | Publish the coverage/completion limitation without a stronger claim |

These remain Bittype contract-design cases, not completed Bittype integration
tests. The Logos follow-up has shared logical TextBox/Page cases, but no actual
Bittype artifact extractor or native model. The
[Logos typed-checker milestone](#the-logos-checking-prerequisite) remains an
execution dependency for production symbolic certification; it does not prevent
defining this contract. The proposal implements no provider, generator, proof
producer or checking service.

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

Use the narrow Logos-owned **typed, check-only milestone**, admitted through
Logos's own contract and implementation process. The
[follow-up inspection](logos-verification.md#concrete-typed-check-only-contract)
records its landed concrete API separately from the earlier
test-private Rust experiment and the broader semantic-package roadmap. This
prerequisite concerns production symbolic certification, not every Bittype
verification operation. The reviewed implementation supplies:

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

The selected first route is the in-process Go API
`github.com/Bitspark/logos/contract/go/typed.Check`. Bittype does not need to design
a Rust process protocol, service deployment or portable proof decoder for this
slice. The owner's shared/CI workload rule governs development executions; it
does not imply a permanent remote service or a new publishing/artifact workflow.

Today's ground Rule/Builtin/Leaf DAG can transport a theorem established by an
admitted external checker. It does not itself implement universal introduction.
RuleAdmissibility receives no full instantiated derivation and cannot check
scope discharge merely through its current signature.

Keep checking in the selected Logos contract. An unsupported claim must remain
unsupported; it does not authorize a second consumer certificate calculus,
unchecked trusted leaf or weaker automatic fallback.

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

1. Implement the private composition defined by the [joint contract](verification.md)
   against the [Logos contract map](logos-verification.md#ownership-map-for-the-verification-spi).
   Preserve the walkthroughs' exact TextBox/Page statements, artifact associations
   and evidence dependencies without designing another proof or checker registry.
2. Establish the concrete TextBox bool counterexample through the selected typed
   Go checking route and authoritative source refutation. The existing ground
   route remains usable, but this slice does not require two integrations.
3. Model the selected Go fragment's exact/underlying acceptance characterizations
   and field representation. Retain the example's explicit Go 1.25 baseline with
   matching compiler/specification inputs; the moving live specification now
   describes Go 1.27. This document ran neither toolchain.
4. Admit the delivered Logos typed Go dependency. Prepare the precise Horos
   schema bridge and language/mapping inputs.
   Logical conformance cases do not establish artifact or native-model adequacy.
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
| Remove a field, alter optionality or change a helper beyond the admitted law | Fresh preparation/checking rejects stale representation evidence or reports missing coverage |
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
