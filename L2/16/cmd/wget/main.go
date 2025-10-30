package main

import (
	"16/internal/downloader"
	"flag"
	"fmt"
	"os"
)

func main() {
	rFlag := flag.Int("r", 0, "recursion depth")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] URL\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	URL := flag.Args()[0]

	fmt.Printf("Downloading %s with depth %d...\n", URL, *rFlag)

	err := downloader.DownloadPage(URL, *rFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading page: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Download completed successfully!")
}
