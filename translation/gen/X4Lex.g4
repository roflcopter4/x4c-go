lexer grammar X4Lex;

/****************************************************************************************/
/* Lexer */

/* Keywords */
//Keyword_Conditional: 'if' | 'elseif' | 'else' | 'while' ;
BuiltinFunction : 'sin' | 'cos' | 'sqrt' | 'log' | 'exp';

//Keyword: 'in' | 'then' | 'nil' | 'chance';

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

//Operator: '='
//	| ';' | ':' | '.' | ','
//	| '(' | ')' | '{' | '}' | '[' | ']';

/* TOK_TYPEOF: 'typeof'; */

//AdditiveOp       : '+' | '-' ;
//MultiplicativeOp : '*' | '/' | '%' ;
//PowerOp          : '^';
//UnaryPostfixOp   : '?';
//UnaryOp          : '@' | 'typeof' ;
//NegationOp       : 'not' | '!' ;
//ComparitiveOp    : 'le' | 'ge' | 'lt' | 'gt'  | '<' | '>' | '<=' | '>=' ;
//EqualityOp       : '==' | '!=' ;
//AndOp            : 'and' | '&&' ; // I may or may not allow 'C' style logical operators.
//OrOp             : 'or' | '||' ;

/* Numbers and lots of etc */
//TimeValue:     (INT+ | FLOAT) [ ]* ('ms' | 's' | 'min' | 'h');
//DistanceValue: (INT+ | FLOAT) [ ]* ('m'|'km');
//CreditValue:   (INT+ | FLOAT) [ ]* ('ct'|'Cr');
//DegreeValue:   (INT+ | FLOAT) [ ]* ('deg'|'rad');
//HealthValue:   (INT+ | FLOAT) [ ]* 'hp';

PostfixTime:     ('ms' | 's' | 'min' | 'h');
PostfixDistance: ('m' | 'km');
PostfixMoney:    ('ct' | 'Cr');
PostfixAngle:    ('deg' | 'rad');
PostfixHealth:   'hp';
PostfixInteger:  'i'|'L';
PostfixFloat:    'f'|'LF';

/* UnaryPostfixModifier: 'ms' | 's' | 'min' | 'h' | 'm' | 'km' | 'ct' | 'Cr' | 'deg' | 'rad' | 'hp'; */

Float  : FLOAT [ ]* ('f'|'LF')?
       | INT+ [ ]* ('f'|'LF');

Integer: INT+ [iL]?;
SString: ['] ('\\'['] | ~['])* ['];

/* BlankLine: Newline Whitespace* Newline; */
/* BlankLine: Newline Newline; */


//{h}{p}			{ ECHON; SETSTR;  return TOK_HP; }
//{m}{a}{x}		{ ECHON; SETSTR;  return TOK_MAX; }
//{m}{i}{n}		{ ECHON; SETSTR;  return TOK_MIN; }
//{m}{s}			{ ECHON; SETSTR;  return TOK_MS; }
//{k}{m}			{ ECHON; SETSTR;  return TOK_KM; }
//{c}{r}			{ ECHON; SETSTR;  return TOK_CR; }
//{d}{e}{g}		{ ECHON; SETSTR;  return TOK_DEG; }
//{c}{t}			{ ECHON; SETSTR;  return TOK_CT; }
//{l}{f}			{ ECHON; SETSTR;  return TOK_LF; }
//{m}			{ ECHON; SETSTR;  return 'm'; }
//{f}			{ ECHON; SETCHAR; return 'f'; }
//{i}			{ ECHON; SETCHAR; return 'i'; }
//{s}			{ ECHON; SETCHAR; return 's'; }
//{h}			{ ECHON; SETCHAR; return 'h'; }
//{l}			{ ECHON; SETCHAR; return 'L'; }

//TOK_ADD:        'add';
//TOK_BREAK:      'break';
TOK_CHANCE:     'chance';
//TOK_CONTINUE:   'continue';
//TOK_DEBUG:      'debug';
//TOK_DEBUGCHANCE:'debugchance';
TOK_ELSE:       'else';
TOK_ELSEIF:     'elseif';
TOK_FOR:        'for';
TOK_FOREACH:    'foreach';
TOK_IF:         'if';
TOK_IN:         'in';
//TOK_INSERT:     'insert';
//TOK_IS:         'is';
//TOK_ISNOT:      'isnot';
//TOK_LET:        'let';
//TOK_LIST:       'list';
TOK_NOT:        'not';
TOK_NULL:       'null';
//TOK_RANGE:      'range';
//TOK_RESUME:     'resume';
//TOK_RETURN:     'return';
//TOK_REVERSED:   'reversed';
//TOK_SQRT:       'sqrt';
//TOK_SUBTRACT:   'subtract';
TOK_TABLE:      'table';
TOK_THEN:       'then';
TOK_TYPEOF:     'typeof';
//TOK_UNDEF:      'undef';
//TOK_UNLET:      'unlet';
//TOK_WEIGHT:     'weight';
TOK_WHILE:      'while';
TOK_MIN: 'min';
TOK_MAX: 'max';

//TOK_CONST: 'event';
//TOK_CONST: 'this';
//TOK_CONST: 'error';

//OP_FILTER: '>>';
OP_EQ: '==';
OP_NEQ: '!=';
OP_LE: '<=';
OP_GE: '>=';
OP_AND: '&&';
OP_OR: '||';
//OP_ARROW: '=>';

TOK_AND: 'and';
TOK_OR: 'or';
TOK_GE: 'ge';
TOK_GT: 'gt';
TOK_LE: 'le';
TOK_LT: 'lt';

//TOK_DBG_CONTEXT: '#[';

ATSIGN:    '@';
BACKSLASH: '\\';
DOLLAR:    '$';
EQUALS:    '=';
EXCLAM:    '!';
QMARK:     '?';
LBRACKET:  '[';
RBRACKET:  ']';
LBRACE:    '{';
RBRACE:    '}';
LPAREN:    '(';
RPAREN:    ')';
LANGLE:    '<';
RANGLE:    '>';
POWER:     '^';
PLUS:      '+';
MINUS:     '-';
ASTERIX:   '*';
SLASH:     '/';
PERCENT:   '%';
COLON:     ':';
COMMA:     ',';
PERIOD:    '.';
SEMICOLON: ';';
SQUOTE:    ['];
DQUOTE:    '"';


/* And my types... */
Variable: '$' IdentChar+;
BareIdentifier: IdentHead IdentChar*;
AttributeValue: '"' (~'"' | '\\"')* '"';

LineComment:  '//' .*? Newline;
BlockComment: '/*' .*? '*/';

Newline:    ('\r\n' | '\n') -> skip;
Whitespace: [\t ]+          -> skip;

