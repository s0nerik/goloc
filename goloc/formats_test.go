package goloc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatsEmptyData(t *testing.T) {
	var data [][]string

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "", "{}")
	if assert.Error(t, err) {
		assert.IsType(t, &emptySheetError{}, err)
	}
}

func TestFormatsEmptyFirstRow(t *testing.T) {
	data := [][]RawCell{
		{},
		{"x"},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "", "{}")
	if assert.Error(t, err) {
		assert.IsType(t, &firstRowNotFoundError{}, err)
	}
}

func TestFormatsMissingFormatColumn(t *testing.T) {
	data := [][]RawCell{
		{"mock"},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format", "{}")
	if assert.Error(t, err) {
		assert.IsType(t, &noFormatColumnError{}, err)
	}
}

func TestFormatsMissingPlatformColumn(t *testing.T) {
	data := [][]RawCell{
		{"format"},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format", "{}")
	if assert.Error(t, err) {
		assert.IsType(t, &noPlatformColumnError{}, err)
	}
}

func TestFormatsMissingFormatKey(t *testing.T) {
	data := [][]RawCell{
		{"mock", "format"},
		{""},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format", "{}")
	if assert.Error(t, err) {
		assert.IsType(t, &formatKeyNotSpecifiedError{}, err)
	}
}

func TestFormatsMissingFormatValue1(t *testing.T) {
	data := [][]RawCell{
		{"format", "mock"},
		{""},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format", "{}")
	if assert.Error(t, err) {
		assert.IsType(t, &formatValueNotSpecifiedError{}, err)
	}
}

func TestFormatsMissingFormatValue2(t *testing.T) {
	data := [][]RawCell{
		{"format", "mock"},
		{"x", ""},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format", "{}")
	if assert.Error(t, err) {
		assert.IsType(t, &formatValueNotSpecifiedError{}, err)
	}
}

func TestFormatsMissingFormatValue3(t *testing.T) {
	data := [][]RawCell{
		{"format", "mock"},
		{"x", "   "},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format", "{}")
	if assert.Error(t, err) {
		assert.IsType(t, &formatValueNotSpecifiedError{}, err)
	}
}

func TestFormatsFormatValidation(t *testing.T) {
	data := [][]RawCell{
		{"format", "mock"},
		{"x", "s"},
		{"y", "%s"},
		{"z", "%s"},
	}

	platform := newMockPlatform(func(p *mockPlatform) {
		p.On("ValidateFormat", "s").Return(nil)
		p.On("ValidateFormat", "%s").Return(&formatValueInvalidError{})
	})

	_, err := ParseFormats(data, platform, "", "format", "{}")
	if assert.Error(t, err) {
		assert.IsType(t, &formatValueInvalidError{}, err)
	}

	platform.AssertCalled(t, "ValidateFormat", "s")
	platform.AssertCalled(t, "ValidateFormat", "%s")
	platform.AssertNumberOfCalls(t, "ValidateFormat", 2)
}