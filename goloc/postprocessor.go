package goloc

// PostprocessArgs encapsulates arguments for a postprocess function
type PostprocessArgs struct {
	Localizations       Localizations
	Formats             Formats
	FormatArgs          LocalizationFormatArgs
	ResDir              ResDir
	DefaultLocalization Lang
}

type Postprocessor interface {
	Postprocess(args PostprocessArgs) error
}
