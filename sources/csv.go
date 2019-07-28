package sources

import (
	"encoding/csv"
	"github.com/s0nerik/goloc/goloc"
	"os"
)

type csvSource struct {
	localizationsFilePath string
	formatsFilePath       string
}

func CSV(localizationsFilePath string, formatsFilePath string) *csvSource {
	return &csvSource{
		localizationsFilePath: localizationsFilePath,
		formatsFilePath:       formatsFilePath,
	}
}

func readCsv(filePath string) (result [][]string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func (s csvSource) FormatsDocumentName() string {
	return s.formatsFilePath
}

func (s csvSource) LocalizationsDocumentName() string {
	return s.localizationsFilePath
}

func (s csvSource) Formats() ([][]goloc.RawCell, error) {
	return readCsv(s.formatsFilePath)
}

func (s csvSource) Localizations() ([][]goloc.RawCell, error) {
	return readCsv(s.localizationsFilePath)
}
