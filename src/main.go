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
	setting.Print()

	book, err := ParseXls(setting.input)
	if err != nil {
		return err
	}
	//book.Print()

	json_data, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(*setting.pout)
	if err != nil {
		return err
	}

	n, err := file.Write(json_data)
	if err != nil {
		return err
	}
	//fmt.Println(string(json_data))
	fmt.Println(n)

	return nil
}
