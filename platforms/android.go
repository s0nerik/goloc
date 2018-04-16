package platforms

import (
	"github.com/s0nerik/goloc/goloc"
	"fmt"
	"strings"
	"errors"
	"path/filepath"
	"github.com/s0nerik/goloc/platforms/registry"
)

func init() {
	registry.Register(&android{})
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

func (android) Header(lang goloc.Lang) string {
	return "<?xml version=\"1.0\" encoding=\"utf-8\" ?>\n<resources>\n"
}

func (android) Localization(lang goloc.Lang, key goloc.Key, value string) string {
	return fmt.Sprintf("\t<string name=\"%v\">%v</string>\n", key, value)
}

func (android) Footer(lang goloc.Lang) string {
	return "</resources>\n"
}

func (android) ValidateFormat(format string) error {
	if strings.HasPrefix(format, `%`) {
		return errors.New(`format must not start with "%" - it will be added automatically`)
	}
	return nil
}

func (android) IndexedFormatString(index uint, format string) string {
	return fmt.Sprintf(`%%%v$%v`, index+1, format)
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
