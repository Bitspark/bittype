# bittype

A wire-independent contract language. Its parts:

- **A value-type core** (records, unions, aliases, generics) that is usable on its
  own.
- **Abstract operations and callable signatures** that extend the core without
  naming any wire protocol.
- **Tooling:** parsing, resolution and checking; rendering native types per
  language; canonical declaration identity; small runtime libraries for identity
  and presence.
- **The generation kernel**, at first, as a separate module.

**Status: chartered 26 September 2026 by
[Bitwire decision 0010](https://github.com/Bitspark/bitwire/blob/main/docs/decisions/0010-bitwire-holds-the-contract-and-bitruntime-implements-it.md);
no code yet.**

- **The language comes first.** It is a new design, not Nightseam's declaration
  language carried forward.
- **Until then**, Nightseam's generator users stay pinned to its frozen v0.6.0.

## Where it sits

```text
bittype                    depends on nothing else in the project
bitschema                  what a program validates values against; bittype emits its descriptions
Bitlink                    maps bittype operations onto the wire (paths, profile) and generates adapters
bitruntime / Bitwire       the wire contract and its implementation; bittype never names them
```

The rule that places a declaration here:

> Would it still make sense, unchanged, with a different adapter and no Bitwire?

If it would, it belongs in bittype. Paths, the profile and message
representation belong to Bitlink's binding.

## Read first

- [Charter](CHARTER.md): what this repository owns, promises and is checked by.
- [Working here as an agent](AGENTS.md).
- [The contract theory](https://github.com/Bitspark/bittheory): contracts,
  realizations, instances, and generators as functions with contracts of their own.

## Source layout

Read [LAYOUT.md](LAYOUT.md) for the component-first, two-letter language
directory convention and this repository's adoption notes.
