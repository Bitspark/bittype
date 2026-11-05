# Logos proofs for native type preservation

**Proposal, 5 November 2026. Documentation only.** The owner selected a symbolic
pagination proof for any `T : Data`, with native admission proved separately.
This replaces the finite-argument example as the design of this document. It
extends the [paired consultation synthesis](generic-elaboration.md) and its
[pagination walkthrough](generic-elaboration.md#pagination-and-scope).
No provider, transformer, proof bridge, dependency or runtime binding changes.
The Go and Logos examples below are unexecuted proposal notation.

The [type-system modeling proposal](type-system-model.md) defines the next
design prerequisite: describe the languages' terms, judgments, binding and rules
before implementing this bridge. The concrete predicates below illustrate a
preservation problem; they do not replace a reusable language-theory model or
assert that such a model already exists. Its [Horos/Go walkthrough](type-system-model.md#worked-example-horos-textbox-and-go-constraints)
and [language coverage](type-system-model.md#how-the-five-languages-fit) connect
this example to the other native targets.

## Prove preservation under symbolic assumptions

For a fixed source world and selected mapping configuration C, the theorem is:

```text
for every source type T and native type N:
    T : Data and Rep_C(T, N)
        implies Page<T> : Data
        and Rep_C(Page<T>, GoPage<N>)
```

`Rep_C(T, N)` means an established valid selected representation association,
including a well-formed native type. It is stronger than a mapping's preference
for a type name. The theorem quantifies over arbitrary T and N satisfying its
premises; neither is selected from a catalogue. It does not assert that every
Horos Data type already has an implemented Go mapping. Establishing existence
of N for a promised source domain is a separate completeness obligation.

Introduce fresh symbolic types tau and nu, assume `tau : Data` and
`Rep_C(tau, nu)`, then compose the list, optional-field and record laws. Their
derivation need not inspect tau's structure or enumerate possible arguments.
Discharge those assumptions to obtain the theorem. Concrete uses substitute
their own membership and representation proofs into that same theorem.

The canonical declarations are unchanged; `h` denotes an explicit import of
the ordinary standard kinds module:

```text
type Page<T> = { items: list<T>; next?: text; };
type DataPage<D: h.Data>: h.Data = Page<D>;
```

Page itself remains unconstrained. Establishing its conditional Data guarantee
does not add a Data bound to Page. DataPage is an ordinary constrained alias.
Both declarations must already be valid under their original contracts; a
narrower query cannot repair an invalid declaration.

For consistency with the [existing Go representation choice](semantic-generation.md#pagination-subset),
this example selects the presence helper for optional fields:

```go
type BittypeOptional[V any] struct {
    Present bool
    Value   V
}

type Page[N any] struct {
    Items []N
    Next  BittypeOptional[string]
}
```

This is the representation family to reason about. Its `any` parameter does
not enforce Horos Data on native callers. No hypothetical Go Data constraint
or constrained alias is inserted. There are no executable bodies or adapters.

## Ownership and proof inputs

| Owner | Proposed contribution |
| --- | --- |
| Horos integration | Valid source declarations, established kind judgments and closure consequences, original assumptions |
| Native integration | Actual Go AST, resolved type/field/binder identities and applicable native type rules |
| Selected mappings | Admitted representation laws connected to the selected native forms and fixed dependencies |
| Logos integration | Scoped derivation, assumption discharge and independent checking of requested claims |
| Bittype orchestration | Association of source, candidate AST, mappings, requirements and proof result with the same invocation |

Ontos carries ground data; Logos supplies terms with variables, clauses and
accounts. These predicates are consumer-owned logic vocabulary. Neither general
kernel interprets Data, Page, Go, Function, Stream or Cell. A contract library
can supply its mapping laws through the same mechanism.

Native facts must describe the actual candidate AST and resolved dependencies.
Representation laws must apply to the selected mappings. Declaring
`preserves_data` or selecting an implementation ID establishes neither. Changed
list/text/optional choices require applicable laws for the new configuration.

### Concrete scoped facts

Lowercase atoms are ground identities; uppercase names in rules are Logos
variables. `c` identifies one immutable source/target/configuration context;
`p` identifies a proof scope. `rigid(tau)` and `rigid(nu)` are ground terms
encoding arbitrary symbolic types, not concrete types or Logos variables awaiting
discovery. Real implementations retain the actual source and native binders.

```prolog
context(c).
proof_scope(c, p).
requested_kind(c, p, k_data).
checked_kind(c, k_data).

% Paired declarations and actual bound parameters.
definition_pair(c, h_page, g_page).
source_parameter(c, p, h_page_t, rigid(tau)).
native_parameter(c, p, g_page_n, rigid(nu)).

% Local hypotheses, not unconditional global facts.
assumed(c, p, member(rigid(tau), k_data)).
assumed(c, p, representation(rigid(tau), rigid(nu))).

% Established consequences for the ordinary resolved source kind.
horos_member(c, h_text, k_data).
kind_closure(c, k_data, list).
kind_closure(c, k_data, closed_record).

% Admitted laws of the selected representations.
primitive_representation(c, h_text, g_string).
mapping_law(c, list_to_slice).
mapping_law(c, optional_to_presence).
mapping_law(c, closed_record_to_struct).

% Actual templates opened with their corresponding symbolic parameters.
source_body(c, p, h_page,
    record([required(items, list(rigid(tau))), optional(next, h_text)])).
native_body(c, p, g_page,
    structure([field(g_items, slice(rigid(nu))),
               field(g_next, presence(g_string))])).
field_binding(c, p, items, g_items).
field_binding(c, p, next, g_next).
```

These are expected bridge inputs, not results of a proof run. Their admission
has explicit responsibilities:

- The scope and parameter records identify fresh, correctly sorted symbols and
  their real binders; they do not grant kind membership.
- `definition_pair` identifies selected named families and corresponding formal
  layout. It does not equate Go Page with an anonymous struct or prove its body.
- Body facts come from capture-avoiding substitution into the actual captured
  templates. They are not replacement bodies invented by the proof producer.
- `kind_closure` denotes an established Horos consequence, not a configurable
  claim that every kind is closed under lists or records. The closed-record
  consequence covers arbitrary well-formed finite required/optional fields.
- Primitive representation and mapping laws require admitted equations or
  derivations with stated assumptions. They are not self-asserted capabilities.
  The presence law covers the actual helper and its selected type-representation
  meaning, without asserting runtime value conformance.
- Field identities resolve to the actual corresponding fields. `presence(...)`
  denotes application of the selected helper, not a new generic AST kind.

These source consequences do not redefine Data by a smaller constructor set or
replace its general calculus. Tau can denote any Data type, including one whose
construction is absent here, provided its representation premise is established.

## Rules for the source guarantee

```prolog
kind_in(C, P, T, K) :-
    proof_scope(C, P), assumed(C, P, member(T, K)).

kind_in(C, P, T, K) :-
    proof_scope(C, P), horos_member(C, T, K).

kind_in(C, P, list(T), K) :-
    kind_closure(C, K, list), kind_in(C, P, T, K).

fields_in(C, P, [], K) :-
    proof_scope(C, P), checked_kind(C, K).

fields_in(C, P, [required(Label, T)|Rest], K) :-
    kind_in(C, P, T, K), fields_in(C, P, Rest, K).

fields_in(C, P, [optional(Label, T)|Rest], K) :-
    kind_in(C, P, T, K), fields_in(C, P, Rest, K).

kind_in(C, P, record(Fields), K) :-
    kind_closure(C, K, closed_record), fields_in(C, P, Fields, K).
```

For the supplied facts, these derive:

```text
under p:
    tau : Data                                  assumption
    list<tau> : Data                            list closure
    text : Data                                 supplied Horos judgment
    { items: list<tau>; next?: text; } : Data    record closure
```

This walks the fields of one open body, not possible T values. K remains a
variable in the rules; the resolved Data identity appears only in supplied facts
and the requested conclusion. Another kind requires its own established closure
consequences, not a kernel branch.

## Rules for representation composition

```prolog
rep_in(C, P, T, N) :-
    proof_scope(C, P), assumed(C, P, representation(T, N)).

rep_in(C, P, T, N) :-
    proof_scope(C, P), primitive_representation(C, T, N).

rep_in(C, P, list(T), slice(N)) :-
    mapping_law(C, list_to_slice), rep_in(C, P, T, N).

fields_rep(C, P, [], []) :- proof_scope(C, P).

fields_rep(C, P, [required(HName, T)|Hs], [field(GName, N)|Gs]) :-
    field_binding(C, P, HName, GName),
    rep_in(C, P, T, N), fields_rep(C, P, Hs, Gs).

fields_rep(C, P, [optional(HName, T)|Hs],
                 [field(GName, presence(N))|Gs]) :-
    mapping_law(C, optional_to_presence),
    field_binding(C, P, HName, GName),
    rep_in(C, P, T, N), fields_rep(C, P, Hs, Gs).

rep_in(C, P, record(Hs), structure(Gs)) :-
    mapping_law(C, closed_record_to_struct), fields_rep(C, P, Hs, Gs).

symbolic_preservation(C, P, H, G, K) :-
    definition_pair(C, H, G), requested_kind(C, P, K),
    source_body(C, P, H, HBody), native_body(C, P, G, GBody),
    kind_in(C, P, HBody, K), rep_in(C, P, HBody, GBody).
```

This derives the correspondence between `list<tau>` and `[]nu`, the optional
text field and BittypeOptional[string], and the complete record/struct bodies.
Both field lists must be consumed. Missing, extra or differently represented
fields cannot be skipped. This rule uses corresponding field order; another
ordering needs an explicit lawful field correspondence.

The final rule is generic in source and target declarations. Page needs no
hand-authored proof rule: the structural mapping instantiates its ordinary
record/list/optional rules for the opened bodies. Actual binder correspondence
is retained instead of matching same-spelled parameter names.

Representative ground queries, separate from the program, are:

```prolog
kind_in(c, p, list(rigid(tau)), k_data)
rep_in(c, p, list(rigid(tau)), slice(rigid(nu)))
symbolic_preservation(c, p, h_page, g_page, k_data)
```

Their conclusions remain **under p**. No rule drops that scope, promotes a
hypothesis to an unconditional fact, runs hidden analysis or treats an unproved
goal as nonmembership.

## Assumption discharge and the generic theorem

The intended theorem-introduction rule is:

```text
Gamma = { tau : Data, Rep_C(tau, nu) }
Gamma |- Page<tau> : Data and Rep_C(Page<tau>, GoPage<nu>)
tau and nu are arbitrary fresh parameters; the proof is stable under substitution
-------------------------------------------------------------------------------
|- forall T,N. (T : Data and Rep_C(T,N))
       implies (Page<T> : Data and Rep_C(Page<T>,GoPage<N>))
```

The body derivation is one premise. Lifting it to named application families
also requires checked definition/binder correspondence and the selected
named-representation law: substituting represented arguments into the native
definition agrees with representing the source substitution. This preserves
nominal/alias promises without identifying every named type with its body.

The introduction rule must establish:

1. Tau and nu are fresh, correctly sorted parameters. Earlier premises contain
   no concrete choice for either symbol.
2. Recorded assumptions are exactly the theorem's premises. Additional
   assumptions remain in its statement; they cannot disappear at discharge.
3. Used rules and theory results remain valid under admissible substitutions.
   A name test, whitelist or property of a concrete example is insufficient.
4. Opening and closing are capture-avoiding, preserve repeated binders and keep
   caller-local symbols from escaping into helpers or another declaration.
5. The result carries the whole application association and conditional kind
   guarantee. Support for T alone is not support for Page<T>.

Logos's current account has ground conclusions and finite acyclic derivations.
Ground terms here encode a scoped symbolic judgment. A checked query containing
`rigid(tau)` is **not by itself** universal quantification: the explicit
introduction/discharge step is still required. A consumer-owned proof theory or
certificate rule would check it through the existing extension boundary. Its
implementation and exact evidence API remain deferred; no new served Logos
capability is claimed. The check must inspect the scoped derivation and premises,
not accept producer-supplied `fresh` or `parametric` Booleans as sufficient proof.

This is a reusable proof rule, not a Page-aware kernel or another Data evaluator.
Primitive laws remain visible foundational assumptions. A finished account proves
preservation relative to those laws, not their independent real-world truth or
the correctness of a compiler implementation.

## Reuse and nested pagination

For a later Invoice, establish its ordinary `Invoice : Data` judgment and
`Rep_C(Invoice, GoInvoice)`, then substitute those proofs into the theorem. Page's
definition and proof schema need no new Invoice case or whitelist entry.

The theorem also supports nested pagination:

```text
T : Data, Rep_C(T,N)
    -> Page<T> : Data, Rep_C(Page<T>,GoPage<N>)
    -> Page<Page<T>> : Data, Rep_C(Page<Page<T>>,GoPage<GoPage<N>>)
```

The second application uses the first application's **result** proofs, not the
original T evidence. An outer generic likewise retains its caller binder and
assumptions instead of capturing Page's formal environment.

The Data theorem does not restrict ordinary Page. `Page<Cell<Person>>` may use
the unconstrained representation theorem when a valid Cell mapping is supplied.
A Cell containing Wire cannot supply the Data premise for DataPage. These
contracts and their mappings remain library-owned.

This establishes type representation and source membership, not validation of
arbitrary Go values, UTF-8 correctness of every Go string, mutation safety or
implementation of a behavioral contract.

## Native admission is separate

Generic preservation starts with premises. An enforcing native boundary must
establish them for the arguments it admits:

```text
for every A allowed by the native consumer contract:
    NativeAccepts_C(Encoding, A)
        implies an established intended interpretation (T,N)
        and T : Data and Rep_C(T,N)
```

The interpretation is the declared one, not whichever source type makes the
conclusion convenient. A may be a payload, descriptor or another supported
encoding; no universal source/native/proof tuple is prescribed. Aliases,
inference, extension authority and checking points belong to this theorem.

There is deliberately no rule from `native_accepts` to `assumed(member(...))`.
A hypothetical premise is not evidence that a native caller satisfied it.
Go Page[N any] supplies a representation family but no Data gate. Its constraint
rules do not automatically interpret Horos kinds; see the
[Go type-set specification](https://go.dev/ref/spec#General_interfaces).
An enforcing request still needs a valid encoding and the separate admission
proof. The current enforcement profile remains Unsupported.

The [Go admission proposal](go-native-admission.md) supplies a concrete candidate:
function-shaped certificates indexed by source and native types, a private marker,
and a constrained Page certificate constructor. Independent compiler fixtures
exercise nesting, aliases, inference and forgery. They do not establish Horos
premises or replace the formal native-admission theorem.

Neither theorem inherently requires enumeration. Admission mechanisms can also
have symbolic proofs. The former finite union was one limited admission example,
not pagination's definition, the Data domain or a requirement of Logos. It is
not retained as the proof in this document.

For source-checked generation, explicit Horos evidence supplies the premises of
each generated use. This is a representation guarantee, not Go enforcing the
bound on arbitrary handwritten native code. An enforcing result composes the
independent admission theorem with generic preservation. Neither replaces the
other. Completeness is also separate: a conditional theorem over valid
representations does not establish their availability for every source type.

## Integration and outstanding obligations

The proposed flow keeps prerequisites explicit:

```text
captured source + source semantic evidence + selected mappings
    -> candidate native AST and explicit native/body facts
    -> scoped preservation account and assumption discharge
    -> independent native-admission account, when requested
    -> output with the established guarantees, or failure
```

No mapping invokes hidden source analysis. Semantic receipts associate evidence
with captured inputs but do not prove preservation. Generated associations and
accounts belong to invocation/result state, not immutable source Semantics.
ASTSet still contains syntax only; carrying accounts needs an explicit evidence
contract. None is added here.

| Obligation | Required scope |
| --- | --- |
| Source evidence | Actual source world, valid original declarations, resolved kinds and complete assumptions; no foreign Program handle or unproved fact |
| AST correspondence | Complete bodies, helper definitions, field links, selected dependencies and real formal binders |
| Mapping laws | Stated representation meaning, well-formedness and substitution stability; no self-asserted top-level preservation claim |
| Symbolic derivation | Kind and representation conclusions under explicit assumptions, accounting for every field |
| Assumption discharge | Freshness, complete dependencies, lawful introduction and named-family substitution |
| Native admission | Independent proof that native acceptance establishes the intended membership and representation premises |
| Completeness | Coverage of the independently promised source/native domain |
| Emission | Formatter/printer and native-model assumptions preserve the proved AST's meaning; no automatic compiler-correctness claim |

Public Horos supplies Satisfies, Includes and conditional membership facts in its
[Go declaration API](https://github.com/Bitspark/horos/blob/main/core/go/declarations/README.md),
but no public Logos certificate export for general-kind judgments. A first bridge
could admit established facts through a LeafVerifier, retaining Horos in the
trusted foundation. A stronger bridge could supply certificates or a justified
pure Theory. Neither closure-fact notation nor a Proven verdict implements that
bridge.

Horos recursive kinds use greatest-fixed-point reasoning. The sound closure
consequences used here do not translate that semantics into ordinary least-model
Horn recursion. Tau's structure stays abstract. Required source consequences
must be established or requested through an implemented provider capability.

Logos distinguishes Rule derivations, Builtin theories and admitted Leaf premises
in its [account contract](https://github.com/Bitspark/logos/blob/main/docs/spec/logos-contract.md).
Its [domain-blind boundary](https://github.com/Bitspark/logos/blob/main/docs/architecture/decisions/0025-logos-is-domain-blind.md)
keeps these consumer rules outside the kernel. Expressing a claim in Ontos does
not establish a serving engine/checker. Missing capability, unproved claims,
exhausted search and refutation remain distinct.

The next design work begins with the [shared type-system model](type-system-model.md)
and small Horos/native theories. Then specify the primitive representation laws,
actual-AST fact extraction and source bridge, followed by the generic
evidence/result interface. They must preserve the existing source-owned semantics
and stage-first layout.

No implementation, native analyzer, proof-execution workflow, provider-matrix
cell, dependency or runtime behavior is added. Relevant software tests and
typechecking belong to future implementation; this document claims no executed
proof or compiler result.
