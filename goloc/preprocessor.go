package goloc

// PreprocessArgs encapsulates arguments for a preprocess function
type PreprocessArgs struct {
	Localizations       Localizations
	FormatArgs          LocalizationFormatArgs
	ResDir              ResDir
	DefaultLocalization Lang
}

type Preprocessor interface {
	Preprocess(args PreprocessArgs) error
}
