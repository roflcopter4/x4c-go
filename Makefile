#
# I'm very lazy
#

.SUFFIXES: .g4 .go

maingrammar   := translation/untranslate/gen/X4C.g4
scriptgrammar := translation/untranslate/gen/Script.g4

target := x4c
dirs   := config myxml util util/colour translation \
	  translation/ast \
	  translation/translate \
	  translation/untranslate \
	  translation/untranslate/gen/parser \
	  translation/untranslate/gen/scriptparser

all: antlr .WAIT install_dirs
	go build

quick: antlr
	go build

fast:
	go build

skip: install_dirs
	go build

antlr:
	antlr4 -Dlanguage=Go -package parser -long-messages -o "${.CURDIR}/${maingrammar:H}/parser" "${.CURDIR}/${maingrammar}"
	antlr4 -Dlanguage=Go -package scriptparser -visitor -long-messages -o "${.CURDIR}/${scriptgrammar:H}/scriptparser" "${.CURDIR}/${scriptgrammar}"

install_dirs:
.for DIR in ${dirs}
	go install "${.CURDIR}/${DIR}"
.endfor