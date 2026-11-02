# Agnostic semantic model

**2 November 2026. Status: documented design proposal; no shared semantic model
package or semantic-analysis provider SPI is implemented.** The owner requested
this document after reviewing the minimal Go sketch below. The initial scope is
types and declarations. This change adds documentation only.

Bittype should own a shared vocabulary for semantic facts and the contracts for
producing and inspecting them. Language providers own the rules establishing
those facts. A common Analyze method returning an opaque compiler object would
not give consumers a useful agnostic model.

This model complements the [AST](ast.md). Syntax records what was written;
analysis establishes what declarations it refers to, how parameters are bound,
what type expressions denote and which conditions have been checked. Horos
remains the [canonical data-type authority](model.md). A semantic model of Go
source still describes Go meanings; conversion to Horos needs an explicit mapping.

## How the model was derived

The local evidence is bittype at
`ae2466347a5067f27bd74cd1d48e89f3b8622d03` and horos at
`097496351b70336462c659e3deabea833e294b85`, inspected on 2 November:

- Our [generic AST](../trees/ast/go/node.go) records syntax, reference edges and
  origins. It does not record analysis coverage or why a binding is valid.
- [Horos's checked Program](https://github.com/Bitspark/horos/blob/097496351b70336462c659e3deabea833e294b85/core/go/declarations/model.go)
  retains module, declaration and parameter identities and resolved applications.
  These are concrete requirements; its descriptor vocabulary stays Horos-owned.
- Go's [types.Info](https://pkg.go.dev/go/types#Info) exposes definitions, uses,
  types, scopes and generic instances. The
  [TypeScript checker API](https://github.com/microsoft/TypeScript/wiki/Using-the-Compiler-API#using-the-type-checker)
  separates syntax nodes, symbols and types. These are design references, not
  new dependencies or a commitment to a compiler API. The cited TypeScript
  examples describe its pre-7 compiler API.
- Our [Rust AST](../syntax/rust/go/nodes.go) already distinguishes lifetime and
  type parameters. The [Rust reference](https://doc.rust-lang.org/reference/items/generics.html)
  also describes constant parameters. A shared model must not assume every
  argument is a type; constant-parameter support remains outside our current
  Rust syntax slice and this documentation change.

The common abstraction is a graph of semantic facts. It shares representation
and inspection without prescribing one language's lookup or type rules.

## Minimal Go sketch

These illustrative declarations are not a present public package or finalized
storage API. Exported fields make the shape readable. An implementation would
publish checked, immutable snapshots with no mutable maps, slices or handles
leaking into the stored result.

```go
package semantics

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type EntityID uint32
type ScopeID uint32
type TermID uint32

// A vocabulary item owned by a provider.
// Examples: a parameter category or a semantic relation.
type Key struct {
	Owner string
	Name  string
}

type Model struct {
	Language language.ID
	Profile  string
	Snapshot string // Identifies the exact inputs, dependencies and options.

	Entities map[EntityID]Entity
	Scopes   map[ScopeID]Scope
	Terms    map[TermID]Term

	// Established reference bindings. Absence means no recorded answer.
	Bindings map[ast.NodeRef][]EntityID

	Judgments []Judgment
	Completed []Key // Analyses completed under this profile.
}

type Entity struct {
	Name         string // Display name; identity is EntityID.
	Declarations []ast.NodeRef
	DeclaredIn   ScopeID

	Meaning    TermID
	Parameters []Parameter
	Members    []EntityID
}

type Scope struct {
	Declarations []EntityID
	Links        []ScopeLink
}

type ScopeLink struct {
	Relation Key // Provider-defined lookup relationship.
	Target   ScopeID
}

type Parameter struct {
	Entity EntityID
	Sort   Key // For example, a type or lifetime parameter.
}

// A reference to an entity, optionally applied to arguments.
// Head may identify a named declaration, parameter or provider-owned intrinsic.
type Term struct {
	Head      EntityID
	Arguments []TermID

	// Established generic binding: parameter identity -> argument meaning.
	Substitution map[EntityID]TermID
}

type Judgment struct {
	Relation Key
	Operands []TermID
	At       ast.NodeRef
	Verdict  Verdict
	Message  string
}

type Verdict uint8

const (
	Undetermined Verdict = iota
	Holds
	Rejected
	Conditional
	Unsupported
	Incomplete
)
```

## Meaning and invariants

| Concept | Shared meaning | Provider responsibility |
| --- | --- | --- |
| Entity | Identity of a declaration, member, parameter or intrinsic, with declaration sites | Which declarations introduce, merge or alias entities |
| Scope | Declarations and connections to other environments | Lookup, namespaces, visibility and import rules |
| Binding | An established association from a syntax occurrence to target entities | Resolving the occurrence and interpreting valid target sets |
| Term | Meaning formed by referencing or applying entities | Constructor meanings, argument categories and admissible forms |
| Substitution | Parameter identity mapped to argument meaning | Establishing the binding and checking its constraints |
| Judgment | A result for a named relation in this context | Relation rules, premises and supported coverage |

IDs are local to one snapshot. Zero denotes an absent optional reference, never
a fabricated entity, scope or term. Required links must identify records in the
model, including its admitted dependency closure. Equal IDs in different models
establish nothing. Cross-model use needs explicit identity mapping and compatible
context; names and casts cannot establish it.

Names are display metadata. Multiple declaration sites may contribute to one
entity; an intrinsic may have no source declaration. A named alias, its target
and the term it denotes remain distinguishable. In this type-declaration slice,
Meaning is the denoted or defined term. Future value declarations would need
separate accounts of their value meaning and their type.

AST addresses provide provenance, not semantic identity, module pins or ownership
proofs. Analysis does not mutate its input ASTs. Supplied AST reference annotations
must agree with analysis; their presence alone does not establish a binding.

Scope links form a graph. The carrier imposes no universal lexical-parent lookup,
import precedence or acyclicity rule. Horos's revision-1 cycle restrictions remain
Horos rules; native recursion must not inherit those restrictions accidentally.

Bindings contains established targets, not unresolved candidate sets. Multiple
targets require an explicit meaning under the provider's profile. Absence means
no recorded answer, not a missing-name error or successful empty resolution.
Ambiguity, unavailable dependencies and limits need explicit diagnostics; a richer
resolution-result representation is deferred.

An open parameter T is a valid bound parameter, not an unresolved name. Parameter
sorts are provider-owned. Arguments preserve order; Substitution records their
established binding by parameter identity. The engine does not infer bindings
from spelling or execute expansion merely because a substitution is recorded.

Head explicitly identifies the declaration, parameter or intrinsic involved.
The core can inspect operands without defining the meanings of a Go pointer,
Rust reference or Horos descriptor. Similar structure establishes no cross-language
equivalence. This minimal Term covers references and ordinary applications, not
every native language's semantic forms.

Keys name vocabulary under a provider/profile contract. They cause no discovery,
imports or automatic registration. Shared code can check graph integrity and
identity consistency; labels alone cannot establish language truth.

Profile states admitted forms, relations and analyses. Snapshot identifies the
exact AST inputs, dependencies, options and relevant environment. This sketch
defines no content hash, persistent identity or cache-key algorithm. Those need
explicit contracts before caching or interchange.

## Pagination

The existing two-module example can be displayed as:

```text
pagination.Page<T> = record { items: list-of(T), next?: text }
people.Person = record { name: text }
people.PeoplePage = pagination.Page<people.Person>
```

This is explanatory notation, not a new Horos text grammar. Retain distinct
entities for Page, Page's T, Person and PeoplePage. Person's term refers to the
caller's Person entity. The application is:

```go
model.Terms[appliedPage] = Term{
	Head:      pageEntity,
	Arguments: []TermID{personTerm},
	Substitution: map[EntityID]TermID{
		pageParameterT: personTerm,
	},
}
```

These variables stand for allocated IDs in an initialized model; this is an
excerpt, not a standalone builder. PeoplePage's Meaning refers to appliedPage.
The T inside Page's items refers to pageParameterT. A different declaration's T
has a different identity and is unaffected.

The application remains inspectable without expansion. A Horos provider can
use its existing specialization on explicit request to obtain a ground descriptor.
That preserves the open declaration and does not select a native representation
or perform a Horos-to-native transformation.

## Judgments and coverage

Identity, assignability, conversion and constraint satisfaction are distinct
relations, not one Compatible predicate. Go's
[type relation APIs](https://pkg.go.dev/go/types#Identical) illustrate this
distinction and the need for a consistent symbol environment.

Holds and Rejected are completed answers to the named relation. Conditional
depends on explicit assumptions or parameters. Unsupported means the provider
cannot answer that request. Incomplete means necessary context or resources
were unavailable. Undetermined is the safe default. None of the latter states
may silently become Holds or Rejected.

Completed records analyses actually completed, not advertised capabilities.
Finishing resolution does not establish every type constraint. Finishing an
analysis may still produce Rejected judgments. Consumers must require the
coverage and successful judgments their operation needs. Profiles must specify
the meaning and granularity of coverage keys before implementations exchange them.

Judgment has no structured premise or diagnostic-cause graph in this minimal
sketch. Message is explanatory text, not machine-readable proof. Conditional
results must not drive automatic generation until their assumptions are represented
explicitly. Partial facts from a failed attempt must never masquerade as a complete
valid model. The future SPI must specify errors, limits and partial-result handling.

## Current integration and future analysis SPI

Today, [resolve/horos/go](../resolve/horos/go/resolve.go) calls Horos through a
private bridge. The [Horos checker](../check/horos/go/check.go) also reconstructs
module values, resolves them and verifies AST bindings, returning success or error.
Native checkers validate admitted syntax and selected contextual rules. There is
no common semantic-analysis provider exposed by the engine.

The proposed operation is conceptually:

```text
Analyze(AST collection, explicit context)
    -> semantic model, diagnostics and coverage
```

This is not a finalized Go interface. Analysis must accept syntax with unresolved
names; requiring a fully analyzed AST first would be circular. A later change
should separate syntax admission from analysis and adapt Horos's existing semantic
implementation behind the common contract. This document does not change the
current Horos AST profile, Check behavior or registration.

Dependency retrieval remains explicit application work. Providers retain language
rules and model/compiler dependencies; generic infrastructure owns only the shared
representation and SPI. Selection and registration remain caller-controlled.
Future packages must follow the [stage-first boundaries](../LAYOUT.md). No empty
packages, provider interfaces or new dependencies are added by this proposal.

## Scope and next qualification

Start implementation with entities, scopes, established bindings, ordinary generic
terms, judgments and explicit coverage. Qualify those concepts using Horos and at
least one native language. Go is a useful candidate because its checker exposes
these facts; this does not claim native analysis is already implemented in bittype.

Future tests should cover imported pagination, same-spelled declarations and
parameters, aliases, multiple declaration sites, open/applied generics, missing
or ambiguous names, context changes, limits and immutable snapshots. Unknown or
unsupported must remain distinct from false; completed analysis from validation
success. These are proposed acceptance cases, not tests added by this document.

Richer constant terms, dependent projections, structured premises, general
inference and further native semantics need explicit extensions. Distinguishing
an unapplied generic reference from a checked zero-argument application also needs
a precise representation contract before this sketch becomes a public API.

Function, Stream, Cell and future contracts remain ordinary library entities and
applications. Forbidden are core contract kinds, semantic name switches, builtin
imports/registrations, shape fingerprints choosing behavioral ownership, and
automatic equivalence between languages. Moving such dispatch into a semantic
model would violate the same ownership rule as putting it in the runtime.
