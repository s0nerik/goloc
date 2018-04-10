package goloc

import (
	"log"
	"strings"
)

type langColumns = map[int]Lang

func ParseLocalizations(
	rawData [][]interface{},
	platform Platform,
	formats Formats,
	tabName string,
	keyColumn string,
	errorIfMissing bool,
) (loc Localizations, warnings []error, error error) {
	keyColIndex, langColumns, err := localizationColumnIndices(rawData, tabName, keyColumn)
	if err != nil {
		error = err
		return
	}

	loc = Localizations{}
	for index, row := range rawData[1:] {
		actualRow := index + 2
		if keyColIndex >= len(row) || len(strings.TrimSpace(row[keyColIndex].(Key))) == 0 {
			if errorIfMissing {
				error = newKeyMissingError(tabName, actualRow, keyColIndex)
				return
			} else {
				warnings = append(warnings, newKeyMissingError(tabName, actualRow, keyColIndex))
				continue
			}
		}
		key := strings.TrimSpace(row[keyColIndex].(Key))
		if keyLoc, warn, err := keyLocalizations(platform, formats, tabName, actualRow, row, key, langColumns, errorIfMissing); err == nil {
			if len(warn) > 0 {
				warnings = append(warnings, warn...)
			}
			loc[key] = keyLoc
		} else {
			error = err
			return
		}
	}

	return
}

func localizationColumnIndices(
	rawData [][]interface{},
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
		lang := LangColumnNameRegexp().FindStringSubmatch(val.(string))
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

func keyLocalizations(
	platform Platform,
	formats Formats,
	tab string,
	line int,
	row []interface{},
	key Key,
	langColumns langColumns,
	errorIfMissing bool,
) (keyLoc map[Key]string, warnings []error, error error) {
	keyLoc = map[Key]string{}
	for i, lang := range langColumns {
		if i < len(row) {
			val := strings.TrimSpace(row[i].(string))
			if len(val) > 0 {
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
	var index uint = 0
	var err error
	strWithReplacedFormats := FormatRegexp().ReplaceAllStringFunc(str, func(formatName string) string {
		if len(formatName) < 2 {
			log.Fatalf(`%v: something went wrong. Please submit an issue with the values in the problematic row.`, Cell{tab, uint(row), uint(column)})
		}

		name := formatName[1 : len(formatName)-1]
		if format, ok := formats[name]; ok {
			return platform.IndexedFormatString(index, format)
		} else {
			if err == nil {
				err = &formatNotFoundError{Cell{tab, uint(row), uint(column)}, name}
			}
		}

		index += 1
		return str
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
