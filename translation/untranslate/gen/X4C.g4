grammar X4C;

/****************************************************************************************/
/* Lexer */

/* Keywords */
Keyword_If:     'if';
Keyword_Elseif: 'elseif';
Keyword_Else:   'else';
//Keyword_Logic: 'ge' | 'le' | 'gt' | 'lt' | 'and' | 'or' | 'not';
Keyword_Builtin: 'sqrt';

Keyword: 'in' | 'then' | 'nil' | 'chance';

fragment IdentHead: [a-zA-Z_];
fragment IdentChar: [a-zA-Z0-9_];
fragment SP: [ ];
fragment INT: [0-9];
fragment HEX: [0-9a-fA-F];

/* These are the types formally defined in the schemas. */

///* Clear as mud, innit? */
////        { ( Initial char     ) | ( String - may start with # ) }  { ( Non string character              ) | ( Another string?  dot??    ) }/Repeat
//Expression: ( [A-Za-z0-9_$@+\-({[] | ('#'?['] (~['\\] | '.')* [']) )  ( [A-Za-z0-9_$!?@=<>;,.+\-*/%^(){}[\] ] | ('#'?['] (~['\\] | '.')* [']) )*;

Variable: '$' IdentChar+;
TextDbRef: '{' [1-9][0-9]* ',' SP* [1-9][0-9]* '}';

Operator: '[]'
	| '='
	| ';' | ':' | '.' | ','
	| '(' | ')' | '{' | '}' | '[' | ']';

AdditiveOp       : '+' | '-' ;
MultiplicativeOp : '^' | '*' | '/' | '%' ;
UnaryPostfixOp   : '?';
UnaryOp          : '@' | 'typeof' ;
NegationOp       : 'not' | '!' ;
RelationalOp     : '==' | '!=' | 'le' | 'ge' | 'lt' | 'gt' /*| '<=' | '>=' | '<' | '>'*/ ;
LogicalOp        : 'and' | 'or' | '&&' | '||' ;

//DumbExpr: '(' (DumbExpr | ~[()])+ ')';

/* Numbers and lots of etc */
TimeValue:     INT+ ([mM]'in' | [mM]'s' | [hH] | [sS]);
DistanceValue: INT+ ('m' | [kK]'m');
CreditValue:   INT+ [cC]'r';
DegreeValue:   INT+ [dD]'eg';
HealthValue:   INT+ [hH]'p';

Float  : (INT+ '.' INT* | '.' INT+) ([lL]?[fF])?
       | INT+ [lL]?[fF];

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
/* Parser */

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
	| conditionStmt statement //compoundStmt
	| xmlStmt       statement //maybeCompoundStmt
        | compoundStmt
        | ';'
	;

/* Comments. Let's pretend they're statements because I'm lazy and dumb. */
commentStmt
	: BlockComment
	| LineComment
	;

/* Generic XML statement */
xmlStmt
	: Ident=Identifier '<' Lst=attributeList? '>'
	;

attributeList
    : attributeList attribute
    | attribute
    ;

attribute
	: Ident=specialXmlIdentifier '=' Val=AttributeValue
	;

specialXmlIdentifier
	: Identifier (':' Identifier)?
        | keywordClash
	;

keywordClash
        : 'chance' | 'in'
        ;

/* Condition statement: if/elseif/else. Sanity checking the if/else chain is
 * handled in the code because I couldn't think of a way to do it here. */
conditionStmt
	: Ident='if'     Lst=conditionExpr
	| Ident='elseif' Lst=conditionExpr
	| Ident='else'
	;

/* As a special case if/elseif/else will allow xml style statements for now. */
conditionExpr
	: '<' attributeList '>'
	//| DumbExpr /* FIXME: This sucks */
        | '(' expression ')'
	;

/****************************************************************************************/

expression
        : object                                              # object_expression
        | (AdditiveOp|UnaryOp) expression                     # unary_expression
        |<assoc=right> expression UnaryPostfixOp              # unary_postfix_expression
        | 'if' expression 'then' expression 'else' expression # terniary_expression
        | expression MultiplicativeOp expression              # multiplicative_expression
        | expression AdditiveOp expression                    # additive_expression
        | NegationOp expression                               # negation_expression
        | expression RelationalOp expression                  # relational_expression
        | expression LogicalOp expression                     # logical_expression
	;

predicate
        : 'chance' expression
        ;

object
        : first_obj_fragment ('.' obj_fragment)*
        | literal
        ;

first_obj_fragment
        : '(' expression ')' 
        | myterminal
        ;

obj_fragment
        : first_obj_fragment
        | '{' expression '}'
        ;

myterminal
        : identifier
	//| '[' expression_list ']'
	//| table_expression
	//| list_expression
	| literal
	;

literal
	: SString
	| Float
	| Integer
	| DistanceValue
	| DegreeValue
	| HealthValue
	| TimeValue
	| CreditValue
        | TextDbRef
	;

identifier
	: Identifier
	| Variable
	;


