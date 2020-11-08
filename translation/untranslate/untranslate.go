package untranslate

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/davecgh/go-spew/spew"

	"github.com/roflcopter4/x4c-go/translation/ast"
	"github.com/roflcopter4/x4c-go/translation/untranslate/parser"
)

func init() {
	spew.Config.Indent = "  "
}

func Translate(outfp *os.File, fname string) {
	tree := parse_file(fname)
	doc := create_xml(tree)
	out := doc.Dump(true)
	doc.Free()

	out = strings.ReplaceAll(out, "&#10;", "\n")
	outfp.WriteString(out)
}

/****************************************************************************************/

type listener struct {
	*parser.BaseX4CListener

	a     ast.AST
	cur   ast.Node
	block ast.Node
}

func TestLexer(str string) {
	// var (
	//       is    = antlr.NewInputStream(str)
	//       lexer = parser.NewX4CLexer(is)
	// )
	//
	// for {
	//       t := lexer.NextToken()
	//       if t.GetTokenType() == antlr.TokenEOF {
	//             break
	//       }
	//       fmt.Printf("%s (%q)\n", lexer.SymbolicNames[t.GetTokenType()], t.GetText())
	// }
	//

	lexer := get_lexer(str)
	for {
		t := lexer.NextToken()
		if t.GetTokenType() == antlr.TokenEOF {
			break
		}
		fmt.Printf("%s (%q)\n",
			lexer.SymbolicNames[t.GetTokenType()], t.GetText())
	}
}

func parse_file(fname string) ast.AST {
	var (
		lexer  = get_lexer(fname)
		stream = antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
		p      = parser.NewX4CParser(stream)
		l      = new(listener)
	)

	l.a = ast.NewAst()
	l.cur = l.a.Root()
	l.block = l.a.Root()

	antlr.ParseTreeWalkerDefault.Walk(l, p.Document())

	return l.a
}

func (l *listener) EnterCompoundStmt(c *parser.CompoundStmtContext) {
	l.block = l.cur
}

func (l *listener) ExitCompoundStmt(c *parser.CompoundStmtContext) {
	l.block = l.block.Parent()
}

func (l *listener) ExitXmlStmt(c *parser.XmlStmtContext) {
	stmt := l.block.AddXMLStatement(c.GetIdent().GetText())
	l.cur = stmt

	if lst := c.GetLst(); lst != nil {
		add_attrs(stmt, lst.GetChildren())
	}
}

func (l *listener) ExitConditionStmt(c *parser.ConditionStmtContext) {
	stmt := l.block.AddXMLStatement("do_" + c.GetIdent().GetText())
	l.cur = stmt

	if lst := c.GetLst(); lst != nil {
		ctx := lst.(*parser.ConditionExprContext)
		if ctx.AttributeList() != nil {
			add_attrs(stmt, ctx.AttributeList().GetChildren())
		} else {
			dumb := ctx.DumbExpr()
			val := dumb.GetText()
			val = val[1 : len(val)-1]
			stmt.AddAttribute("value", ast.NewExpression(val))
		}
	}
}

func (l *listener) ExitCommentStmt(c *parser.CommentStmtContext) {
	tok := c.GetChild(0).(antlr.TerminalNode).GetSymbol()
	txt := tok.GetText()

	switch tok.GetTokenType() {
	case parser.X4CParserLineComment:
		txt = strings.TrimPrefix(txt, "//")
		txt = strings.TrimRight(txt, "\r\n")
		if txt[0] == ' ' {
			txt += " "
		}

	case parser.X4CParserBlockComment:
		txt = strings.TrimPrefix(txt, "/*")
		txt = strings.TrimSuffix(txt, "*/")
	}

	l.block.AddComment(txt)
}

/****************************************************************************************/

func add_attrs(stmt *ast.XMLStatement, lst []antlr.Tree) {
	for _, child := range lst {
		switch a := child.(type) {
		case parser.IAttributeContext:
			val := a.GetVal().GetText()
			val = strings.Trim(val, "\"")
			expr := ast.NewExpression(val)
			stmt.AddAttribute(a.GetIdent().GetText(), expr)

		case parser.IAttributeListContext:
			add_attrs(stmt, a.GetChildren())
		}
	}
}

func get_lexer(fname string) *parser.X4CLexer {
	var charstream antlr.CharStream

	if fname == "-" {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		charstream = antlr.NewInputStream(string(b))
	} else {
		fs, err := antlr.NewFileStream(fname)
		if err != nil {
			panic(err)
		}
		charstream = fs
	}

	return parser.NewX4CLexer(charstream)
}
