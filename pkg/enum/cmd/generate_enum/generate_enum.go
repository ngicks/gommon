package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

var (
	pkgName          = flag.String("pkg-name", "", "[required] package name of output")
	typeName         = flag.String("type-name", "Enum", "output typename")
	ty               = flag.String("ty", "", "[required] types to enumerate. comma seperated.")
	matcherReturns   = flag.Bool("matcher-returns", true, "whether matcher returns type [T any].")
	outFilename      = flag.String("o", "", "out filename. stdout if empty.")
	disableGoimports = flag.Bool("disable-goimports", false, "disable applying goimports after code generation. forced to be false if outFilename is empty.")
	panicOnNoMatch   = flag.Bool("panic-on-no-match", true, "whether panic on non-exhaustive match")
)

const autoGenerationNotice = "// Code generated by github.com/ngicks/gommon/pkg/enum/cmd/generate_enum/generate_enum.go. DO NOT EDIT."

var enumFuncMap = template.FuncMap{
	"capitalize": capitalize,
	"makeVariants": func(str string) string {
		return joinDot(removeStar(str))
	},
}

var enumTemplate = template.Must(template.New("v").Funcs(enumFuncMap).Parse(`
type {{.TypeName}}{{if .MatcherReturns}}[T any]{{end}} struct {
	data any
}

{{range $index, $element := .Types}}
func {{$.TypeName}}{{makeVariants $element}}{{if $.MatcherReturns}}[T any]{{end}}(val {{$element}})  {{$.TypeName}}{{if $.MatcherReturns}}[T]{{end}} {
	return {{$.TypeName}}{{if $.MatcherReturns}}[T]{{end}}{
		data: val,
	}
}
{{end}}

type {{.TypeName}}Matcher{{if .MatcherReturns}}[T any]{{end}} struct {
{{- range $index, $typName := .Types}}
    {{makeVariants $typName}} func({{$typName}}) {{if $.MatcherReturns}}T{{end}}
{{- end}} 
	Any func() {{if $.MatcherReturns}}T{{end}}
}

func (e {{.TypeName}}{{if .MatcherReturns}}[T]{{end}}) Match(m {{.TypeName}}Matcher{{if .MatcherReturns}}[T]{{end}}) {{if .MatcherReturns}}T{{end}} {
	{{- if .MatcherReturns}}var ret T{{end}}
	switch x := e.data.(type) {
{{- range $index, $typName := .Types}}
	case {{$typName}}:
		if m.{{makeVariants $typName}} != nil {
			{{if $.MatcherReturns}}ret = {{end}}m.{{makeVariants $typName}}(x) 
			return {{if $.MatcherReturns}} ret {{end}}
		}
{{- end}}
	}

	if m.Any != nil {
		{{if .MatcherReturns}}ret = {{end}}m.Any()
		return {{if .MatcherReturns}} ret {{end}}
	}

	{{- if .PanicOnNoMatch}}

	panic("non exhaustive match")
	{{- else}}
		{{- if .MatcherReturns}}

			var zeroValue T
			return zeroValue
		{{- end}}
	{{- end}}
}
`))

func main() {
	if err := _main(); err != nil {
		panic(fmt.Sprintf("%+v\n", err))
	}
}

func _main() (err error) {
	flag.Parse()

	if *pkgName == "" {
		return fmt.Errorf("empty pkgName")
	}
	if *ty == "" {
		return fmt.Errorf("empty ty")
	}

	types := strings.Split(*ty, ",")
	if len(types) == 0 || every(types, func(s string) bool { return len(s) == 0 }) {
		return fmt.Errorf("malformed ty = %+v", types)
	}

	outFile, err := parseOut(*outFilename)
	if err != nil {
		return errors.WithStack(err)
	}
	if closer, ok := outFile.(io.Closer); ok {
		defer closer.Close()
	}

	_, err = fmt.Fprintf(outFile, "%s\n", autoGenerationNotice)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = fmt.Fprintf(outFile, "package %s\n\n", *pkgName)
	if err != nil {
		return errors.WithStack(err)
	}

	variable := struct {
		TypeName       string
		MatcherReturns bool
		Types          []string
		PanicOnNoMatch bool
	}{
		TypeName:       capitalize(*typeName),
		MatcherReturns: *matcherReturns,
		Types:          types,
		PanicOnNoMatch: *panicOnNoMatch,
	}

	fmt.Fprintf(os.Stderr, "%+v\n", variable)

	err = enumTemplate.
		Funcs(enumFuncMap).
		Execute(outFile, variable)
	if err != nil {
		return errors.WithStack(err)
	}

	if *outFilename != "" && !*disableGoimports {
		out, err := exec.Command("goimports", "-w", *outFilename).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %+v, out=%s\n", err, string(out))
			return errors.WithStack(err)
		}
	}

	return nil
}

func every(strSl []string, predicate func(string) bool) bool {
	for _, v := range strSl {
		if !predicate(v) {
			return false
		}
	}
	return true
}

func parseOut(outFilename string) (io.Writer, error) {
	if outFilename == "" {
		return os.Stdout, nil
	} else {
		file, err := os.Create(outFilename)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

func capitalize(str string) string {
	head := str[:1]
	rest := str[1:]

	return strings.ToUpper(head) + rest
}

func joinDot(str string) string {
	out := []string{}
	for _, v := range strings.Split(str, ".") {
		out = append(out, capitalize(v))
	}
	return strings.Join(out, "")
}

func removeStar(str string) string {
	return strings.ReplaceAll(str, "*", "")
}
