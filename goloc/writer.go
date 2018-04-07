package goloc

import (
	"os"
	"bufio"
	"errors"
	"fmt"
	"path/filepath"
	"path"
)

func WriteLocalizations(
	platform Platform,
	dir ResDir,
	localizations Localizations,
	defLocLang Lang,
	defLocPath string,
) error {
	// Make sure the the resources dir exists
	file, err := os.Open(dir)
	if err != nil {
		return err
	}

	writers := map[Lang]*bufio.Writer{}

	// For each localization: create a writer if needed and write each localized string
	for key, keyLoc := range localizations {
		for lang, value := range keyLoc {
			if writer, ok := writers[lang]; !ok { // Create a new writer and write a header to it if needed
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
				file, err = os.Create(filepath.Join(resDir, fileName))
				// noinspection GoDeferInLoop
				defer file.Close()
				if err != nil {
					return err
				}

				// Open a new writer for the localization file
				writer := bufio.NewWriter(file)
				if writer == nil {
					return errors.New(fmt.Sprintf(`can't create a bufio.Writer for %v`, file))
				}

				writers[lang] = writer

				// Write a header
				_, err = writer.WriteString(platform.Header(lang))
				if err != nil {
					return errors.New(fmt.Sprintf(`can't write header to %v, reason: %v`, file, err))
				}
			} else { // Use an existing writer to write another localized string
				localizedString := platform.Localization(lang, key, value)
				if len(localizedString) < 1 {
					return errors.New(fmt.Sprintf(`can't write a new line to %v, reason: %v`, file, err))
				}

				_, err := writer.WriteString(localizedString)
				if err != nil {
					return errors.New(fmt.Sprintf(`can't write a new line to %v, reason: %v`, file, err))
				}
			}
		}
	}

	// For each writer: write a footer and flush everything
	for lang, writer := range writers {
		_, err := writer.WriteString(platform.Footer(lang))
		if err != nil {
			return errors.New(fmt.Sprintf(`can't write footer to %v, reason: %v`, file, err))
		}

		err = writer.Flush()
		if err != nil {
			return errors.New(fmt.Sprintf(`can't write to %v, reason: %v`, file, err))
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
		fileName = path.Base(fileName)
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
