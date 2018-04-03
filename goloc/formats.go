package goloc

import (
	"errors"
	"fmt"
	"strings"
)

func ParseFormats(
	rawData [][]interface{},
	platform Platform,
	formatsTabName string,
	formatNameColumn string,
) (Formats, error) {
	firstRow := rawData[0]
	if firstRow == nil {
		return nil, errors.New(fmt.Sprintf(`there's no first row in the "%v" tab`, formatsTabName))
	}

	var formatColIndex = -1
	var platformColIndex = -1
	var actualPlatformName = ``
	for i, val := range firstRow {
		if val == formatNameColumn {
			formatColIndex = i
		}
		for _, name := range platform.Names() {
			if val == name {
				platformColIndex = i
				actualPlatformName = name
			}
		}
	}

	if formatColIndex == -1 {
		return nil, errors.New(fmt.Sprintf(`"%v" column not found in the first row of "%v" tab`, formatNameColumn, formatsTabName))
	}

	if platformColIndex == -1 {
		return nil, errors.New(fmt.Sprintf(`can't find any of %v columns in the first row of "%v" tab`, platform.Names(), formatsTabName))
	}

	formats := Formats{}
	for rowIndex, row := range rawData[1:] {
		actualRowIndex := uint(rowIndex + 2)
		if formatColIndex >= len(row) {
			return nil, &formatNameNotSpecifiedError{tab: formatsTabName, row: actualRowIndex}
		}
		if platformColIndex >= len(row) {
			return nil, &formatValueNotSpecifiedError{tab: formatsTabName, row: actualRowIndex, platformName: actualPlatformName}
		}
		if key, ok := row[formatColIndex].(FormatKey); ok {
			if val, ok := row[platformColIndex].(string); ok {
				trimmedVal := strings.TrimSpace(val)
				if len(trimmedVal) == 0 {
					return nil, &formatValueNotSpecifiedError{tab: formatsTabName, row: actualRowIndex, platformName: actualPlatformName}
				}
				formats[key] = trimmedVal
			} else {
				return nil, errors.New(fmt.Sprintf(`%v!%v: wrong value type`, formatsTabName, actualRowIndex))
			}
		} else {
			return nil, errors.New(fmt.Sprintf(`%v!%v: wrong key type`, formatsTabName, actualRowIndex))
		}
	}

	return formats, nil
}

type formatNameNotSpecifiedError struct {
	tab string
	row uint
}

func (e *formatNameNotSpecifiedError) Error() string {
	return fmt.Sprintf(`%v!%v: format name is not specified`, e.tab, e.row)
}

type formatValueNotSpecifiedError struct {
	tab          string
	row          uint
	platformName string
}

func (e *formatValueNotSpecifiedError) Error() string {
	return fmt.Sprintf(`%v!%v: value for "%v" platform is not specified`, e.tab, e.row, e.platformName)
}
