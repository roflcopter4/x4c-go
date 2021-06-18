parser grammar X4Parse;

@header {
//import (
//)

//@members {
//var x int
//}
}


/*
 * NOTE MUST BE IN THE SAME ORDER AS LISTED IN THE LEXER!!!!!!
 * Look. You're lazy. Yes, you. You know it, I know it, and you know I know it.
 * Here's a one liner to extract the things. You're welcome. It'll probably
 * miss any with the colon not on the same line as the token ID though.
 *
 * perl -wnE 'my @m = m/^([A-Z_]\w+?)\s*:/; print "@m, " if @m;' ./translation/gen/X4Lex.g4 && echo
 */

tokens {
    BuiltinFunction, TextDbRef, Float, Integer, SString, /*AdditiveOp,
    MultiplicativeOp, PowerOp, NegationOp, ComparitiveOp, EqualityOp,
    AndOp, OrOp, UnaryPostfixOp, UnaryOp, UnaryPostfixModifier, Postfix_Distance,
    Postfix_Money, Postfix_Time, Postfix_Angle, Postfix_Health, Postfix_Integer,
    Postfix_Float,*/ ATSIGN, BACKSLASH, DOLLAR, EQUALS, EXCLAM, QMARK,
    LBRACKET, RBRACKET, LBRACE, RBRACE, LPAREN, RPAREN, LANGLE, RANGLE,
    POWER, PLUS, MINUS, ASTERIX, SLASH, PERCENT, COLON, COMMA, PERIOD,
    SEMICOLON, SQUOTE, DQUOTE, TOK_CHANCE, TOK_ELSE, TOK_ELSEIF, TOK_FOR,
    TOK_FOREACH, TOK_IF, TOK_IN, TOK_LIST, TOK_NOT, TOK_NULL, TOK_TABLE,
    TOK_THEN, TOK_TYPEOF, TOK_WHILE, TOK_MIN, TOK_MAX, OP_EQ, OP_NEQ,
    OP_LE, OP_GE, OP_AND, OP_OR, TOK_AND, TOK_OR, TOK_GE, TOK_GT, TOK_LE,
    TOK_LT, TOK_MS, TOK_S, TOK_H, TOK_M, TOK_KM, TOK_CT, TOK_CR, TOK_DEG,
    TOK_RAD, TOK_HP, TOK_I, TOK_L, TOK_F, TOK_LF, Variable, BareIdentifier,
    AttributeValue, LineComment, BlockComment, Newline, Whitespace, Garbage
}

/****************************************************************************************/
/* Parser */

document:       fileTypeStmt commentStmt* EOF;
debugStatement: expression EOF;

fileTypeStmt
	: xmlStmt compoundStmt
	;

compoundStmt
	: LBRACE pseudoStatement* RBRACE
	;

pseudoStatement
	: commentStmt
	| statement
	;

statement
	: conditionStmt statementEnding
	| xmlStmt       statementEnding
	;

statementEnding
	: compoundStmt
	| SEMICOLON
	;


/*--------------------------------------------------------------------------------------*/


/* Comments. Let's pretend they're statements because I'm lazy and dumb. */
commentStmt
	: BlockComment
	| LineComment
	;

/* Generic XML statement */
xmlStmt
	: Ident=BareIdentifier LANGLE Lst=attributeList? RANGLE
	;

attributeList
	: attribute+
	;

attribute
	: Ident=specialXmlIdentifier EQUALS Val=AttributeValue
	;

specialXmlIdentifier
	: BareIdentifier (COLON BareIdentifier)?
	| keywordClash
	;

keywordClash
	: TOK_CHANCE | TOK_IN | TOK_TABLE | TOK_LIST | builtinFunction
	;

builtinFunction
	: BuiltinFunction
	| TOK_MIN | TOK_MAX
	;

/* Condition statement: if/elseif/else/while. Sanity checking the if/else chain
 * is handled in the code because I couldn't think of a way to do it here. */
conditionStmt
	: Ident=TOK_IF     Lst=conditionExpr # ifStmt
	| Ident=TOK_ELSEIF Lst=conditionExpr # elseifStmt
	| Ident=TOK_WHILE  Lst=conditionExpr # whileStmt
	| Ident=TOK_ELSE                     # elseStmt
	;

/* As a special case conditions will allow xml style statements for now. */
conditionExpr
	: LANGLE attributeList RANGLE
	| LPAREN expression RPAREN
	;


/****************************************************************************************/


expressionList
	: expression (COMMA expression)*
	;

expression
	: negationOp* unaryOp? subexpression
	;

//expression
//	: negationOp expression # exprNegation
//	| maybeUnaryExpression  # exprNEXT
//	;
//
//maybeUnaryExpression
//	: unaryOp subexpression # unaryExpression
//	| subexpression         # unaryNEXT
//	;

subexpression
	: Left=subexpression Op=powerOp          Right=expression  # subexprPower
	| Left=subexpression Op=multiplicativeOp Right=expression  # subexprMultiplicative
	| Left=subexpression Op=additiveOp       Right=expression  # subexprAdditive
	| Left=subexpression Op=comparitiveOp    Right=expression  # subexprComparitive
	| Left=subexpression Op=equalityOp       Right=expression  # subexprEquality
	| Left=subexpression Op=andOp            Right=expression  # subexprAnd
	| Left=subexpression Op=orOp             Right=expression  # subexprOr
	| TOK_IF First=expression TOK_THEN Second=expression
	  (TOK_ELSE Third=expression)?                             # subexprTerniary
	| object                                                   # subexprObject
	;


object
	: primaryObject (PERIOD secondaryObject)* unaryPostfixOp? unaryPostfixModifier?
	;

primaryObject
	: parenthetical
	| tableDefinition
	| listDefinition
	| exprBuiltinFunction
	| primaryTerminal
	;

secondaryObject
	: LBRACE expression RBRACE
	| listDefinition   // Format operations. FIXME: This is downright ugly.
	| secondaryTerminal
	;

tableDefinition
	: TOK_TABLE LBRACKET (tableAssignment (COMMA tableAssignment)* COMMA?)? RBRACKET
	;

tableAssignment
	: LBRACE object RBRACE EQUALS expression
	| Variable             EQUALS expression
	;

listDefinition
	: LBRACKET (expressionList COMMA?)? RBRACKET
	;

parenthetical
	: LPAREN expression RPAREN
	;

exprBuiltinFunction
	: builtinFunction LPAREN expressionList? RPAREN
	;

primaryTerminal
	: literal
	| keywordClash
	| identifier
	;

secondaryTerminal  /* ie. no literals */
	: keywordClash
	| identifier
	;

literal
	: SString
	| Float
	| Integer
	| TextDbRef
	| TOK_NULL
	;

identifier
	: BareIdentifier
	| Variable
	;

unaryPostfixModifier
	: (TOK_M | TOK_KM)                   # postfix_distance
	| (TOK_CR | TOK_CT)                  # postfix_money
	| (TOK_MS | TOK_S | TOK_MIN | TOK_H) # postfix_time
	| (TOK_DEG | TOK_RAD)                # postfix_angle
	| (TOK_HP)                           # postfix_health
	| (TOK_I | TOK_L)                    # postfix_integer
	| (TOK_F | TOK_LF)                   # postfix_float
	;


/*--------------------------------------------------------------------------------------*/


comparitiveOp: groupGT | groupLT | groupGE | groupLE;

additiveOp:       PLUS | MINUS;
andOp:            TOK_AND | OP_AND ;
equalityOp:       OP_EQ | OP_NEQ;
multiplicativeOp: ASTERIX | SLASH | PERCENT;
negationOp:       TOK_NOT | EXCLAM;
orOp:             TOK_OR | OP_OR;
powerOp:          POWER;
unaryOp:          PLUS | MINUS | ATSIGN | TOK_TYPEOF;
unaryPostfixOp:   QMARK;


groupGT: TOK_GT | RANGLE;
groupLT: TOK_LT | LANGLE;
groupGE: TOK_GE | OP_GE; 
groupLE: TOK_LE | OP_LE; 


/*--------------------------------------------------------------------------------------*/
// vim: tw=0 sts=0 sw=0 noet ts=8
