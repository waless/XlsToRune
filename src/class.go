package main

import "os"

func OutputClassString(book RuneTypeBook, out_dir string) error {
	for _, sheet := range book.Sheets {
		err := outputSheet(sheet, out_dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func outputSheet(sheet RuneTypeSheet, out_dir string) error {
	for _, table := range sheet.Tables {
		err := outputTable(table, out_dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func outputTable(table RuneTypeTable, out_dir string) error {
	class_str := "using UnityEngine;\n"
	class_str += "\n"
	class_str += addRuneClassName(table.Name)

	for _, t := range table.Types {
		switch t.TypeName.Kind {
		case SString:
			class_str += addRuneString(t.TypeName)
		case SInt:
			class_str += addRuneInteger(t.TypeName)
		case SFloat:
			class_str += addRuneFloat(t.TypeName)
		}
	}

	class_str += "}\n"
	return write(table.Name, class_str, out_dir)
}

func write(class_name string, class_str string, out_dir string) error {
	path := out_dir + "/" + class_name + ".cs"

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(class_str))
	if err != nil {
		return err
	}

	return nil
}

func addRuneClassName(type_name string) string {
	str := "public class " + type_name + " : RuneScriptableObject\n"
	str += "{\n"

	return str
}

func addRuneString(type_name RuneTypeName) string {
	return "    public string " + type_name.Value + ";\n"
}

func addRuneInteger(type_name RuneTypeName) string {
	return "    public int " + type_name.Value + ";\n"
}

func addRuneFloat(type_name RuneTypeName) string {
	return "    public float " + type_name.Value + ";\n"
}
