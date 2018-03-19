package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	credentials  = kingpin.Flag("credentials", "Credentials to access a spreadsheet.").Short('c').Default("client_secret.json").String()
	platform     = kingpin.Flag("platform", "Target platform name.").Short('p').Required().String()
	sheetId      = kingpin.Flag("spreadsheet", "Spreadsheet ID.").Short('s').Required().String()
	resDir       = kingpin.Flag("resources", "Path to the resources folder in the project.").Short('r').Required().String()
	tabName      = kingpin.Flag("tab", "Localizations tab name.").Short('t').Required().String()
	defLoc       = kingpin.Flag("default-localization", "Default localization language (e.g. \"en\").").Default("en").String()
	defLocPath   = kingpin.Flag("default-localization-path", "Full path to default localization file.").String()
	preferDefLoc = kingpin.Flag("prefer-default-localization", "Don't provide language-specific resources for default language.").Default("true").Bool()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	println(*credentials, *platform, *sheetId, *resDir, *tabName, *defLoc, *defLocPath, *preferDefLoc)
}
