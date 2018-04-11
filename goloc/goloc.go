package goloc

import (
	"io/ioutil"
	"log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"golang.org/x/net/context"
	"regexp"
	"sync"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"strings"
	"sort"
)

type Lang = string
type Key = string
type Localizations = map[Key]map[Lang]string

type FormatKey = string
type Formats = map[FormatKey]string

type ResDir = string

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

func fetchRawValues(api *sheets.SpreadsheetsService, sheetId string, tab string) ([][]interface{}, error) {
	resp, err := api.Values.Get(sheetId, tab).Do()
	if err != nil {
		return nil, err
	}
	return resp.Values, nil
}

func fetchEverythingRaw(
	api *sheets.SpreadsheetsService,
	sheetId string,
	formatsTab string,
	localizationsTab string,
) (rawFormats, rawLocalizations [][]interface{}, err error) {
	var formatsError error
	var localizationsError error

	var wg sync.WaitGroup
	go func() {
		defer wg.Done()
		rawFormats, formatsError = fetchRawValues(api, sheetId, formatsTab)
	}()
	go func() {
		defer wg.Done()
		rawLocalizations, localizationsError = fetchRawValues(api, sheetId, localizationsTab)
	}()

	wg.Add(2)
	wg.Wait()

	if formatsError != nil {
		return nil, nil, formatsError
	}
	if localizationsError != nil {
		return nil, nil, localizationsError
	}

	return
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
	defaultLocalization string,
	defaultLocalizationPath string,
	stopOnMissing bool,
	reportMissingLocalizations bool,
) {
	api := sheetsApi(credFilePath)

	rawFormats, rawLocalizations, err := fetchEverythingRaw(api, sheetId, formatsTabName, tabName)
	if err != nil {
		log.Fatalf(`Can't fetch data from "%v" sheet. Reason: %v.`, sheetId, err)
	}

	formats, err := ParseFormats(rawFormats, platform, formatsTabName, formatNameColumn)
	if err != nil {
		log.Fatal(err)
	}

	localizations, warn, err := ParseLocalizations(rawLocalizations, platform, formats, tabName, keyColumn, stopOnMissing)
	if err != nil {
		log.Fatal(err)
	} else {
		if reportMissingLocalizations {
			reportMissingLanguages(warn)
			return
		} else {
			for _, w := range warn {
				log.Println(w)
			}
		}
	}

	err = WriteLocalizations(platform, resDir, localizations, defaultLocalization, defaultLocalizationPath)
	if err != nil {
		log.Fatalf(`Can't write localizations. Reason: %v.`, err)
	}
}

func reportMissingLanguages(warnings []error) {
	rowWarnings := map[uint][]*localizationMissingError{}
	for _, w := range warnings {
		if w, ok := w.(*localizationMissingError); ok {
			rowWarnings[w.cell.row] = append(rowWarnings[w.cell.row], w)
		}
	}

	type kv struct {
		row      uint
		warnings []*localizationMissingError
	}

	var sortedRowWarnings []kv
	for k, v := range rowWarnings {
		sortedRowWarnings = append(sortedRowWarnings, kv{k, v})
	}

	sort.Slice(sortedRowWarnings, func(i, j int) bool {
		return sortedRowWarnings[i].row < sortedRowWarnings[j].row
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Row", "Key", "Missing localizations"})
	for _, kv := range sortedRowWarnings {
		row := kv.warnings[0].cell.row
		key := kv.warnings[0].key

		var missingLanguages []string
		for _, w := range kv.warnings {
			missingLanguages = append(missingLanguages, w.lang)
		}

		table.Append([]string{strconv.Itoa(int(row)), key, strings.Join(missingLanguages, ",")})
	}
	table.Render()
}
