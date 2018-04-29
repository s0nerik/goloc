package goloc

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func localizationsCount(localizations Localizations) map[Lang]int {
	result := map[Lang]int{}
	for _, keyLoc := range localizations {
		for lang := range keyLoc {
			result[lang]++
		}
	}
	return result
}

// WriteLocalizations writes localization files into platform-defined directories.
func WriteLocalizations(
	platform Platform,
	dir ResDir,
	localizations Localizations,
	defLocLang Lang,
	defLocPath string,
) error {
	// Make sure we can access resources dir
	if _, err := os.Stat(dir); err != nil {
		return err
	}

	writers := map[Lang]*bufio.Writer{}

	locIndices := map[Lang]int{}
	locCounts := localizationsCount(localizations)
	locStringArgs := &LocalizedStringArgs{}
	headerArgs := &HeaderArgs{}
	footerArgs := &FooterArgs{}

	// For each localization: create a writer if needed and write each localized string
	for key, keyLoc := range localizations {
		for lang, value := range keyLoc {
			if _, ok := writers[lang]; !ok { // Create a new writer and write a header to it if needed
				// Get actual resource file dir and name
				resDir, fileName, err := localizationFilePath(platform, dir, lang, defLocLang, defLocPath)
				if err != nil {
					return err
				}

				// Create all intermediate directories
				err = os.MkdirAll(resDir, os.ModePerm)
				if err != nil {
					return err
				}

				// Create actual localization file
				file, err := os.Create(filepath.Join(resDir, fileName))
				// noinspection GoDeferInLoop
				defer file.Close()
				if err != nil {
					return err
				}

				// Create a new writer for the localization file
				writer := bufio.NewWriter(file)
				writers[lang] = writer

				// Write a header
				headerArgs.Lang = lang
				_, err = writer.WriteString(platform.Header(headerArgs))
				if err != nil {
					return err
				}
			}

			writer := writers[lang]

			// Update arguments
			locStringArgs.Index = locIndices[lang]
			locStringArgs.IsLast = locIndices[lang]+1 >= locCounts[lang]
			locStringArgs.Key = key
			locStringArgs.Lang = lang
			locStringArgs.Value = value

			// Write a localized string
			localizedString := platform.LocalizedString(locStringArgs)
			_, err := writer.WriteString(localizedString)
			if err != nil {
				return err
			}
			locIndices[lang]++
		}
	}

	// For each writer: write a footer and flush everything
	for lang, writer := range writers {
		footerArgs.Lang = lang
		_, err := writer.WriteString(platform.Footer(footerArgs))
		if err != nil {
			return err
		}

		err = writer.Flush()
		if err != nil {
			return err
		}
	}

	return nil
}

func localizationFilePath(platform Platform, dir ResDir, lang Lang, defLocLang Lang, defLocPath string) (resDir string, fileName string, err error) {
	// Handle default language
	if len(defLocLang) > 0 && lang == defLocLang && len(defLocPath) > 0 {
		resDir = path.Dir(defLocPath)
		fileName = path.Base(defLocPath)
	} else {
		filePath := platform.LocalizationFilePath(lang, dir)
		if len(filePath) == 0 {
			return "", "", &emptyLocalizationFilePath{}
		}

		resDir = path.Dir(filePath)
		fileName = path.Base(filePath)
	}
	return
}

// region Errors

type emptyLocalizationFilePath struct {
}

func (e *emptyLocalizationFilePath) Error() string {
	return fmt.Sprintf("empty localization file path")
}

// endregion
