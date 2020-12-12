package untranslate

import (
	// XML "github.com/lestrrat-go/libxml2"
	XMLtypes "github.com/lestrrat-go/libxml2/types"
	// XMLxsd "github.com/lestrrat-go/libxml2/xsd"
	XMLdom "github.com/lestrrat-go/libxml2/dom"
	// XMLclib "github.com/lestrrat-go/libxml2/clib"

	// XMLreader "github.com/roflcopter4/xml_addon/reader"

	"github.com/roflcopter4/x4c-go/translation/ast"
	"github.com/roflcopter4/x4c-go/util"
)

type cur_data struct {
	tree ast.AST
	doc  *XMLdom.Document
	cur  XMLtypes.Node
}

func create_xml(tree ast.AST) *XMLdom.Document {
	data := new(cur_data)
	data.doc = XMLdom.NewDocument("1.0", "UTF-8")
	data.tree = tree
	data.cur = nil

	{
		node := tree.GetRoot().GetChildren()[0].(*ast.XMLStatement)
		element := data.handle_xml_statement(node)
		data.doc.SetDocumentElement(element)
		data.cur = element
		for _, ch := range node.GetChildren() {
			data.walk_tree(ch)
		}
	}

	return data.doc
}

func (data *cur_data) walk_tree(node ast.Node) {
	cur := data.handle_node(node)

	if node.HasFlags(ast.NFlagBlock) {
		tmp := data.cur
		data.cur = cur
		for _, ch := range node.GetChildren() {
			data.walk_tree(ch)
		}
		data.cur = tmp
	}
}

func (data *cur_data) handle_node(node ast.Node) XMLtypes.Node {
	var ret XMLtypes.Node

	switch n := node.(type) {
	case *ast.RootNode:
		panic("wtf")

	case *ast.XMLStatement:
		el := data.handle_xml_statement(n)
		data.cur.AddChild(el)
		ret = el

	case *ast.XMLComment:
		com, err := data.doc.CreateCommentNode(n.Text)
		if err != nil {
			panic(err)
		}
		data.cur.AddChild(com)
		ret = com

	case *ast.ConditionStatement:
		panic("IM LAZY")

	case *ast.AstNode:
		util.Die(1, "Invalid type somehow %+v\n", n)
	default:
		util.Die(1, "Invalid type somehow %+v\n", n)
	}

	ret.MakeMortal()
	return ret
}

func (data *cur_data) handle_xml_statement(node *ast.XMLStatement) XMLtypes.Element {
	el, err := data.doc.CreateElement(node.Name)
	if err != nil {
		panic(err)
	}

	for _, attr := range node.Attributes {
		el.SetAttribute(attr.Name, attr.Val.Raw)
	}

	return el
}
