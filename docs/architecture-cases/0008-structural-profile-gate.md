# bittype/ACX-008: The shared horos mapping policy admitted one profile

**Status:** Corrected, with limits; state of 12 November 2026. A mapping
states its profile in its own description, no shared code lists, writes or
refuses a profile, and the two identities are `structural-representation/1`
and `structural-enforcement/1`. A first form of the correction kept a refusal
of one name in the adapter and was replaced before merging. One independent
review read that first form. The reworked form gets another. The
[limits](#limits) say what the correction does not establish.
**Scope:** The common horos mapping policy,
[`transform/syntax/horos/mapping/go`](https://github.com/Bitspark/bittype/blob/ca5dc80/transform/syntax/horos/mapping/go/mapping.go),
at `ca5dc80`, and the four `horos-to-<native>` mapping packages that build
their libraries through it. Identified on 11 November 2026. Tracked in
[#216](https://github.com/Bitspark/bittype/issues/216).

## Rules involved

[AC-03](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-03)
requires that adding a contribution change only its owning implementation and
the consumer's composition. A representation profile is such a contribution:
the [accepted families design](../representation-families.md#domains) names the
lossless structural profile and "an explicitly named bounded profile" beside it.
The shared policy's `Library.Describe` returned `ErrUnsupported` for every
profile whose name was not the one constant `RepresentationProfile`. A second
profile that reused the horos patterns and the declaration recipe therefore
needed an edit to this shared file.

[AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)
asks reviewers to inspect helpers outside the kernel. That is where this gate
sat. The kernel was not at fault: it treats a profile as an opaque identity and
selects only the mappings registered under the requested one
([registry](https://github.com/Bitspark/bittype/blob/ca5dc80/transform/syntax/shared/generation/go/registry.go#L185)).

The profile's name was a second, related fault, and not a rule violation by
itself. Both identities began with the word `pagination`, which names the
first fixture and not what the profile promises
([the two constants at `ca5dc80`](https://github.com/Bitspark/bittype/blob/ca5dc80/transform/syntax/horos/mapping/go/mapping.go#L15-L16)).
It was corrected with the gate because the plan records profile identities
([#195](https://github.com/Bitspark/bittype/issues/195)).

A third fault was found on 12 November 2026, while the first form of the
correction was under review, and it is the cause of the gate.
[AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)
keeps a meaning with its semantic owner. A profile is a promise about the
native type, and every mapping constructor took it as an argument, so the
composition stated the promise of the mapping it registered. The gate existed
to stop a composition from registering a constraint-erasing library under the
enforcement name.

## Intended boundary

The kernel's selection by profile identity stays invariant. Profiles vary. Each
profile's owner supplies its mapping libraries under its own identity, and the
consumer's request selects one. The common horos policy supplies patterns and
the declaration recipe to any of them. An unsupported profile fails because no
mapping is registered for it, or because its owner refuses it, never because a
shared file lists the admitted names.

## Actual implementation and consequence

At `ca5dc80` the adapter began with this check:

```go
func (p Library[T]) Describe(from, to representation.Language) (...) {
	if p.Profile != RepresentationProfile {
		return ..., g.ErrUnsupported
	}
```

Every constructor of a native mapping library takes a `profile` argument, for
example `text.Library(profile)` and `record.Library(profile)` in the packages
under `transform/syntax/horos-to-<native>/mapping/<name>/go`. At `ca5dc80`
these were `Text(profile)` and `Libraries(profile)` in one package per pair;
[bittype/ACX-009](0009-package-boundaries.md) records that split. Exactly one
value of that argument worked. The documented purpose of the check was to
refuse the enforcement profile, whose encoding does not exist
([GEN-02](../generation-gaps.md#gen-02)). It was written as an allowlist of
one, so it refused every other profile as well.

Consequences:

- A bounded profile, or any second structural profile, could not be registered
  through the common policy without changing it.
- Integers, maps, unions and every later family would have been recorded in
  plan entries under an identity that says pagination. Renaming it after
  [#207](https://github.com/Bitspark/bittype/issues/207) writes the first entry
  would have changed an identity that binding backends pin.

This was architectural and not missing coverage: the missing profile was not
another contribution behind a sound contract. The shared helper had to change
to admit it.

## What worked and what it did not establish

All generation tests selected the one profile, and they passed. An independent
library could enter the kernel directly: an independent contribution in the
[composition test](../../integration/go/horos_generation_composition_test.go)
implements the mapping contract itself and sets the profile on its own
description. Those results established that the kernel is open. They did not
establish that the common horos policy was, because no test registered a second
profile through it.

## Required correction

1. The common policy holds no list of admitted profile names.
2. A request for the enforcement profile still fails with `ErrUnsupported` and
   never yields representation output. A request for a profile with no
   registered mapping fails.
3. The structural and enforcement profiles get identities that state their
   promise and name no fixture, with no alias for the old names. The owner of
   research 0015 chooses them.
4. The native declarations generated for the existing fixtures stay
   byte-identical.
5. Added on 12 November 2026: a mapping states its own profile. A composition
   chooses which mappings it registers and which profile it requests, and it
   does not hand a mapping its profile.

## Correction and qualification status

Corrected on 12 November 2026 by pull request
[#265](https://github.com/Bitspark/bittype/pull/265), which closes
[#216](https://github.com/Bitspark/bittype/issues/216). It is ordered before
#207 under the owner's direction of 11 November 2026 that a misplaced or
wrongly named part changes first and missing coverage is added afterwards
([AC-10](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-10)).

### The first form, and why it was replaced before merging

The first form of #265, at `fd75abe`, renamed the identities and removed the
closed list. It kept the constructors `Library(profile)`: the composition
passed a profile, and the adapter copied it into the description. Because a
composition could then build the constraint-erasing libraries under the
enforcement name, the adapter kept a refusal of that one name,
`EnforcementProfile`.

That form was not merged. The owner asked what the adapter is needed for, and
the holder of #216 decided on 12 November 2026 that the profile moves to the
mapping that promises it. The decision and its sources are comments on #265
and on [#196](https://github.com/Bitspark/bittype/issues/196). The reasons:

- **The refusal guarded a misplacement and did not remove it.** A profile says
  what the native type retains. Only the mapping's owner can promise that
  ([AC-02](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-02)).
  Research 0014 lists "Representation context: target and supported profile"
  among the elements that a contribution publishes. With `Library(profile)` the
  composition stated the mapping's promise.
- **It was a guard and not a guarantee.** An independent review of `fd75abe`
  registered the builtin TypeScript libraries under `EnforcementProfile`
  through a five-line wrapper that uses only exported API, and a request for
  that profile returned representation output. The check compared one exact
  string. The record of the first form said that representation output never
  carries the enforcement identity, which the code did not give.
- **It reserved an identity in a shared helper.** Implementing enforcement
  ([#162](https://github.com/Bitspark/bittype/issues/162)) would have needed an
  edit to that helper.
- **Its evidence for a second profile was a renamed copy.** The test
  registered the builtin libraries under another string. That does not
  establish an extension boundary.

### Correction

| Was | Is |
| --- | --- |
| Both identities began with `pagination` | `RepresentationProfile` is `structural-representation/1` and `EnforcementProfile` is `structural-enforcement/1`. The owner of research 0015 chose them and recorded the reason on #216. The Go names are unchanged. No alias exists |
| Every mapping constructor took a `profile` argument, and the adapter wrote it into the description | A mapping states its profile itself. Each builtin mapping writes `RepresentationProfile` into its own description, and its constructor takes no profile: `text.Library()`, `record.Library()`, and for Go `record.Library(key)` and `optional.Library(key)` with the helper key. The common declaration recipe states the profile for the four declarations mappings |
| The adapter refused every profile but `RepresentationProfile` | No shared code lists, writes or refuses a profile. `Library` in `transform/syntax/horos/mapping/go` is a function type with the method `Describe`. It returns what the function returns |
| The enforcement profile was refused as one of all other names | No mapping in bittype states `EnforcementProfile`. The registry finds no mapping for a request that names it and fails with `ErrUnsupported` |
| The integration fixture repeated the identity as a string | [`horosgeneration.Request`](../../integration/go/testdata/horosgeneration/request.go) takes it from the horos mapping package. One test beside the declaration pins both values on purpose. No other test or fixture spells either identity |

**Where each responsibility sits.**

- **The identities** stay in `transform/syntax/horos/mapping/go`. They are the
  two promises that the horos mapping side declares.
- **The promise** sits with the mapping that makes it. The profile is a field
  of the description that a library returns, and the library fills it.
- **Selection** sits with the registry, unchanged: it plans with the mappings
  that state the requested profile. The registry and `contracts/**` are not
  edited.
- **The choice** sits with the consumer: which libraries it registers and
  which profile it requests. It cannot choose what a builtin library promises.

**What remains of the adapter.** One thing: a description function becomes a
`generation.MappingLibrary`. The type has no field and no profile. It stays a
shared, exported function type for a reason found during the rework. The
provider matrix identifies a mapping library by a concrete exported type or by
a factory's concrete result type, and by its own rule a returned interface
does not establish one. A form in which every mapping returned the contract
interface from an unexported type removed all builtin mappings from the matrix
and failed `TestCurrentRepositoryProviderDiscovery`.

**Behaviour at the edges.** A nil function fails with `ErrInvalid`, in
`Describe` and at registration. A description that states no profile is
refused at registration by the contract's `FreezeDescription`, with
`ErrInvalid`. On `main` the gate answered both with `ErrUnsupported`. No
caller relied on either.

### Qualification

The tests are in
[`integration/go/horos_generation_profile_test.go`](../../integration/go/horos_generation_profile_test.go)
and beside the declaration in
[`mapping_test.go`](../../transform/syntax/horos/mapping/go/mapping_test.go).
The fixture contribution is
[`integration/go/testdata/readonlyprofile`](../../integration/go/testdata/readonlyprofile/library.go).

- **Every builtin mapping states the structural profile.**
  `TestBuiltinMappingsStateTheStructuralProfile` describes every builtin
  library of the four pairs and compares the stated profile. It checks that no
  constructor has a parameter of a string type, and that the lists it reads
  have as many entries as the pair has mapping packages. The calls
  `text.Library()` in the test compositions compile only while a constructor
  takes no profile.
- **A second profile as a contribution of its own.**
  `TestRequestSelectsTheContributionsOfItsProfile` registers the builtin
  TypeScript family and the fixture contribution in one registry. The fixture
  lives in its own package, states its own profile, maps a list to
  `ReadonlyArray<E>` and maps no record. For one module the structural request
  yields `string[]` and the fixture's request yields `ReadonlyArray<string>`,
  and each plan names its own list binding. A module with records fails under
  the fixture's profile although the builtin record library is in the same
  registry: a family is not completed from another profile. The horos mapping
  package and the registry know nothing of the fixture.
- **The enforcement profile.** `TestEnforcementProfileIsRefused` requests it
  from the registry that has just generated under the structural profile, in
  all four targets. The request fails with `ErrUnsupported` with no plan and
  no output.
- **An unregistered profile.** `TestUnregisteredProfileIsRefused` covers an
  empty registry and the four builtin sets.
- **The identities.** `TestProfileIdentitiesArePinned` spells both values
  beside their declaration. On the first form no test failed when a value was
  changed.
- **The function type and registration.**
  `TestLibraryReturnsWhatItsFunctionDescribes` and
  `TestRegistrationNeedsAStatedProfile` cover the edges above.
- **Seeded faults.** A builtin constructor that takes and uses a caller's
  profile breaks the build of the tests. The same with a variadic parameter
  fails the signature check. A builtin mapping that states the enforcement
  profile fails the profile comparison. The fixture stating the builtin
  profile fails the second-profile test, and with the test's own guard removed
  the registry answers "ambiguous generation mapping". Either identity changed
  fails the pin.
- **Output.** Between `main` at `6c788d8` and the correction, the 79 texts
  that the existing integration tests format have the same test names, SHA-256
  and lengths. So has the text of three fixtures in the four targets, twelve
  files, written by a dump helper that is not committed. The comparison covers
  formatted text. It does not cover the output trees that tests read without
  formatting them.
- **The whole suite.** `go test ./...` on the author's Windows machine gives
  the same outcome for every test that existed at `6c788d8`. 23 test outcomes
  are new and pass. The eight failing outcomes are the same before and after:
  the Rust semantic tests have no rustc helper on that machine.
- **Imports.** `TestProductionArchitecture` passes with no new entry and no
  new exception. No package and no production import is added. The fixture
  package under `testdata` and the test that uses it import what they need.

### Limits

- **Review.** One independent review read the first form at `fd75abe`, and
  its findings are in this record. The reworked form gets another review,
  which this record does not contain yet.
- **The registry does not judge a stated profile.** A contribution that states
  the enforcement profile is selected by a request for it. A wrapper that
  takes a builtin library's description and states another profile is such a
  contribution, and it yields representation output under that profile.
  `TestRegistryDoesNotJudgeAStatedProfile` reproduces this. It is a
  contribution of its own author, who states its profile and answers for it:
  selecting a contribution does not establish its claim
  ([AC-05](https://github.com/Bitspark/platform/blob/main/docs/concepts/architectural-correctness.md#ac-05)).
  What bittype gives is narrower and holds: no mapping in bittype states the
  enforcement profile, and no builtin constructor takes a profile.
- **A profile's family is complete in itself.** A builtin mapping states the
  structural profile, so another profile cannot register it as its own. The
  common patterns and the naming function stay usable by any contribution.
  The common declaration recipe states the structural profile. A further
  profile that wants the same recipe writes its own, as the fixture does, or
  the recipe is offered as a function it can call. That is an addition and is
  not made here.
- **The fixture is small.** It is one TypeScript library with three bindings
  and an alias recipe. No bounded profile with a narrower domain exists yet,
  and no second profile exists for Go, Rust or C++.
- **The signature check reads strings.** A constructor that took a profile
  through a type that is not a string would pass it. The comparison of the
  stated profile would then catch only what the listed calls produce.
- **One matrix entry per pair.** The provider matrix shows the builtin
  mappings of a pair as one entry, the shared function type, although the pair
  has five or six mapping packages. That was so before this change and is not
  examined here.
- **The names in plan entries.** No plan entry exists yet
  ([#207](https://github.com/Bitspark/bittype/issues/207)). That the new
  identities are what the entries record is for that work to show.
- **Verbatim records keep the old names.** The documents under `research-docs/`
  are submitted consultations and their replies and are unchanged.

Implementing enforcement ([#162](https://github.com/Bitspark/bittype/issues/162))
no longer needs an edit to shared code: an enforcing mapping states
`EnforcementProfile` in its own description.

## Assessment lesson

The check was reviewed as the refusal of one unsupported profile, which is
correct behaviour. Its form went unexamined: a refusal written as "not the one
admitted name" is a closed catalogue of one. When a helper takes an identity as
a parameter, ask how many values of that parameter can succeed, and who would
have to edit what to add another.

The first form of the correction repeated the mistake one level up. It asked
who should refuse the enforcement name and did not ask why a composition could
state a mapping's promise at all. When a correction needs a guard, ask which
misplacement makes the guard necessary, and move that instead. And state a
guarantee where it lives: a refusal by an exact name in a helper that a caller
can go around is a guard.
