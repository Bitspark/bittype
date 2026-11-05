# Logos proofs for native type preservation

**Proposal, 5 November 2026. Documentation only.** The owner proposed expressing
Horos requirements, native guarantees and their implication in Ontos/Logos. This
extends the [paired consultation synthesis](generic-elaboration.md); it does not
change an implemented provider, transformer, profile or dependency. The examples
below have not been executed as Logos programs or compiled as Go. No generator,
proof bridge, native semantic analyzer or runtime binding is implemented here.

## The claim to prove

For a fixed source world, generated declaration, mapping configuration and native
environment, prove:

```text
for every native argument in the declared consumer domain:
    native acceptance
        implies validity of its declared Horos application
        and the selected representation of that application's result
```

The interpretation must identify the intended source application. Finding some
other valid Horos type with the same native payload does not establish this claim.
Keep source membership, selected representation and runtime value conformance
separate. This proposal concerns type applications; it does not establish that
arbitrary runtime values satisfy Horos, or that Function/Stream/Cell behave
correctly.

Sound admission and relative completeness are separate obligations. Completeness
requires a native application for every source application in an independently
declared supported domain. An encoding that rejects everything is insufficient.
A finite domain is legitimate when explicitly requested; it cannot silently
replace a request to support every argument admitted by a canonical kind.

An open generic needs a theorem under its original parameter assumptions, or a
coverage proof for an explicitly complete finite domain. Several successful
ground applications do not establish the open theorem. A ground Logos conclusion
such as `bound_sound(context, source, target)` can name a universal property, but
its rules must justify that property. Naming it or asserting it as a fact is not
such a justification.

## Ownership and connection to the pipeline

| Component | Proposed responsibility |
| --- | --- |
| Horos integration | Supply established source judgments with their real binders, assumptions, subjects and source dependencies |
| Native language integration | Describe native acceptance and type identity in a specified fragment, compiler/language configuration and consumer-operation domain |
| Type-owned mapping | Supply representation equations, composition obligations and their derivations for the actual selected native form |
| Logos integration | Derive/check the requested claims against those programs, facts and explicitly admitted primitive theories |
| Bittype orchestration | Keep the source, candidate AST, mappings, requirements and proof result associated with the same invocation |

Ontos carries ground data; Logos supplies terms with variables, clauses and
accounts. These predicate names and rules belong to consumer programs. Neither
the Logos kernel nor Bittype's generic infrastructure gains a Data, Go, Page,
Function, Stream or Cell case. A contract library may supply its own mapping
laws; the native integration supplies language laws. Exact mapping selection and
explicit overrides continue to precede generation.

The proposed flow is:

```text
captured source + explicit source semantic evidence
        + selected mappings and requested guarantees
    -> candidate native AST
    -> explicit native analysis and representation obligations
    -> Logos account for the requested preservation claims
    -> admit the candidate output under those claims, or report failure
```

This is requested library behavior, not an additional repository delivery check.
No hidden source analysis is invoked from a mapping. A future integration must
request the relevant native/source capabilities explicitly. Existing semantic
receipts associate evidence with captured inputs; they do not prove preservation.
Generated associations and accounts belong to invocation/result state, not to the
immutable source Semantics graph. The existing ASTSet output contains syntax
only; transporting accounts would require an explicit result/evidence contract.
This document does not add that SPI or pretend the current adapter carries it.

Native facts must be extracted from the actual candidate AST and resolved
dependencies. A plugin's assertion `enforces_data` is not an acceptable substitute.
The native model must cover every consumer operation admitted by its promise,
including applicable aliases, inference and extension mechanisms. A model that
omits a permitted bypass cannot establish the promised guarantee.

Logos checks reasoning relative to admitted rules and premises. A native language
model remains a foundational assumption unless its adequacy is established
separately; a proof about a model does not independently prove a compiler correct.
This is consistent with [Logos's account boundary](https://github.com/Bitspark/logos/blob/main/README.md)
and [domain-blind ownership](https://github.com/Bitspark/logos/blob/main/docs/architecture/decisions/0025-logos-is-domain-blind.md).

## Worked example: an explicitly bounded Go pagination view

The canonical declarations use the existing Horos syntax; `h` denotes an explicit
import of the standard kinds module:

```text
type Person = { name: text; active: bool; };
type Page<T> = { items: list<T>; next?: text; };
type DataPage<D: h.Data>: h.Data = Page<D>;
```

For this example the consumer independently requests the source argument domain
`Person`, `text`, and `bool`. This is a native view of DataPage over those three
types, not a definition of Data or complete native support for it. The illustrative
target uses Go 1.25 language semantics and a fixed generated package:

```go
package pagination

type Person struct {
    Name   string
    Active bool
}

type DataItem interface {
    Person | string | bool
}

type Page[T any] struct {
    Items []T
    Next  *string
}

type DataPage[T DataItem] = Page[T]
type PeoplePage = DataPage[Person]
```

The candidate exploits an exact finite union, with no `~` terms, method-based
seal or generic private evidence constructor. The native-domain lemma concerns
closed concrete type arguments at instantiation. Aliases are normalized to their
actual Go type identities; an alias of Person adds no new case. A distinct
defined type with Person's underlying structure is a different identity.
Symbolic forwarding through other Go generic declarations requires an additional
constraint-inclusion theorem and is outside this first example's claim.

The language basis is Go's [general-interface type sets](https://go.dev/ref/spec#General_interfaces),
[constraint satisfaction](https://go.dev/ref/spec#Satisfying_a_type_constraint),
and [instantiation](https://go.dev/ref/spec#Instantiations).
Parameterized aliases are a standard language feature from
[Go 1.24](https://go.dev/doc/go1.24#language). A future model must identify its
actual supported language/compiler configuration; these links alone are not an
implemented model or a compilation result.

The mapping selects slices for lists, pointers for optional fields and exported
Go fields for canonical field names. At the type-representation level, `Next`
stands for the optional text field. This makes no claim that every Go string is
valid UTF-8, that pointer mutation preserves a value invariant, or that a Go value
has been validated as a Horos record.

### Concrete input facts

The following is proposed Logos consumer-program notation using ordinary
predicates, compounds and lists. Lowercase atoms abbreviate resolved identities;
uppercase identifiers are logical variables. `c` names one immutable context
containing both captured inputs, their dependencies and the selected mapping.
It is not a new digest or provenance mechanism.

These are **expected bridge inputs**, not observations from an executed proof.
Horos facts require source evidence; native facts require actual AST extraction.
Configuration facts specify intent and never independently prove enforcement.

```prolog
context(c).
native_model(c, go_1_25_exact_union).

% Established Horos facts; h_data is the resolved ordinary imported kind.
h_bound(c, h_data_page, h_data).
h_result_bound(c, h_data_page, h_data).
h_alias(c, h_data_page, h_page).
h_arrow(c, h_page, h_data, h_data).
h_member(c, h_person, h_data).
h_member(c, h_text, h_data).
h_member(c, h_bool, h_data).

% Consumer scope and intended interpretation, fixed before proving coverage.
source_domain(c, [h_person, h_text, h_bool]).
interpretation_table(c, [pair(g_person, h_person),
                         pair(g_string, h_text), pair(g_bool, h_bool)]).
selected_representation(c, h_person, g_person).
selected_representation(c, h_text, g_string).
selected_representation(c, h_bool, g_bool).
declaration_mapping(c, h_data_page, g_data_page).

% Complete normalized observations of the actual target declarations.
g_exact_union(c, g_data_item, [g_person, g_string, g_bool]).
g_alias(c, g_data_page, g_data_item, g_page).
g_unconstrained_parameter(c, g_page).
```

`g_exact_union` is issued only for a well-formed complete finite union of exact,
resolved, non-interface types. An arbitrary list of selected cases, an unresolved
type name, a `~` term or a partial scan cannot produce that fact. `g_alias` records
the complete one-parameter constrained alias and its unchanged argument use in
the base application. The declared native model supplies the admission lemma
for this exact syntactic fragment; merely inserting `native_model` does not
install an implementation of that lemma.

`selected_representation` records a choice. It does not assert that Person's
fields were generated correctly. Result-representation obligations below cover
the selected element mappings and the Page body separately.

### Rules for admission and the source result

```prolog
% Ordinary list traversal; the term fragment grounds each proof instance.
member(X, [X|Rest]).
member(X, [Head|Rest]) :- member(X, Rest).

% Native language lemma specialized to the admitted fragment.
native_domain(C, G, Ns) :-
    native_model(C, go_1_25_exact_union),
    g_alias(C, G, Constraint, Base),
    g_exact_union(C, Constraint, Ns),
    g_unconstrained_parameter(C, Base).

native_accepts(C, G, N) :-
    native_domain(C, G, Ns), member(N, Ns).

interpretation(C, N, S) :-
    interpretation_table(C, Rows), member(pair(N, S), Rows).

source_accepts(C, H, S) :-
    h_bound(C, H, K), h_member(C, S, K).

source_result(C, H, S, Kout) :-
    source_accepts(C, H, S),
    h_alias(C, H, Base),
    h_bound(C, H, Kin), h_arrow(C, Base, Kin, Kout).

safe_argument(C, H, G, N) :-
    declaration_mapping(C, H, G),
    native_accepts(C, G, N), interpretation(C, N, S),
    selected_representation(C, S, N), source_accepts(C, H, S).
```

The admission characterization has two directions in the native model: every
concrete Go type accepted by this exact constraint is in `Ns`, and each listed
type can instantiate it. The first direction is essential. Merely writing the
positive `native_accepts` rule does not establish that Go accepts nothing else.
This model-adequacy obligation is stated explicitly below.

`h_arrow` is a checked Horos conditional result guarantee. It includes the whole
Page body, including `next`, and its original parameter domain. It is not a rule
that every collection preserves membership in every kind. All these rules carry
the context; additional symbolic facts would also carry the full assumption
environment and real binder identities.

### A finite coverage account, not sample testing

```prolog
% A functional intended interpretation, using ground structural inequality.
different_keys(N, []).
different_keys(N, [pair(M, S)|Rest]) :-
    N \== M, different_keys(N, Rest).

unique_keys([]).
unique_keys([pair(N, S)|Rest]) :-
    different_keys(N, Rest), unique_keys(Rest).

all_safe(C, H, G, []) :- declaration_mapping(C, H, G).
all_safe(C, H, G, [N|Rest]) :-
    safe_argument(C, H, G, N), all_safe(C, H, G, Rest).

bound_sound(C, H, G) :-
    declaration_mapping(C, H, G),
    interpretation_table(C, Rows), unique_keys(Rows),
    native_domain(C, G, Ns), all_safe(C, H, G, Ns).

represented_source(C, H, G, S) :-
    safe_argument(C, H, G, N), interpretation(C, N, S).

all_represented(C, H, G, []) :- declaration_mapping(C, H, G).
all_represented(C, H, G, [S|Rest]) :-
    represented_source(C, H, G, S), all_represented(C, H, G, Rest).

domain_complete(C, H, G) :-
    source_domain(C, Ss), all_represented(C, H, G, Ss).
```

The recursion consumes the complete lists. `bound_sound` covers every identity
the native-domain lemma admits, not just the example PeoplePage. `domain_complete`
covers the independently supplied source list; it does not redefine that list
as whatever Go happened to accept. `unique_keys` prevents choosing whichever
source interpretation makes a desired conclusion true.

This is the positive Horn term/list fragment with ground structural inequality.
Clauses with variables are instantiated in a finite ground account. No use of
negation-as-failure converts missing membership into nonmembership, and no
absence of a counterexample is treated as a universal proof. Selecting and
implementing a serving Logos integration is future work.

### The representation obligation

The rules above prove an admission property relative to their model. They do not
yet prove that the generated body represents Page. That is a separate obligation
over the actual source and native ASTs. For this example the normalized body
observations would include:

```prolog
h_page_shape(c, h_page, ht,
    record([required(items, list(parameter(ht))), optional(next, h_text)])).
g_page_shape(c, g_page, gt,
    record([field(g_items, slice(parameter(gt))),
            field(g_next, pointer(g_string))])).
field_binding(c, h_page, items, g_page, g_items).
field_binding(c, h_page, next, g_page, g_next).
```

Here `g_items` and `g_next` resolve to the actual Go Items and Next fields.
The selected structural mapping supplies a composition rule of this shape:

```prolog
page_representation_law(C, HBase, GBase) :-
    h_page_shape(C, HBase, HB,
        record([required(items, list(parameter(HB))), optional(next, h_text)])),
    g_page_shape(C, GBase, GB,
        record([field(g_items, slice(parameter(GB))),
                field(g_next, pointer(g_string))])),
    field_binding(C, HBase, items, GBase, g_items),
    field_binding(C, HBase, next, GBase, g_next),
    list_slice_law(C), optional_pointer_law(C),
    element_representation(C, h_text, g_string).

result_representation(C, H, G, S, N) :-
    safe_argument(C, H, G, N), interpretation(C, N, S),
    h_alias(C, H, HBase), g_alias(C, G, Constraint, GBase),
    page_representation_law(C, HBase, GBase),
    element_representation(C, S, N).
```

`list_slice_law`, `optional_pointer_law` and `element_representation` are **proof
premises to discharge**, deliberately not unconditional facts in the input block.
The primitive mapping library must define their precise type-representation
meaning and supply admitted axioms/derivations. Person's instance must account for
its required Name/Active field mappings. The optional-pointer law is scoped to
this type-level representation promise and is not a general value-conformance
law. Page composition then works under a symbolic element-association assumption;
the syntactic binders HB and GB remain distinct and correspond by this rule.

The displayed Page rule is a readable instance of the structural mapping's rule,
owned outside the kernel. The structural mapper would instantiate the supplied
record/list/optional laws for this body; an author need not write a Page-specific
mapping or proof rule. Production rules must use actual resolved source operands,
native nodes and field bindings; the spellings Page, Data, Items or Next must never
trigger generic runtime dispatch. The same record/list composition must work for
another declaration with the corresponding admitted structure.

The complete requested claim combines the independent obligations:

```prolog
all_results(C, H, G, []) :- declaration_mapping(C, H, G).
all_results(C, H, G, [N|Rest]) :-
    interpretation(C, N, S),
    h_result_bound(C, H, K), source_result(C, H, S, K),
    result_representation(C, H, G, S, N), all_results(C, H, G, Rest).

preserves_declared_domain(C, H, G) :-
    bound_sound(C, H, G), domain_complete(C, H, G),
    native_domain(C, G, Ns), all_results(C, H, G, Ns).
```

### Queries and expected proof outline

The requested queries, written separately from the program, are:

```prolog
native_accepts(c, g_data_page, g_person)
source_accepts(c, h_data_page, h_person)
source_result(c, h_data_page, h_person, h_data)
bound_sound(c, h_data_page, g_data_page)
domain_complete(c, h_data_page, g_data_page)
result_representation(c, h_data_page, g_data_page, h_person, g_person)
preserves_declared_domain(c, h_data_page, g_data_page)
```

For PeoplePage, an account would have the following dependencies:

```text
native union/alias observations + Go admission lemma
    -> native_accepts(DataPage, Person)
source bound + Person : Data
    -> source_accepts(DataPage, Person)
source acceptance + Page : Data -> Data
    -> source_result(DataPage<Person>, Data)
native acceptance + fixed interpretation + selected association + source acceptance
    -> safe_argument(DataPage, Person)
three safe_argument instances + complete native domain + unique interpretation
    -> bound_sound(DataPage)
three represented source instances + independent source domain
    -> domain_complete(DataPage, declared source domain)
Page AST shapes + primitive/element representation laws + safe argument
    -> result_representation(DataPage<Person>, Go DataPage[Person])
```

No account is produced by this document. The full preservation claim requires
result representation for **every** admitted argument, in addition to the bound
and completeness accounts. The `all_results` fold covers those three result
obligations; the single Person query above is an illustration, not a replacement
for that coverage. Missing primitive or element laws prevent completion of the
top-level account even when its bound-soundness subgoal is established.

## Obligations before any enforcing claim

| Obligation | What must be established |
| --- | --- |
| Source admission | Horos judgments concern the exact input, original parameter contracts and actual assumptions; no foreign Program handle or unproved query becomes a fact |
| Native-model adequacy | For the declared Go fragment, actual acceptance is equivalent to the extracted exact union; aliases and applicable inference cannot add an unmodeled concrete case |
| AST correspondence | Native facts describe the complete candidate declarations and dependency identities, including the bound, base use and all fields |
| Interpretation | Every admitted argument has the intended, functional source interpretation and the selected representation association |
| Bound soundness | Every admitted native argument satisfies the corresponding Horos argument contract |
| Relative completeness | Every source argument in the independent declared domain has a supported native argument |
| Result preservation | Primitive, element and structural mapping laws justify the entire result family, under symbolic assumptions where applicable |
| Composition and emission | Definition/use substitution agrees; changed imported fixed dependencies invalidate inapplicable laws; AST-to-source emission preserves the proved declaration |

The proof may treat a language-model law or primitive representation equation as
an explicit foundational axiom. Its authority and meaning remain visible. It must
not disguise the desired top-level preservation conclusion as a primitive axiom.
Likewise, proving facts for a candidate AST does not by itself prove its formatter,
printer or the native compiler implementation correct.

If the required domain includes all Data types, this finite example is incomplete.
For instance, `DataPage<DataPage<Person>>` is source-valid but its inner native
Page type is outside DataItem. That request must fail or use a separately justified
larger encoding. No automatic finite-domain fallback is permitted.

The unconstrained Go Page remains usable with other types, just like canonical
Page. Its existence does not grant a Data guarantee. Successful DataPage alias
instantiation then erases the alias identity: assigning a Page value is not proof
of passage through DataPage. Retained identity, arbitrary native generic
forwarding, cross-package extension and runtime conformance need separate claims.

Changing exact terms to `~` underlying-type terms, widening the union,
changing an element mapping, or altering the optional/list representation changes
the obligations. The former private singleton-evidence scheme is still withdrawn;
moving it into this proof vocabulary does not fix its inference hole.

## Existing foundations and deferred integration

Current public Horos exposes Satisfies, Includes and conditional membership facts
through its [Go declaration API](https://github.com/Bitspark/horos/blob/main/core/go/declarations/README.md).
It has no public Logos certificate export for these general-kind judgments.
Earlier private Horos proof integration concerned descriptor conformance, not
the native-preservation theorem proposed here.

A first bridge could admit Horos-established facts through an explicit
LeafVerifier bound to the supplied source evidence. That retains Horos's decision
procedure in the trusted foundation; it is not an independent Logos replay of
Horos reasoning. A stronger bridge could export certificates or bind a justified
pure Theory for selected judgments. Recursive Horos kinds use greatest-fixed-point
reasoning; ordinary least-model Horn recursion is not a semantics-preserving
replacement. No required missing fact triggers hidden reanalysis or a duplicated
Horos evaluator in generation.

Logos already separates Rule derivations, Builtin theories and admitted Leaf
premises in its [account contract](https://github.com/Bitspark/logos/blob/main/docs/spec/logos-contract.md).
Encoding a new claim as Ontos data does not mean that a packaged engine or checker
already supports its reasoning. Missing capability, an unproved goal and exhausted
search remain distinct from a counterexample. None permits an enforcing success.

The next design work is to specify the primitive representation laws, the native
fact extractor/model and the source bridge, then a generic evidence/result
interface. They must preserve the existing source-owned semantics and stage-first
layout. No new source packages, dependency pins, registry entries, provider-matrix
cells, proof workflows or runtime behavior are introduced by this proposal.
