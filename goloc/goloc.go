package goloc

import (
	"io/ioutil"
	"log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"golang.org/x/net/context"
	"regexp"
)

type Lang = string
type Key = string
type Localizations = map[Key]map[Lang]string

type FormatKey = string
type Formats = map[FormatKey]string

type ResDir = string

type Platform interface{
	Names() []string
	ReplacementChars() map[string]string
	Header(lang Lang) string
	Footer(lang Lang) string
	IndexedFormatString(index int, format string) string
	LangResPath(lang Lang, resDir ResDir) string
	LocalizationFileName(lang Lang) string
}

func FormatRegexp() (*regexp.Regexp, error) {
	r, err := regexp.Compile(`\{([^\{^\}]*)\}`)
	if err != nil {
		log.Fatalf("Can't create a format regexp: %v", err)
	}
	return r, nil
}

func sheetsApi(credFilePath string) *sheets.SpreadsheetsService {
	ctx := context.Background()

	sec, err := ioutil.ReadFile(credFilePath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(sec, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	s, err := sheets.New(config.Client(ctx))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	return s.Spreadsheets
}

func Run(platform Platform, credFilePath string, sheetId string, tabName string, formatsTabName string) {
	api := sheetsApi(credFilePath)

	formats, err := ParseFormats(api, platform, sheetId, formatsTabName, "format")
	if err != nil {
		log.Fatalf(`Can't parse formats from the "%v" tab. Reason: %v.`, formatsTabName, err)
	}

	log.Println(formats)
}