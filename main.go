package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	ini "github.com/ochinchina/go-ini"
	ini2 "gopkg.in/ini.v1"
)

var fileErrorName = "Error.log"

type Student struct {
	XMLName xml.Name `xml:"student"`
	Name    string   `xml:"name,attr"`
	Mark    int      `xml:"mark,attr"`
}

type Course struct {
	XMLName  xml.Name  `xml:"course"`
	Name     string    `xml:"name,attr"`
	Students []Student `xml:"student"`
}

type Courses struct {
	XMLName xml.Name `xml:"courses"`
	Cours   []Course `xml:"course"`
}

func main() {

	// Считать путь к папке с файлами xml
	var pathToDirectory string
	fmt.Print("Введите путь к папке с xml-файлами: ")
	fmt.Scanf("%s", &pathToDirectory)

	// Получить список xml-файлов
	listFiles := ListXmlFile(pathToDirectory)

	if len(listFiles) == 0 { // при отсутствии файлов xml  в папке вывести ошибку
		Logs("В папке \"" + pathToDirectory + "\" не найдены файлы xml")
	} else {

		iniFiles := CreateListFile("ini", listFiles)

		jsonFiles := CreateListFile("json", iniFiles)

		fmt.Print("Успешно созданы файлы: ")
		fmt.Println(jsonFiles)
	}
}

/*
 * Создание списка файлов дл конвертации
 * @param typeFile - тип создаваемого списка файлов
 * @param listFiles - список файлов для конвертации
 * @return - список преобразованных файлов
 */
func CreateListFile(typeFile string, listFiles []string) []string {

	var f func(fileName string) string

	switch typeFile {
	case "ini":
		f = ConvertXmlToIni
	case "json":
		f = ConvertIniToJson
	default:
		{
			Logs("")
			os.Exit(1)
		}
	}

	var files []string
	for _, file := range listFiles {
		fileName := f(file)
		if fileName != "" {
			files = append(files, fileName)
		}
	}

	if len(files) == 0 {
		Logs("")
		os.Exit(1)
	}

	return files
}

/*
 * Получение списка xml-файлов в папке
 * @param Directory - путь к папке для проверки
 * @return список имен файлов, если файлы есть, иначе - nil
 */
func ListXmlFile(Directory string) (listFiles []string) {

	// считывание файлов из папки
	files, err := ioutil.ReadDir(Directory)
	// если путь не является папкой, или нет такого пути
	if err != nil {
		// записать ошибку в файл
		Logs("Указанный путь \"" + Directory + "\" не является директорией или его не существует!")
		return // выход из функции
	}

	// исследование всех файлов в папке
	for _, file := range files {
		// считать имя файла
		fileName := file.Name()
		// если файл содержит расширение xml
		if strings.HasSuffix(fileName, ".xml") {
			// то добавить в список найденных файлов
			listFiles = append(listFiles, fileName)
		}
	}
	return listFiles
}

/*
 *  Перенесение информации из xml в ini
 * @param xmlFileName - имя xml-файла
 * @return - имя ini-файла или пустую строку при ошибки считывания
 */
func ConvertXmlToIni(xmlFileName string) string {

	// открытие xml-файла
	xmlFile, err := os.Open(xmlFileName)
	if err != nil {
		Logs("Ошибка считывания " + xmlFileName)
		return ""
	}
	defer xmlFile.Close()

	// считать структуру файла
	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var c Courses
	// записать данные из файла в структуру
	xml.Unmarshal(XMLdata, &c)

	// создание ini файла
	ini := ini.NewIni()

	// сохранение в файл
	for _, course := range c.Cours {
		for _, student := range course.Students {
			section := ini.NewSection(student.Name)
			section.Add(course.Name, strconv.Itoa(student.Mark))
		}
	}

	// замена расширения xml на ini
	iniFileName := strings.Replace(xmlFileName, ".xml", ".ini", 1)
	ini.WriteToFile(iniFileName) // запись данных в файл
	return iniFileName
}

/*
 *  Перенесение информации из ini в json
 * @param iniFileName - имя ini-файла
 * @return - имя json-файла или пустую строку при ошибки считывания
 */
func ConvertIniToJson(iniFileName string) string {
	cfg, err := ini2.Load(iniFileName)
	if err != nil {
		Logs("Ошибка открытия файла " + iniFileName)
		return ""
	}

	// список названий секций
	sectionsName := cfg.SectionStrings()
	// ссылка на секции
	section := cfg.Sections()

	// замена расширения ini на json
	jsonFileName := strings.Replace(iniFileName, ".ini", ".json", 1)
	file, err := os.Create(jsonFileName)
	if err != nil {
		Logs("gbctw")
	}
	defer file.Close()
	file.WriteString("[\n")

	for i := 1; i < len(section); i++ {

		//oneStudent.Student = sectionsName[i]
		file.WriteString("{\n\"Student\": " + "\"" + sectionsName[i] + "\",\n")
		keysName := section[i].KeyStrings() // список названий ключей
		keys := section[i].Keys()           // ссылка на ключи

		for j, value := range keys {
			v := value.Value() // получение значения ключ
			file.WriteString("\"Course\":" + "\"" + keysName[j] + "\",\n")
			file.WriteString("\"Mark\":" + "\"" + v + "\",\n")
		}
		file.WriteString("}")

		if i != len(section) {
			file.WriteString(",\n")
		} else {
			file.WriteString("\n")
		}
	}

	file.WriteString("]")
	return jsonFileName

}

/*
 * Запись ошибки в файл
 * @param error - описание ошибки
 */
func Logs(error string) {

	// открыть файл для добавления записи
	file, err := os.OpenFile(fileErrorName, os.O_APPEND, 0666)
	if err != nil { // если файла не существует
		// то, создать файл для записи ошибок
		file, err = os.OpenFile(fileErrorName, os.O_CREATE, 0666)
	}
	defer file.Close() // закрытие файла
	date := time.Now() // текущая дата
	// занести в файл дату и причину ошибку
	file.WriteString(date.Format("2006-01-02 15:04:05") + ": " + error + "\n")
}
