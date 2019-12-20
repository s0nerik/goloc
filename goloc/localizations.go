package goloc

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/s0nerik/goloc/goloc/re"
)

type langColumns = map[int]Lang

var DefaultEmptyLocRegexp, _ = regexp.Compile("^$")

// ParseLocalizations parses formats given the raw table data and returns, if successful, mappings
// for each localized string in different languages.
func ParseLocalizations(
	rawData [][]RawCell,
	platform Platform,
	formats Formats,
	tabName string,
	keyColumn string,
	errorIfMissing bool,
	emptyLocalizationRegexp *regexp.Regexp,
) (loc Localizations, formatArgs LocalizationFormatArgs, warnings []error, error error) {
	formatArgs = LocalizationFormatArgs{}

	if emptyLocalizationRegexp == nil {
		emptyLocalizationRegexp = DefaultEmptyLocRegexp
	}

	keyColIndex, langCols, err := localizationColumnIndices(rawData, tabName, keyColumn)
	if err != nil {
		error = err
		return
	}

	loc = Localizations{}
	for index, row := range rawData[1:] {
		actualRow := index + 2
		if keyColIndex >= len(row) || len(strings.TrimSpace(row[keyColIndex])) == 0 {
			if errorIfMissing {
				error = newKeyMissingError(tabName, actualRow, keyColIndex)
				return
			}
			warnings = append(warnings, newKeyMissingError(tabName, actualRow, keyColIndex))
			continue
		}
		key := strings.TrimSpace(row[keyColIndex])
		if keyLoc, warn, err := keyLocalizations(platform, formats, tabName, actualRow, row, key, langCols, errorIfMissing, emptyLocalizationRegexp); err == nil {
			if len(warn) > 0 {
				warnings = append(warnings, warn...)
			}
			loc[key] = keyLoc
		} else {
			error = err
			return
		}

		formatArgs[key], error = keyFormatArgs(platform, tabName, actualRow, row, key, langCols, emptyLocalizationRegexp)
		if error != nil {
			return
		}
	}

	return
}

func localizationColumnIndices(
	rawData [][]RawCell,
	tabName string,
	keyColumn string,
) (keyColIndex int, langCols langColumns, err error) {
	keyColIndex = -1
	langCols = langColumns{}

	if len(rawData) == 0 {
		err = &emptySheetError{tab: tabName}
		return
	}

	firstRow := rawData[0]
	if len(firstRow) == 0 {
		err = &firstRowNotFoundError{Cell{tabName, uint(1), uint(0)}}
		return
	}

	for i, val := range firstRow {
		if val == keyColumn {
			keyColIndex = i
		}
		lang := re.LangColumnNameRegexp().FindStringSubmatch(val)
		if lang != nil {
			langCols[i] = lang[1]
		}
	}

	if keyColIndex == -1 {
		err = &columnNotFoundError{Cell{tabName, uint(1), uint(keyColIndex)}, keyColumn}
		return
	}

	if len(langCols) == 0 {
		err = &langColumnsNotFoundError{Cell{tabName, uint(1), uint(0)}}
		return
	}

	return
}

func keyFormatArgs(
	platform Platform,
	tab string,
	line int,
	row []string,
	key Key,
	langColumns langColumns,
	emptyLocalizationRegexp *regexp.Regexp,
) (formatArgs []string, err error) {
	langFormatArgs := map[int][]string{}

	for i := range langColumns {
		if i < len(row) {
			val := strings.TrimSpace(row[i])
			valWithoutSpecChars := withReplacedSpecialChars(platform, val)
			if !emptyLocalizationRegexp.MatchString(valWithoutSpecChars) {
				langFormatArgs[i] = FormatArgs(valWithoutSpecChars)
			}
		}
	}

	if len(langFormatArgs) > 1 {
		first := true
		var prevArgs []string
		for col, fArgs := range langFormatArgs {
			if !first && !reflect.DeepEqual(fArgs, prevArgs) {
				err = newFormatArgsDifferentError(tab, line, col, key, langColumns[col])
				return
			}
			prevArgs = fArgs
			first = false
		}
	}

	if len(langFormatArgs) > 0 {
		for _, v := range langFormatArgs {
			formatArgs = v
			break
		}
	}

	return
}

func keyLocalizations(
	platform Platform,
	formats Formats,
	tab string,
	line int,
	row []string,
	key Key,
	langColumns langColumns,
	errorIfMissing bool,
	emptyLocalizationRegexp *regexp.Regexp,
) (keyLoc map[Lang]string, warnings []error, error error) {
	keyLoc = map[Lang]string{}
	for i, lang := range langColumns {
		keyLoc[lang] = ""
		if i < len(row) {
			val := strings.TrimSpace(row[i])
			if match := emptyLocalizationRegexp.MatchString(val); !match {
				valWithoutSpecChars := withReplacedSpecialChars(platform, val)
				finalValue, err := withReplacedFormats(platform, valWithoutSpecChars, formats, tab, line, i)
				if err != nil {
					error = err
					return
				}
				keyLoc[lang] = finalValue
			} else if errorIfMissing {
				error = newLocalizationMissingError(tab, line, i, key, lang)
				return
			} else {
				warnings = append(warnings, newLocalizationMissingError(tab, line, i, key, lang))
			}
		} else if errorIfMissing {
			error = newLocalizationMissingError(tab, line, i, key, lang)
			return
		} else {
			warnings = append(warnings, newLocalizationMissingError(tab, line, i, key, lang))
		}
	}
	return
}

func withReplacedFormats(platform Platform, str string, formats Formats, tab string, row int, column int) (string, error) {
	var index int
	var err error
	formatStringArgs := &FormatStringArgs{}
	strWithReplacedFormats := re.FormatRegexp().ReplaceAllStringFunc(str, func(formatName string) string {
		defer func() { index++ }()
		if len(formatName) < 2 {
			err = fmt.Errorf(`%v: something went wrong, please submit an issue with the values in the problematic row`, Cell{tab, uint(row), uint(column)})
			return ""
		}

		name := formatName[1 : len(formatName)-1]
		// Check if format specification exist and report an error if not
		if _, ok := formats[name]; !ok {
			if err == nil {
				err = &formatNotFoundError{Cell{tab, uint(row), uint(column)}, name}
			}
			return ""
		}

		formatStringArgs.Index = index
		formatStringArgs.Format = formats[name]
		return platform.FormatString(formatStringArgs)
	})

	return strWithReplacedFormats, err
}

func withReplacedSpecialChars(platform Platform, str string) string {
	specChars := platform.ReplacementChars()

	replacements := make([]string, 0, len(specChars))
	for orig, repl := range specChars {
		replacements = append(replacements, orig)
		replacements = append(replacements, repl)
	}

	return strings.NewReplacer(replacements...).Replace(str)
}

func (loc Localizations) Count() map[Lang]int {
	result := map[Lang]int{}
	for _, keyLoc := range loc {
		for lang := range keyLoc {
			result[lang]++
		}
	}
	return result
}

func (loc Localizations) Locales() (locales []string) {
	langs := map[Lang]bool{}
	for _, v := range loc {
		for lang := range v {
			langs[lang] = true
		}
	}

	for k := range langs {
		locales = append(locales, k)
	}

	return locales
}

func (loc Localizations) SortedKeys() (sortedKeys []Key) {
	for k := range loc {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}
