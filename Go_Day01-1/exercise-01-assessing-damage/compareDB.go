package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

type Ingredients struct {
	Ingredient_name  string
	Ingredient_count string
	Ingredient_unit  string
}

type Cake struct {
	Name        string
	Time        string
	Ingredients []Ingredients
}

type JSON struct {
	Cake []Cake
}

type XML struct {
	Cake []struct {
		Name        string `xml:"name"`
		Stovetime   string `xml:"stovetime"`
		Ingredients struct {
			Item []struct {
				Itemname  string `xml:"itemname"`
				Itemcount string `xml:"itemcount"`
				Itemunit  string `xml:"itemunit"`
			} `xml:"item"`
		} `xml:"ingredients"`
	} `xml:"cake"`
}

type data struct {
	xml_file      XML
	json_file     JSON
	type_original string
}

func ParseXml(filename string) XML {
	f_xml, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f_xml.Close()
	xml_data, err := io.ReadAll(f_xml)
	if err != nil {
		log.Fatal(err)
	}
	var xml_result XML
	xmlErr := xml.Unmarshal(xml_data, &xml_result)
	if xmlErr != nil {
		log.Fatal(xmlErr)
	}
	return xml_result
}

func ParseJson(filename string) JSON {
	f_json, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f_json.Close()
	json_data, err := io.ReadAll(f_json)
	if err != nil {
		log.Fatal(err)
	}
	var json_result JSON
	jsonErr := json.Unmarshal(json_data, &json_result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return json_result
}

func CheckFileExtension(str string) {
	switch str {
	case ".xml":
	case "json":
	default:
		fmt.Println("File missing")
		os.Exit(1)
	}
}

func ParsingArguments() (map[string]string, string) {
	count_args := len(os.Args[1:])
	if count_args != 4 {
		fmt.Println("Enter file names with --old or --new flags")
		os.Exit(1)
	}
	files := map[string]string{
		"old": "",
		"new": "",
	}
	for i, e := range os.Args[1:] {
		if i%2 == 0 && (e == "--old" || e == "--new") {
			if files[e[2:]] == "" {
				files[e[2:]] = os.Args[2+i]
			} else {
				fmt.Println("There cannot be two new or old flags")
				os.Exit(1)
			}
		} else if i%2 == 1 {
			continue
		} else {
			fmt.Println("Enter file names with --old or --new flags")
			os.Exit(1)
		}
	}
	len_old := len(files["old"])
	len_new := len(files["new"])
	if len_old < 5 || len_new < 5 {
		fmt.Println("Invalid value entered")
		os.Exit(1)
	}
	if files["old"] == files["new"] {
		fmt.Println("Same file entered")
		os.Exit(1)
	}
	type_original := files["old"][len_old-4:]
	type_new := files["new"][len_new-4:]
	CheckFileExtension(type_original)
	CheckFileExtension(type_new)
	if type_original == type_new {
		fmt.Println("Enter files with different types: xml and json")
		os.Exit(1)
	}
	return files, type_original
}

func AddCake(info data) {
	for _, stolen := range info.json_file.Cake {
		flag := 0
		for _, original := range info.xml_file.Cake {
			if original.Name == stolen.Name {
				flag = 1
				break
			}
		}
		if flag == 0 {
			fmt.Printf("ADDED cake \"%s\"\n", stolen.Name)
		}
	}
}

func RmCake(info data) {
	for _, original := range info.xml_file.Cake {
		flag := 0
		for _, stolen := range info.json_file.Cake {
			if original.Name == stolen.Name {
				flag = 1
				break
			}
		}
		if flag == 0 {
			fmt.Printf("REMOVED cake \"%s\"\n", original.Name)
		}
	}
}

func ChangedCookingTime(info data) {
	for _, original := range info.xml_file.Cake {
		for _, stolen := range info.json_file.Cake {
			if original.Name == stolen.Name && original.Stovetime != stolen.Time {
				fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", original.Name, stolen.Time, original.Stovetime)
			}
		}
	}
}

func AddIngredient(info data) {
	for _, original := range info.xml_file.Cake {
		for _, stolen := range info.json_file.Cake {
			if original.Name == stolen.Name {
				for _, ingredients_stolen := range stolen.Ingredients {
					flag := 0
					for _, ingredients_original := range original.Ingredients.Item {
						if ingredients_original.Itemname == ingredients_stolen.Ingredient_name {
							flag = 1
							break
						}
					}
					if flag == 0 {
						fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", ingredients_stolen.Ingredient_name, original.Name)
					}
				}
			}
		}
	}
}

func RmIngredient(info data) { // лишние дупликаты без проверки на имя
	for _, original := range info.xml_file.Cake {
		for _, ingredients_original := range original.Ingredients.Item {
			flag := 0
			for _, stolen := range info.json_file.Cake {
				for _, ingredients_stolen := range stolen.Ingredients {
					if ingredients_original.Itemname == ingredients_stolen.Ingredient_name {
						flag = 1
						break
					}
				}
				if flag == 0 {
					fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", ingredients_original.Itemname, original.Name)
				}
			}
		}
	}
}

func ChangedUnitCount(info data) {
	for _, stolen := range info.json_file.Cake {
		for _, ingredients_stolen := range stolen.Ingredients {
			for _, original := range info.xml_file.Cake {
				for _, ingredients_original := range original.Ingredients.Item {
					if ingredients_original.Itemname == ingredients_stolen.Ingredient_name && ingredients_original.Itemcount != ingredients_stolen.Ingredient_count {
						fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake"+
							"  \"%s\" - \"%s\" instead of"+
							" \"%s\"\n", ingredients_stolen.Ingredient_name, original.Name, ingredients_stolen.Ingredient_count, ingredients_original.Itemcount)
					}
				}
			}
		}
	}
}

func RmAddChangedUnit(info data) {
	for _, stolen := range info.json_file.Cake {
		for _, ingredients_stolen := range stolen.Ingredients {
			for _, original := range info.xml_file.Cake {
				for _, ingredients_original := range original.Ingredients.Item {
					if ingredients_original.Itemname == ingredients_stolen.Ingredient_name && ingredients_original.Itemunit != ingredients_stolen.Ingredient_unit {
						if ingredients_stolen.Ingredient_unit == "" {
							fmt.Printf("REMOVED unit \"%s\" for ingredient"+
								" \"%s\" for cake  \"%s\"\n", ingredients_original.Itemunit, ingredients_original.Itemname, original.Name)
						} else if ingredients_original.Itemunit == "" {
							fmt.Printf("ADDED unit \"%s\" for ingredient"+
								" \"%s\" for cake  \"%s\"\n", ingredients_stolen.Ingredient_unit, ingredients_original.Itemname, original.Name)
						} else {
							fmt.Printf("CHANGED unit for ingredient \"%s\" for cake "+
								"\"%s\" - \"%s\" instead of "+
								"\"%s\"\n", ingredients_stolen.Ingredient_name, original.Name, ingredients_stolen.Ingredient_unit, ingredients_original.Itemunit)
						}
					}
				}
			}
		}
	}
}

func Calculate(info data) {
	AddCake(info)
	RmCake(info)
	ChangedCookingTime(info)
	AddIngredient(info)
	RmIngredient(info)
	ChangedUnitCount(info)
	RmAddChangedUnit(info)
}

func main() {
	files, type_original := ParsingArguments()
	var data_files data
	data_files.type_original = type_original
	if type_original == ".xml" {
		data_files.xml_file = ParseXml(files["old"])
		data_files.json_file = ParseJson(files["new"])
		Calculate(data_files)
	} else {
		fmt.Println("Original database must be in xml format")
		os.Exit(1)
	}
}
