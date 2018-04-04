package platforms

import (
	"github.com/s0nerik/goloc/platforms/registry"
	"github.com/s0nerik/goloc/goloc"
)

func init() {
	registry.Register(&ios{})
}

type ios struct{}

func (ios) String() string {
	panic("implement me")
}

func (ios) Names() []string {
	panic("implement me")
}

func (ios) ReplacementChars() map[string]string {
	panic("implement me")
}

func (ios) Header(lang goloc.Lang) string {
	panic("implement me")
}

func (ios) Localization(lang goloc.Lang, key goloc.Key, value string) string {
	panic("implement me")
}

func (ios) Footer(lang goloc.Lang) string {
	panic("implement me")
}

func (ios) ValidateFormat(format string) error {
	panic("implement me")
}

func (ios) IndexedFormatString(index uint, format string) string {
	panic("implement me")
}

func (ios) LocalizationDirPath(lang goloc.Lang, resDir goloc.ResDir) string {
	panic("implement me")
}

func (ios) LocalizationFileName(lang goloc.Lang) string {
	panic("implement me")
}
