package goloc

import (
	"fmt"
	"strings"
)

func ParseFormats(
	rawData [][]interface{},
	platform Platform,
	formatsTabName string,
	formatColumnTitle string,
) (Formats, error) {
	formatColIndex, platformColIndex, actualPlatformName, err := columnIndices(platform, rawData, formatsTabName, formatColumnTitle)
	if err != nil {
		return nil, err
	}

	formats := Formats{}
	for rowIndex, row := range rawData[1:] {
		actualRowIndex := uint(rowIndex + 2)
		if formatColIndex >= len(row) {
			return nil, &formatKeyNotSpecifiedError{
				cell: *NewCell(formatsTabName, actualRowIndex, uint(formatColIndex)),
			}
		}
		if platformColIndex >= len(row) {
			return nil, &formatValueNotSpecifiedError{
				cell:         *NewCell(formatsTabName, actualRowIndex, uint(platformColIndex)),
				platformName: actualPlatformName,
			}
		}
		if key, ok := row[formatColIndex].(FormatKey); ok {
			if val, ok := row[platformColIndex].(string); ok {
				trimmedVal := strings.TrimSpace(val)
				if len(trimmedVal) == 0 {
					return nil, &formatValueNotSpecifiedError{
						cell:         *NewCell(formatsTabName, actualRowIndex, uint(platformColIndex)),
						platformName: actualPlatformName,
					}
				}
				err := platform.ValidateFormat(trimmedVal)
				if err != nil {
					return nil, &formatValueInvalidError{
						cell:         *NewCell(formatsTabName, actualRowIndex, uint(platformColIndex)),
						platformName: actualPlatformName,
						formatValue:  trimmedVal,
						reason:       err,
					}
				}
				formats[key] = trimmedVal
			} else {
				return nil, &wrongValueTypeError{
					cell: *NewCell(formatsTabName, actualRowIndex, uint(platformColIndex)),
				}
			}
		} else {
			return nil, &wrongKeyTypeError{
				cell: *NewCell(formatsTabName, actualRowIndex, uint(formatColIndex)),
			}
		}
	}

	return formats, nil
}

func columnIndices(
	platform Platform,
	rawData [][]interface{},
	formatsTabName string,
	formatColumnTitle string,
) (formatColIndex int, platformColIndex int, actualPlatformName string, err error) {
	formatColIndex = -1
	platformColIndex = -1
	actualPlatformName = ``

	if len(rawData) == 0 {
		err = &emptySheetError{tab: formatsTabName}
		return
	}

	firstRow := rawData[0]
	if len(firstRow) == 0 {
		err = &emptyFirstRowError{tab: formatsTabName}
		return
	}

	for i, val := range firstRow {
		if val == formatColumnTitle {
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
		err = &noFormatColumnError{tab: formatsTabName, requiredColumnTitle: formatColumnTitle}
	}

	if platformColIndex == -1 {
		err = &noPlatformColumnError{tab: formatsTabName, platformNames: platform.Names()}
	}

	return
}

// region Errors

type emptySheetError struct {
	tab string
}

type emptyFirstRowError struct {
	tab string
}

type noFormatColumnError struct {
	tab                 string
	requiredColumnTitle string
}

type noPlatformColumnError struct {
	tab           string
	platformNames []string
}

type formatKeyNotSpecifiedError struct {
	cell Cell
}

type formatValueNotSpecifiedError struct {
	platformName string
	cell         Cell
}

type formatValueInvalidError struct {
	cell         Cell
	platformName string
	formatValue  string
	reason       error
}

type wrongValueTypeError struct {
	cell Cell
}

type wrongKeyTypeError struct {
	cell Cell
}

func (e *emptySheetError) Error() string {
	return fmt.Sprintf(`%v!A1: sheet is empty`, e.tab)
}

func (e *emptyFirstRowError) Error() string {
	return fmt.Sprintf(`%v!A1: first row is required`, e.tab)
}

func (e *noFormatColumnError) Error() string {
	return fmt.Sprintf(`%v!A1: "%v" column is missing in the first row`, e.tab, e.requiredColumnTitle)
}

func (e *noPlatformColumnError) Error() string {
	return fmt.Sprintf(`%v!A1: can't find any of %v columns in the first row`, e.tab, e.platformNames)
}

func (e *formatKeyNotSpecifiedError) Error() string {
	return fmt.Sprintf(`%v: format name is not specified`, e.cell)
}

func (e *formatValueNotSpecifiedError) Error() string {
	return fmt.Sprintf(`%v: value for "%v" platform is not specified`, e.cell, e.platformName)
}

func (e *formatValueInvalidError) Error() string {
	return fmt.Sprintf(`%v: format "%v" is invalid for platform "%v" (%v)`, e.cell, e.formatValue, e.platformName, e.reason)
}

func (e *wrongValueTypeError) Error() string {
	return fmt.Sprintf(`%v: wrong value type`, e.cell)
}

func (e *wrongKeyTypeError) Error() string {
	return fmt.Sprintf(`%v: wrong key type`, e.cell)
}

// endregion
