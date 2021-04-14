#
# Makefile for BSD Make (bmake)
# 
# 'Installing' all the packages is critical for my wonky Neovim highlighting
# plugin. This just makes it easier, and also generates the necessary antlr
# output. BSD Make has a functional for loop, saving several whole seconds (!)
# of copy-pasting. Never leave home without it.
#

.SUFFIXES: .g4 .go

local_go_args := #-compiler gccgo -gccgoflags '-O3 -march=native -g'

antlrdir     := ${.CURDIR}/translation/gen
lexergrammar := ${antlrdir}/X4Lex.g4
parsergrammar:= ${antlrdir}/X4Parse.g4
combinedgrammar := ${antlrdir}/X4C.g4

target := x4c
dirs   := config myxml util util/color \
	  translation/gen/sepParser \
	  translation/gen/sepLexer \
	  translation \
	  translation/ast \
	  translation/handle_script \
	  translation/handle_xml \
	  translation/newast

dirs_uscore := ${dirs:S/\//___/g}

all: antlr .WAIT install_dirs
	go build ${local_go_args}

fast:
	go build ${local_go_args}

parser: antlr
	go build ${local_go_args}

dirs: install_dirs
	go build ${local_go_args}

antlr:
	antlr4 -Xexact-output-dir -Dlanguage=Go -package parser    -no-listener -no-visitor -long-messages -o "${antlrdir}/combined"  "${combinedgrammar}"
	antlr4 -Xexact-output-dir -Dlanguage=Go -package sepParser -no-listener -no-visitor -long-messages -o "${antlrdir}/sepParser" "${parsergrammar}"
	antlr4 -Xexact-output-dir -Dlanguage=Go -package sepLexer                           -long-messages -o "${antlrdir}/sepLexer"  "${lexergrammar}"
#	antlr4 -Xexact-output-dir -Dlanguage=Go -package parser    -Xlog -atn -long-messages -o "${antlrdir}/combined"  "${combinedgrammar}"

install_dirs: ${dirs_uscore}

.for DIR_USCORE in ${dirs_uscore}
${DIR_USCORE}:
	go install ${local_go_args} "${.CURDIR}/${DIR_USCORE:S/___/\//g}"
.endfor

.PHONY: ${dirs_uscore} all fast parser dirs antlr install_dirs
