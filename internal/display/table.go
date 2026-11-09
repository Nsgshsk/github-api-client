package display

import (
	"io"
	"strings"
	"unicode/utf8"
)

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignRight
)

type Column struct {
	Header string
	Align  Alignment
}

type Table struct {
	Columns []Column
	Rows    [][]string
}

func NewTable(columns ...Column) *Table {
	return &Table{Columns: columns}
}

func (t *Table) AddRow(cells ...string) {
	t.Rows = append(t.Rows, cells)
}

func (t *Table) Render(w io.Writer) error {
	if len(t.Columns) == 0 {
		return nil
	}

	widths := make([]int, len(t.Columns))
	for i, column := range t.Columns {
		widths[i] = utf8.RuneCountInString(normalizeCell(column.Header))
	}

	rows := make([][]string, len(t.Rows))
	for i, row := range t.Rows {
		normalized := make([]string, len(t.Columns))
		for j := range t.Columns {
			if j < len(row) {
				normalized[j] = normalizeCell(row[j])
			}
			if width := utf8.RuneCountInString(normalized[j]); width > widths[j] {
				widths[j] = width
			}
		}
		rows[i] = normalized
	}

	var builder strings.Builder
	writeBorder(&builder, "┌", "┬", "┐", widths)
	writeRow(&builder, headerCells(t.Columns), t.Columns, widths)
	writeBorder(&builder, "├", "┼", "┤", widths)
	for _, row := range rows {
		writeRow(&builder, row, t.Columns, widths)
	}
	writeBorder(&builder, "└", "┴", "┘", widths)

	_, err := io.WriteString(w, builder.String())
	return err
}

func headerCells(columns []Column) []string {
	cells := make([]string, len(columns))
	for i, column := range columns {
		cells[i] = column.Header
	}
	return cells
}

func writeBorder(builder *strings.Builder, left, middle, right string, widths []int) {
	builder.WriteString(left)
	for i, width := range widths {
		builder.WriteString(strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			builder.WriteString(middle)
		}
	}
	builder.WriteString(right)
	builder.WriteString("\n")
}

func writeRow(builder *strings.Builder, cells []string, columns []Column, widths []int) {
	builder.WriteString("│")
	for i, column := range columns {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}
		padding := widths[i] - utf8.RuneCountInString(cell)
		builder.WriteString(" ")
		if column.Align == AlignRight {
			builder.WriteString(strings.Repeat(" ", padding))
			builder.WriteString(cell)
		} else {
			builder.WriteString(cell)
			builder.WriteString(strings.Repeat(" ", padding))
		}
		builder.WriteString(" │")
	}
	builder.WriteString("\n")
}

func normalizeCell(cell string) string {
	return strings.Join(strings.Fields(cell), " ")
}
