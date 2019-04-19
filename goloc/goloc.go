package goloc

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// Lang represents a language code.
type Lang = string

// Key represents a localized string key.
type Key = string

// Localizations represents a mapping between a localized string key and it's values for different languages.
type Localizations map[Key]map[Lang]string

// LocalizationFormatArgs represents a mapping between a localized string key and its format arguments.
type LocalizationFormatArgs map[Key][]FormatKey

// FormatKey represents a name of the format.
type FormatKey = string

// Formats represents a mapping between format names and platform-specific format descriptions.
type Formats = map[FormatKey]string

// ResDir represents a resources directory path.
type ResDir = string

func sheetsAPI(credFilePath string) *sheets.SpreadsheetsService {
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

func fetchRawValues(api *sheets.SpreadsheetsService, sheetID string, tab string) ([][]interface{}, error) {
	resp, err := api.Values.Get(sheetID, tab).Do()
	if err != nil {
		return nil, err
	}
	return resp.Values, nil
}

func fetchEverythingRaw(
	api *sheets.SpreadsheetsService,
	sheetID string,
	formatsTab string,
	localizationsTab string,
) (rawFormats, rawLocalizations [][]interface{}, err error) {
	var formatsError error
	var localizationsError error

	var wg sync.WaitGroup
	go func() {
		defer wg.Done()
		rawFormats, formatsError = fetchRawValues(api, sheetID, formatsTab)
	}()
	go func() {
		defer wg.Done()
		rawLocalizations, localizationsError = fetchRawValues(api, sheetID, localizationsTab)
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

// Run launches the actual process of fetching, parsing and writing the localization files.
func Run(
	platform Platform,
	resDir string,
	credFilePath string,
	sheetID string,
	tabName string,
	keyColumn string,
	formatsTabName string,
	formatNameColumn string,
	defaultLocalization string,
	defaultLocalizationPath string,
	stopOnMissing bool,
	reportMissingLocalizations bool,
	defFormatName string,
	emptyLocalizationMatch *regexp.Regexp,
) {
	api := sheetsAPI(credFilePath)

	rawFormats, rawLocalizations, err := fetchEverythingRaw(api, sheetID, formatsTabName, tabName)
	if err != nil {
		log.Fatalf(`Can't fetch data from "%v" sheet. Reason: %v.`, sheetID, err)
	}

	formats, err := ParseFormats(rawFormats, platform, formatsTabName, formatNameColumn, defFormatName)
	if err != nil {
		log.Fatal(err)
	}

	localizations, fArgs, warn, err := ParseLocalizations(rawLocalizations, platform, formats, tabName, keyColumn, stopOnMissing, emptyLocalizationMatch)
	if err != nil {
		log.Fatal(err)
	} else {
		if reportMissingLocalizations {
			reportMissingLanguages(warn)
			return
		}
		for _, w := range warn {
			log.Println(w)
		}
	}

	// Make sure we can access resources dir
	if _, err := os.Stat(resDir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(resDir, 0755)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	if p, ok := platform.(Preprocessor); ok {
		err := p.Preprocess(PreprocessArgs{ResDir: resDir, Localizations: localizations, Formats: formats, FormatArgs: fArgs, DefaultLocalization: defaultLocalization})
		if err != nil {
			log.Fatal(err)
		}
	}

	err = WriteLocalizations(platform, resDir, localizations, fArgs, defaultLocalization, defaultLocalizationPath)
	if err != nil {
		log.Fatalf(`Can't write localizations. Reason: %v.`, err)
	}

	if p, ok := platform.(Postprocessor); ok {
		err := p.Postprocess(PostprocessArgs{ResDir: resDir, Localizations: localizations, Formats: formats, FormatArgs: fArgs, DefaultLocalization: defaultLocalization})
		if err != nil {
			log.Fatal(err)
		}
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
