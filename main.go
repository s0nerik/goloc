package main

import (
	"github.com/s0nerik/goloc/goloc"
	"github.com/s0nerik/goloc/platforms/resolver"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
)

var (
	credentials                = kingpin.Flag(`credentials`, `Credentials to access a spreadsheet.`).Short('c').Default(`client_secret.json`).String()
	platformName               = kingpin.Flag(`platform`, `Target platform name.`).Short('p').Required().String()
	sheetID                    = kingpin.Flag(`spreadsheet`, `Spreadsheet ID.`).Short('s').Required().String()
	resDir                     = kingpin.Flag(`resources`, `Path to the resources folder in the project.`).Short('r').Required().String()
	tabName                    = kingpin.Flag(`tab`, `Localizations tab name.`).Short('t').Required().String()
	keyColumn                  = kingpin.Flag(`key-column`, `Title of the key column.`).Default(`key`).String()
	stopOnMissing              = kingpin.Flag(`stop-on-missing`, `Stop execution if missing localization is found.`).Default(`false`).Bool()
	formatsTabName             = kingpin.Flag(`formats-tab`, `Formats tab name.`).Short('f').Default(`formats`).String()
	formatNameColumn           = kingpin.Flag(`format-name-column`, `Title of the format name column.`).Default(`format`).String()
	defFormatName              = kingpin.Flag(`default-format-name`, `Name of the format to be used in place of "{}"`).Default("").String()
	defLoc                     = kingpin.Flag(`default-localization`, `Default localization language (e.g. "en"). Specifying this doesn't have any effect if the "--default-localization-file-path" is not specified.`).Default(`en`).String()
	defLocPath                 = kingpin.Flag(`default-localization-file-path`, `Full path to the default localization file. Specify this if you want to write a default localization into a specific file (ignoring the localization path generation logic for a language specified in "--default-localization").`).String()
	missingLocalizationsReport = kingpin.Flag(`missing-localizations-report`, `Specify this flag if you want to only pretty-print missing localizations without generating the actual localization files.`).Default(`false`).Bool()
	emptyLocalizationMatch     = kingpin.Flag(`empty-localization-match`, `Regex for empty localization string.`).Default(`^$`).Regexp()
)

func main() {
	kingpin.Version("0.9.3")
	kingpin.Parse()

	platform, err := resolver.FindPlatform(*platformName)
	if err != nil {
		log.Fatalf(`Platform "%v" is not supported.`, *platformName)
	}

	goloc.Run(platform, *resDir, *credentials, *sheetID, *tabName, *keyColumn, *formatsTabName, *formatNameColumn, *defLoc, *defLocPath, *stopOnMissing, *missingLocalizationsReport, *defFormatName, *emptyLocalizationMatch)
}
