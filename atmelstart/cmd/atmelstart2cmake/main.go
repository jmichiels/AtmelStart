package main

import (
	"fmt"
	"os"

	"path"

	"github.com/jmichiels/AtmelStart/atmelstart"
	"github.com/pkg/errors"
)

func main() {
	// Assume this is run in the atmel start directory.
	dirAtmelStart := "."
	if len(os.Args) > 1 {
		// The directory can also be specified as first argument.
		dirAtmelStart = os.Args[1]
	}
	pathConfig := path.Join(dirAtmelStart, "atmel_start_config.atstart")
	if _, err := os.Stat(pathConfig); os.IsNotExist(err) {
		// Could not find config file: not an atmel start directory.
		usage()
	}
	pathMakefile := path.Join(dirAtmelStart, "gcc/Makefile")
	pathToolchain := path.Join(dirAtmelStart, "toolchain.cmake")

	makefile, err := os.Open(pathMakefile)
	if err != nil {
		if os.IsNotExist(err) {
			fatal(errors.New("makefile missing: make sure to include it when you export the Atmel Start project"))
		}
		fatal(errors.Wrap(err, "open makefile"))
	}
	defer makefile.Close()
	data := atmelstart.Data{}
	if err := data.ReadMakefile(makefile); err != nil {
		fatal(errors.Wrap(err, "read makefile"))
	}

	// Create toolchain.cmake
	toolchainFile, err := os.Create(pathToolchain)
	if err != nil {
		fatal(errors.Wrap(err, "open toolchain"))
	}
	defer toolchainFile.Close()
	if err := data.WriteToolchain(toolchainFile); err != nil {
		fatal(errors.Wrap(err, "write toolchain"))
	}
	fmt.Println("generated 'toolchain.cmake'")
}

func usage() {
	fmt.Println("usage: atmelstart2cmake [atmel_start_directory]")
	os.Exit(0)
}

func fatal(err error) {
	fmt.Println("error:", err.Error())
	os.Exit(1)
}
