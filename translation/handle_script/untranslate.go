package handle_script

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/roflcopter4/x4c-go/translation/ast"
	"github.com/roflcopter4/x4c-go/translation/gen/parser"
)

func Translate(outfp *os.File, fname string) {
	tree := parse_file(fname)
	doc := create_xml(tree)
	out := doc.Dump(true)
	doc.Free()

	out = strings.ReplaceAll(out, "&#10;", "\n")
	outfp.WriteString(out)
}

func TestLexer(str string, isfile bool) {
	var (
		chs   antlr.CharStream
		lexer *parser.X4CLexer
	)

	if isfile {
		chs, lexer = get_lexer(str)

	} else {
		chs = antlr.NewInputStream(str)
		lexer = parser.NewX4CLexer(chs)
	}

	for {
		t := lexer.NextToken()
		if t.GetTokenType() == antlr.TokenEOF {
			break
		}
		ind := t.GetStart()
		fmt.Printf("%s (%q) -> (%d: %v)\n",
			lexer.SymbolicNames[t.GetTokenType()], t.GetText(), ind,
			chs.GetText(t.GetStart(), t.GetStop()))
	}
	fmt.Println(chs.GetText(0, chs.Size()))
}

/****************************************************************************************/

type listener struct {
	*parser.BaseX4CListener

	a     ast.AST
	cur   ast.Node
	block ast.Node

	chs antlr.CharStream
	lex *parser.X4CLexer
	par *parser.X4CParser
}

func parse_file(fname string) ast.AST {
	var (
		chs, lexer = get_lexer(fname)
		stream     = antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
		par        = parser.NewX4CParser(stream)

		l = &listener{
			a:   ast.NewAst(),
			chs: chs,
			lex: lexer,
			par: par,
		}
	)

	l.cur = l.a.GetRoot()
	l.block = l.a.GetRoot()

	antlr.ParseTreeWalkerDefault.Walk(l, par.Document())

	return l.a
}

func get_lexer(fname string) (antlr.CharStream, *parser.X4CLexer) {
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

	return charstream, parser.NewX4CLexer(charstream)
}

/****************************************************************************************/

func (l *listener) EnterCompoundStmt(c *parser.CompoundStmtContext) {
	l.block = l.cur
}

func (l *listener) ExitCompoundStmt(c *parser.CompoundStmtContext) {
	l.block = l.block.GetParent()
}

func (l *listener) EnterXmlStmt(c *parser.XmlStmtContext) {
	stmt := l.block.AddXMLStatement(c.GetIdent().GetText())
	l.cur = stmt

	if lst := c.GetLst(); lst != nil {
		add_attrs(stmt, lst.GetChildren())
	}
}

func (l *listener) EnterCommentStmt(c *parser.CommentStmtContext) {
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

func (l *listener) EnterBlankLine(c *parser.BlankLineContext) {
	l.block.AddTextNode("")
}

/****************************************************************************************/

func (l *listener) handle_conditional_statement(lst_i parser.IConditionExprContext, ctype int, ident string) {
	var (
		expr = new(ast.Expression)
		lst  = lst_i.(*parser.ConditionExprContext)
	)

	if lst.AttributeList() != nil {
		// Just treat the condition as a generic XML statement for now.
		expr.XML = ast.NewXMLStatement(ident)
		add_attrs(expr.XML, lst.AttributeList().GetChildren())
	} else {
		cond := lst.Expression()
		val := l.chs.GetText(cond.GetStart().GetStart(), cond.GetStop().GetStop())
		expr.Raw = val
	}

	stmt := l.block.AddConditionStatement(expr, ctype)
	l.cur = stmt
}

func (l *listener) EnterIfStmt(c *parser.IfStmtContext) {
	l.handle_conditional_statement(c.GetLst(), ast.ConditionIf, "do_if")
}

func (l *listener) EnterElseifStmt(c *parser.ElseifStmtContext) {
	l.handle_conditional_statement(c.GetLst(), ast.ConditionElseif, "do_elseif")
}

func (l *listener) EnterWhileStmt(c *parser.WhileStmtContext) {
	l.handle_conditional_statement(c.GetLst(), ast.ConditionWhile, "do_while")
}

func (l *listener) EnterElseStmt(c *parser.ElseStmtContext) {
	stmt := l.block.AddConditionStatement(nil, ast.ConditionElse)
	l.cur = stmt
}

//func (l *listener) ExitConditionStmt(c *parser.ConditionStmtContext) {
//	// stmt := l.block.AddXMLStatement("do_" + c.GetIdent().GetText())
//	// stmt = l.block.AddConditionStatement
//	// l.cur = stmt
//
//	if lst := c.GetLst(); lst != nil {
//		ctx := lst.(*parser.ConditionExprContext)
//		if ctx.AttributeList() != nil {
//			// Just treat the condition as a generic XML statement for now.
//			stmt := l.block.AddXMLStatement("do_" + c.GetIdent().GetText())
//			l.cur = stmt
//			add_attrs(stmt, ctx.AttributeList().GetChildren())
//		} else {
//			// Actually handle the condition properly.
//			expr := ctx.Expression()
//			val := l.chs.GetText(expr.GetStart().GetStart(), expr.GetStop().GetStop())
//			stmt.AddAttribute("value", ast.NewExpression(val))
//		}
//	} else {
//		// Can only be an `else` statement
//		stmt := l.block.AddConditionStatement(nil, ast.ConditionElse)
//	}
//}

/****************************************************************************************/

func add_attrs(stmt *ast.XMLStatement, lst []antlr.Tree) {
	for _, child := range lst {
		switch a := child.(type) {
		case parser.IAttributeContext:
			val := strings.Trim(a.GetVal().GetText(), "\"")
			expr := ast.NewExpression(val)
			stmt.AddAttribute(a.GetIdent().GetText(), expr)

		case parser.IAttributeListContext:
			add_attrs(stmt, a.GetChildren())
		}
	}
}