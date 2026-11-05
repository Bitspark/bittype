# General structural kinds

**Implemented in Go source.** This replaces the closed
Type/Data bound system. Horos owns the canonical definitions and their meaning;
Bittype projects them into syntax and provider-owned semantic facts. The
[research decision](../research-docs/0005-general-kinds.md) explains the expert
recommendations and adaptations. The canonical grammar belongs in
[Horos's declaration specification](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md).

## What a kind describes

A type describes values. A kind describes types or type constructors. For
example, the type `text` describes text values, while the kind `Scalar` describes
the scalar type descriptions selected by its ordinary definition. A kind can
constrain the nested types as well as the outer form:

```text
kind Flat = Scalar | List<Scalar>
kind Rows = Record<Flat>
kind Tree = nu X. Scalar | List<X> | Record<X>
```

These equations show the meaning; the [source grammar](horos-syntax.md) supplies
the concrete spelling. `Rows` includes closed records whose field types are
scalar types or lists of scalar types. `Tree` also permits deeper nesting. Neither
definition needs an implementation branch or registration in the semantic
engine.

Kinds classify **checked structural type descriptions**. They do not infer
properties from all possible runtime inhabitants. A type-level union has a union
root; it does not thereby belong to the kind union `Scalar | Record`. A record
containing a Wire field stays outside `Data` even when another field makes it
impossible to inhabit.

Type and kind expressions occupy different levels:

| Expression | Meaning |
| --- | --- |
| `list-of(text)` | One type describing lists of text values |
| `List<Scalar>` | A kind containing list types whose element type is scalar |
| A union of two types | One type with a union descriptor |
| A union of two kinds | Every type or constructor belonging to either kind |
| `Page<Person>` | Apply a type constructor to a type |
| `Record<Flat>` | Apply a kind definition to another kind |

## Sorts keep the levels well formed

Kinds have finite simple sorts:

```text
Sort = P | (Sort, ...) -> Sort
```

The proper sort is the domain of types. An arrow sort is the domain of
constructors with a particular input/result shape. It distinguishes a field type
from a constructor awaiting arguments. It imposes no data-only rule.

Every sort has top and bottom. An ordinary prelude declaration defines `Type`
as top at the proper sort. `Type` is not a reserved semantic name. A local kind
with that spelling has the same binding rules as any other name.

Ground Horos `any` is one type description, not kind top. Likewise, the empty
type union describes no values but is still a type; kind bottom contains no
types. A constructor sort can itself return a constructor sort. The sort of a
declaration's value is independent of how many parameters that declaration has:

```text
Alias : Type -> Type = Page
Pick<F: Type -> Type> : Type -> Type = F
Pick<Page><text>
```

`Alias` has no parameters and denotes a constructor. Each application supplies
one complete, nonempty parameter list. Partial applications and type-level
lambdas are outside this language revision.

`top(P -> P)` contains every constructor of that sort. The stronger contract
`top(P) -> top(P)` promises admission and a proper result for every proper input.
The two expressions must not be conflated.

## The canonical vocabulary

Revision 4 keeps the existing module envelope: revision, module name, imports
and declarations. Imports keep their ordinary module identity. There is no new
schema-digest field, schema registry or independently managed formation artifact.
The language revision fixes primitive forms, metadata and type-child positions.

One declarations map contains both kind definitions and type definitions:

```text
KindDefinition = parameters with sorts + result sort + kind expression
TypeDefinition = parameters with kind contracts + result kind + type expression
```

The kind vocabulary is general:

| Form | Meaning |
| --- | --- |
| Top / bottom at a sort | All / no values of that type-level sort |
| Kind parameter | An ordered, lexically bound kind-definition parameter |
| Local / imported reference | An ordinary kind declaration |
| Kind application | Exactly saturated application of a named kind definition |
| Union / intersection | Alternatives / simultaneous requirements at one sort |
| Node pattern | A primitive form, metadata selectors and a child-sequence pattern |
| Recursive group | A guarded simultaneous greatest fixed point of proper kinds |
| Recursive variable | A lexical group depth and slot |
| Arrow | A constructor's accepted argument kinds and guaranteed result kind |

Child-sequence patterns provide empty sequence, one item, concatenation,
alternatives and finite repetition. A node pattern matches the entire ordered
sequence of type children. Repetition admits arbitrarily long finite sequences,
never infinitely many children.

The scalar form's selector key is `"selector"` (`int`, `utf8-text`, `bool`);
the record form's key is `"mode"` (`open`, `closed`). Other primitive metadata remains
part of the admitted descriptor, but this kind language supplies no arbitrary
predicates over it. There are no callback predicates, subtree-equality tests or
unrestricted refinements hidden behind node selectors.

Named kind references form an acyclic dependency graph. Explicit recursive
groups provide recursion; every recursive cycle must cross a node's type-child
descent. Boolean operators and aliases do not guard recursion. Parameters remain
fixed inside a group, excluding endlessly growing kind applications.

## Ordinary families and Data

The bundled `horos.kinds` module supplies normal definitions. Schematically:

```text
List<K>       = node(list-of, {}, item(K))
Tuple<K>      = node(tuple, {}, star(item(K)))
Record<K>     = node(record, {}, star(item(K)))
Map<K,V>      = node(map-of, {}, concat(item(K), item(V)))
Wire<K>       = node(wire, {}, item(K))
UnionOf<K>    = node(union, {}, star(item(K)))
```

`Scalar` explicitly enumerates scalar selectors and the supported text
refinement forms. Atom refinements have their own family. `Exact` has no observed
type children: a parameter-looking tuple inside its payload remains literal
data. An exact text singleton does not automatically become a scalar descriptor.

`Data` is a recursive union using those ordinary families:

```text
Data = nu D.
    Any | AtomFamily | Scalar | Exact
  | Tuple<D> | Compound<D> | Seq<D>
  | List<D> | Set<D> | Map<D,D>
  | ClosedRecord<D> | OpenRecord<D>
  | UnionOf<D> | AndOf<D> | Not<D>
```

All listed alternatives expand into ordinary kind expressions. There is no
hidden data-leaf flag. Adding `Rows`, importing an equivalent definition under a
different name, or defining another recursive kind uses the same machinery.

Authors explicitly import the prelude identity to refer to its names. Horos
also uses that bundled module for primitive operand contracts. These are normal
kind contracts, bound to their declaration identities; a user's unrelated kind
called `Data` cannot change primitive formation rules.

Ground `rec` and `ref` become lexical backedges in the observation graph. They
are not extra alternatives in `Data`. The graph is finite, but its unfolding can
be infinite, which is why recursive kinds use greatest fixed points. Ground
values remain finite, and ground Horos retains its own runtime-value descent
rule. Kind guarding and ground recursion guarding solve different problems.

An open record contributes an implicit `any` child for additional fields, as
well as every declared field, including optional fields. Thus an open record
with only scalar named fields does not establish `Record<Scalar>`. An optional
Wire field prevents Data membership just as a required Wire field does.

## Generic bodies and constructor contracts

A constructor contract `A -> B` promises that arguments in kind `A` are admitted
and results belong to `B`. It makes no promise about other arguments. Its domain
is contravariant and result covariant when establishing a compatible contract.

Known constructors retain their checked bodies. Given a requested contract,
Horos first establishes that the requested domain is admitted by the original
parameter contract, then reasons about the body using fresh rigid inputs with
those assumptions. The original open declaration must already be valid; a
narrower query cannot repair it.

For pagination:

```text
Page<T> = closed record {
  items: list-of(T)
  next?: text
}
```

The same retained body can establish:

```text
Page : Data -> Data
Page : Scalar -> Rows
Page : K -> Record<List<K> | Text>
```

The last line describes a query for each particular suitable kind `K`, not a
universally quantified constructor contract. `Page : Flat -> Rows` does not
generally hold: substituting `list-of(text)` makes the items field a nested list.

Abstract constructors expose their declared guarantees. For `F : Flat -> Rows`
and `T : Flat`, `F<T>` has a Rows guarantee. Wrapping it in a record yields
`Record<Rows>`, which is different from Rows.

Intersections and unions of arrows remain distinct:

```text
(A -> B) & (A -> C)     applied to A guarantees B & C
(A -> B) | (A -> C)     applied to A guarantees B | C
```

Every possible callee alternative must admit the argument. Intersect promises
that apply jointly; union results across alternatives. Repeated occurrences of a
rigid parameter retain the same identity when assumptions split into cases.

The result model is a checked symbolic value plus assumptions. It is not a
principal inferred kind. For example, `Duplicate<T> = tuple(T,T)` preserves an
equality relationship between two unknown subtrees that regular kind patterns
cannot generally express. Keeping the body avoids discarding that relationship
without promising a complete relational solver.

## Membership, inclusion and limits

Closed membership operates on pairs of observation nodes and kind states. The
positive Boolean obligations are solved at their greatest fixed point, matching
each finite child list in full. Merely encountering a previously visited pair
does not prove membership.

Proper-kind inclusion uses positive recursive hedge automata. Horizontal finite
word automata describe child sequences; intersections require the same child to
satisfy all simultaneous obligations. First compute which states have a finite
or infinite regular continuation. Then permit an internal cut leaf at those live
states and decide inclusion of the resulting finite prefixes. Cuts are added
after intersections are composed, since incompatible kinds must not acquire a
common inhabitant merely because both separately allow an unknown prefix.

This finite-cut route follows the positive greatest-fixed-point setting of
[Schmidt-Schauß, Sabel and Schütz](https://www.numdam.org/item/10.1051/ita:2007010.pdf),
extended with finite child sequences and the fixed metadata partitions. It does
not require a general parity-automaton engine. The implementation must not claim
the cited paper's complexity bound automatically for its additional constructions.

Inclusion is over a free structural observation universe containing admitted
Horos observations. A failed inclusion need not identify a realizable Horos
descriptor. For example, `nu X. Wire<X>` has an infinite structural model, while
current ground-recursion admission cannot produce it. This distinction affects
the public outcome:

| Outcome | Meaning |
| --- | --- |
| Proven | The requested guarantee has been established |
| Refuted | A closed value fails membership, or an admitted counterexample establishes failure |
| Unproved | The available assumptions or qualified proof rules do not establish the claim |
| Invalid | The syntax, binding, sort, ownership or request is malformed |
| Exhausted | The operation reached a resource limit |

Constructor-contract reasoning is sound and deliberately incomplete. Resource
exhaustion never becomes Refuted. Required obligations must be established before
Horos publishes a checked program. Existing budgets bound construction,
automata, environments, comparisons and expanded output, with operation-local
caches and immutable Programs.

## Bittype and native generation

The [native constraint-preservation proposal](native-kind-constraints.md) explains
how source kind obligations could survive TypeScript, Go, Rust and C++ generation
through direct bounds or source witnesses. Its encodings remain proposed; the
implemented representation/enforcement distinction below is unchanged.

The shared Bittype semantic graph already supplies entities, parameters,
references, applications and subject-associated facts. Horos-specific terms
represent kind expressions, so the kernel requires no `Data`, `Scalar`, `List`
or `Record` cases.

The Horos provider publishes these separately:

- The semantic sort of a declaration or expression.
- Its declared kind contract.
- Its checked symbolic body.
- Individual established memberships, with their assumptions.
- Completion of the required obligations.

There may be several membership facts for one subject. Missing facts assert
nothing. A narrower query proving `Page : Scalar -> Rows` must not attach an
unconditional Rows fact to the body under its original unrestricted parameter.
Completion never means enumeration of every true membership.

Kind definitions are compile-time declarations, not native type exports. Native
generation retains the open pagination type, ordinary bindings and named
applications. It uses an explicit **representation** profile: generated types
represent uses already admitted by Horos. This is different from promising that
arbitrary native code can never violate the source kind constraint. An
**enforcement** profile needs a specified native encoding or returns Unsupported.

The current native subset remains closed records, text, bool, lists, optional
fields and ordinary proper generic parameters/applications. Residual constructor
parameters/results/aliases, unused proper parameters, canonical integers and
other unsupported descriptor forms remain located Unsupported throughout the
exported type closure. No hidden specialization or omission repairs them.

Ground descriptor extraction requires a proper closed value, membership in the
ordinary ground kind, materialization and the unchanged ground WF judgment.
This bridge does not hydrate endpoints, choose codecs or authorize behavior.

## Boundaries

Function, Stream and Cell stay ordinary library-owned declarations. Establishing
membership in a kind never selects one of their runtime implementations or
native representations. Adding those kinds or user-defined categories does not
change the generic generator, engine, Wire or runtime.

Adding an actual primitive descriptor form is different: it changes Horos's
versioned syntax, metadata, observed type children and operand contracts. It must
not be disguised as an ordinary library definition.

The replacement removes Claim/Interface, DataSummary, the dedicated Data
evaluator and revision-3 readers. It does not retain a compatibility path under
new names. Quantified constructor schemes, growing kind instantiations,
type-level lambdas, partial application and arbitrary subtree predicates need
separate concrete requirements before extending this model.
