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
) (Localizations, error) {
	keyColIndex, langColumns, err := localizationColumnIndices(rawData, tabName, keyColumn)
	if err != nil {
		return nil, err
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
	if firstRow == nil {
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
