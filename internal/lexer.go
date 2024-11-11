package internal

import (
	"io"
	"log/slog"
	"unicode"

	"github.com/berquerant/ybase"
)

func SetDebug(level int) {
	yyDebug = level
}

func ParseTags(lex *TagLexer) int {
	return yyParse(lex)
}

func NewTagLexer(r io.Reader) *TagLexer {
	yyErrorVerbose = true
	lex := &TagLexer{}
	lex.Lexer = ybase.NewLexer(
		ybase.NewScanner(
			ybase.NewReader(r, slog.Debug),
			lex.ScanFunc,
		),
	)
	return lex
}

var (
	_ yyLexer = &TagLexer{}
)

type TagLexer struct {
	ybase.Lexer
	Result *Tags

	expectString      bool
	expectStringClose bool
}

func (lex *TagLexer) Lex(lval *yySymType) int {
	return lex.DoLex(func(tok ybase.Token) {
		lval.token = tok
	})
}

func (lex *TagLexer) eof() int {
	return ybase.EOF
}

func (lex *TagLexer) ScanFunc(r ybase.Reader) int {
	slog.Debug("TagLexer",
		slog.Bool("expectString", lex.expectString),
		slog.Bool("expectStringClose", lex.expectStringClose),
	)

	switch {
	case r.Peek() == ybase.EOF:
		return lex.eof()
	case lex.expectStringClose:
		lex.expectStringClose = false
		if r.Peek() == '"' {
			_ = r.Next()
			return DQUOTE
		}
		lex.Error("TagLexer: expect string close failure")
		return lex.eof()
	case lex.expectString:
		lex.expectString = false
		if lex.scanString(r) {
			lex.expectStringClose = true
			return STRING
		}
		lex.Error("TagLexer: expect string failure")
		return lex.eof()
	case unicode.IsSpace(r.Peek()):
		r.NextWhile(unicode.IsSpace)
		return SPACES
	}

	switch r.Peek() {
	case ':':
		_ = r.Next()
		return COLON
	case '"':
		_ = r.Next()
		lex.expectString = true
		return DQUOTE
	default:
		if lex.scanIdent(r) {
			return IDENT
		}
		return lex.eof()
	}
}

func (lex *TagLexer) scanIdent(r ybase.Reader) bool {
	slog.Debug("TagLexer: scanIdent")
	defer func() {
		slog.Debug("TagLexer: end scanIdent")
	}()
	r.NextWhile(func(x rune) bool {
		return x != ':' && x != ybase.EOF
	})
	return true
}

func (lex *TagLexer) scanString(r ybase.Reader) bool {
	slog.Debug("TagLexer: scanString")
	defer func() {
		slog.Debug("TagLexer: end scanString")
	}()
	escape := false
	for {
		x := r.Peek()
		switch {
		case x == ybase.EOF:
			return false
		case x == '\\':
			escape = true
			_ = r.Next()
		case x == '"':
			if escape {
				escape = false
				_ = r.Next()
				continue
			}
			return true
		default:
			escape = false
			_ = r.Next()
		}
	}
}
