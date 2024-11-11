package internal_test

import (
	"bytes"
	"testing"

	"github.com/berquerant/structfields/internal"
	"github.com/stretchr/testify/assert"
)

func TestTag(t *testing.T) {
	// internal.SetDebug(3)
	// slog.SetLogLoggerLevel(slog.LevelDebug)

	for _, tc := range []struct {
		title string
		input string
	}{
		{
			title: "not key-value",
			input: "value",
		},
		{
			title: "empty",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			r := bytes.NewBufferString(tc.input)
			lex := internal.NewTagLexer(r)
			_ = internal.ParseTags(lex)
			assert.NotNil(t, lex.Err())
		})
	}

	for _, tc := range []struct {
		title string
		input string
		want  *internal.Tags
	}{
		{
			title: "1 tag",
			input: `KEY:"VALUE"`,
			want: &internal.Tags{
				List: []*internal.Tag{
					&internal.Tag{
						Key:   "KEY",
						Value: "VALUE",
					},
				},
			},
		},
		{
			title: "2 tags",
			input: `KEY:"VALUE" KEY2:"VALUE2"`,
			want: &internal.Tags{
				List: []*internal.Tag{
					&internal.Tag{
						Key:   "KEY",
						Value: "VALUE",
					},
					&internal.Tag{
						Key:   "KEY2",
						Value: "VALUE2",
					},
				},
			},
		},
		{
			title: "1 tag quoted",
			input: `KEY:"VALUE\"2\""`,
			want: &internal.Tags{
				List: []*internal.Tag{
					&internal.Tag{
						Key:   "KEY",
						Value: `VALUE\"2\"`,
					},
				},
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			r := bytes.NewBufferString(tc.input)
			lex := internal.NewTagLexer(r)
			_ = internal.ParseTags(lex)
			if !assert.Nil(t, lex.Err()) {
				t.Errorf(lex.Err().Error())
				return
			}
			assert.Equal(t, tc.want, lex.Result)
		})
	}
}
