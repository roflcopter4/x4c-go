grammar X4C;

/****************************************************************************************/
/* Tokens */

/* Keywords */
Keyword_If:     'if';
Keyword_Elseif: 'elseif';
Keyword_Else:   'else';
Keyword_Logic: 'ge' | 'le' | 'gt' | 'lt' | 'and' | 'or' | 'not';
Keyword_Builtin: 'sqrt';

Keyword: 'in' | 'then' | 'nil';

fragment IdentHead: [a-zA-Z_];
fragment IdentChar: [a-zA-Z0-9_];
fragment SP: [ ];
fragment INT: [0-9];

/* These are the types formally defined in the schemas. */

///* Clear as mud, innit? */
////        { ( Initial char     ) | ( String - may start with # ) }  { ( Non string character              ) | ( Another string?  dot??    ) }/Repeat
//Expression: ( [A-Za-z0-9_$@+\-({[] | ('#'?['] (~['\\] | '.')* [']) )  ( [A-Za-z0-9_$!?@=<>;,.+\-*/%^(){}[\] ] | ('#'?['] (~['\\] | '.')* [']) )*;

Variable: '$' IdentChar+;
TextDbRef: '{' [1-9][0-9]* ',' SP* [1-9][0-9]* '}';

Operator: '==' | '!=' | '[]'
	| '+' | '-' | '*' | '/' | '%' | '^' | '!' | '(' | ')' | '{' | '}' | ';' | '.'
	| '@' | '[' | ']' | '?' | '=' | ',' | ':';

/* Stupidity: '(<' | '>)'; */
DumbExpr: '(' (DumbExpr | ~[()])+ ')';

/* Numbers and lots of etc */
TimeValue:     INT+ ([mM]'in' | [mM]'s' | [hH] | [sS]);
DistanceValue: INT+ ('m' | [kK]'m');
CreditValue:   INT+ [cC]'r';
DegreeValue:   INT+ [dD]'eg';
HealthValue:   INT+ [hH]'p';

Float: INT+ '.' INT* ([lL]?[fF])?;
Integer: INT+;
SString: ['] ('\\'['] | ~['])* ['];

/* And my types... */
Identifier: IdentHead IdentChar*;

AttributeValue: '"' (~'"' | '\\"')* '"';

LineComment:  '//' .*? Newline;
BlockComment: '/*' .*? '*/';

Newline:    ('\r\n' | '\n') -> skip;
Whitespace: [\t ]+          -> skip;


/****************************************************************************************/
/* Grammar */

document
	: fileTypeStmt commentStmt* EOF
	;

fileTypeStmt
	: xmlStmt compoundStmt
	;

maybeCompoundStmt
	: ';'
	| compoundStmt
	;

compoundStmt
	: '{' statement* '}'
	;

statement
	: commentStmt
	| conditionStmt compoundStmt
	| xmlStmt       maybeCompoundStmt
	;

/* Comments. Let's pretend they're statements because I'm lazy and dumb. */
commentStmt
	: BlockComment
	| LineComment
	;

/* Generic XML statement */
xmlStmt
	: Ident=Identifier '<' Lst=attributeList '>'
	;

//attributeList
//	: attribute*
//	;

attributeList
    : attributeList ',' attribute
    | attribute
    | /* empty */
    ;

attribute
	: Ident=specialXmlIdentifier '=' Val=AttributeValue
	;

specialXmlIdentifier
	: Identifier (':' Identifier)?
	;

/* Condition statement: if/elseif/else. Sanity checking the if/else chain is
 * handled in the code because I couldn't think of a way to do it here. */
conditionStmt
	: Ident='if'     Lst=conditionExpr
	| Ident='elseif' Lst=conditionExpr
	| Ident='else'
	;

conditionExpr
	: '<' attributeList '>'
	| DumbExpr { fmt.Printf("%v\n", $DumbExpr); }
	;


//conditionStmt
//	: Ident='if'     '(' Lst=conditionExpr ')'
//	| Ident='elseif' '(' Lst=conditionExpr ')'
//	| Ident='else'
//	;

//conditionExpr
//	:  '<' Lst=attributeList '>'
//	| Expr=expression {
//	    //cur := $Expr.ctx.(antlr.ParseTree)
//	    //for cur.GetChildCount() > 0 {
//	    //    cur = cur.GetChild(0).(antlr.ParseTree)
//	    //    fmt.Printf("%#T  ", cur)
//            //}
//	    fmt.Printf("got:\t< %v >\n", $Expr.text);
//	}
//	;

/*======================================================================================*/
/* Attempt number 1 at doing actual expressions. This will fail. */

//expression_list
//	: expression_list ',' expression
//	| expression
//	;
//
//expression
//	: builtin_function '(' expression ')'
//	| primary_expression
//	;
//
//primary_expression
//	: primary_expression logical_op relational_expression
//	| relational_expression
//	;
//
//relational_expression
//	: relational_expression relational_op additive_expression
//	| additive_expression
//	;
//
//additive_expression 
//	: additive_expression additive_op multiplicative_expression
//	| multiplicative_expression
//	;
//
//multiplicative_expression
//	: multiplicative_expression multiplicative_op unary_expression
//	| unary_expression
//	;
//
///* terniary_expression        */
///*         | unary_expression */
///*         ;                  */
//
//unary_expression
//	: unary_op unary_expression
//	| tern
//	/* | identifier_clash */
//	;
//
//tern
//	: 'if' expression 'then' expression 'else' expression
//	/* : TOK_IF '(' expression ')' "then" '(' expression ')' "else" '(' expression ')' */
//	| identifier
//	;
//
//identifier
//	: identifier '.' terminal
//	| terminal
//	;
//
//terminal
//	: '(' expression ')'
//	| '[' expression_list ']'
//	| '{' expression_list '}'
//	//| table_expression
//	//| list_expression
//	| identifier_terminal
//	| literal
//	| terminal post_op
//	;
//
//literal
//	: SString
//	| Float
//	| Integer
//	//| TOK_NULL { $$ = b_fromlit("null"); }
//	| DistanceValue
//	| DegreeValue
//	| HealthValue
//	| TimeValue
//	| CreditValue
//	//| TOK_EMPTY_ARRAY
//	;
//
//identifier_terminal
//	: Identifier
//	| Variable
//	//| TOK_CONST
//	//| identifier_clash
//	;
//
//post_op: '?';
//
//builtin_function  : 'sqrt' ;
//additive_op       : '+' | '-' ;
//unary_op          : '+' | '-' | '@' | '!' | 'typeof' | 'not' ;
//multiplicative_op : '^' | '*' | '/' | '%' ;
//relational_op     : '==' | '!=' | 'le' | 'ge' | 'lt' | 'gt' ;
//logical_op        : 'and' | 'or'  ;
