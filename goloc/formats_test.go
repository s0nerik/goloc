package goloc

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestEmptyData(t *testing.T) {
	var data [][]interface{}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "")
	if assert.Error(t, err) {
		assert.IsType(t, &emptySheetError{}, err)
	}
}

func TestEmptyFirstRow(t *testing.T) {
	data := [][]interface{}{
		{},
		{"x"},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "")
	if assert.Error(t, err) {
		assert.IsType(t, &emptyFirstRowError{}, err)
	}
}

func TestMissingFormatColumn(t *testing.T) {
	data := [][]interface{}{
		{"mock"},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &noFormatColumnError{}, err)
	}
}

func TestMissingPlatformColumn(t *testing.T) {
	data := [][]interface{}{
		{"format"},
	}

	platform := newMockPlatform(nil)
	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &noPlatformColumnError{}, err)
	}
}

func TestMissingFormatKey(t *testing.T) {
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

func TestMissingFormatValue1(t *testing.T) {
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

func TestMissingFormatValue2(t *testing.T) {
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

func TestMissingFormatValue3(t *testing.T) {
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

func TestFormatValidation(t *testing.T) {
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

func TestWrongValueType(t *testing.T) {
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

func TestWrongKeyType(t *testing.T) {
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