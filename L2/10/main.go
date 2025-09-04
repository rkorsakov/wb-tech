package main

import (
	"10/unixsort"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	KFlag int
	NFlag bool
	RFlag bool
	UFlag bool
)

func main() {
	parseCombinedFlags()

	var input io.Reader
	var filename string

	if len(flag.Args()) > 0 {
		filename = flag.Args()[0]
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	} else {
		input = os.Stdin
	}

	opts := unixsort.Options{
		Column:  KFlag,
		Numeric: NFlag,
		Reverse: RFlag,
		Unique:  UFlag,
	}

	sorted := unixsort.SortInput(input, opts)
	for _, line := range sorted {
		fmt.Println(line)
	}
}

func parseCombinedFlags() {
	args := os.Args[1:]
	var filteredArgs []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && len(arg) > 2 {
			for _, char := range arg[1:] {
				filteredArgs = append(filteredArgs, "-"+string(char))
			}
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	os.Args = []string{os.Args[0]}
	os.Args = append(os.Args, filteredArgs...)

	flag.IntVar(&KFlag, "k", 0, "Column to sort")
	flag.BoolVar(&NFlag, "n", false, "Sort by numeric values")
	flag.BoolVar(&RFlag, "r", false, "Reverse")
	flag.BoolVar(&UFlag, "u", false, "Only unique")
	flag.Parse()
}
