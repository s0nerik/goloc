package goloc

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"fmt"
)

type mockPlatform struct {
	mock.Mock
}

func (p mockPlatform) Names() []string {
	return []string{"mock"}
}

func (p mockPlatform) LocalizationFilePath(lang Lang, resDir ResDir) string {
	return ""
}

func (p mockPlatform) Header(lang Lang) string {
	return ""
}

func (p mockPlatform) Localization(lang Lang, key Key, value string) string {
	return fmt.Sprintf("%v = %v\n", key, value)
}

func (p mockPlatform) Footer(lang Lang) string {
	return ""
}

func (p mockPlatform) ValidateFormat(format string) error {
	return nil
}

func (p mockPlatform) IndexedFormatString(index uint, format string) string {
	return format
}

func (p mockPlatform) ReplacementChars() map[string]string {
	return nil
}

func TestEmptyData(t *testing.T) {
	var data [][]interface{}

	platform := &mockPlatform{}
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

	platform := &mockPlatform{}
	_, err := ParseFormats(data, platform, "", "")
	if assert.Error(t, err) {
		assert.IsType(t, &emptyFirstRowError{}, err)
	}
}

func TestMissingFormatColumn(t *testing.T) {
	data := [][]interface{}{
		{"mock"},
	}

	platform := &mockPlatform{}
	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &noFormatColumnError{}, err)
	}
}

func TestMissingPlatformColumn(t *testing.T) {
	data := [][]interface{}{
		{"format"},
	}

	platform := &mockPlatform{}
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

	platform := &mockPlatform{}
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

	platform := &mockPlatform{}
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

	platform := &mockPlatform{}
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

	platform := &mockPlatform{}
	_, err := ParseFormats(data, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &formatValueNotSpecifiedError{}, err)
	}
}