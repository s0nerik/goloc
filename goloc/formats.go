package goloc

import (
	"google.golang.org/api/sheets/v4"
	"errors"
	"fmt"
)

func ParseFormats(
	api *sheets.SpreadsheetsService,
	platform Platform,
	sheetId string,
	formatsTabName string,
	formatNameColumn string,
) (Formats, error) {
	resp, err := api.Values.Get(sheetId, formatsTabName).Do()
	if err != nil {
		return nil, err
	}

	firstRow := resp.Values[0]
	if firstRow == nil {
		return nil, errors.New(fmt.Sprintf(`there's no first row in the "%v" tab`, formatsTabName))
	}

	var formatColIndex = -1
	var platformColIndex = -1
	for i, val := range firstRow {
		if val == formatNameColumn {
			formatColIndex = i
		}
		for _, name := range platform.Names() {
			if val == name {
				platformColIndex = i
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
	for _, row := range resp.Values[1:] {
		formats[row[formatColIndex].(FormatKey)] = row[platformColIndex].(string)
	}

	return formats, nil
}
