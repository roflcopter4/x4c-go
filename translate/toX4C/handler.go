package toX4C

import (
	"fmt"
	"sort"
	"strings"

	// XML "github.com/lestrrat-go/libxml2"
	// XMLxsd "github.com/lestrrat-go/libxml2/xsd"
	// XMLclib "github.com/lestrrat-go/libxml2/clib"

	XMLtypes "github.com/lestrrat-go/libxml2/types"
	// "github.com/roflcopter4/x4c/myxml"
	// "github.com/roflcopter4/x4c/util"
)

//========================================================================================

func handle_generic(out *output, node XMLtypes.Node) {
	var (
		str   = make_indent(out) + out.rd.Name() + "<"
		nattr = out.rd.AttributeCount()
	)

	for i := 0; i < nattr; i++ {
		out.rd.MoveToAttributeNo(i)
		node, _ := out.rd.CurrentNode()
		if node != nil {
			str += fmt.Sprintf("%s=\"%s\" ", out.rd.Name(), out.rd.Value())
		}
	}

	if nattr > 0 {
		str = str[:len(str)-1]
	}

	if node.HasChildNodes() {
		out.depth++
		str += "> {"
	} else {
		str += ">;"
	}

	out.lines = append(out.lines, str)
}

func handle_end_element(out *output) {
	out.depth--
	str := make_indent(out) + "}"
	out.lines = append(out.lines, str)
}

//========================================================================================

func handle_text(out *output, node XMLtypes.Node) {
	// fmt.Print(node.NodeValue())
	// return
	// str := replace_whitespace(node.NodeValue())
	// fmt.Printf("`%s'\n", str)
}

func handle_comment(out *output, node XMLtypes.Node) {
	// add_space(depth)
	// str := node.NodeValue()
	// fmt.Printf("%s", str)
	str := make_indent(out) + fmt.Sprintf("/*%s*/", node.NodeValue())
	out.lines = append(out.lines, str)
}

func handle_start_element(out *output, node XMLtypes.Node) {
	nn := node.NodeName()
	i := sort.SearchStrings(xs_eids, nn)

	if i < len(xs_eids) && xs_eids[i] == nn {
		// Handle special recognized keywords
		var (
			str             = make_indent(out)
			elem            = node.(XMLtypes.Element)
			attributes, err = elem.Attributes()
		)
		if err != nil {
			panic(err)
		}

		if node.HasChildNodes() {
			out.depth++
		}

		switch nn {
		case "do_if":
			str += fmt.Sprintf("if (%s) {", attributes)
		case "do_elseif":
			str += fmt.Sprintf("elseif (%s) {", attributes)
		case "do_else":
			str += fmt.Sprintf("else {")
		default:
			panic("Impossible!")
		}

		out.lines = append(out.lines, str)
	} else {
		handle_generic(out, node)
	}
}

//========================================================================================
// Util

func make_indent(out *output) string {
	return strings.Repeat(" ", out.spaces*out.depth)
}

func replace_whitespace(str string) string {
	str = strings.ReplaceAll(str, "\r\n", "\\n")
	str = strings.ReplaceAll(str, "\n", "\\n")
	str = strings.ReplaceAll(str, "\t", "\\t")
	return str
}
