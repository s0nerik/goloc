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

func (p *mockPlatform) LocalizationFilePath(lang Locale, resDir ResDir) string {
	args := p.Called(lang, resDir)
	return args.String(0)
}

func (p *mockPlatform) Header(args *HeaderArgs) string {
	ret := p.Called(args)
	return ret.String(0)
}

func (p *mockPlatform) LocalizedString(args *LocalizedStringArgs) string {
	ret := p.Called(args)
	return ret.String(0)
}

func (p *mockPlatform) Footer(args *FooterArgs) string {
	ret := p.Called(args)
	return ret.String(0)
}

func (p *mockPlatform) ValidateFormat(format string) error {
	args := p.Called(format)
	return args.Error(0)
}

func (p *mockPlatform) FormatString(args *FormatStringArgs) string {
	ret := p.Called(args)
	return ret.String(0)
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
	p.On("Header", mock.AnythingOfType("*goloc.HeaderArgs")).Return("")
	p.On("LocalizedString", mock.AnythingOfType("*goloc.LocalizedStringArgs")).Return("")
	p.On("Footer", mock.AnythingOfType("*goloc.FooterArgs")).Return("")
	p.On("ValidateFormat", mock.AnythingOfType("string")).Return(nil)
	p.On("FormatString", mock.AnythingOfType("*goloc.FormatStringArgs")).Return("")
	p.On("ReplacementChars").Return(map[string]string{
		`~`: `tilde`,
	})
	return p
}
