# Proposed canonical type model

**Date:** 31 October 2026. **Status:** proposal for review; no canonical language
implementation or published language revision. The generic AST and pipeline
foundation is implemented separately; see [the generation design](generation.md).

The owner selected bittype as the home for canonical type definitions and type
generation for TypeScript, Go, Rust and C++. The first milestone is types only.
The [charter](../CHARTER.md) records that boundary; this document proposes the
model, authoring form and implementation stages.

## 1. Direction and scope

A versioned, language-independent definition is the authoritative input.
TypeScript, Go, Rust and C++ are intended generated representations of that definition:

```text
JSON source -> checked canonical AST -> target-language AST -> native CST -> source
```

JSON is the initial syntax for authoring definitions. It does not make JSON the
value model or prescribe a serialization format for instances. A byte sequence,
for example, is a byte sequence even though its type is declared in a JSON file.

This milestone covers data shapes, declarations, imports, generic parameters and
applications. Output consists of native type declarations, any necessary support
data types for presence/nullness, and a generation manifest.

It does not include adapter generation, callable or stream implementations,
application-value validators, serialization, Wire protocols, service operation
tables, package publication or consumer migration. Checking a type declaration
is different from validating an application value.

## 2. Model layers

| Layer | Contents |
| --- | --- |
| Module | Language revision, logical module name, imports and named declarations |
| Declaration | Name, ordered generic parameter list and a type expression |
| Type expression | A primitive, structural composition, parameter, reference or application |

Documentation and target-specific naming are separate metadata. They do not
alter a type's value meaning. Source locations accompany diagnostics but are
not part of that meaning.

A resolved reference retains its declaring module and declaration, rather than
only a bare name or flattened structural body. A parameter reference points to
its lexical binder. An application retains its constructor and resolved
arguments, as well as any normalized view needed by a consumer.

The planned canonical AST is the common input to all four AST-to-AST generators.
Their spelling and memory-layout choices do not rewrite its semantics.
Each target AST is formatted into a native CST, then printed by the shared
printer, independently of the canonical model. The toolchain implementation
language is Go. The [transformation proposal](generation.md) uses one
ASTTransformer contract for both generation and explicit native-to-canonical
adoption; reverse mappings cannot infer semantic information absent from source.

## 3. Initial value-type vocabulary

| Form | Proposed source expression | Meaning |
| --- | --- | --- |
| Primitive | `"string"` | One precisely specified scalar domain |
| Record | `{"record": {"name": {"type": "string"}}}` | A closed set of named fields |
| Variant | `{"variant": {"ok": "string", "none": {"record": {}}}}` | Exactly one named alternative with its payload |
| List | `{"list": "string"}` | A finite ordered sequence of one element type |
| Tuple | `{"tuple": ["string", "i64"]}` | A fixed ordered sequence of potentially different types |
| Map | `{"map": "string"}` | A finite unordered map from string keys to values of one type |
| Nullable | `{"nullable": "string"}` | A value of the underlying type or explicit null |
| Reference | `{"ref": "Person"}` | A nongeneric local or imported declaration |
| Parameter | `{"parameter": "T"}` | A slot bound by the current declaration |
| Application | `{"apply": "Page", "with": {"T": {"ref": "Person"}}}` | A generic declaration with all arguments explicitly supplied |

Only primitive names are bare strings in type-expression positions. Parameters
and declarations have explicit forms, so name resolution does not guess which
sort a string denotes. An expression selects exactly one form.

The primitive proposal is:

| Name | Value domain | Initial native representation |
| --- | --- | --- |
| `boolean` | True or false | TypeScript `boolean`; Go and Rust `bool` |
| `string` | Finite Unicode scalar sequences; no implicit Unicode normalization | TypeScript/Go strings; Rust `String` |
| `bytes` | Finite sequences of octets | `Uint8Array`, `[]byte`, `Vec<u8>` |
| `i64` | Integers from -2^63 through 2^63-1 | `bigint`, `int64`, `i64` |
| `u64` | Integers from 0 through 2^64-1 | `bigint`, `uint64`, `u64` |
| `f64` | Finite IEEE 754 binary64 values; no NaN or infinities | `number`, `float64`, `f64` |
| `null` | The single explicit null value | Target-specific null data representation |

There is no platform-sized integer or unspecified `number` primitive. TypeScript
must not approximate a 64-bit integer with `number`. Native strings and numeric
types may admit values outside these domains; types-only generation does not
claim to check those values at runtime. The treatment of signed floating zero
in equality and canonical encoding must be fixed before the language is frozen.

A record's field names and a variant's alternative labels carry meaning.
Their source order does not. Tuple positions, list order and declaration
parameter order do carry meaning. Map keys are unique strings; map iteration
order is not part of its value.

Records are closed in this first proposal. Undeclared fields do not become
semantic members merely because a native value can store them. Variants are
closed and have at least one alternative. Their payload can be any admitted
value type, including a scalar, record, collection or another variant.

A variant label is not the name of a serialized discriminator property.
Properties such as `tag`, `kind` and `value` belong to a selected representation;
they are not additional fields in the canonical variant.

### Derived types

- Unit is the empty record.
- An enum is a variant whose alternatives all carry unit.
- An alias is a declaration whose body refers to another type. An alias adds
  no value-level wrapper merely by naming the type.
- `Page<T>`, `Result<T, E>` and a labelled `Option<T>` are ordinary declarations
  composed from the vocabulary above. They are not additional core kinds.
- A labelled `Option<Option<T>>` preserves its nested alternatives. It is not
  automatically the same type as nullable `T`.

Dates, timestamps, application identifiers and other semantic value domains
need explicit library definitions or a later refinement proposal. They must not
be silently approximated by a primitive with a different meaning.

## 4. Proposed authoring form

This complete example is proposed bittype syntax, not an existing parser API:

```json
{
  "bittype": 1,
  "module": "example.people",
  "types": {
    "Person": {
      "type": {
        "record": {
          "name": {
            "type": "string"
          }
        }
      }
    },
    "Page": {
      "parameters": ["T"],
      "type": {
        "record": {
          "items": {
            "type": {
              "list": {"parameter": "T"}
            }
          },
          "next": {
            "type": {"nullable": "string"},
            "optional": true
          }
        }
      }
    },
    "PeoplePage": {
      "type": {
        "apply": "Page",
        "with": {
          "T": {"ref": "Person"}
        }
      }
    }
  }
}
```

Omitted `parameters` means an empty list. Omitted `optional` means false.
The parser normalizes those defaults into the checked model. Missing `type`
is an error; it never means an unrestricted value.

`Page` remains an open generic declaration. `PeoplePage` is its application
to `Person`. Consumers do not author an expanded page interface or a separate
runtime description of its item type.

An `imports` member will associate local import names with immutable module
references. An imported declaration is addressed through that import, for
example `{"ref": "people.Person"}`. Dependency resolution must pin and verify
the actual input content; a mutable package label alone is insufficient.
The release-manifest and import-pin spelling remains to be specified before
the parser's input format is frozen.

## 5. Presence and null

Presence belongs to a field. Nullability belongs to its value type.

| Field | Absent | Present null | Present non-null value of T |
| --- | --- | --- | --- |
| Required T | Invalid | Invalid unless T already admits null | Valid |
| Optional T | Valid | Invalid unless T already admits null | Valid |
| Required nullable T | Invalid | Valid | Valid |
| Optional nullable T | Valid | Valid | Valid |

The example's `next` therefore has three states: absent, present null, and
present string. The empty string is a present string. Empty bytes, lists and
maps likewise do not imply absence.

Nullable adds the null value; applying it again introduces no new semantic
state. Alias resolution and generic substitution must respect that rule.
A backend must not accidentally make nested native wrappers introduce extra
canonical null states. This is distinct from explicit labelled optional
variants, whose nesting remains observable.

For TypeScript, the proposed field projections are `field: T`,
`field?: T`, `field: T | null` and `field?: T | null`. They require
`strictNullChecks` and
[`exactOptionalPropertyTypes`](https://www.typescriptlang.org/tsconfig/exactOptionalPropertyTypes.html);
explicit `undefined` is not an extra canonical field state.

## 6. Generics, references and recursion

The initial generic mechanism is first-order: parameters stand for value types.
A constructor such as `Page` is not itself a value type or a higher-kinded
argument. Open declarations are allowed; each application supplies every named
parameter exactly once. Missing, extra and unbound arguments are errors.

Parameter names are lexical binders. Renaming a binder and its uses preserves
the value-type family. Substitution must avoid capture across declarations and
imports. Resolution converts named arguments to the declaration's ordered slots
without losing their source locations or declaration origins.

The implementation must demonstrate:

- A separately rendered open generic and its concrete application describe the
  same applied value shape as capture-avoiding substitution.
- Nested applications preserve their argument types and declaration origins.
- Two imports using the same local parameter spelling do not capture each other.
- Unknown references and unresolved dependencies produce diagnostics rather
  than an `any`, `unknown`, empty record or placeholder identity fallback.

Recursive types describe finite values. A tree with a possibly empty list of
children, an optional next node, or a variant with a terminating alternative
can have finite instances. A required self-reference with no terminating case
is rejected. Direct alias cycles are also rejected.

The first revision should admit regular recursive generics whose recursive
references preserve the recursive group's parameter vector. Expansions that
continually grow their arguments, such as a recursive application from `F<T>`
to `F<List<T>>`, are outside the initial finite-description model.

Mutually recursive declarations belong in one module initially; module imports
are acyclic. Host-language cyclic objects are not automatically canonical values.

Pointer or `Box` insertion is a native layout choice. It must not prohibit a
valid canonical recursive type merely because inline storage would be infinite.
Go's [recursive struct rules](https://go.dev/ref/spec#Struct_types) illustrate
the distinction between the value model and a target layout.

## 7. Declaration identity and canonical form

Here, canonical means the authoritative, checked model from which native
representations are derived. It does not mean that a hash encoding has already
been chosen or published.

Keep three observations separate:

1. Which declaration an import or application refers to.
2. Which value shape results from resolving and applying that declaration.
3. Whether a behavioral contract or native binding is compatible with another.

An alias can add no value-level wrapper while its origin remains available to
tools. Two differently owned declarations can have the same value shape.
That fact does not authorize substituting their future library-owned bindings.

Defaults, binders, aliases, nullable composition, map ordering and recursive
graphs need specified normalization rules. A later canonical encoding must
state exactly which of those views it encodes, what names it retains and what
equivalence it promises. Sorting source JSON alone does not settle recursive
graph identity or behavioral compatibility.

Documentation, source locations and native spelling belong outside value
meaning. Preserve them separately for diagnostics and generated-code provenance.
Do not erase declaration origins to manufacture a universal structural identity
for Function, Stream, Cell or another behavioral contract.

## 8. Native type generation

All backends consume the same checked model and explicit target configuration.
The following initial representation proposals cover TypeScript, Go and Rust.
C++ has now been added as a target AST vocabulary; its canonical representation
mapping remains to be specified. This slice adds no canonical AST or generator.

| Canonical form | TypeScript | Go | Rust |
| --- | --- | --- | --- |
| Record | Interface or object type | Struct | Struct |
| Variant | Discriminated union | Tagged struct representation | Enum |
| List of T | `Array<T>` | `[]T` | `Vec<T>` |
| Tuple | Tuple type | Struct with positional fields | Tuple |
| String-keyed map | String-keyed record/map type | `map[string]T` | String-keyed map type |
| Signed 64-bit integer | `bigint` | `int64` | `i64` |
| Generic page | `Page<T>` | `Page[T]` | `Page<T>` |
| Optional nullable string field | `next?: string | null` | `Presence[Nullable[string]]` | `Option<Option<String>>` |

The Go support names are illustrative data types, not existing APIs.
The concrete optional-nullable string row preserves all three states.
General nested-nullable applications still require the normalization rule in
section 5; the row is not permission to add a wrapper for every source alias.

Rust's [payload-carrying enums](https://doc.rust-lang.org/reference/types/enum.html)
naturally represent variants. Go needs an explicit tag and payload layout.
In this types-only milestone, that layout can admit inconsistent combinations;
generated declarations alone do not enforce every canonical invariant.

Likewise, TypeScript structural typing does not prove exact closed records,
native numeric types can exceed a declared domain, and native strings may need
additional checks to meet the Unicode domain. Record these limitations in the
target contract. Do not describe a compilation pass as application-value
validation, or silently widen the canonical model to match a native limitation.

Target rules must cover identifier escaping, collisions, anonymous type naming,
recursive indirection, presence and null representations, and generic aliases.
Output naming choices are deterministic and kept outside value semantics.

A generation manifest records the input and generator revisions, target
configuration and supported language/toolchain version. Generated files do not
import the generator, a TypeScript compiler, Wire, bitruntime or a default
contract library.

## 9. Library ownership

Function, Stream, Cell and future behavioral contracts remain ordinary libraries.
Their behavior, protocols, adapters and special native realizations belong to
those libraries. This milestone implements none of them.

The core vocabulary above describes value construction. A library's declaration
uses references and applications through the same machinery as any other
declaration. The checked graph retains those associations so a later, separately
scoped consumer can select library-owned realizations explicitly.

The following are forbidden:

- Core `function`, `stream`, `cell` or other contract-specific kinds.
- Core name switches, hardcoded library imports, default registrations or a
  central catalogue that must change when a library is added.
- Inferring library ownership or protocol compatibility from matching shapes
  or a universal structural fingerprint.
- Erasing imports, aliases or generic applications in a way that loses the
  declaration needed to select its library-owned realization.
- Substituting TypeScript extraction for canonical input.
- Adding adapter, RPC, stream or Cell behavior under this types-only milestone.

## 10. Implementation stages and evidence

This document specifies the model proposal. The Go foundation contains
[generic AST data types](ast.md), a [typed stage contract](generation.md) and
[four target vocabularies with language identity](syntax.md). The canonical AST
is explicitly deferred. Native parsers, syntax checks, AST-to-CST formatters and
the shared CST printer are implemented separately. Canonical parsing/checking,
actual AST-to-AST mappings and public builders remain planned. The syntax engine
and typed transformation frame are implemented separately; see [engine.md](engine.md).

1. **Specify the model and fixtures.** Resolve the outstanding syntax and
   normalization choices below. Write language-neutral valid/invalid examples,
   including the expected resolved structures and substitution results.
2. **Implement one frontend.** Parse strictly, resolve imports and declarations,
   check binders, applications and recursion, and produce the shared model with
   useful source diagnostics. Reject duplicate JSON object keys before an object
   map can discard them, unknown semantic forms and incomplete inputs.
3. **Implement all four AST-to-AST generators.** Use the same fixtures and
   canonical AST to produce TypeScript, Go, Rust and C++ ASTs. Add only the support
   data declarations their representations need. Do not add adapters or runtime
   validators, and do not emit source strings from these transformations.
4. **Compose the native text stages.** The implemented native formatters consume
   their target AST and construct CST; the printer emits source. They do not
   depend on canonical input or make semantic representation choices. Compose
   these stages with the selected AST transformations at the application root.
5. **Qualify the representations.** Compile generated declarations and native
   use examples with pinned toolchains on shared build workers or GitHub
   runners. Compare intended value distinctions, not merely generated strings.

Acceptance examples include open and applied pagination; nested generics;
same-spelled parameters in separate scopes; scalar and record variant payloads;
tuple and map values; recursive trees; all presence/null cases; integer
boundaries; nullable normalization after application; and invalid references,
duplicates, arity and recursion.

A successful types-only result establishes checked declarations and compilable
native representations with documented limits. It does not establish lawful
behavior, value validation, transport compatibility or a deployed consumer.

## 11. Decisions still needed before freezing revision 1

- The immutable import and release-manifest syntax, identifier grammar and
  target naming/collision rules.
- Exact normalization and identity content, including recursive graph rules,
  signed floating zero and unused generic parameters. A hash format and its
  independent vectors follow those decisions.
- Each backend's detailed representation contract, especially Go variants,
  presence/null data types and nested nullable generic applications.
- Supported target-toolchain versions. The Go host language and component
  boundaries are selected in the [generation design](generation.md).

Refinements and arbitrary predicates, extensible records and variants,
inheritance, nominal wrappers, higher-kinded parameters and generic inference
are outside this initial proposal. They require a concrete semantic extension,
not an undocumented permissive fallback.

The proposal does not close the historical language-design issue or qualify a
released canonical encoding. It supplies the current reviewable model and
scoped path to the four type generators.
