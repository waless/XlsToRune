package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

// ---定数----
//
// ファイル拡張子
const RuneExt = ".rune"

type SettingParam struct {
	input string
	pout  *string
}

func ParseArgs() SettingParam {
	var result SettingParam

	// 入力ファイルパス(必須)
	// 先頭に引数名なしで書かれている事を意図している
	var in = ""
	if len(os.Args) > 1 {
		in = os.Args[1]
	}
	result.input = in

	// 補助引数(なくても良い)
	var out_default = makeOutputDefaultPath(in)
	result.pout = flag.String("o", out_default, "出力ファイルパス")

	flag.Parse()

	return result
}

func (c *SettingParam) Print() {
	fmt.Println("----Setting----")
	fmt.Printf("input  : %s\n", c.input)
	fmt.Printf("output : %s\n", *c.pout)
}

func makeOutputDefaultPath(input_path string) string {
	var in_ext = path.Ext(input_path)
	var none_ext = strings.Replace(input_path, in_ext, "", 1)
	return none_ext + RuneExt
}
