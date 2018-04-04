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
	firstRow := rawData[0]
	if firstRow == nil {
		return nil, &noFirstRowError{tab: formatsTabName}
	}

	var formatColIndex = -1
	var platformColIndex = -1
	var actualPlatformName = ``
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
		return nil, &noFormatColumnError{tab: formatsTabName, requiredColumnTitle: formatColumnTitle}
	}

	if platformColIndex == -1 {
		return nil, &noPlatformColumnError{tab: formatsTabName, platformNames: platform.Names()}
	}

	formats := Formats{}
	for rowIndex, row := range rawData[1:] {
		actualRowIndex := uint(rowIndex + 2)
		if formatColIndex >= len(row) {
			return nil, &formatKeyNotSpecifiedError{
				cell: cell{tab: formatsTabName, row: actualRowIndex, column: uint(formatColIndex)},
			}
		}
		if platformColIndex >= len(row) {
			return nil, &formatValueNotSpecifiedError{
				cell:         cell{tab: formatsTabName, row: actualRowIndex, column: uint(platformColIndex)},
				platformName: actualPlatformName,
			}
		}
		if key, ok := row[formatColIndex].(FormatKey); ok {
			if val, ok := row[platformColIndex].(string); ok {
				trimmedVal := strings.TrimSpace(val)
				if len(trimmedVal) == 0 {
					return nil, &formatValueNotSpecifiedError{
						cell:         cell{tab: formatsTabName, row: actualRowIndex, column: uint(platformColIndex)},
						platformName: actualPlatformName,
					}
				}
				err := platform.ValidateFormat(trimmedVal)
				if err != nil {
					return nil, &formatValueInvalidError{
						cell:         cell{tab: formatsTabName, row: actualRowIndex, column: uint(platformColIndex)},
						platformName: actualPlatformName,
						formatValue:  trimmedVal,
						reason:       err,
					}
				}
				formats[key] = trimmedVal
			} else {
				return nil, &wrongValueTypeError{
					cell: cell{tab: formatsTabName, row: actualRowIndex, column: uint(platformColIndex)},
				}
			}
		} else {
			return nil, &wrongKeyTypeError{
				cell: cell{tab: formatsTabName, row: actualRowIndex, column: uint(formatColIndex)},
			}
		}
	}

	return formats, nil
}

type noFirstRowError struct {
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
	cell cell
}

type formatValueNotSpecifiedError struct {
	platformName string
	cell         cell
}

type formatValueInvalidError struct {
	cell         cell
	platformName string
	formatValue  string
	reason       error
}

type wrongValueTypeError struct {
	cell cell
}

type wrongKeyTypeError struct {
	cell cell
}

func (e *noFirstRowError) Error() string {
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
