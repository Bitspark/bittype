package rust

import (
	g "github.com/Bitspark/bittype/internal/grammar"
	"strings"
)

var reserved = map[string]bool{}

func init() {
	for _, word := range strings.Fields("as break const continue crate else enum extern false fn for if impl in let loop match mod move mut pub ref return self Self static struct super trait true type unsafe use where while async await dyn abstract become box do final macro override priv typeof unsized virtual yield try gen") {
		reserved[word] = true
	}
}
func Identifier(v string) bool {
	if strings.HasPrefix(v, "r#") {
		name := v[2:]
		return g.Identifier(name) && name != "_" && name != "self" && name != "Self" && name != "super" && name != "crate"
	}
	return g.Identifier(v) && !reserved[v]
}
func Declaration(v string) bool { return v != "_" && Identifier(v) }
func LifetimeParameter(v string) bool { return v != "_" && g.Identifier(v) && !reserved[v] }
func LifetimeReference(v string) bool {
	return v == "static" || v == "_" || LifetimeParameter(v)
}
func PathSegment(v string) bool {
	return v == "self" || v == "Self" || v == "super" || v == "crate" || Identifier(v)
}
