package goloc

type Source interface {
	Formats() ([][]string, error)
	Localizations() ([][]string, error)
}