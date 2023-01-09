package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ERuneType int

const (
	EType ERuneType = iota
	EEnum
	EInt
	EFloat
	EString
)

const (
	SRuneType = "RuneType"
	SType     = "type"
	SEnum     = "enum"
	SString   = "string"
	SInt      = "int"
	SFloat    = "float"
	SComment  = "#"
)

func (c *ERuneType) ToString() string {
	switch *c {
	case EType:
		return SType
	case EEnum:
		return SEnum
	case EString:
		return SString
	case EInt:
		return SInt
	}

	return ""
}

type RuneTypeName struct {
	kind  ERuneType
	value string
}

func (c *RuneTypeName) Print() {
	fmt.Println("--- RuneTypeName ---")
	fmt.Printf("kind  : %s\n", c.kind.ToString())
	fmt.Printf("value : %s\n", c.value)
}

type RuneTypeValue struct {
	typeName   RuneTypeName
	valueArray []string
	colIndex   int
}

func (c *RuneTypeValue) Print() {
	fmt.Println("--- RuneTypeValue ---")
	fmt.Printf("col_index   : %d\n", c.colIndex)
	fmt.Printf("value_array : %s\n", c.valueArray)
	c.typeName.Print()
}

type RuneTypeTable struct {
	name   string
	values []RuneTypeValue
}

func (c *RuneTypeTable) Print() {
	fmt.Println("--- RuneTypeTable ---")
	fmt.Printf("name : %s\n", c.name)
	for _, v := range c.values {
		v.Print()
	}
}

func (c *RuneTypeTable) FindTypeValueFromColIndex(col_index int) int {
	for i, v := range c.values {
		if v.colIndex == col_index {
			return i
		}
	}

	return -1
}

type RuneTypeSheet struct {
	name   string
	tables []RuneTypeTable
}

func (c *RuneTypeSheet) Print() {
	fmt.Println("--- RuneTypeSheet ---")
	fmt.Printf("name : %s\n", c.name)
	for _, v := range c.tables {
		v.Print()
	}
}

type RuneTypeBook struct {
	name   string
	sheets []RuneTypeSheet
}

func (c *RuneTypeBook) Print() {
	fmt.Println("--- RuneTypeBook ---")
	fmt.Printf("name : %s\n", c.name)
	for _, v := range c.sheets {
		v.Print()
	}
}

type contextType int

const (
	ContextNone contextType = iota
	ContextTable
)

type context struct {
	runeBook     RuneTypeBook
	currentType  contextType
	currentSheet RuneTypeSheet
	currentTable RuneTypeTable
	rowIndex     int
	colIndex     int
}

var gctx context

func ParseXls(path string) (RuneTypeBook, error) {
	file, err := excelize.OpenFile(path)
	if err != nil {
		return gctx.runeBook, err
	}

	name := strings.Split(filepath.Base(path), ".")[0]
	gctx.runeBook.name = name

	sheets := file.GetSheetList()
	for i := 0; i < len(sheets); i++ {
		name := sheets[i]

		rows, err := file.Rows(name)
		if err != nil {
			return gctx.runeBook, err
		}

		gctx.rowIndex = 0

		if len(gctx.currentSheet.tables) > 0 {
			gctx.runeBook.sheets = append(gctx.runeBook.sheets, gctx.currentSheet)
		}
		gctx.currentSheet = RuneTypeSheet{}
		gctx.currentSheet.name = name
		for rows.Next() {
			cols, err := rows.Columns()
			if err != nil {
				return gctx.runeBook, err
			}

			err = parseCols(cols)
			if err != nil {
				return gctx.runeBook, err
			}
			gctx.rowIndex++
		}
	}

	return gctx.runeBook, nil
}

func parseCols(cols []string) error {
	gctx.colIndex = 0

	switch gctx.currentType {
	case ContextNone:
		return parseColsForNone(cols)
	case ContextTable:
		return parseColsForTable(cols)
	}

	return nil
}

func parseColsForNone(cols []string) error {
	if len(cols) <= 0 {
		return nil
	}

	col := cols[0]
	if col == SRuneType {
		err := newCurrentTable(cols)
		if err != nil {
			return err
		}
		gctx.currentType = ContextTable
	}

	return nil
}

func parseColsForTable(cols []string) error {
	current_table := gctx.currentTable
	for i, col := range cols {
		index := current_table.FindTypeValueFromColIndex(i)
		if index < 0 {
			continue
		}

		type_value := current_table.values[index]
		type_value.valueArray = append(type_value.valueArray, col)
	}

	return nil
}

func newCurrentTable(cols []string) error {
	if len(gctx.currentTable.values) > 0 {
		gctx.currentSheet.tables = append(gctx.currentSheet.tables, gctx.currentTable)
	}
	gctx.currentTable = RuneTypeTable{}
	gctx.colIndex = 0
	for _, col := range cols {
		value := strings.Trim(col, " ")

		if isComment(value) {
		} else if strings.Contains(value, SEnum) {
			v := parseSEnum(col)
			gctx.currentTable.values = append(gctx.currentTable.values, v)
		} else {
			err := checkTypeValidity(value)
			if err != nil {
				return err
			}

			if strings.Contains(value, SType) {
				parseSType(col)
			} else if strings.Contains(value, SString) {
				parseSString(col)
			} else if strings.Contains(value, SInt) {
				parseSInt(col)
			} else if strings.Contains(value, SFloat) {
				parseSFloat(col)
			}
		}

		gctx.colIndex++
	}

	return nil
}

func checkTypeValidity(str string) error {
	strs := strings.Split(str, ":")
	str_len := len(strs)

	validity := true
	if str_len == 1 {
		str := strs[0]
		if strings.Contains(str, SRuneType) {
			validity = true
		} else {
			validity = false
		}
	} else if str_len == 2 {
		validity = true
	} else {
		validity = false
	}

	if validity {
		return nil
	} else {
		return fmt.Errorf("%s : invalid type string : %s", makeRuneErrorPrefix(), str)
	}
}

func makeRuneErrorPrefix() string {
	return fmt.Sprintf("row, %d, col %d", gctx.rowIndex, gctx.colIndex)
}

func parseSType(str string) error {
	str_array := parseTypeString(str)
	if len(str_array) > 1 {
		name := str_array[0]
		if str_array[1] != SType {
			return fmt.Errorf("row %d, col %d, not found RuneType", gctx.rowIndex, gctx.colIndex)
		}
		gctx.currentTable.name = name
	}

	return nil
}

func parseSEnum(str string) RuneTypeValue {
	result := RuneTypeValue{}
	result.typeName.kind = EType

	return result
}

func parseSString(str string) RuneTypeValue {
	return parseTypeValue(str, EString)
}

func parseSInt(str string) RuneTypeValue {
	return parseTypeValue(str, EInt)
}

func parseSFloat(str string) RuneTypeValue {
	return parseTypeValue(str, EFloat)
}

func isComment(str string) bool {
	return strings.Contains(str, SComment)
}

func parseTypeValue(str string, t ERuneType) RuneTypeValue {
	result := RuneTypeValue{}
	result.typeName.kind = t

	strs := parseTypeString(str)
	result.typeName.value = strs[0]
	result.colIndex = gctx.colIndex

	return result
}

func parseTypeString(str string) []string {
	str = strings.Trim(str, " ")
	return strings.Split(str, ":")
}
