package integration_test

import (
 "fmt"
 goast "go/ast"
 "go/importer"
 "go/parser"
 "go/token"
 "go/types"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

// Import actual production sources, caching both standard and local identities.
// Fixtures are independent packages; no manually redefined public API is tested.
type engineImporter struct {
 standard types.Importer
 packages map[string]*types.Package
}
func (i *engineImporter) Import(path string)(*types.Package,error){
 if p,ok:=i.packages[path]; ok { return p,nil }
 if !strings.HasPrefix(path,modulePath+"/") { p,err:=i.standard.Import(path); if err==nil { i.packages[path]=p }; return p,err }
 relative:=strings.TrimPrefix(path,modulePath+"/")
 dir:=filepath.Join("..","..",filepath.FromSlash(relative))
 entries,err:=os.ReadDir(dir); if err!=nil { return nil,err }
 positions:=token.NewFileSet()
 var files []*goast.File
 for _,entry:=range entries {
  if !strings.HasSuffix(entry.Name(),".go") || strings.HasSuffix(entry.Name(),"_test.go") { continue }
  file,err:=parser.ParseFile(positions,filepath.Join(dir,entry.Name()),nil,parser.AllErrors)
  if err!=nil { return nil,err }; files=append(files,file)
 }
 config:=types.Config{GoVersion:"go1.27",Importer:i}
 p,err:=config.Check(path,positions,files,nil)
 if err==nil { i.packages[path]=p }
 return p,err
}
func TestEngineIndependentPublicContracts(t *testing.T){
 imports:=&engineImporter{standard:importer.Default(),packages:make(map[string]*types.Package)}
 cases:=[]struct{name,want string}{{"consumer",""},{"wrong_stage","cannot use"},{"wrong_view","cannot use"},{"missing_endpoints","FromID"},{"private_envelope","unexported"}}
 for _,tc:=range cases { t.Run(tc.name,func(t *testing.T){
  positions:=token.NewFileSet()
  file,err:=parser.ParseFile(positions,filepath.Join("testdata","engine",tc.name+".go"),nil,parser.AllErrors)
  if err!=nil { t.Fatal(err) }
  config:=types.Config{GoVersion:"go1.27",Importer:imports}
  _,err=config.Check("example.invalid/engine-consumer/"+tc.name,positions,[]*goast.File{file},nil)
  if tc.want=="" { if err!=nil { t.Fatal(err) }; return }
  if err==nil || !strings.Contains(err.Error(),tc.want) { t.Fatal(fmt.Sprintf("want %q: %v",tc.want,err)) }
 })}
}
