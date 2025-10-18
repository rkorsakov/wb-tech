package cut

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Options struct {
	Fields    string
	Delimiter string
	Separated bool
}

func InputCut(input *io.Reader, options *Options) ([]string, error) {
	numericFields, err := fieldsToNumeric(options)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(*input)
	var result []string
	for scanner.Scan() {
		line := scanner.Text()
		processedLine := process(line, options, numericFields)
		result = append(result, processedLine)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func fieldsToNumeric(options *Options) ([]int, error) {
	str := options.Fields
	if str == "" {
		return []int{}, nil
	}
	fields := strings.Split(str, ",")
	var numbers []int
	for _, field := range fields {
		if strings.Contains(field, "-") {
			rangeParts := strings.Split(field, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid field range: %s", field)
			}
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid start of range: %s", rangeParts[0])
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid end of range: %s", rangeParts[1])
			}
			if start > end {
				return nil, fmt.Errorf("invalid range: start (%d) greater than end (%d)", start, end)
			}
			for i := start; i <= end; i++ {
				numbers = append(numbers, i)
			}
		} else {
			num, err := strconv.Atoi(field)
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", field)
			}
			numbers = append(numbers, num)
		}
	}

	return numbers, nil
}

func process(line string, options *Options, fields []int) string {
	parts := strings.Split(line, options.Delimiter)
	if options.Separated && len(parts) == 1 {
		return ""
	}
	if len(fields) == 0 {
		return line
	}
	var selectedFields []string
	for _, fieldNum := range fields {
		index := fieldNum - 1
		if index >= 0 && index < len(parts) {
			selectedFields = append(selectedFields, parts[index])
		}
	}
	if len(selectedFields) == 0 {
		return ""
	}
	return strings.Join(selectedFields, options.Delimiter)
}