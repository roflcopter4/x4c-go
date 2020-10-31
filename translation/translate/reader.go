package translate

import (
	"fmt"
	"os"
	"sort"
	"strings"

	// XML "github.com/lestrrat-go/libxml2"
	// XMLxsd "github.com/lestrrat-go/libxml2/xsd"

	XMLclib "github.com/lestrrat-go/libxml2/clib"
	XMLtypes "github.com/lestrrat-go/libxml2/types"

	XMLreader "github.com/roflcopter4/xml_addon/reader"

	"github.com/roflcopter4/x4c-go/myxml"
	"github.com/roflcopter4/x4c-go/translation/ast"
	"github.com/roflcopter4/x4c-go/util"
)

//========================================================================================

type builder struct {
	doc   myxml.DocWrapper
	rd    XMLreader.TextReader
	a     ast.AST
	cur   ast.Node
	block ast.Node
}

func Translate(outfp *os.File, doc myxml.DocWrapper) {
	reader, err := XMLreader.NewTextReaderFromDoc(doc.Doc())
	if err != nil {
		util.DieE(1, err)
	}
	defer reader.Free()

	b := &builder{
		doc: doc,
		rd:  reader,
		a:   ast.NewAst(),
	}
	b.cur = b.a.Root()
	b.block = b.a.Root()

	for reader.TextRead() != 0 {
		node, _ := reader.CurrentNode()

		switch node.NodeType() {
		case XMLclib.CommentNode:
			b.CommentNode(node)
			continue
		}

		switch reader.NodeType() {
		case XMLreader.Reader_Text:

		case XMLreader.Reader_Comment:
			b.CommentNode(node)

		case XMLreader.Reader_Element:
			b.StartElement(node)

		case XMLreader.Reader_EndElement:
			b.EndElement(node)
		}
	}

	lines := make_output(b.a)
	fmt.Fprintln(outfp, strings.Join(lines, "\n"))
}

//========================================================================================

func (b *builder) StartElement(node XMLtypes.Node) {
	nn := node.NodeName()
	i := sort.SearchStrings(xs_eids, nn)

	if i < len(xs_eids) && xs_eids[i] == nn {
		switch nn {
		case "do_if":
			b.Conditional(node, ast.ConditionIf)
		case "do_elseif":
			b.Conditional(node, ast.ConditionElseif)
		case "do_else":
			b.Conditional(node, ast.ConditionElse)
		default:
			panic("Impossible!")
		}
	} else {
		b.GenericXML(node)
	}
}

func (b *builder) EndElement(node XMLtypes.Node) {
	b.block = b.block.Parent()
}

//========================================================================================

func (b *builder) CommentNode(node XMLtypes.Node) {
	b.block.AddComment(node.NodeValue())
}

func (b *builder) GenericXML(node XMLtypes.Node) {
	stmt := b.block.AddXMLStatement(b.rd.Name())
	b.cur = stmt
	b.add_attributes(stmt)

	if node.HasChildNodes() {
		b.block = stmt
	}
}

func (b *builder) Conditional(node XMLtypes.Node, ctype int) {
	var (
		expr *ast.Expression = nil
		// elem                 = node.(XMLtypes.Element)
	)

	if ctype == ast.ConditionIf || ctype == ast.ConditionElseif {
		// attr, _ := elem.Attributes()
		// expr = ast.NewExpression(get_attr_string(attr))

		// expr.XML = new(ast.XMLStatement);
		expr = new(ast.Expression)
		nattr := b.rd.AttributeCount()
		b.rd.MoveToAttributeNo(0)

		if nattr == 1 && b.rd.Name() == "value" {
			expr.Raw = b.rd.Value()
		} else {
			attr, _ := node.(XMLtypes.Element).Attributes()
			expr.Raw = get_attr_string(attr)
			expr.XML = make([]*ast.XMLAttribute, nattr)
			for i := 0; i < nattr; i++ {
				b.rd.MoveToAttributeNo(i)
				if node, _ := b.rd.CurrentNode(); node != nil {
					e := ast.XMLAttribute{Name: b.rd.Name(), Val: ast.NewExpression(b.rd.Value())}
					expr.XML[i] = &e
					// expr := ast.NewExpression(b.rd.Value())
					// stmt.AddAttribute(b.rd.Name(), expr)
				}
			}
		}
	}

	stmt := b.block.AddConditionStatement(expr, ctype)
	b.cur = stmt
	b.block = stmt
}

//========================================================================================

func (b *builder) add_attributes(stmt *ast.XMLStatement) {
	for i, nattr := 0, b.rd.AttributeCount(); i < nattr; i++ {
		b.rd.MoveToAttributeNo(i)
		if node, _ := b.rd.CurrentNode(); node != nil {
			expr := ast.NewExpression(b.rd.Value())
			stmt.AddAttribute(b.rd.Name(), expr)
		}
	}
}
