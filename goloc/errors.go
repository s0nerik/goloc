package goloc

import "fmt"

type emptySheetError struct {
	tab string
}

type firstRowNotFoundError struct {
	cell Cell
}

type noFormatColumnError struct {
	tab                 string
	requiredColumnTitle string
}

type noPlatformColumnError struct {
	tab           string
	platformNames []string
}

type formatKeyNotSpecifiedError struct {
	cell Cell
}

type formatValueNotSpecifiedError struct {
	platformName string
	cell         Cell
}

type formatValueInvalidError struct {
	cell         Cell
	platformName string
	formatValue  string
	reason       error
}

type formatArgsDifferentError struct {
	cell Cell
	key  Key
	lang string
}

type wrongValueTypeError struct {
	cell Cell
}

type wrongKeyTypeError struct {
	cell Cell
}

type columnNotFoundError struct {
	cell   Cell
	column string
}

type langColumnsNotFoundError struct {
	cell Cell
}

type localizationMissingError struct {
	cell Cell
	key  Key
	lang string
}

type keyMissingError struct {
	cell Cell
}

type formatNotFoundError struct {
	cell       Cell
	formatName string
}

func newFormatArgsDifferentError(tab string, row int, col int, key Key, lang string) *formatArgsDifferentError {
	return &formatArgsDifferentError{
		cell: *NewCell(tab, uint(row), uint(col)),
		key:  key,
		lang: lang,
	}
}

func newLocalizationMissingError(tab string, row int, col int, key Key, lang string) *localizationMissingError {
	return &localizationMissingError{
		cell: *NewCell(tab, uint(row), uint(col)),
		key:  key,
		lang: lang,
	}
}

func newKeyMissingError(tab string, row int, col int) *keyMissingError {
	return &keyMissingError{
		cell: *NewCell(tab, uint(row), uint(col)),
	}
}

func (e *emptySheetError) Error() string {
	return fmt.Sprintf(`%v!A1: sheet is empty`, e.tab)
}

func (e *firstRowNotFoundError) Error() string {
	return fmt.Sprintf(`%v: there's no first row in the tab`, e.cell)
}

func (e *noFormatColumnError) Error() string {
	return fmt.Sprintf(`%v!A1: "%v" column is missing in the first row`, e.tab, e.requiredColumnTitle)
}

func (e *noPlatformColumnError) Error() string {
	return fmt.Sprintf(`%v!A1: can't find any of %v columns in the first row`, e.tab, e.platformNames)
}

func (e *formatKeyNotSpecifiedError) Error() string {
	return fmt.Sprintf(`%v: format name is not specified`, e.cell)
}

func (e *formatValueNotSpecifiedError) Error() string {
	return fmt.Sprintf(`%v: value for "%v" platform is not specified`, e.cell, e.platformName)
}

func (e *formatValueInvalidError) Error() string {
	return fmt.Sprintf(`%v: format "%v" is invalid for platform "%v" (%v)`, e.cell, e.formatValue, e.platformName, e.reason)
}

func (e *formatArgsDifferentError) Error() string {
	return fmt.Sprintf(`%v: format arguments must be the same for each language`, e.cell)
}

func (e *wrongValueTypeError) Error() string {
	return fmt.Sprintf(`%v: wrong value type`, e.cell)
}

func (e *wrongKeyTypeError) Error() string {
	return fmt.Sprintf(`%v: wrong key type`, e.cell)
}

func (e *columnNotFoundError) Error() string {
	return fmt.Sprintf(`%v: "%v" column not found in the first row`, e.cell, e.column)
}

func (e *langColumnsNotFoundError) Error() string {
	return fmt.Sprintf(`%v: language columns are not found`, e.cell)
}

func (e *localizationMissingError) Error() string {
	return fmt.Sprintf(`%v: "%v" is missing for "%v" language`, e.cell, e.key, e.lang)
}

func (e *keyMissingError) Error() string {
	return fmt.Sprintf(`%v: key name is missing, ignoring this string...`, e.cell)
}

func (e *formatNotFoundError) Error() string {
	return fmt.Sprintf(`%v: no such format - "%v"`, e.cell, e.formatName)
}
