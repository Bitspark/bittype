package integration_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	resolution "github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	resolvegolang "github.com/Bitspark/bittype/resolve/files/golang/provider/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
)

const goCancellationChild = "BITTYPE_GO_CANCELLATION_CHILD"

// endlessSuppliedPackage is valid Go that go/types cannot finish checking. Two
// alias chains are structurally equal, and the assignment makes the checker
// compare them member by member: every step doubles the work. Depth 26 took 11
// seconds on the reviewing workstation, so depth 64 does not end.
func endlessSuppliedPackage() n.Package {
	const depth = 64
	chain := func(prefix string) string {
		var b strings.Builder
		b.WriteString("type " + prefix + "0 = struct{ x, y int }\n")
		for i := 1; i <= depth; i++ {
			b.WriteString("type " + prefix + strconv.Itoa(i) + " = struct{ a, b " + prefix + strconv.Itoa(i-1) + " }\n")
		}
		return b.String()
	}
	text := "package slow\n\n" + chain("A") + chain("B") + "var x A" + strconv.Itoa(depth) + "\nvar y B" + strconv.Itoa(depth) + " = x\n"
	return n.Package{ImportPath: "example.org/slow", Files: []n.PackageFile{{Name: "slow.go", Text: text}}}
}

// The resolution runs in a child process, because the abandoned check never
// ends and only the end of its process stops it. The parent waits for the
// child's exit. Its bound turns a call that waits for the check into a failure
// and is not a measurement.
//
// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=context-ends-during-supplied-check expects=context-error-without-waiting
func TestGoImportsResolveReturnsWhenContextEndsDuringSuppliedCheck(t *testing.T) {
	if os.Getenv(goCancellationChild) != "" {
		snapshot, _ := goPackageCapture(t, map[string]string{"a/a.go": "package a\n\ntype A string\n"}, true)
		resolver, err := resolution.NewResolver[s.Language](snapshot.Language(), n.ResolutionContract(), resolvegolang.Provider{})
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		_, err = resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: goImportSettings(t, endlessSuppliedPackage())})
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatal("resolution under an ended context returned", err)
		}
		return
	}
	bound, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	child := exec.CommandContext(bound, os.Args[0], "-test.run=^"+t.Name()+"$", "-test.count=1")
	child.Env = append(os.Environ(), goCancellationChild+"=1")
	output, err := child.CombinedOutput()
	if bound.Err() != nil {
		t.Fatalf("the resolution waited for a check that cannot end:\n%s", output)
	}
	if err != nil {
		t.Fatalf("%v\n%s", err, output)
	}
	if !strings.Contains(string(output), "PASS") {
		t.Fatalf("the child did not run the case:\n%s", output)
	}
}
