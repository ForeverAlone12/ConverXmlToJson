package main

import "testing"

func TestListXmlFileSuccess(t *testing.T) {
	var goodResult []string

	goodResult = append(goodResult, "testFile.xml")

	res := ListXmlFile("D:\\GitHub\\ConverXmlToJson")

	if len(res) != len(goodResult) {
		t.Fatal("Количество файлов не совпадает")
	} else {
		for i := 0; i < len(res); i++ {
			countEqualFile := 0
			for j := 0; j < len(goodResult); j++ {
				if res[i] != goodResult[j] {
					countEqualFile++
				}
				if countEqualFile == len(goodResult) {
					t.Fatal("Данные не найдены")
				}
			}
		}
	}
}

func TestListXmlFileFail(t *testing.T) {
	res := ListXmlFile("E:\\")

	if res == nil {
		t.Fatal("Неверный путь")
	}
}
