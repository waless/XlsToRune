package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const ToolVersion string = "1.00"

type RuneTypeInfo struct {
	ToolVersion string
}

type RuneTypeJson struct {
	Info RuneTypeInfo
	Book RuneTypeBook
}

func main() {
	err := mainImpl()
	if err != nil {
		fmt.Println(err)
	}
}

func mainImpl() error {
	var setting = ParseArgs()

	if *setting.pinput == "" {
		return fmt.Errorf("入力ファイル指定がありません")
	}

	setting.Print()

	book, err := ParseXls(*setting.pinput)
	if err != nil {
		return err
	}

	output_json := RuneTypeJson{}
	output_json.Book = book
	output_json.Info.ToolVersion = ToolVersion

	json_data, err := json.MarshalIndent(output_json, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(*setting.pout), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(*setting.pout)
	if err != nil {
		return err
	}

	_, err = file.Write(json_data)
	if err != nil {
		return err
	}

	err = OutputClassString(book, *setting.poutClass)
	if err != nil {
		return err
	}

	err = OutputEnum(book, *setting.penumNS, *setting.poutEnum)
	if err != nil {
		return err
	}

	return nil
}
