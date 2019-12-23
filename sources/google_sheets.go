package sources

import (
	"fmt"
	"io/ioutil"

	"github.com/s0nerik/goloc/goloc"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type googleSheets struct {
	sheetID          string
	formatsTab       string
	localizationsTab string

	sheetsAPI *sheets.SpreadsheetsService
}

func GoogleSheets(
	credFilePath string,
	sheetID string,
	formatsTab string,
	localizationsTab string,
) (*googleSheets, error) {
	sheetsAPI, err := sheetsAPI(credFilePath)
	if err != nil {
		return nil, err
	}

	return &googleSheets{
		sheetID:          sheetID,
		formatsTab:       formatsTab,
		localizationsTab: localizationsTab,
		sheetsAPI:        sheetsAPI,
	}, nil
}

func sheetsAPI(credFilePath string) (*sheets.SpreadsheetsService, error) {
	ctx := context.Background()

	sec, err := ioutil.ReadFile(credFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	config, err := google.JWTConfigFromJSON(sec, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	s, err := sheets.New(config.Client(ctx))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets Client: %w", err)
	}

	return s.Spreadsheets, nil
}

func fetchRawValues(api *sheets.SpreadsheetsService, sheetID string, tab string) ([][]interface{}, error) {
	resp, err := api.Values.Get(sheetID, tab).Do()
	if err != nil {
		return nil, err
	}
	return resp.Values, nil
}

func fetchRawStringValues(api *sheets.SpreadsheetsService, sheetID string, tab string) ([][]string, error) {
	values, err := fetchRawValues(api, sheetID, tab)
	result := make([][]string, len(values))
	for i, row := range values {
		result[i] = make([]string, len(row))
		for j, col := range row {
			result[i][j] = col.(string)
		}
	}
	return result, err
}

func (s googleSheets) FormatsDocumentName() string {
	return s.formatsTab
}

func (s googleSheets) LocalizationsDocumentName() string {
	return s.localizationsTab
}

func (s googleSheets) Formats() ([][]goloc.RawCell, error) {
	return fetchRawStringValues(s.sheetsAPI, s.sheetID, s.formatsTab)
}

func (s googleSheets) Localizations() ([][]goloc.RawCell, error) {
	return fetchRawStringValues(s.sheetsAPI, s.sheetID, s.localizationsTab)
}
