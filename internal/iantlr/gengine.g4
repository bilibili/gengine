grammar gengine;

primary: ruleEntity+;

ruleEntity:  RULE ruleName ruleDescription? salience? BEGIN ruleContent END;
ruleName : stringLiteral;
ruleDescription : stringLiteral;
salience : SALIENCE integer;
ruleContent : statements;
statements: statement* returnStmt? ;

statement : ifStmt | functionCall | methodCall  | threeLevelCall | assignment | concStatement  | forStmt | breakStmt |forRangeStmt | continueStmt;

concStatement : CONC LR_BRACE ( functionCall | methodCall | threeLevelCall | assignment )* RR_BRACE;

expression : mathExpression
            | expression comparisonOperator expression
            | expression logicalOperator expression
            | notOperator ? expressionAtom
            | notOperator ? LR_BRACKET expression  RR_BRACKET
            ;

mathExpression : mathExpression  mathMdOperator mathExpression
               | mathExpression  mathPmOperator mathExpression
               | expressionAtom
               | LR_BRACKET mathExpression RR_BRACKET
               ;

expressionAtom
    : functionCall
    | methodCall
    | threeLevelCall
    | constant
    | mapVar
    | variable
    ;

assignment : (mapVar | variable) assignOperator (mathExpression| expression);

returnStmt : RETURN expression?;

ifStmt : IF expression LR_BRACE statements RR_BRACE elseIfStmt*  elseStmt? ;

elseIfStmt : ELSE IF expression LR_BRACE statements RR_BRACE;

elseStmt : ELSE LR_BRACE statements RR_BRACE;

forStmt : FOR assignment SEMICOLON expression SEMICOLON assignment LR_BRACE statements RR_BRACE;

breakStmt: BREAK;

forRangeStmt : FORRANGE variable rangeOperator variable LR_BRACE statements RR_BRACE;

continueStmt:  CONTINUE;

constant
    : booleanLiteral
    | integer
    | realLiteral
    | stringLiteral
    | atName
    | atId
    | atDesc
    | atSal
    ;

functionArgs
    : (constant | variable  | functionCall | methodCall | threeLevelCall | mapVar | expression)  (','(constant | variable | functionCall | methodCall | threeLevelCall | mapVar | expression))*
    ;

integer : MINUS? INT;

realLiteral : MINUS? REAL_LITERAL;

stringLiteral: DQUOTA_STRING ;

booleanLiteral : TRUE | FALSE;

functionCall : SIMPLENAME LR_BRACKET functionArgs? RR_BRACKET;

methodCall : DOTTEDNAME LR_BRACKET functionArgs? RR_BRACKET;

threeLevelCall : DOUBLEDOTTEDNAME LR_BRACKET functionArgs? RR_BRACKET;

variable :  SIMPLENAME | DOTTEDNAME | DOUBLEDOTTEDNAME;

mathPmOperator : PLUS | MINUS ;

mathMdOperator : MUL | DIV ;

comparisonOperator : GT | LT | GTE | LTE | EQUALS | NOTEQUALS ;

logicalOperator : AND | OR ;

assignOperator: ASSIGN | SET | PLUSEQUAL | MINUSEQUAL | MULTIEQUAL | DIVEQUAL ;

rangeOperator: ASSIGN;

notOperator: NOT;

mapVar: variable LSQARE (integer |stringLiteral | variable ) RSQARE;

atName : '@name';
atId : '@id';
atDesc : '@desc';
atSal : '@sal';

fragment DEC_DIGIT          : [0-9];
fragment A                  : [aA] ;
fragment B                  : [bB] ;
fragment C                  : [cC] ;
fragment D                  : [dD] ;
fragment E                  : [eE] ;
fragment F                  : [fF] ;
fragment G                  : [gG] ;
fragment H                  : [hH] ;
fragment I                  : [iI] ;
fragment J                  : [jJ] ;
fragment K                  : [kK] ;
fragment L                  : [lL] ;
fragment M                  : [mM] ;
fragment N                  : [nN] ;
fragment O                  : [oO] ;
fragment P                  : [pP] ;
fragment Q                  : [qQ] ;
fragment R                  : [rR] ;
fragment S                  : [sS] ;
fragment T                  : [tT] ;
fragment U                  : [uU] ;
fragment V                  : [vV] ;
fragment W                  : [wW] ;
fragment X                  : [xX] ;
fragment Y                  : [yY] ;
fragment Z                  : [zZ] ;
fragment EXPONENT_NUM_PART  : ('E'| 'e') '-'? DEC_DIGIT+;
fragment HEX_DIGIT          : [0-9a-fA-F];
fragment OCTAL_DIGIT        : [0-7];

fragment ESCAPED_VALUE
    : '\\' ('u' HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT
           | 'U' HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT HEX_DIGIT
           | [abfnrtv\\'"]
           | OCTAL_DIGIT OCTAL_DIGIT OCTAL_DIGIT
           | 'x' HEX_DIGIT HEX_DIGIT)
    ;

NIL                         : N I L;
RULE                        : R U L E  ;
AND                         : '&&' ;
OR                          : '||' ;

CONC                        : C O N C;
IF                          : I F;
ELSE                        : E L S E;
RETURN                      : R E T U R N;
FOR                         : F O R;
BREAK                       : B R E A K;
FORRANGE                    : F O R R A N G E;
CONTINUE                    : C O N T I N U E;

TRUE                        : T R U E ;
FALSE                       : F A L S E ;
NULL_LITERAL                : N U L L ;
SALIENCE                    : S A L I E N C E ;
BEGIN                       : B E G I N;
END                         : E N D;

SIMPLENAME :  ('a'..'z' |'A'..'Z'| '_')+ ( ('0'..'9') | ('a'..'z' |'A'..'Z') | '_' )* ;

INT : '0'..'9' + ;
PLUS                        : '+' ;
MINUS                       : '-' ;
DIV                         : '/' ;
MUL                         : '*' ;

EQUALS                      : '==' ;
GT                          : '>' ;
LT                          : '<' ;
GTE                         : '>=' ;
LTE                         : '<=' ;
NOTEQUALS                   : '!=' ;
NOT                         : '!' ;

ASSIGN                      : ':=' ;
SET                         : '=';
PLUSEQUAL                   : '+=';
MINUSEQUAL                  : '-=';
MULTIEQUAL                  : '*=';
DIVEQUAL                    : '/=';

LSQARE    : '[' ;
RSQARE    : ']' ;

SEMICOLON                   : ';' ;
LR_BRACE                    : '{';
RR_BRACE                    : '}';
LR_BRACKET                  : '(';
RR_BRACKET                  : ')';
DOT                         : '.' ;
DQUOTA_STRING               : '"' (~["\\\r\n] | ESCAPED_VALUE)*  '"';
DOTTEDNAME                  : SIMPLENAME DOT SIMPLENAME  ;
DOUBLEDOTTEDNAME            : SIMPLENAME DOT SIMPLENAME DOT SIMPLENAME;
REAL_LITERAL                : (DEC_DIGIT+)? '.' DEC_DIGIT+
                            | DEC_DIGIT+ '.' EXPONENT_NUM_PART
                            | (DEC_DIGIT+)? '.' (DEC_DIGIT+ EXPONENT_NUM_PART)
                            | DEC_DIGIT+ EXPONENT_NUM_PART
                            ;

SL_COMMENT: '//' .*? '\n' -> skip ;
WS  :   [ \t\n\r]+ -> skip ;