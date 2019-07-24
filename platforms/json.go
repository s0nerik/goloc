package platforms

import (
	"fmt"
	"path/filepath"

	"github.com/s0nerik/goloc/goloc"
	"github.com/s0nerik/goloc/registry"
)

func init() {
	registry.RegisterPlatform(&json{})
}

type json struct{}

func (json) Names() []string {
	return []string{
		"json",
		"JSON",
	}
}

func (json) LocalizationFilePath(lang goloc.Lang, resDir goloc.ResDir) string {
	return filepath.Join(resDir, fmt.Sprintf("%v.json", lang))
}

func (json) Header(args *goloc.HeaderArgs) string {
	return "{\n"
}

func (json) LocalizedString(args *goloc.LocalizedStringArgs) string {
	if args.IsLast {
		return fmt.Sprintf("\t\"%v\": \"%v\"\n", args.Key, args.Value)
	}
	return fmt.Sprintf("\t\"%v\": \"%v\",\n", args.Key, args.Value)
}

func (json) Footer(args *goloc.FooterArgs) string {
	return "}"
}

func (json) ValidateFormat(format string) error {
	return nil
}

func (json) FormatString(args *goloc.FormatStringArgs) string {
	return args.Format
}

func (json) ReplacementChars() map[string]string {
	return map[string]string{
		"\n": `\n`,
		"\t": `\t`,
		`"`:  `\"`,
		`\`:  `\\`,
	}
}
