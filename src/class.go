package main

import (
	"os"
	"strconv"
)

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
	class_name := "Rune_" + table.Name

	class_str := "using System;\n"
	class_str += "using UnityEngine;\n"
	class_str += "using UnityEngine.AddressableAssets;\n"
	class_str += "using UnityEngine.ResourceManagement.AsyncOperations;\n"
	class_str += "using RuneImporter;\n"
	class_str += "\n"
	class_str += addRuneClassName(class_name, len(table.Values))

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

	class_str += "    }\n"

	class_str += "\n"
	class_str += "    public static AsyncOperationHandle<" + class_name + "> LoadInstanceAsync() {\n"
	class_str += "        var path = Config.ScriptableObjectDirectory + \"" + class_name + ".asset\";\n"
	class_str += "        var handle = Addressables.LoadAssetAsync<" + class_name + ">(path);\n"
	class_str += "        handle.Completed += (handle) => { instance = handle.Result; };\n"
	class_str += "\n"
	class_str += "        return handle;\n"
	class_str += "    }\n"

	class_str += "}\n"

	file_name := "Rune_" + table.Name
	return write(file_name, class_str, out_dir)
}

func write(class_name string, class_str string, out_dir string) error {
	path := out_dir + "/" + class_name + ".cs"

	err := os.MkdirAll(out_dir, os.ModePerm)
	if err != nil {
		return err
	}

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

func addRuneClassName(type_name string, value_length int) string {
	str := "public class " + type_name + " : RuneScriptableObject\n"
	str += "{\n"
	str += "    public static " + type_name + " instance { get; private set; }\n"
	str += "\n"
	str += "    [SerializeField]\n"
	str += "    public Value[] ValueList = new Value[" + strconv.Itoa(value_length) + "];\n"
	str += "\n"
	str += "    [Serializable]\n"
	str += "    public struct Value\n"
	str += "    {\n"

	return str
}

func addRuneString(type_name RuneTypeName) string {
	return "        public string " + type_name.Value + ";\n"
}

func addRuneInteger(type_name RuneTypeName) string {
	return "        public int " + type_name.Value + ";\n"
}

func addRuneFloat(type_name RuneTypeName) string {
	return "        public float " + type_name.Value + ";\n"
}
