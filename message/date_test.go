package message

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestNewDatetime(t *testing.T) {
	dt := NewDatetime("foo", ShortDatetimeFormat, language.English)
	assert.NotNil(t, dt)

	got, err := dt.Eval(Context{"foo": time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)})
	require.NoError(t, err)
	assert.Equal(t, "4/12/61, 6:07 AM", got)
}

func TestNewTime(t *testing.T) {
	dt := NewTime("foo", ShortDatetimeFormat, language.English)
	assert.NotNil(t, dt)

	got, err := dt.Eval(Context{"foo": time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)})
	require.NoError(t, err)
	assert.Equal(t, "6:07 AM", got)
}

func TestNewDate(t *testing.T) {
	dt := NewDate("foo", ShortDatetimeFormat, language.English)
	assert.NotNil(t, dt)

	got, err := dt.Eval(Context{"foo": time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)})
	require.NoError(t, err)
	assert.Equal(t, "4/12/61", got)
}

func TestDatetime_Eval(t *testing.T) {
	type fields struct {
		argName string
		lang    language.Tag
		format  DatetimeFormat
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     Context
		want    string
		wantErr bool
	}{
		{
			"error if no argument",
			fields{
				argName: "foo",
				lang:    language.Tag{},
				format:  0,
			},
			Context{},
			"",
			true,
		},
		{
			"error if argument is not time",
			fields{
				argName: "foo",
				lang:    language.Tag{},
				format:  0,
			},
			Context{"foo": 42},
			"",
			true,
		},
		{
			"datetime none",
			fields{
				argName: "foo",
				lang:    language.English,
				format:  NoneDatetimeFormat,
			},
			Context{"foo": time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)},
			"",
			false,
		},
		{
			"datetime short",
			fields{
				argName: "foo",
				lang:    language.English,
				format:  ShortDatetimeFormat,
			},
			Context{"foo": time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)},
			"4/12/61, 6:07 AM",
			false,
		},
		{
			"datetime medium",
			fields{
				argName: "foo",
				lang:    language.English,
				format:  MediumDatetimeFormat,
			},
			Context{"foo": time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)},
			"Apr 12, 1961, 6:07:03 AM",
			false,
		},
		{
			"datetime long",
			fields{
				argName: "foo",
				lang:    language.English,
				format:  LongDatetimeFormat,
			},
			Context{"foo": time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)},
			"April 12, 1961 at 6:07:03 AM UTC",
			false,
		},
		{
			"datetime full",
			fields{
				argName: "foo",
				lang:    language.English,
				format:  FullDatetimeFormat,
			},
			Context{"foo": time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)},
			"Wednesday, April 12, 1961 at 6:07:03 AM UTC",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := NewDatetime(
				tt.fields.argName,
				tt.fields.format,
				tt.fields.lang,
			)

			got, err := dt.Eval(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Datetime.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Datetime.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
