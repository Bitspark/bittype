package rustsem

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	parser "github.com/Bitspark/bittype/parse/syntax/rust/provider/go"
	n "github.com/Bitspark/bittype/semantics/rust/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
)

func compilerFixture(t *testing.T) ([]s.AST, *files.Context) {
	t.Helper()
	text := "pub type Byte = u8;"
	parsed, err := parser.Parse(context.Background(), parse.Input{Text: text, SourceID: "test", TreeID: "test"})
	if err != nil {
		t.Fatal(err)
	}
	return []s.AST{parsed.ASTView()}, &files.Context{Files: files.Files{Entries: []files.File{{ID: "test", Path: "test.rs", Text: text}}}, Links: []files.SyntaxLink{{File: "test", Tree: "test", Source: "test"}}}
}

var environment = Environment{Edition: n.Edition2024, Target: n.TargetLinuxX64, CrateName: "pages", CrateKind: n.KindLibrary}

func TestCompilerFailureAndProtocol(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("qualified helper execution platform is Linux")
	}
	trees, source := compilerFixture(t)
	ctx := context.Background()
	sysroot := t.TempDir()
	// Every response states the helper profile, so that each is refused for
	// what it shows and not for a missing profile.
	profile := `"Profile":"` + n.HelperProfile + `",`
	for _, response := range []string{`{}`, `{"Version":1,` + profile + `"Compiler":"rustc 0.0.0 (x)","Status":"ok","Facts":{}}`, `{"Version":1,` + profile + `"Compiler":"rustc 1.90.0 (x)","Status":"ok"}`, `{"Version":1,` + profile + `"Compiler":"rustc 1.90.0 (x)","Status":"unsupported","Facts":{}}`, `{"Version":1,` + profile + `"Compiler":"rustc 1.90.0 (x)","Status":"ok","Facts":{},"unexpected":true}`} {
		path := filepath.Join(t.TempDir(), "helper")
		if err := os.WriteFile(path, []byte("#!/bin/sh\ncat >/dev/null\nprintf '%s' '"+response+"'\n"), 0700); err != nil {
			t.Fatal(err)
		}
		if _, err := Check(ctx, path, sysroot, environment, trees, source); !errors.Is(err, n.ErrCompiler) {
			t.Fatal(response, err)
		}
	}
	if _, err := Check(ctx, filepath.Join(t.TempDir(), "missing"), sysroot, environment, trees, source); !errors.Is(err, n.ErrCompiler) {
		t.Fatal(err)
	}
	if _, err := Check(ctx, os.Getenv("BITTYPE_RUST_HELPER"), sysroot, environment, trees, source); !errors.Is(err, n.ErrCompiler) {
		t.Fatal("wrong sysroot", err)
	}
	t.Setenv("BITTYPE_SHOULD_NOT_LEAK", "secret-placeholder")
	path := filepath.Join(t.TempDir(), "helper")
	script := "#!/bin/sh\ncat >/dev/null\ntest -z \"$BITTYPE_SHOULD_NOT_LEAK\" || exit 1\nprintf '%s' '{\"Version\":1," + profile + "\"Compiler\":\"rustc 1.90.0 (x)\",\"Status\":\"incomplete\"}'\n"
	if err := os.WriteFile(path, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	if out, err := Check(ctx, path, sysroot, environment, trees, source); err != nil || out.Status != "incomplete" {
		t.Fatal(out, err)
	}
}

// The bridge takes a response only under the helper profile that the
// vocabulary publishes. A helper binary of another change of the facts states
// another profile, and one of before this rule states none. Each is a helper
// the bridge cannot use, and the refusal is the one for a helper that was
// built with another compiler. The same response under the published profile
// is taken.
func TestBridgeRefusesAnotherHelperProfile(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("qualified helper execution platform is Linux")
	}
	trees, source := compilerFixture(t)
	respond := func(member string) (Output, error) {
		path := filepath.Join(t.TempDir(), "helper")
		script := "#!/bin/sh\ncat >/dev/null\nprintf '%s' '{\"Version\":1," + member + "\"Compiler\":\"rustc " + n.CompilerVersion + " (x)\",\"Status\":\"incomplete\"}'\n"
		if err := os.WriteFile(path, []byte(script), 0700); err != nil {
			t.Fatal(err)
		}
		return Check(context.Background(), path, t.TempDir(), environment, trees, source)
	}
	if out, err := respond(`"Profile":"` + n.HelperProfile + `",`); err != nil || out.Status != statusIncomplete || out.Profile != n.HelperProfile {
		t.Fatal("a response under the published profile", out, err)
	}
	for name, c := range map[string]struct{ member, stated string }{
		"no profile, as the helper before this rule":     {``, `""`},
		"the profile before the primitive type bindings": {`"Profile":"captured-types-2",`, `"captured-types-2"`},
		"a later profile":                   {`"Profile":"captured-types-4",`, `"captured-types-4"`},
		"an empty profile":                  {`"Profile":"",`, `""`},
		"the profile with a space after it": {`"Profile":"` + n.HelperProfile + ` ",`, `"` + n.HelperProfile + ` "`},
	} {
		out, err := respond(c.member)
		if !errors.Is(err, n.ErrCompiler) || errors.Is(err, semantic.ErrInvalidRequest) || errors.Is(err, n.ErrUnsupported) || out.Status != "" {
			t.Errorf("%s: status %q and the error %v, want the refusal of a helper", name, out.Status, err)
			continue
		}
		// The refusal names what the helper stated and what the bridge reads.
		if !strings.Contains(err.Error(), "states the profile "+c.stated) || !strings.Contains(err.Error(), `"`+n.HelperProfile+`"`) {
			t.Errorf("%s: the refusal does not name the two profiles: %v", name, err)
		}
	}
}

// The helper states the published profile in every response: with a graph,
// for a source that the compiler rejects, and for a request it does not
// cover. This reads what the helper writes, before the bridge judges it, so
// the name is pinned from the helper's side too.
func TestHelperStatesItsProfile(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("qualified helper execution platform is Linux")
	}
	for _, c := range []struct {
		name, status string
		units        []Unit
	}{
		{"a source with a graph", "ok", helperUnits(t, "pub type Flag = bool;")},
		{"a source that the compiler rejects", "compiler-errors", helperUnits(t, "pub type Missing = Unknown;")},
		{"a request without a unit", "unsupported", []Unit{}},
	} {
		response := helperResponse(t, c.units)
		if string(response["Status"]) != `"`+c.status+`"` {
			t.Errorf("%s: the status %s, want %s", c.name, response["Status"], c.status)
		}
		if string(response["Profile"]) != `"`+n.HelperProfile+`"` {
			t.Errorf("%s: the helper states the profile %s, and the vocabulary publishes %q", c.name, response["Profile"], n.HelperProfile)
		}
	}
}

// helperUnits is the one unit that the bridge sends the helper for a text.
func helperUnits(t *testing.T, text string) []Unit {
	t.Helper()
	parsed, err := parser.Parse(context.Background(), parse.Input{Text: text, SourceID: "test", TreeID: "test"})
	if err != nil {
		t.Fatal(err)
	}
	source := &files.Context{Files: files.Files{Entries: []files.File{{ID: "test", Path: "test.rs", Text: text}}}, Links: []files.SyntaxLink{{File: "test", Tree: "test", Source: "test"}}}
	units, err := prepare(context.Background(), []s.AST{parsed.ASTView()}, source)
	if err != nil {
		t.Fatal(err)
	}
	return units
}

// helperResponse runs the helper on units as the bridge does and returns the
// members of its response as it wrote them.
func helperResponse(t *testing.T, units []Unit) map[string]json.RawMessage {
	t.Helper()
	encoded, err := json.Marshal(input{Version: 1, Environment: environment, Units: units})
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	command := exec.Command(os.Getenv("BITTYPE_RUST_HELPER"), os.Getenv("BITTYPE_RUST_SYSROOT"))
	command.Env = []string{}
	command.Stdin = bytes.NewReader(encoded)
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil {
		t.Fatal(err, stderr.String())
	}
	var response map[string]json.RawMessage
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatal(err, stdout.String())
	}
	return response
}

// The helper binds every primitive type node or publishes no graph. The parser
// gives the helper no primitive type node that the compiler does not resolve,
// so this sends one through the private protocol: a node of that kind over
// the text of a keyword. A graph with such a node and no binding at it would
// hold a use whose capture nobody can see.
func TestHelperRefusesUnboundPrimitiveType(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("qualified helper execution platform is Linux")
	}
	units := helperUnits(t, "pub type Flag = bool;")
	if response := helperResponse(t, units); string(response["Status"]) != `"ok"` {
		t.Fatal("the unit as the bridge sends it", string(response["Status"]), string(response["Diagnostics"]))
	}
	units[0].Nodes = append(units[0].Nodes, Node{ID: units[0].Nodes[len(units[0].Nodes)-1].ID + 1, Kind: s.KindPrimitiveType, Start: 0, End: len("pub")})
	response := helperResponse(t, units)
	if string(response["Status"]) != `"unsupported"` || string(response["Facts"]) != "null" || !strings.Contains(string(response["Diagnostics"]), "unmapped primitive type") {
		t.Errorf("a primitive type node without a binding: status %s, facts %.80s, diagnostics %s", response["Status"], response["Facts"], response["Diagnostics"])
	}
}

// The helper passes the environment of the request to the compiler and judges
// no value of it. Which values a profile supports is decided in one place, the
// Rust vocabulary, and the helper holds no second copy of them.
func TestHelperObeysEnvironment(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("qualified helper execution platform is Linux")
	}
	check := func(text string, change func(*Environment)) (Output, error) {
		parsed, err := parser.Parse(context.Background(), parse.Input{Text: text, SourceID: "test", TreeID: "test"})
		if err != nil {
			t.Fatal(err)
		}
		source := &files.Context{Files: files.Files{Entries: []files.File{{ID: "test", Path: "test.rs", Text: text}}}, Links: []files.SyntaxLink{{File: "test", Tree: "test", Source: "test"}}}
		stated := environment
		change(&stated)
		return Check(context.Background(), os.Getenv("BITTYPE_RUST_HELPER"), os.Getenv("BITTYPE_RUST_SYSROOT"), stated, []s.AST{parsed.ASTView()}, source)
	}
	// A path that starts with :: names the crate root in edition 2015 and
	// another crate since edition 2018.
	rooted := "pub struct Local {} pub type Alias = ::Local;"
	if out, err := check(rooted, func(e *Environment) { e.Edition = "2015" }); err != nil || out.Status != StatusOK {
		t.Error("edition 2015 was not passed on", out.Status, out.Diagnostics, err)
	}
	if out, err := check(rooted, func(*Environment) {}); err != nil || out.Status != StatusInvalid {
		t.Error("edition 2024", out.Status, out.Diagnostics, err)
	}
	// The pinned sysroot holds the standard library of one target.
	if _, err := check("pub type Byte = u8;", func(e *Environment) { e.Target = "aarch64-unknown-linux-gnu" }); !errors.Is(err, n.ErrCompiler) {
		t.Error("another target was not passed on", err)
	}
	// A binary crate needs a main function, which a library does not.
	if out, err := check("pub type Byte = u8;", func(e *Environment) { e.CrateKind = "bin" }); err == nil && out.Status == StatusOK {
		t.Error("the crate kind was not passed on")
	}
	if out, err := check("pub type Byte = u8;", func(e *Environment) { e.CrateName = "other" }); err != nil || out.Status != StatusOK || out.Facts.Entities[0].Name != "other::Byte" {
		t.Error("the crate name was not passed on", out, err)
	}
}

// The helper forces one lint and withholds every warning from its error
// stream. The stream carries nothing for a source that the compiler accepts,
// and the compiler's errors for one that it rejects. This runs the helper as
// the bridge does and reads the stream, which Check does not return.
func TestHelperWithholdsWarnings(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("qualified helper execution platform is Linux")
	}
	run := func(text string) (string, string) {
		t.Helper()
		parsed, err := parser.Parse(context.Background(), parse.Input{Text: text, SourceID: "test", TreeID: "test"})
		if err != nil {
			t.Fatal(err)
		}
		source := &files.Context{Files: files.Files{Entries: []files.File{{ID: "test", Path: "test.rs", Text: text}}}, Links: []files.SyntaxLink{{File: "test", Tree: "test", Source: "test"}}}
		units, err := prepare(context.Background(), []s.AST{parsed.ASTView()}, source)
		if err != nil {
			t.Fatal(err)
		}
		encoded, err := json.Marshal(input{Version: 1, Environment: environment, Units: units})
		if err != nil {
			t.Fatal(err)
		}
		var stdout, stderr bytes.Buffer
		command := exec.Command(os.Getenv("BITTYPE_RUST_HELPER"), os.Getenv("BITTYPE_RUST_SYSROOT"))
		command.Env = []string{}
		command.Stdin = bytes.NewReader(encoded)
		command.Stdout, command.Stderr = &stdout, &stderr
		if err := command.Run(); err != nil {
			t.Fatal(err, stderr.String())
		}
		return stdout.String(), stderr.String()
	}
	// The forced lint twice, and the compiler's count of the warnings.
	stdout, stderr := run("pub struct Public { pub field: Hidden, pub other: Hidden } struct Hidden {}")
	if stderr != "" || !strings.Contains(stdout, `"Status":"ok"`) || strings.Count(stdout, `"Name":"private-interface"`) != 2 {
		t.Errorf("a source with the forced lint: the error stream holds %d bytes: %.300s\n%.300s", len(stderr), stderr, stdout)
	}
	stdout, stderr = run("pub type Missing = Unknown;")
	if !strings.Contains(stderr, `"level":"error"`) || strings.Contains(stderr, `"level":"warning"`) || !strings.Contains(stdout, `"Status":"compiler-errors"`) {
		t.Errorf("a source with a compiler error: %.300s\n%.300s", stderr, stdout)
	}
}
func TestCaptureLimitsAndOffsets(t *testing.T) {
	trees, source := compilerFixture(t)
	source.Files.Entries[0].Text = "\xff"
	if _, err := prepare(context.Background(), trees, source); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal(err)
	}
	text := "// café\r\npub type Byte = u8;\r\n"
	parsed, err := parser.Parse(context.Background(), parse.Input{Text: text, TreeID: "test", SourceID: "test"})
	if err != nil {
		t.Fatal(err)
	}
	source.Files.Entries[0].Text = text
	units, err := prepare(context.Background(), []s.AST{parsed.ASTView()}, source)
	if err != nil {
		t.Fatal(err)
	}
	if units[0].Text != "// café\npub type Byte = u8;\n" {
		t.Fatal(units[0].Text)
	}
	ctx, cancel := context.WithCancel(context.Background())
	b := boundedBuffer{remaining: 3, cancel: cancel}
	if _, err := b.Write([]byte("four")); !errors.Is(err, semantic.ErrLimit) || !b.exceeded || ctx.Err() != context.Canceled {
		t.Fatal("unbounded output")
	}
}
