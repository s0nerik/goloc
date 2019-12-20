package re

import (
	"fmt"
	"regexp"
	"sync"
)

var formatRegexpInitializer sync.Once
var formatRegexp *regexp.Regexp

// FormatRegexp returns a regexp that is used to distinguish between a normal text and a format occurrence.
// Each occurrence of a "{format_name}" in a localized string is considered a format occurrence and is going to be
// replaced with the actual platform-specific format definition at the time of parsing.
func FormatRegexp() *regexp.Regexp {
	formatRegexpInitializer.Do(func() {
		r, err := regexp.Compile(`\{([^\{^\}]*)\}`)
		if err != nil {
			panic(fmt.Errorf("Can't create a regexp for format string. Reason: %w. Please, submit an issue with the execution logs here: https://github.com/s0nerik/goloc", err))
		}
		formatRegexp = r
	})
	return formatRegexp
}

var langColumnRegexpInitializer sync.Once
var langColumnRegexp *regexp.Regexp

// LangColumnNameRegexp returns a regexp that is used to distinguish between localization columns and non-localization columns.
func LangColumnNameRegexp() *regexp.Regexp {
	langColumnRegexpInitializer.Do(func() {
		r, err := regexp.Compile("lang_([a-z]{2})")
		if err != nil {
			panic(fmt.Errorf("Can't create a regexp for lang column name. Reason: %w. Please, submit an issue with the execution logs here: https://github.com/s0nerik/goloc", err))
		}
		langColumnRegexp = r
	})
	return langColumnRegexp
}

var sprintRegexpInitializer sync.Once
var sprintfRegexp *regexp.Regexp

// SprintfRegexp returns a regexp that is used to distinguish between value type for type-safe format arguments
// Group #5 includes value type information
func SprintfRegexp() *regexp.Regexp {
	sprintRegexpInitializer.Do(func() {
		r, err := regexp.Compile(`%(?:(\d+)\$)?([\+\-\#0 ]*)(\d+|\*)?(?:\.(\d+|\*))?([a-z%])`)
		if err != nil {
			panic(fmt.Errorf("Can't create a regext for 'sprintf' format. Reason: %v. Please, submit an issue with the execution logs here: https://github.com/s0nerik/goloc", err))
		}
		sprintfRegexp = r
	})
	return sprintfRegexp
}
