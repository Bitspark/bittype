package integration_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Throwaway probe for row 24 of research 0016. Removed by the next commit.
// It asks the pinned compiler which lint flag makes the private-interface
// diagnostic appear beside "-A warnings", and at which level.
func TestProbeRow24LintFlags(t *testing.T) {
	rustc := filepath.Join(os.Getenv("BITTYPE_RUST_SYSROOT"), "bin", "rustc")
	if _, err := os.Stat(rustc); err != nil {
		t.Skip("no pinned compiler")
	}
	source := filepath.Join(t.TempDir(), "probe.rs")
	if err := os.WriteFile(source, []byte("pub struct Public { pub field: Hidden } struct Hidden {}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	version, _ := exec.Command(rustc, "--version").CombinedOutput()
	fmt.Println("### compiler:", strings.TrimSpace(string(version)))
	for _, flags := range [][]string{
		{"-A", "warnings"},
		{"-A", "warnings", "-W", "private_interfaces"},
		{"-W", "private_interfaces", "-A", "warnings"},
		{"-A", "warnings", "--force-warn", "private_interfaces"},
		{"--force-warn", "private_interfaces", "-A", "warnings"},
		{"-A", "warnings", "-D", "private_interfaces"},
	} {
		args := append([]string{"--crate-name", "captured", "--crate-type", "lib", "--edition", "2024", "--emit", "metadata", "--error-format", "short", "-o", filepath.Join(t.TempDir(), "out.rmeta"), source}, flags...)
		out, err := exec.Command(rustc, args...).CombinedOutput()
		fmt.Printf("### flags: %s\n  exit error: %v\n", strings.Join(flags, " "), err)
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if at := strings.LastIndex(line, "probe.rs"); at >= 0 {
				line = line[at:]
			}
			if line != "" {
				fmt.Println("  " + line)
			}
		}
	}
}
