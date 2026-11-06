package integration_test

import (
	"context"
	"strings"
	"testing"
	generationfixture "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	gp "github.com/Bitspark/bittype/provider/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func TestAutomaticGenerationIndependentComposition(t *testing.T){
	f,err:=generationfixture.New(context.Background());if err!=nil{t.Fatal(err)}
	out,err:=f.Automatic(context.Background());if err!=nil{t.Fatal(err)}
	tree,err:=representation.Unpack[s.Language](f.Target(),out.EntryAST());if err!=nil{t.Fatal(err)}
	text:=formatted(t,gp.Provider{},tree)
	for _,want:=range []string{"type Page[","type Outer[","type Nested[","type Cursor ","ExternalBox["}{if !strings.Contains(text,want){t.Fatalf("missing %q:\n%s",want,text)}}
	if strings.Count(text,"type Cursor ")!=1{t.Fatal("helper defined more than once")}
}
func TestAutomaticGenerationKeepsWholeResultAssociation(t *testing.T){
	f,err:=generationfixture.NewArgumentOnlySupport(context.Background());if err!=nil{t.Fatal(err)}
	out,err:=f.Automatic(context.Background());if err==nil||out.LanguageID()!=""{t.Fatalf("incorrect nested support accepted: %v",err)}
}
