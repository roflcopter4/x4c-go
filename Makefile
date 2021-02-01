#
# Makefile for BSD Make (bmake)
# 
# 'Installing' all the packages is critical for my wonky Neovim highlighting
# plugin. This just makes it easier, and also generates the necessary antlr
# output. BSD Make has a functional for loop, saving several whole seconds (!)
# of copy-pasting. Never leave home without it.
#

.SUFFIXES: .g4 .go

maingrammar   := translation/gen/X4C.g4
scriptgrammar := translation/gen/Script.g4

target := x4c
dirs   := config myxml util util/color \
	  translation/gen/parser \
	  translation/gen/scriptparser \
	  translation \
	  translation/ast \
	  translation/newast \
	  translation/translate \
	  translation/untranslate

all: antlr .WAIT install_dirs
	go build

quick: antlr
	go build

fast:
	go build

skip: install_dirs
	go build

antlr:
	antlr4 -Dlanguage=Go -package parser       -visitor -long-messages -o "${.CURDIR}/${maingrammar:H}/parser" "${.CURDIR}/${maingrammar}"
	antlr4 -Dlanguage=Go -package scriptparser -visitor -long-messages -o "${.CURDIR}/${scriptgrammar:H}/scriptparser" "${.CURDIR}/${scriptgrammar}"

install_dirs:
.for DIR in ${dirs}
	go install "${.CURDIR}/${DIR}"
.endfor
