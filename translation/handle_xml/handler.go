package handle_xml

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
	"github.com/roflcopter4/x4c-go/util"
)

//========================================================================================

type output struct {
	rd     XMLreader.TextReader
	lines  []string
	depth  int
	spaces int
}

func TestReader(outfp *os.File, d myxml.DocWrapper) {
	reader, err := XMLreader.NewTextReaderFromDoc(d.Doc())
	if err != nil {
		util.DieE(1, err)
	}
	defer reader.Free()

	out := &output{
		rd:     reader,
		lines:  []string{},
		depth:  0,
		spaces: 2,
	}

	for reader.TextRead() != 0 {
		node, _ := reader.CurrentNode()

		switch node.NodeType() {
		case XMLclib.TextNode:
			handle_text(out, node)
			continue
		case XMLclib.CommentNode:
			handle_comment(out, node)
			continue
		}

		switch reader.NodeType() {
		case XMLreader.Reader_Text:
			handle_text(out, node)

		case XMLreader.Reader_Comment:
			handle_comment(out, node)

		case XMLreader.Reader_Element:
			// handle_generic(out, node)
			handle_start_element(out, node)

		case XMLreader.Reader_EndElement:
			handle_end_element(out)
		}
	}

	fmt.Fprintln(outfp, strings.Join(out.lines, "\n"))
}

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
			if i > 0 {
				str += ", "
			}
			str += fmt.Sprintf("%s=\"%s\"", out.rd.Name(), out.rd.Value())
		}
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
}

func handle_comment(out *output, node XMLtypes.Node) {
	str := make_indent(out) + fmt.Sprintf("/*%s*/", node.NodeValue())
	out.lines = append(out.lines, str)
}

func handle_start_element(out *output, node XMLtypes.Node) {
	nn := node.NodeName()
	i := sort.SearchStrings(special_idents, nn)

	if i < len(special_idents) && special_idents[i] == nn {
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
			str += fmt.Sprintf("if (%s) {", get_attr_string(attributes))
		case "do_elseif":
			str += fmt.Sprintf("elseif (%s) {", get_attr_string(attributes))
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
