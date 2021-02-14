package handle_xml

import (

	// XML "github.com/lestrrat-go/libxml2"
	// XMLxsd "github.com/lestrrat-go/libxml2/xsd"
	// XMLclib "github.com/lestrrat-go/libxml2/clib"
	// XMLtypes "github.com/lestrrat-go/libxml2/types"
	// XMLdom "github.com/lestrrat-go/libxml2/dom"
	// XMLtypes "github.com/lestrrat-go/libxml2/types"

	"fmt"
	"strings"

	"github.com/roflcopter4/x4c-go/translation/ast"
	"github.com/roflcopter4/x4c-go/util"
)

const DEFAULT_INDENT = 4

var Indent_Size = DEFAULT_INDENT

type cur_data struct {
	lines       []string
	tree        ast.AST
	depth       int
	indent_size int
}

func make_output(tree ast.AST) []string {
	data := &cur_data{
		lines:       make([]string, 0, 256),
		tree:        tree,
		depth:       0,
		indent_size: Indent_Size,
	}

	data.walk_tree(tree.StartNode())
	return data.lines
}

func (data *cur_data) walk_tree(node ast.Node) {
	str := data.handle_node(node)

	switch {
	case node.HasFlags(ast.NFlagBlock):
		data.lines = append(data.lines, str+" {")
		data.depth++

		for _, child := range node.GetChildren() {
			data.walk_tree(child)
		}

		data.depth--
		data.lines = append(data.lines, data.get_indent()+"}")

	case node.HasFlags(ast.NFlagText):
		data.lines = append(data.lines, str)

	default:
		data.lines = append(data.lines, str+";")
	}
}

//========================================================================================

func (data *cur_data) handle_node(node ast.Node) string {
	var str string

	switch n := node.(type) {
	case *ast.XMLStatement:
		str = data.handle_xml_statement(n)

	case *ast.XMLComment:
		str = fmt.Sprintf("%s/*%s*/", data.get_indent(), n.Text)

	case *ast.ConditionStatement:
		str = data.handle_conditional_statement(n)

	default:
		util.Die(1, "Invalid type somehow (%[1]T):\n(%#[1]v)", n)
	}

	return str
}

func (data *cur_data) handle_xml_statement(node *ast.XMLStatement) string {
	str := data.get_indent() + node.Name + "<"

	for i, attr := range node.Attributes {
		if i == 0 {
			str += fmt.Sprintf("%s=\"%s\"", attr.Name, attr.Val.Raw)
		} else {
			str += fmt.Sprintf(" %s=\"%s\"", attr.Name, attr.Val.Raw)
		}
	}

	return str + ">"
}

func (data *cur_data) handle_conditional_statement(node *ast.ConditionStatement) string {
	str := data.get_indent()

	switch node.Type {
	case ast.ConditionIf:
		str += "if " + condition_expression(node.Expr)
	case ast.ConditionElseif:
		str += "elseif " + condition_expression(node.Expr)
	case ast.ConditionElse:
		str += "else"
	case ast.ConditionWhile:
		str += "while " + condition_expression(node.Expr)
	default:
		panic(fmt.Errorf("Invalid condition type (%d)", node.Type))
	}

	return str
}

//========================================================================================

func (data *cur_data) get_indent() string {
	return strings.Repeat(" ", data.depth*data.indent_size)
}

func condition_expression(expr *ast.Expression) string {
	var str string
	if expr.XML == nil {
		str = fmt.Sprintf("(%s)", expr.Raw)
	} else {
		str = fmt.Sprintf("<%s>", expr.Raw)
	}
	return str
}
