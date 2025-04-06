package pytemplate_test

import (
	"errors"
	"testing"

	"github.com/KSpaceer/pytemplate"
)

func TestTemplate_Invalid(t *testing.T) {
	_, err := pytemplate.New("invalid template $")
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
}

func TestTemplate_Substitute(t *testing.T) {
	tests := []struct {
		name     string
		template string

		mapping map[string]string
		mapper  pytemplate.Mapper

		expected    string
		expectedErr error
	}{
		{
			name:     "simple",
			template: "Hello, $friend! How are you?",
			mapping:  map[string]string{"friend": "World"},
			expected: "Hello, World! How are you?",
		},
		{
			name:     "simple braced",
			template: "Hello, ${friend}! How are you?",
			mapping:  map[string]string{"friend": "World"},
			expected: "Hello, World! How are you?",
		},
		{
			name:     "simple with mapper",
			template: "Hello, ${friend}! How are you?",
			mapper: pytemplate.MapperFunc(func(s string) (string, bool) {
				return s + "-san", true
			}),
			expected: "Hello, friend-san! How are you?",
		},
		{
			name:     "multiple substitutions",
			template: "Hello, ${friend}! Tell me about $first, $second and ${third}",
			mapping: map[string]string{
				"friend": "Bob",
				"first":  "history",
				"second": "math",
				"third":  "biology",
			},
			expected: "Hello, Bob! Tell me about history, math and biology",
		},
		{
			name: "multiple substitutions: mapping+mapper combined",
			template: "Measurements: length - ${length}, weight - ${weight}, density - ${density}, " +
				"temperature - ${temperature}",
			mapping: map[string]string{
				"length":      "5m",
				"temperature": "450K",
			},
			mapper: pytemplate.MapperFunc(func(string) (string, bool) {
				return "unknown", true
			}),
			expected: "Measurements: length - 5m, weight - unknown, density - unknown, temperature - 450K",
		},
		{
			name:     "missing mapping",
			template: "My favorite ${thing} is ${MISSING_NO}",
			mapping: map[string]string{
				"thing": "pokemon",
			},
			expectedErr: &pytemplate.SubstitutionFailedError{Variable: "MISSING_NO"},
		},
		{
			name:     "missing mapper",
			template: "My favorite ${thing} is ${MISSING_NO}",
			mapping: map[string]string{
				"thing": "pokemon",
			},
			mapper: pytemplate.MapperFunc(func(string) (string, bool) {
				return "", false
			}),
			expectedErr: &pytemplate.SubstitutionFailedError{Variable: "MISSING_NO"},
		},
		{
			name:     "escaped",
			template: "I really-really like $parameter and $$money$$!!!",
			mapping:  map[string]string{"parameter": "wealth"},
			expected: "I really-really like wealth and $money$!!!",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpl, err := pytemplate.New(tc.template)
			if err != nil {
				t.Fatalf("unexpected error in New(): %s", err)
			}

			result, err := tmpl.Substitute(
				pytemplate.WithMapping(tc.mapping),
				pytemplate.WithMapper(tc.mapper),
			)

			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("failed to get expected error %s, got %s", tc.expectedErr, err)
			}

			if tc.expectedErr == nil && result != tc.expected {
				t.Fatalf("result does not match expected value %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestTemplate_SafeSubstitute(t *testing.T) {
	tests := []struct {
		name     string
		template string

		mapping map[string]string
		mapper  pytemplate.Mapper

		expected string
	}{
		{
			name:     "simple",
			template: "Hello, $friend! How are you?",
			mapping:  map[string]string{"friend": "World"},
			expected: "Hello, World! How are you?",
		},
		{
			name:     "simple braced",
			template: "Hello, ${friend}! How are you?",
			mapping:  map[string]string{"friend": "World"},
			expected: "Hello, World! How are you?",
		},
		{
			name:     "simple with mapper",
			template: "Hello, ${friend}! How are you?",
			mapper: pytemplate.MapperFunc(func(s string) (string, bool) {
				return s + "-san", true
			}),
			expected: "Hello, friend-san! How are you?",
		},
		{
			name:     "multiple substitutions",
			template: "Hello, ${friend}! Tell me about $first, $second and ${third}",
			mapping: map[string]string{
				"friend": "Bob",
				"first":  "history",
				"second": "math",
				"third":  "biology",
			},
			expected: "Hello, Bob! Tell me about history, math and biology",
		},
		{
			name: "multiple substitutions: mapping+mapper combined",
			template: "Measurements: length - ${length}, weight - ${weight}, density - ${density}, " +
				"temperature - ${temperature}",
			mapping: map[string]string{
				"length":      "5m",
				"temperature": "450K",
			},
			mapper: pytemplate.MapperFunc(func(string) (string, bool) {
				return "unknown", true
			}),
			expected: "Measurements: length - 5m, weight - unknown, density - unknown, temperature - 450K",
		},
		{
			name:     "missing mapping",
			template: "My favorite ${thing} is ${MISSING_NO}",
			mapping: map[string]string{
				"thing": "pokemon",
			},
			expected: "My favorite pokemon is ${MISSING_NO}",
		},
		{
			name:     "missing mapper",
			template: "My favorite ${thing} is ${MISSING_NO}",
			mapping: map[string]string{
				"thing": "pokemon",
			},
			mapper: pytemplate.MapperFunc(func(string) (string, bool) {
				return "", false
			}),
			expected: "My favorite pokemon is ${MISSING_NO}",
		},
		{
			name:     "escaped",
			template: "I really-really like $parameter and $$money$$!!!",
			mapping:  map[string]string{"parameter": "wealth"},
			expected: "I really-really like wealth and $money$!!!",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpl, err := pytemplate.New(tc.template)
			if err != nil {
				t.Fatalf("unexpected error in New(): %s", err)
			}

			result := tmpl.SafeSubstitute(
				pytemplate.WithMapping(tc.mapping),
				pytemplate.WithMapper(tc.mapper),
			)

			if result != tc.expected {
				t.Fatalf("result does not match expected value %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestTemplate_TextUnmarshalMarshal(t *testing.T) {
	tcases := []struct {
		name     string
		template string
	}{
		{
			name:     "full featured",
			template: "We have $named, ${braced}, $$escaped",
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			var tmpl pytemplate.Template

			if err := tmpl.UnmarshalText([]byte(tc.template)); err != nil {
				t.Fatalf("unexpected error %s during unmarshalling", err)
			}

			substituted := tmpl.SafeSubstitute(pytemplate.WithMapper(pytemplate.MapperFunc(func(string) (string, bool) {
				return "replaced", true
			})))

			t.Logf("substituted from unmarshalled template: %s", substituted)

			marshalled, err := tmpl.MarshalText()
			if err != nil {
				t.Fatalf("unexpected error %s during marshalling", err)
			}

			if string(marshalled) != tc.template {
				t.Fatalf("round-tripped template %q is not equal to the original one %q", marshalled, tc.template)
			}
		})
	}
}
