#
# I'm very lazy
#

.SUFFIXES: .peg .go

peg   := translate/toXML/pig.pigeon
dirs  := ast config myxml translate/toX4C translate/toXML util util/colour
target:= x4c

all: run_peg .WAIT install_dirs
	go build

install_dirs:
.for DIR in ${dirs}
	go install "${.CURDIR}/${DIR}"
.endfor

run_peg:
.for f in ${peg}
	pigeon -o "${f}.go" "${f}"
#	peg -inline -switch -noast -strict "${f}"
.endfor
