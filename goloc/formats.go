package goloc

import (
	"github.com/s0nerik/goloc/goloc/re"
	"strings"
)

// ParseFormats parses formats given the raw table data and returns, if successful, mappings
// to the actual platform format for each format name.
func ParseFormats(
	rawData [][]interface{},
	platform Platform,
	formatsTabName string,
	formatColumnTitle string,
	defFormatName string,
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

	// Handle default format ("{}")
	if defFormatName != "" {
		formats[""] = defFormatName
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
		err = &firstRowNotFoundError{Cell{formatsTabName, uint(1), 0}}
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

func FormatArgs(str string) []string {
	return re.FormatRegexp().FindAllString(str, -1)
}