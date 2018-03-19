package goloc

type Lang = string
type Key = string
type Localizations = map[Key]map[Lang]string

type FormatKey = string
type Formats = map[FormatKey]string

type ResDir = string

type Platform interface{
	Names() []string
	ReplacementChars() map[string]string
	Header(lang Lang) string
	Footer(lang Lang) string
	IndexedFormatString(index int, format string) string
	LangResPath(lang Lang, resDir ResDir) string
	LocalizationFileName(lang Lang) string
}