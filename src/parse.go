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
	case EFloat:
		return SFloat
	}

	return ""
}

type RuneTypeName struct {
	Kind  ERuneType
	Value string
}

func (c *RuneTypeName) Print() {
	fmt.Println("--- RuneTypeName ---")
	fmt.Printf("kind  : %s\n", c.Kind.ToString())
	fmt.Printf("value : %s\n", c.Value)
}

type RuneTypeValue struct {
	TypeName   RuneTypeName
	ValueArray []string
	ColIndex   int
}

func (c *RuneTypeValue) Print() {
	fmt.Println("--- RuneTypeValue ---")
	fmt.Printf("col_index   : %d\n", c.ColIndex)
	fmt.Printf("value_num   : %d\n", len(c.ValueArray))
	fmt.Printf("value_array : %s\n", c.ValueArray)
	c.TypeName.Print()
}

type RuneTypeTable struct {
	name   string
	values []RuneTypeValue
}

func (c *RuneTypeTable) Print() {
	fmt.Println("--- RuneTypeTable ---")
	fmt.Printf("name        : %s\n", c.name)
	fmt.Printf("value count : %d\n", len(c.values))
	for _, v := range c.values {
		v.Print()
	}
}

func (c *RuneTypeTable) FindTypeValueFromColIndex(col_index int) int {
	for i, v := range c.values {
		if v.ColIndex == col_index {
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
	Name   string
	Sheets []RuneTypeSheet
}

func (c *RuneTypeBook) Print() {
	fmt.Println("--- RuneTypeBook ---")
	fmt.Printf("name : %s\n", c.Name)
	for _, v := range c.Sheets {
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
	gctx.runeBook = RuneTypeBook{}
	gctx.runeBook.Name = name

	sheets := file.GetSheetList()
	for i := 0; i < len(sheets); i++ {
		name := sheets[i]

		rows, err := file.Rows(name)
		if err != nil {
			return gctx.runeBook, err
		}

		gctx.rowIndex = 0

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

		if len(gctx.currentSheet.tables) > 0 {
			gctx.runeBook.Sheets = append(gctx.runeBook.Sheets, gctx.currentSheet)
		}
	}

	return gctx.runeBook, nil
}

func parseCols(cols []string) error {
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

	gctx.colIndex = 0
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
	current_table := &gctx.currentTable
	col_len := len(cols)
	value_len := len(current_table.values)
	if value_len <= 0 || col_len <= 0 {
		return nil
	}

	gctx.colIndex = 0

	begin_index := current_table.values[0].ColIndex
	end_index := current_table.values[value_len-1].ColIndex

	if end_index >= col_len {
		end_index = col_len - 1
	}

	for i := begin_index; i <= end_index; i++ {
		col := cols[i]

		index := current_table.FindTypeValueFromColIndex(i)
		if index < 0 {
			continue
		}

		type_value := &current_table.values[index]
		type_value.ValueArray = append(type_value.ValueArray, col)

		gctx.colIndex++
	}

	value_count := end_index - begin_index
	blank_count := value_len - value_count - 1
	blank_begin_index := end_index + 1
	blank_end_index := blank_begin_index + blank_count - 1

	for i := blank_begin_index; i <= blank_end_index; i++ {
		index := current_table.FindTypeValueFromColIndex(i)
		type_value := &current_table.values[index]
		type_value.ValueArray = append(type_value.ValueArray, " ")
	}

	return nil
}

func newCurrentTable(cols []string) error {
	gctx.currentTable = RuneTypeTable{}
	gctx.colIndex = 0

	current_table := &gctx.currentTable
	for _, col := range cols {
		value := strings.Trim(col, " ")

		if isComment(value) {
		} else if strings.Contains(value, SEnum) {
			v := parseSEnum(col)
			current_table.values = append(current_table.values, v)
		} else {
			err := checkTypeValidity(value)
			if err != nil {
				return err
			}

			if strings.Contains(value, SType) {
				err = parseSType(col)
			} else if strings.Contains(value, SString) {
				v := parseSString(col)
				current_table.values = append(current_table.values, v)
			} else if strings.Contains(value, SInt) {
				v := parseSInt(col)
				current_table.values = append(current_table.values, v)
			} else if strings.Contains(value, SFloat) {
				v := parseSFloat(col)
				current_table.values = append(current_table.values, v)
			}

			if err != nil {
				return err
			}
		}

		gctx.colIndex++
	}

	if len(gctx.currentTable.values) > 0 {
		gctx.currentSheet.tables = append(gctx.currentSheet.tables, gctx.currentTable)
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
	result.TypeName.Kind = EEnum
	result.ColIndex = gctx.colIndex

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
	result.TypeName.Kind = t

	strs := parseTypeString(str)
	result.TypeName.Value = strs[0]
	result.ColIndex = gctx.colIndex

	return result
}

func parseTypeString(str string) []string {
	str = strings.Trim(str, " ")
	strs := strings.Split(str, ":")
	for i, s := range strs {
		strs[i] = strings.Trim(s, " ")
	}

	return strs
}
