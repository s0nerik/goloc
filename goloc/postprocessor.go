package goloc

// PostprocessArgs encapsulates arguments for a postprocess function
type PostprocessArgs struct {
	Localizations       Localizations
	Formats             Formats
	FormatArgs          LocalizationFormatArgs
	ResDir              ResDir
	DefaultLocalization Locale
}

type Postprocessor interface {
	Postprocess(args PostprocessArgs) error
}
