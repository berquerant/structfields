%{
package internal

import  "github.com/berquerant/ybase"
%}

%union{
  result *Tags

  tags []*Tag
  tag *Tag

  token ybase.Token
}

%type <result> result
%type <tags> tags
%type <tag> tag

%token <token> IDENT
%token <token> COLON
%token <token> DQUOTE
%token <token> STRING
%token <token> SPACES

%%

result:
  tags {
    r := &Tags{
      List: $1,
    }
    yylex.(*TagLexer).Result = r
    $$ = r
  }

tags:
  tag {
    $$ = []*Tag{$1}
  }
  | tag SPACES tags {
    $$ = append([]*Tag{$1}, $3...)
  }

tag:
  IDENT COLON DQUOTE STRING DQUOTE {
    $$ = &Tag{
      Key: $1.Value(),
      Value: $4.Value(),
    }
  }
