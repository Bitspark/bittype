// Command horos-module lowers one Horos source module to its canonical Ontos
// encoding. It does not resolve imports or establish semantic validity.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	parsecontract "github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/internal/horosmodel"
	"github.com/Bitspark/bittype/internal/lex"
	horosparse "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	codec "github.com/bitspark/ontos/codec/go"
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stderr); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintln(os.Stderr, "horos-module:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, diagnostics io.Writer) error {
	flags := flag.NewFlagSet("horos-module", flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	output := flags.String("o", "", "output canonical Ontos module file")
	flags.Usage = func() {
		fmt.Fprintln(diagnostics, "Usage: horos-module -o module.ontos source.horos")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 || *output == "" {
		return errors.New("expected one source file and -o output file")
	}
	filename := flags.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	source, err := io.ReadAll(io.LimitReader(file, int64(lex.MaxBytes)+1))
	if err != nil {
		return fmt.Errorf("read source: %w", err)
	}
	if len(source) > lex.MaxBytes {
		return fmt.Errorf("source exceeds %d bytes", lex.MaxBytes)
	}
	tree, err := horosparse.Parse(ctx, parsecontract.Input{
		SourceID: ast.SourceID(filename),
		TreeID:   "module",
		Text:     string(source),
	})
	if err != nil {
		return err
	}
	if len(tree.Abstract.Nodes[tree.Abstract.Root-1].Children) != 1 {
		return errors.New("expected exactly one module section")
	}
	modules, err := horosmodel.SourceModules(tree.Abstract)
	if err != nil {
		return err
	}
	if err := os.WriteFile(*output, codec.Encode(modules[0].Value), 0o644); err != nil {
		return fmt.Errorf("write canonical module: %w", err)
	}
	return nil
}
