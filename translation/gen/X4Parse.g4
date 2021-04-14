parser grammar X4Parse;

//@header {
//import (
//    "fmt"
//)
//}
//
//@members {
//var x int
//}

//tokens {
//    TextDbRef, PostfixTime, PostfixDistance, PostfixMoney, PostfixAngle,
//    PostfixHealth, PostfixInteger, PostfixFloat, Integer, SString, TOK_ELSE,
//    TOK_ELSEIF, TOK_FOR, TOK_FOREACH, TOK_IF, TOK_LIST, TOK_NOT, TOK_TABLE,
//    TOK_THEN, TOK_WHILE, OP_EQ, OP_NEQ, OP_LE, OP_GE, OP_AND, OP_OR, TOK_AND,
//    TOK_OR, TOK_GE, TOK_GT, TOK_LE, TOK_LT, ATSIGN, BACKSLASH, DOLLAR, EQUALS,
//    EXCLAM, QMARK, LBRACKET, RBRACKET, LBRACE, RBRACE, LPAREN, RPAREN, RANGLE,
//    LANGLE, POWER, PLUS, MINUS, ASTERIX, SLASH, PERCENT, COLON, COMMA, PERIOD,
//    SEMICOLON, SQUOTE, DQUOTE, Variable, BareIdentifier, AttributeValue,
//    LineComment, BlockComment, Newline, Whitespace, TOK_NULL, TOK_CHANCE, TOK_IN,
//    TOK_MIN, TOK_MAX, BuiltinFunction, Float, TOK_TYPEOF
//}

tokens {
    BuiltinFunction, TextDbRef, PostfixTime, PostfixDistance, PostfixMoney,
    PostfixAngle, PostfixHealth, PostfixInteger, PostfixFloat, Float, Integer, SString,
    TOK_CHANCE, TOK_ELSE, TOK_ELSEIF, TOK_FOR, TOK_FOREACH, TOK_IF, TOK_IN,
    TOK_NOT, TOK_NULL, TOK_TABLE, TOK_THEN, TOK_TYPEOF, TOK_WHILE, TOK_MIN, TOK_MAX,
    OP_EQ, OP_NEQ, OP_LE, OP_GE, OP_AND, OP_OR, TOK_AND, TOK_OR, TOK_GE, TOK_GT, TOK_LE,
    TOK_LT, ATSIGN, BACKSLASH, DOLLAR, EQUALS, EXCLAM, QMARK, LBRACKET, RBRACKET, LBRACE,
    RBRACE, LPAREN, RPAREN, LANGLE, RANGLE, POWER, PLUS, MINUS, ASTERIX, SLASH, PERCENT,
    COLON, COMMA, PERIOD, SEMICOLON, SQUOTE, DQUOTE, Variable, BareIdentifier,
    AttributeValue, LineComment, BlockComment, Newline, Whitespace
}

/****************************************************************************************/
/* Parser */

document
	: fileTypeStmt commentStmt* EOF
	;

debugStatement: expression EOF ;

fileTypeStmt
	: xmlStmt compoundStmt
	;

compoundStmt
	: RBRACE statement* LBRACE
	;

statement
	: commentStmt
	| conditionStmt statement
	| xmlStmt       statement
        | compoundStmt
        | SEMICOLON
	;

//blankLine
//        : BlankLine
//        ;

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
//      : attributeList attribute
//      | attribute              
        ;

attribute
	: Ident=specialXmlIdentifier EQUALS Val=AttributeValue
	;

specialXmlIdentifier
	: BareIdentifier (COLON BareIdentifier)?
        | keywordClash
	;

keywordClash
        : TOK_CHANCE | TOK_IN | TOK_TABLE | BuiltinFunction
        | TOK_MIN | TOK_MAX
//        | 'ms' | 's' | 'min' | 'h' | 'm' | 'km' | 'ct' | 'Cr' | 'deg' | 'rad' | 'hp'
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
	: RANGLE attributeList LANGLE
        | LPAREN expression RPAREN
	;

/****************************************************************************************/



expression
        : unaryOp expression                                           # subexpr_unary
        //|<assoc=right> expression unaryPostfixOp                       # subexpr_unary_postfix
        //|<assoc=right> expression unaryPostfixModifier                 # subexpr_unary_postfix_modifier
        | negationOp expression                                        # subexpr_negation
        //| BuiltinFunction LPAREN expression RPAREN                     # subexpr_builtin_function
        | expression powerOp expression                                # subexpr_power
        | expression multiplicativeOp expression                       # subexpr_multiplicative
        | expression additiveOp expression                             # subexpr_additive
        | expression comparitiveOp expression                          # subexpr_comparitive
        | expression equalityOp expression                             # subexpr_equality
        | expression andOp expression                                  # subexpr_and
        | expression orOp expression                                   # subexpr_or
        | TOK_IF expression TOK_THEN expression (TOK_ELSE expression)? # subexpr_terniary
        | object_expr_unary_postfix                                    # subexpr_object
	;


object_expr_unary_postfix
        : object_expr_unary_postfix_modifier unaryPostfixOp
        | object_expr_unary_postfix_modifier
        ;

object_expr_unary_postfix_modifier
        : object unaryPostfixModifier
        | object
        ;


object
        : primary_object (PERIOD secondary_object)*    # expr_primary_object
        | BuiltinFunction LPAREN expression RPAREN     # expr_builtin_function
        ;

primary_object
        : parenthetical
        | table_definition
        | list_definition
        | simple_terminal
        ;

secondary_object
        : LBRACE expression RBRACE
        | list_definition   // Format operations. FIXME: This is downright ugly.
        | simple_terminal
        ;

table_definition
        : TOK_TABLE LBRACKET (table_assignment (COMMA table_assignment)* COMMA?)? RBRACKET
        ;

table_assignment
        : LBRACE object RBRACE EQUALS expression
        | Variable             EQUALS expression
        ;

list_definition
        : LBRACKET (expression (COMMA expression)* COMMA?)? RBRACKET
        ;

parenthetical
	: LPAREN expression RPAREN
        ;

simple_terminal
        : literal
        | keywordClash
        | identifier
//        | blankLine  // FIXME: UGH
        ;

literal
	: SString
	| Float
	| Integer
//	| DistanceValue
//	| DegreeValue
//	| HealthValue
//	| TimeValue
//	| CreditValue
        | TextDbRef
        | TOK_NULL
	;

identifier
	: BareIdentifier
	| Variable
	;

//unary_postfix_modifier_expression_impl
//        :<assoc=right> expression PostfixTime     # postfix_time_expression
//        |<assoc=right> expression PostfixDistance # postfix_distance_expression
//        |<assoc=right> expression PostfixMoney    # postfix_money_expression
//        |<assoc=right> expression PostfixAngle    # postfix_angle_expression
//        |<assoc=right> expression PostfixHealth   # postfix_health_expression
//        ;

//        : ('ms' | 's' | 'min' | 'h') # postfix_time
//        | ('m' | 'km')               # postfix_distance
//        | ('ct' | 'Cr')              # postfix_money
//        | ('deg' | 'rad')            # postfix_angle
//        | 'hp'                       # postfix_hp

unaryPostfixModifier
        : PostfixTime
        | PostfixDistance
        | PostfixMoney
        | PostfixAngle
        | PostfixHealth
        | PostfixInteger
        | PostfixFloat
        ;


additiveOp: PLUS | MINUS;
multiplicativeOp: ASTERIX | SLASH | PERCENT ;
powerOp: POWER ;
unaryPostfixOp: QMARK ;
unaryOp: PLUS | MINUS | ATSIGN | TOK_TYPEOF ;
negationOp: TOK_NOT | EXCLAM ;
comparitiveOp : rule_gt | rule_lt | rule_le | rule_ge;
equalityOp: OP_EQ | OP_NEQ ;
andOp: TOK_AND | OP_AND ;
orOp: TOK_OR | OP_OR ;


rule_gt: TOK_GT | RANGLE;
rule_lt: TOK_LT | LANGLE;
rule_ge: TOK_GE | OP_GE;
rule_le: TOK_LE | OP_LE;

//        : ('ms' | 's' | 'min' | 'h') # postfix_time
//        | ('m' | 'km')               # postfix_distance
//        | ('ct' | 'Cr')              # postfix_money
//        | ('deg' | 'rad')            # postfix_angle
//        | 'hp'                       # postfix_hp
//        ;

// vim: tw=0
