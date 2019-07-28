package goloc

type Source interface {
	FormatsDocumentName() string
	LocalizationsDocumentName() string

	Formats() ([][]RawCell, error)
	Localizations() ([][]RawCell, error)
}