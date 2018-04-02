package platforms

import (
	"github.com/s0nerik/goloc/goloc"
	"fmt"
	"strings"
	"errors"
)

type Android struct{}

func (it *Android) Names() []string {
	return []string{
		"android",
		"Android",
	}
}

func (it *Android) Header(lang goloc.Lang) string {
	return "<?xml version=\"1.0\" encoding=\"utf-8\" ?>\n<resources>\n"
}

func (it *Android) Footer(lang goloc.Lang) string {
	return "</resources>\n"
}

func (it *Android) IndexedFormatString(index uint, format string) (string, error) {
	if strings.HasPrefix(format, `%`) {
		return ``, errors.New(`format string shouldn't start with "%" sign`)
	}
	return fmt.Sprintf(`%%%v$%v`, index+1, format), nil
}

func (it *Android) LocalizationFileName(lang goloc.Lang) string {
	return "localized_strings.xml"
}

func (it *Android) LangResPath(lang goloc.Lang, resDir goloc.ResDir) string {
	if resDir != "" {
		return fmt.Sprintf("%v/values-%v", resDir, lang)
	} else {
		return fmt.Sprintf("src/main/res/values-%v", lang)
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

