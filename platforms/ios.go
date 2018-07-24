package platforms

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/s0nerik/goloc/goloc"
	"github.com/s0nerik/goloc/platforms/registry"
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
	}
	return filepath.Join("Resources", "Localization", targetDir, fileName)
}

func (ios) Header(args *goloc.HeaderArgs) string {
	return ""
}

func (ios) LocalizedString(args *goloc.LocalizedStringArgs) string {
	return fmt.Sprintf("\"%v\" = \"%v\";\n", args.Key, args.Value)
}

func (ios) Footer(args *goloc.FooterArgs) string {
	return ""
}

func (ios) ValidateFormat(format string) error {
	if strings.HasPrefix(format, `%`) {
		return errors.New(`format must not start with "%" - it will be added automatically`)
	}
	return nil
}

func (ios) FormatString(args *goloc.FormatStringArgs) string {
	return fmt.Sprintf(`%%%v`, args.Format)
}

func (ios) ReplacementChars() map[string]string {
	return map[string]string {
		`'`:  `\'`,
		`"`:  `\"`,
		"\n": `\n`,
	}
}
