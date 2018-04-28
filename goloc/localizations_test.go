package goloc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func formats() Formats {
	data := [][]interface{}{
		{"format", "mock"},
		{"x", "s"},
		{"y", "%s"},
		{"z", "%s"},
	}

	platform := newMockPlatform(nil)
	formats, _ := ParseFormats(data, platform, "", "format")
	return formats
}

func TestLocalizationsEmptyData(t *testing.T) {
	var data [][]interface{}

	p := newMockPlatform(nil)
	_, _, err := ParseLocalizations(data, p, formats(), "", "", true)
	assert.Error(t, err)
	assert.IsType(t, &emptySheetError{}, err)
}

func TestLocalizationsEmptyFirstRow(t *testing.T) {
	data := [][]interface{}{
		{},
		{"x"},
	}

	p := newMockPlatform(nil)
	_, _, err := ParseLocalizations(data, p, formats(), "", "", true)
	assert.IsType(t, &firstRowNotFoundError{}, err)
}

func TestLocalizationsNoKeyColumn(t *testing.T) {
	data := [][]interface{}{
		{"", "lang_en"},
	}

	p := newMockPlatform(nil)
	_, _, err := ParseLocalizations(data, p, formats(), "", "key", true)
	assert.Error(t, err)
	assert.IsType(t, &columnNotFoundError{}, err)
}

func TestLocalizationsNoLangColumns(t *testing.T) {
	data := [][]interface{}{
		{"key", "something", "something else"},
	}

	p := newMockPlatform(nil)
	_, _, err := ParseLocalizations(data, p, formats(), "", "key", true)
	assert.Error(t, err)
	assert.IsType(t, &langColumnsNotFoundError{}, err)
}

func TestLocalizationsMissingKey(t *testing.T) {
	dataBad := [][][]interface{}{
		{
			{"key", "lang_en"},
			{"", "something"},
		},
		{
			{"key", "lang_en"},
			{" ", "something"},
		},
		{
			{"lang_en", "key"},
			{"something"},
		},
		{
			{"lang_en", "key"},
			{"something", ""},
		},
		{
			{"lang_en", "key"},
			{"something", " "},
		},
	}

	dataGood := [][][]interface{}{
		{
			{"key", "lang_en"},
			{"k", "something"},
		},
		{
			{"key", "lang_en"},
			{"k", "something {x}"},
		},
		{
			{"lang_en", "key"},
			{"something", "k"},
		},
		{
			{"lang_en", "key"},
			{"something {x}", "k"},
		},
	}

	for _, d := range dataBad {
		p := newMockPlatform(nil)

		_, _, err := ParseLocalizations(d, p, formats(), "", "key", true)
		assert.Error(t, err)
		assert.IsType(t, &keyMissingError{}, err)

		_, warn, err := ParseLocalizations(d, p, formats(), "", "key", false)
		assert.Nil(t, err)
		assert.Len(t, warn, 1)
		assert.IsType(t, &keyMissingError{}, warn[0])
	}

	for _, d := range dataGood {
		p := newMockPlatform(nil)
		_, warn, err := ParseLocalizations(d, p, formats(), "", "key", true)
		assert.Nil(t, err)
		assert.Empty(t, warn)
	}
}

func TestLocalizationsMissingFormat(t *testing.T) {
	dataBad := [][][]interface{}{
		{
			{"key", "lang_en"},
			{"p1", "something {x}"},
			{"p2", "something {y}"},
			{"m", "something {missing}"},
			{"p3", "something {z}"},
		},
	}

	dataGood := [][][]interface{}{
		{
			{"key", "lang_en"},
			{"p1", "something {x}"},
			{"p2", "something {y}"},
			{"p3", "something {z}"},
		},
	}

	for _, d := range dataBad {
		p := newMockPlatform(nil)

		_, _, err := ParseLocalizations(d, p, formats(), "", "key", true)
		assert.Error(t, err)
		assert.IsType(t, &formatNotFoundError{}, err)
	}

	for _, d := range dataGood {
		p := newMockPlatform(nil)
		_, warn, err := ParseLocalizations(d, p, formats(), "", "key", true)
		assert.Nil(t, err)
		assert.Empty(t, warn)
	}
}

func TestLocalizationsMissingLocalization(t *testing.T) {
	dataBad := [][][]interface{}{
		{
			{"key", "lang_en", "lang_ru"},
			{"m", "something {y}"},
		},
		{
			{"key", "lang_en", "lang_ru"},
			{"m", "something {y}", ""},
		},
		{
			{"key", "lang_en", "lang_ru"},
			{"m", "something {y}", " "},
		},
		{
			{"key", "lang_en", "lang_ru"},
			{"m", "", "что-то {y}"},
		},
		{
			{"key", "lang_en", "lang_ru"},
			{"m", " ", "что-то {y}"},
		},
	}

	dataGood := [][][]interface{}{
		{
			{"key", "lang_en", "lang_ru"},
			{"p", "something {x}", "что-то {x}"},
		},
	}

	for _, d := range dataBad {
		p := newMockPlatform(nil)

		_, _, err := ParseLocalizations(d, p, formats(), "", "key", true)
		assert.Error(t, err)
		assert.IsType(t, &localizationMissingError{}, err)

		_, warn, err := ParseLocalizations(d, p, formats(), "", "key", false)
		assert.Nil(t, err)
		assert.Len(t, warn, 1)
		assert.IsType(t, &localizationMissingError{}, warn[0])
	}

	for _, d := range dataGood {
		p := newMockPlatform(nil)
		_, warn, err := ParseLocalizations(d, p, formats(), "", "key", true)
		assert.Nil(t, err)
		assert.Empty(t, warn)
	}
}
