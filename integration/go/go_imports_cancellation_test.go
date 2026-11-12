package integration_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
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

const (
	goCancellationChild  = "BITTYPE_GO_CANCELLATION_CHILD"
	goCancellationMarker = "abandoned: go/types was still comparing after the resolution returned"
)

// goTypesComparing matches a go/types frame that compares two types. The
// comparison of two structurally equal chains is what never ends below.
var goTypesComparing = regexp.MustCompile(`go/types\.[^\n]*identical`)

func goTypesStillComparing() bool {
	stacks := make([]byte, 4<<20)
	return goTypesComparing.Match(stacks[:runtime.Stack(stacks, true)])
}

// equalChains declares two chains of struct aliases, A0 to A<depth> and B0 to
// B<depth>, that are structurally equal. go/types compares two struct types
// member by member and remembers no pair, so comparing the two ends doubles
// the work with every step: depth 26 took 11 seconds on the reviewing
// workstation, and neither depth used here ends.
func equalChains(depth int, field string) string {
	var b strings.Builder
	for _, prefix := range []string{"A", "B"} {
		b.WriteString("type " + prefix + "0 = struct {\n\tX " + field + "\n\tY " + field + "\n}\n")
		for i := 1; i <= depth; i++ {
			previous := prefix + strconv.Itoa(i-1)
			b.WriteString("type " + prefix + strconv.Itoa(i) + " = struct {\n\tA " + previous + "\n\tB " + previous + "\n}\n")
		}
	}
	return b.String()
}

// goCancellationCases holds the one input whose first comparison never ends: a
// supplied package that assigns the end of one chain to the end of the other.
//
// Captured code reaches the same cost in the export, which compares every pair
// of exported types: the same chains as captured declarations took 0.6 seconds
// at depth 20 and 8.8 seconds at depth 24 in the base profile. There the cost
// grows comparison by comparison and the work stops between two of them once
// its context has ended, so no moment is known at which a comparison that
// cannot end is running. The export runs inside the same unit as the checks,
// which TestBesideReturnsWhenContextEnds and this case pin at its two ends.
var goCancellationCases = map[string]func(t *testing.T) (semantic.Snapshot, []semantic.Setting){
	"supplied package check": func(t *testing.T) (semantic.Snapshot, []semantic.Setting) {
		const depth = 64
		text := "package slow\n\n" + equalChains(depth, "int") + "var x A" + strconv.Itoa(depth) + "\nvar y B" + strconv.Itoa(depth) + " = x\n"
		snapshot, _ := goPackageCapture(t, map[string]string{"a/a.go": "package a\n\ntype A string\n"}, true)
		return snapshot, goImportSettings(t, n.Package{ImportPath: "example.org/slow", Files: []n.PackageFile{{Name: "slow.go", Text: text}}})
	},
}

// Each case runs in a child process, because the abandoned work never ends and
// only the end of its process stops it. The child ends the context once it has
// seen go/types comparing, expects the resolution back, and then sees go/types
// still comparing: the resolution returned without waiting. No step depends on
// a duration. The bounds turn a resolution that waits into a failure.
//
// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=context-ends-during-native-work expects=context-error-without-waiting
func TestGoResolveReturnsWhenContextEndsDuringNativeWork(t *testing.T) {
	if name := os.Getenv(goCancellationChild); name != "" {
		goCancellationInChild(t, name)
		return
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for name := range goCancellationCases {
		t.Run(name, func(t *testing.T) {
			bound, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			child := exec.CommandContext(bound, executable, "-test.run=^TestGoResolveReturnsWhenContextEndsDuringNativeWork$", "-test.count=1", "-test.v")
			child.Env = append(os.Environ(), goCancellationChild+"="+name)
			output, err := child.CombinedOutput()
			if bound.Err() != nil {
				t.Fatalf("the child did not end:\n%s", output)
			}
			if err != nil {
				t.Fatalf("%v\n%s", err, output)
			}
			if !strings.Contains(string(output), goCancellationMarker) {
				t.Fatalf("the child did not reach the abandon:\n%s", output)
			}
		})
	}
}

func goCancellationInChild(t *testing.T, name string) {
	build, ok := goCancellationCases[name]
	if !ok {
		t.Fatal("unknown case", name)
	}
	snapshot, settings := build(t)
	resolver, err := resolution.NewResolver[s.Language](snapshot.Language(), n.ResolutionContract(), resolvegolang.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	returned := make(chan error, 1)
	go func() {
		_, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: settings})
		returned <- err
	}()
	patience := time.After(time.Minute)
	for !goTypesStillComparing() {
		select {
		case err := <-returned:
			t.Fatal("the resolution ended before go/types was seen comparing:", err)
		case <-patience:
			t.Fatal("go/types was never seen comparing")
		default:
			runtime.Gosched()
		}
	}
	cancel()
	select {
	case err := <-returned:
		if !errors.Is(err, context.Canceled) {
			t.Fatal("the resolution under an ended context returned", err)
		}
	case <-time.After(time.Minute):
		t.Fatal("the resolution waited for native work that cannot end")
	}
	if !goTypesStillComparing() {
		t.Fatal("go/types was not comparing after the resolution returned")
	}
	fmt.Println(goCancellationMarker)
}
