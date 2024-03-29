package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
)

const autoGenerationNotice = "// Code generated by github.com/ngicks/gommon/pkg/atomicstate/cmd/generate_state_impl/generate_state_impl.go. DO NOT EDIT."

const importPackages = `import "sync/atomic"`

var stateTemplate = template.Must(template.New("v").Parse(`
// {{.StateName}}State is atomic state primitive.
// It holds a boolean state corresponds to its name.
type {{.StateName}}State struct {
	s uint32
}

// Is{{.StateName}} is atomic state checker.
// It returns true if state is set, and vice versa.
func (s *{{.StateName}}State) Is{{.StateName}}() bool {
	return atomic.LoadUint32(&s.s) == 1
}

// Set{{.StateName}} is atomic state setter.
// It tries to set its state based on to.
// If first element of to is false, it tries set it to false,
// true otherwise.
//
// swapped is true when it is successfully set, false if it is already set to the state.
func (s *{{.StateName}}State) Set{{.StateName}}(to ...bool) (swapped bool) {
	setTo := true
	if len(to) > 0 {
		setTo = to[0]
	}

	if setTo {
		return atomic.CompareAndSwapUint32(&s.s, 0, 1)
	} else {
		return atomic.CompareAndSwapUint32(&s.s, 1, 0)
	}
}

// New{{.StateName}}State builds splitted {{.StateName}}State wrapper.
// Either or both can be embedded and/or used as unexported member to hide its setter.
func New{{.StateName}}State() (checker *{{.StateName}}StateChecker, setter *{{.StateName}}StateSetter) {
	s := new({{.StateName}}State)
	checker = &{{.StateName}}StateChecker{s}
	setter = &{{.StateName}}StateSetter{s}
	return
}

// {{.StateName}}StateSetter is simple wrapper of {{.StateName}}State.
// It only exposes Is{{.StateName}}.
type {{.StateName}}StateChecker struct {
	s *{{.StateName}}State
}

// Is{{.StateName}} is atomic state checker.
// It returns true if state is set, and vice versa.
func (s *{{.StateName}}StateChecker) Is{{.StateName}}() bool {
	return s.s.Is{{.StateName}}()
}

// {{.StateName}}StateSetter is simple wrapper of {{.StateName}}State.
// It only exposes Set{{.StateName}}. 
type {{.StateName}}StateSetter struct {
	s *{{.StateName}}State
}

// Set{{.StateName}} is atomic state setter.
// It tries to set its state based on to.
// If first element of to is false, it tries set it to false,
// true otherwise.
//
// swapped is true when it is successfully set, false if it is already set to the state.
func (s *{{.StateName}}StateSetter) Set{{.StateName}}(to ...bool) (swapped bool) {
	return s.s.Set{{.StateName}}(to...)
}`))

var (
	pkgName     = flag.String("pkg-name", "", "package name of output")
	stateName   = flag.String("state", "", "output state name list. comma-separated.")
	outFilename = flag.String("o", "", "out filename. stdout if empty.")
)

func main() {
	if err := _main(); err != nil {
		panic(err)
	}
}

func _main() (err error) {
	flag.Parse()

	if *pkgName == "" {
		return fmt.Errorf("empty pkgName")
	}

	stateNames := strings.Split(*stateName, ",")

	var outFile io.Writer
	if *outFilename == "" {
		outFile = os.Stdout
	} else {
		file, err := os.Create(*outFilename)
		if err != nil {
			return err
		}
		defer file.Close()
		outFile = file
	}

	_, err = fmt.Fprintf(outFile, "%s\n", autoGenerationNotice)
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(outFile, "package %s\n\n", *pkgName)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(outFile, "%s\n\n", importPackages)
	if err != nil {
		return
	}
	for _, stateName := range stateNames {
		err = stateTemplate.Execute(outFile, struct{ StateName string }{capitalize(stateName)})
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(outFile, "\n")
		if err != nil {
			return
		}
	}

	return nil
}

func capitalize(str string) string {
	head := str[:1]
	rest := str[1:]

	return strings.ToUpper(head) + strings.ToLower(rest)
}
