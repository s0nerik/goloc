package goloc

import (
	"io/ioutil"
	"log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"golang.org/x/net/context"
	"regexp"
	"sync"
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
	IndexedFormatString(index uint, format string) (string, error)
	LangResPath(lang Lang, resDir ResDir) string
	LocalizationFileName(lang Lang) string
}

var formatRegexpInitializer sync.Once
var formatRegexp *regexp.Regexp

var langColumnRegexpInitializer sync.Once
var langColumnRegexp *regexp.Regexp

func FormatRegexp() *regexp.Regexp {
	formatRegexpInitializer.Do(func() {
		r, err := regexp.Compile(`\{([^\{^\}]*)\}`)
		if err != nil {
			log.Fatalf("Can't create a regexp for format string. Reason: %v. Please, submit an issue with the execution logs here: https://github.com/s0nerik/goloc", err)
		}
		formatRegexp = r
	})
	return formatRegexp
}

func LangColumnNameRegexp() *regexp.Regexp {
	langColumnRegexpInitializer.Do(func() {
		r, err := regexp.Compile("lang_([a-z]{2})")
		if err != nil {
			log.Fatalf("Can't create a regexp for lang column name. Reason: %v. Please, submit an issue with the execution logs here: https://github.com/s0nerik/goloc", err)
		}
		langColumnRegexp = r
	})
	return langColumnRegexp
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

func Run(
	platform Platform,
	resDir string,
	credFilePath string,
	sheetId string,
	tabName string,
	keyColumn string,
	formatsTabName string,
	formatNameColumn string,
	stopOnMissing bool,
) {
	api := sheetsApi(credFilePath)

	formats, err := ParseFormats(api, platform, sheetId, formatsTabName, formatNameColumn)
	if err != nil {
		log.Fatalf(`Can't parse formats from the "%v" tab. Reason: %v.`, formatsTabName, err)
	}

	loc, err := ParseLocalizations(api, platform, formats, sheetId, tabName, keyColumn, stopOnMissing)
	if err != nil {
		log.Fatalf(`Can't parse localizations from the "%v" tab. Reason: %v.`, tabName, err)
	}

	for k, v := range formats {
		log.Printf(`FORMAT %v: %v`, k, v)
	}

	for k, v := range loc {
		log.Printf(`LOCALIZATION %v`, k)
		for k, v := range v {
			log.Printf(`LOCALIZATION %v: %v`, k, v)
		}
	}
}