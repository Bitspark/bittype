# Working here as an agent

Read the [charter](CHARTER.md), Bitwire's
[decision 0010](https://github.com/Bitspark/bitwire/blob/main/docs/decisions/0010-bitwire-holds-the-contract-and-bitruntime-implements-it.md)
and the [contract theory](https://github.com/Bitspark/bittheory) before changing
this tree.

- Keep the language wire-independent. Paths, profiles and message representation
  are Bitlink's. Test each declaration with: "would it make sense, unchanged, with a
  different adapter and no Bitwire?"
- Depend on nothing else in the project. The generation kernel module may not
  import the declaration language.
- Keep these distinct:
  - declaration identity;
  - profile revision;
  - generated-code compatibility;
  - behavior the declaration does not represent.

  A digest names declared content, not behavior.
- Hold every identity library to the shared canonicalization vectors. Two
  implementations of one encoding in one language are a defect.
- This is a new language, designed rather than migrated. Nightseam's language and
  the family design are evidence to examine, not specifications to copy.
- After the initial bootstrap, use a branch or worktree and a pull request.
  Squash a green change onto `main`.
- Do not add private checkout dependencies, local orchestration state or
  credentials.
