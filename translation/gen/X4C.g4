grammar X4C;

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
fragment FLOAT: (((INT+ '.' INT* | '.' INT+)('e' INT+)?) | (INT+ 'e' INT+));

/* These are the types formally defined in the schemas. */

///* Clear as mud, innit? */
////        { ( Initial char     ) | ( String - may start with # ) }  { ( Non string character              ) | ( Another string?  dot??    ) }/Repeat
//Expression: ( [A-Za-z0-9_$@+\-({[] | ('#'?['] (~['\\] | '.')* [']) )  ( [A-Za-z0-9_$!?@=<>;,.+\-*/%^(){}[\] ] | ('#'?['] (~['\\] | '.')* [']) )*;

TextDbRef: '{' [1-9][0-9]* ',' SP* [1-9][0-9]* '}';

Operator: '='
	| ';' | ':' | '.' | ','
	| '(' | ')' | '{' | '}' | '[' | ']';

AdditiveOp       : '+' | '-' ;
MultiplicativeOp : '*' | '/' | '%' ;
PowerOp          : '^';
UnaryPostfixOp   : '?';
UnaryOp          : '@' | 'typeof' ;
NegationOp       : 'not' | '!' ;
ComparitiveOp    : 'le' | 'ge' | 'lt' | 'gt' | '<=' | '>=' | '<' | '>' ;
EqualityOp       : '==' | '!=' ;
AndOp            : 'and' | '&&' ; // I may or may not allow 'C' style logical operators.
OrOp             : 'or' | '||' ;

/* Numbers and lots of etc */
TimeValue:     (INT+ | FLOAT) [ ]* ('ms' | 's' | 'min' | 'h');
DistanceValue: (INT+ | FLOAT) [ ]* ('m'|'km');
CreditValue:   (INT+ | FLOAT) [ ]* ('ct'|'Cr');
DegreeValue:   (INT+ | FLOAT) [ ]* ('deg'|'rad');
HealthValue:   (INT+ | FLOAT) [ ]* 'hp';

Float  : FLOAT [ ]* ('f'|'LF')?
       | INT+ [ ]* ('f'|'LF');

Integer: INT+ [iL]?;
SString: ['] ('\\'['] | ~['])* ['];

/* And my types... */
Variable: '$' IdentChar+;
BareIdentifier: IdentHead IdentChar*;
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

compoundStmt
	: '{' statement* '}'
	;

statement
	: commentStmt
	| conditionStmt statement
	| xmlStmt       statement
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
	: Ident=BareIdentifier '<' Lst=attributeList? '>'
	;

attributeList
        : attribute+
//      : attributeList attribute
//      | attribute              
        ;

attribute
	: Ident=specialXmlIdentifier '=' Val=AttributeValue
	;

specialXmlIdentifier
	: BareIdentifier (':' BareIdentifier)?
        | keywordClash
	;

keywordClash
        : 'chance' | 'in' | 'table'
        ;

/* Condition statement: if/elseif/else/while. Sanity checking the if/else chain
 * is handled in the code because I couldn't think of a way to do it here. */
conditionStmt
	: Ident='if'     Lst=conditionExpr # ifStmt
	| Ident='elseif' Lst=conditionExpr # elseifStmt
        | Ident='while'  Lst=conditionExpr # whileStmt
	| Ident='else'                     # elseStmt
	;

/* As a special case conditions will allow xml style statements for now. */
conditionExpr
	: '<' attributeList '>'
        | '(' expression ')'
	;

/****************************************************************************************/

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

//predicate
//        : 'chance' expression
//        ;

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
        : '[' (expression (',' expression)*)? ']'
        | 'table' '[' (table_assignment (',' table_assignment)*)? ']'
        | identifier
	| literal
	;

table_assignment
        : '{' literal '}' '=' expression
        | Variable '=' expression
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
	: BareIdentifier
	| Variable
	;

// vim: tw=0
