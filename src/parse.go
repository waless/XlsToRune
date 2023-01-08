package main

import (
	"fmt"
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

type RuneTypeName struct {
	kind  ERuneType
	value string
}

type RuneTypeValue struct {
	typeName   RuneTypeName
	valueArray []string
	colIndex   int
}

type RuneTypeTable struct {
	name   string
	values []RuneTypeValue
}

type RuneTypeSheet struct {
	name   string
	tables []RuneTypeTable
}

type RuneTypeBook struct {
	name   string
	sheets []RuneTypeSheet
}

type contextType int

const (
	ContextNone contextType = iota
	ContextTable
)

type context struct {
	runeBook     RuneTypeBook
	currentType  contextType
	currentTable RuneTypeTable
	rowIndex     int
	colIndex     int
}

var gctx context

func ParseXls(path string) RuneTypeBook {
	file, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return gctx.runeBook
	}

	sheets := file.GetSheetList()
	for i := 0; i < len(sheets); i++ {
		name := sheets[i]

		rows, err := file.Rows(name)
		if err != nil {
			fmt.Println(err)
			return gctx.runeBook
		}

		gctx.rowIndex = 0
		for rows.Next() {
			cols, err := rows.Columns()
			if err != nil {
				fmt.Println(err)
				return gctx.runeBook
			}

			parseCols(cols)
			gctx.rowIndex++
		}
	}

	return gctx.runeBook
}

func parseCols(cols []string) {
	gctx.colIndex = 0

	switch gctx.currentType {
	case ContextNone:
		parseColsForNone(cols)
	case ContextTable:
		parseColsForTable(cols)
	}
}

func parseColsForNone(cols []string) {
	for _, col := range cols {
		switch col {
		case SRuneType:
			newCurrentTable(cols)
			return
		}
		gctx.colIndex++
	}
}

func parseColsForTable(cols []string) {
}

func newCurrentTable(cols []string) error {
	gctx.currentTable = RuneTypeTable{}
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
	}

	return nil
}

func checkTypeValidity(str string) error {
	strs := strings.Split(str, ":")
	if len(strs) == 2 {
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

	return result
}

func parseTypeString(str string) []string {
	str = strings.Trim(str, " ")
	return strings.Split(str, ":")
}
