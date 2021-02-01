package untranslate

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/davecgh/go-spew/spew"
	script "github.com/roflcopter4/x4c-go/translation/gen/scriptparser"
)

type ScriptVisitor struct {
	*script.BaseScriptVisitor
}

type ScriptListener struct {
	*script.BaseScriptListener
	level int
}

var evil [3]string

const (
	_ = iota
	_EXPRESSION_ENTER
	_EXPRESSION_EXIT
)

func (l *ScriptListener) get_indent() int {
	return l.level * 3
}

func (l *ScriptListener) mydump(ctx antlr.ParseTree, state int) {
	// if ctx.GetChildCount() == 1 {
	// if _, ok := ctx.(antlr.ParserRuleContext); !ok {
	//       return
	// }

	var (
		prefix, color string
		colon         = "\033[0;1m:\033[0m"
	)

	output := func() {
		start := fmt.Sprintf("%T", ctx)
		start = start[14 : len(start)-7]
		space := strings.Repeat(" ", (30+26)-len(start)-(l.get_indent()))
		// fmt.Printf("%s%s%s%s%s\n", prefix, start, space, mytext(ctx))

		if _, ok := ctx.(*script.Object_expressionContext); ok {
			if state == _EXPRESSION_ENTER {
				fmt.Println(prefix + start + colon + space + mytext(ctx))
			}
		} else {
			fmt.Println(prefix + color + start + colon + space + mytext(ctx))
		}
	}

	switch state {
	case _EXPRESSION_ENTER:
		color = "\033[1;32m"
		prefix = strings.Repeat(" ", l.get_indent())
		output()
		l.level++
	case _EXPRESSION_EXIT:
		l.level--
		color = "\033[1;31m"
		prefix = strings.Repeat(" ", l.get_indent())
		output()
	default:
		panic("Impossible!")
	}

	// txt := mytext(ctx)
	// if txt != evil[:2] {
	//       fmt.Printf("%s:%s%s\n", evil[0], evil[1], evil[2])
	//       evil[2] = txt
	// }
	// evil[0] = fmt.Sprintf("%T", ctx)
	// evil[1] = strings.Repeat(" ", 48-len(evil[0]))
}

func mytext(prc antlr.ParseTree) string {
	var s string

	if prc.GetChildCount() == 0 {
		s = prc.GetText() + " "
	} else {
		for _, child := range prc.GetChildren() {
			s += mytext(child.(antlr.ParseTree))

			// if x, ok := child.(*script.MyterminalContext); ok {
			//       s += x.GetText() + " "
			// } else {
			//       s += mytext(child.(antlr.ParseTree))
			// }
		}
	}

	return s
}

func (l *ScriptListener) EnterStart(ctx *script.StartContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterStatement(ctx *script.StatementContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterSimple_statement(ctx *script.Simple_statementContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterPredicate(ctx *script.PredicateContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterAssignment_statement(ctx *script.Assignment_statementContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterExpression(ctx *script.ExpressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterObject_expression(ctx *script.Object_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterUnary_expression(ctx *script.Unary_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterUnary_postfix_expression(ctx *script.Unary_postfix_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterTerniary_expression(ctx *script.Terniary_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterMultiplicative_expression(ctx *script.Multiplicative_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterAdditive_expression(ctx *script.Additive_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterNegation_expression(ctx *script.Negation_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterComparitive_expression(ctx *script.Comparitive_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterEquality_expression(ctx *script.Equality_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterAnd_expression(ctx *script.And_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterOr_expression(ctx *script.Or_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterPower_expression(ctx *script.Power_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}
func (l *ScriptListener) EnterBuiltin_function_expression(ctx *script.Builtin_function_expressionContext) {
	l.mydump(ctx, _EXPRESSION_ENTER)
}

func (l *ScriptListener) ExitStart(ctx *script.StartContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitStatement(ctx *script.StatementContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitSimple_statement(ctx *script.Simple_statementContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitPredicate(ctx *script.PredicateContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitAssignment_statement(ctx *script.Assignment_statementContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitExpression(ctx *script.ExpressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitObject_expression(ctx *script.Object_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitUnary_expression(ctx *script.Unary_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitUnary_postfix_expression(ctx *script.Unary_postfix_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitTerniary_expression(ctx *script.Terniary_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitMultiplicative_expression(ctx *script.Multiplicative_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitAdditive_expression(ctx *script.Additive_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitNegation_expression(ctx *script.Negation_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitComparitive_expression(ctx *script.Comparitive_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitEquality_expression(ctx *script.Equality_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitAnd_expression(ctx *script.And_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitOr_expression(ctx *script.Or_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitPower_expression(ctx *script.Power_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}
func (l *ScriptListener) ExitBuiltin_function_expression(ctx *script.Builtin_function_expressionContext) {
	l.mydump(ctx, _EXPRESSION_EXIT)
}

//
//func (v *ScriptListener) EnterParen_expression(ctx *script.Paren_expressionContext) { mydump(ctx) }
//
//func (v *ScriptListener) EnterExpression(ctx *script.ExpressionContext)                 { mydump(ctx) }
//func (v *ScriptListener) EnterPrimary_expression(ctx *script.Primary_expressionContext) { mydump(ctx) }
//
//// func (v *ScriptListener) EnterUnary_expression(ctx *script.Unary_expressionContext)     { mydump(ctx) }
//func (v *ScriptListener) EnterIdentifier(ctx *script.IdentifierContext)                 { mydump(ctx) }
//func (v *ScriptListener) EnterMyterminal(ctx *script.MyterminalContext)                 { mydump(ctx) }
//func (v *ScriptListener) EnterLiteral(ctx *script.LiteralContext)                       { mydump(ctx) }
//func (v *ScriptListener) EnterPostfix_expression(ctx *script.Postfix_expressionContext) { mydump(ctx) }
//func (v *ScriptListener) EnterPostfix(ctx *script.PostfixContext)                       { mydump(ctx) }
//func (v *ScriptListener) EnterDotpostfix(ctx *script.DotpostfixContext)                 { mydump(ctx) }

//func (v *ScriptListener) EnterAdditiveOp(ctx *script.AdditiveOpContext)             { mydump(ctx) }
//func (v *ScriptListener) EnterMultiplicativeOp(ctx *script.MultiplicativeOpContext) { mydump(ctx) }
//func (v *ScriptListener) EnterUnaryPostfixOp(ctx *script.UnaryPostfixOpContext)     { mydump(ctx) }

// func (v *ScriptListener) EnterUnaryOp(ctx *script.UnaryOpContext) { mydump(ctx) }

//func (v *ScriptListener) EnterRelationalOp(ctx *script.RelationalOpContext)         { mydump(ctx) }
//func (v *ScriptListener) EnterLogicalOp(ctx *script.LogicalOpContext)               { mydump(ctx) }
//func (v *ScriptListener) EnterNegationOp(ctx *script.NegationOpContext)              { mydump(ctx) }

// func (v *ScriptListener) EnterBuiltinFunction(ctx *script.BuiltinFunctionContext)   { mydump(ctx) }

func TestScriptLexer(str string, isfile bool) {
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

	var charstream antlr.CharStream

	if isfile {
		if str == "" || str == "-" {
			b, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}
			charstream = antlr.NewInputStream(string(b))
		} else {
			fs, err := antlr.NewFileStream(str)
			if err != nil {
				panic(err)
			}
			charstream = fs
		}
	} else {
		charstream = antlr.NewInputStream(str)
	}

	lexer := script.NewScriptLexer(charstream)
	// lexer := script.NewScript(charstream)

	//for {
	//	t := lexer.NextToken()
	//	if t.GetTokenType() == antlr.TokenEOF {
	//		break
	//	}
	//	fmt.Printf("%s (%q)\n",
	//		lexer.SymbolicNames[t.GetTokenType()], t.GetText())
	//}

	var (
		stream = antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
		p      = script.NewScriptParser(stream)
	)
	// p.Expression()

	// p.Vis

	vis := new(ScriptListener)

	spew.Config.MaxDepth = 1
	// vis.VisitExpression(p.Expression().(*script.ExpressionContext))
	antlr.ParseTreeWalkerDefault.Walk(vis, p.Start())

	// if evil[2] != "" {
	//       fmt.Printf("%s:%s%s\n", evil[0], evil[1], evil[2])
	// }

	spew.Config.MaxDepth = 0

	// antlr.ParseTreeVisitorDefault.Walk(nil, p.Document())
	// script.NewScriptV
	//
	// return l.a
}
