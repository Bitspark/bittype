// Command payload-schema makes or checks the pinned bootstrap artifact of a
// payload schema in a vocabulary package: the canonical module of
// payload.horos and the file that states its pin.
//
//	payload-schema -package record <directory>           writes the artifact
//	payload-schema -package record -check <directory>    fails if a file differs
//
// It is the public route from a schema source to the pinned artifact
// (docs/contribution-contract.md, O39). The library behind it is the package
// tools/payload-schema/artifact/go.
package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Bitspark/bittype/tools/payload-schema/artifact/go"
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stderr); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintln(os.Stderr, "payload-schema:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, diagnostics io.Writer) error {
	flags := flag.NewFlagSet("payload-schema", flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	pkg := flags.String("package", "", "name of the vocabulary package")
	check := flags.Bool("check", false, "compare the committed files and write nothing")
	flags.Usage = func() {
		fmt.Fprintln(diagnostics, "Usage: payload-schema -package name [-check] directory")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 || *pkg == "" {
		return errors.New("expected -package and one directory")
	}
	directory := flags.Arg(0)
	source, err := os.ReadFile(filepath.Join(directory, artifact.SourceFile))
	if err != nil {
		return err
	}
	module, pin, err := artifact.Module(ctx, string(source))
	if err != nil {
		return err
	}
	files := map[string][]byte{artifact.ModuleFile: module, artifact.PinFile: []byte(artifact.Pin(*pkg, pin))}
	for name, content := range files {
		path := filepath.Join(directory, name)
		if !*check {
			if err := os.WriteFile(path, content, 0o644); err != nil {
				return err
			}
			continue
		}
		committed, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if name != artifact.ModuleFile {
			// A checkout may have given a text file other line ends.
			committed = bytes.ReplaceAll(committed, []byte("\r\n"), []byte("\n"))
		}
		if !bytes.Equal(committed, content) {
			return fmt.Errorf("%s is not what %s gives", path, artifact.SourceFile)
		}
	}
	return nil
}
