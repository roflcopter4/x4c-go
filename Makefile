#
# I'm very lazy
#

.SUFFIXES: .g4 .go

target := x4c
file   := translation/untranslate/X4C.g4
dirs   := config myxml util util/colour translation \
	  translation/ast \
	  translation/translate \
	  translation/untranslate \
	  translation/untranslate/parser

all: antlr .WAIT install_dirs
	go build

quick: antlr
	go build

fast:
	go build

skip: install_dirs
	go build

antlr:
	antlr4 -Dlanguage=Go -long-messages -o "${.CURDIR}/${file:H}/parser" "${.CURDIR}/${file}"

install_dirs:
.for DIR in ${dirs}
	go install "${.CURDIR}/${DIR}"
.endfor
