package integration_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"
)

// These fixtures qualify a proposed Go encoding directly with the Go compiler.
// They do not invoke Bittype generation, Horos analysis, or a Logos checker.
func TestGoNativeAdmission(t *testing.T) {
	root := t.TempDir()
	fixtures := filepath.Join("testdata", "go-admission")
	write := func(t *testing.T, name string, data []byte) {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	copyFixture := func(name string) {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(fixtures, name))
		if err != nil {
			t.Fatal(err)
		}
		write(t, name, data)
	}
	write(t, "go.mod", []byte("module admission.example\n\ngo 1.25.0\n"))
	for _, name := range []string{"model/model.go", "admission/admission.go", "foreign/foreign.go"} {
		copyFixture(name)
	}

	goBinary := filepath.Join(runtime.GOROOT(), "bin", "go")
	if runtime.GOOS == "windows" {
		goBinary += ".exe"
	}
	run := func(t *testing.T, args ...string) ([]byte, error) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx, goBinary, args...)
		cmd.Dir = root
		// The fixture has only local imports. Do not fetch dependencies or change
		// toolchains; use the exact compiler running this test.
		cmd.Env = append(os.Environ(),
			"GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off", "GOSUMDB=off", "GOFLAGS=",
		)
		output, err := cmd.CombinedOutput()
		if ctx.Err() != nil {
			t.Fatalf("Go compiler timed out: %v\n%s", args, output)
		}
		return output, err
	}
	version, err := run(t, "version")
	if err != nil {
		t.Fatalf("compiler version: %v\n%s", err, version)
	}
	t.Logf("compiler: %s; fixture language: go1.25.0", strings.TrimSpace(string(version)))
	if output, err := run(t, "build", "./model", "./admission", "./foreign"); err != nil {
		t.Fatalf("fixture dependencies must compile before testing consumers: %v\n%s", err, output)
	}

	for _, category := range []string{"positive", "negative"} {
		paths, err := filepath.Glob(filepath.Join(fixtures, category, "*.go"))
		if err != nil || len(paths) == 0 {
			t.Fatalf("%s fixtures: %v", category, err)
		}
		for _, path := range paths {
			name := strings.TrimSuffix(filepath.Base(path), ".go")
			t.Run(category+"/"+name, func(t *testing.T) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				directory := filepath.Join("cases", category, name)
				write(t, filepath.Join(directory, "consumer.go"), data)
				output, err := run(t, "build", "./"+filepath.ToSlash(directory))
				if category == "positive" {
					if err != nil {
						t.Fatalf("valid consumer rejected: %v\n%s", err, output)
					}
					return
				}
				if err == nil {
					t.Fatalf("invalid consumer unexpectedly compiled:\n%s", data)
				}
				var exit *exec.ExitError
				if !errors.As(err, &exit) {
					t.Fatalf("compiler did not produce a diagnostic: %v\n%s", err, output)
				}
				const prefix = "// want-error: "
				var expected string
				for _, line := range strings.Split(string(data), "\n") {
					if strings.HasPrefix(line, prefix) {
						if expected != "" {
							t.Fatal("fixture must declare exactly one expected diagnostic")
						}
						expected = strings.TrimSpace(strings.TrimPrefix(line, prefix))
					}
				}
				if expected == "" {
					t.Fatal("negative fixture lacks an expected diagnostic")
				}
				pattern, err := regexp.Compile(expected)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(string(output), "consumer.go:") || !pattern.Match(output) {
					t.Fatalf("wanted consumer diagnostic matching %q; got:\n%s", expected, output)
				}
				t.Logf("expected rejection:\n%s", output)
			})
		}
	}
}
