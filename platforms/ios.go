package platforms

import (
	"github.com/s0nerik/goloc/platforms/registry"
	"github.com/s0nerik/goloc/goloc"
	"fmt"
	"strings"
	"errors"
	"path/filepath"
)

func init() {
	registry.Register(&ios{})
}

type ios struct{}

func (ios) Names() []string {
	return []string{
		"ios",
		"iOS",
	}
}

func (ios) LocalizationFilePath(lang goloc.Lang, resDir goloc.ResDir) string {
	fileName := "Localizable.strings"
	targetDir := fmt.Sprintf("%v.lproj", lang)
	if resDir != "" {
		return filepath.Join(resDir, targetDir, fileName)
	} else {
		return filepath.Join("Resources", "Localization", targetDir, fileName)
	}
}

func (ios) Header(lang goloc.Lang) string {
	return ""
}

func (ios) Localization(lang goloc.Lang, key goloc.Key, value string) string {
	return fmt.Sprintf("\"%v\" = \"%v\";\n", key, value)
}

func (ios) Footer(lang goloc.Lang) string {
	return ""
}

func (ios) ValidateFormat(format string) error {
	if strings.HasPrefix(format, `%`) {
		return errors.New(`format must not start with "%" - it will be added automatically`)
	}
	return nil
}

func (ios) IndexedFormatString(index uint, format string) string {
	return fmt.Sprintf(`%%%v`, format)
}

func (ios) ReplacementChars() map[string]string {
	return nil
}
