package golangsem

import (
	"context"
	"errors"
	goast "go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"
	"time"
)

// heldImporter stops the checker inside an import until the test releases it.
type heldImporter struct{ entered, release chan struct{} }

func (h heldImporter) Import(string) (*types.Package, error) {
	close(h.entered)
	<-h.release
	return nil, errors.New("released")
}

type panickingImporter struct{}

func (panickingImporter) Import(string) (*types.Package, error) { panic("importer panic") }

func importingFile(t *testing.T, fset *token.FileSet) []*goast.File {
	t.Helper()
	file, err := parser.ParseFile(fset, "a.go", "package a\n\nimport \"held\"\n\ntype A = held.T\n", parser.SkipObjectResolution)
	if err != nil {
		t.Fatal(err)
	}
	return []*goast.File{file}
}

// The checker is provably still running, held inside an import, when the
// context ends. The call returns the context's error without waiting for it.
func TestCheckPackageReturnsWhenContextEnds(t *testing.T) {
	held := heldImporter{entered: make(chan struct{}), release: make(chan struct{})}
	defer close(held.release) // lets the abandoned check end
	fset := token.NewFileSet()
	files := importingFile(t, fset)
	conf := types.Config{Importer: held}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	returned := make(chan error, 1)
	go func() {
		_, err := checkPackage(ctx, &conf, "a", fset, files, nil)
		returned <- err
	}()
	// The bounds only turn a hang into a failure; nothing is inferred from time.
	select {
	case <-held.entered:
	case <-time.After(time.Minute):
		t.Fatal("the checker never reached the import")
	}
	cancel()
	select {
	case err := <-returned:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Minute):
		t.Fatal("the call waited for a check that cannot end")
	}
}

// A check that ends is returned as it is, type errors included.
func TestCheckPackageReturnsTheCheckersResult(t *testing.T) {
	fset := token.NewFileSet()
	good, err := parser.ParseFile(fset, "a.go", "package a\n\ntype A string\n", parser.SkipObjectResolution)
	if err != nil {
		t.Fatal(err)
	}
	conf := types.Config{}
	pkg, err := checkPackage(context.Background(), &conf, "a", fset, []*goast.File{good}, nil)
	if err != nil || pkg.Scope().Lookup("A") == nil {
		t.Fatal(pkg, err)
	}
	bad, err := parser.ParseFile(fset, "b.go", "package b\n\ntype B Missing\n", parser.SkipObjectResolution)
	if err != nil {
		t.Fatal(err)
	}
	var typeError types.Error
	if _, err := checkPackage(context.Background(), &conf, "b", fset, []*goast.File{bad}, nil); !errors.As(err, &typeError) {
		t.Fatal(err)
	}
}

// A panic inside the checker surfaces in the caller's goroutine, where it did
// before the check ran beside the call.
func TestCheckPackagePanicReachesTheCaller(t *testing.T) {
	fset := token.NewFileSet()
	files := importingFile(t, fset)
	conf := types.Config{Importer: panickingImporter{}}
	var recovered any
	func() {
		defer func() { recovered = recover() }()
		_, _ = checkPackage(context.Background(), &conf, "a", fset, files, nil)
	}()
	if recovered != "importer panic" {
		t.Fatal("recovered", recovered)
	}
}
