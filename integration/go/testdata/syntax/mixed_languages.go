package consumer

import (
	"github.com/Bitspark/bittype/syntax/cpp/go"
	"github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/syntax/typescript/go"
)

// Four deliberately invalid assignments, checked without executing the fixture.
var wrongTypeScript typescript.AST = golang.AST{}
var wrongGo golang.AST = rust.AST{}
var wrongRust rust.AST = cpp.AST{}
var wrongCPP cpp.AST = typescript.AST{}
