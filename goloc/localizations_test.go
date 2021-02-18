package goloc

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func formats() Formats {
	data := [][]RawCell{
		{"format",				"mock"},
		{"x",					"s"},
		{"y",					"%s"},
		{"z",					"%s"},
	}

	platform := newMockPlatform(nil)
	formats, _ := ParseFormats(data, platform, "", "format", "{}")
	return formats
}

func parseTestLocalizations(data [][]RawCell, errorIfMissing bool, missingRegexp *regexp.Regexp) (loc Localizations, warnings []error, error error) {
	p := newMockPlatform(nil)
	loc, _, warnings, error = ParseLocalizations(data, p, formats(), "", "key", errorIfMissing, missingRegexp)
	return
}

func TestLocalizationsEmptyData(t *testing.T) {
	var data [][]RawCell

	_, _, err := parseTestLocalizations(data, true, nil)
	assert.Error(t, err)
	assert.IsType(t, &emptySheetError{}, err)
}

func TestLocalizationsEmptyFirstRow(t *testing.T) {
	data := [][]RawCell{
		{},
		{"x"},
	}

	_, _, err := parseTestLocalizations(data, true, nil)
	assert.IsType(t, &firstRowNotFoundError{}, err)
}

func TestLocalizationsNoKeyColumn(t *testing.T) {
	data := [][]RawCell{
		{"", "lang_en"},
	}

	_, _, err := parseTestLocalizations(data, true, nil)
	assert.Error(t, err)
	assert.IsType(t, &columnNotFoundError{}, err)
}

func TestLocalizationsNoLangColumns(t *testing.T) {
	data := [][]RawCell{
		{"key", "something", "something else"},
	}

	_, _, err := parseTestLocalizations(data, true, nil)
	assert.Error(t, err)
	assert.IsType(t, &langColumnsNotFoundError{}, err)
}

func TestLocalizationsMissingKey(t *testing.T) {
	dataBad := [][][]RawCell{
		{
			{"key",				"lang_en"},
			{"",				"something"},
		},
		{
			{"key",				"lang_en"},
			{" ",				"something"},
		},
		{
			{"lang_en",			"key"},
			{"something"				},
		},
		{
			{"lang_en",			"key"},
			{"something", 		""},
		},
		{
			{"lang_en",			"key"},
			{"something",		" "},
		},
	}

	dataGood := [][][]RawCell{
		{
			{"key",				"lang_en"},
			{"k",				"something"},
		},
		{
			{"key",				"lang_en"},
			{"k",				"something {x}"},
		},
		{
			{"lang_en",			"key"},
			{"something",		"k"},
		},
		{
			{"lang_en",			"key"},
			{"something {x}",	"k"},
		},
	}

	for _, d := range dataBad {
		_, _, err := parseTestLocalizations(d, true, nil)
		assert.Error(t, err)
		assert.IsType(t, &keyMissingError{}, err)

		_, warn, err := parseTestLocalizations(d, false, nil)
		assert.Nil(t, err)
		assert.Len(t, warn, 1)
		assert.IsType(t, &keyMissingError{}, warn[0])
	}

	for _, d := range dataGood {
		_, warn, err := parseTestLocalizations(d, true, nil)
		assert.Nil(t, err)
		assert.Empty(t, warn)
	}
}

func TestLocalizationsMissingFormat(t *testing.T) {
	dataBad := [][][]RawCell{
		{
			{"key",				"lang_en"},
			{"p1",				"something {x}"},
			{"p2",				"something {y}"},
			{"m",				"something {missing}"},
			{"p3",				"something {z}"},
		},
	}

	dataGood := [][][]RawCell{
		{
			{"key",				"lang_en"},
			{"p1",				"something {x}"},
			{"p2",				"something {y}"},
			{"p3",				"something {z}"},
		},
	}

	for _, d := range dataBad {
		_, _, err := parseTestLocalizations(d, true, nil)
		assert.Error(t, err)
		assert.IsType(t, &formatNotFoundError{}, err)
	}

	for _, d := range dataGood {
		_, warn, err := parseTestLocalizations(d, true, nil)
		assert.Nil(t, err)
		assert.Empty(t, warn)
	}
}

func TestLocalizationsMissingLocalization(t *testing.T) {
	dataBad := [][][]RawCell{
		{
			{"key",			"lang_en",			"lang_ru"},
			{"m",			"something {y}"					},
		},
		{
			{"key",			"lang_en",			"lang_ru"},
			{"m",			"something {y}",	""},
		},
		{
			{"key",			"lang_en",			"lang_ru"},
			{"m",			"something {y}",	" "},
		},
		{
			{"key",			"lang_en",			"lang_ru"},
			{"m",			"",					"что-то {y}"},
		},
		{
			{"key",			"lang_en",			"lang_ru"},
			{"m",			" ",				"что-то {y}"},
		},
	}

	dataGood := [][][]RawCell{
		{
			{"key",			"lang_en",			"lang_ru"},
			{"p",			"something {x}",	"что-то {x}"},
		},
	}

	for _, d := range dataBad {
		_, _, err := parseTestLocalizations(d, true, nil)
		assert.Error(t, err)
		assert.IsType(t, &localizationMissingError{}, err)

		_, warn, err := parseTestLocalizations(d, false, nil)
		assert.Nil(t, err)
		assert.Len(t, warn, 1)
		assert.IsType(t, &localizationMissingError{}, warn[0])
	}

	for _, d := range dataGood {
		_, warn, err := parseTestLocalizations(d, true, nil)
		assert.Nil(t, err)
		assert.Empty(t, warn)
	}
}

func TestLocalizationsMissingRegexpLocalization(t *testing.T) {
	re, _ := regexp.Compile("^x$")

	dataBad := [][][]RawCell{
		{
			{"key",			"lang_en"},
			{"m",			"x"},
		},
		{
			{"key",			"lang_en"},
			{"m",			"  x  "},
		},
	}

	dataGood := [][][]RawCell{
		{
			{"key",			"lang_en"},
			{"m",			"xy"},
		},
		{
			{"key",			"lang_en"},
			{"m",			" xyz   "},
		},
	}

	for _, d := range dataBad {
		_, _, err := parseTestLocalizations(d, true, re)
		assert.Error(t, err)
		assert.IsType(t, &localizationMissingError{}, err)

		_, warn, err := parseTestLocalizations(d, false, re)
		assert.Nil(t, err)
		assert.Len(t, warn, 1)
		assert.IsType(t, &localizationMissingError{}, warn[0])
	}

	for _, d := range dataGood {
		_, warn, err := parseTestLocalizations(d, true, re)
		assert.Nil(t, err)
		assert.Empty(t, warn)
	}
}
