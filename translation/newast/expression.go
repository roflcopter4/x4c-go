package newast

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/roflcopter4/x4c-go/translation/gen/sepLexer"
	parser "github.com/roflcopter4/x4c-go/translation/gen/sepParser"
	"github.com/roflcopter4/x4c-go/util/color"
)

type whatever struct {
	p     *parser.X4Parse
	depth int
	name  string
	val   string
	str   string
}

var OUTPUT_FILE = os.Stderr

// var ERROR_FILE = os.Stderr
var ERROR_FILE = OUTPUT_FILE

func init() {
	runtime.GOMAXPROCS(12)
}

type MyErrorListener struct {
	*antlr.DefaultErrorListener
	fname string
}

func (c *MyErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// fmt.Fprintln(os.Stderr, "line "+strconv.Itoa(line)+":"+strconv.Itoa(column)+" "+msg)
	fmt.Fprintf(ERROR_FILE, "\033[1;31mERROR:\033[0m \033[1mFile \"%s\"\033[0m:\tline %d: %d - %s\n", c.fname, line, column, msg)
}

func Parse_Expression(wg *sync.WaitGroup, fname, name, str string) {
	// fuckme(fname, name, str)
	if wg != nil {
		wg.Add(1)
		go do_parse_expression(wg, fname, name, str)
	} else {
		do_parse_expression(wg, fname, name, str)
	}
}

func do_parse_expression(wg *sync.WaitGroup, fname, name, str string) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	var (
		cs     = antlr.NewInputStream(str)
		lex    = sepLexer.NewX4Lex(cs)
		stream = antlr.NewCommonTokenStream(lex, antlr.TokenDefaultChannel)
		par    = parser.NewX4Parse(stream)
		d      = new(whatever)
		el     = new(MyErrorListener)
	)

	el.fname = fname
	par.RemoveErrorListeners()
	par.AddErrorListener(el)
	d.p = par
	d.depth = 0
	d.name = name
	d.val = str

	// d.mydump(par.DebugStatement())

	tree := par.DebugStatement()
	fmt.Fprintf(OUTPUT_FILE, "Looking at expression `%s=\"%s\"`\n%s\n", d.name, d.val, treesStringTree(tree, d.p.RuleNames, d.p, true))
}

/***************************************************************************************/
/* Debugging - New */

func treesStringTree(tree antlr.Tree, ruleNames []string, recog antlr.Recognizer, pretty bool) string {
	if ruleNames == nil && recog != nil {
		ruleNames = recog.GetRuleNames()
	}
	var depth = 0
	var recurse func(treeNode antlr.Tree) string

	recurse = func(treeNode antlr.Tree) string {
		s := treesGetNodeText(treeNode, ruleNames, nil)
		s = antlr.EscapeWhitespace(s, false)

		c := treeNode.GetChildCount()
		if c == 0 {
			return s
		}

		var res string
		if pretty {
			if depth > 0 {
				res = "\n" + strings.Repeat("    ", depth)
			}
			depth++
		}
		res += "(" + s + " "

		if c > 0 {
			res += recurse(treeNode.GetChild(0))
			for i := 1; i < c; i++ {
				res += " " + recurse(treeNode.GetChild(i))
			}
		}

		if pretty {
			depth--
		}

		res += ")"
		return res
	}

	return recurse(tree)
}

func treesGetNodeText(t antlr.Tree, ruleNames []string, recog antlr.Parser) string {
	if ruleNames == nil && recog != nil {
		ruleNames = recog.GetRuleNames()
	}

	if ruleNames != nil {
		switch t2 := t.(type) {
		case antlr.RuleNode:
			return color.Cyan(get_true_name(t2, ruleNames))
		case antlr.ErrorNode:
			return color.BRed(fmt.Sprint(t2))
		case antlr.TerminalNode:
			if t2.GetSymbol() != nil {
				return color.BMagenta(t2.GetSymbol().GetText())
			}
		}
	}

	// no recog for rule names
	payload := t.GetPayload()
	if p2, ok := payload.(antlr.Token); ok {
		return "BLEGH " + p2.GetText()
	}

	return fmt.Sprint("BOO", t.GetPayload())
}

func get_true_name(ctx antlr.RuleNode, ruleNames []string) string {
	typename := fmt.Sprintf("%T", ctx)
	typename = strings.TrimPrefix(typename, "*parser.")
	typename = strings.TrimPrefix(typename, "*sepParser.")
	typename = strings.TrimSuffix(typename, "Context")
	return typename
}

/***************************************************************************************/
/* Debugging */

var (
	colon, quot string
)

func init() {
	colon = color.Bold(":")
	quot = color.BOrange("\"")
}

func (d *whatever) mydump(ctx antlr.ParseTree) {
	d.do_mydump(ctx)
	fmt.Fprintf(OUTPUT_FILE, "Looking at expression `%s=\"%s\"`\n%s\n\n", d.name, d.val, d.str)
}

func (d *whatever) do_mydump(ctx antlr.ParseTree) {
	d.depth++
	defer func() { d.depth-- }()

	if n, ok := ctx.(antlr.TerminalNode); ok && n.GetSymbol().GetTokenType() == parser.X4ParseEOF {
		return
	}
	var (
		typename     string = fmt.Sprintf("%T", ctx)
		typename_len int
	)

	if strings.HasPrefix(typename, "*antlr.") {
		if typename[7:] == "TerminalNodeImpl" {
			typename = d.p.GetSymbolicNames()[ctx.(*antlr.TerminalNodeImpl).GetSymbol().GetTokenType()]

			if typename == "" {
				typename = "Unclassified Terminal"
				typename_len = len(typename)
				typename = color.BBlue(typename)
			} else {
				typename_len = len(typename)
				typename = color.BMagenta(typename)
			}
		} else {
			typename_len = len(typename)
			typename = color.BRed(typename)
		}
	} else {
		typename = strings.TrimPrefix(typename, "*parser.")
		typename = strings.TrimPrefix(typename, "*sepParser.")
		typename = strings.TrimSuffix(typename, "Context")
		typename_len = len(typename)
		typename = color.Green(typename)
	}

	//	do_print := func() {
	//		var (
	//			ns    = (30) - typename_len
	//			space = strings.Repeat(" ", ns)
	//		)
	//		thing := quot + strings.TrimRight(mytext(ctx), " ") + quot
	//		fmt.Fprintf(OUTPUT_FILE, "\033[1m%2d\033[0m:  %s\n", d.depth, typename+colon+space+thing)
	//	}
	//
	//	do_print()

	var (
		ns    = (50) - typename_len
		space = strings.Repeat(" ", ns)
		thing = quot + strings.TrimRight(mytext(ctx), " ") + quot
	)

	d.str += fmt.Sprintf("\033[1m%2d\033[0m:  %s\n", d.depth, typename+colon+space+thing)

	//switch ctx.(type) {
	//case parser.IExpressionContext:
	//	if _, ok := ctx.(*parser.Object_expressionContext); !ok {
	//		do_print()
	//	}
	//case antlr.TerminalNode:
	//	do_print()
	//}

	for _, child := range ctx.GetChildren() {
		d.do_mydump(child.(antlr.ParseTree))
	}
}

func mytext(prc antlr.ParseTree) string {
	var s string

	if prc.GetChildCount() == 0 {
		switch x := prc.(type) {

		case parser.IComparitiveOpContext:
			switch x.GetText() {
			case "gt":
				s = ">"
			case "lt":
				s = "<"
			case "ge":
				s = ">="
			case "le":
				s = "<="
			default:
				s = x.GetText()
			}

		//case antlr.TerminalNode:
		//	if x.GetSymbol().GetTokenType() == parser.X4CParserComparitiveOp {
		//		switch prc.GetText() {
		//		case "gt":
		//			s = ">"
		//		case "ge":
		//			s = ">="
		//		case "lt":
		//			s = "<"
		//		case "le":
		//			s = "<="
		//		default:
		//			s = prc.GetText()
		//		}
		//	} else {
		//		s = prc.GetText()
		//	}

		default:
			s = prc.GetText()
		}
		s += " "
	} else {
		for _, child := range prc.GetChildren() {
			s += mytext(child.(antlr.ParseTree))
		}
	}

	return s
}
