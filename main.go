package main

import (
	"encoding/json"
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

type Data struct {
	Course string
	Mark   int
}

type jsonFile struct {
	Student string
	data    []Data
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

// Перенесение информации из xml в ini
// xmlFileName - имя xml-файла
// возвращвет - имя ini-файла или пустую строку при ошибки считывания
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

	var allStudent []jsonFile

	for i := 1; i < len(section); i++ {
		oneStudent := jsonFile{}
		//	da := Data{}
		//var dat []Data
		oneStudent.Student = sectionsName[i]
		keysName := section[i].KeyStrings() // список названий ключей
		keys := section[i].Keys()           // ссылка на ключи
		//data := map[string]int{}
		for j, value := range keys {
			v := value.Value()           // получение значения ключа
			chislo, _ := strconv.Atoi(v) // перевод строки в число
			//da.Course = keysName[j]
			//da.Mark = chislo
			//dat = append(dat, da)
			//data[keysName[j]] = chislo
			oneStudent.Course = keysName[j]
			oneStudent.Mark = chislo
		}

		//oneStudent.data = dat
		allStudent = append(allStudent, oneStudent)
	}

	fmt.Print("Данные student")
	fmt.Println(allStudent)

	jsonData, err := json.Marshal(allStudent)

	if err != nil {
		Logs("")
		return ""
	}

	// sanity check - JSON level
	fmt.Print("Данные json: ")
	fmt.Println(string(jsonData))

	// now write to JSON file

	// замена расширения ini на json
	jsonFileName := strings.Replace(iniFileName, ".ini", ".json", 1)

	jsonFile, err := os.Create(jsonFileName)

	if err != nil {
		Logs("")
		return ""
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()

	return jsonFileName
}

// запись ошибки в файл
// @param error - описание ошибки
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
