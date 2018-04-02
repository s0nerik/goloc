package goloc

import (
	"google.golang.org/api/sheets/v4"
	"log"
	"errors"
	"fmt"
	"strings"
)

func ParseLocalizations(
	api *sheets.SpreadsheetsService,
	platform Platform,
	sheetId string,
	resDir string,
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
	var langColumns = map[int]Lang{}
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
	for line, row := range resp.Values[1:] {
		if keyColIndex >= len(row) || len(strings.TrimSpace(row[keyColIndex].(Key))) == 0 {
			if errorIfMissing {
				return nil, &keyMissingError{tab: tabName, line: line+2}
			} else {
				log.Println(&keyMissingError{tab: tabName, line: line+2})
				continue
			}
		}
		key := strings.TrimSpace(row[keyColIndex].(Key))
		keyLoc := map[Key]string{}
		for i, lang := range langColumns {
			if i < len(row) {
				val := strings.TrimSpace(row[i].(string))
				if len(val) != 0 {
					keyLoc[lang] = val
				} else if errorIfMissing {
					return nil, &localizationMissingError{key: key, lang: lang}
				} else {
					log.Println(&localizationMissingError{key: key, lang: lang})
				}
			} else if errorIfMissing {
				return nil, &localizationMissingError{key: key, lang: lang}
			} else {
				log.Println(&localizationMissingError{key: key, lang: lang})
			}
		}
		loc[key] = keyLoc
	}

	log.Println(fmt.Sprintf(`keyColIndex: %v, langColumns: %v`, keyColIndex, langColumns))

	return loc, nil
}

type localizationMissingError struct {
	key  string
	lang string
}

func (e *localizationMissingError) Error() string {
	return fmt.Sprintf(`"%v" is missing for "%v" language`, e.key, e.lang)
}

type keyMissingError struct {
	line int
	tab string
}

func (e *keyMissingError) Error() string {
	return fmt.Sprintf(`key name is missing for line %v in "%v" tab`, e.line, e.tab)
}