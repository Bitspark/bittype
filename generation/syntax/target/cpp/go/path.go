package target
import(
g "github.com/Bitspark/bittype/contracts/generation/go"
m "github.com/Bitspark/bittype/generation/syntax/target/go"
s "github.com/Bitspark/bittype/syntax/cpp/go"
"github.com/Bitspark/bittype/trees/ast/go"
)
// Path is native syntax; selecting a library artifact and its meaning is the caller's responsibility.
func(b *Backend[F]) Path(names []string,args []g.Fragment[s.Language],origin ast.NodeRef,absolute bool)(g.Fragment[s.Language],error){
if len(names)==0{return g.Fragment[s.Language]{},g.ErrInvalid};var segments []m.Child[s.Language]
for i,name:=range names{var children []m.Child[s.Language];if i==len(names)-1{children=edges(s.RoleArgument,args)};p,err:=b.Node(m.Node(s.KindNameSegment,m.Attributes("name",name),origin),children...);if err!=nil{return p,err};segments=append(segments,m.Child[s.Language]{Role:s.RoleSegment,Fragment:p})}
var attrs []ast.Attribute;if absolute{attrs=m.Attributes("global","true")};return b.Node(m.Node(s.KindNamedType,attrs,origin),segments...)
}

