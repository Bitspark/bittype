// Package report is the composable, tooling-only provider report host. It owns
// discovery, association and rendering; consumers supply definition extractors.
package report

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	capability "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
	observation "github.com/Bitspark/bittype/tools/provider-matrix/observation/go"
)

type Binding struct {
	Format    string
	Extractor capability.Extractor
}

// Config selects one source module in the existing repository layout. Imported
// export data supplies types, not ambient dependency sources or documentation.
// Extractors are trusted Go tooling, never inspected provider implementations.
type Config struct {
	Root        string
	MatrixPath  string
	DetailsPath string
	SourceURL   string
	Extractors  []Binding
	LoadLimits  LoadLimits
	Limits      capability.Limits
	BuildFlags  []string
	Context     *observation.SourceContext
	Executions  []observation.Execution
}

// Documents contains fully rendered outputs. Generation does not modify files.
// Write publishes neither Git commits nor execution artifacts automatically.
type Documents struct {
	Matrix                        string
	Capabilities                  string
	root, matrixPath, detailsPath string
}

func Generate(ctx context.Context, config Config) (Documents, error) {
	if ctx == nil {
		return Documents{}, fmt.Errorf("nil reporting context")
	}
	root, err := filepath.Abs(config.Root)
	if err != nil {
		return Documents{}, err
	}
	config.Root = root
	if config.MatrixPath == "" {
		config.MatrixPath = "docs/providers.md"
	}
	if config.DetailsPath == "" {
		config.DetailsPath = "docs/provider-capabilities.md"
	}
	for _, path := range []string{config.MatrixPath, config.DetailsPath} {
		if filepath.IsAbs(path) || !contained(root, filepath.Join(root, path)) {
			return Documents{}, fmt.Errorf("output escapes selected module: %s", path)
		}
	}
	if filepath.Clean(config.MatrixPath) == filepath.Clean(config.DetailsPath) {
		return Documents{}, fmt.Errorf("matrix and details require different paths")
	}
	p, err := load(ctx, root, config.LoadLimits, config.BuildFlags)
	if err != nil {
		return Documents{}, err
	}
	return generate(ctx, p, config)
}

func generate(ctx context.Context, p *program, config Config) (Documents, error) {
	if config.Context != nil && config.Context.ModulePath != p.module {
		return Documents{}, fmt.Errorf("report context module does not match loaded source")
	}
	if len(config.Executions) > 0 && config.Context == nil {
		return Documents{}, fmt.Errorf("execution reports need an independently captured source context")
	}
	m, err := p.discover()
	if err != nil {
		return Documents{}, err
	}
	a, err := p.readAnnotations()
	if err != nil {
		return Documents{}, err
	}
	if err := associateCases(m, a); err != nil {
		return Documents{}, err
	}
	sourceURL := config.SourceURL
	if sourceURL == "" {
		if !strings.HasPrefix(p.module, "github.com/") {
			return Documents{}, fmt.Errorf("supply SourceURL for a non-GitHub module")
		}
		sourceURL = "https://" + p.module + "/blob/main"
	}
	detailLink, err := filepath.Rel(filepath.Dir(config.MatrixPath), config.DetailsPath)
	if err != nil {
		return Documents{}, err
	}
	providers, transformers, err := p.render(m, sourceURL, filepath.ToSlash(detailLink))
	if err != nil {
		return Documents{}, err
	}
	original, err := os.ReadFile(filepath.Join(p.root, config.MatrixPath))
	if err != nil {
		return Documents{}, err
	}
	document := strings.ReplaceAll(string(original), "\r\n", "\n")
	document, err = replaceRegion(document, "providers", providers)
	if err != nil {
		return Documents{}, err
	}
	document, err = replaceRegion(document, "transformers", transformers)
	if err != nil {
		return Documents{}, err
	}
	details, err := p.renderDetails(ctx, m, a, config, sourceURL)
	if err != nil {
		return Documents{}, err
	}
	if err := ctx.Err(); err != nil {
		return Documents{}, err
	}
	if bytes.Contains(original, []byte("\r\n")) {
		document = strings.ReplaceAll(document, "\n", "\r\n")
	}
	return Documents{Matrix: document, Capabilities: details, root: p.root, matrixPath: config.MatrixPath, detailsPath: config.DetailsPath}, nil
}

// Write writes the already rendered pair and returns every write failure. Two
// writes are not a transaction: the calling publication step must stop on error.
// The existing updater publishes the pair together in one Git commit.
func Write(documents Documents) error {
	return writeDocuments(documents, os.WriteFile)
}

func writeDocuments(documents Documents, write func(string, []byte, os.FileMode) error) error {
	if documents.root == "" || documents.matrixPath == "" || documents.detailsPath == "" {
		return fmt.Errorf("documents were not produced by Generate")
	}
	for _, item := range []struct{ path, text string }{{documents.matrixPath, documents.Matrix}, {documents.detailsPath, documents.Capabilities}} {
		path := filepath.Join(documents.root, item.path)
		if !contained(documents.root, path) {
			return fmt.Errorf("output path escapes selected module")
		}
		// Do not follow an existing output symlink or a parent-directory escape.
		realRoot, err := filepath.EvalSymlinks(documents.root)
		if err != nil {
			return err
		}
		parent, err := filepath.EvalSymlinks(filepath.Dir(path))
		if err != nil || !contained(realRoot, parent) {
			return fmt.Errorf("output directory unavailable or outside module: %s", item.path)
		}
		if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("output must not be a symlink: %s", item.path)
		}
		previous, err := os.ReadFile(path)
		if err == nil && string(previous) == item.text {
			continue
		}
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if err := write(path, []byte(item.text), 0644); err != nil {
			return err
		}
	}
	return nil
}
