package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ERuneType int

const (
	Type ERuneType = iota
	Enum
	Int
	Float
	String
)

type RuneTypeName struct {
	kind  ERuneType
	value string
}

type RuneTypeValue struct {
	typeName   RuneTypeName
	valueArray []string
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

func ParseXls(path string) RuneTypeBook {
	var result RuneTypeBook

	file, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return result
	}

	sheets := file.GetSheetList()
	for i := 0; i < len(sheets); i++ {
		name := sheets[i]

		rows, err := file.Rows(name)
		if err != nil {
			fmt.Println(err)
			return result
		}

		for rows.Next() {
			cols, err := rows.Columns()
			if err != nil {
				fmt.Println(err)
				return result
			}

			for _, col := range cols {
				fmt.Printf("row: %v\n", col)
			}
		}
	}
}
