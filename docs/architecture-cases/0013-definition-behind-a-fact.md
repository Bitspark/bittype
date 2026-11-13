# bittype/ACX-013: Two analyses published a declaration's definition behind a fact

**Status.** Identified on 13 November 2026 at main
`3998db5`, by the two other generation pairs of
[#196](https://github.com/Bitspark/bittype/issues/196) and
[decision O45](https://github.com/Bitspark/bittype/issues/196#issuecomment-6854327944).
Corrected for Go and for TypeScript in the change of
[#310](https://github.com/Bitspark/bittype/issues/310), with the two
[decided points](#decided-points) below. **Pending:** [Rust](#rust) has the same
shape and is corrected next under the same issue, in its own pull request.
The case stays pending until then. The correction was not read by an
independent reviewer when this was written.

**Scope.** The graphs that the Go, TypeScript and Rust semantic analyses
publish, and the readers of the three analysed vocabularies. The generation
contract and its kernel are not part of the defect and did not change.

## Intended boundary

The shared semantic graph has one place for a declaration's definition, the
body of its entity. [The semantic model](../semantics.md#representation) says
"`Owner` and `Body` are optional links, not a universal namespace or type-of
rule", and reads the body as the definition in its example: "PeoplePage's body
is an application with a Reference to Page as its callee".

The generation kernel walks that structure for every source.
[GEN-01](../generation-gaps.md#gen-01), completion criterion 1: "The kernel
consumes declaration, reference and application structure through the shared
graph. No source supplies a traversal". It "opens declarations from their body
and real binders", and it "never reads entity kinds, names, operators,
payloads, locator meanings or facts"
([generation SPI](../generation-spi.md#semantic-traversal-and-source-interpretation)).

What varies with a language is the vocabulary: its entity kinds, operators,
payloads and facts. Where the definition of a declaration sits does not vary.

## Actual implementation and consequence

Two analyses used the body as the type of the declared name and put the
definition behind a fact of their own vocabulary.

- **Go.** `internal/golangsem/export.go` set `body := e.typ(o.Type())` for
  every object. For a defined type that is a Reference to the entity itself.
  The struct was the value of the fact `underlying` on that Reference term.
- **TypeScript.** `internal/typescriptsem/compiler.cjs` set
  `graph.Entities[id - 1].Body = type`, where `type` is the declared type of
  the symbol. For a named declaration that is a Reference to its own entity.
  The object was the value of the facts `properties` and `structure`. A
  comment stated the choice: "References break recursive type graphs;
  definition structure is a fact, not a containment edge back into the named
  type".
- **Go's readers** read the body as the type of the name:
  `semantics/golang/go/model.go` in `lookup` and in `Type.Parameters`.

To the kernel every Go or TypeScript declaration was therefore a use of
itself. The first other pair, `golang-to-typescript`, was refused in every
test with an evidence cycle before any field became a use. The pairs lane then
made the smallest change that let it continue: an optional interface of a
source interpretation, `Definitions[F]`, through which a source names the
term at which the kernel starts a declaration's body (commit `3452622` on its
branch). O45 refused that commit: a source that names the starting term
supplies the first step of a traversal from one of its facts, which GEN-01
excludes, and the kernel could check only that the term is in the graph.

This is no missing coverage. Nothing was absent from the two graphs. The
definition was published, in the wrong place, and the place the shared
mechanism reads held something else.

## Rules involved

- [AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02):
  the meaning of the body belongs to the shared graph, whose mechanisms depend
  on it. Two contributions gave that link a meaning of their own, the type of
  the name, and kept the shared meaning in their private vocabulary.
- [AC-03](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-03):
  adding a materially different source did not compose additively. It needed
  either an edit of the shared contract (`Definitions[F]`) or this change of
  two shared analyses.
- [AC-07](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-07):
  GEN-01 was qualified with the horos source and the independent fixture. Both
  publish the definition as the body. That evidence did not establish that the
  Go and TypeScript graphs can be walked, and the migration's assumption "a
  source has an analysis that publishes a semantic graph" was true and not
  sufficient.
- [AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
  governs the correction: the misplaced definition moves first, as its own
  change, before the pairs that need it. The cases that the change cannot
  decide are refused and not guessed.
- [AC-11](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-11)
  is preserved by the correction and was not violated by the defect: the
  reader of "the one fact of a subject", which the Go vocabulary now needs as
  the TypeScript vocabulary did, has one agnostic owner.

## What worked and what it did not establish

Every consumer of the two analysed vocabularies worked: verification, the
TypeScript Files formatter's semantic checks, the Go import environment and
the native queries. They start from a declaration or a type and read facts.
That established the vocabularies for such consumers. It did not establish
that a consumer which is handed a graph and no vocabulary, the kernel, can
read a declaration.

## Correction

**Go.** The body of a captured type declaration or alias is the term of the
type that its written right-hand side denotes. The bridge already kept the
type expression of every captured declaration. A type parameter, a predeclared
type, a supplied declaration and a field have no body. The type of a type name
and of a field is the new entity fact `declared-type`, which
`Model.DeclaredType`, `Lookup`, `LookupIn` and `Type.Parameters` read.
[Go bodies and declared types](../native-semantics.md#go-bodies-and-declared-types)
states every kind.

**TypeScript.** The body of an alias is the term of its written right-hand
side: the Reference or Application of the source expression, or the structural
expression. The body of an interface without an `extends` clause is the
`object` of its members. The body of an interface with one is an `extension`,
whose operands are the base type terms as written and then the `object` of
its own members. A type parameter and a property have no body. `declared-type`
was an entity fact before, and a property carries it now. The vocabulary
gained the reader `Model.Extension`, and no existing reader changed.
[TypeScript bodies and declared types](../native-semantics.md#typescript-bodies-and-declared-types)
states every kind.

**Recursion.** A field that names its own declaration is a Reference to the
entity, which is no containment. Term containment stays a DAG, and graph
admission checks it as before.

**Versions.** Each profile has its second version, because a published member
changed: `go1.25/captured-types-2` and `typescript-5.9.3/captured-types-2`. A
receipt of the earlier version is refused.

**The requirement on any source** is stated in the
[generation SPI](../generation-spi.md#semantic-traversal-and-source-interpretation)
and in one sentence of [the semantic model](../semantics.md#representation).

**One reader for the one fact of a subject.** `semantics.Semantics[L].OnlyFact`
is in the agnostic model. The Go and the TypeScript readers use it.

## Qualification

- [One reader in a test file](../../integration/go/definition_body_test.go)
  walks a graph as the kernel does: body, operands, reference targets and
  direct parameter bindings, and no fact, kind, operator or payload. It is one
  function for both languages. For each language it reaches every field of a
  struct or object type, a reference from one declaration to another, a
  recursive declaration, a use of a generic declaration's parameter and an
  application paired with that parameter. All ten subtests were red on main
  before the change.
- [Tests of each declaration kind](../../integration/go/definition_body_kinds_test.go)
  state the body and the declared type per kind, for both languages, a field
  and a property among them, and that a supplied Go declaration has no body.
  They state for every entity of the graph whether it has a body and whether
  it carries the fact.
- For TypeScript the walk-like reader also reaches, from an interface with
  two base types, each base declaration, the actual of the base with type
  arguments and each member that the interface declares itself.
- [One check for both languages](../../integration/go/earlier_profile_test.go)
  refuses a receipt of the earlier profile.
- Every test that passed on main passes by the same name.

Limits of this evidence:

- **The kernel itself did not run** over a Go or a TypeScript graph here. The
  reader in the test follows the walk's steps and is not the walk. The run of
  the two pairs without `Definitions[F]` belongs to the pairs lane of #196.
- **The kernel cannot check that a body is the definition.** It is the
  analysis's statement. The tests check it for the forms they name.
- **The Rust tests did not run** on the workstation of this change. The
  hosted check runs them. Of a body they assert one thing: the entity of a
  primitive type has none.

## Decided points

The first form of this correction left two points open, each with its options.
Both were decided on 13 November 2026
([decisions D1 and D2](https://github.com/Bitspark/bittype/issues/310#issuecomment-6857138520))
and are implemented in the same change.

- **Decided, D1: a field and a property have no body.** A field has a name
  and a type and no definition, so a body that holds its type was the type-of
  reading for one more entity kind. Its type is the `declared-type` fact of
  its entity. Nothing is lost for a reader that walks: the type is an operand
  of the struct or object term. `Body` now has one meaning in a Go and in a
  TypeScript graph.
- **Decided, D2: an interface with base types has an `extension` as its
  body.** Its operands are the base type terms as written and then the
  `object` of its own members, and its data states the number of bases.
  [The form](../native-semantics.md#the-extension-of-an-interface). It is not
  the `object` of all members, which is the checked type and hides the base
  declaration from the walk, not the own members alone, which would drop the
  inherited ones without a refusal, and not an `intersection`, which treats
  conflicting members differently.
- **Still outside the profile:** interface declaration merging. The helper
  refuses two declarations of one interface as unsupported coverage and
  publishes no graph, so the body of a merged interface is not defined.

## Rust

Read on 13 November 2026 in `internal/rustsem/helper/src/main.rs` and
`semantics/rust/go`, and not run.

- The link from a declared name to its own type is an entity fact,
  `declared-type`, and no reader of `semantics/rust/go` reads a body. That
  requirement is met.
- The body of a struct and of an enum is the declared type: a Reference to
  the entity itself, or an Application of the entity to its own parameters.
  The definition is the value of the term facts `fields` and `variants`. This
  is the same defect.
- The body of an alias is the type that the right-hand side denotes after the
  compiler's normalization. An alias of an alias has the normalized type, not
  a Reference to the alias it names.
- A field's body is its type, which Go and TypeScript no longer publish
  (decision D1).

**Open violation, corrected next.** Not corrected in the change for Go and
TypeScript. It is corrected next under
[#310](https://github.com/Bitspark/bittype/issues/310), in its own pull
request, with a new version of the Rust profile
([decision D3 of 13 November 2026](https://github.com/Bitspark/bittype/issues/310#issuecomment-6857138520)):
it is the same small correction a third time, and
[AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)
permits deferring a move only when its destination needs a huge effort. This
case stays pending until that change. Until then a Rust struct or enum is, to
the kernel, a use of itself.

## Assessment lesson

A shared slot was filled by every provider, validated by graph admission and
read by tests, and still meant two different things. Admission checks that a
body is a term of the graph. It cannot check what the term means. The walk was
qualified against the two sources that were written together with it. When a
shared mechanism reads a link of a shared graph, state what the link means for
every publisher where a publisher reads it, and qualify the mechanism against
a publisher that was written for another consumer.
