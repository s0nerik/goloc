package goloc

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/olekukonko/tablewriter"
)

type RawCell = string

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

func fetchEverythingRaw(source Source) (rawFormats, rawLocalizations [][]string, err error) {
	var formatsError error
	var localizationsError error

	var wg sync.WaitGroup
	go func() {
		defer wg.Done()
		rawFormats, formatsError = source.Formats()
	}()
	go func() {
		defer wg.Done()
		rawLocalizations, localizationsError = source.Localizations()
	}()

	wg.Add(2)
	wg.Wait()

	if formatsError != nil {
		return nil, nil, fmt.Errorf(`can't load formats (%w)`, formatsError)
	}
	if localizationsError != nil {
		return nil, nil, fmt.Errorf(`can't load localizations (%w)`, formatsError)
	}

	return
}

// Run launches the actual process of fetching, parsing and writing the localization files.
func Run(
	source Source,
	platform Platform,
	resDir string,
	keyColumn string,
	formatNameColumn string,
	defaultLocalization string,
	defaultLocalizationPath string,
	stopOnMissing bool,
	reportMissingLocalizations bool,
	defFormatName string,
	emptyLocalizationMatch *regexp.Regexp,
) error {
	rawFormats, rawLocalizations, err := fetchEverythingRaw(source)
	if err != nil {
		return fmt.Errorf(`can't fetch data, reason: %w`, err)
	}

	formats, err := ParseFormats(rawFormats, platform, source.FormatsDocumentName(), formatNameColumn, defFormatName)
	if err != nil {
		return err
	}

	localizations, fArgs, warn, err := ParseLocalizations(rawLocalizations, platform, formats, source.LocalizationsDocumentName(), keyColumn, stopOnMissing, emptyLocalizationMatch)
	if err != nil {
		return err
	}

	if reportMissingLocalizations {
		reportMissingLanguages(warn)
		return errors.New("found missing localizations")
	}

	for _, w := range warn {
		log.Println(w)
	}

	// Make sure we can access resources dir
	if _, err := os.Stat(resDir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(resDir, 0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if p, ok := platform.(Preprocessor); ok {
		err := p.Preprocess(PreprocessArgs{ResDir: resDir, Localizations: localizations, Formats: formats, FormatArgs: fArgs, DefaultLocalization: defaultLocalization})
		if err != nil {
			return err
		}
	}

	err = WriteLocalizations(platform, resDir, localizations, fArgs, defaultLocalization, defaultLocalizationPath)
	if err != nil {
		return fmt.Errorf(`can't write localizations, reason: %w`, err)
	}

	if p, ok := platform.(Postprocessor); ok {
		err := p.Postprocess(PostprocessArgs{ResDir: resDir, Localizations: localizations, Formats: formats, FormatArgs: fArgs, DefaultLocalization: defaultLocalization})
		if err != nil {
			return err
		}
	}

	return nil
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
