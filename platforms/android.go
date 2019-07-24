package platforms

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/s0nerik/goloc/goloc"
	"github.com/s0nerik/goloc/registry"
)

func init() {
	registry.RegisterPlatform(&android{})
}

type android struct{}

func (android) Names() []string {
	return []string{
		"android",
		"Android",
	}
}

func (android) LocalizationFilePath(lang goloc.Lang, resDir goloc.ResDir) string {
	fileName := "localized_strings.xml"
	targetDir := fmt.Sprintf("values-%v", lang)
	if resDir != "" {
		return filepath.Join(resDir, targetDir, fileName)
	}
	return filepath.Join("src", "main", "res", targetDir, fileName)
}

func (android) Header(args *goloc.HeaderArgs) string {
	return "<?xml version=\"1.0\" encoding=\"utf-8\" ?>\n<resources>\n"
}

func (android) LocalizedString(args *goloc.LocalizedStringArgs) string {
	return fmt.Sprintf("\t<string name=\"%v\">%v</string>\n", args.Key, args.Value)
}

func (android) Footer(args *goloc.FooterArgs) string {
	return "</resources>\n"
}

func (android) ValidateFormat(format string) error {
	return nil
}

func (android) FormatString(args *goloc.FormatStringArgs) string {
	return fmt.Sprintf(`%%%v$%v`, args.Index+1, strings.TrimPrefix(args.Format, `%`))
}

func (android) ReplacementChars() map[string]string {
	return map[string]string{
		`\`:  `\\`,
		`'`:  `\'`,
		`"`:  `\"`,
		"\n": `\n`,
		`?`:  `\?`,
		`@`:  `\@`,
		`<`:  `&lt;`,
		`>`:  `&gt;`,
		`&`:  `&amp;`,
	}
}
