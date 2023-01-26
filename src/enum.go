package main

import "os"

func OutputEnum(book RuneTypeBook, enum_ns string, out_dir string) error {
	for _, sheet := range book.Sheets {
		for _, table := range sheet.Tables {
			outputEnumFromTable(table, enum_ns, out_dir)
		}
	}

	return nil
}

func outputEnumFromTable(table RuneTypeTable, enum_ns string, out_dir string) error {
	enum_index := -1
	for i, t := range table.Types {
		if t.TypeName.Kind == SEnum {
			enum_index = i
			break
		}
	}

	// 出力するenumはなかったので何もしない
	if enum_index < 0 || len(table.Values) <= 0 {
		return nil
	}

	enum_name := "e" + table.Name
	path := out_dir + "/" + enum_name + ".cs"

	enum_str := ""

	if len(enum_ns) > 0 {
		enum_str += "namespace " + enum_ns + "\n"
		enum_str += "{\n\n"
	}

	enum_str += "public enum " + enum_name + "\n"
	enum_str += "{\n"

	for _, v := range table.Values {
		element_name := v.Values[enum_index]
		enum_str += "    " + element_name + ",\n"
	}

	enum_str += "\n"
	enum_str += "    Count,\n"

	enum_str += "}\n"

	if len(enum_ns) > 0 {
		enum_str += "\n"
		enum_str += "}\n"
	}

	err := os.MkdirAll(out_dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(enum_str))
	if err != nil {
		return err
	}

	return nil
}
