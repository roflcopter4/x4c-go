grammar Script;

/****************************************************************************************/
/* Lexer */

/* Keywords */
Keyword_Conditional: 'if' | 'elseif' | 'else' | 'while' ;
BuiltinFunction : 'sin' | 'cos' | 'sqrt' | 'log' | 'exp';
Keyword: 'in' | 'then' | 'nil' | 'chance';

fragment IdentHead: [a-zA-Z_];
fragment IdentChar: [a-zA-Z0-9_];
fragment SP: [ ];
fragment INT: [0-9];
fragment HEX: [0-9a-fA-F];

TextDbRef: '{' [1-9][0-9]* ',' SP* [1-9][0-9]* '}';

//Operator: '==' | '!=' | '[]'
//	| '+' | '-' | '*' | '/' | '%' | '^' | '='
//	| ';' | ':' | '!' | '@' | '?' | '.' | ','
//	| '(' | ')' | '{' | '}' | '[' | ']';

Operator: '='
	| ';' | ':' | '.' | ','
	| '(' | ')' | '{' | '}' | '[' | ']';

AdditiveOp       : '+' | '-' ;
MultiplicativeOp : '*' | '/' | '%' ;
PowerOp          : '^';
UnaryPostfixOp   : '?';
UnaryOp          : '@' | 'typeof' ; // Supposedly negation has the same precidence as other unary operators? Madness. MADNESS!
NegationOp       : 'not' | '!' ;
ComparitiveOp    : 'le' | 'ge' | 'lt' | 'gt' | '<=' | '>=' | '<' | '>' ;
EqualityOp       : '==' | '!=' ;
AndOp            : 'and' | '&&' ;
OrOp             : 'or' | '||' ;

//BuiltinFunction  : 'sqrt' ;

/* Numbers and lots of etc */
TimeValue:     INT+ ('ms' | 's' | 'min' | 'h');
DistanceValue: INT+ ('m'|'km');
CreditValue:   INT+ ('ct'|'Cr');
DegreeValue:   INT+ ('deg'|'rad');
HealthValue:   INT+ 'hp';

Float  : (INT+ '.' INT* | '.' INT+) ('f'|'LF')?
       | INT+ ('f'|'LF');

Integer: INT+ [iL]?;
SString: ['] ('\\'['] | ~['])* ['];

/* And my types... */
Variable: '$' IdentChar+;
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

//expression
//        : 'if' expression 'then' expression 'else' expression # terniary_expression
//        |<assoc=right> expression LogicalOp expression                     # logical_expression
//        |<assoc=right> expression RelationalOp expression                  # relational_expression
//        |<assoc=right> expression AdditiveOp expression                    # additive_expression
//        |<assoc=right> expression MultiplicativeOp expression              # multiplicative_expression
//        |<assoc=right> expression PowerOp expression                       # power_expression
//        |<assoc=right> BuiltinFunction '(' expression ')'                  # builtin_function_expression
//        |<assoc=right> NegationOp expression                               # negation_expression
//        |<assoc=right> expression UnaryPostfixOp              # unary_postfix_expression
//        |<assoc=right> (AdditiveOp|UnaryOp) expression                     # unary_expression
//        | object                                              # object_expression
//	;

expression
        : object                                                 # object_expression
        | (AdditiveOp|UnaryOp) expression                        # unary_expression
        |<assoc=right> expression UnaryPostfixOp                 # unary_postfix_expression
        | NegationOp expression                                  # negation_expression
        | BuiltinFunction '(' expression ')'                     # builtin_function_expression
        | expression PowerOp expression                          # power_expression
        | expression MultiplicativeOp expression                 # multiplicative_expression
        | expression AdditiveOp expression                       # additive_expression
        | expression ComparitiveOp expression                    # comparitive_expression
        | expression EqualityOp expression                       # equality_expression
        | expression AndOp expression                            # and_expression
        | expression OrOp expression                             # or_expression
        | 'if' expression 'then' expression ('else' expression)? # terniary_expression
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
	: literal
        | identifier
        | '[' (expression (',' expression)*)? ']'
        | 'table' '[' (expression (',' expression)*)? ']'
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
        | 'null'
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
