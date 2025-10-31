package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type UI struct {
	reader *bufio.Reader
}

func NewUI() *UI {
	return &UI{reader: bufio.NewReader(os.Stdin)}
}

func (ui *UI) Read() ([]byte, error) {
	data, err := ui.reader.ReadBytes('\n')
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}
	return data, nil
}
