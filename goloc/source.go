package goloc

type Source interface {
	Name() string
	Formats() ([][]string, error)
	Localizations() ([][]string, error)
}