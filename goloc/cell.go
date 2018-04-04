package goloc

import (
	"fmt"
	"github.com/s0nerik/goloc/utils"
)

type cell struct {
	tab    string
	row    uint
	column uint
}

func (c cell) String() string {
	return fmt.Sprintf(`%v!%v%v`, c.tab, utils.ColumnName(c.column), c.row)
}