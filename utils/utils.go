package utils

// ColumnName returns a column name in Google Sheets format.
func ColumnName(i uint) string {
	return columnName("", i)
}

func columnName(orig string, i uint) string {
	const LettersNum = 26
	const ASCIIFirstLetterIndex = 65
	if i < LettersNum {
		asciiIndex := ASCIIFirstLetterIndex + i
		return orig + string(asciiIndex)
	}
	return columnName(columnName(orig, i/LettersNum-1), i%LettersNum)
}
