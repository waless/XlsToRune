package main

import "fmt"

func main() {
	var setting = ParseArgs()
	setting.Print()

	_, err := ParseXls(setting.input)
	if err != nil {
		fmt.Println(err)
	}
}
