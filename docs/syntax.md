# Target-language AST declarations

**Status, 31 October 2026:** TypeScript, Go, Rust and C++ syntax vocabularies
are implemented in Go on the [generic AST](ast.md). Each package supplies
`Language`, `AST`, `Node`, `Kind`, and constants for kinds, roles and attributes.

The shared language.Provider interface, native text providers, parsers,
AST-to-CST formatters, checks and generic printer are implemented. The
[Horos profile](horos.md) adds canonical declarations and [readable syntax](horos-syntax.md);
actual language mappings and public builders remain deferred. The syntax engine
now provides [explicit registration and route execution](engine.md).
See the [native text design](parsing-formatting.md).

The exported carrier fields remain unchecked; the separate language checker
validates supported type syntax. These are not complete compiler ASTs.

## Identity and registration

[contracts/language/go](../contracts/language/go/language.go) declares the agnostic interface:

```go
type ID string

type Provider interface {
    ID() ID
    Name() string
}
```

| Provider | ID | Display name |
| --- | --- | --- |
| [typescript.Language](../syntax/typescript/go/nodes.go) | `typescript` | TypeScript |
| [golang.Language](../syntax/golang/go/nodes.go) | `go` | Go |
| [rust.Language](../syntax/rust/go/nodes.go) | `rust` | Rust |
| [cpp.Language](../syntax/cpp/go/nodes.go) | `cpp` | C++ |
| [horos.Language](../syntax/horos/go/nodes.go) | `horos` | Horos |

The zero-value providers implement this interface. The methods return fixed
metadata and perform no IO, discovery or registration. IDs are exact,
case-sensitive language-family keys; names are display text. They are not
toolchain versions, schema versions or structural fingerprints.

**Owner correction, 31 October 2026:** there will be a central, agnostic language
AST registry. The earlier statement that there would be no language registry
is superseded. The syntax engine implements it through typed/profile descriptors
and explicit native-provider adaptation, without concrete syntax imports.
Duplicate keys and directed edges are rejected. Configured descriptors require
construction; the native markers above keep usable zero values. There are no
hidden defaults or automatic init-time registrations.

Each provider's `Language` is also a distinct named Go marker:

```go
type AST = ast.AST[Language]
type Node = ast.Node[Language]
type Kind = ast.Kind[Language]
```

These aliases use the existing carrier, not parallel tree storage.
The generic AST remains unconstrained by the identity interface. An independent
consumer can implement the interface and use its own marker without editing
either foundation. Equal ID strings do not establish Go type identity or
grammar compatibility, and do not permit converting an `AST[L]` into another
language's AST.

## Shared notation and ownership

The tables below specify the intended data shapes, not executable checkers.
Kind names omit the exported `Kind` prefix. Attribute and role names show
their actual string values; the packages export their corresponding
`Attribute...` and `Role...` constants.

- A role without a suffix means exactly one child; `?` means zero or one,
  `*` means zero or more, and `+` means one or more. Repeated children retain
  order. A type expression means one of that language's type-expression kinds.
- Attributes are unique by name. Names, qualified-name segments and package
  qualifiers are individual identifiers unless a row explicitly permits other
  text. Compound types and generic arguments belong in child nodes.
- Boolean attributes are the text `true` or `false`; absence means false.
  Other defaults are stated below. No arbitrary source-code attribute exists.
- An optional reference with role `target` on a named type use, name segment
  or lifetime use can retain its resolved declaration/import/binder. It is a
  qualified `ast.NodeRef`, not a containment edge. Omitting it means no resolved
  link is supplied, not an invented target. Resolution is not implemented.
- Generic parameter declarations and their uses remain separate nodes.
  Same-spelled parameters in different scopes can point to different binders.
  Applications retain their constructor/path and ordered arguments.
- All nodes can carry the generic `Origin`. Neither origin nor a target address
  proves semantic identity, ownership or compatible behavior.

Only the listed forms and attributes belong to this vocabulary. Separate native
checkers enforce structure, roles, lexical atoms and supported contextual rules.
Complete compiler/type checking and cross-tree resolution remain outside them.
Parsing retains original source; formatting applies the documented
[normalizations and layout](parsing-formatting.md).

Function-type and method-signature nodes describe native language syntax.
They contain no executable bodies. They do not give the Function library,
Stream, Cell or another behavioral contract a privileged kind in the generic
AST, identity interface, pipeline or registry.

## TypeScript

Source vocabulary: [nodes.go](../syntax/typescript/go/nodes.go).
Native reference: [TypeScript object types](https://www.typescriptlang.org/docs/handbook/2/objects.html).

| Kind | Attributes | Children |
| --- | --- | --- |
| SourceFile | — | `import*` Import; `declaration*` TypeAlias/Interface |
| Import | `module`: decoded module specifier | `binding+` ImportBinding, or one NamespaceImport |
| ImportBinding | `name`: imported identifier; optional `alias`: local identifier | — |
| NamespaceImport | `name`: local namespace identifier | — |
| TypeAlias | `name`; `exported` | `type-parameter*`; `type` |
| Interface | `name`; `exported` | `type-parameter*`; `base*` named/applied types; `member*` |
| TypeParameter | `name` | `constraint?`; `default?` types |
| TypeReference | `name` | — |
| QualifiedType | `name`: final identifier | `qualifier` TypeReference/QualifiedType |
| TypeApplication | — | `constructor` TypeReference/QualifiedType; `argument+` types |
| KeywordType | `name`: native type keyword | — |
| LiteralType | `literal-kind`; `text`: exact literal spelling | — |
| ArrayType | — | `element` type |
| TupleType | — | `element*` TupleElement |
| TupleElement | optional `name`; `optional`; `rest` | `type` |
| ObjectType | — | `member*` Property/Method/CallSignature/IndexSignature |
| UnionType, IntersectionType | — | `type+` (at least two) |
| ConditionalType | — | `check` type; `constraint` type; `true` type; `false` type |
| InferType | — | `type-parameter` TypeParameter without a default |
| MappedType | optional `readonly-modifier`, `optional-modifier`: `add` or `remove` | `type-parameter` TypeParameter with a constraint, without a default; `name-type?` type; `type` value |
| FunctionType | — | `type-parameter*`; `parameter*`; `result` type |
| Parameter | `name`; `optional`; `rest` | `type` |
| Property | `name`: exact property text; `optional`; `readonly` | `type` |
| Method | `name`: exact property text; `optional` | `signature` FunctionType |
| CallSignature | — | `signature` FunctionType |
| IndexSignature | `readonly` | `parameter` Parameter; `type` |
| TypeOperator | `operator`: `keyof`, `readonly` or `unique` | `type` |
| IndexedAccessType | — | `object` type; `index` type |

Imports are type-only imports. An absent binding alias means its imported name.
Interface members use the same member forms as ObjectType. Property text is
not an expression: a later formatter can quote a non-identifier property name.
Computed property keys are outside this subset.

KeywordType admits `any`, `unknown`, `never`, `void`, `undefined`,
`boolean`, `number`, `bigint`, `string`, `symbol`, `object` and `this`.
Literal kinds are `string`, `number`, `bigint`, `boolean` and `null`.
Their text is one matching native literal (including its string quotes or
numeric sign), never a template interpolation or arbitrary expression.
Checks cover rest/optional ordering, tuple naming consistency and direct operator
applicability. Semantic index-key and imported-type constraints remain compiler
obligations. Numeric property keys are admitted only as canonical safe decimal
integers; other spellings need numeric-to-property-text decoding.

ConditionalType represents `T extends U ? A : B` with four structured operands.
Conditional branches nest to the right; conditionals used as a check or extends
operand, union/intersection member, array element, operator operand or indexed
object are parenthesized when formatting requires it. Tuple wrappers such as
`[T] extends [U]` remain explicit tuple nodes. They are not discarded as redundant
grouping because TypeScript gives them different conditional behavior.

InferType introduces an ordinary TypeParameter child, optionally constrained by
`extends`, with no default. It is admitted only inside the extends operand of an
enclosing conditional. This checks syntactic placement; it does not resolve
uses of the inferred name or infer its type. Constrained infer operands are
grouped where needed to avoid absorbing a surrounding union/intersection.

```typescript
type Flatten<T> = T extends (infer Item)[] ? Item : T;
type FirstText<T> = T extends [infer Head extends string, ...unknown[]]
  ? Head
  : never;
```

MappedType is a type expression, separate from ObjectType and IndexSignature.
Its TypeParameter's `constraint` stores the source key type following `in`.
The optional `name-type` stores key remapping after `as`, and the final `type`
stores the value type. These children can themselves contain conditional or
mapped types. Same-spelled binders remain separate nodes; consumers can attach
ordinary `target` references to their uses without changing the shared carrier.

```typescript
type TextFields<T> = {
  [K in keyof T as T[K] extends string ? K : never]: T[K];
};
type MutableRequired<T> = {
  -readonly [K in keyof T]-?: T[K];
};
```

Absent mapped modifiers retain the original modifier behavior. `add` prints
`readonly` or `?`; `remove` prints `-readonly` or `-?`. The parser retains
explicit `+readonly` / `+?` spelling and trivia in source CST, while formatting
normalizes those plus signs away. Removing a modifier is therefore distinct
from leaving it unspecified. Mapping requires an explicit value type, admits
one binder, and cannot contain additional object members. Template-literal key
remapping remains outside the subset; ordinary type applications, indexed
access and conditional remapping are supported.

The language-owned stages follow the supported forms in the TypeScript
[conditional type](https://www.typescriptlang.org/docs/handbook/2/conditional-types.html)
and [mapped type](https://www.typescriptlang.org/docs/handbook/2/mapped-types.html)
documentation. Syntax admission accepts unresolved references. It does not
evaluate conditionals, infer substitutions, validate key membership, or prove
distributivity or canonical constraint preservation. Compiler fixtures test
representative formatted outputs, including type-only assertions about their
native results; that does not add a semantic provider or generator mapping.

Template-literal types, value declarations, class bodies, decorators and
executable expressions remain outside this slice. See the
[text-tooling boundary](parsing-formatting.md#checks-and-admitted-syntax).

## Go

Source vocabulary: [nodes.go](../syntax/golang/go/nodes.go).
Native reference: [Go types](https://go.dev/ref/spec#Types).

| Kind | Attributes | Children |
| --- | --- | --- |
| File | `package`: package identifier | `import*`; `declaration*` TypeDeclaration/TypeAlias |
| Import | `path`: decoded import path; optional `alias` | — |
| TypeDeclaration, TypeAlias | `name` | `type-parameter*`; `type` |
| TypeParameter | `name` | `constraint` type |
| NamedType | `name`: local or predeclared identifier | — |
| QualifiedType | `qualifier`: package identifier; `name` | — |
| TypeApplication | — | `constructor` NamedType/QualifiedType; `argument+` types |
| StructType | — | `field*` Field/EmbeddedField |
| Field | `name`; optional `tag`: decoded tag text | `type` |
| EmbeddedField | optional `tag` | `type` |
| InterfaceType | — | `member*` Method/EmbeddedType |
| Method | `name` | `signature` FunctionType |
| EmbeddedType | — | `type`: type or constraint TypeUnion |
| FunctionType | — | `parameter*` Parameter; `result*` Parameter |
| Parameter | optional `name`; `variadic` | `type` |
| PointerType, SliceType | — | `element` type |
| ArrayType | `length`: nonnegative decimal integer | `element` type |
| MapType | — | `key` type; `element` type |
| ChannelType | `direction`: `both` (default), `send` or `receive` | `element` type |
| TypeUnion | — | `term+`: types or TypeApproximation |
| TypeApproximation | — | `type` |

Names such as `any`, `comparable` and `string` are NamedType uses, not
canonical primitive kinds. An import alias may also be `_` or `.`.
Defined types and aliases remain different kinds. Each field or parameter
node represents one field/parameter, rather than a grouped name list.

A variadic parameter represents `...T` with `T` as its type child. Only the
last input parameter may be variadic; results may not be. Interface method
signatures cannot introduce method-specific type parameters. Constraint unions
and approximations belong to constraint positions, not arbitrary value unions.

Array lengths are literal integers in this initial vocabulary. Constant
expressions, constant/variable/function declarations, method implementations
and statements are outside this slice.

## Rust

Source vocabulary: [nodes.go](../syntax/rust/go/nodes.go).
Native reference: [Rust types](https://doc.rust-lang.org/reference/types.html).
Qualified paths and generic arguments follow the supported subset of
[Rust path syntax](https://doc.rust-lang.org/reference/paths.html).

| Kind | Attributes | Children |
| --- | --- | --- |
| File | — | `item*` Module/Use/TypeAlias/Struct/TupleStruct/Enum/Trait/TraitImpl |
| Module | `name`; `visibility` | `item*` |
| Use | optional `alias` | `path` PathType without generic arguments |
| TypeAlias | `name`; `visibility` | `generic-parameter*`; `where*`; `type` |
| Struct | `name`; `visibility` | `generic-parameter*`; `where*`; `field*` Field |
| TupleStruct | `name`; `visibility` | `generic-parameter*`; `where*`; `field*` TupleField |
| Enum | `name`; `visibility` | `generic-parameter*`; `where*`; `variant*` |
| UnitVariant | `name` | — |
| TupleVariant | `name` | `field*` TupleField |
| RecordVariant | `name` | `field*` Field |
| Field | `name`; `visibility` | `type` |
| TupleField | `visibility` | `type` |
| Trait | `name`; `visibility` | `generic-parameter*`; `bound*`; `where*`; `member*` MethodSignature/AssociatedType |
| TraitImpl | — | `generic-parameter*`; `trait` PathType; `subject` type; `where*`; `member*` TypeAlias |
| AssociatedType | `name` | `generic-parameter*`; `bound*`; `where*` |
| MethodSignature | `name` | `generic-parameter*`; `where*`; `receiver?`; `parameter*`; `result?` type |
| Receiver | `receiver-form`: `value`, `shared` or `mutable` | `lifetime?` LifetimeReference for a borrowed receiver |
| Parameter | `name` for a method; optional in FunctionType | `type` |
| TypeParameter | `name` | `bound*`; `default?` type |
| LifetimeParameter | `name` without apostrophe | `bound*` LifetimeReference |
| LifetimeReference | `name` without apostrophe | — |
| WherePredicate | — | `subject` type/lifetime; `bound+` |
| RelaxedBound | — | `path` PathType |
| PathType | `absolute` | `segment+` PathSegment |
| QualifiedPathType | — | `subject` type; `trait?` PathType; `segment+` PathSegment |
| PathSegment | `name` | `argument*` types/LifetimeReference/AssociatedTypeBinding |
| AssociatedTypeBinding | — | `path` PathSegment; `type` |
| PrimitiveType | `name`: native primitive | — |
| ReferenceType | `mutable` | `element` type; `lifetime?` |
| RawPointerType | `mutable` | `element` type |
| SliceType | — | `element` type |
| ArrayType | `length`: nonnegative decimal integer | `element` type |
| TupleType | — | `element*` types |
| NeverType | — | — |
| FunctionType | `unsafe`; optional `abi`: decoded ABI name | `parameter*`; `result?` type |
| TraitObjectType | — | `bound+` |

Generic parameter children are TypeParameter or LifetimeParameter. Bounds are
PathType or LifetimeReference. Type parameters, associated-type declarations
and type-subject where predicates also admit RelaxedBound, which wraps the path
in `?Sized` (or an unresolved native name). This is a native bound modifier,
not a type expression or a special canonical kind. Lifetime binders and
lifetime-subject where predicates only have lifetime bounds. RelaxedBound is
not admitted in supertrait or trait-object lists.
`visibility` is `private` (default), `public` or `crate`. Enum variant
fields have inherited visibility, so only the default is admitted there.
Primitive names are `bool`, `char`, `str`, `f32`, `f64`,
`i8/i16/i32/i64/i128/isize` and `u8/u16/u32/u64/u128/usize`.
Unit is an empty TupleType. A missing result denotes unit.
Reference mutability and raw-pointer mutability are explicit; a false raw
pointer flag denotes `*const`. Path segments may use native `crate`,
`self`, `super` and `Self` names where Rust permits them.

FunctionType is a bare function-pointer signature. Trait method receivers are
separate nodes, never ordinary first parameters. TraitObjectType represents
explicit `dyn` bounds. Default generic types belong only where Rust allows
them, not on method-specific, associated-type or impl parameters.

AssociatedType declares a trait member, with optional type/lifetime parameters,
result bounds and where predicates. A TraitImpl supplies a trait path and
implementing subject, including any type/lifetime parameters and where predicates.
Its members reuse TypeAlias for associated definitions. The parent context
distinguishes these from module aliases; there is no duplicate definition node.
An associated definition requires a RHS type, inherits visibility and cannot
declare result bounds. The corresponding declaration has no RHS or default.
These forms follow [Rust associated items](https://doc.rust-lang.org/reference/items/associated-items.html)
and [type aliases](https://doc.rust-lang.org/reference/items/type-aliases.html).

```rust
trait Family {
    type Item<T: Clone>: Clone where T: Default;
}

struct VectorFamily {}

impl Family for VectorFamily {
    type Item<T: Clone> = Vec<T> where T: Default;
}

type Bytes = <VectorFamily as Family>::Item<u8>;
```

Lifetimes and where clauses compose through the same nodes:

```rust
trait Lending {
    type Item<'a>: 'a where Self: 'a;
}

struct Holder<T> { value: T }

impl<T> Lending for Holder<T> {
    type Item<'a> = &'a T where Self: 'a;
}
```

AST child order for TypeAlias stays `generic-parameter*`, `where*`, `type` in both
contexts. The formatter places an impl member's where clause after its RHS,
while module aliases keep their existing layout. The parser also accepts Rust's
deprecated pre-equals placement in an impl, preserves it in source CST, and
normalizes it on formatting; two where clauses on one definition are rejected.
Formatting builds CST from the structured children directly. Declaration and
definition parameters have distinct node identities; existing target references
and origins remain intact.

TraitImpl admits empty impls and associated type definitions only. It accepts
unresolved trait and subject names. Syntax checking does not establish member
completeness, trait coherence, bound satisfaction, required GAT lifetime clauses,
or agreement between declaration and definition. Native compiler tests qualify
representative complete examples separately. No Rust semantic provider or
generation SPI behavior is introduced by these syntax nodes.

QualifiedPathType represents both `<T as Collection>::Item` and `<T>::Item`.
The subject remains a complete type child; the optional trait is an ordinary
PathType, including an absolute path and its arguments. Each suffix segment
retains its own name and arguments. Nested projections therefore compose without
flattening the qualified subject into a name:

```rust
type Element<T: Collection> = <T as Collection>::Item;
type Inferred<T: Collection> = <T>::Item;
type Nested<T: Collection> where <T as Collection>::Item: Collection =
    <<T as Collection>::Item as Collection>::Item;
trait Bytes: Iterator<Item = u8> {}
trait Linked<T: Collection>: Iterator<Item = <T as Collection>::Item> {}
```

AssociatedTypeBinding is an equality argument, not a standalone type expression
or an associated-type declaration. For `Iterator<Item = T>`, its `path` child is
the single PathSegment `Item` and its `type` child is `T`. The single segment
can itself carry arguments, as in `Family<Item<U> = V>` or
`Lending<Item<'a> = &'a T>`. Generic associated-type uses, declarations and
definitions compose directly, including `<T as Family>::Item<U>`.
Raw identifiers are retained. The binding name
must be an associated identifier, not `_`, `Self`, `self`, `super` or `crate`.
Path segments, including binding names, can retain ordinary `target` references
and origins through the existing carrier.

Within each argument list, lifetime arguments precede positional type arguments,
which precede equality bindings. The parser retains the original source in CST;
formatting normalizes optional type-path turbofish separators (`Trait::<T>` to
`Trait<T>`) and whitespace. Empty generic argument lists remain unsupported.
See [text stages and their limits](parsing-formatting.md#checks-and-admitted-syntax).

These are syntax rules only. Unresolved traits and associated names are admitted;
no trait membership, associated-type existence, equality, lifetime relationship
or projection normalization is established. Rust compiler fixtures separately
qualify representative formatted output, including independently authored ASTs.

Unit structs have no distinct AST form and are rejected.
Const generics, higher-ranked binders, associated-type argument bounds
(`Iterator<Item: Clone>`), associated-type defaults, `impl Trait`, macros,
attributes, inherent/unsafe/negative/specializing impls, associated constants,
impl methods and executable bodies are outside this slice. Type-only trait
impls are the supported implementation form.

## C++

Source vocabulary: [nodes.go](../syntax/cpp/go/nodes.go).
Native reference: [C++ declaration grammar](https://eel.is/c++draft/dcl.pre).

| Kind | Attributes | Children |
| --- | --- | --- |
| TranslationUnit | — | `include*`; `declaration*` Namespace/Record/TypeAlias/Enum |
| Include | `path`: header text without delimiters; `system` | — |
| Namespace | `name` | `declaration*` |
| Record | `name`; `record-key`: `struct` or `class`; `forward`; optional `access` when nested | `template-parameter*`; `base*`; `member*` |
| TypeAlias | `name`; optional `access` when nested | `template-parameter*`; `type` |
| Enum | `name`; `scoped`; `forward`; optional `access` when nested | `underlying?` type; `enumerator*` |
| Enumerator | `name`; optional `value`: signed decimal integer | — |
| Base | optional `access`; `virtual` | `type` NamedType |
| Field | `name`; optional `access` | `type` |
| Method | `name`; optional `access`; `const`; `volatile`; `virtual`; `pure`; `noexcept`; `ref-qualifier`; `variadic` | `parameter*`; `result` type |
| Destructor | optional `access`; `virtual`; `pure`; `noexcept` | — |
| Parameter | optional `name` | `type` |
| TypeParameter | `name`; `parameter-key`: `typename` (default) or `class` | `default?` type |
| NamedType | `global`; `typename` | `segment+` NameSegment |
| NameSegment | `name`; `template` | `argument*` types |
| BuiltinType | `name`: native builtin spelling | — |
| PointerType, LValueReferenceType, RValueReferenceType | — | `element` type |
| CVQualifiedType | `const`; `volatile` (at least one true) | `element` type |
| ArrayType | `length`: positive decimal integer | `element` type |
| FunctionType | `noexcept`; `variadic` | `parameter*`; `result` type |
| MemberPointerType | — | `owner` NamedType; `element` type |

Record members are Field, Method, Destructor, TypeAlias, Record or Enum.
Access is `public`, `protected` or `private`; absent member/base access uses
the enclosing record's native default. A true system include uses angle
delimiters; false uses quotes. Namespace nesting uses child nodes. The global
flag means leading `::`; name segments retain their own template arguments.
The `typename` and `template` flags describe dependent-name disambiguators.
A forward Record has no bases or members; a forward Enum has no enumerators.
Unscoped enum forward declarations require an underlying type. Forward forms
preserve the distinction between naming a type and defining an empty type.

Builtin spellings cover `void`, `bool`, the native character types, signed
and unsigned integer type specifiers, and `float`, `double`, `long double`.
Library types such as `std::int64_t` use NamedType; this AST assigns no
fixed-width meaning to native `long`.
`ref-qualifier` is `none` (default), `lvalue` or `rvalue`.
Pure declarations must be virtual. A destructor takes its name from its record;
this vocabulary supplies a declaration only, not a destructor implementation.
The noexcept flag requests an unconditional noexcept specification; omission
leaves the native default rather than asserting a particular exception outcome.

Pointer and qualifier nesting preserves declarator structure: a pointer to
const T differs from a const pointer to T. NamedType represents applications
such as `std::vector<T>` using separate name segments and argument nodes.
There is no source-string fallback for a difficult declarator.

Non-type/template-template parameters, parameter packs, specialization,
constraints/requires expressions, computed array bounds, initializers,
constructors, operator declarations, macros and function bodies are outside
this slice. Qualified member-function-pointer types are rejected because native
FunctionType has no method CV/ref fields. The maintained native text gates
provide separately scoped compiler qualification.

## Evidence and limits

[Consumer fixtures](../integration/go/testdata/syntax/typescript.go) and the
equivalent [Go](../integration/go/testdata/syntax/golang.go),
[Rust](../integration/go/testdata/syntax/rust.go) and
[C++](../integration/go/testdata/syntax/cpp.go) fixtures construct native
`Page<T>`, a `Person`, an applied `PeoplePage`, and a generic reader
signature. They keep same-spelled T binders distinct through target links.
They are independently authored examples, not canonical-to-native mappings.

The [tests](../integration/go/syntax_test.go) exercise provider metadata through
the common interface, type-check those consumers, reject mixed target ASTs,
and check every pair of public marker/AST/node/kind types for non-assignability.
A [separately owned provider](../integration/go/testdata/syntax/extension.go)
implements the same interface without any central language list.
The [architecture checks](../integration/go/architecture_test.go) enforce that
syntax packages depend only on the generic AST and language identity contracts
plus the generic CST carrier and standard-library utilities. The three foundations do not depend
on each other or on concrete languages.

The original contract tests establish Go API usability and isolation. Additional
native text tests cover checking, retained source and independent AST formatting.
Native compiler fixtures run separately on remote runners. These claims do not
establish executable behavior, registry operation or canonical equivalence.
The pipeline contract is unchanged; AST adds a borrowed view interface.

See [native semantic providers](native-semantics.md)
for the bounded Go resolver/analyzer, immutable native facts and explicit
admission queries. TypeScript, Rust and C++ semantic providers remain planned.
