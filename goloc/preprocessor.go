package goloc

// PreprocessArgs encapsulates arguments for a preprocess function
type PreprocessArgs struct {
	Localizations       Localizations
	Formats             Formats
	FormatArgs          LocalizationFormatArgs
	ResDir              ResDir
	DefaultLocalization Locale
}

type Preprocessor interface {
	Preprocess(args PreprocessArgs) error
}
