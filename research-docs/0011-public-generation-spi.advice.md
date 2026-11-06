# Advice for consult
**Research doc:** doc_9d0d3c9e-4487-4300-8e9d-3f72381774aa.md
**Document:** doc_9d0d3c9e-4487-4300-8e9d-3f72381774aa
**Advisor:** ChatGPT, model "6 Pro" as labelled on the answer page
**Session:** nsess_7165b985-e5bc-4856-b6ef-d843e10dd77d
**Run:** run_971a5b3d-7108-49fe-936a-f2ce60d7c222
**Chat:** https://chatgpt.com/c/6aed6f2c-3294-83eb-a5a7-ef2ef2784dfe
**Verify:** nightfall consult verify run_971a5b3d-7108-49fe-936a-f2ce60d7c222 --output <this file>

I would choose **an integration-owned association token inside a host-owned, request-bound product handle**. The integration decides what connects a particular source use, under particular assumptions, to a particular native fragment. The host decides whether that association may be used here. Give the host a way to ask the configured integration to authenticate an association; give it no constructor that takes a source, a fragment, and independent support.

For generics, retain **one callee binding per canonical open declaration and selected convention in an invocation**. It owns the authoritative native layout. `Introduce` builds it from the open declaration; `Define` consumes its formal environment; `Apply` consumes the caller’s actuals. `Use` and `Publish` can both acquire it. Neither can reconstruct it from the first specialization encountered.

There is also one concrete prerequisite in [research document 0011](https://chatgpt.com/library?file_id=libfile_fde3_f-7042f42b03608191abc44ccc9d376da1_13e49f0e3cac8191bc5a253deec2d8a4&account_id=personal): **preserve the original `semantic.Snapshot` through the typed transformer adapter.** The supplied `transform.Request[L]` exposes unpacked source, entry, inputs, and evidence, but not the snapshot itself. Store `input.Source` in the request and add `Snapshot() semantic.Snapshot`, or start the generation driver at the boundary that still has `transform.Input`. Do not recapture `Inputs()`. Do not recover the snapshot from `Evidence(0)`: an empty evidence contract is legitimate, and a receipt should not be the only route to the identity of the input.

The design below is a proposed public contract and its implementation acceptance criteria. The repository names and existing behavior I refer to are from the supplied baseline; the proposed API names are new.

## 1. What a product means, and who can issue one

Separate three objects:

| Object | Issuer | Meaning |
|---|---|---|
| **Admitted use** | Configured source integration, accepted by the driver | This complete source occurrence denotes this effective subject under these original assumptions; the source obligations required by this profile were met using the supplied evidence. |
| **Association token `P`** | Selected mapping or convention’s association authority | This exact native fragment represents this effective subject under those assumptions, according to this integration’s rule and support. It may be conditional on open formals. |
| **Accepted `Product`** | Generation driver | The configured authority authenticated the token, the target admitted its fragment in the required frame, and the driver accepted it for this operation. |

**Do not give `P` a self-certifying `ValidFor(request)` method that the host treats as authoritative.** The driver chooses the authority from the frozen configuration. Otherwise the candidate effectively chooses its own validator. A target helper can construct a perfectly good `NativePage<N>` fragment; that gives it no authority to say whether the fragment represents `U`, `Page<U>`, or anything else.

Keep **operation lifetime** separate from **lexical lifetime**. An `Introduce` callback loses the power to issue when it returns. Its accepted formal products must survive for `Define`. A child `Use` callback ends before its parent `Apply` consumes the child product. Conversely, a `Publication` from a previous callback must not satisfy a new callback, even in the same lexical frame. A single “scope ID” cannot express all three cases.

I would use these checks:

- **On issuing or returning an operation result:** exact live invocation, configuration, operation ticket, selected capability, and expected result role.
- **On consuming an accepted product:** live invocation, compatible configured authority, exact effective subject and source assumptions, and a permitted native frame. Check the producer’s layout where the rule requires it. Do not require a child product to have its parent’s operation ticket or its parent’s callee layout.
- **On committing output:** every requested root has a valid exposure, every required owned reservation has a definition, the target can resolve the symbolic obligations, and the invocation has not failed or been cancelled.

For the first slice, use **exact lexical-frame equality** for products. It is conservative and easy to test. A closed primitive can be issued in each frame that needs it. General weakening, lifting, or transporting an association across frames is a separate integration operation if a later example requires it.

## 2. A concrete public construction and access boundary

Here is one way to close the API without asking the host to clone an arbitrary `P any`. Use **owner-qualified, typed references into immutable provider records**. This is an implementation choice for the first slice, not a requirement to build a common semantic model.

The entire common identity mechanism can be this small:

```go
type domain[K any] struct{ identity byte }

type Space[K any] struct{ owner *domain[K] }

type Ref[K any] struct {
	owner *domain[K]
	slot  uint64
}

func NewSpace[K any]() Space[K] {
	return Space[K]{owner: &domain[K]{identity: 1}}
}

func (s Space[K]) Ref(slot uint64) (Ref[K], bool) {
	if s.owner == nil || slot == 0 {
		return Ref[K]{}, false
	}
	return Ref[K]{owner: s.owner, slot: slot}, true
}

func (s Space[K]) Index(r Ref[K]) (uint64, bool) {
	if s.owner == nil || r.owner != s.owner || r.slot == 0 {
		return 0, false
	}
	return r.slot, true
}
```

A `Space` stays with its issuer. A `Ref` is an address, **not an admission**: the owner must still find a live, published record of the correct role at that address. The owner inserts each frozen record once and never replaces it or recycles its slot during the invocation. Another package can create a space of the same Go type, but it cannot make its references belong to the configured owner. The identity allocation is deliberately nonempty: Go permits pointers to distinct zero-size variables to compare equal. Go also permits copied slices, pointers, functions, maps, and interface values to retain shared mutable state. Neither a struct copy nor a generic type constraint solves that problem. [The Go Programming Language](https://go.dev/ref/spec?utm_source=chatgpt.com)

For readability, the signatures below are shown **after instantiation for one source/target pair**. These abbreviations have concrete types:

| Name | Concrete carrier | What the owner stores |
|---|---|---|
| `D`, `B` | Distinct `Ref` types for source declarations and binders, parameterized by the source marker | The exact model owner and real declaration or binder; never a spelling or native parameter number. |
| `U`, `E` | Distinct `Ref` types for a source use and source environment | Immutable source descriptors, effective substitutions, retained assumptions, and references to supplied query results. |
| `Z`, `H`, `N` | Distinct `Ref` types for a target destination, layout, and fragment, parameterized by the target marker | Local symbol, external export, or inline template; frozen formal arrangement; reusable native syntax and symbolic references. |
| `P` | `Ref` for an association, parameterized by both endpoint markers | The issuing rule’s private immutable association record. Different rules can have different private record types. |
| `Site`, `Want`, `Callee`, `Product`, `Introduction`, `Publication`, `Definition` | Concrete structs with private host state | Invocation identities, expected binding, saved callee, or accepted operation result. Only the driver’s lifecycle methods return nonzero instances. |

For example, `B` can be `Ref[BinderRole[Horos]]`, `N` can be `Ref[FragmentRole[Go]]`, and `P` can be `Ref[AssociationRole[Horos,Go]]`. These role markers distinguish carriers. They do not enumerate kinds, representations, laws, or proof forms. A package with an existing suitable opaque token can use that instead; there is no benefit in translating a good immutable source model into a second graph.

### Source admission

The source protocol needs an open-declaration operation and a whole-use operation:

```go
type Formal struct {
	Binder  B
	Subject U
}

type Actual struct {
	Binder B
	Use    U // effective actual interpreted in the caller's environment
}

type OpenData struct {
	Decl    D
	Env     E // canonical open assumptions, independent of the caller
	Body    U
	Formals []Formal
}

type UseData struct {
	Use     U // complete admitted use; retains its evidence and occurrence
	Subject U // effective subject used for association matching
	Env     E
	Callee  D // zero for a use with no declaration callee
	Bound   B // nonzero only for a source-certified binder lookup
	Actuals []Actual
}

type Source interface {
	CheckInput(context.Context, semantic.Snapshot, []semantic.Result) error
	RootEnv() E
	Open(context.Context, D) (OpenData, error)
	Admit(context.Context, E, U) (UseData, error)
}
```

`OpenData` and `UseData` are **construction records returned by the configured source**, not freely admissible claims. The driver calls `semantic.AdmitEvidence` at entry, checks the source integration against that exact snapshot and those admitted receipts, and copies the returned record slices into private storage. It verifies membership of addresses in the configured owner, duplicate binders, and agreement with the saved open binder set. The source integration is responsible for the meaning of the substitution, arity, saturation, applicability, and required judgments. The driver does not infer any of those from the shape of `semantics.Term`.

The independent source package needs an ordinary public factory. For the fixture, give it the following concrete operations:

| Public operation | Responsibility |
|---|---|
| `horossource.Capture(snapshot, admittedReceipts, suppliedQueries)` | Check the owners, model, configuration and subjects of the supplied immutable query results, and capture any mutable construction storage. Return the `Source` implementation and its read-only catalogue. |
| `catalog.Declaration(entity)` / `catalog.Use(subject, occurrence)` | Return a `D` or `U` backed by a record in that catalogue. These are locators; possession does not establish whole-use admission. Permit model-only subjects where there is no syntax occurrence. |
| `source.Open(decl)` | Obtain the canonical declaration, its real binders and its own open assumptions. There is **no caller actual argument** to this operation. |
| `source.Admit(env, use)` | Authenticate the use in that environment, including ignored arguments, and return effective actuals and the complete admitted subject. Missing required evidence is an error. |

The distinction between `Use` and `Subject` matters. Two occurrences can both denote `Outer.U`; requiring identical `ast.NodeRef`s would reject valid reuse. Only the source authority can establish that both occurrences denote the same effective subject under the retained assumptions. The host never equates them because of identical native text or local numeric IDs.

### Association acceptance

Use **one request-completion style**. In this design providers return a host-issued `Product` obtained from the request. They do not return an independently assembled payload/support pair.

```go
// Binding contains only immutable handles and identities, copied by value.
// The operation ticket is on Want and on the accepted result, not part
// of the lexical identity of reusable support.
type Binding struct {
	Run       RunID
	Config    ConfigID
	Issuer    IssuerID
	Subject   U
	Env       E
	Frame     FrameID
	Signature SignatureID // nonzero when the producing rule uses a layout
}

type Captured struct {
	Binding Binding
	Token   P
	Payload N
}

type Associations interface {
	Capture(context.Context, Want, P) (Captured, error)
}

type Native interface {
	CheckDestination(context.Context, Site, Z, bool) error
	CheckLayout(context.Context, Site, Z, H) error
	CheckProduct(context.Context, Want, N) error
	CheckBody(context.Context, Callee, N) error
	Assemble(context.Context, []Publication, []Definition) (
		representation.ASTSet, error)
}
```

`RunID`, `ConfigID`, `IssuerID`, `FrameID`, `SignatureID`, and the operation ticket are distinct opaque host identities. `Want` exposes read-only access to the site and operation ticket, original snapshot, admitted use, effective subject, environment, selected issuer, caller frame, and saved callee/signature when applicable. `Callee` exposes its destination, layout, declaration identity, and ownership. It does **not** give a caller the callee’s formal-product environment as a substitute for actual arguments.

`r.Return(p)` has a specific implementation:

1. Check the request ticket, cancellation, budget, and invocation state.
2. Call the **configured** `Associations.Capture(ctx, r.Want(), p)`. The authority looks up `p` in its own store, checks its rule and premises, and returns metadata and the payload **from that record**. Copying the expected binding from `Want` into the reply without authenticating the record is a broken implementation.
3. Compare the returned binding with the expected binding for this result. Ask the configured target to admit the **returned fragment itself** in the expected frame. For an application, check its association with the exact saved layout and its complete native application or inline expansion.
4. Store that captured record in a new host `Product` with the current operation ticket. Recheck cancellation/liveness before accepting the callback’s return. Only the product issued for that callback can be returned as its result.

There is no public `AcceptCaptured(Captured)` entry point. `Captured` is temporary storage on a call **made by the driver**. It is harmless that an ordinary package can write a `Captured` literal: there is nothing it can do with that literal that mints a product.

Give accepted products read-only `Binding()`, `Payload()`, and `Token()` accessors. Here `Token()` returns only the immutable address `P`. The owning integration can expose `ReadPremises(product)` or a narrower domain-specific query. It checks ownership and returns frozen handles or fresh copies. No accessor should expose a mutable proof tree, slice backing array, builder, or closure that consults whatever frame happens to be current later.

### The five provider capabilities and their constructors

The following method sets are the public callback surface. Request implementations and result state are private. A `Site` identifies the current invocation, operation, selected plan, and lexical frame; it is not a public constructor for any result.

```go
type User interface {
	Use(context.Context, UseRequest) (Product, error)
}
type Publisher interface {
	Publish(context.Context, PublishRequest) (Publication, error)
}
type Definer interface {
	Define(context.Context, DefineRequest) (Definition, error)
}
type Convention interface {
	Introduce(context.Context, IntroduceRequest) (Introduction, error)
	Apply(context.Context, ApplyRequest) (Product, error)
}

type UseRequest interface {
	Want() Want
	Child(U) (Product, error)       // admit and use in this caller frame
	Bound() (Product, error)        // only for a source-certified binder use
	Apply() (Product, error)        // acquire this callee; invoke saved convention
	RequireOwned(D) (Callee, error) // separate preselected owned dependency
	Return(P) (Product, error)
}

type IntroduceRequest interface {
	Site() Site
	Open() OpenData                 // owned copy of slices; immutable handles
	Destination() Z
	FreezeLayout(H) (FormalRequest, error)
}

type FormalBinding struct {
	Binder B
	Token  P
}

type FormalRequest interface {
	WantFor(B) (Want, error)        // exact open binder, frame and frozen layout
	Finish([]FormalBinding) (Introduction, error)
}

type ApplyRequest interface {
	Want() Want
	Callee() Callee
	Actual(B) (U, error)            // all effective actuals are available
	Represent(B) (Product, error)   // demand this actual in the caller frame
	Return(P) (Product, error)
}

type PublishRequest interface {
	Site() Site
	Root() D
	Acquire() (Callee, error)
	Expose(Callee) (Publication, error)
}

type DefineRequest interface {
	Site() Site
	Callee() Callee
	Open() OpenData
	Bound(B) (Product, error)       // this declaration's saved formal environment
	Use(U) (Product, error)         // in this definition's own frame
	CompleteBody(N) (Definition, error)
}
```

Three details make this more than a list of callbacks:

- **`FreezeLayout` precedes `WantFor`.** The convention builds the whole target layout, including extra native formals and constraints, before associating source binders with products derived from it. `Finish` admits each token against the expectation for its real binder, checks exact-once coverage, and freezes the environment. There is no assumption that binder *i* represents native parameter *i*.
- **`UseRequest.Apply()` and `Bound()` are narrow forwarding routes.** `Apply()` can return the delegated result for the same whole use; `Bound()` can forward the saved product for the source-certified effective binder. The driver checks subject, assumptions, frame, and provenance and issues the current operation’s return envelope. It preserves the underlying association and its issuer. There is no `ForwardAs(otherSource, product)`.
- **`CompleteBody` accepts a body only.** The target checks it against the callee and later assembles it with the saved header. The definer cannot return a different declaration with fewer parameters or weaker constraints. `Expose` similarly checks that the callee actually exposes the requested root and that the configured target supports that exposure. An inline representation alone does not publish a root.

### How an independent integration actually mints the inputs to those methods

The concrete Go target should expose `BeginLayout(site, destination)`, builder operations for its own formals and constraints, and `Freeze() (H, error)`. Its expression constructors return `N`, including a formal reference, a reference to a reserved symbol, a complete application of a saved layout, and the few structural forms the fixture uses. These constructors capture target-owned storage. A read operation that unpacks a fragment returns a fresh copy. The agnostic SPI does not gain a `Text`, `List`, `Record`, or `Cell` method.

The concrete association owner should expose **rule constructors**, for example:

```go
Formal(want, nativeFormal, admittedFormalPremise) (P, error)
Text(want) (P, error)
PageResult(want, constructorGuarantee, actualPremise, appliedFragment) (P, error)
IgnoreResult(want, admittedWholeUse, textFragment) (P, error)
```

These are methods in their owning integration packages, not methods required of every provider. `PageResult` checks that the constructor guarantee is for this family and kind, the actual premise belongs to the admitted actual, the native application was built for this layout using the corresponding native arguments, and the conclusion is about the **whole application**. It freezes those records and mints a slot in its own association space. `Formal` introduces a conditional association under the declaration’s formal assumptions. It does not claim a native type parameter independently proves Horos membership.

At the driver boundary, the fixture supplies an immutable plan: optional `User`, `Publisher`, `Definer`, paired `Convention`, association authority, and target acquisition function. The acquisition function returns a checked destination and whether it creates a local definition obligation. A nongeneric owned plan supplies a target nullary-header constructor instead of a convention. Use an explicit table of fixture-selected routes to these plans; assign identities to the bindings when capturing the configuration.

The public driving sequence can then be as small as `Begin(ctx, input, config, limits)`, `Use(rootUse)`, `Acquire(declaration)`, `Publish(root)`, `Define(callee)`, and `Finish()`. `Begin` admits receipts and captures the configuration. `Acquire` is idempotent for the same binding within the invocation. `Define` is permitted once for an owned reservation. `Finish` is the only operation that returns usable output. The fixture drives these operations explicitly; the host keeps the obligation ledger, without having to implement a definition scheduler.

**Do not require all mappings in a program to manufacture one centrally defined `horosgo.Product` variant.** That would make adding `Cell` require editing the package that owns the variant. With owner-qualified tokens, a library can supply a rule, its own frozen support records, and a configured reader for a shared premise such as “membership and representation for this actual.” A consuming convention uses that reader on an accepted product. An unsupported producer/consumer combination fails explicitly. For the first slice, one preconfigured compatible set of readers is sufficient; a general protocol-discovery system is unnecessary.

## 3. Callee acquisition and authoritative layouts

The host’s acquisition key should include the exact source session/model, canonical declaration, selected plan/convention binding, target configuration, and invocation. Do not key it on a spelling, printed signature, or actual argument. In particular, `Page<Person>` and `Page<U>` acquire the same open `Page` binding in one configuration.

| Route | What acquisition does | Definition obligation |
|---|---|---|
| **Owned generic** | Allocate the stable declaration symbol and independent formal frame; run `Introduce` on the canonical open declaration; freeze the layout and formal products. | The owner must define the body before commit. |
| **External generic** | Authenticate the configured external artifact/export and its signature; bind the same convention to a target-owned view of that signature. `Introduce` can validate/import the existing layout. | None locally. |
| **Inline generic** | Freeze an inline template layout in its own formal frame. `Apply` returns the complete expansion and its result association. | None for the inline callee; it may require a separate owned helper. |
| **Nongeneric owned declaration** | Reserve the symbol and freeze the target’s explicit nullary header. | Define the body. No fake `Introduce` or `Apply`. |
| **Direct primitive or structural use** | Call `Use` for the admitted subject. No callee acquisition is needed unless the mapping requests a dependency. | None by default. |

For an owned callee the useful states are `Reserved → Introducing → Ready → Defining → Defined`. Publication is a separate exposure record: it may happen before definition. Allocate the symbol before introduction, but **reject application of a signature that is still being introduced**. Once the signature is ready, recursive references from a body can refer to the reserved symbol without recursively expanding the body. This gives the slice a precise cycle policy without requiring recursive signature inference.

An owned helper requested from `Outer` gets a fresh declaration frame. The caller can explicitly apply that helper to `Outer.U`; its definition cannot look up `Outer.U` through a captured caller frame. If the helper is purely native, give its preselected target plan its own native layout and definition obligation, with zero source binders. Do not invent a Horos declaration or binder for a target implementation detail.

## 4. Why `Outer`, `Nested`, and `Ignore` work

The source substitution and the representation environment answer different questions:

| Environment | Example inside `Outer<U> = Page<U>` |
|---|---|
| Source substitution for this application | `Page.T ↦ Outer.U`, interpreted under `Outer`’s original assumptions |
| Native product for the effective actual | `Page.T ↦ pOuterU`, where `pOuterU` represents `Outer.U` in `Outer`’s frame |
| Callee’s saved open formal environment | `Page.T ↦ pPageT`, used to define the open `Page` body |

**`Outer`.** `ApplyRequest.Actual(Page.T)` returns the effective source actual `Outer.U`. `Represent(Page.T)` resolves that actual in the caller’s frame and returns `pOuterU`. The selected convention builds the complete native application using the saved `Page` layout. Its result token associates `Page<Outer.U>` with that application. `pPageT` does not cross into the caller.

**`Nested`.** For `Page<Page<U>>`, the inner call consumes the premise for `U` and returns `pInner`, whose subject is `Page<U>` and whose payload is `NativePage<N>`. The outer call’s effective actual is `Page<U>`, so its premise reader must accept `pInner`. It must reject `pU`, even if someone separately supplies the beautifully formed fragment `NativePage<N>`. The constructor guarantee supplies the source membership step; the representation rule supplies the representation step. The result of the outer call is about `Page<Page<U>>`.

A Page convention’s `Apply` can therefore have this shape:

```go
func (p PageRule) Apply(ctx context.Context, r ApplyRequest) (Product, error) {
	b := p.RealPageBinder()
	u, err := r.Actual(b)
	if err != nil {
		return Product{}, err
	}
	arg, err := r.Represent(b)
	if err != nil {
		return Product{}, err
	}
	premise, err := p.ReadActualPremise(r.Want(), u, arg)
	if err != nil {
		return Product{}, err
	}
	n, err := p.Target.ApplyPage(r.Callee(), r.Want(), arg.Payload())
	if err != nil {
		return Product{}, err
	}
	result, err := p.PageResult(r.Want(), p.Guarantee(), premise, n)
	if err != nil {
		return Product{}, err
	}
	return r.Return(result)
}
```

Here `ReadActualPremise` reads the **accepted argument product** through the configured owner’s reader. The convention never receives a free-standing payload and an unrelated premise to zip together by position.

**`Ignore`.** Admit `Ignore<A>` as a whole before invoking its mapping or demanding a child product. If `A:K` is not established as required, fail even though the target expression would be `string`. If it is established and the selected Ignore rule needs no representation premise for `A`, call `Represent` zero times. Issue an association for **`Ignore<A>` represented by `string`**, retaining the source assumptions on which that conclusion depends.

There is a useful limit here: total formal binding at `Introduce` is a different requirement from mandatory representation traversal at `Apply`. Preserve total formal binding. An inline convention can have a legitimate abstract representation formal that its expansion never uses. If a proposed zero-native-formal convention cannot legitimately supply a product for a real open binder, it is unsupported under this first contract. Do not invent `Rep(T,string)` to make the binder table total. A later source-only erased-binding protocol would be a deliberate extension.

`Ignore` also gives the strongest scope-escape test: **a fragment can have no free native symbols while its support still depends on `U:K`.** Target dependency checking alone will not catch that escape. The source assumption check must fail independently.

## 5. Fragments, final verification, and the release gate

For the Go fixture, a fragment can be a frozen target AST subtree plus typed symbolic relocations owned by the Go target. No second grammar is necessary. The target must derive the actual symbol dependencies from its representation, or construct fragments so that every symbolic occurrence is registered. A caller-supplied `FreeSymbols` list that can omit a reference is insufficient.

Materializing a fragment twice allocates two sets of containment nodes. References to the same local declaration may resolve to the same final declaration node. External exports remain target-owned artifact/export descriptors until target assembly emits the appropriate import and native syntax. They do not acquire a fabricated resolved `ast.NodeRef`. Those rules follow directly from the containment/reference distinction in the attachment’s `trees/ast/go` extracts.

**Before assembly, the promise is scoped symbolic representation under the selected integration’s stated assumptions.** It is not a judgment about future target bytes, all constructible native values, native enforcement of `K`, or equivalence of source and target programs.

After successful assembly, capture the actual output, project that capture, and construct the expected verification statement from the completed participants. A generation trace may identify the selected family, premises, source origins, and materialization sites to inspect. It cannot serve as the verification result. In the supplied PR101 pattern, a previously checked theorem is retained as a definition and used in checking a new claim. Reuse that pattern only after establishing correspondence between the theorem’s subjects and the **actual newly captured output**. Successful Go compilation and successful Logos checking answer different questions.

I would make these the remote acceptance gates:

| Gate | Observable result required |
|---|---|
| **Independent-package construction** | Source adapter, target helper, association owner, mapping, and driving fixture use public APIs only. A library-owned mapping can be added through consumer configuration without editing a central type switch. |
| **Canonical introduction** | Visit `Page<Person>` first and last. Exactly one open layout is introduced for the selected binding, with no `Person` capture. Test renamed binders and two unrelated binders both spelled `T`. |
| **Non-positional layout** | At least one fixture has **more native formals than real source binders**, or a justified non-one-to-one arrangement. `Define` uses the frozen header; `Apply` supplies the complete native layout. A one-to-one Page fixture alone cannot detect a regression to binder indices. |
| **Whole-use admission and result support** | Invalid `Ignore<A>` fails before child representation or mapping execution. A valid ignored actual need not be represented. Nested Page succeeds with the inner result and fails with argument-only support beside the same Page-shaped payload. |
| **Identity and scope** | Reject equal syntax recaptured separately, model-local IDs from another receipt, a different configuration instance with identical settings, a lookalike layout, a foreign formal, and a syntactically closed product carrying another frame’s assumptions. |
| **Lifetime and immutability** | Formal products survive `Introduce` into `Define`; child products survive into parent `Apply`. Retained requests cannot issue later. Old publications/definitions cannot be replayed. Mutating construction inputs or unpacked views cannot change accepted records. |
| **Ownership and assembly** | External generic use succeeds without a local definer. An inline use can require an independently scoped owned helper. A missing root publisher or unfinished owned reservation fails the invocation. Repeated materialization creates fresh containment nodes. |
| **Failure and resource bounds** | Exercise cancellation before dispatch, after a callback, and before commit; operation/depth limits; and output expansion limits. Every failing path yields no usable AST set. Charge **materialized nodes**, not just unique fragments. |
| **Actual target artifact** | Assemble a real Go type AST, run the existing syntax/format/print path, and remotely compile/type-check the generated fixture. Record the result as ordinary Go syntax and typing for the fixture’s supported profile. Native source-kind enforcement remains `Unsupported`. |

Bound source arguments, binders, layouts, reservations, operations, recursion depth, captured support/fragment size, and total materialized output. Check limits before allocations where the host or target controls them. These counters bound work performed through the protocol; an arbitrary trusted in-process callback can still ignore cancellation or allocate outside it. Freeze configuration as carefully as layouts: retaining the same convention pointer does not prevent its implementation from changing its policy between `Introduce` and `Apply`.

**My release recommendation:** ship the public issuance paths, retained callee binding, body-only definition, immutable fragment boundary, and explicitly driven fixtures together. Keep the four production transformers on their existing path until the replacement milestone. The decisive demonstration is an independent package that can produce `Page<Page<U>>`, reject the wrong inner support, and refuse an escaping `Ignore<U>` product—while the host itself knows nothing about `Page`, `K`, `Rep`, or Go’s type rules.