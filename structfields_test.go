package structfields_test

import (
	"testing"

	"github.com/berquerant/structfields"
	"github.com/stretchr/testify/assert"
)

func TestStruct(t *testing.T) {
	// internal.SetDebug(3)
	// slog.SetLogLoggerLevel(slog.LevelDebug)

	type Empty struct{}
	type EmptyTag struct {
		A int ``
	}
	type NonKeyValue struct {
		A int `string` // TODO: go vet ignore this
	}
	type KeyValue struct {
		A int `key:"value"`
	}
	type KeyValues struct {
		A int `key:"value"`
		B int `key:"value" key2:"value2"`
	}

	for _, tc := range []struct {
		title string
		v     any
		want  *structfields.Struct
		isErr bool
	}{
		{
			title: "key-values",
			v:     KeyValues{},
			want: &structfields.Struct{
				Name: "KeyValues",
				Fields: []structfields.Field{
					{
						Name:    "A",
						Tag:     `key:"value"`,
						TagKeys: []string{"key"},
					},
					{
						Name:    "B",
						Tag:     `key:"value" key2:"value2"`,
						TagKeys: []string{"key", "key2"},
					},
				},
			},
		},
		{
			title: "a key-value",
			v:     KeyValue{},
			want: &structfields.Struct{
				Name: "KeyValue",
				Fields: []structfields.Field{
					{
						Name:    "A",
						Tag:     `key:"value"`,
						TagKeys: []string{"key"},
					},
				},
			},
		},
		{
			title: "not key-value",
			v:     NonKeyValue{},
			isErr: true,
		},
		{
			title: "empty",
			v:     Empty{},
			want: &structfields.Struct{
				Name:   "Empty",
				Fields: []structfields.Field{},
			},
		},
		{
			title: "emptytag",
			v:     EmptyTag{},
			want: &structfields.Struct{
				Name: "EmptyTag",
				Fields: []structfields.Field{
					{
						Name: "A",
						Tag:  "",
					},
				},
			},
		},
		{
			title: "not a struct",
			v:     10,
			isErr: true,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got, err := structfields.New(tc.v)
			if tc.isErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.want.Name, got.Name)
			assert.Equal(t, len(tc.want.Fields), len(got.Fields))
			for i, w := range tc.want.Fields {
				g := got.Fields[i]
				assert.Equal(t, w.Name, g.Name)
				assert.Equal(t, w.Tag, g.Tag)
				assert.Equal(t, w.TagKeys, g.TagKeys)
			}
		})
	}
}
