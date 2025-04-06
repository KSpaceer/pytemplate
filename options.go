package pytemplate

type substituteOptions struct {
	mapping map[string]string
	mapper  Mapper
	safe    bool
}

// SubstituteOption allows to tune Template.Substitute and Template.SafeSubstitute behavior.
type SubstituteOption func(*substituteOptions)

// WithMapping sets map of mappings between variables and their replacing values.
// WithMapping has higher priority than WithMapper, meaning provided mapping will be used first,
// and only when there is no match, Mapper.Map will be called.
func WithMapping(mapping map[string]string) SubstituteOption {
	return func(so *substituteOptions) {
		so.mapping = mapping
	}
}

// Mapper implements custom mapping logic when substituting variables.
type Mapper interface {

	// Map maps given variable to actual string value.
	// The second boolean parameter shows if mapping was successful.
	Map(s string) (string, bool)
}

// MapperFunc implements Mapper interface.
type MapperFunc func(string) (string, bool)

func (f MapperFunc) Map(s string) (string, bool) {
	return f(s)
}

// WithMapper sets Mapper to apply custom mapping logic to variables.
// WithMapper has lower priority than WithMapping, meaning Mapper.Map will be called
// only if there is no match in WithMapping mapping.
func WithMapper(mapper Mapper) SubstituteOption {
	return func(so *substituteOptions) {
		so.mapper = mapper
	}
}

// WithSafeSubstitution enables error-safe mode.
// When variable cannot be substituted, it is original text is used.
func WithSafeSubstitution() SubstituteOption {
	return func(so *substituteOptions) {
		so.safe = true
	}
}
