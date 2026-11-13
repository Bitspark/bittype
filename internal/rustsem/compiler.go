package rustsem

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/rust/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type Output struct {
	Version int
	// Profile is the helper profile that the helper states for itself: the
	// name of the facts it writes. Check accepts only rust.HelperProfile.
	Profile     string
	Compiler    string
	Status      string
	Facts       *semantics.Data
	Diagnostics []semantic.Diagnostic
}

// The values of Output.Status that more than one place names. Readers in other
// packages compare with the three exported ones. The helper writes "ok",
// "unsupported" and "incomplete" and spells them itself, because it cannot
// import this package. It never writes "invalid": Check sets that status when
// the helper reports compiler errors.
const (
	StatusOK          = "ok"
	StatusInvalid     = "invalid"
	StatusUnsupported = "unsupported"
	statusIncomplete  = "incomplete"
)

// boundedBuffer collects what the helper process writes to one of its streams,
// up to a bound. Write is its only method. os/exec copies a stream with
// io.Copy, which calls the writer's ReadFrom when it has one, so a writer that
// embeds bytes.Buffer is filled past the bound without a call of Write. The
// first write that would pass the bound cancels the process, and every later
// write is refused too.
type boundedBuffer struct {
	data      []byte
	remaining int
	cancel    context.CancelFunc
	exceeded  bool
}

func (b *boundedBuffer) Write(data []byte) (int, error) {
	if b.exceeded || len(data) > b.remaining {
		b.exceeded = true
		b.cancel()
		return 0, semantic.ErrLimit
	}
	b.remaining -= len(data)
	b.data = append(b.data, data...)
	return len(data), nil
}
func ValidCompiler(helper, sysroot string) bool {
	return filepath.IsAbs(helper) && filepath.IsAbs(sysroot)
}
func Check(ctx context.Context, helper, sysroot string, environment Environment, trees []s.AST, source *files.Context) (Output, error) {
	if ctx == nil || !ValidCompiler(helper, sysroot) {
		return Output{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return Output{}, err
	}
	units, err := prepare(ctx, trees, source)
	if errors.Is(err, semantic.ErrLimit) {
		return Output{Status: statusIncomplete, Diagnostics: []semantic.Diagnostic{diagnostic("limit", err.Error())}}, nil
	}
	if err != nil {
		return Output{}, err
	}
	if units == nil {
		return Output{Status: "missing-context", Diagnostics: []semantic.Diagnostic{diagnostic("missing-context", "Rust checking requires captured source and AST correspondence")}}, nil
	}
	encoded, err := json.Marshal(input{Version: 1, Environment: environment, Units: units})
	if err != nil {
		return Output{}, err
	}
	if len(encoded) > 2<<20 {
		return Output{Status: statusIncomplete, Diagnostics: []semantic.Diagnostic{diagnostic("limit", "Rust compiler request exceeds the protocol bound")}}, nil
	}
	run, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	stdout := boundedBuffer{remaining: 16 << 20, cancel: cancel}
	stderr := boundedBuffer{remaining: 256 << 10, cancel: cancel}
	command := exec.CommandContext(run, helper, sysroot)
	command.Env = []string{}
	command.Stdin = bytes.NewReader(encoded)
	command.Stdout = &stdout
	command.Stderr = &stderr
	command.WaitDelay = time.Second
	err = command.Run()
	if stdout.exceeded || stderr.exceeded {
		return Output{}, semantic.ErrLimit
	}
	if ctx.Err() != nil {
		return Output{}, ctx.Err()
	}
	if run.Err() != nil {
		return Output{Status: statusIncomplete, Diagnostics: []semantic.Diagnostic{diagnostic("limit", "Rust compiler exceeded the provider deadline")}}, nil
	}
	if err != nil {
		return Output{}, fmt.Errorf("%w: %v: %s", n.ErrCompiler, err, stderr.data)
	}
	var out Output
	decoder := json.NewDecoder(bytes.NewReader(stdout.data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&out); err != nil {
		return Output{}, fmt.Errorf("%w: malformed response: %v", n.ErrCompiler, err)
	}
	var trailing any
	if decoder.Decode(&trailing) != io.EOF || out.Version != 1 || !strings.HasPrefix(out.Compiler, "rustc "+n.CompilerVersion+" (") {
		return Output{}, n.ErrCompiler
	}
	// The helper states which facts it writes. A helper of another profile,
	// or one that states none, is a helper this bridge cannot use, like one
	// that was built with another compiler: its graph is not the vocabulary
	// that the configuration of the receipt names.
	if out.Profile != n.HelperProfile {
		return Output{}, fmt.Errorf("%w: the helper states the profile %q, and this bridge reads %q", n.ErrCompiler, out.Profile, n.HelperProfile)
	}
	switch out.Status {
	case StatusOK:
		if out.Facts == nil {
			return Output{}, n.ErrCompiler
		}
	case StatusUnsupported, statusIncomplete:
		if out.Facts != nil {
			return Output{}, n.ErrCompiler
		}
	case "compiler-errors":
		if out.Facts != nil {
			return Output{}, n.ErrCompiler
		}
		out.Status = StatusInvalid
		out.Diagnostics = compilerDiagnostics(stderr.data, units)
		if len(out.Diagnostics) == 0 {
			return Output{}, fmt.Errorf("%w: %s", n.ErrCompiler, stderr.data)
		}
	case "mismatch":
		return Output{}, fmt.Errorf("%w: %v", semantic.ErrInvalidRequest, out.Diagnostics)
	case "environment":
		return Output{}, fmt.Errorf("%w: %v", n.ErrCompiler, out.Diagnostics)
	default:
		return Output{}, n.ErrCompiler
	}
	return out, nil
}
func compilerDiagnostics(data []byte, units []Unit) []semantic.Diagnostic {
	var diagnostics []semantic.Diagnostic
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		var d struct {
			Level   string
			Message string
			Code    *struct{ Code string }
			Spans   []struct {
				FileName string `json:"file_name"`
				Start    int    `json:"byte_start"`
				End      int    `json:"byte_end"`
				Primary  bool   `json:"is_primary"`
			}
		}
		if json.Unmarshal(line, &d) != nil || d.Level != "error" {
			continue
		}
		// rustc can attach missing/incompatible sysroot metadata errors to a
		// captured span. Those are environment failures, never type rejection.
		if d.Code != nil {
			switch d.Code.Code {
			case "E0460", "E0461", "E0462", "E0463", "E0514", "E0786":
				return nil
			}
		}
		for _, span := range d.Spans {
			if !span.Primary || span.FileName != "<captured.rs>" {
				continue
			}
			unit := units[0]
			var best *Node
			for i, node := range unit.Nodes {
				if node.Start <= span.Start && node.End >= span.End && (best == nil || node.End-node.Start < best.End-best.Start) {
					best = &unit.Nodes[i]
				}
			}
			if best != nil && len(diagnostics) == 0 {
				at := ast.NodeRef{Tree: unit.Tree, Node: best.ID}
				diagnostics = []semantic.Diagnostic{{Code: n.Key("compiler-error"), Subject: n.Validity(), Message: d.Message, At: &at}}
			}
		}
	}
	return diagnostics
}
