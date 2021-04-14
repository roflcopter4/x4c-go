package handle_xml

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	// XML "github.com/lestrrat-go/libxml2"
	// XMLxsd "github.com/lestrrat-go/libxml2/xsd"

	XMLclib "github.com/lestrrat-go/libxml2/clib"
	XMLtypes "github.com/lestrrat-go/libxml2/types"

	XMLreader "github.com/roflcopter4/xml_addon/reader"

	"github.com/roflcopter4/x4c-go/myxml"
	"github.com/roflcopter4/x4c-go/translation/ast"
	"github.com/roflcopter4/x4c-go/translation/newast"

	"github.com/roflcopter4/x4c-go/util"
)

//========================================================================================

type builder struct {
	doc   myxml.DocWrapper
	rd    XMLreader.TextReader
	a     ast.AST
	cur   ast.Node
	block ast.Node

	wg sync.WaitGroup
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

		/* Sometimes text nodes and comments have to be identified via the node,
		 * other times things work just fine via the reader. I dunno. I'm probably
		 * doing something wrong but who knows what. */
		switch node.NodeType() {
		case XMLclib.TextNode:
			b.TextNode(node)

		case XMLclib.CommentNode:
			b.CommentNode(node)

		default:
			switch reader.NodeType() {
			case XMLreader.Reader_Text:
				b.TextNode(node)

			case XMLreader.Reader_Comment:
				b.CommentNode(node)

			case XMLreader.Reader_Element:
				b.StartElement(node)

			case XMLreader.Reader_EndElement:
				b.EndElement(node)
			}
		}

	}

	if outfp == os.Stdout {
		if outfp, err = os.Open(os.DevNull); err != nil {
			panic(err)
		}
	}

	lines := make_output(b.a)
	fmt.Fprintln(outfp, strings.Join(lines, "\n"))
	b.wg.Wait()
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

func (b *builder) TextNode(node XMLtypes.Node) {
	for nnl := strings.Count(node.NodeValue(), "\n"); nnl > 1; nnl-- {
		b.block.AddTextNode("")
	}
}

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
			b.new_dump("value", b.rd.Value())
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
						//Val:  ast.NewExpression(b.rd.Value()),
						Val: b.lazy_lazy_lazy(b.rd.Name(), b.rd.Value()),
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
			// expr := ast.NewExpression(b.rd.Value())
			expr := b.lazy_lazy_lazy(b.rd.Name(), b.rd.Value())
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

func (b *builder) new_dump(name, val string) {
	newast.Parse_Expression(&b.wg, b.doc.Name().Base(), name, val)
}

func (b *builder) lazy_lazy_lazy(name, val string) *ast.Expression {
	if !util.StrEqAny(name, "comment", "xmlns:xsi", "xsi:noNamespaceSchemaLocation") {
		b.new_dump(name, val)
	}
	return ast.NewExpression(val)
}
