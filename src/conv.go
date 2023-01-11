package main

func RuneBookToJson(book RuneTypeBook) string {
	return toBookJson(book)
}

func toBookJson(book RuneTypeBook) string {
	result := ""
	for _, sheet := range book.sheets {
		result += toSheetJson(sheet)
	}

	return result
}

func toSheetJson(sheet RuneTypeSheet) string {
	result := ""
	for _, table := range sheet.tables {
		result += toTableJson(table)
	}

	return result
}

func toTableJson(table RuneTypeTable) string {
	return ""
}