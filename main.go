package main

import (
	"fmt"
	"log"

	"github.com/s0nerik/goloc/goloc"
	"github.com/s0nerik/goloc/registry"
	"github.com/s0nerik/goloc/sources"
	"gopkg.in/alecthomas/kingpin.v2"

	// Register all supported platforms
	_ "github.com/s0nerik/goloc/platforms"
)

const remoteSources = `google_sheets`
const localSources = `csv`

var availableSources = fmt.Sprintf(`%v, %v`, remoteSources, localSources)

var (
	// Basic params
	source       = kingpin.Flag(`source`, fmt.Sprintf(`Data source. Available sources: %v`, availableSources)).Default(`google_sheets`).String()
	platformName = kingpin.Flag(`platform`, `Target platform name.`).Short('p').Required().String()
	resDir       = kingpin.Flag(`resources`, `Path to the resources folder in the project.`).Short('r').Required().String()

	// Local source params
	locFilePath     = kingpin.Flag(`localizations-file-path`, fmt.Sprintf(`Localizations file path. Required for sources: %v`, localSources)).String()
	formatsFilePath = kingpin.Flag(`formats-file-path`, fmt.Sprintf(`Formats file path. Required for sources: %v`, localSources)).String()

	// Google Sheets params
	sheetID        = kingpin.Flag(`spreadsheet`, `Spreadsheet ID. Required if selected source is 'google_sheets'`).Short('s').String()
	credentials    = kingpin.Flag(`credentials`, `Credentials to access a spreadsheet.`).Short('c').Default(`client_secret.json`).String()
	tabName        = kingpin.Flag(`tab`, `Localizations tab name.`).Short('t').Default(`localizations`).String()
	formatsTabName = kingpin.Flag(`formats-tab`, `Formats tab name.`).Short('f').Default(`formats`).String()

	// Advanced configuration
	keyColumn              = kingpin.Flag(`key-column`, `Title of the key column.`).Default(`key`).String()
	stopOnMissing          = kingpin.Flag(`stop-on-missing`, `Stop execution if missing localization is found.`).Default(`false`).Bool()
	formatNameColumn       = kingpin.Flag(`format-name-column`, `Title of the format name column.`).Default(`format`).String()
	defFormatName          = kingpin.Flag(`default-format-name`, `Name of the format to be used in place of "{}"`).Default("").String()
	defLoc                 = kingpin.Flag(`default-localization`, `Default localization language (e.g. "en"). Specifying this doesn't have any effect if the "--default-localization-file-path" is not specified.`).Default(`en`).String()
	defLocPath             = kingpin.Flag(`default-localization-file-path`, `Full path to the default localization file. Specify this if you want to write a default localization into a specific file (ignoring the localization path generation logic for a language specified in "--default-localization").`).String()
	emptyLocalizationMatch = kingpin.Flag(`empty-localization-match`, `Regex for empty localization string.`).Default(`^$`).Regexp()

	// Extra features
	missingLocalizationsReport = kingpin.Flag(`missing-localizations-report`, `Specify this flag if you want to only pretty-print missing localizations without generating the actual localization files.`).Default(`false`).Bool()
)

func main() {
	kingpin.Version("0.9.9")
	kingpin.Parse()

	platform := registry.GetPlatform(*platformName)
	if platform == nil {
		log.Fatalf(`Platform "%v" is not supported.`, *platformName)
	}

	src := resolveSource()
	if src == nil {
		log.Fatalf(`"%v" is not a supported source. Supported sources: %v.`, *source, availableSources)
	}

	err := goloc.Run(
		src,
		platform,
		*resDir,
		*keyColumn,
		*formatNameColumn,
		*defLoc,
		*defLocPath,
		*stopOnMissing,
		*missingLocalizationsReport,
		*defFormatName,
		*emptyLocalizationMatch,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func resolveSource() goloc.Source {
	switch *source {
	case "google_sheets":
		if sheetID == nil {
			log.Fatalf(`"--spreadsheet" parameter must be specified`)
		}
		if credentials == nil {
			log.Fatalf(`"--credentials" parameter must be specified`)
		}
		if *tabName == "" {
			log.Fatalf(`"--tab" parameter cannot be empty`)
		}
		if *formatsTabName == "" {
			log.Fatalf(`"--formats-tab" parameter cannot be empty`)
		}

		source, err := sources.GoogleSheets(*credentials, *sheetID, *formatsTabName, *tabName)
		if err != nil {
			log.Fatalf("can't create googlesheets source, %v", err.Error())
		}

		return source
	case "csv":
		if locFilePath == nil {
			log.Fatalf(`"--localizations-file-path" parameter must be specified`)
		}
		if *locFilePath == "" {
			log.Fatalf(`"--localizations-file-path" must be a valid file path`)
		}
		if formatsFilePath == nil {
			log.Fatalf(`"--formats-file-path" parameter must be specified`)
		}
		if *formatsFilePath == "" {
			log.Fatalf(`"--formats-file-path" must be a valid file path`)
		}
		return sources.CSV(*locFilePath, *formatsFilePath)
	}
	return nil
}
