package goloc

import (
	"fmt"

	"github.com/s0nerik/goloc/utils"
)

// Cell represents a table cell.
type Cell struct {
	tab    string
	row    uint
	column uint
}

// NewCell creates and returns a new Cell instance given tab name, row and column indices.
func NewCell(tab string, row uint, column uint) *Cell {
	return &Cell{tab: tab, row: row, column: column}
}

func (c Cell) String() string {
	return fmt.Sprintf(`%v!%v%v`, c.tab, utils.ColumnName(c.column), c.row)
}
