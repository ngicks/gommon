# enum

Enum-like struct generator.

```
# go run ./cmd/generate_enum/generate_enum.go --help
Usage of /tmp/go-build581005351/b001/exe/generate_enum:
  -disable-goimports
        disable applying goimports after code generation. forced to be false if outFilename is empty.
  -matcher-returns
        whether matcher returns type [T any]. (default true)
  -o string
        out filename. stdout if empty.
  -panic-on-no-match
        whether panic on non-exhaustive match (default true)
  -pkg-name string
        [required] package name of output
  -ty string
        [required] types to enumerate. comma separated.
  -type-name string
        output typename (default "Enum")
```

It behaves like Tagged union. It actually is not.

see [./example](./example/). It has go:generate.

```
//go:generate go run ../cmd/generate_enum/generate_enum.go -o enum_a.go -pkg-name example -type-name EnumA -ty int,string,*os.File -matcher-returns=true -panic-on-no-match=false -disable-goimports=false
//go:generate go run ../cmd/generate_enum/generate_enum.go -o enum_b.go -pkg-name example -type-name EnumB -ty int,string,*os.File -matcher-returns=false -panic-on-no-match=true -disable-goimports=false
```

result is ./example/enum_a.go and ./example/enum_b.go.

see an usage example in ./main.go.

```go
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
```

Idea based on https://qiita.com/sxarp/items/cd528a546d1537105b9d.
