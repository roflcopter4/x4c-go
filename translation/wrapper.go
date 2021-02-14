package translation

import (
	"os"

	"github.com/roflcopter4/x4c-go/myxml"
	"github.com/roflcopter4/x4c-go/translation/handle_script"
	"github.com/roflcopter4/x4c-go/translation/handle_xml"
)

func TestTranslate(outfp *os.File, doc myxml.DocWrapper) {
	handle_xml.TestReader(outfp, doc)
}

func Translate_XML(outfp *os.File, doc myxml.DocWrapper) {
	handle_xml.Translate(outfp, doc)
}

func Translate_Script(outfp *os.File, fname string) {
	handle_script.Translate(outfp, fname)
}

func TestLexer(str string, isfile bool) {
	handle_script.TestLexer(str, isfile)
}

func TestScriptLexer(str string, isfile bool) {
	handle_script.TestScriptLexer(str, isfile)
}

func SetIndent(indent int) {
	handle_xml.Indent_Size = indent
}
