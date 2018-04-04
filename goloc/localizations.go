package goloc

import (
	"log"
	"fmt"
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
) (Localizations, error) {
	firstRow := rawData[0]
	if firstRow == nil {
		return nil, &firstRowNotFoundError{Cell{tabName, uint(1), uint(0)}}
	}

	var keyColIndex = -1
	var langColumns = langColumns{}
	for i, val := range firstRow {
		if val == keyColumn {
			keyColIndex = i
		}
		lang := LangColumnNameRegexp().FindStringSubmatch(val.(string))
		if lang != nil {
			langColumns[i] = lang[1]
		}
	}

	if keyColIndex == -1 {
		return nil, &columnNotFoundError{Cell{tabName, uint(1), uint(keyColIndex)}, keyColumn}
	}

	if len(langColumns) == 0 {
		return nil, &langColumnsNotFoundError{Cell{tabName, uint(1), uint(0)}}
	}

	loc := Localizations{}
	for index, row := range rawData[1:] {
		actualRow := index + 2
		if keyColIndex >= len(row) || len(strings.TrimSpace(row[keyColIndex].(Key))) == 0 {
			if errorIfMissing {
				return nil, newKeyMissingError(tabName, actualRow, keyColIndex)
			} else {
				log.Println(newKeyMissingError(tabName, actualRow, keyColIndex))
				continue
			}
		}
		key := strings.TrimSpace(row[keyColIndex].(Key))
		if keyLoc, err := keyLocalizations(platform, formats, tabName, actualRow, row, key, langColumns, errorIfMissing); err == nil {
			loc[key] = keyLoc
		} else {
			return nil, err
		}
	}

	return loc, nil
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
) (map[Key]string, error) {
	keyLoc := map[Key]string{}
	for i, lang := range langColumns {
		if i < len(row) {
			val := strings.TrimSpace(row[i].(string))
			if len(val) > 0 {
				valWithoutSpecChars := withReplacedSpecialChars(platform, val)
				finalValue, err := withReplacedFormats(platform, valWithoutSpecChars, formats, tab, line, i)
				if err != nil {
					return nil, err
				}
				keyLoc[lang] = finalValue
			} else if errorIfMissing {
				return nil, newLocalizationMissingError(tab, line, i, key, lang)
			} else {
				log.Println(newLocalizationMissingError(tab, line, i, key, lang))
			}
		} else if errorIfMissing {
			return nil, newLocalizationMissingError(tab, line, i, key, lang)
		} else {
			log.Println(newLocalizationMissingError(tab, line, i, key, lang))
		}
	}
	return keyLoc, nil
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

type firstRowNotFoundError struct {
	cell Cell
}

func (e *firstRowNotFoundError) Error() string {
	return fmt.Sprintf(`%v: there's no first row in the tab`, e.cell)
}

type columnNotFoundError struct {
	cell   Cell
	column string
}

func (e *columnNotFoundError) Error() string {
	return fmt.Sprintf(`%v: "%v" column not found in the first row`, e.cell, e.column)
}

type langColumnsNotFoundError struct {
	cell Cell
}

func (e *langColumnsNotFoundError) Error() string {
	return fmt.Sprintf(`%v: language columns are not found`, e.cell)
}

type localizationMissingError struct {
	cell     Cell
	key      string
	lang     string
	platform Platform
}

func newLocalizationMissingError(tab string, row int, col int, key string, lang string) *localizationMissingError {
	return &localizationMissingError{
		cell: *NewCell(tab, uint(row), uint(col)),
		key:  key,
		lang: lang,
	}
}

func (e *localizationMissingError) Error() string {
	return fmt.Sprintf(`%v: "%v" is missing for "%v" language`, e.cell, e.key, e.lang)
}

type keyMissingError struct {
	cell Cell
}

func newKeyMissingError(tab string, row int, col int) *keyMissingError {
	return &keyMissingError{
		cell: *NewCell(tab, uint(row), uint(col)),
	}
}

func (e *keyMissingError) Error() string {
	return fmt.Sprintf(`%v: key name is missing, ignoring this string...`, e.cell)
}

type formatNotFoundError struct {
	cell       Cell
	formatName string
}

func (e *formatNotFoundError) Error() string {
	return fmt.Sprintf(`%v: no such format - "%v"`, e.cell, e.formatName)
}
