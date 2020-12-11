grammar Script;

/****************************************************************************************/
/* Lexer */

/* Keywords */
//Keyword_If:     'if';
//Keyword_Elseif: 'elseif';
//Keyword_Else:   'else';
//Keyword_Logic: 'ge' | 'le' | 'gt' | 'lt' | 'and' | 'or' | 'not';
//Keyword_Builtin: 'sqrt';

Keyword : 'in' | 'then' | 'nil'
        //| 'ge' | 'le' | 'gt' | 'lt' | 'and' | 'or' | 'not'
        | 'if' | 'elseif' | 'else'
        | 'chance'
        ;

fragment IdentHead: [a-zA-Z_];
fragment IdentChar: [a-zA-Z0-9_];
fragment SP: [ ];
fragment INT: [0-9];
fragment HEX: [0-9a-fA-F];

Variable: '$' IdentChar+;
TextDbRef: '{' [1-9][0-9]* ',' SP* [1-9][0-9]* '}';

//Operator: '==' | '!=' | '[]'
//	| '+' | '-' | '*' | '/' | '%' | '^' | '='
//	| ';' | ':' | '!' | '@' | '?' | '.' | ','
//	| '(' | ')' | '{' | '}' | '[' | ']';

Operator: '[]'
	| '='
	| ';' | ':' | '.' | ','
	| '(' | ')' | '{' | '}' | '[' | ']';

AdditiveOp       : '+' | '-' ;
MultiplicativeOp : '^' | '*' | '/' | '%' ;
UnaryPostfixOp   : '?';
UnaryOp          : '@' | 'typeof' ;
NegationOp       : 'not' | '!' ;
RelationalOp     : '==' | '!=' | 'le' | 'ge' | 'lt' | 'gt' | '<=' | '>=' | '<' | '>' ;
LogicalOp        : 'and' | 'or' | '&&' | '||' ;

//BuiltinFunction  : 'sqrt' ;

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

LineComment:  '//' .*? Newline;
BlockComment: '/*' .*? '*/';

Newline:    ('\r\n' | '\n') -> skip;
Whitespace: [\t ]+          -> skip;

/****************************************************************************************/
/* Parser */

start : statement EOF;

statement
        :<assoc=right> object '=' expression  # assignment_statement
        | expression predicate?               # simple_statement
        ;

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

//additiveOp       : '+' | '-' ;
//multiplicativeOp : '^' | '*' | '/' | '%' ;
//unaryPostfixOp   : '?';
//unaryOp          : '+' | '-' | '@' | 'typeof' ;
//relationalOp     : '==' | '!=' | 'le' | 'ge' | 'lt' | 'gt'
//                 | '<=' | '>=' | '<' | '>' ;
//logicalOp        : 'and' | 'or' | '&&' | '||' ;
/* negationOp       : 'not' | '!' ; */

//builtinFunction  : 'sqrt' ;
