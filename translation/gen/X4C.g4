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
ComparitiveOp    : 'le' | 'ge' | 'lt' | 'gt'  | '<' | '>' | '<=' | '>=' ;
EqualityOp       : '==' | '!=' ;
AndOp            : 'and' | '&&' ; // I may or may not allow 'C' style logical operators.
OrOp             : 'or' | '||' ;

/* Numbers and lots of etc */
TimeValue:     (INT+ | FLOAT) [ ]* ('ms' | 's' | 'min' | 'h');
DistanceValue: (INT+ | FLOAT) [ ]* ('m'|'km');
CreditValue:   (INT+ | FLOAT) [ ]* ('ct'|'Cr');
DegreeValue:   (INT+ | FLOAT) [ ]* ('deg'|'rad');
HealthValue:   (INT+ | FLOAT) [ ]* 'hp';

PostfixTime:     ('ms' | 's' | 'min' | 'h');
PostfixDistance: ('m' | 'km');
PostfixMoney:    ('ct' | 'Cr');
PostfixAngle:    ('deg' | 'rad');
PostfixHealth:   'hp';

/* UnaryPostfixModifier: 'ms' | 's' | 'min' | 'h' | 'm' | 'km' | 'ct' | 'Cr' | 'deg' | 'rad' | 'hp'; */

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

BlankLine: Newline Whitespace* Newline;

Newline:    ('\r\n' | '\n') -> skip;
Whitespace: [\t ]+          -> skip;


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
	: '{' statement* '}'
	;

statement
	: blankLine
        | commentStmt
	| conditionStmt statement
	| xmlStmt       statement
        | compoundStmt
        | ';'
	;

blankLine
        : BlankLine
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
        //| 'ms' | 's' | 'min' | 'h' | 'm' | 'km' | 'ct' | 'Cr' | 'deg' | 'rad' | 'hp'
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
        : object                                                 # expr_object
        | (AdditiveOp|UnaryOp) expression                        # expr_unary
        |<assoc=right> expression UnaryPostfixOp                 # expr_unary_postfix
//        | unary_postfix_modifier_expression_impl                 # expr_unary_postfix_modifier
        |<assoc=right> expression unaryPostfixModifier           # expr_unary_postfix_modifier
        | NegationOp expression                                  # expr_negation
        | BuiltinFunction '(' expression ')'                     # expr_builtin_function
        | expression PowerOp expression                          # expr_power
        | expression MultiplicativeOp expression                 # expr_multiplicative
        | expression AdditiveOp expression                       # expr_additive
        | expression comparitiveOp expression                    # expr_comparitive
        | expression EqualityOp expression                       # expr_equality
        | expression AndOp expression                            # expr_and
        | expression OrOp expression                             # expr_or
        | 'if' expression 'then' expression ('else' expression)? # expr_terniary
	;

//predicate
//        : 'chance' expression
//        ;

object
        : first_obj_fragment ('.' obj_fragment)*
        ;

first_obj_fragment
        : '(' expression ')' 
        | base_object
        ;

obj_fragment
        : first_obj_fragment
        | '{' expression '}'
        ;

base_object
        : '[' (expression (',' expression)*)? ']'                     # square_bracket_expr
        | 'table' '[' (table_assignment (',' table_assignment)*)? ']' # table_definition
        | (identifier | literal)                                      # base_token
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
        ;

comparitiveOp : 'le' | 'ge' | 'lt' | 'gt'  | '<' | '>' | '<=' | '>=' ;

//        : ('ms' | 's' | 'min' | 'h') # postfix_time
//        | ('m' | 'km')               # postfix_distance
//        | ('ct' | 'Cr')              # postfix_money
//        | ('deg' | 'rad')            # postfix_angle
//        | 'hp'                       # postfix_hp
//        ;

// vim: tw=0
