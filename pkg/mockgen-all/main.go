package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	inputDir  = flag.String("t", ".", "target directory. cwd if empty.")
	outputDir = flag.String("o", ".", "output directory of generated mock.")
)

func main() {
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}

func _main() error {
	ensureMockgen()

	flag.Parse()

	dirents, err := os.ReadDir(*inputDir)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, dirent := range dirents {
		if dirent.IsDir() ||
			!strings.HasSuffix(dirent.Name(), ".go") ||
			strings.HasSuffix(dirent.Name(), "_test.go") {
			continue
		}

		targetSrc := filepath.Join(*inputDir, dirent.Name())
		destination := filepath.Join(*outputDir, dirent.Name())
		cmd := exec.Command("mockgen", append([]string{"-source", targetSrc, "-destination", destination}, flag.Args()...)...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			return errors.WithStack(err)
		}
	}
	return nil
}

func ensureMockgen() {
	cmd := exec.Command("mockgen", "--help")
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
}
