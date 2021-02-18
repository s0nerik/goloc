package goloc

import "time"

// LocalizedStringArgs encapsulates arguments to a function that returns the actual localized string for a given platform.
type LocalizedStringArgs struct {
	Index      int
	IsLast     bool
	Lang       Locale
	Key        Key
	Value      string
	FormatArgs []string
}

// FormatStringArgs encapsulates arguments to a function that returns the actual format specification for a given platform.
type FormatStringArgs struct {
	Index  int
	Format string
}

// HeaderArgs encapsulates arguments to a function that returns a localization file header for a given platform.
type HeaderArgs struct {
	Lang Locale
	Time time.Time
}

// FooterArgs encapsulates arguments to a function that returns a localization file footer for a given platform.
type FooterArgs struct {
	Lang Locale
}

// Platform represents an object responsible for specifying a format of the resulting localization file.
type Platform interface {
	// Returns platform names that can be used to identify it in the sheet.
	Names() []string

	// Returns a full relative path to localization file for a given language.
	LocalizationFilePath(lang Locale, resDir ResDir) string

	// Returns header text. Returned string can be empty. Newlines must be included here if localization format requires them.
	Header(args *HeaderArgs) string

	// Returns actual localization binding for a given language. Newlines must be included here if localization format requires them.
	LocalizedString(args *LocalizedStringArgs) string

	// Returns footer text. Can be an empty string. Newlines must be included here if localization format requires them.
	Footer(args *FooterArgs) string

	// Returns nil if format is valid and non-nil error otherwise
	ValidateFormat(format string) error

	// Returns an actual format string taking the argument position into consideration.
	// Example 1: format strings on Android are positional (with position starting from 1). In this case invocation of IndexedFormatString(0, "s") would return "%1$s".
	// Example 2: format strings on iOS aren't positional. In this case invocation of IndexedFormatString(0, "@") would return "%@".
	FormatString(args *FormatStringArgs) string

	// Returns replacement characters for any special character that needs to be guarded in the platform resources.
	ReplacementChars() map[string]string
}

type FallbackStringWriter interface {
	FallbackString(args *LocalizedStringArgs) string
}
