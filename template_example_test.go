package pytemplate_test

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/KSpaceer/pytemplate"
)

func ExampleTemplate_Substitute() {
	tmpl, err := pytemplate.New("ERROR: expected ${expected}, got ${actual}")
	if err != nil {
		log.Fatal(err)
	}

	substituted, err := tmpl.Substitute(pytemplate.WithMapping(map[string]string{
		"expected": "nil",
		"actual":   "ENOENT",
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(substituted)
	// Output: ERROR: expected nil, got ENOENT
}

func ExampleTemplate_SafeSubstitute() {
	tmpl, err := pytemplate.New("Welcome, $user")
	if err != nil {
		log.Fatal(err)
	}

	substituted := tmpl.SafeSubstitute()
	fmt.Println(substituted)
	// Output: Welcome, $user
}

func ExampleWithMapping() {
	tmpl, err := pytemplate.New("Shows how mapping works: key is var, value is ${var}")
	if err != nil {
		log.Fatal(err)
	}

	substituted, err := tmpl.Substitute(pytemplate.WithMapping(map[string]string{"var": "var_value"}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(substituted)
	// Output: Shows how mapping works: key is var, value is var_value
}

func ExampleWithMapper() {
	tmpl, err := pytemplate.New("Shows how mapper works: reversed palyndrome is not palyndrome, but ${palyndrome}")
	if err != nil {
		log.Fatal(err)
	}

	substituted, err := tmpl.Substitute(pytemplate.WithMapper(
		pytemplate.MapperFunc(func(s string) (string, bool) {
			var sb strings.Builder
			for {
				r, size := utf8.DecodeLastRuneInString(s)
				if size == 0 {
					break
				}
				sb.WriteRune(r)
				s = s[:len(s)-size]
			}
			return sb.String(), true
		}),
	))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(substituted)
	// Output: Shows how mapper works: reversed palyndrome is not palyndrome, but emordnylap
}

func ExampleWithSafeSubstitution() {
	tmpl, err := pytemplate.New("safe substitution: ${WILL_BE_REPLACED} and ${WONT_BE_REPLACED}")
	if err != nil {
		log.Fatal(err)
	}

	mapping := map[string]string{"WILL_BE_REPLACED": "replacement"}

	_, err = tmpl.Substitute(pytemplate.WithMapping(mapping))
	if err == nil {
		log.Fatal("expected error, got nil")
	}
	fmt.Println(err)

	substituted, err := tmpl.Substitute(pytemplate.WithMapping(mapping), pytemplate.WithSafeSubstitution())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(substituted)
	// Output:
	// failed to substitute variable WONT_BE_REPLACED
	// safe substitution: replacement and ${WONT_BE_REPLACED}
}
