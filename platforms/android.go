package platforms

import (
	"github.com/s0nerik/goloc/goloc"
	"fmt"
	"strings"
	"errors"
	"path/filepath"
)

type Android struct {}

func (it *Android) String() string {
	return "android"
}

func (it *Android) Names() []string {
	return []string{
		"android",
		"Android",
	}
}

func (it *Android) Header(lang goloc.Lang) string {
	return "<?xml version=\"1.0\" encoding=\"utf-8\" ?>\n<resources>\n"
}

func (it *Android) Localization(lang goloc.Lang, key goloc.Key, value string) string {
	return fmt.Sprintf("\t<string name=\"%v\">%v</string>\n", key, value)
}

func (it *Android) Footer(lang goloc.Lang) string {
	return "</resources>\n"
}

func (it *Android) ValidateFormat(format string) error {
	if strings.HasPrefix(format, `%`) {
		return errors.New(`format must not start with "%" - it will be added automatically`)
	}
	return nil
}

func (it *Android) IndexedFormatString(index uint, format string) (string, error) {
	return fmt.Sprintf(`%%%v$%v`, index+1, format), nil
}

func (it *Android) LocalizationFileName(lang goloc.Lang) string {
	return "localized_strings.xml"
}

func (it *Android) LocalizationDirPath(lang goloc.Lang, resDir goloc.ResDir) string {
	targetDir := fmt.Sprintf("values-%v", lang)
	if resDir != "" {
		return filepath.Join(resDir, targetDir)
	} else {
		return filepath.Join("src", "main", "res", targetDir)
	}
}

func (it *Android) ReplacementChars() map[string]string {
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
