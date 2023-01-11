package main

import "fmt"

func main() {
	var setting = ParseArgs()
	setting.Print()

	book, err := ParseXls(setting.input)
	if err != nil {
		fmt.Println(err)
	}

	json := RuneBookToJson(book)
	fmt.Println(json)
}
