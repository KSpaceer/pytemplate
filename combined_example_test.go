package pytemplate_test

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/KSpaceer/pytemplate"
)

type hostinfo struct {
	user     string
	time     time.Time
	cpu      int
	ramBytes int
}

func getInfo() hostinfo {
	return hostinfo{
		user:     "default user",
		time:     time.Now(),
		cpu:      33,
		ramBytes: 2_030_163_000,
	}
}

func Example_mappingAndMapperCombined() {
	tmpl, err := pytemplate.New("Good morning, ${user}! Current time is $TIME, CPU percentage is ${CPU}% and RAM usage is $RAM")
	if err != nil {
		log.Fatal(err)
	}

	overloads := map[string]string{
		"TIME": "today =)",
		"RAM":  "100%",
	}

	info := getInfo()

	mapper := pytemplate.MapperFunc(func(s string) (string, bool) {
		switch s {
		case "user":
			return info.user, true
		case "TIME":
			return info.time.String(), true
		case "CPU":
			return strconv.Itoa(info.cpu), true
		case "RAM":
			return strconv.FormatFloat(float64(info.ramBytes)/8_000_000_000, 'f', 2, 64) + "%", true
		default:
			return "", false
		}
	})

	substituted := tmpl.SafeSubstitute(
		pytemplate.WithMapping(overloads),
		pytemplate.WithMapper(mapper),
	)
	fmt.Println(substituted)
	// Output: Good morning, default user! Current time is today =), CPU percentage is 33% and RAM usage is 100%
}
