package observation

import (
	"reflect"
	"slices"
	"strings"
)

func sourceProblems(c SourceContext) []string {
	var problems []string
	if !c.Complete {
		problems = append(problems, "source context capture is incomplete")
	}
	if c.Revision == "" || c.ModulePath == "" || c.ModuleRevision == "" || c.Harness == "" {
		problems = append(problems, "source, module, or harness identity is missing")
	}
	if c.GoVersion == "" || c.GOOS == "" || c.GOARCH == "" || (c.CGOEnabled != "0" && c.CGOEnabled != "1") {
		problems = append(problems, "Go build selection is incomplete")
	}
	if c.InputState != "clean" && c.InputState != "modified" {
		problems = append(problems, "source input state is unknown")
	}
	mainFound := false
	seen := map[string]bool{}
	for _, module := range c.Modules {
		if module.Path == "" || seen[module.Path] {
			problems = append(problems, "module selection is ambiguous")
		}
		seen[module.Path] = true
		if module.Path == c.ModulePath {
			mainFound = true
		}
	}
	if !mainFound {
		problems = append(problems, "selected main module is missing")
	}
	replacements := map[string]bool{}
	for _, replacement := range c.Replacements {
		key := replacement.Path + "\x00" + replacement.Version
		if replacement.Path == "" || replacement.Replacement == "" || replacements[key] {
			problems = append(problems, "module replacement identity is ambiguous")
		}
		replacements[key] = true
		if replacement.WithVersion == "" {
			problems = append(problems, "local replacement inputs are not qualified")
		}
	}
	return append(problems, c.Problems...)
}

func normalized(c SourceContext) SourceContext {
	c.BuildFlags = slices.Clone(c.BuildFlags)
	c.Modules = slices.Clone(c.Modules)
	c.Replacements = slices.Clone(c.Replacements)
	slices.SortFunc(c.Modules, func(a, b ModuleSelection) int { return strings.Compare(a.Path, b.Path) })
	slices.SortFunc(c.Replacements, func(a, b Replacement) int {
		if a.Path != b.Path {
			return strings.Compare(a.Path, b.Path)
		}
		return strings.Compare(a.Version, b.Version)
	})
	if len(c.BuildFlags) == 0 {
		c.BuildFlags = nil
	}
	if len(c.Replacements) == 0 {
		c.Replacements = nil
	}
	c.InputChanges = nil
	c.Problems = nil
	return c
}

// Associate joins observations only to exact, unambiguous top-level declarations
// under an independently supplied matching source/build context. Equal profile
// labels, a parent pass, or a package pass cannot supply a missing observation.
func Associate(snapshot Snapshot, execution Execution) Assessment {
	result := Assessment{Status: Matched}
	for _, item := range []struct {
		label   string
		context SourceContext
	}{{"declaration", snapshot.Context}, {"invocation", execution.Invocation.Context}} {
		for _, reason := range sourceProblems(item.context) {
			result.Reasons = append(result.Reasons, item.label+": "+reason)
		}
	}
	if !execution.ProtocolComplete || !execution.Invocation.CaptureComplete || execution.Invocation.ExitCode == nil || len(execution.Diagnostics) != 0 {
		result.Reasons = append(result.Reasons, "execution capture is incomplete")
	}
	if execution.Invocation.Schema != Schema || execution.Invocation.CommandID == "" || len(execution.Invocation.Args) == 0 {
		result.Reasons = append(result.Reasons, "invocation capture metadata is invalid")
	}
	if execution.Invocation.RunID == "" || execution.Invocation.Attempt == "" || execution.Invocation.Job == "" {
		result.Reasons = append(result.Reasons, "invocation identity is incomplete")
	}
	if len(result.Reasons) != 0 {
		result.Status = Insufficient
	} else if snapshot.Context.InputState != "clean" || execution.Invocation.Context.InputState != "clean" || len(snapshot.Context.InputChanges) != 0 || len(execution.Invocation.Context.InputChanges) != 0 {
		result.Status = Mismatched
		result.Reasons = append(result.Reasons, "modified source inputs do not establish the recorded revision")
	} else if !reflect.DeepEqual(normalized(snapshot.Context), normalized(execution.Invocation.Context)) {
		result.Status = Mismatched
		result.Reasons = append(result.Reasons, "source revision, module/harness, or build selection differs")
	}
	index := map[testKey][]TestDeclaration{}
	ids := map[string]bool{}
	for _, declaration := range snapshot.Declarations {
		if declaration.ID == "" || ids[declaration.ID] {
			result.Status = Insufficient
			result.Reasons = append(result.Reasons, "test declaration identities are missing or duplicated")
		}
		ids[declaration.ID] = true
		key := testKey{declaration.TestedPackage, declaration.Name}
		index[key] = append(index[key], declaration)
	}
	for _, run := range execution.Tests {
		association := Association{Run: run, Status: result.Status}
		if result.Status != Matched {
			association.Reason = strings.Join(result.Reasons, "; ")
		} else if strings.Contains(run.Name, "/") {
			association.Status = Unassociated
			association.Reason = "subtest observation; cases are declared on top-level test functions"
		} else if candidates := index[testKey{run.Package, run.Name}]; len(candidates) > 1 {
			association.Status = Ambiguous
			association.Reason = "multiple source declarations share this tested-package/name pair"
			for _, candidate := range candidates {
				association.Candidates = append(association.Candidates, candidate.ID)
			}
			slices.Sort(association.Candidates)
		} else if len(candidates) == 0 {
			association.Status = Unassociated
			association.Reason = "no exact top-level test declaration was supplied"
		} else if candidates[0].ID == "" || candidates[0].File == "" || (candidates[0].Variant != "internal" && candidates[0].Variant != "external") {
			association.Status = Insufficient
			association.Reason = "test declaration identity is incomplete"
		} else if !run.Started || run.Occurrence < 1 || (run.Outcome != Passed && run.Outcome != Failed && run.Outcome != Skipped) {
			association.Status = Insufficient
			association.Reason = "test execution is incomplete"
		} else {
			association.DeclarationID = candidates[0].ID
		}
		result.Observations = append(result.Observations, association)
	}
	return result
}
