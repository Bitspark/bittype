package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
	"slices"
	"strings"
	"testing"
)

// Root-module architecture is qualified by its own public gate. The optional
// composition has a separate explicit dependency boundary, not an exemption.
func TestVerificationArchitecture(t *testing.T){
	command:=exec.Command("go","list","-deps","-json","./...");command.Dir="../.."
	output,err:=command.Output();if err!=nil{t.Fatal(err)}
	type pkg struct{ImportPath string;Standard bool;Imports []string}
	packages:=map[string]pkg{};decoder:=json.NewDecoder(bytes.NewReader(output))
	for{var p pkg;err:=decoder.Decode(&p);if err==io.EOF{break};if err!=nil{t.Fatal(err)};packages[p.ImportPath]=p}
	const root="github.com/Bitspark/bittype/"
	const module=root+"verification/"
	rules:=map[string][]string{
		"files/source/horos/go":{root+"contracts/semantic/go",root+"contracts/resolution/go",root+"internal/horosmodel",root+"resolve/files/horos/go",root+"syntax/horos/go",root+"trees/ast/go",root+"trees/representation/go","github.com/Bitspark/horos/core/go/declarations","github.com/bitspark/ontos/core/go","github.com/bitspark/ontos/data/go"},
		"files/target/golang/go":{root+"contracts/semantic/go",root+"semantics/go",root+"syntax/golang/go",root+"trees/ast/go",root+"trees/representation/go"},
		"files/horos-to-go/go":{root+"contracts/semantic/go",root+"contracts/verification/go",root+"semantics/go",root+"trees/ast/go",module+"files/source/horos/go",module+"files/target/golang/go","github.com/Bitspark/horos/core/go/declarations","github.com/bitspark/ontos/core/go","github.com/bitspark/ontos/data/go","github.com/Bitspark/logos/contract/go/typed"},
	}
	for path,p:=range packages{
		if p.Standard || !strings.HasPrefix(path,module) || path==module+"integration/go"{continue}
		allowed,ok:=rules[strings.TrimPrefix(path,module)];if !ok{t.Errorf("unclassified verification production package: %s",path);continue}
		for _,dependency:=range p.Imports{if !packages[dependency].Standard && !slices.Contains(allowed,dependency){t.Errorf("%s imports %s",path,dependency)}}
	}
}
