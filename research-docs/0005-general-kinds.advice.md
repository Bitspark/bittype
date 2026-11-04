# Expert advice: general structural kinds

Advice for [research 0005](0005-general-kinds.md). The research document records
the accepted decisions and adaptations; this response is not the implementation
specification.

## Recommendation

Adopt **structural kinds over checked type descriptions**, with ordinary named and parameterized definitions, unions and intersections, guarded recursive definitions, and explicit constructor contracts. Retain checked constructor bodies so that the same kind checker can establish more precise results when requested.

`Data` should be an ordinary recursive definition in a pinned prelude. `Flat`, `Rows`, `Tree`, and imported user definitions should use exactly the same evaluator, inclusion procedure, and evidence representation.

I recommend three explicit guarantees for the first implementation:

1. **Exact structural membership for checked, closed proper types.**
2. **Sound structural inclusion over a specified shape universe**, with conservative treatment of shapes that no admitted Horos descriptor realizes.
3. **Sound, deliberately incomplete constructor-contract checking**, using declared guarantees for abstract constructors and checked bodies for known constructors.

This answers the owner’s generalization requirement without promising unrestricted semantic subtyping. The implementation baseline below is the [supplied research brief](0005-general-kinds.md).

## 1. Meaning, grammar, and ordinary definitions

### Classify type descriptions structurally

For a checked proper expression \(T\), construct an observation graph \(O(T)\):

- Resolve aliases and evaluate known applications.
- Preserve primitive form tags and relevant metadata.
- Follow only schema-designated type-child positions.
- Keep `exact` payloads opaque.
- Represent ground `rec`/`ref` using lexically resolved backedges.

The graph is finite; its unfolding may be infinite. Its runtime inhabitants remain finite values.

A proper kind denotes a language of these structural observations:

\[
T:K \quad\Longleftrightarrow\quad
T\text{ is admitted and }O(T)\in\llbracket K\rrbracket.
\]

This deliberately gives `union(text, record(...))` a **Union** root. It does not automatically belong to `Scalar | Record`.

Structural refinement of types has a direct precedent in Caires and Toninho. Their broader refinement logic uses theories that are generally undecidable, so importing that full logic would exceed the proposed qualification boundary. [arxiv.org](https://arxiv.org/html/1908.00441v1?utm_source=chatgpt.com)

### Necessary primitive distinction: sorts

Use finite simple sorts:

\[
\sigma ::= P \mid (\sigma_1,\ldots,\sigma_n)\to\sigma,\qquad n\ge1.
\]

`P` identifies the domain of proper type descriptions. Arrow sorts identify constructor values. Each kind has exactly one sort; unions and intersections require equal sorts.

This distinction is necessary to reject a constructor where a field type is required and to determine application arity. It carries no `Data` policy.

Provide `⊤σ` and `⊥σ` at each sort. A library may define:

```text
kind Type = ⊤P
```

The checker recognizes the primitive top operator, never the name `Type`. Also, ground `any` is one descriptor; it is not `⊤P`. Kind bottom contains no types, whereas `union()` is a type describing no values.

### Exact proposed canonical core

Reuse the brief’s `D`, `T`, `I`, `L`, and canonical-map notation. Let `H` be a 32-byte digest. Lists below are finite.

```text
S = D("sort-proper")
  | D("sort-arrow", L(S, ...), S)

KR = D("kind-local", T(name))
   | D("kind-imported", T(alias), T(name))

K = D("top", S)
  | D("bottom", S)
  | D("kind-param", I(slot))
  | KR
  | D("kind-apply", KR, L(K, ...))
  | D("or", S, L(K, ...))
  | D("and", S, L(K, ...))
  | D("node", T(form), selectorMap, Q)
  | D("nu", L(K, ...), I(entry))
  | D("recursive-var", I(depth), I(slot))
  | D("arrow", L(K, ...), K)

Q = D("epsilon")
  | D("item", K)
  | D("concat", L(Q, ...))
  | D("alternative", L(Q, ...))
  | D("star", Q)

KindParameter = D("kind-parameter", T(name), S)
TypeParameter = D("parameter", T(name), K)

Declaration =
    D("kind", L(KindParameter, ...), S, K)
  | D("type", L(TypeParameter, ...), K, E)

Module =
    D("module", I(4), T(moduleName), H(signatureDigest),
      importsMap, declarationsMap)

Import =
    D("import", T(expectedModuleName), H(moduleDigest))

E = existing primitive descriptor syntax,
      with E only in designated type-child positions
  | Tuple(A("horos/wire"), E)
  | D("param", I(slot))
  | D("local", T(name))
  | D("imported", T(alias), T(name))
  | D("apply", E, L(E, ...))
```

The formation rules complete this grammar:

- `node`, `nu`, and recursive variables have sort `P`. `item` contains a proper kind.
- A node matches its form, selected metadata, and **entire ordered sequence** of observed type children against `Q`.
- Selector maps test only schema-enumerated metadata such as scalar selector and record mode. No arbitrary predicates, callbacks, or subtree equality.
- Arrow lists and applications are nonempty. `or`, `and`, `concat`, and `alternative` have at least two operands.
- A kind reference without application names a zero-parameter definition. Kind applications are exactly saturated and first order.
- A `nu` binds a nonempty simultaneous group of proper-kind equations; `entry` selects its result. Recursive variables use lexical group depth and slot.
- Named kind-definition references form a DAG. Recursion occurs explicitly inside `nu`; its parameters remain fixed.
- Every recursive cycle crosses a node’s type-child observation. Boolean operators and aliases do not guard cycles. Reject malformed recursion before simplification.
- Type declarations may return any sort. Their result annotation is a kind of that sort.
- Declaration names need no capitalization convention. References must resolve to the appropriate declaration category.

`Q` describes finite sequences, including unbounded finite widths through `star`. It does not admit an infinite tuple or record.

These restrictions exclude growing instantiations such as `Grow<K> = List<Grow<List<K>>>`. They also avoid recursive arrow contracts and quantified constructor schemes in the first release.

### Structural families and recursive kinds

Define the familiar families through ordinary `node` expressions:

```text
List<K>       = node(list-of, {}, item(K))
Tuple<K>      = node(tuple, {}, star(item(K)))
UnionOf<K>    = node(union, {}, star(item(K)))
Map<K,V>      = node(map-of, {}, concat(item(K), item(V)))
Wire<K>       = node(wire, {}, item(K))
WireFamily    = Wire<⊤P>
```

`Scalar` should have an explicit structural definition: the three scalar embedding selectors and existing text-refinement forms. Raw atoms and their refinements can have a separate `AtomFamily`. An exact-text singleton belongs to `Exact`, without automatic Scalar membership.

Define `Record<K>` to constrain every declared field, including optional fields, **and the implicit additional-field type of an open record**. That additional type is the ground descriptor `any`; expose it as an implicit observation child. Consequently, an open record with scalar named fields does not establish `Record<Scalar>`.

Labels, presence markers, bounds, and literal payloads remain validated and preserved. A uniform family query simply makes no label-specific or bound-specific guarantee.

Then:

```text
kind Flat = Scalar | List<Scalar>
kind Rows = Record<Flat>
kind Tree = ν X. Scalar | List<X> | Record<X>
```

Use **greatest fixed points** for recursive kinds. For example, the valid descriptor `rec r. list<ref r>` has an infinite regular unfolding and can belong to `Tree`. A least fixed point over finite unfolded type trees would exclude it.

Keep the two guards separate:

- Ground recursion requires strict **runtime-value descent**.
- Kind recursion requires structural **type-description descent**.

`Not<K>` is a positive structural node constructor, despite the value-complement meaning of `not(T)`. Do not impose kind-positivity rules on existing ground recursion. CDuce’s distinction between structural constructors and Boolean connectives is a useful precedent for guarding, but does not determine Horos’s separate value-descent rule. [cduce.org](https://www.cduce.org/manual_types_patterns.html?utm_source=chatgpt.com)

### A complete ordinary `Data`

Using ordinary node-family abbreviations:

```text
GroundLeaves =
    Any
  | AtomFamily
  | Scalar
  | Exact

Data = ν D.
    GroundLeaves
  | Tuple<D>
  | Compound<D>
  | Seq<D>
  | List<D>
  | Set<D>
  | Map<D,D>
  | ClosedRecord<D>
  | OpenRecord<D>
  | UnionOf<D>
  | AndOf<D>
  | Not<D>
```

`AtomFamily` expands to `atom` and every existing atom-refinement case. `Scalar` includes every existing text-refinement case. These groups must expand into ordinary node selectors; there must be no hidden “data leaf” flag.

Variadic families constrain every member of a finite child sequence, including an empty sequence. `Exact` has no observed type children.

`rec` and `ref` need no alternatives after graph lowering. Their original lexical checks and formation obligations still apply.

Separate three immutable inputs to avoid a bootstrap dependency cycle:

| Input | Responsibility |
|---|---|
| Shape signature `Σ` | Forms, metadata validation, child positions, binding and descent rules |
| Ordinary prelude `P` | Kind definitions against `Σ`, including `Data` |
| Formation profile `F` | Pins `Σ` and `P`; assigns ordinary kind contracts to operands |

Modules pin `Σ`; the resolver’s explicit configuration pins `F`. All participate in receipts. The prelude does not depend on the formation profile.

The profile preserves current admission:

| Forms | Operand requirement |
|---|---|
| Ground leaves, `exact` | No type operands |
| Tuple, compound, list-of, closed record, union | Proper operands |
| Wire | Proper message type |
| Seq, set-of, map-of, open record, and, not | Operands satisfy the pinned ordinary `Data` definition |
| Rec | Existing guard validation; body satisfies that same definition |
| Ref | Valid lexical backedge |

Renaming the pinned definition changes no rule. A user declaring a different kind named `Data` changes no primitive contract.

### Required counterexamples

| Case | Judgment |
|---|---|
| `list<wire<Person>>` | In `List<⊤P>`; outside `Data` |
| Closed record with Wire field | In `Record<⊤P>`; outside `Data` |
| `union(text, {name:text})` | In `UnionOf<Data>`; outside `Scalar \| Record<⊤P>` |
| Record with one scalar and one Wire field | In `Record<Scalar \| WireFamily>`; outside `Record<Scalar> \| Record<WireFamily>` |
| `exact(D("param", I(0)))` | In `Exact` and `Data`; payload remains opaque |
| `union()` and `and()` | In `Data` |
| `wire<union()>` | In `WireFamily`; outside `Data` |
| Empty closed record | In `Record<K>` for every proper kind `K` |
| Optional Wire field | Prevents `Data` membership |
| Wire field plus impossible required field | Still structurally outside `Data` |
| `not(text)` | In `Data`; unrelated to kind complement |
| `ν X.X` | Invalid kind definition |

Removing only outer Wire shapes from `⊤P` would still admit nested Wire. No kind-difference operator is needed for the proposed `Data`.

## 2. Constructors and general result information

### Arrow contracts, variance, and saturation

An arrow

\[
(A_1,\ldots,A_n)\Rightarrow B
\]

means that a constructor admits every argument tuple in those input kinds and returns a result in `B`. It promises nothing about other inputs.

The sufficient variance rule is:

\[
(A_i)\Rightarrow B \;\le\; (C_i)\Rightarrow D
\quad\text{if}\quad
C_i\le A_i\ \text{for every }i,\quad B\le D.
\]

Apply it recursively to constructor-valued inputs and results.

Permit constructor results now:

```text
Pick<F: (⊤P => ⊤P)> : (⊤P => ⊤P) = F
Pick<Page><text>
```

Each application saturates one declared parameter list. Empty application remains invalid. Partial application and captured anonymous constructors can be separate term-language extensions.

Track semantic sort independently of declaration parameter count: a zero-parameter alias `Alias = Page` denotes a constructor while being referenced directly. GHC likewise documents kind signatures and application arity as separate concerns; its rules are a precedent for making that distinction explicit. [Glasgow Haskell Compiler 9.14.1 User's Guide](https://downloads.haskell.org/ghc/latest/docs/users_guide/exts/poly_kinds.html?utm_source=chatgpt.com)

### Known constructors: check the retained body

To establish a requested arrow for a known constructor:

1. Prove that each requested input kind is admitted by the corresponding **original parameter contract**.
2. Introduce fresh rigid parameters with those requested kinds.
3. Check the retained body under those assumptions.
4. Establish the requested result kind through ordinary kind checking.

Every body must already have passed checking under its original declaration. A narrower query cannot repair an invalid open definition.

For every suitable closed kind `K`, Page’s body can establish:

```text
Page : K => Record<List<K> | Text>
```

Here `Text` is the structural family containing its concrete `text` descriptor. This yields:

```text
Page : Data   => Data
Page : Scalar => Rows
```

It does **not** generally yield `Flat => Rows`. Taking `T = list<text>` makes `items` have type `list<list<text>>`, which is outside Flat.

`Ignore` establishes a scalar result for every admitted input. `Send` establishes `Wire<K>` for a `K` message. Both still validate and admit every supplied argument.

This is general result reasoning: the retained object is the checked expression and its lexical environment, not a Data-specific summary language.

### Abstract constructors: use their contracts

For:

```text
F : Flat => Rows
A : Flat
```

derive:

```text
F<A>             : Rows
{value: F<A>}     : Record<Rows>
```

The wrapper does not establish Rows: that would require its field type to belong to Flat. A constructor constantly returning an empty record supplies a concrete counterexample.

If `F` guarantees only `Flat => Flat`, `F<A>` cannot establish Rows: it may return text.

Likewise, an appropriately declared composition with `G : A => B` and `F : B => C` establishes `A => C`.

Preserve original higher-order admission. The existing:

```text
Wrapped<F: ⊤P => ⊤P, T: ⊤P>
```

cannot accept a constructor guaranteed only `Data => Data`. A preservation query must supply both guarantees:

```text
F : (⊤P => ⊤P) & (Data => Data)
T : Data
```

Alternatively, the author must declare a different input contract.

### Joint arrows and alternative arrows

Use finite disjunctive normal form, explored lazily:

```text
union of alternatives
each alternative = intersection of arrow contracts
```

For an abstract application:

1. Split top-level unions of certified argument bounds into cases.
2. Consider every argument case and every possible callee alternative.
3. Collect all arrows in that alternative whose entire input tuple admits the arguments.
4. Require at least one applicable arrow in every case.
5. Intersect applicable result guarantees within a case; union results across cases.

Thus:

```text
(A => B) & (A => C)     applied to A guarantees B & C
(A => B) | (A => C)     applied to A guarantees B | C
```

An alternative `(A => B) | (C => D)` does not generally accept an `A | C` argument: the unknown constructor alternative may reject it.

Do not distribute `Record<A | B>` into `Record<A> | Record<B>`. When splitting assumptions about a repeated rigid parameter, retain its identity consistently across occurrences.

For constructor-kind inclusion, require each source alternative to establish a target alternative. Check each target arrow by symbolic application to fresh target-domain inputs and proving the result obligation. Recursion through arrows decreases the finite sort structure.

These are finite, sound proof rules. They do not decide every true set-theoretic relation between constructors.

### Quantification and principal kinds

Kind parameters in definitions such as:

```text
kind Preserves<K> = K => K
```

are ordinary definition parameters. They do not assert `∀K. K => K`.

The first release needs no quantified constructor schemes. Known bodies can establish each requested instantiated contract; abstract constructors expose finite joint guarantees.

The cost is real: an abstract identity cannot advertise preservation of every future kind in one finite contract. Add explicitly scoped rank-one schemes when a concrete consumer requires that capability across an abstraction boundary.

Do not promise a principal inferred regular kind. Consider:

```text
Duplicate<T> = tuple(T, T)
```

Its exact result family requires equality between arbitrary subtrees, which general regular-tree languages cannot express. Retain symbolic structure and prove requested memberships instead of forcing all information into one inferred kind.

## 3. Algorithms, boundedness, identity, and extension

### Checking pipeline

Use one authoritative sequence:

1. **Raw admission:** canonical structure, metadata, sorts, slots, import pins, and resource limits.
2. **Dependency checking:** include every reference and argument, even ignored ones. Reject named cycles; admit only explicit, valid recursive-kind binders.
3. **Ground-binding validation:** resolve `rec`/`ref` and establish the existing strict-descent condition.
4. **Kind compilation:** expand acyclic kind definitions into shared structural states and finite sequence automata.
5. **Open-body checking:** discharge primitive formation contracts and declared results using rigid inputs.
6. **Application and queries:** preserve closure scopes and recheck requested memberships through the same machinery.
7. **Publication:** publish only after every required obligation succeeds.

A kind name supplies identity, not trust. Declaration identity remains `(module digest, name)`; parameter identity adds the ordered slot. Equivalent definitions may establish mutual inclusion without becoming the same declaration.

### Concrete structural membership

For a checked closed descriptor graph, create Boolean obligations indexed by:

```text
(descriptor node, instantiated kind state)
```

Node tests compare form and selected metadata; sequence transitions check actual finite children. Union produces disjunction and intersection produces conjunction.

Solve recursive obligations at their greatest fixed point: initialize candidate pairs as true, then propagate failures to stability. Solve the complete equations for each component. Encountering a cached pair does not itself prove membership or validate recursion.

With unlimited resources, this decides membership in the proposed proper-kind fragment.

### Structural inclusion

Define a free observation universe `O` containing all admitted descriptor observations. It allows regular, potentially infinite vertical unfoldings, but each child sequence is finite. The finite selector vocabulary determines the alphabet abstraction.

Decide:

\[
K\le_O L
\quad\Longleftrightarrow\quad
\llbracket K\rrbracket_O\cap
\overline{\llbracket L\rrbracket_O}
=\varnothing.
\]

A concrete implementation route is:

- Encode child sequences with explicitly finite sibling spines.
- Translate structural states and sequence automata to alternating parity tree automata.
- Dualize the right-hand automaton for inclusion.
- Solve the resulting emptiness problem.

Complement is internal machinery; public kind negation is unnecessary.

Finite XML-type inclusion supplies relevant techniques, but does not alone justify infinite descriptor unfoldings. Schmidt-Schauß, Sabel, and Schütz establish decidable, DEXPTIME-complete inclusion for a related greatest-fixed-point language over also-infinite trees. The Horos translation, metadata treatment, and sequence encoding still require their own correctness argument. [cis.upenn.edu](https://www.cis.upenn.edu/~bcpierce/papers/regsub-toplas.pdf?utm_source=chatgpt.com)

The modal μ-calculus/alternating-parity-automaton connection supplies the more general algorithmic route, including internal least fixed points needed to enforce finite sibling lists. [arxiv.org](https://arxiv.org/pdf/2207.00517?utm_source=chatgpt.com)

**Qualification:** inclusion over `O` is sufficient for inclusion among admitted Horos types. Failed inclusion may expose an unrealizable shape, so it does not automatically establish an actual type counterexample.

### Result statuses

| Status | Meaning |
|---|---|
| Proven | The requested guarantee is established |
| Refuted | Concrete closed nonmembership or an independently validated admissible counterexample |
| Unproved | Open assumptions or qualified proof rules are insufficient; an abstract counterexample may be attached |
| Invalid | Malformed syntax, binding, sort, ownership, or request |
| Exhausted | A resource budget prevented completion |

A required obligation returning Unproved prevents publication. Exhaustion never becomes Refuted.

Without resource limits, concrete structural membership and free-shape inclusion terminate. The stated higher-order proof procedure also terminates within its restrictions, while remaining incomplete for semantic constructor inclusion.

### Budgets and extension

Retain separate raw-admission, aggregate-work, and expanded-output budgets. Charge automaton states, sequence states, Boolean alternatives, environments, cache keys, proof comparisons, and shared structure construction. Expanded output charges every repeated occurrence.

Use operation-local caches and immutable Programs. Apply depth limits to parsing and semantic nesting, with iterative worklists for graph algorithms. Preserve the brief’s default limits and absolute depth ceiling.

Adding `Rows`, renaming `Data`, or importing another recursive category requires only canonical definitions. Adding a primitive form changes `Σ`, its metadata and child-position semantics, and its formation profile; that is a separately versioned Horos change. Bittype’s generic engine acquires no corresponding semantic branch.

## 4. Migration, semantic facts, native consumers, and qualification

### Horos replacement and API boundary

Replace `Claim`, `Interface`, `DataSummary`, the summary evaluator, `descriptorData`, and `RequiresData` with:

- Sorts.
- Canonical kind expressions and compiled kind handles.
- Checked symbolic type-level values.
- Ordinary membership/inclusion queries.
- Primitive operand contracts containing ordinary kind expressions.

A possible public shape, after those semantics are fixed:

```go
func (p *Program) CompileKind(expr core.Value) (KindHandle, error)

func (p *Program) Satisfies(
    value TypeValue, expected KindHandle,
) (Verdict, error)

func (p *Program) Includes(
    actual, expected KindHandle,
) (Verdict, error)

func (p *Program) Apply(
    callee TypeValue, arguments []TypeValue,
) (TypeValue, error)

func (p *Program) GroundDescriptor(
    value TypeValue,
) (core.Value, error)
```

Handles remain sealed and Program-owned. `Apply` may return a constructor. `GroundDescriptor` requires proper sort, ordinary membership in the formation profile’s pinned ground kind, materialization, and existing ground WF.

The release must establish this bridge property, independently of operational resource limits:

> For an admitted closed expression, membership in the specified ordinary Data definition is equivalent to its fully materialized descriptor being accepted by unchanged ground Horos.

Both checks may separately exhaust. A subset proof never rewrites descriptor syntax.

Ground Ontos and ground Horos semantics remain unchanged. The declaration implementation remains Go-only until separately ported.

### Bittype evidence

Keep the generic semantic graph unchanged. Replace Horos-owned syntax, projections, facts, and native readers.

Publish separately:

- Semantic sort.
- Declared kind contract.
- Checked body or symbolic result.
- Individual established memberships.
- Required-obligation completion.

An entity may have multiple membership facts. Their absence does not assert nonmembership. No single fact should be advertised as a principal kind.

Evidence introduced under stronger query assumptions must retain those assumptions. A proof of `Page : Scalar => Rows` must not become an unconditional Rows fact on Page’s body under its original unrestricted parameter.

“Complete evidence” should mean that every required obligation and every published guarantee is established against the exact snapshot. It cannot mean enumeration of every true membership.

Receipts must pin kind definitions, schema, prelude, formation profile, resolver configuration, and supplied dependencies. Keep explicit resolution-to-analysis handoff and projection correspondence checks. Membership never authorizes selecting Function, Stream, Cell, or other behavioral-library implementations.

### Native generation

The native reader should consume declared contracts, established memberships, ordinary applications, and the checked structural information required by its explicit representation profile.

Separate two profile promises:

- **Representation:** generated types represent Horos-checked uses.
- **Enforcement:** arbitrary native clients are prevented from violating the Horos constraint.

Erasing a richer kind annotation is permitted only when the selected profile explicitly authorizes the first promise. If native enforcement is required and unavailable, return located Unsupported.

Preserve one open Page declaration and ordinary applications. Retain whole-exported-closure rejection for residual constructor parameters, unused proper parameters, and unsupported forms. Keep canonical integer mapping Unsupported in the existing profile. Do not silently omit declarations.

The replacement remains types/interfaces generation. Richer native encodings should begin only when a concrete required declaration cannot pass the existing profile and a target-specific representation contract has been specified.

### Qualification that could falsify the design

Run new tests on the maintained remote/shared gates, as requested.

| Area | Required falsification cases |
|---|---|
| Generality | Rename and import Data-equivalent definitions; add unrelated Flat/Rows/Tree variants; ensure no name-sensitive behavior |
| Structural membership | Independent small-graph oracle, including exhaustive fixed-point valuation checks |
| Inclusion | Independently validate witnesses; distinguish abstract shapes from admissible canonical descriptors |
| Recursion | Reject `νX.X`, growing kind instantiations, and unguarded ground refs; preserve valid ground recursion through complement |
| Metadata and literals | Opaque exact payloads, every scalar/refinement case, open tails, optional fields, invalid metadata |
| Constructors | Joint versus alternative arrows, variance, constructor results, constructor aliases, original-bound admission |
| Correlation | Repeated parameters and examples where a regular approximation loses their relationship |
| Integrity | Ignored arguments, foreign Program handles, changed pins/configuration, stale receipts, concurrent reads |
| Resources | Exponential alternatives, repeated output, cache/environment charges, separate exhaustion diagnostics |
| Native output | Four pagination outputs and compilers, Unsupported/no-output cases, independently installed public Go consumer |

The independent oracle must not reuse production transition construction. Preserve historical evidence as history; qualify the replacement on its own exact revision.

Ship one new declaration revision and one Horos language implementation. Remove the revision-3 reader and dedicated Data route together.

## 5. Open assessment and strongest limitations

The main missing distinction is between **general definability** and **complete inference**. Users can define an open-ended structural kind vocabulary without the checker deciding every semantic implication or producing a principal kind.

The simplest sufficient construction is therefore the one above: positive recursive structural schemas, finite constructor contracts, and retained checked bodies. Universal schemes, arbitrary refinement predicates, rows, and type-level lambdas are separable extensions.

Two counterexamples define its limits particularly clearly.

### Unrealizable structural observations

Consider:

```text
kind EndlessWire = ν X. Wire<X>
```

The free shape universe contains an infinite Wire chain. Current Horos admission cannot realize it: named declaration cycles are forbidden, and ground recursive bodies must satisfy Data.

Therefore the structural inclusion engine cannot establish:

```text
EndlessWire ≤ ⊥P
```

even though no currently admitted proper expression inhabits EndlessWire. Reporting an actual Horos counterexample here would be a bug. Eliminating this conservatism requires integrating descriptor realizability into inclusion and qualifying that larger algorithm.

### Lost relationships between repeated subtrees

`Duplicate<T> = tuple(T,T)` preserves an equality relationship that ordinary regular shape schemas cannot generally express. A coarse structural result can lose useful information.

Retaining the symbolic body avoids losing it prematurely; it does not make every relational query decidable. If a future requirement needs arbitrary field-preserving transformations or equality between unknown subtrees, that requirement should explicitly trigger a richer calculus.

These boundaries still allow the owner’s central requirement to be met: **Data becomes one ordinary recursive kind, and every required Scalar/List/Record, Flat/Rows/Tree, Page, and abstract-constructor example is checked through the same general machinery.**
