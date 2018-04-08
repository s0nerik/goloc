package goloc

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFormatsEmptyData(t *testing.T) {
	var data [][]interface{}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "")
	if assert.Error(t, err) {
		assert.IsType(t, &emptySheetError{}, err)
	}
}

func TestFormatsEmptyFirstRow(t *testing.T) {
	data := [][]interface{}{
		{},
		{"x"},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "")
	if assert.Error(t, err) {
		assert.IsType(t, &firstRowNotFoundError{}, err)
	}
}

func TestFormatsMissingFormatColumn(t *testing.T) {
	data := [][]interface{}{
		{"mock"},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &noFormatColumnError{}, err)
	}
}

func TestFormatsMissingPlatformColumn(t *testing.T) {
	data := [][]interface{}{
		{"format"},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &noPlatformColumnError{}, err)
	}
}

func TestFormatsMissingFormatKey(t *testing.T) {
	data := [][]interface{}{
		{"mock", "format"},
		{""},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &formatKeyNotSpecifiedError{}, err)
	}
}

func TestFormatsMissingFormatValue1(t *testing.T) {
	data := [][]interface{}{
		{"format", "mock"},
		{""},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &formatValueNotSpecifiedError{}, err)
	}
}

func TestFormatsMissingFormatValue2(t *testing.T) {
	data := [][]interface{}{
		{"format", "mock"},
		{"x", ""},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &formatValueNotSpecifiedError{}, err)
	}
}

func TestFormatsMissingFormatValue3(t *testing.T) {
	data := [][]interface{}{
		{"format", "mock"},
		{"x", "   "},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &formatValueNotSpecifiedError{}, err)
	}
}

func TestFormatsFormatValidation(t *testing.T) {
	data := [][]interface{}{
		{"format", "mock"},
		{"x", "s"},
		{"y", "%s"},
		{"z", "%s"},
	}

	platform := newMockPlatform(func(p *mockPlatform) {
		p.On("ValidateFormat", "s").Return(nil)
		p.On("ValidateFormat", "%s").Return(&formatValueInvalidError{})
	})

	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &formatValueInvalidError{}, err)
	}

	platform.AssertCalled(t, "ValidateFormat", "s")
	platform.AssertCalled(t, "ValidateFormat", "%s")
	platform.AssertNumberOfCalls(t, "ValidateFormat", 2)
}

func TestFormatsWrongValueType(t *testing.T) {
	data := [][]interface{}{
		{"format", "mock"},
		{"x", "s"},
		{"y", 1},
		{"z", "%s"},
	}

	platform := newMockPlatform(nil)

	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &wrongValueTypeError{}, err)
	}

	platform.AssertCalled(t, "ValidateFormat", "s")
	platform.AssertNotCalled(t, "ValidateFormat", "%s")
	platform.AssertNumberOfCalls(t, "ValidateFormat", 1)
}

func TestFormatsWrongKeyType(t *testing.T) {
	data := [][]interface{}{
		{"format", "mock"},
		{"x", "s"},
		{1, "%s"},
		{"z", "%s"},
	}

	platform := newMockPlatform(nil)

	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &wrongKeyTypeError{}, err)
	}

	platform.AssertCalled(t, "ValidateFormat", "s")
	platform.AssertNotCalled(t, "ValidateFormat", "%s")
	platform.AssertNumberOfCalls(t, "ValidateFormat", 1)
}