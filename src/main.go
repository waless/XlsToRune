package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	var setting = ParseArgs()
	setting.Print()

	book, err := ParseXls(setting.input)
	if err != nil {
		fmt.Println(err)
	}

	json, err := json.Marshal(book)
	fmt.Printf("%s", json)

	//json := RuneBookToJson(book)
	//[fmt.Println(json)
}
