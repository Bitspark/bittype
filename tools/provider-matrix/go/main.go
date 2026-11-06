// Command provider-matrix composes source discovery, supplied extractors and
// report publication. Inspected provider code is never executed.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	grammar "github.com/Bitspark/bittype/tools/provider-matrix/extract/grammar/go"
	observation "github.com/Bitspark/bittype/tools/provider-matrix/observation/go"
	report "github.com/Bitspark/bittype/tools/provider-matrix/report/go"
)

type values []string

func (v *values) String() string     { return fmt.Sprint([]string(*v)) }
func (v *values) Set(s string) error { *v = append(*v, s); return nil }

func run() error {
	var resultPaths, buildFlags values
	root := flag.String("root", ".", "selected source module root")
	output := flag.String("output", "docs/providers.md", "matrix Markdown path relative to root")
	details := flag.String("details", "docs/provider-capabilities.md", "details Markdown path relative to root")
	sourceURL := flag.String("source-url", "", "source URL prefix; defaults to the module GitHub main branch")
	contextPath := flag.String("context", "", "independently captured current source context JSON")
	executionOutput := flag.String("execution-output", "", "write an execution report here, leaving committed reports unchanged")
	timeout := flag.Duration("timeout", 5*time.Minute, "whole-command cooperative deadline (outer worker deadline remains required)")
	flag.Var(&resultPaths, "results", "invocation.json path (repeatable)")
	flag.Var(&buildFlags, "build-flag", "explicit source-selection Go build flag (repeatable)")
	flag.Parse()
	if flag.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	if *timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	config := report.Config{Root: *root, MatrixPath: *output, DetailsPath: *details, SourceURL: *sourceURL, BuildFlags: buildFlags, Extractors: []report.Binding{{Format: grammar.Format, Extractor: grammar.Extractor{}}}}
	if *contextPath != "" {
		file, err := os.Open(*contextPath)
		if err != nil {
			return err
		}
		defer file.Close()
		decoder := json.NewDecoder(io.LimitReader(file, 1<<20))
		decoder.DisallowUnknownFields()
		var source observation.SourceContext
		if err := decoder.Decode(&source); err != nil {
			return err
		}
		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			return fmt.Errorf("trailing or oversized source context")
		}
		config.Context = &source
		if len(config.BuildFlags) > 0 {
			return fmt.Errorf("context supplies build flags; do not also pass -build-flag")
		}
		config.BuildFlags = source.BuildFlags
		if config.SourceURL == "" && source.Revision != "" {
			config.SourceURL = "https://" + source.ModulePath + "/blob/" + source.Revision
		}
	}
	for _, path := range resultPaths {
		execution, err := observation.Load(ctx, path, observation.Limits{})
		if err != nil {
			return err
		}
		config.Executions = append(config.Executions, execution)
	}
	if (len(resultPaths) > 0 || config.Context != nil) && *executionOutput == "" {
		return fmt.Errorf("observations require -execution-output; they are never written into committed reports")
	}
	documents, err := report.Generate(ctx, config)
	if err != nil {
		return err
	}
	if *executionOutput != "" {
		if err := os.WriteFile(*executionOutput, []byte(documents.Capabilities), 0644); err != nil {
			return err
		}
		fmt.Println("Wrote invocation report", *executionOutput)
		return nil
	}
	if err := report.Write(documents); err != nil {
		return err
	}
	fmt.Println("Generated provider matrices and capability details.")
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
