package pytemplate

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const delimiter = "$"

var (
	pattern = regexp.MustCompile(`\$(?:(?P<escaped>\$)|(?P<named>[_a-zA-z][_a-zA-Z0-9]*)|` +
		`\{(?P<braced>[_a-zA-Z][_a-zA-Z0-9]*)\}|(?P<invalid>))`)

	escapedIdx = pattern.SubexpIndex("escaped")
	namedIdx   = pattern.SubexpIndex("named")
	bracedIdx  = pattern.SubexpIndex("braced")
	invalidIdx = pattern.SubexpIndex("invalid")
)

type SubstitutionFailedError struct {
	Variable string
}

func (e *SubstitutionFailedError) Error() string {
	return "failed to substitute variable " + e.Variable
}

func (e *SubstitutionFailedError) Is(err error) bool {
	sfe, ok := err.(*SubstitutionFailedError)
	if !ok {
		return false
	}
	return e.Variable == sfe.Variable
}

type Template struct {
	parts []templatePart
}

func New(template string) (*Template, error) {
	parts, err := parseTemplate(template)
	if err != nil {
		return nil, err
	}
	return &Template{parts: parts}, nil
}

func parseTemplate(template string) ([]templatePart, error) {
	var (
		prevIdx int
		parts   []templatePart
	)

	for _, submatches := range pattern.FindAllStringSubmatchIndex(template, -1) {
		var (
			replacedName string
			braced       bool
		)

		switch {
		case submatches[namedIdx*2] != -1:
			replacedName = template[submatches[2*namedIdx]:submatches[2*namedIdx+1]]
		case submatches[bracedIdx*2] != -1:
			replacedName = template[submatches[2*bracedIdx]:submatches[2*bracedIdx+1]]
			braced = true
		case submatches[escapedIdx*2] != -1:
			parts = append(parts, templatePart{
				value: template[prevIdx:submatches[0]] + delimiter,
			})
			prevIdx = submatches[1]
			continue
		case submatches[invalidIdx*2] != -1:
			return nil, fmt.Errorf("invalid placeholder in template string at index %d", submatches[invalidIdx*2])
		}

		constPrefix := template[prevIdx:submatches[0]]
		if constPrefix != "" {
			parts = append(parts, templatePart{
				value: constPrefix,
			})
		}
		parts = append(parts, templatePart{
			value:      replacedName,
			isVariable: true,
			braced:     braced,
		})

		prevIdx = submatches[1]
	}

	if constSuffix := template[prevIdx:]; constSuffix != "" {
		parts = append(parts, templatePart{value: constSuffix})
	}
	return parts, nil
}

func (t *Template) Substitute(opts ...SubstituteOption) (string, error) {
	var so substituteOptions

	for _, opt := range opts {
		opt(&so)
	}

	var sb strings.Builder
	for i := range t.parts {
		if !t.parts[i].isVariable {
			sb.WriteString(t.parts[i].value)
			continue
		}

		mapped, ok := substituteVariable(&so, t.parts[i].value)

		if !ok {
			if so.safe {
				sb.WriteString(t.parts[i].String())
			} else {
				return "", &SubstitutionFailedError{Variable: t.parts[i].value}
			}
		} else {
			sb.WriteString(mapped)
		}
	}

	return sb.String(), nil
}

func (t *Template) SafeSubstitute(opts ...SubstituteOption) string {
	result, _ := t.Substitute(append(opts, WithSafeSubstitution())...)
	return result
}

func substituteVariable(so *substituteOptions, name string) (string, bool) {
	var (
		mapped string
		ok     bool
	)

	if so.mapping != nil {
		mapped, ok = so.mapping[name]
	}

	if !ok && so.mapper != nil {
		mapped, ok = so.mapper.Map(name)
	}

	return mapped, ok
}

func (t Template) MarshalText() (text []byte, err error) {
	var bb bytes.Buffer
	for i := range t.parts {
		bb.WriteString(t.parts[i].String())
	}
	return bb.Bytes(), nil
}

func (t *Template) UnmarshalText(text []byte) error {
	parts, err := parseTemplate(string(text))
	if err != nil {
		return err
	}
	t.parts = parts
	return nil
}

type templatePart struct {
	value      string
	isVariable bool
	braced     bool
}

func (tp templatePart) String() string {
	switch {
	case !tp.isVariable:
		return strings.ReplaceAll(tp.value, delimiter, delimiter+delimiter)
	case tp.braced:
		return delimiter + "{" + tp.value + "}"
	default:
		return delimiter + tp.value
	}
}
