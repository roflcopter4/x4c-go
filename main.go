package main

import (
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/pborman/getopt"

	"github.com/roflcopter4/x4c-go/myxml"
	"github.com/roflcopter4/x4c-go/translation"
	"github.com/roflcopter4/x4c-go/util"
)

var opt struct {
	infname string
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
	getopt.StringVarLong(&opt.infname, "file", 'f', "Input filename")
	getopt.StringVarLong(&opt.outfile.name, "out", 'o', "Output filename")

	getopt.EnumVarLong(&opt.operation, "operation", 'x', []string{"u", "t", "q", "Q", "s", "S", "N"}, "Operation")
}

/***************************************************************************************/

func main() {
	args := handle_opts()

	switch opt.operation {
	case "u":
		translation.Translate_Script(opt.outfile.fp, opt.infname)
	case "t":
		do_translate()
	case "q":
		translation.TestLexer(args[0], true)
	case "Q":
		translation.TestLexer(args[0], false)
	case "s":
		translation.TestScriptLexer(args[0], true)
	case "S":
		translation.TestScriptLexer(args[0], false)
	case "N":
		translation.TestNewAst(args[0])
	default:
	}
}

func handle_opts() []string {
	getopt.Parse()
	args := getopt.Args()

	if opt.help {
		getopt.Usage()
		os.Exit(0)
	}

	if opt.infname == "" {
		if getopt.NArgs() == 0 {
			getopt.Usage()
			os.Exit(1)
		}
		opt.infname = args[0]
	}

	if opt.operation == "" {
		switch filepath.Ext(opt.infname) {
		case ".xml":
			opt.operation = "t"
		case ".x4c":
			opt.operation = "u"
		default:
			util.Die(1, "Can't identify input file \"%s\". Please specify type.", opt.infname)
		}
	}

	if util.StrEqAny(opt.outfile.name, "", "-") {
		opt.outfile.fp = os.Stdout
	} else {
		fp, err := os.Create(opt.outfile.name)
		if err != nil {
			util.PanicUnwrap(err)
		}
		opt.outfile.fp = fp
	}

	return args
}

func do_translate() {
	doc, err := myxml.New_Document(opt.infname)
	if err != nil {
		util.DieE(2, err)
		// panic(err)
	}
	defer doc.Free()

	if opt.validate {
		err = doc.GetSchema()
		if err != nil {
			util.DieE(2, err)
			// panic(err)
		}

		err = doc.ValidateSchema()
		if err != nil {
			util.DieE(2, err)
			// panic(err)
		}
	}

	// translate.TestReader(opt.outfile.fp, doc)
	// translation.TestTranslate(opt.outfile.fp, doc)
	translation.SetIndent(2)
	translation.Translate_XML(opt.outfile.fp, doc)
}

/***************************************************************************************/

// Initialize spew config
func init() {
	spew.Config.Indent = "  "
	spew.Config.DisableCapacities = true
	// spew.Config.DisableMethods = true
	// spew.Config.DisablePointerMethods = true
	spew.Config.DisablePointerAddresses = true
}
