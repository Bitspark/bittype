# Horos kinds and data guarantees

**4 November 2026. Revision-3 source contract.** Horos owns kind checking,
constructor interfaces, data guarantees and specialization. Bittype projects
that model into [readable syntax](horos-syntax.md), [ASTs](horos.md) and
[provider-owned facts](semantics.md). The [consultation delivery table](../research-docs/0004-kind-system.md#applied)
records accepted work, validation and explicit deferrals. Publication and
exact-revision qualification are recorded in [dependency evidence](dependencies.md).

## Types, constructors and Data

A value inhabits a type. A proper type describes values, while a type constructor
accepts type-level arguments and produces a proper type. Constructors may accept
other constructors. Every constructor application supplies all its arguments;
these are finite, n-ary interfaces, not curried functions.

`Type` classifies proper types, including data-only and capability-bearing types.
`Data` is a refined kind, or subkind, of `Type`: the category of proper types with
an established compositional data-only certificate. It collects the admitted
data-only cases; it is not a new value-level union descriptor. A Data obligation
asks for that certificate. The checker does not solve whether every possible
inhabitant of an arbitrary descriptor happens to be data through emptiness or
other semantic equivalences.

The implementation exposes three views of this one model:

| View | What it establishes |
| --- | --- |
| Erased shape kind | Whether an expression is proper or a constructor with specified input shapes |
| Declared interface | Which arguments are admitted and which proper result is guaranteed: Type or Data |
| Checked Data summary | Whether the body establishes Data, possibly after substituting known arguments |

Erasing Data to Type recursively retains application shape and forgets the
refinement. That implementation view does not deny that Data is a kind.
`Type` is not the ground descriptor `any`: `any` still describes only data values.
All canonical descriptions are themselves ground Ontos values, even descriptions
of Wire. Encoding a description as data proves nothing about its instances.

## Page, Wrapped and Compose

The readable source uses explicit interfaces where needed. Ordinary parameters
and declaration results default to Type.

```text
module example.kinds;

type Page<T> = { items: list<T>; next?: text; };
type Person = { name: text; };
type DataPage<D: Data> = Page<D>;
type Ignore<T> = text;
type Send<T> = wire<T>;

type Wrapped<F: Type -> Type, T> = { value: F<T>; };
type Compose<F: Type -> Type, G: Type -> Type, T> = F<G<T>>;
type DataWrapped<F: Data -> Data, D: Data>: Data = { value: F<D>; };

type WrappedPage = Wrapped<Page, Person>;
type Nested = Compose<Page, Page, Person>;
```

Page is authored once. Its Type input admits either a data or a capability type;
its result is data exactly when its value-bearing input establishes Data. Ignore
has an unconditional data result, but every supplied argument must still be valid.
Send has a Wire result and does not establish Data.

| Application | Canonical judgment |
| --- | --- |
| `Page<Person>` | Valid, Data |
| `Page<wire<Person>>` | Valid Type, no Data certificate |
| `DataPage<wire<Person>>` | Rejected: the argument does not satisfy Data |
| `Wrapped<Page, Person>` | Valid, Data after substituting the known Page |
| `Wrapped<Ignore, wire<Person>>` | Valid, Data; Ignore does not use its admitted argument |
| `Wrapped<Send, Person>` | Valid Type, no Data certificate |
| `Compose<Page, Page, Person>` | Valid, Data through both known applications |
| `DataWrapped<Page, Person>` | Valid, Data; Page establishes the required preservation interface |
| `DataWrapped<Send, Person>` | Rejected: Send cannot guarantee a Data result |
| `DataWrapped<Ignore, wire<Person>>` | Rejected: the second argument must satisfy Data |
| `Wrapped<DataPage, Person>` | Rejected: DataPage cannot accept arbitrary Type inputs |

The known Page declaration can satisfy `Data -> Data` by checking its body under
a fresh symbolic Data input. This is an additional established interface of the
same declaration, not a copied Page or a different identity. Ignore can establish
`Type -> Data`. The original admission contract is always checked first; ignoring
an argument does not widen a Data-only input to Type.

An abstract constructor parameter exposes only its declared interface:

```text
type Bad<F: Type -> Type, D: Data> = DataPage<F<D>>;
type Good<F: Data -> Data, D: Data> = DataPage<F<D>>;
type BadInput<F: Data -> Data> = F<wire<Person>>;
```

Bad is invalid because F could be Send. Good is valid because F promises a Data
result for its admitted Data input. BadInput is invalid even if a proposed caller
would supply Page. A later specialization cannot repair an invalid open body.

Constructor interface comparison requires matching shape and arity, contravariant
input admission and covariant result guarantees. `Type -> Data` can satisfy
`Data -> Type`; `Data -> Data` cannot satisfy `Type -> Type`. These judgments
concern constructor substitutability, not native container variance or general
descriptor subtyping.

## Exact canonical revision 3

The authoritative grammar lives in
[Horos's declaration specification](https://github.com/Bitspark/horos/blob/main/docs/spec/declarations.md).
Here `A(s)` is a UTF-8 atom, `T(s)` a text embedding, `I(n)` an integer embedding,
`L(...)` a list embedding and `D(tag, ...)` the tuple whose first atom is
`horos/decl/` followed by tag. Maps have canonical text keys.

```text
Module = D("module", I(3), T(name), importsMap, declarationsMap)
Import = D("import", T(expectedModuleName), Atom(32-byte pin))

Claim = A("type") | A("data")
Interface = D("proper", Claim)
          | D("constructor", L(Interface, ...), Claim)
            // constructor input list is nonempty; result is proper
Parameter = D("parameter", T(name), Interface)
Declaration = D("type", L(Parameter, ...), Claim, E)

E = Horos descriptor with E only in designated child-type positions
  | Tuple(A("horos/wire"), E)
  | D("param", I(zeroBasedSlot))
  | D("local", T(name))
  | D("imported", T(alias), T(name))
  | D("apply", E, L(E, ...)) // nonempty, fully supplied arguments
```

The constructor's inputs can recursively be constructor interfaces. Its result,
and every declaration body, is proper. Store the declared interface once and
derive its erased shape; there is no second authoritative kind field. An omitted
source annotation elaborates to explicit `type`. A stronger inferred Data summary
does not rewrite that authored result claim.

The Wire form belongs to the unified static declaration language. Ground Horos
WF/Conforms and ground Ontos values remain data-only.

A declaration reference denotes a constructor when it has parameters and a proper
type otherwise. A zero-parameter declaration is referenced directly, not applied
to an empty argument list. Application has an expression callee, preserving its
syntax and binding. An application used as another application's callee fails
semantic kind checking in this revision because its result is proper.

Declaration identity remains the exact canonical module digest plus declaration
name; parameter identity adds its ordered slot. Import aliases retain the expected
module name and immutable pin. A supplied constructor keeps its definition scope.
All named references contribute dependency edges, including ignored constructor
arguments. The reader admits expression structure and descriptor skeletons,
then rejects declaration cycles before checking body interface and Data
obligations. Ground `rec`/`ref` remain separately checked recursion forms.

Only schema-designated expression positions are traversed. A literal tuple inside
`exact(...)` that resembles `param`, `apply`, `local` or `constructor` remains
literal data. Labels, presence markers and other metadata preserve their bytes
and ordering.

Revision 3 replaces revision 2 cleanly. Its reader rejects the old declaration
format. Reissue canonical modules and repin their dependents from rewritten
dependency leaves; readable ordinary Page source remains unchanged. Historical
releases, consultation input and validation evidence remain immutable. There is
no compatibility decoder, bound-only API or template-lifting path.

## Data summaries, specialization and limits

Horos retains symbolic summaries with True, False, Need, And and Probe forms.
Need refers to a proper parameter; Probe retains a checked constructor call and
its argument structure. For Wrapped the summary must retain F's application so
that substituting Page, Ignore or Send can establish different answers. A generic
summary is never overwritten with one specialization's result.

The public Data status distinguishes `proven`, `failed` and `conditional`.
Failed means the completed compositional criterion fails; it is not an assertion
that a non-data inhabitant exists. Conditional means a remaining
parameter or constructor dependence. Constructor expressions have no proper-type
Data summary or status. Limit exhaustion is a separate incomplete operation, not
any of those completed classifications.

Specialization accepts validated proper arguments and sealed constructor handles
from the same checked Program. It checks every argument before substituting,
preserves declaration/binder identity and definition scope, expands the proper
result and rechecks it. A foreign Program's handle is not accepted merely because
its module bytes match. `DataDescriptor` additionally establishes Data before
returning a descriptor for ground WF/Conforms. Open declarations stay immutable.

Admission, aggregate semantic work and expanded output have separate budgets.
Interfaces, comparisons, summary environments and cache growth consume work;
materialization charges repeated output occurrences even when internal structure
is shared. A compact summary can succeed while expansion exceeds its budget.
Failures return no usable partial result. Data preservation does not establish
recursion guardedness; application-argument closure and final descriptor checks
remain independent.

## Bittype and native generation

Bittype's syntax profile is `horos-declarations-3/ast-1`; resolution and analysis
use the corresponding `resolution-1` and `analysis-1` profiles. The frontend
exports checked Horos interfaces, shapes, summaries and identities through
structured provider-owned facts. Its `type-facts` obligation is required for
model publication. The [agnostic graph](semantics.md) stores expression callees
and subject-associated facts without interpreting Type, Data, Wire or libraries.

The existing [pagination mappings](semantic-generation.md) still emit one open
Page and preserve ordinary applications in TypeScript, Go, Rust and C++. All four
return located `Unsupported` diagnostics and no usable output for residual
constructor parameters or higher-kinded forms. They check the complete exported
closure, even for a closed use such as WrappedPage; they neither silently
specialize it nor discard an unsupported declaration. Native compilation does
not establish Horos Data for arbitrary native arguments.

## Explicit deferrals

| Deferred feature | Trigger for a separately qualified extension |
| --- | --- |
| Partial application, constructor-valued results, bare constructor aliases and source lambdas | An accepted case must construct or return a constructor and cannot use named saturated declarations |
| Quantified kind variables | A declaration must abstract over several shape kinds, beyond nested fixed constructor inputs |
| Conditional or intersected abstract interfaces | The same abstract F must both admit arbitrary Type inputs and guarantee Data for Data inputs |
| C++ template-template, Rust family/GAT or TypeScript encodings | A consumer selects a representation with complete target syntax and native compiler qualification |
| Explicit constructor binding or residualization | A native consumer supplies named roots, bindings and an output identity policy; the result needs a fresh checked document and receipt |

For now an abstract `F: Data -> Data` admits only Data calls, while
`F: Type -> Type` promises no Data result. Known Page can establish both uses
independently. A parameterized forwarding declaration such as
`type P<T> = Page<T>;` remains valid and has its own declaration identity; a bare
`type P = Page;` has a constructor-valued body and is rejected.

No kind validity, Data certificate or generated type establishes a runtime
protocol, authority or behavioral equivalence. Function, Stream and Cell remain
ordinary libraries; their names and structural fingerprints select no core rule.
