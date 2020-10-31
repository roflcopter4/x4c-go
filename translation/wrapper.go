package translation

import (
	"os"

	"github.com/roflcopter4/x4c-go/myxml"
	"github.com/roflcopter4/x4c-go/translation/translate"
	"github.com/roflcopter4/x4c-go/translation/untranslate"
)

func TestTranslate(outfp *os.File, doc myxml.DocWrapper) {
	translate.TestReader(outfp, doc)
}

func Translate(outfp *os.File, doc myxml.DocWrapper) {
	translate.Translate(outfp, doc)
}

func UnTranslate(outfp *os.File, fname string) {
	untranslate.Translate(outfp, fname)
}

func TestLexer(str string) {
	untranslate.TestLexer(str)
}
