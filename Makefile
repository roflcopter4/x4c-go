#
# Makefile for BSD Make (bmake)
# 
# 'Installing' all the packages is critical for my wonky Neovim highlighting
# plugin. This just makes it easier, and also generates the necessary antlr
# output. BSD Make has a functional for loop, saving several whole seconds (!)
# of copy-pasting. Never leave home without it.
#

.SUFFIXES: .g4 .go
.MAIN: all

TARGET := x4c

local_go_args := #-compiler gccgo -gccgoflags '-O3 -march=native -g'
gccgo_CFLAGS  := -O3 -ftree-vectorize -fdiagnostics-color=always -Wall -march=native -g3 -gdwarf-5 -ffunction-sections -flto
gccgo_args    := -compiler gccgo -gccgoflags '-O3 -ftree-vectorize -fdiagnostics-color=always -Wall -march=native'
gccgo_args_g  := -x -compiler gccgo

antlrdir     := ${.CURDIR}/translation/gen
lexergrammar := ${antlrdir}/X4Lex.g4
parsergrammar:= ${antlrdir}/X4Parse.g4
combinedgrammar := ${antlrdir}/X4C.g4

target := x4c
dirs   := config myxml util util/color \
	  translation                  \
	  translation/ast              \
	  translation/newast           \
	  translation/handle_xml       \
	  translation/handle_script    \
	  translation/gen/sepLexer     \
	  translation/gen/sepParser

dirs_uscore := ${dirs:S/\//___/g}

# all: antlr .WAIT install_dirs .WAIT target
all: antlr .WAIT target

target:
	go build
#	go build -linkshared -a ${local_go_args} -o "${TARGET}"
#	go build -linkshared -a ${gccgo_args_g} -o "${TARGET}-gccgo"

fast:  target
quick: fast
build: fast

parser: antlr .WAIT target

dirs: install_dirs .WAIT target

antlr:
	antlr4 -Xexact-output-dir -Dlanguage=Go -package sepParser              -no-visitor -long-messages -o "${antlrdir}/sepParser" "${parsergrammar}"
	antlr4 -Xexact-output-dir -Dlanguage=Go -package sepLexer                           -long-messages -o "${antlrdir}/sepLexer"  "${lexergrammar}"

#	antlr4 -Xexact-output-dir -Dlanguage=Go -package parser    -Xlog -atn -long-messages -o "${antlrdir}/combined"  "${combinedgrammar}"

install_dirs: ${dirs_uscore}

.for DIR_USCORE in ${dirs_uscore}
${DIR_USCORE}:
	go install ${local_go_args} "${.CURDIR}/${DIR_USCORE:S/___/\//g}"
	go install ${gccgo_args_g} "${.CURDIR}/${DIR_USCORE:S/___/\//g}"
.endfor

.PHONY: ${dirs_uscore} all fast quick build parser dirs antlr install_dirs
