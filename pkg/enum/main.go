package main

import (
	"fmt"

	"github.com/ngicks/gommon/pkg/enum/example"
)

func matcher(enumA example.EnumA[string]) {
	result := enumA.Match(example.EnumAMatcher[string]{
		Int: func(v int) string {
			return fmt.Sprint(v)
		},
		String: func(v string) string {
			return v
		},
		Default: func(v any) string {
			return "default"
		},
	})

	fmt.Println(result)
}

func main() {
	matcher(example.EnumAInt[string](123))
	matcher(example.EnumAString[string]("foo"))
	matcher(example.EnumAOsFile[string](nil))
	// 123
	// foo
	// default
}
