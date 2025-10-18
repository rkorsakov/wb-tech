package main

import (
	"13/cut"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	options := parseFlags()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: cut -f list [-s] [-d delim] [file]")
		os.Exit(1)
	}
	var input io.Reader
	if len(args) == 1 {
		filename := args[0]
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
	ans, err := cut.InputCut(&input, options)
	if err != nil {
		fmt.Println("Error:", err)
	}
	for _, line := range ans {
		fmt.Println(line)
	}
}

func parseFlags() *cut.Options {
	options := &cut.Options{}
	flag.StringVar(&options.Fields, "f", "", "select only these fields; also print any line that contains no delimiter character, unless the -s option is specified")
	flag.StringVar(&options.Delimiter, "d", "\t", "use delimiter instead of TAB for field delimiter")
	flag.BoolVar(&options.Separated, "s", false, "do not print lines not containing delimiters")
	flag.Parse()
	return options
}
