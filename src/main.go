package main

import "fmt"

func main() {
	var setting = ParseArgs()
	setting.Print()

	book, err := ParseXls(setting.input)
	if err != nil {
		fmt.Println(err)
	}

	book.Print()
}
