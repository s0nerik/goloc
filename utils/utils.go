package utils

func ColumnName(i uint) string {
	return columnName("", i)
}

func columnName(orig string, i uint) string {
	const LettersNum = 26
	const AsciiFirstLetterIndex = 65
	if i < LettersNum {
		asciiIndex := AsciiFirstLetterIndex + i
		return orig + string(asciiIndex)
	} else {
		return columnName(columnName(orig, i/LettersNum-1), i%LettersNum)
	}
}
