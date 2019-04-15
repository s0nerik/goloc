package goloc

// PreprocessArgs encapsulates arguments for a preprocess function
type PreprocessArgs struct {
	Localizations Localizations
	ResDir        ResDir
}

type Preprocessor interface {
	Preprocess(args PreprocessArgs) error
}