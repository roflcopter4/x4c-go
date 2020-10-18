package main

import (
	"os"

	"github.com/roflcopter4/x4c/myxml"
	"github.com/roflcopter4/x4c/util"

	"github.com/roflcopter4/x4c/translate/toX4C"
	"github.com/roflcopter4/x4c/translate/toXML"

	"github.com/pborman/getopt"
)

var opt struct {
	infile  string
	outfile struct {
		name string
		fp   *os.File
	}
	help     bool
	dump     bool
	validate bool

	operation string
}

func init() {
	getopt.BoolVarLong(&opt.help, "help", 'h', "Display help")
	getopt.BoolVarLong(&opt.dump, "dump", 'D', "Dump the tree")
	getopt.BoolVarLong(&opt.validate, "validate", 'v', "Validate input with schema")
	getopt.StringVarLong(&opt.infile, "file", 'f', "Input filename")
	getopt.StringVarLong(&opt.outfile.name, "out", 'o', "Output filename")

	getopt.EnumVarLong(&opt.operation, "operation", 'x', []string{"c", "d"}, "Operation")
}

func main() {
	handle_opts()

	switch opt.operation {
	case "c":
		compile()
	case "d":
		do_everything()
	default:
		util.Die(1, "Must specify an operation. Dipshit.")
	}
}

func handle_opts() {
	getopt.Parse()
	if opt.help {
		getopt.Usage()
		os.Exit(0)
	}

	if opt.infile == "" {
		if getopt.NArgs() == 0 {
			getopt.Usage()
			os.Exit(1)
		}
		opt.infile = getopt.Args()[0]
	}

	if str_eq_any(opt.outfile.name, "", "-") {
		opt.outfile.fp = os.Stdout
	} else {
		fp, err := os.Create(opt.outfile.name)
		if err != nil {
			util.PanicE(err)
		}
		opt.outfile.fp = fp
	}
}

func compile() {
	toXML.Translate(opt.infile)
}

func do_everything() {
	doc, err := myxml.New_Document(opt.infile)
	if err != nil {
		util.DieE(2, err)
	}
	defer doc.Free()

	if opt.validate {
		err = doc.GetSchema()
		if err != nil {
			util.DieE(2, err)
		}

		err = doc.ValidateSchema()
		if err != nil {
			util.DieE(2, err)
		}
	}

	// doc.GetSchema()
	// d2, err := myxml.New_Document(opt.infile)
	// d2.SetSchema(doc.Schema())
	// doc.SetSchema(nil)
	// doc.Free()
	//
	// dumb.TestReader(d2)
	//
	// d2.Free()

	toX4C.TestReader(opt.outfile.fp, doc)

	// err = dumb.Dumb(doc, opt.outfile.fp)
	// if err != nil {
	//       util.PanicE(err)
	// }
}

func str_eq_any(cmp string, lst ...string) bool {
	for _, s := range lst {
		if cmp == s {
			return true
		}
	}
	return false
}

func str_eq_all(cmp string, lst ...string) bool {
	for _, s := range lst {
		if cmp != s {
			return false
		}
	}
	return true
}
