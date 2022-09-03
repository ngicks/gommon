package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

var (
	inputDir  = flag.String("t", ".", "target directory. cwd if empty.")
	outputDir = flag.String("o", ".", "output directory of generated mock.")
)

func main() {
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
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
		var destination string
		if *inputDir == *outputDir {
			destination = filepath.Join(*outputDir, "__mock", dirent.Name())
		} else {
			destination = filepath.Join(*outputDir, dirent.Name())
		}

		if !isDestMockgen(destination) {
			fmt.Fprintf(os.Stderr, "ignoring: dest exists but not mockgen generated file. dest = %s\n", destination)
		}

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

const mockgenWarning = "// Code generated by MockGen. DO NOT EDIT."

func isDestMockgen(dest string) bool {
	f, err := os.Open(dest)
	if err != nil {
		if !errors.Is(err, syscall.ENOENT) {
			panic(err)
		}
		return true
	}

	scanner := bufio.NewScanner(f)

	if scanner.Scan() && scanner.Text() == mockgenWarning {
		return true
	}
	return false
}
