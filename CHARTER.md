# Charter

[Bitwire decision 0010](https://github.com/Bitspark/bitwire/blob/main/docs/decisions/0010-bitwire-holds-the-contract-and-bitruntime-implements-it.md)
requires every new repository to open with answers to four questions. A
repository owns an independently useful compatibility commitment. A module owns a
coherent semantic decision.

## 1. What decisions does it own?

- **The contract language:** its syntax and semantics for value types, abstract
  operations, callable signatures, generics and imports. It covers a contract's
  shape and the declared behavior the language can state.
- **Canonical declaration identity:**
  - the normative canonicalization specification;
  - its shared vectors;
  - one library per supported language.

  Before an encoding is fixed, an audit settles exactly what it contains.
- **Native type rendering per language.** This is the theory's type generator:
  given a contract, it produces a native realization together with the contract
  it came from.
- **Runtime support for presence and absence** used by generated types.
- **The generation kernel's interface**, at first, in a separate module that may
  not import the declaration language. It moves to its own repository once its
  interface stabilizes and an independent pipeline needs it.

It does **not** own:

- mapping operations onto the wire, profiles or adapters, which are Bitlink's;
- descriptors and validation, which are bitschema's;
- the runtime, which is bitruntime's;
- the theory, which is bittheory's.

## 2. What does it promise consumers, and how is that versioned?

- A published language revision is immutable. A change to what is accepted,
  rejected or meant is a new revision.
- The canonical encoding is versioned by name. Any change to identity content is an
  explicit transition, never a side effect of packaging.
- Generated output records its language, encoding and generator versions in a
  generation manifest beside it, not in file headers.

Before 1.0 there is no compatibility promise.

## 3. What independently written evidence checks the promise?

- Canonicalization vectors shared by every language's identity library. Compile-time
  specialization and run-time generic application must reach the same identity.
- Tests of the generic-application laws at the semantic level, not only on
  generated text.
- Construction-square fixtures that Bitlink uses: instantiating a generic and then
  adapting it agrees with adapting it and then instantiating.
- bitschema's description tables: the descriptions bittype emits must equal the
  ones bitschema validates against, with neither repository importing the other.
- bittheory's integration checks against released versions.

## 4. Which real change becomes easier or safer because it is separate?

- The language can evolve, and be used for types without protocols (by storage,
  UI or command-line generators), without pulling in a wire runtime.
- Declaration identity gets one implementation per language instead of two
  diverging ones.
- Declaring a model stays independent of how it is carried: a model operation is
  not a network protocol.

## The first milestone

The language design, starting from its value-type core. BitTree is the first
generated consumer once bittype, bitschema and Bitlink together produce a working
slice. That slice must cover absence and null, errors, reverse calls and generic
identity.
