package pytemplate

type substituteOptions struct {
	mapping map[string]string
	mapper  Mapper
	safe    bool
}

type SubstituteOption func(*substituteOptions)

func WithMapping(mapping map[string]string) SubstituteOption {
	return func(so *substituteOptions) {
		so.mapping = mapping
	}
}

type Mapper interface {
	Map(s string) (string, bool)
}

type MapperFunc func(string) (string, bool)

func (f MapperFunc) Map(s string) (string, bool) {
	return f(s)
}

func WithMapper(mapper Mapper) SubstituteOption {
	return func(so *substituteOptions) {
		so.mapper = mapper
	}
}

func WithSafeSubstitution() SubstituteOption {
	return func(so *substituteOptions) {
		so.safe = true
	}
}
