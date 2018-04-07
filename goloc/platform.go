package goloc

type Platform interface {
	// Returns platform names that can be used to identify it in the sheet.
	Names() []string

	// Returns replacement characters for any special character that needs to be guarded in the platform resources.
	ReplacementChars() map[string]string

	// Returns header text. Can be an empty string. Newlines must be included here if localization format requires them.
	Header(lang Lang) string

	// Returns actual localization binding for a given language. Newlines must be included here if localization format requires them.
	Localization(lang Lang, key Key, value string) string

	// Returns footer text. Can be an empty string. Newlines must be included here if localization format requires them.
	Footer(lang Lang) string

	// Returns nil if format is valid and non-nil error otherwise
	ValidateFormat(format string) error

	// Returns an actual format string taking the argument position into consideration.
	// Example 1: format strings on Android are positional (with position starting from 1). In this case invocation of IndexedFormatString(0, "s") would return "%1$s".
	// Example 2: format strings on iOS aren't positional. In this case invocation of IndexedFormatString(0, "@") would return "%@".
	IndexedFormatString(index uint, format string) string

	// Returns a full relative path to localization file for a given language.
	LocalizationFilePath(lang Lang, resDir ResDir) string
}