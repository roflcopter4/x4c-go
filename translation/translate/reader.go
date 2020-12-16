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
	b.cur = b.a.GetRoot()
	b.block = b.a.GetRoot()

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
	i := sort.SearchStrings(special_idents, nn)

	if i < len(special_idents) && special_idents[i] == nn {
		switch nn {
		case "do_if":
			b.Conditional(node, ast.ConditionIf)
		case "do_elseif":
			b.Conditional(node, ast.ConditionElseif)
		case "do_else":
			b.Conditional(node, ast.ConditionElse)
		case "do_while":
			b.Conditional(node, ast.ConditionWhile)

		case "do_for_each":
			b.GenericXML(node)
		case "do_all":
			b.GenericXML(node)

		default:
			panic("Impossible!")
		}
	} else {
		b.GenericXML(node)
	}
}

func (b *builder) EndElement(node XMLtypes.Node) {
	b.block = b.block.GetParent()
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
	var expr *ast.Expression = nil

	// if ctype == ast.ConditionIf || ctype == ast.ConditionElseif || ctype == ast.ConditionWhile {
	if nattr := b.rd.AttributeCount(); nattr > 0 {
		expr = new(ast.Expression)
		b.rd.MoveToAttributeNo(0)

		if nattr == 1 && b.rd.Name() == "value" {
			expr.Raw = b.rd.Value()
		} else {
			attr, _ := node.(XMLtypes.Element).Attributes()
			expr.Raw = get_attr_string(attr)
			expr.XML = new(ast.XMLStatement)
			expr.XML.Attributes = make([]*ast.XMLAttribute, nattr)

			for i := 0; i < nattr; i++ {
				b.rd.MoveToAttributeNo(i)
				if node, _ := b.rd.CurrentNode(); node != nil {
					expr.XML.Attributes[i] = &ast.XMLAttribute{
						Name: b.rd.Name(),
						Val:  ast.NewExpression(b.rd.Value()),
					}
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

func get_attr_string(lst []XMLtypes.Attribute) (ret string) {
	for i, a := range lst {
		if i > 0 {
			ret += " "
		}
		ret += fmt.Sprintf("%s=\"%s\"", a.NodeName(), a.NodeValue())
	}
	return
}
