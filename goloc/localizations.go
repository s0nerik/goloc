package goloc

import (
	"google.golang.org/api/sheets/v4"
	"log"
	"errors"
	"fmt"
	"strings"
)

type langColumns = map[int]Lang

func ParseLocalizations(
	api *sheets.SpreadsheetsService,
	platform Platform,
	formats Formats,
	sheetId string,
	tabName string,
	keyColumn string,
	errorIfMissing bool,
) (Localizations, error) {
	resp, err := api.Values.Get(sheetId, tabName).Do()
	if err != nil {
		return nil, err
	}

	firstRow := resp.Values[0]
	if firstRow == nil {
		return nil, errors.New(fmt.Sprintf(`there's no first row in the "%v" tab`, tabName))
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
		return nil, errors.New(fmt.Sprintf(`"%v" column not found in the first row of "%v" tab`, keyColumn, tabName))
	}

	if len(langColumns) == 0 {
		return nil, errors.New(fmt.Sprintf(`language columns are not found in the "%v" tab`, tabName))
	}

	loc := Localizations{}
	for index, row := range resp.Values[1:] {
		actualRow := index + 2
		if keyColIndex >= len(row) || len(strings.TrimSpace(row[keyColIndex].(Key))) == 0 {
			if errorIfMissing {
				return nil, &keyMissingError{tab: tabName, line: actualRow}
			} else {
				log.Println(&keyMissingError{tab: tabName, line: actualRow})
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
				finalValue := withReplacedFormats(platform, valWithoutSpecChars, formats, tab, line)
				keyLoc[lang] = finalValue
			} else if errorIfMissing {
				return nil, &localizationMissingError{tab: tab, line: line, key: key, lang: lang}
			} else {
				log.Println(&localizationMissingError{tab: tab, line: line, key: key, lang: lang})
			}
		} else if errorIfMissing {
			return nil, &localizationMissingError{tab: tab, line: line, key: key, lang: lang}
		} else {
			log.Println(&localizationMissingError{tab: tab, line: line, key: key, lang: lang})
		}
	}
	return keyLoc, nil
}

func withReplacedFormats(platform Platform, str string, formats Formats, tab string, line int) string {
	var index uint = 0
	return FormatRegexp().ReplaceAllStringFunc(str, func(formatName string) string {
		if len(formatName) < 2 {
			log.Fatalf(`%v!%v: something went wrong. Please submit an issue with the values in the problematic row.`, tab, line)
		}

		name := formatName[1 : len(formatName)-1]
		if format, ok := formats[name]; ok {
			str, err := platform.IndexedFormatString(index, format)
			if err != nil {
				log.Fatalf(`%v!%v: can't use the "%v" format. Reason: %v`, tab, line, name, err)
			}
			return str
		} else {
			log.Fatalf(`%v!%v: no such format - "%v".`, tab, line, name)
		}

		index += 1
		return str
	})
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

type localizationMissingError struct {
	tab  string
	line int
	key  string
	lang string
}

func (e *localizationMissingError) Error() string {
	return fmt.Sprintf(`%v!%v: "%v" is missing for "%v" language`, e.tab, e.line, e.key, e.lang)
}

type keyMissingError struct {
	tab  string
	line int
}

func (e *keyMissingError) Error() string {
	return fmt.Sprintf(`%v!%v: key name is missing, ignoring this string...`, e.tab, e.line)
}
