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
	var emptyData [][]interface{}

	platform := &mockPlatform{}
	_, err := ParseFormats(emptyData, platform, "", "")
	if assert.Error(t, err) {
		assert.IsType(t, &emptySheetError{}, err)
	}
}

func TestEmptyFirstRow(t *testing.T) {
	emptyFirstRowData := [][]interface{}{
		{},
		{"x"},
	}

	platform := &mockPlatform{}
	_, err := ParseFormats(emptyFirstRowData, platform, "", "")
	if assert.Error(t, err) {
		assert.IsType(t, &emptyFirstRowError{}, err)
	}
}

func TestMissingFormatColumn(t *testing.T) {
	emptyFirstRowData := [][]interface{}{
		{"mock"},
	}

	platform := &mockPlatform{}
	_, err := ParseFormats(emptyFirstRowData, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &noFormatColumnError{}, err)
	}
}

func TestMissingPlatformColumn(t *testing.T) {
	emptyFirstRowData := [][]interface{}{
		{"format"},
	}

	platform := &mockPlatform{}
	_, err := ParseFormats(emptyFirstRowData, platform, "", "format")
	if assert.Error(t, err) {
		assert.IsType(t, &noPlatformColumnError{}, err)
	}
}
