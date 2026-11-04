# Public model dependencies

The 4 November 2026 alignment pins `github.com/Bitspark/horos/core/go v0.8.0`
and `github.com/bitspark/ontos v0.16.0`. The Horos source pseudo-version and
Ontos 0.15 dependency are removed. Upgrade consumers together; there are no
compatibility aliases, alternate version branches, private checkouts or local
module replacements.

Horos 0.8.0 removes the template-lifting API and its exclusive helpers. Bittype
removes the corresponding transformer and has no compatibility path to 0.7.0.
The release includes ground declarations and the optional static Wire declaration
profile in Go. The two Bittype frontends and their private projection helpers use that model. Generic AST/CST, contracts,
semantic facts, engine and native stages keep their existing import boundaries.
The architecture tests continue to enforce those boundaries.

Ontos 0.16.0 also publishes the optional mixed value carrier in Go and TypeScript.
Bitwire 0.7.0 and Bitruntime 0.8.0 use its native structure; Bittype's static
frontends need only ground Ontos module descriptions. There is no production
dependency here on Ontos/Wire, Bitwire or Bitruntime, and no generated runtime
adapter or live mixed-value checker. Function, Stream and Cell remain library-owned.

The maintained [Go workflow](../.github/workflows/go.yml) runs the existing
provider, semantic, pagination, generic binding, specialization and direct optional declaration
tests against these pins. It also copies an [independent consumer](../integration/go/testdata/installed/consumer_test.go)
outside the checkout, downloads the exact public Bittype revision with clean
module/build caches, selects Horos 0.8.0 and Ontos 0.16.0, rejects module
replacements and runs ground pagination and directly authored optional pagination with data
and capability-type specialization. This checks the installed module graph, rather than relying on
a workspace replacement. Native compiler, race, fuzz and documentation gates remain.

Horos separately qualifies TypeScript constructor identity with public Ontos
0.16.0 and Bitwire 0.7.0 in its installed package gate. Neither qualification is
a service deployment claim. Historical release records and research evidence
retain their original inspected versions.
