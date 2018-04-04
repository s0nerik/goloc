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

func (it *android) String() string {
	return "android"
}

func (it *android) Names() []string {
	return []string{
		"android",
		"Android",
	}
}

func (it *android) Header(lang goloc.Lang) string {
	return "<?xml version=\"1.0\" encoding=\"utf-8\" ?>\n<resources>\n"
}

func (it *android) Localization(lang goloc.Lang, key goloc.Key, value string) string {
	return fmt.Sprintf("\t<string name=\"%v\">%v</string>\n", key, value)
}

func (it *android) Footer(lang goloc.Lang) string {
	return "</resources>\n"
}

func (it *android) ValidateFormat(format string) error {
	if strings.HasPrefix(format, `%`) {
		return errors.New(`format must not start with "%" - it will be added automatically`)
	}
	return nil
}

func (it *android) IndexedFormatString(index uint, format string) string {
	return fmt.Sprintf(`%%%v$%v`, index+1, format)
}

func (it *android) LocalizationFileName(lang goloc.Lang) string {
	return "localized_strings.xml"
}

func (it *android) LocalizationDirPath(lang goloc.Lang, resDir goloc.ResDir) string {
	targetDir := fmt.Sprintf("values-%v", lang)
	if resDir != "" {
		return filepath.Join(resDir, targetDir)
	} else {
		return filepath.Join("src", "main", "res", targetDir)
	}
}

func (it *android) ReplacementChars() map[string]string {
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
