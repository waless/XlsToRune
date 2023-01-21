package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	err := mainImpl()
	if err != nil {
		fmt.Println(err)
	}
}

func mainImpl() error {
	var setting = ParseArgs()

	if setting.input == "" {
		return fmt.Errorf("入力ファイル指定がありません")
	}

	setting.Print()

	book, err := ParseXls(setting.input)
	if err != nil {
		return err
	}

	json_data, err := json.MarshalIndent(book, "", "  ")
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

	return nil
}
