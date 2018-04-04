package utils

func ColumnName(i uint) string {
	return columnName("", i)
}

func columnName(orig string, i uint) string {
	if i < 26 {
		asciiIndex := 65 + i
		return orig + string(asciiIndex)
	} else {
		return columnName(columnName(orig, i/26-1), i%26)
	}
}