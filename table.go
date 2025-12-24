package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Table struct {
	Headers []string
	Rows    [][]string
	Widths  []int
}

func NewTable(headers []string) *Table {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = utf8.RuneCountInString(h)
	}

	return &Table{
		Headers: headers,
		Rows:    [][]string{},
		Widths:  widths,
	}
}

func (t *Table) AddRow(row []string) {
	t.Rows = append(t.Rows, row)
	for i, cell := range row {
		cellLen := utf8.RuneCountInString(stripAnsi(cell))
		if i < len(t.Widths) && cellLen > t.Widths[i] {
			t.Widths[i] = cellLen
		}
	}
}

func (t *Table) drawLine(left, mid, right, fill string) string {
	parts := make([]string, len(t.Widths))
	for i, w := range t.Widths {
		parts[i] = strings.Repeat(fill, w+2)
	}
	return left + strings.Join(parts, mid) + right
}

func (t *Table) drawRow(cells []string) string {
	parts := make([]string, len(t.Widths))
	for i, w := range t.Widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}

		padding := w - utf8.RuneCountInString(stripAnsi(cell))
		if padding < 0 {
			padding = 0
		}
		parts[i] = " " + cell + strings.Repeat(" ", padding) + " "
	}
	return "|" + strings.Join(parts, "|") + "|"
}

func (t *Table) Print() {
	if len(t.Rows) == 0 {
		return
	}

	// Top border
	fmt.Println(t.drawLine("┌", "┬", "┐", "─"))

	// Header row
	fmt.Println(t.drawRow(t.Headers))

	// Header separator
	fmt.Println(t.drawLine("├", "┼", "┤", "─"))

	// Data rows
	for _, row := range t.Rows {
		fmt.Println(t.drawRow(row))
	}

	// Bottom border
	fmt.Println(t.drawLine("└", "┴", "┘", "─"))
}

func stripAnsi(str string) string {
	re := regexp.MustCompile(`\033\[[0-9;]*m`)
	return re.ReplaceAllString(str, "")
}
