package main

import (
	"flag"
	"fmt"
	"path"
	"strings"
)

// ---定数----
//
// ファイル拡張子
const RuneExt = ".rune"

type SettingParam struct {
	pinput    *string
	pout      *string
	poutClass *string
	poutEnum  *string
	penumNS   *string
}

func ParseArgs() SettingParam {
	var result SettingParam

	result.pinput = flag.String("i", "", "入力ファイルパス")

	// 補助引数(なくても良い)
	result.pout = flag.String("o", "", "出力ファイルパス")

	result.poutClass = flag.String("class", "", "クラス出力ディレクトリ")
	result.poutEnum = flag.String("enum", "", "enum出力ディレクトリ")
	result.penumNS = flag.String("enum-namespace", "", "enumのnamespace")

	flag.Parse()

	out_default := makeOutputDefaultPath(*result.pinput)
	out_dir := path.Dir(out_default)

	if len(*result.pout) <= 0 {
		*result.pout = out_default
	}
	if len(*result.poutClass) <= 0 {
		*result.poutClass = out_dir
	}
	if len(*result.poutEnum) <= 0 {
		*result.poutEnum = out_dir
	}

	return result
}

func (c *SettingParam) Print() {
	fmt.Println("----Setting----")
	fmt.Printf("input          : %s\n", *c.pinput)
	fmt.Printf("output         : %s\n", *c.pout)
	fmt.Printf("class          : %s\n", *c.poutClass)
	fmt.Printf("enum           : %s\n", *c.poutEnum)
	fmt.Printf("enum-namespace : %s\n", *c.penumNS)
}

func makeOutputDefaultPath(input_path string) string {
	var in_ext = path.Ext(input_path)
	var none_ext = strings.Replace(input_path, in_ext, "", 1)
	return none_ext + RuneExt
}
