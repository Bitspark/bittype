package typescript

import (
	"strings"
	g "github.com/Bitspark/bittype/internal/grammar"
)
var reserved = map[string]bool{}
func init() {
	for _,word := range strings.Fields("break case catch class const continue debugger default delete do else enum export extends false finally for function if import in instanceof new null return super switch throw true try typeof var void while with implements interface let package private protected public static yield await this") {reserved[word]=true}
}
func Identifier(v string)bool{return g.Identifier(strings.ReplaceAll(v,"$","_"))&&!reserved[v]}
