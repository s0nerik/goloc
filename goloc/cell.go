package goloc

import (
	"fmt"
	"github.com/s0nerik/goloc/utils"
)

type Cell struct {
	tab    string
	row    uint
	column uint
}

func NewCell(tab string, row uint, column uint) *Cell {
	return &Cell{tab: tab, row: row, column: column}
}

func (c Cell) String() string {
	return fmt.Sprintf(`%v!%v%v`, c.tab, utils.ColumnName(c.column), c.row)
}
