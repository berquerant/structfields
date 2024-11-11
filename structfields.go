package structfields

import (
	"bytes"
	"errors"
	"log/slog"
	"reflect"

	"github.com/berquerant/structfields/internal"
)

type Field struct {
	// Name is the field name.
	Name    string
	Type    reflect.Type      // field type
	Tag     reflect.StructTag // field tag
	TagKeys []string          // field tag keys
}

// Struct is the metadata of struct.
type Struct struct {
	// Name is the struct name.
	Name   string
	Fields []Field // struct fields info
}

var (
	ErrNotStruct = errors.New("NotStruct")
	ErrParseTag  = errors.New("ParseTag")
)

// New returns a new Struct.
//
// If v is not a struct, return ErrNotStruct.
// If a tag of v is not space-separated key:"value" pairs, return ErrParseTag.
// If v has no tags or its values are empty, ignore them.
func New(v any) (*Struct, error) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	fields := []Field{}
	for i := range t.NumField() {
		f := t.Field(i)
		x := Field{
			Name: f.Name,
			Type: f.Type,
			Tag:  f.Tag,
		}

		if string(f.Tag) == "" {
			// no keys
			fields = append(fields, x)
			continue
		}

		slog.Debug("ParseTag", slog.String("tag", string(f.Tag)))
		r := bytes.NewBufferString(string(f.Tag))
		lex := internal.NewTagLexer(r)
		_ = internal.ParseTags(lex)
		if err := lex.Err(); err != nil {
			// not key-values
			slog.Debug("End ParseTag", slog.String("tag", string(f.Tag)), slog.String("err", err.Error()))
			return nil, errors.Join(ErrParseTag, err)
		}

		keys := make([]string, len(lex.Result.List))
		for i, k := range lex.Result.List {
			slog.Debug("End ParseTag", slog.String("tag", string(f.Tag)), slog.Int("index", i), slog.String("key", k.Key))
			keys[i] = k.Key
		}
		x.TagKeys = keys
		fields = append(fields, x)
	}

	return &Struct{
		Name:   t.Name(),
		Fields: fields,
	}, nil
}
