package grep

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

type Options struct {
	After      int
	Before     int
	Context    int
	Count      bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	LineNum    bool
	Pattern    string
}

func InputGrep(options *Options, pattern string, input *io.Reader) ([]string, error) {
	var results []string
	lines, err := readLines(input)
	if err != nil {
		return nil, fmt.Errorf("error reading input file: %v", err)
	}
	var regex *regexp.Regexp
	if options.Fixed {
		pattern = regexp.QuoteMeta(pattern)
	}
	if options.IgnoreCase {
		regex = regexp.MustCompile("(?i)" + pattern)
	} else {
		regex = regexp.MustCompile(pattern)
	}

	if options.Count {
		count := 0
		for _, line := range lines {
			matched := regex.MatchString(line)
			if (matched && !options.Invert) || (!matched && options.Invert) {
				count++
			}
		}
		return []string{fmt.Sprintf("%d", count)}, nil
	}
	alreadyPrinted := make(map[int]bool)
	for i, line := range lines {
		matched := regex.MatchString(line)
		if (matched && !options.Invert) || (!matched && options.Invert) {
			start := max(0, i-options.Before)
			for j := start; j < i; j++ {
				if !alreadyPrinted[j] {
					if options.LineNum {
						results = append(results, fmt.Sprintf("%d:%s", j+1, lines[j]))
					} else {
						results = append(results, lines[j])
					}
					alreadyPrinted[j] = true
				}
			}
			if !alreadyPrinted[i] {
				if options.LineNum {
					results = append(results, fmt.Sprintf("%d:%s", i+1, line))
				} else {
					results = append(results, line)
				}
				alreadyPrinted[i] = true
			}
			end := min(len(lines)-1, i+options.After)
			for j := i + 1; j <= end; j++ {
				if !alreadyPrinted[j] {
					if options.LineNum {
						results = append(results, fmt.Sprintf("%d:%s", j+1, lines[j]))
					} else {
						results = append(results, lines[j])
					}
					alreadyPrinted[j] = true
				}
			}
		}
	}
	return results, nil
}

func readLines(input *io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(*input)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
