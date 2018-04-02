package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/s0nerik/goloc/goloc"
	"github.com/s0nerik/goloc/platforms/resolver"
	"log"
)

var (
	credentials      = kingpin.Flag(`credentials`, `Credentials to access a spreadsheet.`).Short('c').Default(`client_secret.json`).String()
	platformName     = kingpin.Flag(`platform`, `Target platform name.`).Short('p').Required().String()
	sheetId          = kingpin.Flag(`spreadsheet`, `Spreadsheet ID.`).Short('s').Required().String()
	resDir           = kingpin.Flag(`resources`, `Path to the resources folder in the project.`).Short('r').Required().String()
	tabName          = kingpin.Flag(`tab`, `Localizations tab name.`).Short('t').Required().String()
	keyColumn        = kingpin.Flag(`key-column`, `Title of the key column`).Default(`key`).String()
	stopOnMissing    = kingpin.Flag(`stop-on-missing`, `Title of the key column`).Default(`false`).Bool()
	formatsTabName   = kingpin.Flag(`formats-tab`, `Formats tab name`).Short('f').Default(`formats`).String()
	formatNameColumn = kingpin.Flag(`format-name-column`, `Title of the format name column`).Default(`format`).String()
	defLoc           = kingpin.Flag(`default-localization`, `Default localization language (e.g. "en").`).Default(`en`).String()
	defLocPath       = kingpin.Flag(`default-localization-path`, `Full path to default localization file.`).String()
	preferDefLoc     = kingpin.Flag(`prefer-default-localization`, `Don't provide language-specific resources for default language.`).Default(`true`).Bool()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	platform, err := resolver.FindPlatform(*platformName)
	if err != nil {
		log.Fatalf("Platform %v is not supported.", *platformName)
	}

	goloc.Run(platform, *resDir, *credentials, *sheetId, *tabName, *keyColumn, *formatsTabName, *formatNameColumn, *stopOnMissing)
}
