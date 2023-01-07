package main

func main() {
	var setting = ParseArgs()
	setting.Print()

	ParseXls(setting.input)
}
