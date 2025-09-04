package unixsort

import (
	"bufio"
	"io"
	"sort"
	"strconv"
	"strings"
)

type Options struct {
	Column  int
	Numeric bool
	Reverse bool
	Unique  bool
}

func SortInput(input io.Reader, opts Options) []string {
	scanner := bufio.NewScanner(input)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	sorted := sortLines(lines, opts)
	if opts.Unique {
		sorted = removeDuplicates(sorted)
	}
	return sorted
}

func sortLines(lines []string, opts Options) []string {
	sort.Slice(lines, func(i, j int) bool {
		valI := getComparisonValue(lines[i], opts.Column)
		valJ := getComparisonValue(lines[j], opts.Column)
		var result bool
		if opts.Numeric {
			numI, errI := strconv.ParseFloat(valI, 64)
			numJ, errJ := strconv.ParseFloat(valJ, 64)
			if errI == nil && errJ == nil {
				result = numI < numJ
			} else if errI != nil && errJ != nil {
				result = valI < valJ
			} else {
				result = errI == nil
			}
		} else {
			result = valI < valJ
		}
		if opts.Reverse {
			return !result
		}
		return result
	})
	return lines
}

func getComparisonValue(line string, column int) string {
	if column <= 0 {
		return line
	}
	columns := strings.Split(line, "\t")
	if column-1 < len(columns) {
		return columns[column-1]
	}
	return ""
}

func removeDuplicates(lines []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}
	return result
}