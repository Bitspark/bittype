# Go native admission through compositional certificates

**Proposal and compiler qualification, 5 November 2026.** The owner requested
documentation and Go compiler tests for this encoding, including nested pagination,
aliases, inference and evidence forgery. No generator, native semantic provider,
Logos bridge or runtime binding is implemented. The existing enforcing generation
profile remains Unsupported.

The [symbolic pagination proof](logos-preservation.md) establishes a result under
source membership and selected representation assumptions. This document proposes
a native admission boundary that requires those assumptions before a Go consumer
instantiates the bounded alias. It is one candidate within the
[type-system architecture's native-admission boundary](type-system-architecture.md#native-admission-remains-independent).
It does not implement the proposed theory, evidence or bridge interfaces.

## Guarantee and scope

Fix the source world, selected mappings and a generated Go admission package.
For the supported consumer operations, the intended claim is:

```text
NativeAccepts_C(S,N,E)  // Go admits DataPage[S,N,E]
    implies there exists T:
        Denotes_C(S,T) and Member(T,Data) and Rep_C(T,N)
```

This specializes the [architecture's admission statement](type-system-architecture.md#native-admission-remains-independent)
to the resolved ordinary Data kind. S identifies the intended canonical subject
through the fixed Denotes_C interpretation; N alone need not identify it. E is a
certificate type whose issuance and constraints must justify that exact S/N
association under C. The interpretation is fixed before checking, not chosen to
make the conclusion true. This is a Go convention, not a kernel-wide tuple.

The Go compiler checks the relationships between certificate types. Horos and
the proposed Logos integration establish the source judgments and representation
laws that authorize their declarations. Go does not independently infer Horos
membership from native fields.

This is admission at type application. The bounded alias subsequently denotes
the ordinary payload type. It is not a nominal barrier against assigning an
otherwise compatible payload, a restriction on unconstrained Page, or validation
of arbitrary runtime values.

Consumers are outside the fixed issuing package. Editing that package changes
the proof inputs and requires justification again. This is a compiler/type-model
boundary, not a runtime security or cryptographic certificate scheme.

## Canonical declarations and payloads

The source assumes an explicit import of the ordinary kinds module as h:

```text
type Page<T> = { items: list<T>; next?: text; };
type DataPage<D: h.Data>: h.Data = Page<D>;
type Person = { name: text; active: bool; };
```

The selected Go payloads remain ordinary declarations:

```go
package model

type Person struct {
    Name   string
    Active bool
}

type Optional[T any] struct {
    Present bool
    Value   T
}

type Page[T any] struct {
    Items []T
    Next  Optional[string]
}
```

Page is unconstrained in both languages. The fixture uses the same explicit
presence representation as the [structural profile](semantic-generation.md#pagination-subset).
The spelling Optional is local to this independently authored fixture.

## Encoding the admission relation

The following declarations belong to one generated admission package. It imports
the selected payload packages; payload packages do not import it.

```go
package admission

import "example/model"

type PersonSource struct{}
type PageSource[S any] struct{}

type DataEvidence[S, N any] interface {
    ~func(S) N
    dataEvidence()
}

type PersonEvidence func(PersonSource) model.Person

func (PersonEvidence) dataEvidence() {}

type PageEvidence[S, N any, E DataEvidence[S, N]] func(PageSource[S]) model.Page[N]

func (PageEvidence[S, N, E]) dataEvidence() {}

type DataPage[S, N any, E DataEvidence[S, N]] = model.Page[N]
```

S identifies the source subject, N its selected native payload, and E the evidence
for their membership/representation association. The source identifiers are
distinct defined types, not aliases of one common empty struct. PageSource[S]
denotes the actual source application Page<S>; it does not itself grant membership.

The function signature indexes the exact pair. It denotes no conversion,
serialization or behavior. No evidence value must be created or passed in the
generated API. Empty marker methods exist only to make the certificate types
satisfy the private method requirement; they need never be called. This encoding
therefore requires marker-method emission beyond strictly type/interface-only
output. The [architecture records that scope extension](type-system-architecture.md#go-marker-method-scope-extension)
and its ownership: the selected Go convention supplies the empty declarations
through native AST/CST tooling, with no runtime calls, adapters or validation
behavior. The tests author those methods directly; documenting this proposal
does not implement production emission, which remains deferred.

DataEvidence is a nonbasic constraint interface, not a runtime proof object.
An arbitrary function value cannot supply its private method. Creating a nil or
other value of an already authorized certificate type grants no new source/native
association.

## Issuance and generic composition

PersonEvidence may be emitted only after establishing both Person : Data and
Rep_C(Person, model.Person). Its compiler fixture takes those judgments as explicit
foundational assumptions; compilation does not establish them.

PageEvidence encodes the independently justified conditional theorem:

```text
E supports T : Data and Rep_C(T,N)
    implies PageEvidence<S,N,E> supports
            Page<T> : Data and Rep_C(Page<T>,model.Page<N>)
```

The generic constraint retains the premise. There is no unconditional generic
Evidence<S,N> that issues certificates for arbitrary pairs, public or private.
The selected Page theorem must concern the actual source/native declarations,
optional/list mappings and binder correspondence.

PageEvidence is a proposed generated consequence of the ordinary structural
proof, not a hand-authored Page branch in Bittype or Logos. Data's Go integration
owns this evidence convention. Other library-owned kinds supply their own
judgments and justified conventions through the same proposed composition
boundaries. The agnostic kernel interprets none of their names.

Adding a source type can add its normal payload declaration and a certificate
justified by source analysis and the selected mapping. That is not an enumeration
of the generic domain. The Page family needs no new case, union member or nesting
limit. Recursive source types can use established Horos judgments as leaf
premises; the Go type system does not replace Horos's fixed-point reasoning.

## Pagination, nesting and caller scope

An independent consumer can instantiate the bounded alias:

```go
type People = admission.DataPage[
    admission.PersonSource, model.Person, admission.PersonEvidence,
]

type PeopleEvidence = admission.PageEvidence[
    admission.PersonSource, model.Person, admission.PersonEvidence,
]

type PeoplePages = admission.DataPage[
    admission.PageSource[admission.PersonSource], People, PeopleEvidence,
]
```

People denotes model.Page[model.Person]. PeoplePages denotes
model.Page[model.Page[model.Person]]. The outer application receives evidence for
the inner result, not the original Person evidence.

The source argument can remain open:

```go
type Nested[S, N any, E DataEvidence[S, N]] = DataPage[
    PageSource[S], model.Page[N], PageEvidence[S, N, E],
]
```

The caller's S/N/E remain the inputs. Page's formal environment is not captured.
The same constructor can be applied again for a third page or any finite depth.
There is no finite union of admitted Person/string/bool arguments.

The low-level Go convention expands one source parameter into three linked
native parameters. Generated named applications can hide them. This is an
explicit representation convention, not a requirement for every language or
every future encoding.

## Why both the signature and private method matter

The intended native-model argument uses Go's
[type sets](https://go.dev/ref/spec#General_interfaces),
[method sets](https://go.dev/ref/spec#Method_sets),
[private identifier identity](https://go.dev/ref/spec#Uniqueness_of_identifiers)
and [defined types](https://go.dev/ref/spec#Type_definitions):

- A consumer's identically spelled private method is a different identifier.
- A newly defined function type does not inherit the issued type's methods.
- Embedding issued evidence can promote its method into a struct, but that
  struct cannot have underlying type func(S) N.
- A pointer to issued evidence also fails the function-shape constraint.
- An alias retains the original type identity and indexed signature; it does
  not rebind the certificate's source or payload.
- Only actual issued leaves and justified constrained constructors in the fixed
  admission package may provide the private method on function-shaped types.

A model of those declarations should establish that every accepted concrete
certificate originates in an authorized leaf or composition. The Page constructor
then reduces its result claim to its child certificate, whose claim is a premise.
Open generic forwarding retains that premise; it does not establish a concrete
certificate for an arbitrary subject.

This differs from the withdrawn private singleton proposal: hiding a universally
available private generic instantiation did not restrict its issuance. The present
scheme must have no such unrestricted constructor. Inference is permitted to
discover a valid certificate; the tests check that it cannot recover an
unauthorized one merely from an alias's erased result.

The formal admission lemma remains to be proved in the proposed Go model.
Compiler examples exercise these rules; they are not an exhaustive proof over
all possible consumer programs.

## Compiler qualification

The [compiler test](../integration/go/go_native_admission_test.go) copies the
[independent fixtures](../integration/go/testdata/go-admission/) into a temporary
module. It compiles model, admission and foreign-issuer packages first, then
builds each consumer separately with the actual Go compiler. Expected rejection
must match a consumer diagnostic; a broken dependency or unrelated tool failure
cannot count as a successful rejection.

The fixture declares Go language version 1.25.0. The maintained Go contracts
workflow selects compiler Go 1.27.3. The test records its actual compiler version,
disables toolchain switching and dependency fetching, and runs in the existing
suite without an opt-in flag:

```sh
go test -count=1 -v -run '^TestGoNativeAdmission$' ./integration/go
```

Remote execution on Go 1.27.3 linux/amd64 passed all 29 admission
cases: 8 positive consumers compiled and all 21 negative consumers produced the
expected rejection. The fixture language was Go 1.25.0; this is not a run of the
Go 1.25 compiler. These observations qualify the tested compiler cases, not the
universal admission theorem or the source judgments assumed by the issuer.

| Cases | Intended outcome |
| --- | --- |
| Concrete pagination and raw-payload identity | Accept |
| Two/three nested pages and open Nested | Accept |
| Aliases and cross-package reexports preserving the exact identities | Accept |
| Generic forwarding, argument inference and function-assignment inference | Accept |
| Values of already authorized evidence types | Accept without adding any new judgment |
| Unconstrained Page of a capability, function or channel | Accept |
| Wrong source, payload, nested payload or inner-result evidence | Reject |
| Unsealed functions, counterfeit private methods and newly defined copies | Reject |
| Struct embedding, embedded constraint interfaces and pointer evidence | Reject |
| Unauthorized PageEvidence and misleading generic aliases | Reject |
| Alias-result erasure, argument/assignment inference without valid evidence | Reject |
| Defined source/payload copies and a foreign issuer's same-spelled method | Reject |
| Recursive alias certificates and methods added to a nonlocal type | Reject |
| An invalid argument whose constrained alias ignores every parameter | Reject |

The fixtures contain helper functions to exercise compiler inference, but they
are built rather than run. No application behavior or runtime proof validation is
implemented. No generated provider-matrix cell changes.

## Composition and remaining boundaries

One consuming build selects an issuing package and its source/mapping context.
Library integrations contribute the laws and declarations needed for its
certificates. A new native mapping requires its own valid association; evidence
for a prior payload cannot be reused for a different native type.

An arbitrary downstream package cannot self-issue new evidence. Importing its
payloads and establishing their certificates in the selected admission package
is supported by this proposal. Federating independent issuing packages is not
solved by copying a marker spelling or embedding an interface; it remains a
separate design obligation. The compiler fixture explicitly rejects the foreign
issuer and accepts a true reexport of the original certificate.

The following remain separate from successful compiler qualification:

- Horos membership and representation proofs authorizing every leaf.
- Logos assumption discharge and the generic Page theorem authorizing composition.
- A formal account of the Go admission relation, inference and permitted operations.
- Representation completeness for the independently promised source domain.
- Source/result association, native emission and compatibility with selected aliases.
- Runtime value conformance, protocol behavior and adapters.

The production generator, formatter and runtime are unchanged. The production
enforcement profile remains Unsupported until the missing semantic/proof and
generation work is implemented and qualified. There is no runtime registry,
new core kind, compatibility alias or handwritten membership catalogue.
