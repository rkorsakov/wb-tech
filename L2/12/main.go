package main

import (
	"12/grep"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	options := parseFlags()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: grep [flags] pattern [filename]")
		os.Exit(1)
	}
	var input io.Reader
	if len(args) > 1 {
		filename := args[1]
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
	pattern := args[0]
	results, err := grep.InputGrep(options, pattern, &input)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}

func parseFlags() *grep.Options {
	options := &grep.Options{}
	flag.IntVar(&options.After, "A", 0, "N lines after match")
	flag.IntVar(&options.Before, "B", 0, "N lines before match")
	flag.IntVar(&options.Context, "C", 0, "N lines of context")
	flag.BoolVar(&options.Count, "c", false, "Count only")
	flag.BoolVar(&options.IgnoreCase, "i", false, "Ignore case")
	flag.BoolVar(&options.Invert, "v", false, "Invert match")
	flag.BoolVar(&options.Fixed, "F", false, "Fixed string")
	flag.BoolVar(&options.LineNum, "n", false, "Line numbers")
	flag.Parse()
	if options.Context > 0 {
		options.After = options.Context
		options.Before = options.Context
	}
	return options
}
