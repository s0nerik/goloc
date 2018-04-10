package goloc

import (
	"github.com/stretchr/testify/mock"
)

type mockPlatform struct {
	mock.Mock
}

func (p *mockPlatform) Names() []string {
	args := p.Called()
	return args.Get(0).([]string)
}

func (p *mockPlatform) LocalizationFilePath(lang Lang, resDir ResDir) string {
	args := p.Called(lang, resDir)
	return args.String(0)
}

func (p *mockPlatform) Header(lang Lang) string {
	args := p.Called(lang)
	return args.String(0)
}

func (p *mockPlatform) Localization(lang Lang, key Key, value string) string {
	args := p.Called(lang, key, value)
	return args.String(0)
}

func (p *mockPlatform) Footer(lang Lang) string {
	args := p.Called(lang)
	return args.String(0)
}

func (p *mockPlatform) ValidateFormat(format string) error {
	args := p.Called(format)
	return args.Error(0)
}

func (p *mockPlatform) IndexedFormatString(index uint, format string) string {
	args := p.Called(index, format)
	return args.String(0)
}

func (p *mockPlatform) ReplacementChars() map[string]string {
	args := p.Called()
	return args.Get(0).(map[string]string)
}

func newMockPlatform(customMocksProvider func(p *mockPlatform)) *mockPlatform {
	p := &mockPlatform{}
	if customMocksProvider != nil {
		customMocksProvider(p)
	}
	p.On("Names").Return([]string{"mock"})
	p.On("LocalizationFilePath", mock.AnythingOfType("Lang"), mock.AnythingOfType("ResDir")).Return("")
	p.On("Header", mock.AnythingOfType("Lang")).Return("")
	p.On("Localization", mock.AnythingOfType("Lang"), mock.AnythingOfType("Key"), mock.AnythingOfType("string")).Return("")
	p.On("Footer", mock.AnythingOfType("Lang")).Return("")
	p.On("ValidateFormat", mock.AnythingOfType("string")).Return(nil)
	p.On("IndexedFormatString", mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return("")
	p.On("ReplacementChars").Return(map[string]string{
		`~`: `tilde`,
	})
	return p
}