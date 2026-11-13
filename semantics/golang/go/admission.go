package golang

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/Bitspark/bittype/contracts/semantic/go"
)

// The two functions below state which file names and which package names the
// first Go Files slice admits. They are published here because the consumer
// that names files and packages in a placement and the Go Files formatter that
// refuses them must agree on the rule. They judge names only: they read no
// file and start no process.
//
// A nil result admits the name. An error that wraps ErrUnsupported refuses a
// name that Go itself accepts and that this slice does not cover; a caller
// reports it as unsupported coverage. An error that wraps
// semantic.ErrInvalidRequest refuses what is no such name at all.
//
// A later slice may admit a name that is refused here. It may not give an
// admitted name another meaning.

// AdmitFileName answers whether the first Go Files slice admits the base name
// of filePath. A build selects a Go file by its name: a suffix such as
// "_linux" or "_amd64" is a build constraint, "_test" makes a test file, and
// the go tool ignores a name that starts with "_" or ".". The rule refuses
// more than those cases, so it needs no list of systems and architectures: a
// base name with an underscore, or one that starts with a dot, is unsupported.
// The directory is not judged.
func AdmitFileName(filePath string) error {
	base := filePath[strings.LastIndex(filePath, "/")+1:]
	switch {
	case !strings.HasSuffix(base, ".go") || strings.ContainsAny(base, "\\\x00"):
		return fmt.Errorf("%w: %q is not the path of a Go file", semantic.ErrInvalidRequest, filePath)
	case strings.HasPrefix(base, "."):
		return fmt.Errorf("%w: the go tool ignores the file name %q, which starts with a dot", ErrUnsupported, base)
	case strings.Contains(base, "_"):
		return fmt.Errorf("%w: the file name %q has an underscore and may carry a build constraint, name a test file or be ignored", ErrUnsupported, base)
	}
	return nil
}

// AdmitPackageName answers whether the first Go Files slice admits a package
// with this name. Ordinary packages are admitted. A command, which is package
// main, and a test package, whose name ends in "_test", are unsupported.
func AdmitPackageName(name string) error {
	switch {
	case !token.IsIdentifier(name) || name == "_":
		return fmt.Errorf("%w: %q is not a Go package name", semantic.ErrInvalidRequest, name)
	case name == "main":
		return fmt.Errorf("%w: package main is a command", ErrUnsupported)
	case strings.HasSuffix(name, "_test"):
		return fmt.Errorf("%w: package %s is a test package", ErrUnsupported, name)
	}
	return nil
}
