# Horos module encoder

`horos-module` turns one ordinary Horos source module into its canonical Ontos
encoding. It uses the existing Horos parser and the same private module lowering
that supplies source pins. No module name is special to this command.

```sh
go build -o horos-module ./tools/horos-module/go
./horos-module -o module.ontos source.horos
```

The source must contain exactly one `module` section and fit the parser's 4 MiB
input limit. Imports and declaration references remain unresolved. Syntax
admission and canonical encoding do not establish semantic validity: consumers
admit the result with Horos's ordinary declaration resolver. The command never
calls semantic resolution or evaluates the installed prelude against its input.

The output is binary. Its parent directory must exist. Parsing and lowering
finish before the output file is opened, so a source error leaves an existing
output untouched. No bundle, extra identity, source map or manifest is emitted.

## Remote preparation

The manual [Horos module workflow](../../.github/workflows/horos-module.yml)
builds this tool on a GitHub runner and lowers a source file from an explicitly
selected repository revision. It returns two artifacts: `horos-module-tool`
(the Linux executable) and `canonical-module` (the file `module.ontos`).
This is generation only; it runs no test or validation workflow.

For example, select the Horos repository and the source path
`core/go/declarations/standard/horos.kinds.horos`. Place the resulting module at
`core/go/declarations/standard/horos.kinds.ontos` in that repository. The command
and workflow also accept other ordinary module sources.

## Prelude bootstrap

Build the encoder against an available Horos release before preparing its
successor. Use that executable to lower the successor's maintained source.
The successor ships the generated asset and decodes it through Ontos; ordinary
Horos builds and runtime use do not invoke Bittype.

The encoder may link the available Horos library, but it does not resolve the
input against that library's prelude. Do not rebuild the encoder through a
workspace override pointing to an incomplete successor missing its own asset.
Update Bittype's dependency after the new Horos release is available.
