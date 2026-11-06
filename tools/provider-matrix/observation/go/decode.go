package observation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrInvocation = errors.New("invalid test invocation capture")
	ErrProtocol   = errors.New("invalid Go test event stream")
	ErrLimit      = errors.New("test observation limit exceeded")
)

func (l Limits) bounded() (Limits, error) {
	if l.StdoutBytes < 0 || l.StderrBytes < 0 || l.Events < 0 || l.StdoutBytes > 1<<30 || l.StderrBytes > 1<<30 || l.Events > 10000000 {
		return l, ErrLimit
	}
	if l.StdoutBytes == 0 {
		l.StdoutBytes = 64 << 20
	}
	if l.StderrBytes == 0 {
		l.StderrBytes = 16 << 20
	}
	if l.Events == 0 {
		l.Events = 200000
	}
	return l, nil
}

// Load reads the versioned invocation file and its fixed raw-output siblings.
// Imported JSON cannot select additional paths. Raw artifacts remain useful if
// decoding fails; the partial Execution never gains ProtocolComplete.
func Load(ctx context.Context, path string, limits Limits) (Execution, error) {
	var invocation Invocation
	f, err := os.Open(path)
	if err != nil {
		return Execution{}, err
	}
	defer f.Close()
	limited := &io.LimitedReader{R: f, N: (1 << 20) + 1}
	decoder := json.NewDecoder(limited)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&invocation); err != nil {
		return Execution{}, fmt.Errorf("%w: %v", ErrInvocation, err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return Execution{}, fmt.Errorf("%w: trailing invocation content", ErrInvocation)
	}
	if limited.N <= 0 {
		return Execution{}, ErrLimit
	}
	stdout, err := os.Open(filepath.Join(filepath.Dir(path), "stdout.jsonl"))
	if err != nil {
		return Execution{Invocation: invocation}, err
	}
	defer stdout.Close()
	stderr, err := os.Open(filepath.Join(filepath.Dir(path), "stderr.log"))
	if err != nil {
		return Execution{Invocation: invocation}, err
	}
	defer stderr.Close()
	return Decode(ctx, invocation, stdout, stderr, limits)
}

type event struct {
	Action      string
	Package     string
	Test        string
	Elapsed     float64
	Output      string
	ImportPath  string
	FailedBuild string
}

type testKey struct {
	pkg, name string
}

// Decode consumes the documented Go test/build JSON protocol, never Output text.
// It retains distinct repetitions and enclosing failures. Cancellation is checked
// between bounded reads; this is not a process-isolation or hard-deadline API.
func Decode(ctx context.Context, invocation Invocation, stdout, stderr io.Reader, limits Limits) (Execution, error) {
	result := Execution{Invocation: invocation}
	fail := func(code string, at int, err error) (Execution, error) {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Code: code, Message: err.Error(), Event: at})
		return result, err
	}
	if invocation.Schema != Schema || invocation.CommandID == "" || len(invocation.Args) == 0 || stdout == nil || stderr == nil {
		return fail("invocation", 0, ErrInvocation)
	}
	limits, err := limits.bounded()
	if err != nil {
		return fail("limit", 0, err)
	}
	if err := ctx.Err(); err != nil {
		return fail("canceled", 0, err)
	}
	data, err := io.ReadAll(io.LimitReader(stderr, limits.StderrBytes+1))
	if int64(len(data)) > limits.StderrBytes {
		result.Stderr = string(data[:limits.StderrBytes])
		return fail("limit", 0, ErrLimit)
	}
	result.Stderr = string(data)
	if err != nil {
		return fail("stderr", 0, err)
	}
	reader := &io.LimitedReader{R: stdout, N: limits.StdoutBytes + 1}
	decoder := json.NewDecoder(reader)
	activeTests := map[testKey]int{}
	counts := map[testKey]int{}
	activePackages := map[string]int{}
	packageCounts := map[string]int{}
	for at := 1; ; at++ {
		if err := ctx.Err(); err != nil {
			return fail("canceled", at, err)
		}
		var e event
		err := decoder.Decode(&e)
		if reader.N <= 0 {
			return fail("limit", at, ErrLimit)
		}
		if err == io.EOF {
			if len(activeTests) != 0 || len(activePackages) != 0 {
				return fail("incomplete", at, fmt.Errorf("%w: stream ended with unfinished executions", ErrProtocol))
			}
			result.ProtocolComplete = true
			return result, nil
		}
		if err != nil {
			return fail("protocol", at, fmt.Errorf("%w at event %d: %v", ErrProtocol, at, err))
		}
		if at > limits.Events {
			return fail("limit", at, ErrLimit)
		}
		if e.Elapsed < 0 {
			return fail("protocol", at, fmt.Errorf("%w: negative elapsed time", ErrProtocol))
		}
		if e.Action == "build-output" || e.Action == "build-fail" {
			if e.ImportPath == "" || e.Test != "" || e.Package != "" {
				return fail("protocol", at, fmt.Errorf("%w: invalid build-event identity", ErrProtocol))
			}
			result.Builds = append(result.Builds, BuildEvent{ImportPath: e.ImportPath, Action: e.Action, Output: e.Output})
			continue
		}
		if e.Package == "" {
			return fail("protocol", at, fmt.Errorf("%w: absent tested package", ErrProtocol))
		}
		key := testKey{e.Package, e.Test}
		switch e.Action {
		case "start":
			if e.Test != "" {
				return fail("protocol", at, fmt.Errorf("%w: named package start", ErrProtocol))
			}
			if _, found := activePackages[e.Package]; found {
				return fail("protocol", at, fmt.Errorf("%w: overlapping package executions", ErrProtocol))
			}
			packageCounts[e.Package]++
			activePackages[e.Package] = len(result.Packages)
			result.Packages = append(result.Packages, PackageRun{Package: e.Package, Occurrence: packageCounts[e.Package], Started: true, Outcome: Incomplete})
		case "run":
			if e.Test == "" {
				return fail("protocol", at, fmt.Errorf("%w: unnamed test run", ErrProtocol))
			}
			if _, found := activeTests[key]; found {
				return fail("protocol", at, fmt.Errorf("%w: overlapping identical test names", ErrProtocol))
			}
			counts[key]++
			activeTests[key] = len(result.Tests)
			result.Tests = append(result.Tests, TestRun{Package: e.Package, Name: e.Test, Occurrence: counts[key], Started: true, Outcome: Incomplete})
		case "pass", "fail", "skip":
			if e.Test == "" {
				i, found := activePackages[e.Package]
				if !found {
					packageCounts[e.Package]++
					i = len(result.Packages)
					result.Packages = append(result.Packages, PackageRun{Package: e.Package, Occurrence: packageCounts[e.Package]})
				}
				result.Packages[i].Outcome = Outcome(e.Action)
				result.Packages[i].Elapsed = e.Elapsed
				result.Packages[i].FailedBuild = e.FailedBuild
				delete(activePackages, e.Package)
			} else {
				i, found := activeTests[key]
				if !found {
					counts[key]++
					i = len(result.Tests)
					result.Tests = append(result.Tests, TestRun{Package: e.Package, Name: e.Test, Occurrence: counts[key]})
				}
				result.Tests[i].Outcome = Outcome(e.Action)
				result.Tests[i].Elapsed = e.Elapsed
				delete(activeTests, key)
			}
		case "pause", "cont":
			if _, found := activeTests[key]; e.Test == "" || !found {
				return fail("protocol", at, fmt.Errorf("%w: pause/continue without a running test", ErrProtocol))
			}
		case "output", "bench":
			// Raw output is retained in stdout.jsonl. It cannot establish outcomes.
		default:
			return fail("protocol", at, fmt.Errorf("%w: unknown action %q", ErrProtocol, e.Action))
		}
	}
}
