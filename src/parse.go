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
	TypeName RuneTypeName
	colIndex int
}

func (c *RuneTypeValue) Print() {
	fmt.Println("--- RuneTypeValue ---")
	fmt.Printf("col_index   : %d\n", c.colIndex)
	c.TypeName.Print()
}

type RuneTypeTable struct {
	Name        string
	Types       []RuneTypeValue
	Values      [][]string
	typeIndex   int
	ignoreIndex []int
}

func (c *RuneTypeTable) Print() {
	fmt.Println("--- RuneTypeTable ---")
	fmt.Printf("name        : %s\n", c.Name)
	fmt.Printf("type  count : %d\n", len(c.Types))
	fmt.Printf("value count : %d\n", len(c.Values))
	for _, t := range c.Types {
		t.Print()
	}
}

func (c *RuneTypeTable) IsIgnoreIndex(index int) bool {
	for _, v := range c.ignoreIndex {
		if v == index {
			return true
		}
	}

	return false
}

type RuneTypeSheet struct {
	Name   string
	Tables []RuneTypeTable
}

func (c *RuneTypeSheet) Print() {
	fmt.Println("--- RuneTypeSheet ---")
	fmt.Printf("name : %s\n", c.Name)
	for _, v := range c.Tables {
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
	runeBook      RuneTypeBook
	currentType   contextType
	currentSheet  RuneTypeSheet
	pcurrentTable *RuneTypeTable
	rowIndex      int
	colIndex      int
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
		gctx.currentSheet.Name = name
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

		if len(gctx.currentSheet.Tables) > 0 {
			gctx.runeBook.Sheets = append(gctx.runeBook.Sheets, gctx.currentSheet)
		}
	}

	return gctx.runeBook, nil
}

func parseCols(cols []string) error {
	col_len := len(cols)

	switch gctx.currentType {
	case ContextNone:
		return parseColsForNone(cols)
	case ContextTable:
		if col_len > 0 {
			if cols[0] == SRuneType {
				return parseColsForNone(cols)
			} else {
				return parseColsForTable(cols)
			}
		}
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
	result := []string{}

	table_len := len(gctx.currentSheet.Tables)
	if table_len <= 0 {
		return nil
	}

	current_table := &gctx.currentSheet.Tables[table_len-1]
	col_len := len(cols)
	type_len := len(current_table.Types)
	if type_len <= 0 || col_len <= 0 {
		return nil
	}

	gctx.colIndex = 0

	begin_index := current_table.typeIndex + 1
	end_index := begin_index + type_len
	for i := begin_index; i <= end_index; i++ {
		if current_table.IsIgnoreIndex(i) {
			continue
		}

		if i < col_len {
			col := cols[i]
			result = append(result, col)
		} else {
			result = append(result, " ")
		}

		gctx.colIndex++
	}
	current_table.Values = append(current_table.Values, result)

	return nil
}

func newCurrentTable(cols []string) error {
	table := RuneTypeTable{}
	gctx.colIndex = 0
	gctx.pcurrentTable = &table

	for i, col := range cols {
		value := strings.Trim(col, " ")

		if isComment(value) {
			table.ignoreIndex = append(table.ignoreIndex, i)
		} else if strings.Contains(value, SEnum) {
			v := parseSEnum(col)
			table.Types = append(table.Types, v)
		} else {
			err := checkTypeValidity(value)
			if err != nil {
				return err
			}

			if strings.Contains(value, SType) {
				err = parseSType(col)
				table.typeIndex = i
			} else if strings.Contains(value, SString) {
				v := parseSString(col)
				table.Types = append(table.Types, v)
			} else if strings.Contains(value, SInt) {
				v := parseSInt(col)
				table.Types = append(table.Types, v)
			} else if strings.Contains(value, SFloat) {
				v := parseSFloat(col)
				table.Types = append(table.Types, v)
			}

			if err != nil {
				return err
			}
		}

		gctx.colIndex++
	}

	if len(table.Types) > 0 {
		gctx.currentSheet.Tables = append(gctx.currentSheet.Tables, table)
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
		gctx.pcurrentTable.Name = name
	}

	return nil
}

func parseSEnum(str string) RuneTypeValue {
	result := RuneTypeValue{}
	result.TypeName.Kind = EEnum
	result.colIndex = gctx.colIndex

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
	result.colIndex = gctx.colIndex

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
