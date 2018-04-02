package platforms

import (
	"github.com/s0nerik/goloc/goloc"
	"fmt"
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

func (it *Android) IndexedFormatString(index int, format string) string {
	return "%"+string(index)+"$"+format
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

