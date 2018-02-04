package main

import "testing"

func TestListXmlFileSuccess(t *testing.T) {
	var goodResult []string

	goodResult = append(goodResult, "testFile.xml")

	res := ListXmlFile("D:\\GitHub\\ConverXmlToJson")

	if res != goodResult {
		t.Fatal("Error!!!!! данные не совпадают")
	}
}

func TestListXmlFileFail(t *testing.T) {
	res := ListXmlFile("E:\\")

	if res == nil {
		t.Fatal("must return error")
	}
}
