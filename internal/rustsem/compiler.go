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
	Version     int
	Compiler    string
	Status      string
	Facts       *semantics.Data
	Diagnostics []semantic.Diagnostic
}
type boundedBuffer struct {
	bytes.Buffer
	remaining int
	cancel    context.CancelFunc
	exceeded  bool
}

func (b *boundedBuffer) Write(data []byte) (int, error) {
	if len(data) > b.remaining {
		b.exceeded = true
		b.cancel()
		return 0, semantic.ErrLimit
	}
	b.remaining -= len(data)
	return b.Buffer.Write(data)
}
func ValidCompiler(helper, sysroot string) bool {
	return filepath.IsAbs(helper) && filepath.IsAbs(sysroot)
}
func Check(ctx context.Context, helper, sysroot string, trees []s.AST, source *files.Context) (Output, error) {
	if ctx == nil || !ValidCompiler(helper, sysroot) {
		return Output{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return Output{}, err
	}
	units, err := prepare(ctx, trees, source)
	if errors.Is(err, semantic.ErrLimit) {
		return Output{Status: "incomplete", Diagnostics: []semantic.Diagnostic{diagnostic("limit", err.Error())}}, nil
	}
	if err != nil {
		return Output{}, err
	}
	if units == nil {
		return Output{Status: "missing-context", Diagnostics: []semantic.Diagnostic{diagnostic("missing-context", "Rust checking requires captured source and AST correspondence")}}, nil
	}
	encoded, err := json.Marshal(input{Version: 1, Units: units})
	if err != nil {
		return Output{}, err
	}
	if len(encoded) > 2<<20 {
		return Output{Status: "incomplete", Diagnostics: []semantic.Diagnostic{diagnostic("limit", "Rust compiler request exceeds the protocol bound")}}, nil
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
		return Output{Status: "incomplete", Diagnostics: []semantic.Diagnostic{diagnostic("limit", "Rust compiler exceeded the provider deadline")}}, nil
	}
	if err != nil {
		return Output{}, fmt.Errorf("%w: %v: %s", n.ErrCompiler, err, stderr.String())
	}
	var out Output
	decoder := json.NewDecoder(bytes.NewReader(stdout.Bytes()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&out); err != nil {
		return Output{}, fmt.Errorf("%w: malformed response: %v", n.ErrCompiler, err)
	}
	var trailing any
	if decoder.Decode(&trailing) != io.EOF || out.Version != 1 || !strings.HasPrefix(out.Compiler, "rustc "+n.CompilerVersion+" (") {
		return Output{}, n.ErrCompiler
	}
	switch out.Status {
	case "ok":
		if out.Facts == nil {
			return Output{}, n.ErrCompiler
		}
	case "unsupported", "incomplete":
		if out.Facts != nil {
			return Output{}, n.ErrCompiler
		}
	case "compiler-errors":
		if out.Facts != nil {
			return Output{}, n.ErrCompiler
		}
		out.Status = "invalid"
		out.Diagnostics = compilerDiagnostics(stderr.Bytes(), units)
		if len(out.Diagnostics) == 0 {
			return Output{}, fmt.Errorf("%w: %s", n.ErrCompiler, stderr.String())
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
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		var d struct {
			Level   string
			Message string
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
			if best != nil {
				at := ast.NodeRef{Tree: unit.Tree, Node: best.ID}
				return []semantic.Diagnostic{{Code: n.Key("compiler-error"), Subject: n.Validity(), Message: d.Message, At: &at}}
			}
		}
	}
	return nil
}
