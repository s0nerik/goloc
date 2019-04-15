package goloc

// PostprocessArgs encapsulates arguments for a postprocess function
type PostprocessArgs struct {
	Localizations Localizations
	ResDir        ResDir
}

type Postprocessor interface {
	Postprocess(args PostprocessArgs) error
}
