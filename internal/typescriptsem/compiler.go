package typescriptsem

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
)

//go:embed compiler.cjs
var helper string

type Output struct { Version int; Status string; Facts *semantics.Data; Diagnostics []semantic.Diagnostic; Assignable *bool }
type boundedBuffer struct { bytes.Buffer; remaining int; cancel context.CancelFunc; exceeded bool }
func (b *boundedBuffer) Write(data []byte) (int,error) {
	if len(data)>b.remaining { b.exceeded=true; b.cancel(); return 0,semantic.ErrLimit }; b.remaining-=len(data); return b.Buffer.Write(data)
}
func ValidCompiler(nodePath,modulePath string) bool { return filepath.IsAbs(nodePath) && filepath.IsAbs(modulePath) }
func Check(ctx context.Context,nodePath,modulePath string,trees []s.AST,source *files.Context,query *Query) (Output,error) {
	if ctx==nil || !ValidCompiler(nodePath,modulePath) { return Output{},semantic.ErrInvalidRequest }
	if err:=ctx.Err(); err!=nil { return Output{},err }
	units,err:=prepare(ctx,trees,source)
	if errors.Is(err,semantic.ErrLimit) { return Output{Status:"incomplete",Diagnostics:[]semantic.Diagnostic{{Code:n.Key("limit"),Subject:n.Coverage(),Message:err.Error()}}},nil }
	if err!=nil { return Output{},err }
	if units==nil { return Output{Status:"missing-context",Diagnostics:[]semantic.Diagnostic{{Code:n.Key("missing-context"),Subject:n.Coverage(),Message:"TypeScript checking requires captured source files and their AST correspondence"}}},nil }
	encoded,err:=json.Marshal(input{Version:1,Units:units,Query:query}); if err!=nil { return Output{},err }
	run,cancel:=context.WithTimeout(ctx,30*time.Second); defer cancel()
	stdout:=boundedBuffer{remaining:16<<20,cancel:cancel}; stderr:=boundedBuffer{remaining:16<<10,cancel:cancel}
	command:=exec.CommandContext(run,nodePath,"--input-type=commonjs","--eval",helper,modulePath)
	// No NODE_OPTIONS, NODE_PATH, package-manager credentials or ambient hooks.
	command.Env=[]string{}
	for _,key:=range []string{"SystemRoot","WINDIR","TEMP","TMP"} { if value,ok:=os.LookupEnv(key); ok { command.Env=append(command.Env,key+"="+value) } }
	command.Stdin=bytes.NewReader(encoded); command.Stdout=&stdout; command.Stderr=&stderr; command.WaitDelay=time.Second
	err=command.Run()
	if stdout.exceeded || stderr.exceeded { return Output{},semantic.ErrLimit }
	if ctx.Err()!=nil { return Output{},ctx.Err() }
	if run.Err()!=nil { return Output{Status:"incomplete",Diagnostics:[]semantic.Diagnostic{{Code:n.Key("limit"),Subject:n.Coverage(),Message:"TypeScript compiler exceeded the provider deadline"}}},nil }
	if err!=nil { return Output{},fmt.Errorf("%w: %v: %s",n.ErrCompiler,err,stderr.String()) }
	var out Output; decoder:=json.NewDecoder(bytes.NewReader(stdout.Bytes())); decoder.DisallowUnknownFields()
	if err:=decoder.Decode(&out); err!=nil { return Output{},fmt.Errorf("%w: malformed response: %v",n.ErrCompiler,err) }
	var trailing any; if decoder.Decode(&trailing)!=io.EOF || out.Version!=1 { return Output{},n.ErrCompiler }
	switch out.Status {
	case "ok": if out.Facts==nil || (query!=nil && out.Assignable==nil) { return Output{},n.ErrCompiler }
	case "invalid","unsupported","incomplete": if out.Facts!=nil || out.Assignable!=nil { return Output{},n.ErrCompiler }
	case "mismatch": return Output{},semantic.ErrInvalidRequest
	case "environment": return Output{},fmt.Errorf("%w: %v",n.ErrCompiler,out.Diagnostics)
	default: return Output{},n.ErrCompiler
	}
	return out,nil
}
