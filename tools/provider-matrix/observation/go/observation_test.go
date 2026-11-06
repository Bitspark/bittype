package observation

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func knownContext() SourceContext {
	return SourceContext{
		Revision: "revision-a", ModulePath: "example.test/reports", ModuleRevision: "revision-a", Harness: "repository",
		GoVersion: "go1.27.3", GOOS: "linux", GOARCH: "amd64", CGOEnabled: "1", InputState: "clean", Complete: true,
		Modules: []ModuleSelection{{Path: "example.test/reports", Revision: "revision-a"}, {Path: "example.test/dependency", Version: "v1.0.0"}},
	}
}

func invocation() Invocation {
	zero := 0
	return Invocation{Schema: Schema, Context: knownContext(), RunID: "17", Attempt: "2", Job: "contracts", CommandID: "contracts", Args: []string{"go", "test", "-json", "-count=1", "./..."}, ExitCode: &zero, CaptureComplete: true}
}

func declaration(name string) TestDeclaration {
	return TestDeclaration{ID: "decl/" + name, TestedPackage: "example.test/reports", Name: name, Variant: "external", File: "reports_test.go"}
}

func decoded(t *testing.T, capture Invocation, input string) Execution {
	t.Helper()
	execution, err := Decode(context.Background(), capture, strings.NewReader(input), strings.NewReader(""), Limits{})
	if err != nil {
		t.Fatal(err)
	}
	return execution
}

func TestNamedOutcomesPreserveInvocationAndSubtests(t *testing.T) {
	capture := invocation()
	failure := 1
	capture.ExitCode = &failure
	capture.Args = []string{"go", "test", "-json", "-count=1", "-run=TestParent/selected|TestFails", "./..."}
	execution := decoded(t, capture, `
{"Action":"start","Package":"example.test/reports"}
{"Action":"run","Package":"example.test/reports","Test":"TestParent"}
{"Action":"run","Package":"example.test/reports","Test":"TestParent/selected"}
{"Action":"skip","Package":"example.test/reports","Test":"TestParent/selected"}
{"Action":"pass","Package":"example.test/reports","Test":"TestParent"}
{"Action":"run","Package":"example.test/reports","Test":"TestFails"}
{"Action":"fail","Package":"example.test/reports","Test":"TestFails"}
{"Action":"fail","Package":"example.test/reports"}
`)
	assessment := Associate(Snapshot{Context: knownContext(), Declarations: []TestDeclaration{declaration("TestParent"), declaration("TestFails")}}, execution)
	if len(assessment.Observations) != 3 || assessment.Observations[0].Status != Matched || assessment.Observations[0].Run.Outcome != Passed || assessment.Observations[1].Status != Unassociated || assessment.Observations[1].Run.Outcome != Skipped || assessment.Observations[2].Run.Outcome != Failed {
		t.Fatal(assessment)
	}
	if execution.Packages[0].Outcome != Failed || *execution.Invocation.ExitCode != 1 || execution.Invocation.Args[4] != "-run=TestParent/selected|TestFails" {
		t.Fatal("invocation failure or subtest selection lost", execution)
	}
}

func TestRepeatedExecutionsNeverOverwriteFailure(t *testing.T) {
	capture := invocation()
	capture.Args = []string{"go", "test", "-json", "-count=2", "./..."}
	execution := decoded(t, capture, `
{"Action":"run","Package":"example.test/reports","Test":"TestAgain"}
{"Action":"fail","Package":"example.test/reports","Test":"TestAgain"}
{"Action":"run","Package":"example.test/reports","Test":"TestAgain"}
{"Action":"pass","Package":"example.test/reports","Test":"TestAgain"}
`)
	if len(execution.Tests) != 2 || execution.Tests[0].Occurrence != 1 || execution.Tests[1].Occurrence != 2 || execution.Tests[0].Outcome != Failed || execution.Tests[1].Outcome != Passed {
		t.Fatal(execution.Tests)
	}
}

func TestBuildEventsAndPackagePassDoNotInventNamedEvidence(t *testing.T) {
	execution, err := Decode(context.Background(), invocation(), strings.NewReader(`
{"Action":"build-output","ImportPath":"example.test/broken [example.test/broken.test]","Output":"compile error\n"}
{"Action":"output","Package":"example.test/reports","Output":"--- PASS: TestInvented (0.00s)\n"}
{"Action":"pass","Package":"example.test/reports"}
{"Action":"build-fail","ImportPath":"example.test/broken [example.test/broken.test]"}
{"Action":"fail","Package":"example.test/broken","FailedBuild":"example.test/broken [example.test/broken.test]"}
`), strings.NewReader("non-JSON early stderr\n"), Limits{})
	if err != nil || len(execution.Tests) != 0 || len(execution.Builds) != 2 || len(execution.Packages) != 2 || execution.Stderr != "non-JSON early stderr\n" {
		t.Fatal(execution, err)
	}
	if execution.Builds[1].ImportPath != execution.Packages[1].FailedBuild {
		t.Fatal("build identity was confused with tested package", execution)
	}
	if got := Associate(Snapshot{Context: knownContext(), Declarations: []TestDeclaration{declaration("TestInvented")}}, execution); len(got.Observations) != 0 {
		t.Fatal(got)
	}
}

func TestIncompleteAndMalformedExecutionCannotMatch(t *testing.T) {
	for _, input := range []string{
		`{"Action":"run","Package":"example.test/reports","Test":"TestOne"}`,
		`{"Action":"pass","Package":"example.test/reports","Test":"TestOne"}`,
		"{\"Action\":\"run\",\"Package\":\"example.test/reports\",\"Test\":\"TestOne\"}\n{",
		`{"Action":"start","Package":"example.test/reports"}{"Action":"run","Package":"example.test/reports","Test":"TestOne"}{"Action":"pass","Package":"example.test/reports","Test":"TestOne"}`,
	} {
		execution, _ := Decode(context.Background(), invocation(), strings.NewReader(input), strings.NewReader(""), Limits{})
		assessment := Associate(Snapshot{Context: knownContext(), Declarations: []TestDeclaration{declaration("TestOne")}}, execution)
		if len(assessment.Observations) != 1 || assessment.Observations[0].Status != Insufficient || assessment.Observations[0].DeclarationID != "" {
			t.Fatal(assessment)
		}
	}
	for _, input := range []string{
		`{"Action":"future-action","Package":"example.test/reports"}`,
		`{"Action":"build-fail","Package":"example.test/reports"}`,
		`{"Action":"pass","Test":"TestOne"}`,
		`{"Action":"run","Package":"example.test/reports","Test":"TestOne"}{"Action":"run","Package":"example.test/reports","Test":"TestOne"}`,
	} {
		execution, err := Decode(context.Background(), invocation(), strings.NewReader(input), strings.NewReader(""), Limits{})
		if !errors.Is(err, ErrProtocol) || execution.ProtocolComplete {
			t.Fatal(execution, err)
		}
	}
}

func passing(t *testing.T) Execution {
	return decoded(t, invocation(), `{"Action":"run","Package":"example.test/reports","Test":"TestOne"}
{"Action":"pass","Package":"example.test/reports","Test":"TestOne"}`)
}

func TestSourceAndBuildSelectionMustMatchIndependently(t *testing.T) {
	cases := []struct {
		name   string
		change func(*Execution)
		want   AssociationStatus
	}{
		{"revision", func(e *Execution) { e.Invocation.Context.Revision = "revision-b" }, Mismatched},
		{"module", func(e *Execution) { e.Invocation.Context.ModuleRevision = "revision-b" }, Mismatched},
		{"harness", func(e *Execution) { e.Invocation.Context.Harness = "installed.consumer" }, Mismatched},
		{"race", func(e *Execution) { e.Invocation.Context.BuildFlags = []string{"-race"} }, Mismatched},
		{"toolchain", func(e *Execution) { e.Invocation.Context.GoVersion = "go1.28.0" }, Mismatched},
		{"dependency", func(e *Execution) { e.Invocation.Context.Modules[1].Version = "v2.0.0" }, Mismatched},
		{"replacement", func(e *Execution) { e.Invocation.Context.Replacements = []Replacement{{Path: "example.test/dependency", Version: "v1.0.0", Replacement: "elsewhere", WithVersion: "v1.0.0"}} }, Mismatched},
		{"modified", func(e *Execution) { e.Invocation.Context.InputState = "modified"; e.Invocation.Context.InputChanges = []string{" M go.sum"} }, Mismatched},
		{"unknown", func(e *Execution) { e.Invocation.Context.InputState = "unknown" }, Insufficient},
		{"context-incomplete", func(e *Execution) { e.Invocation.Context.Complete = false }, Insufficient},
		{"capture-incomplete", func(e *Execution) { e.Invocation.CaptureComplete = false }, Insufficient},
		{"exit-absent", func(e *Execution) { e.Invocation.ExitCode = nil }, Insufficient},
		{"invocation-absent", func(e *Execution) { e.Invocation.RunID = "" }, Insufficient},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			execution := passing(t)
			tc.change(&execution)
			assessment := Associate(Snapshot{Context: knownContext(), Declarations: []TestDeclaration{declaration("TestOne")}}, execution)
			if assessment.Status != tc.want || len(assessment.Observations) != 1 || assessment.Observations[0].Status != tc.want || assessment.Observations[0].DeclarationID != "" {
				t.Fatal(assessment)
			}
		})
	}
	context := knownContext()
	slices.Reverse(context.Modules)
	if assessment := Associate(Snapshot{Context: context, Declarations: []TestDeclaration{declaration("TestOne")}}, passing(t)); assessment.Status != Matched || assessment.Observations[0].DeclarationID != "decl/TestOne" {
		t.Fatal("module ordering changed identity", assessment)
	}
}

func TestAmbiguousInternalAndExternalNamesRemainUnassigned(t *testing.T) {
	internal := declaration("TestOne")
	internal.ID, internal.Variant, internal.File = "internal/TestOne", "internal", "internal_test.go"
	execution := decoded(t, invocation(), `
{"Action":"run","Package":"example.test/reports","Test":"TestOne"}
{"Action":"pass","Package":"example.test/reports","Test":"TestOne"}
{"Action":"run","Package":"example.test/reports","Test":"TestOne#01"}
{"Action":"pass","Package":"example.test/reports","Test":"TestOne#01"}
`)
	assessment := Associate(Snapshot{Context: knownContext(), Declarations: []TestDeclaration{declaration("TestOne"), internal}}, execution)
	if assessment.Observations[0].Status != Ambiguous || len(assessment.Observations[0].Candidates) != 2 || assessment.Observations[0].DeclarationID != "" || assessment.Observations[1].Status != Unassociated {
		t.Fatal("guessed test ownership from event order or suffix", assessment)
	}
}

func TestInvocationCapturesStaySeparate(t *testing.T) {
	ordinary := passing(t)
	ordinary.Tests[0].Outcome = Skipped
	native := passing(t)
	native.Invocation.Job = "native"
	native.Invocation.Environment = []Setting{{Name: "BITTYPE_NATIVE", Value: "1"}}
	snapshot := Snapshot{Context: knownContext(), Declarations: []TestDeclaration{declaration("TestOne")}}
	first, second := Associate(snapshot, ordinary), Associate(snapshot, native)
	if first.Observations[0].Run.Outcome != Skipped || second.Observations[0].Run.Outcome != Passed || ordinary.Invocation.Job == native.Invocation.Job {
		t.Fatal(first, second)
	}
}

func TestObservationBoundsAndCancellation(t *testing.T) {
	for _, tc := range []struct {
		stdout, stderr string
		limits         Limits
	}{
		{"{}", "", Limits{StdoutBytes: 1}},
		{"", "xx", Limits{StderrBytes: 1}},
		{`{"Action":"output","Package":"p"}{"Action":"output","Package":"p"}`, "", Limits{Events: 1}},
	} {
		result, err := Decode(context.Background(), invocation(), strings.NewReader(tc.stdout), strings.NewReader(tc.stderr), tc.limits)
		if !errors.Is(err, ErrLimit) || result.ProtocolComplete {
			t.Fatal(result, err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result, err := Decode(ctx, invocation(), strings.NewReader(""), strings.NewReader(""), Limits{})
	if !errors.Is(err, context.Canceled) || result.ProtocolComplete {
		t.Fatal(result, err)
	}
}

func TestLoadUsesFixedSiblingArtifacts(t *testing.T) {
	dir := t.TempDir()
	data, err := json.Marshal(invocation())
	if err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string][]byte{"invocation.json": data, "stdout.jsonl": {}, "stderr.log": []byte("early failure\n")} {
		if err := os.WriteFile(filepath.Join(dir, name), content, 0600); err != nil {
			t.Fatal(err)
		}
	}
	execution, err := Load(context.Background(), filepath.Join(dir, "invocation.json"), Limits{})
	if err != nil || execution.Stderr != "early failure\n" {
		t.Fatal(execution, err)
	}
	data = append(data[:len(data)-1], []byte(`,"stdout":"elsewhere"}`)...)
	if err := os.WriteFile(filepath.Join(dir, "invocation.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(context.Background(), filepath.Join(dir, "invocation.json"), Limits{}); !errors.Is(err, ErrInvocation) {
		t.Fatal("capture JSON selected another file", err)
	}
}

func TestCaptureHelperPreservesFailureAndIndependentContext(t *testing.T) {
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Fatal("CI capture qualification requires Python 3", err)
	}
	_, filename, _, _ := runtime.Caller(0)
	helper := filepath.Join(filepath.Dir(filename), "..", "..", "capture.py")
	root, artifacts := t.TempDir(), t.TempDir()
	for name, body := range map[string]string{
		"go.mod": "module example.test/capture\n\ngo 1.25.0\n",
		"capture_test.go": "package capture\nimport \"testing\"\nfunc TestPass(t *testing.T) {}\nfunc TestFails(t *testing.T) { t.Fatal(\"fixture failure\") }\nfunc TestSkipped(t *testing.T) { t.Skip(\"fixture skip\") }\n",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
	}
	for _, args := range [][]string{{"init", "-q"}, {"add", "go.mod", "capture_test.go"}, {"-c", "user.name=Capture Test", "-c", "user.email=capture@example.invalid", "-c", "commit.gpgsign=false", "commit", "-qm", "capture fixture"}} {
		command := exec.Command("git", args...)
		command.Dir = root
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatal(string(output), err)
		}
	}
	dir := filepath.Join(artifacts, "invocation")
	command := exec.Command(python, helper, "run", "--root", root, "--output", dir, "--command-id", "fixture", "--", "go", "test", "-json", "-count=1", "./...")
	output, err := command.CombinedOutput()
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 1 {
		t.Fatal("capture lost original test failure", string(output), err)
	}
	execution, err := Load(context.Background(), filepath.Join(dir, "invocation.json"), Limits{})
	if err != nil || !execution.Invocation.CaptureComplete || execution.Invocation.ExitCode == nil || *execution.Invocation.ExitCode != 1 || !execution.Invocation.Context.Complete || execution.Invocation.Context.InputState != "clean" {
		t.Fatal(execution, string(output), err)
	}
	contextPath := filepath.Join(artifacts, "context.json")
	command = exec.Command(python, helper, "context", "--root", root, "--output", contextPath)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatal(string(output), err)
	}
	data, err := os.ReadFile(contextPath)
	if err != nil {
		t.Fatal(err)
	}
	var current SourceContext
	if err := json.Unmarshal(data, &current); err != nil {
		t.Fatal(err)
	}
	var declarations []TestDeclaration
	for _, name := range []string{"TestPass", "TestFails", "TestSkipped"} {
		declarations = append(declarations, TestDeclaration{ID: name, TestedPackage: "example.test/capture", Name: name, Variant: "internal", File: "capture_test.go"})
	}
	assessment := Associate(Snapshot{Context: current, Declarations: declarations}, execution)
	if assessment.Status != Matched || len(assessment.Observations) != 3 {
		t.Fatal(assessment)
	}
	want := []Outcome{Passed, Failed, Skipped}
	for i, observation := range assessment.Observations {
		if observation.Status != Matched || observation.Run.Outcome != want[i] {
			t.Fatal(assessment)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "untracked.go"), []byte("package capture\n"), 0600); err != nil {
		t.Fatal(err)
	}
	command = exec.Command(python, helper, "context", "--root", root, "--output", contextPath)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatal(string(output), err)
	}
	data, err = os.ReadFile(contextPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &current); err != nil {
		t.Fatal(err)
	}
	if assessment := Associate(Snapshot{Context: current, Declarations: declarations}, execution); assessment.Status != Mismatched {
		t.Fatal("changed source was credited to checkout revision", assessment)
	}
}
