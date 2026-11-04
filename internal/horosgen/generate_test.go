package horosgen

import (
 "context"
 "errors"
 "testing"
 "github.com/Bitspark/bittype/contracts/text/go"
 "github.com/Bitspark/bittype/trees/ast/go"
)
func TestArenaBoundAndCancellation(t *testing.T){
 a:=Arena[struct{}]{Ctx:context.Background(),Tree:ast.AST[struct{}]{Nodes:make([]ast.Node[struct{}],MaxNodes)}}
 if id:=a.Add("type",nil,ast.NodeRef{Tree:"source",Node:1});id!=0{t.Fatal("limit allocated node")};var diagnostic *text.Diagnostic;if !errors.As(a.Err,&diagnostic)||diagnostic.Code!=text.LimitExceeded{t.Fatal(a.Err)}
 ctx,cancel:=context.WithCancel(context.Background());cancel();b:=Arena[struct{}]{Ctx:ctx};if b.Add("type",nil,ast.NodeRef{})!=0||!errors.Is(b.Err,context.Canceled){t.Fatal(b.Err)}
}
func TestConservativeIdentifiers(t *testing.T){
 for _,name:=range []string{"","_","a__b","a-b","é"}{if Identifier(name){t.Fatal("unsupported identifier accepted",name)}}
 for _,pair:=range [][2]string{{"page","Page"},{"_page","Type_page"},{"Self","Self_"},{"Page","Page"}}{if got,err:=Exported(pair[0]);err!=nil||got!=pair[1]{t.Fatal(got,err)}}
}
