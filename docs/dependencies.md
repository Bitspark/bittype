# Public model dependencies

Bittype pins public `github.com/Bitspark/horos/core/go v0.9.0` and
`github.com/bitspark/ontos v0.16.0`. Horos's module tag and anonymous Go proxy
identify source `e2896317a1f4b080cf0c4e8e46234d6ad680d141`; checksums come from
the public Go checksum database. The development pseudo-version is removed.
There are no local replacements, private module paths or compatibility aliases.

Horos 0.9.0 removes the separate Wire declaration package and introduces explicit
Type/Data bounds in one language. Bittype has one AST/check/resolve/analyze/provider
family and removes all old frontend exports and profiles. Ground checking and
Ontos value meanings are unchanged. Generic trees, contracts, shared semantics,
engine and native stages keep their model-independent import boundaries.

Ontos 0.16.0 also supplies the optional mixed value carrier. Bitwire 0.7.0 and
Bitruntime 0.8.0 use that structure; static Bittype needs only ground module
descriptions. There is no production dependency on Ontos/Wire, Bitwire or
Bitruntime, no runtime adapter and no live capability validator.

The maintained [Go workflow](../.github/workflows/go.yml) qualifies binding,
specialization, AST and semantic bounds, pagination, limits and provider contracts.
It copies an [independent consumer](../integration/go/testdata/installed/consumer_test.go)
outside the checkout, downloads the exact public Bittype revision with clean
caches, verifies dependency pins and rejects replacements. That consumer applies
one Page to data and Wire types, exercises its DataPage alias, and passes the
data specialization to existing WF/Conforms. Native compiler, race, fuzz and
documentation gates remain required.

Horos separately qualifies public npm, Cargo and Go installations, including
TypeScript constructor identity with Ontos 0.16.0 and Bitwire 0.7.0. These are
package qualification results, not service deployment claims. Historical release
and research evidence retains the originally inspected versions.
