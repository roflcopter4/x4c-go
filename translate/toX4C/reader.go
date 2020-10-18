package toX4C

import (
	"fmt"
	"os"
	"strings"

	// XML "github.com/lestrrat-go/libxml2"
	// XMLxsd "github.com/lestrrat-go/libxml2/xsd"
	XMLclib "github.com/lestrrat-go/libxml2/clib"

	XMLreader "github.com/roflcopter4/xml_addon/reader"

	"github.com/roflcopter4/x4c/myxml"
	"github.com/roflcopter4/x4c/util"
)

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

	// if err = reader.SetSchema(d.Schema()); err != nil {
	//       util.DieE(1, err)
	// }

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
			handle_generic(out, node)
			// handle_start_element(out, node)

		case XMLreader.Reader_EndElement:
			handle_end_element(out)
		}
	}

	fmt.Fprintln(outfp, strings.Join(out.lines, "\n"))
}
