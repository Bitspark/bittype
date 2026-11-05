# Advice for consult
**Research doc:** doc_3aad1791-02c0-47bf-b286-3d0ffe045f41.md
**Document:** doc_3aad1791-02c0-47bf-b286-3d0ffe045f41
**Advisor:** ChatGPT, model "6 Pro" as labelled on the answer page
**Session:** nsess_b68f9798-d04e-4799-8160-0fb33240c869
**Run:** run_aa6f8ad3-4bbc-4b25-a894-fda48f62c8e5
**Chat:** https://chatgpt.com/c/6aecb77d-eb90-83ed-ba80-f65f60403080
**Verify:** nightfall consult verify run_aa6f8ad3-4bbc-4b25-a894-fda48f62c8e5 --output <this file>

# Response to Research 0008

**I recommend a shared generic convention that pairs formal introduction with application elaboration, together with a scoped representation product that carries more than the native payload.** The convention should be owned by the explicitly configured source/target integration. The generation kernel should manage its identity, scope, reservation and lifetime.

Parameter expansion is one possible implementation of that convention. A witness with a representation projection is another. Neither should become the universal model.

The specific change I would make to the proposal is:

- Replace `source binder → target parameter index` with **`source binder → scoped representation product`**.
- Give each represented open declaration **one convention shared by `Use`, `Publish` and `Define`**.
- Make that convention responsible for native formals and bounds, actual arguments, and the **result’s representation association**.
- Admit source applications independently of which arguments their payload generation visits.
- Require explicit contracts for imported generic artifacts, including representation dependencies fixed inside their definitions.

This is a design review of the supplied baseline and proposed contracts in [Research 0008](https://chatgpt.com/library?file_id=libfile_b489cbb406c4819188102ba898367349&account_id=personal). The supplied transformer bridge, semantic facts and real binder identities provide an appropriate foundation. I would keep generated representation associations out of the source semantic graph.

## 1. The conceptual decomposition

Four semantic responsibilities need to remain distinct, with the kernel handling coordination.

| Responsibility | Owner | What it establishes |
|---|---|---|
| Source admission | Horos provider and source adapter | The effective application is admitted, and particular memberships hold under particular assumptions. |
| Representation choice | Exactly selected mapping | Which native family represents a source declaration or builtin under the selected profile. |
| Generic elaboration | Configured source/target convention | How native formals, actual arguments, payload projections and result associations correspond. |
| Native construction | Target integration | Structured native syntax, scopes, imports, anchors and declaration assembly. |
| Coordination | Generation kernel | Selection, invocation identity, scope checks, immutable requests, reservation, work limits and all-or-nothing output. |

The third responsibility is the missing connection. It belongs in a shared integration contract because both sides of a generic boundary must agree on it.

### Why entirely private plugin logic is insufficient

Keeping the kernel payload-only can work **if a shared source/target collaborator already handles the complete convention**. That is a valid placement of the abstraction.

What does not work is independent private bookkeeping in each mapping. If one mapping treats a parameter as a witness with a projection, another expects a payload plus evidence, and a third recognizes evidence by a private package identity, exchanging only `Type[T]` loses the information needed to establish compatibility.

Independent authorship is compatible with a common protocol. Arbitrarily chosen protocols do not automatically compose.

I would initially configure one typed association protocol per transformer invocation. Different declarations can still use different native parameter shapes through that protocol.

### What witnesses are actually necessary for

The brief’s `Leaf` example gives an information-loss argument.

Suppose two admitted source descriptions have the same native payload, but differ in membership:

\[
\operatorname{Payload}(A)=\operatorname{Payload}(B),
\qquad A\in K,\quad B\notin K.
\]

An admission rule that observes only the payload cannot distinguish them. To preserve that distinction, an enforcing encoding needs additional recoverable information, or a domain restriction that excludes the ambiguity.

This establishes a need for **discrimination**, not necessarily a separate witness parameter. The distinction might already reside in a nominal native type, an associated-type carrier, an opaque package of associations, or explicit support parameters.

Accordingly, I would reject a core contract with mandatory fields named `SourceWitness`, `NativePayload` and `MembershipProof`. Those names would prematurely choose an encoding.

### The law that makes the convention coherent

Let a declaration’s convention introduce formal products \(q_b\), one for each real source binder \(b\). For an admitted substitution \(\sigma\), let its application operation produce the corresponding native actuals \(\theta(\sigma)\).

The essential binder law is:

\[
q_b[\theta(\sigma)]
\;\simeq_{\Pi}\;
\operatorname{Represent}_{\rho,\Gamma}(\sigma(b)).
\]

Here:

- \(\rho\) is the fixed effective mapping configuration.
- \(\Gamma\) is the admitted source assumption context.
- \(\Pi\) is the representation promise.
- The left side substitutes native actuals into the formal product.
- Equivalence includes the payload projection and the associations required by the profile.

There must also be a result law: applying the translated generic must agree with representing the source application, including its result identity, association and supported conditional memberships.

This is stronger than agreement on native type text. It also requires a stable supported domain, coherent argument ordering and preservation of assumption dependence.

It is **not** a requirement to unfold every named declaration or make a nominal wrapper identical to its body. The profile must state the permitted representation relation, including whether alias interchangeability is promised.

These laws are obligations on trusted integrations. The kernel cannot establish them by comparing syntax or identifiers.

## 2. A proposed Go contract

I would use a third type parameter, `C`, for an integration-defined support carrier. That gives the common protocol a concrete Go type without imposing a common evidence algebra.

For an erasing integration, `C` can be a deliberately limited carrier. For a Horos-to-Go enforcing integration, it might be an immutable association handle supplied by that integration. It must not be a string-keyed collection of arbitrary values.

The following is proposed API notation; method bodies and copying constructors are omitted.

### Representation results and admitted uses

```go
// C is an immutable carrier defined by the configured integration.
type Mapped[T, C any] struct {
    Payload Type[T]
    Support C
}

// Created by the host after association and provenance checks.
type Product[F, T, C any] struct {
    /* opaque invocation-owned state */
}

func (p Product[F, T, C]) Input() Ref[F]
func (p Product[F, T, C]) Payload() Type[T]
func (p Product[F, T, C]) Support() C

// Created by the source adapter from admitted source evidence.
type AdmittedUse[F any] struct {
    /* exact effective use, judgments and assumption context */
}

func (u AdmittedUse[F]) Input() Ref[F]
func (u AdmittedUse[F]) Arguments() []Argument[F]
```

The source session gains an operation with this meaning:

```go
AdmitUse(context.Context, Ref[F]) (AdmittedUse[F], error)
```

This operation exposes coverage already established by the admitted provider results, including supported substitution of their rigid assumptions. It does not run analysis or prove additional judgments.

That distinction matters because admitting a receipt for the correct snapshot is not itself proof of every contextual judgment generation might request. The supplied `DirectParameterBindings` also explicitly establishes correspondence, not applicability. If the required judgment is unavailable, generation fails and the caller must supply an explicitly enlarged provider contract and receipts. [Research 0008, Appendix A](https://chatgpt.com/library?file_id=libfile_b489cbb406c4819188102ba898367349&account_id=personal)

### Association checking must be explicit

A host must not take arbitrary `Mapped{Payload, Support}` and make it trustworthy merely by attaching the requested source reference.

Otherwise, a mapping could accidentally return Person’s support beside `Cell<Person>`’s payload, and the host would falsely relabel the pair.

Add one explicitly configured association protocol:

```go
type AssociationProtocol[F, T, C any] interface {
    Check(
        context.Context,
        AssociationRequest[F, T, C],
    ) error
}

type AssociationRequest[F, T, C any] struct {
    Source      SourceSession[F]
    Input       Ref[F]
    Mapping     BindingID
    Scope       Scope[T]
    Environment EnvironmentID
    Value       Mapped[T, C]
}

// Identifies the invocation's configuration, profile and formal environment.
type EnvironmentID struct {
    /* opaque host-issued identity */
}
```

The host checks its own invocation, environment and scope identities. The integration checks that its carrier describes the requested source expression and selected representation association.

`C` should therefore be an integration-issued immutable token, or contain such a token. Its constructors are typed integration APIs with restricted construction rules. `Check` is not a universal proof constructor, and it cannot manufacture missing source membership.

The configured integration and mapping implementations remain trusted Go code. This mechanism catches mismatched associations; it does not turn the kernel into a semantic verifier.

### One signature with paired operations

```go
type ConventionID string

type Signature[F, T, C any] struct {
    /* host-interned immutable declaration convention */
}

type GenericPlugin[F, T, C any] interface {
    DescribeGeneric(
        context.Context,
        SignatureRequest[F],
    ) (SignatureSpec[F, T, C], error)
}

type SignatureRequest[F any] struct {
    Declaration SourceDecl[F] // Canonical open declaration.
    Source      SourceSession[F]
}

type SignatureSpec[F, T, C any] struct {
    Contract   ConventionID
    Parameters []ParameterHeader[T]
    Convention Convention[F, T, C]
}

type Convention[F, T, C any] interface {
    Introduce(
        context.Context,
        FormalRequest[F, T, C],
    ) (Introduced[F, T, C], error)

    Apply(
        context.Context,
        ApplyRequest[F, T, C],
    ) (Mapped[T, C], error)
}
```

`DescribeGeneric` determines the convention from the open declaration, admitted facts and fixed configuration. Its request deliberately has no caller’s target formal environment.

The first application may trigger this operation, but it must not influence the signature.

The paired requests and results are:

```go
type FormalRequest[F, T, C any] struct {
    Signature   Signature[F, T, C]
    Declaration SourceDecl[F]
    Output      Owned[T]
    Context     Context[F, T, C]
}

type FormalBinding[F, T, C any] struct {
    Binder Ref[F] // A real source binder.
    Value  Mapped[T, C]
}

type FormalSyntax[T any] struct {
    Parameters []Fragment[T]
    Clauses    []Fragment[T]
}

type Introduced[F, T, C any] struct {
    Bindings []FormalBinding[F, T, C]
    Syntax   FormalSyntax[T]
}

type ApplyRequest[F, T, C any] struct {
    Signature Signature[F, T, C]
    Use       AdmittedUse[F]
    Context   Context[F, T, C]
}
```

`FormalSyntax` contains native syntax: parameter declarations and any associated declaration-level constraint clauses. It is copied and frozen by the host. Its interpretation remains target-owned.

`Apply` produces the **complete mapped result**. It constructs the entire native use, including the actual argument transformation and result support. Returning only a loose slice of additional arguments would leave the most consequential coordination problem unresolved.

### Changes to the existing plugin surface

`Use` now returns the richer result:

```go
type Plugin[F, T, C any] interface {
    ID() PluginID
    Use(
        context.Context,
        UseRequest[F, T, C],
    ) (Mapped[T, C], error)
}
```

`UseRequest` contains the admitted use, the applicable host-issued signature and the caller’s context. For a generic declaration, `Use` delegates to that signature’s `Apply` operation.

The context becomes:

```go
type Context[F, T, C any] interface {
    Source() SourceSession[F]
    Syntax() Builder[T]

    Represent(
        context.Context,
        Ref[F],
    ) (Product[F, T, C], error)

    Lower(
        context.Context,
        Ref[F],
    ) (Type[T], error)

    Own(
        context.Context,
        DefinitionSpec[F, T, C],
    ) (Owned[T], error)

    External(
        context.Context,
        ExternalRef[T],
    ) (Symbol[T], error)
}

type DefinitionSpec[F, T, C any] struct {
    Subject   DefinitionSubject[F]
    Header    DeclarationHeader[T]
    Signature Signature[F, T, C]
}
```

`Lower` is a convenience over `Represent(...).Payload()`. It does not provide a separate path around source admission or association checking.

For a source definition, `Signature` is required, and its layout is authoritative. `Header.Parameters` must agree with it. `BinderLink` disappears.

For an ordinary native helper, `Signature` is absent, its native parameters are explicit, and it receives no inherited source binder environment.

`Publish` and `Define` retain their existing purposes. `PublishRequest` gains the shared signature. `DefineRequest` receives the frozen formal syntax and a context populated by `Introduce`.

Finally, the target’s definition operation needs to preserve that syntax:

```go
Define(Owned[T], FormalSyntax[T], Decl[T]) error
```

The target either assembles the fixed formals into the declaration or checks that the completed declaration contains them with the reserved anchors. This is a structural target check. It must not silently accept a declaration that omitted the convention’s bound.

This extension is justified directly by `Ignore<T: Tree>`: reserving an extra parameter does not ensure that the final declaration constrains it.

### Why this is still a small protocol

The new semantic seam consists of two things:

1. A checked representation product.
2. A declaration convention with `Introduce` and `Apply`.

The remaining types describe their lifecycle and inputs. There is no generic kind AST, evidence registry, alternate dispatch search or serialized elaboration program.

Public factories and issuer interfaces should follow Appendix B’s existing approach. Source adapters issue admitted-use handles; the host issues products and signatures; integrations provide typed support constructors. Independent implementations need no access to private Bittype packages.

## 3. Lifecycle, ownership and scope

I would specify the lifecycle precisely.

1. The host admits the whole effective source use.
2. It selects the mapping by exact identity before expanding the declaration body.
3. If a generic convention is needed, it obtains the signature for the canonical open declaration.
4. `Own` reserves symbols using the fixed native layout and queues an immutable definition request.
5. When the queued definition runs, `Introduce` constructs fresh formal products using those reserved symbols.
6. The host checks exact real-binder coverage and installs the resulting formal environment.
7. `Define` runs in that environment. The target preserves the frozen native formals and constraints.
8. Every application goes through the same convention’s `Apply`.

This remains lazy declaration elaboration within the existing transformer. It does not require a public planning stage or a second semantic provider.

| Object | Lifetime and ownership |
|---|---|
| `AdmittedUse` | Exact source invocation, effective reference and assumption context. |
| `Signature` | Invocation, configured mapping and canonical open source declaration. |
| `Product` | Invocation, effective source reference, representation association and target formal environment. |
| Formal syntax | One reserved declaration scope. |
| Owned helper | Existing mapping-local resource ownership, with explicit native parameters. |
| Export contract | Durable, versioned metadata independent of invocation handles. |

A signature driver may retain immutable open-declaration facts and configuration. It may not retain the first caller’s products, native symbols or formal frame.

Products containing native formal symbols are scoped. A cache keyed only by `Ref` is insufficient: the same source reference can be elaborated in different target formal environments. Any cache must include that environment identity.

The support carrier should not enumerate every kind membership an expression might have. A convention requests the support needed for its admitted obligations through the configured integration. There is no kernel operation that searches all kinds or discovers a “best” evidence provider.

## 4. Walking the examples

### `Page<T>`

A simple convention introduces one native parameter, `N`, and binds the real source entity `Page.T` to a product whose payload is `N`.

Defining:

```text
items: list<T>
```

then proceeds compositionally:

- The structural mapping projects the actual source child.
- Lowering the effective reference to `Page.T` finds its formal product.
- Its payload is `N`.
- The selected List mapping builds its native family around `N`.

For `Page<Cell<Person>>`, Page’s application convention requests the product for `Cell<Person>`. The exactly selected Cell mapping supplies that product. Page consumes the payload required by its convention and constructs its complete application.

Page has no Cell-specific rule.

A richer profile may also require a result association for `Page<Cell<Person>>`. That association is constructed for the Page application. It is not inferred from the native record shape.

### `TreePage<T>`

Two possible conventions illustrate why the core should not prescribe a tuple:

| Convention | Native formals | Payload for real binder `TreePage.T` | Available support |
|---|---|---|---|
| Witness projection | `S` | `Native<S>` or `S::Repr` | Association and the required membership under the formal constraint. |
| Expanded arguments | `S, N, P` | `N` | Correlated association and membership tied to these formals. |

These are protocol examples, conditional on 0009 establishing a valid encoding for the target and domain.

`Introduce` builds the formal product and native constraint syntax together. Extra native parameters do not create additional Horos binders.

For `TreePage<Person>`, `Apply`:

1. Obtains the relevant admitted source membership.
2. Obtains Person’s selected representation product when the encoding needs it.
3. Checks compatibility with the callee’s association protocol.
4. Supplies the appropriate native actuals.
5. Produces a result product for **`TreePage<Person>`**.

That last step matters. Person’s membership support is not automatically support for the resulting page. The source’s conditional result guarantee and the selected representation family must justify the result association.

A useful additional acceptance example is:

```text
TreePage<TreePage<T>>
```

It exposes an implementation that transports argument support correctly but forgets to construct result support.

### `Outer<T>`

`Introduce(Outer)` binds the real entity `Outer.T` to Outer’s formal product.

The admitted use of `TreePage<T>` supplies the substitution:

```text
TreePage.T → reference to Outer.T
```

The callee’s `Apply` requests the representation of the **actual** reference. The caller’s context returns Outer’s product, with Outer’s native symbols and its scoped assumption.

It must not retrieve the product created while defining TreePage, match the two parameters by spelling, or turn Outer’s assumption into an unconditional fact.

The result support for `TreePage<Outer.T>` remains conditional on Outer’s assumptions. If the admitted model does not cover the required transported judgment, generation fails; it does not infer the theorem from the body.

### `Ignore<T: Tree> = text`

This case should be part of the first review.

Source admission must happen before payload-driven child traversal. Otherwise, a mapping returning constant text could bypass the source bound simply by never calling `Lower(T)`.

Under the representation profile, native constraints may be erased to the extent explicitly promised, while preserving the declared open generic interface.

Under an enforcing profile, the native generic must retain a valid admission gate even though its payload ignores the argument. If the target cannot express that interface without a representation change forbidden by the profile, publication is unsupported.

This does **not** require eager payload generation for every source argument. An encoding that needs only source identity and support may work directly from the admitted argument. When it consumes a native representation, however, it must consume the checked association product.

### Helpers and recursive dependencies

A helper that needs a caller’s `S`, `N` or `P` receives them as explicit native arguments. Its deferred definition uses its own native parameters and an immutable recipe.

Retaining an immutable reference to `Outer.T` would still be a capture if the helper expected that reference to resolve later in an inherited formal environment. Immutability alone does not fix scope escape.

Reservation solves symbol identity and emission order. It establishes neither legal native recursion nor recursive source membership. Those remain explicit supported-domain questions.

## 5. Identity, overrides and imported artifacts

### Separate source identity from representation association

The protocol needs distinct identities for distinct purposes:

| Identity | Purpose |
|---|---|
| Exact selector | Selects the mapping. |
| Source expression identity | Identifies what a membership claim concerns. |
| Configured `BindingID` | Identifies the selected implementation/configuration and resource owner. |
| Representation association identity | Identifies the native family associated with that source expression. |
| Convention/support contract identity | Identifies the formal/use transformation, encoding domain and support authority. |

I would initially use a conservative source-owned identity policy: preserve pinned declaration identities, ordered effective applications and real binder identities. Do not introduce alias unfolding or normalization merely to share witness names. Such identities need not assert that differently identified expressions are semantically unequal.

The source adapter owns any later equivalence or normalization policy.

There is also a necessary decision about `Native<S>`:

**A single-valued projection cannot use the same `S` to choose two different representations simultaneously in one native environment.**

If the same source expression can have different associations, either:

- The native witness identifies the association as well as the source expression.
- The projection has a separate representation-family index.
- The profile prohibits coexistence of incompatible associations.

I would initially prefer an association-qualified native witness, while retaining source identity separately in metadata. Calling it only a “source witness” would hide part of its meaning.

### Overrides preserve source facts, not native compatibility

Changing Text from `string` to `app.Text` does not change source membership. It can invalidate a native constraint, projection or evidence association that depended on the previous representation.

An explicit override authorizes replacement. It does not establish compatibility.

Similarly, changing List from a slice to `Collection<T>` is coherent only if it defines the family used by the open declaration and all supported applications. A later concrete argument cannot silently change the storage selected by an already published generic.

Native specialization can be part of a declared family where the profile permits it. It cannot be inferred from the first application.

### An imported export needs a contract

`ExternalRef{Artifact, Export}` identifies native syntax to reference. It does not describe how source generic applications map to that syntax.

An imported generic needs owner-supplied metadata covering:

| Contract information | Why it travels |
|---|---|
| Exact source declaration and signature identity | Connects the export to the right declaration, bounds and result contract. |
| Source identity policy | Defines what its witnesses or association identities mean. |
| Target endpoint and required native feature domain | Establishes compatibility with the target integration. |
| Representation promise and supported domain | Separates erasure, membership soundness, completeness, formation rejection and alias/wrapper behavior. |
| Generic convention identity/version | Identifies native formals, actual transformation and payload projections. |
| Result association behavior | Explains how the resulting application can participate in further composition. |
| Fixed representation dependencies | Records choices baked into the exported definition. |
| Parametric requirements | Records what actual arguments may supply or vary. |
| Support exports and construction authority | Identifies compatible evidence ownership across packages. |
| Exact artifact pin and native export identity | Ties the claim to the supplied artifact. |

The descriptor can identify a compiled Go convention implementation supplied by the integration. There is no need for a generic interpreter for serialized recipes.

These descriptors are trusted producer claims. Their hashes identify the claims and artifacts; they do not prove them.

### The fixed-dependency problem is broader than argument conversion

For:

```text
Page<T> = { items: list<T>; next?: text; }
```

an imported Page can depend on all of these:

- The selected List family.
- The Text representation used by `next`.
- The optional-field representation.
- Any associated support definitions.

The native argument supplied for `T` cannot repair a mismatch in the representation of `next`.

An external mapping may explicitly select a frozen representation family, if the profile permits it. It must not silently claim to implement a locally selected structural family whose fixed dependencies changed.

Check the relevant fixed dependencies and parametric requirements. Requiring equality of the entire effective binding table would unnecessarily reject unrelated overrides.

For the first implementation, exact equality of required contract identities is a sensible conservative rule. Any compatibility between different contracts should be an explicit integration-owned claim.

### Preserve the `ASTSet` boundary honestly

The existing transformer output does not carry a dependency lock or export manifest. The brief correctly calls this out. [Research 0008, Appendix B](https://chatgpt.com/library?file_id=libfile_b489cbb406c4819188102ba898367349&account_id=personal)

Initially, export descriptors can be supplied alongside integrations in the explicitly locked artifact catalogue.

If automatic descriptor emission later becomes a requirement, it needs an explicit generation result or artifact output contract. Do not hide it in mutable callbacks, fake AST nodes or source semantic facts.

### Required failure behavior

The entire invocation should fail when:

- Required source admission or conditional membership coverage is absent.
- An argument’s association cannot satisfy the callee’s convention.
- Source identity matches but representation family or support authority does not.
- An imported artifact has incompatible fixed dependencies.
- A declaration cannot preserve the required native admission gate.
- A wrapper violates the selected alias/interchangeability promise.
- A signature depends on the first application or repeated reservations disagree.
- Support escapes its scope or requires unsupported recursion.
- A requested root cannot be published within its stated profile.

None of these failures should initiate a weaker mapping search.

## 6. The smallest reviewable implementation sequence

I would separate approving this composition seam from approving any enforcing encoding.

### Step 1: Fix the contract meanings

Record:

- Source expression identity policy.
- Representation association identity.
- Formal/application and result laws.
- Source admission coverage.
- Supported-domain and compatibility rules.

This step should establish what the identifiers mean before adding identifiers to interfaces.

### Step 2: Implement the degenerate representation convention

Use the existing proper-parameter, direct, saturated generic domain. Represent the current one-to-one mapping through the richer product and convention.

The decisive review cases are:

| Case | What it distinguishes |
|---|---|
| `Page<Person>` | Basic declaration/application agreement. |
| `Page<Cell<Person>>` | Composition without contract-name dispatch. |
| Text and List overrides | Consistent family selection through an open body. |
| `Ignore<T: Tree>` | Admission independent of payload traversal, with explicit publication limits. |
| The same source reference in two formal scopes | Correct product scope and cache identity. |

This gives the seam practical value before making native enforcement claims.

### Step 3: Add one encoding justified by 0009

Choose one target and one precise domain. Require it to demonstrate:

- A structured payload projection or expanded argument convention.
- Formal assumption forwarding through `Outer`.
- Result support through nested constrained applications.
- Correlation between source identity, payload and membership.
- A compatible construction authority across the supported package boundary.

Where these prerequisites are absent, enforcement remains `Unsupported`.

### Step 4: Exercise one imported generic contract

Include a matching artifact and explicit incompatibilities caused by changed Text, List or support ownership. Review the distinction between fixed dependencies and parametric requirements.

These are proposed future acceptance cases; this consultation requires no implementation or build workflow.

### Keep the initial domain narrow

I would defer constructor parameters, computed callees, arbitrary mixed type/const/lifetime parameters, fragment-local binders, automatic convention adaptation, unrestricted recursion and cross-invocation caching.

A closed generated enforcement domain can be an honest first result. It can still expose open generics over its supported members. It must not be described as complete admission for arbitrary recursive Tree descriptions or arbitrary native types.

## 7. What I would challenge in the brief

### The split with 0009 cannot be completely clean

0008 can define the composition obligations independently. It cannot freeze every public support interface before 0009 determines:

1. Whether membership support is constructible and consumable across independently compiled packages.
2. Whether conditional facts survive symbolic substitution, including result associations.
3. Which identities and native parameter shapes preserve correlation.
4. Which native formation and alias guarantees the encoding actually supports.

These are composition-critical prerequisites. A carrier that transports claims consistently cannot make the claims sound.

### A source fact and its native encoding have different authorities

The source provider can establish that Person satisfies Tree. It does not thereby authorize an arbitrary package to construct native support for every proposed representation of Person.

The integration must establish that the native support is tied to the selected association and accepted authority. An open public `Proof<S,N,K>` constructor would remove the very restriction the enforcing profile is meant to preserve.

### More native parameters do not necessarily improve enforcement

A three-parameter encoding can expose more opportunities to misassociate components. A single witness can lose necessary correlation or introduce unwanted native equivalences.

The review criterion should be preservation of the declared association and admission promises, rather than parameter count.

### Prior art supports the paired operations, with limited reach

MLIR’s signature conversion explicitly relates original inputs to new inputs or replacements, and its conversion framework requires uses to be reconciled when definitions change. That supports the need for a shared formal/use correspondence. It does not transfer MLIR’s fallback selection or materialization mechanisms into this design. [MLIR](https://mlir.llvm.org/docs/DialectConversion/?utm_source=chatgpt.com)

GHC’s evidence API separately represents type and evidence abstractions and applications. This supports explicit treatment of both sides of elaboration, without requiring Haskell dictionaries or instance resolution. [downloads.haskell.org](https://downloads.haskell.org/~ghc/10.0.0.20260917/docs/libraries/ghc-10.0.0.20260917/GHC-Tc-Types-Evidence.html?utm_source=chatgpt.com)

The System FC paper supports keeping established facts distinct from their explicit elaborated representation, and gives a concrete reason to budget generated support: evidence terms can grow substantially. Its equality witnesses do not establish Horos membership soundness. [microsoft.com](https://www.microsoft.com/en-us/research/wp-content/uploads/2016/07/fc-new-tyco.pdf?utm_source=chatgpt.com)

## Decision I would record

**Approve a shared declaration convention with `Introduce` and `Apply`, and a checked representation product carried through real-binder environments.** Make native constraints and result associations part of that convention. Keep source admission independent of payload demand, and require explicit compatibility contracts for imported artifacts.

Implement the representation-only case first. Admit an enforcing profile when a target integration supplies the domain, authority and substitution guarantees required by this protocol.

Do not standardize an `S/N/P` tuple, treat a witness as an automatic proof, or use configuration identity as a substitute for compatibility.